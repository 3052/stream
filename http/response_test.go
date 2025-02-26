package http

import (
   "fmt"
   "net/http"
   "os"
   "testing"
)

func TestResponse(t *testing.T) {
   resp, err := http.Get("http://httpbingo.org/get")
   if err != nil {
      t.Fatal(err)
   }
   err = WriteFile("http.txt", resp)
   if err != nil {
      t.Fatal(err)
   }
   resp1, err := ReadFile("http.txt")
   if err != nil {
      t.Fatal(err)
   }
   err = resp1.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%#v\n", resp1.Request.URL)
}
