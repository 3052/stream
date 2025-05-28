package net

import (
   "io"
   "log"
   "net/http"
   "testing"
)

func TestProgress(t *testing.T) {
   log.SetFlags(log.Ltime)
   var (
      segment   [9]struct{}
      progress1 progress
   )
   progress1.set(len(segment))
   for range segment {
      func() {
         resp, err := http.Get("http://httpbingo.org/drip?delay=0&duration=1")
         if err != nil {
            t.Fatal(err)
         }
         defer resp.Body.Close()
         _, err = io.Copy(io.Discard, resp.Body)
         if err != nil {
            t.Fatal(err)
         }
      }()
      progress1.next()
   }
}
