package net

import (
   "41.neocities.org/dash"
   "errors"
   "fmt"
   "io"
   "net/http"
   "slices"
   "strconv"
   "strings"
)

func (f *Filters) Set(data string) error {
   *f = append(*f, data)
   return nil
}

func (f Filters) String() string {
   var b []byte
   for i, filter := range f {
      if i >= 1 {
         b = append(b, ", "...)
      }
      b = strconv.AppendQuote(b, filter)
   }
   return string(b)
}

// wikipedia.org/wiki/Filter_(higher-order_function)
type Filters []string

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
