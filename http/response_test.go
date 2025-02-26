package http

import (
   "os"
   "testing"
)

func TestResponse(t *testing.T) {
   resp, err := ReadFile(".mpd")
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
