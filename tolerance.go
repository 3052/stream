package net

import (
   "41.neocities.org/dash"
   "fmt"
   "io"
   "net/http"
)

func (e *License) Tolerance(
   resp *http.Response, limit float64, correct ...int,
) error {
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
   var line bool
   for represent := range mpd.Representation() {
      if line {
         fmt.Println()
      } else {
         line = true
      }
      fmt.Println(represent)
      if tolerance(represent, limit, correct...) {
         switch {
         case represent.SegmentBase != nil:
            err = e.segment_base(represent)
         case represent.SegmentList != nil:
            err = e.segment_list(represent)
         case represent.SegmentTemplate != nil:
            err = e.segment_template(represent)
         }
         if err != nil {
            return err
         }
      }
   }
   return nil
}

// github.com/golang/go/blob/go1.24.3/src/math/all_test.go#L2146
// wikipedia.org/wiki/Engineering_tolerance
func tolerance(actual *dash.Representation, limit float64, correct ...int) bool {
   for _, correct1 := range correct {
      variation := actual.Bandwidth - correct1
      if variation < 0 {
         variation = -variation
      }
      if float64(variation) <= float64(correct1)*limit {
         return true
      }
   }
   return false
}
