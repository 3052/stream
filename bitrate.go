package net

import (
   "41.neocities.org/dash"
   "fmt"
   "io"
   "net/http"
   "slices"
)

func (b *Bitrate) String() string {
   var data []byte
   for _, value := range b.Value {
      if data != nil {
         data = append(data, ", "...)
      }
      data = fmt.Appendf(data, "%v-%v", value[0], value[1])
   }
   return string(data)
}

// github.com/golang/go/blob/go1.24.3/src/math/all_test.go#L2146
func (b *Bitrate) contains(actual int) bool {
   for _, correct := range b.Value {
      if actual >= correct[0] {
         if actual <= correct[1] {
            return true
         }
      }
   }
   return false
}

func (b *Bitrate) Set(data string) error {
   var value [2]int
   _, err := fmt.Sscanf(data, "%v-%v", &value[0], &value[1])
   if err != nil {
      return err
   }
   if b.Ok {
      b.Value = append(b.Value, value)
   } else {
      b.Value = [][2]int{value}
      b.Ok = true
   }
   return nil
}

type Bitrate struct {
   Value [][2]int
   Ok    bool
}

func (e *License) Bitrate(resp *http.Response, correct *Bitrate) error {
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
      if correct.contains(represent.Bandwidth) {
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
