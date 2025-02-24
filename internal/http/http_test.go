package http

import (
   "fmt"
   "net/http"
   "os"
   "testing"
)

func Test(t *testing.T) {
   req, err := http.NewRequest("", "http://httpbingo.org/get", nil)
   if err != nil {
      t.Fatal(err)
   }
   err = write(req)
   if err != nil {
      t.Fatal(err)
   }
   resp, err := read()
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%#v\n", resp.Request.URL)
}
