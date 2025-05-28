package net

import (
   "41.neocities.org/dash"
   "errors"
   "flag"
   "fmt"
   "io"
   "net/http"
   "slices"
   "strconv"
   "strings"
)

type bandwidth struct {
   start int
   end   int
}

type filter struct {
   bandwidth bandwidth
   language  string
}

func (f *filter) String() string {
   var b []byte
   if f.bandwidth.start >= 1 {
      b = fmt.Append(b, "bs=", f.bandwidth.start)
   }
   if f.bandwidth.end >= 1 {
      if b != nil {
         b = append(b, ';')
      }
      b = fmt.Append(b, "be=", f.bandwidth.end)
   }
   if f.language != "" {
      if b != nil {
         b = append(b, ';')
      }
      b = fmt.Append(b, "l=", f.language)
   }
   return string(b)
}

func (f *filter) Set(data string) error {
   cookies, err := http.ParseCookie(data)
   if err != nil {
      return err
   }
   for _, cookie := range cookies {
      switch cookie.Name {
      case "bs":
         _, err = fmt.Sscan(cookie.Value, &f.bandwidth.start)
      case "be":
         _, err = fmt.Sscan(cookie.Value, &f.bandwidth.end)
      case "l":
         f.language = cookie.Value
      }
      if err != nil {
         return err
      }
   }
   return nil
}

type filters []filter

func (f filters) String() string {
   var b []byte
   for i, value := range f {
      if i >= 1 {
         b = append(b, ',')
      }
      b = fmt.Append(b, &value)
   }
   return string(b)
}

func (f *filters) Set(data string) error {
   *f = nil
   for _, data := range strings.Split(data, ",") {
      var filter1 filter
      err := filter1.Set(data)
      if err != nil {
         return err
      }
      *f = append(*f, filter1)
   }
   return nil
}

// wikipedia.org/wiki/Filter_(higher-order_function)
const usage = `bs = bandwidth start
be = bandwidth end
l = language
`

func (f Filters) match(represent *dash.Representation) bool {
   return false
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
