package http

import (
   "fmt"
   "net/http"
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
   defer resp.Body.Close()
   var resp1 Response
   err = resp1.New(resp)
   if err != nil {
      t.Fatal(err)
   }
   err = resp1.Write("http")
   if err != nil {
      t.Fatal(err)
   }
   var resp2 Response
   err = resp2.Read("http")
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%s %v\n", resp2.Body, resp2.Url)
}
