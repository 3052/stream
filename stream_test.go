package stream

import (
   "41.neocities.org/dash"
   "fmt"
   "io"
   "log"
   "net/http"
   "testing"
)

func TestExpect(t *testing.T) {
   const expect = 3_300_000
   actual := expected(representations, expect)
   fmt.Printf("%+v %v\n", actual, tolerance(&actual, expect, 0.4))
   fmt.Printf("%+v %v\n", actual, tolerance(&actual, expect, 0.1))
}

func TestProgress(t *testing.T) {
   log.SetFlags(log.Ltime)
   var (
      segment    [9]struct{}
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

var representations = []dash.Representation{
   { Bandwidth: 5_096_445 },
   { Bandwidth: 2_748_690 },
   { Bandwidth: 1_867_586 },
   { Bandwidth: 1278765 },
   { Bandwidth: 772927 },
   { Bandwidth: 402389 },
   { Bandwidth: 102803 },
   { Bandwidth: 1216 },
}
