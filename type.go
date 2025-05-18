package net

import (
   "41.neocities.org/dash"
   "bufio"
   "bytes"
   "fmt"
   "net/http"
   "net/url"
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

func (r Response) marshal() ([]byte, error) {
   var buf bytes.Buffer
   _, err := fmt.Fprintln(&buf, r[0].Request.URL)
   if err != nil {
      return nil, err
   }
   err = r[0].Write(&buf)
   if err != nil {
      return nil, err
   }
   return buf.Bytes(), nil
}

type Response [1]*http.Response

func (r *Response) unmarshal(data []byte) error {
   before, data, _ := bytes.Cut(data, []byte{'\n'})
   var base url.URL
   err := base.UnmarshalBinary(before)
   if err != nil {
      return err
   }
   r[0], err = http.ReadResponse(
      bufio.NewReader(bytes.NewReader(data)), &http.Request{URL: &base},
   )
   if err != nil {
      return err
   }
   return nil
}
