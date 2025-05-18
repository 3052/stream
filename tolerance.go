package net

import (
   "41.neocities.org/dash"
   "fmt"
   "io"
   "log"
   "net/http"
)

// github.com/golang/go/blob/go1.24.3/src/math/all_test.go#L2146
// wikipedia.org/wiki/Engineering_tolerance
func tolerance(actual *dash.Representation, correct []int64, limit float64) bool {
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

func (e *License) Tolerance(
   resp *http.Response, correct []int64, limit float64,
) error {
   for _, correct1 := range correct {
      variation := float64(correct1)*limit
      log.Println(
         "tolerance", correct1-int64(variation), correct1+int64(variation),
      )
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
   var line bool
   for represent := range mpd.Representation() {
      if line {
         fmt.Println()
      } else {
         line = true
      }
      fmt.Println(represent)
      if tolerance(represent, correct, limit) {
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
