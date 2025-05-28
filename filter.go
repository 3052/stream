package net

import (
   "41.neocities.org/dash"
   "errors"
   "fmt"
   "io"
   "net/http"
   "slices"
   "strings"
)

type Filter struct {
   BitrateStart int
   BitrateEnd   int
   Language     string
}

const FilterUsage = `bs = bitrate start
be = bitrate end
l = language
`

func (f *Filter) String() string {
   var b []byte
   if f.BitrateStart >= 1 {
      b = fmt.Append(b, "bs=", f.BitrateStart)
   }
   if f.BitrateEnd >= 1 {
      if b != nil {
         b = append(b, ';')
      }
      b = fmt.Append(b, "be=", f.BitrateEnd)
   }
   if f.Language != "" {
      if b != nil {
         b = append(b, ';')
      }
      b = fmt.Append(b, "l=", f.Language)
   }
   return string(b)
}

func (f *Filter) Set(data string) error {
   cookies, err := http.ParseCookie(data)
   if err != nil {
      return err
   }
   for _, cookie := range cookies {
      switch cookie.Name {
      case "bs":
         _, err = fmt.Sscan(cookie.Value, &f.BitrateStart)
      case "be":
         _, err = fmt.Sscan(cookie.Value, &f.BitrateEnd)
      case "l":
         f.Language = cookie.Value
      }
      if err != nil {
         return err
      }
   }
   return nil
}

type Filters []Filter

func (f Filters) String() string {
   var b []byte
   for i, value := range f {
      if i >= 1 {
         b = append(b, ',')
      }
      b = fmt.Append(b, &value)
   }
   return string(b)
}

func (f *Filters) Set(data string) error {
   *f = nil
   for _, data := range strings.Split(data, ",") {
      var filter1 Filter
      err := filter1.Set(data)
      if err != nil {
         return err
      }
      *f = append(*f, filter1)
   }
   return nil
}

func (f Filters) Filter(resp *http.Response, module *Cdm) error {
   if resp.StatusCode != http.StatusOK {
      var data strings.Builder
      resp.Write(&data)
      return errors.New(data.String())
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   var mpd dash.Mpd
   err = mpd.Unmarshal(data)
   if err != nil {
      return err
   }
   mpd.Set(resp.Request.URL)
   represents := slices.SortedFunc(mpd.Representation(),
      func(a, b *dash.Representation) int {
         return a.Bandwidth - b.Bandwidth
      },
   )
   for i, represent := range represents {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(represent)
      if f.match(represent) {
         switch {
         case represent.SegmentBase != nil:
            err = module.segment_base(represent)
         case represent.SegmentList != nil:
            err = module.segment_list(represent)
         case represent.SegmentTemplate != nil:
            err = module.segment_template(represent)
         }
         if err != nil {
            return err
         }
      }
   }
   return nil
}

func (f Filters) match(r *dash.Representation) bool {
   for _, filter1 := range f {
      if r.Bandwidth >= filter1.BitrateStart {
         if r.Bandwidth <= filter1.BitrateEnd {
            if filter1.Language == "" {
               return true
            }
            if r.GetAdaptationSet().Lang == filter1.Language {
               return true
            }
         }
      }
   }
   return false
}
