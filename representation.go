package net

import (
   "41.neocities.org/dash"
   "fmt"
   "io"
   "os"
   "slices"
)

type representation [1]dash.Representation

func (a *representation) expect(actual []representation, expect int) {
   *a = actual[0]
   for _, b := range actual[1:] {
      if b.variation(expect) < a.variation(expect) {
         *a = b
      }
   }
}

func (r *representation) variation(expect int) int {
   variation := r[0].Bandwidth - expect
   if variation < 0 {
      return -variation
   }
   return variation
}

// github.com/golang/go/blob/go1.24.3/src/math/all_test.go#L2146
func (r *representation) tolerance(expect int, percent float64) bool {
   return float64(r.variation(expect)) <= float64(expect)*percent
}

func (r Response) Mpd(name string) error {
   data, err := r.marshal()
   if err != nil {
      return err
   }
   err = write_file(name, data)
   if err != nil {
      return err
   }
   err = r.unmarshal(data)
   if err != nil {
      return err
   }
   defer r[0].Body.Close()
   data, err = io.ReadAll(r[0].Body)
   if err != nil {
      return err
   }
   var mpd dash.Mpd
   err = mpd.Unmarshal(data)
   if err != nil {
      return err
   }
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
   }
   return nil
}

func (e *License) Download(name, id string) error {
   data, err := os.ReadFile(name)
   if err != nil {
      return err
   }
   var resp Response
   err = resp.unmarshal(data)
   if err != nil {
      return err
   }
   defer resp[0].Body.Close()
   data, err = io.ReadAll(resp[0].Body)
   if err != nil {
      return err
   }
   var mpd dash.Mpd
   err = mpd.Unmarshal(data)
   if err != nil {
      return err
   }
   mpd.Set(resp[0].Request.URL)
   for represent := range mpd.Representation() {
      if represent.Id == id {
         if represent.SegmentBase != nil {
            return e.segment_base(represent)
         }
         if represent.SegmentList != nil {
            return e.segment_list(represent)
         }
         return e.segment_template(represent)
      }
   }
   return nil
}
