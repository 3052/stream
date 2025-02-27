package http

import (
   "net/http"
   "os"
   "strings"
   "testing"
)

func TestResponse(t *testing.T) {
   resp, err := http.Post(
      "http://httpbingo.org/post", "text/plain",
      strings.NewReader("hello world"),
   )
   if err != nil {
      t.Fatal(err)
   }
   data, err := Marshal(resp)
   if err != nil {
      t.Fatal(err)
   }
   resp1, err := Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   err = resp1.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
