package net

import (
   "41.neocities.org/dash"
   "bufio"
   "bytes"
   "fmt"
   "net/http"
   "net/url"
)

func (r response) marshal() ([]byte, error) {
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

type response [1]*http.Response

func (r *response) unmarshal(data []byte) error {
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

func variation(value *dash.Representation, expect int) int {
   variation := value.Bandwidth - expect
   if variation < 0 {
      return -variation
   }
   return variation
}

func expected(values []dash.Representation, expect int) dash.Representation {
   a := values[0]
   for _, b := range values[1:] {
      if variation(&b, expect) < variation(&a, expect) {
         a = b
      }
   }
   return a
}

// github.com/golang/go/blob/go1.24.3/src/math/all_test.go#L2146
func tolerance(value *dash.Representation, expect int, percent float64) bool {
   return float64(variation(value, expect)) <= float64(expect)*percent
}
