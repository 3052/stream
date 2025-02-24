package http

import (
   "fmt"
   "net/http"
   "os"
   "testing"
)

func Test(t *testing.T) {
   resp, err := http.Get("http://httpbingo.org/get")
   if err != nil {
      t.Fatal(err)
   }
   err = Write("http", resp)
   if err != nil {
      t.Fatal(err)
   }
   resp1, err := Read("http")
   if err != nil {
      t.Fatal(err)
   }
   err = resp1.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%#v\n", resp1.Request.URL)
}
