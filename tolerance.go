package net

import (
   "41.neocities.org/dash"
   "fmt"
   "io"
   "net/http"
   "os"
   "strings"
)

func create(represent *dash.Representation) (*os.File, error) {
   var name strings.Builder
   name.WriteString(represent.Id)
   switch *represent.MimeType {
   case "audio/mp4":
      name.WriteString(".m4a")
   case "video/mp4":
      name.WriteString(".m4v")
   }
   return os_create(name.String())
}

func (e *License) Tolerance(
   resp *http.Response, correct int, limit float64,
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

// github.com/golang/go/blob/go1.24.3/src/math/all_test.go#L2146
// wikipedia.org/wiki/Engineering_tolerance
func tolerance(actual *dash.Representation, correct int, limit float64) bool {
   variation := actual.Bandwidth - correct
   if variation < 0 {
      variation = -variation
   }
   return float64(variation) <= float64(correct)*limit
}
