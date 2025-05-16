package dash

import (
   "fmt"
   "io"
   "log"
   "net/http"
   "testing"
)

func TestProgress(t *testing.T) {
   log.SetFlags(log.Ltime)
   var (
      segment    [9]struct{}
      progress1 progress
   )
   progress1.Set(len(segment))
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
      progress1.Next()
   }
}

func TestExpect(t *testing.T) {
   const expect = 3_300_000
   var actual representation
   actual.expect(representations, expect)
   fmt.Printf("%+v %v\n", actual, actual.tolerance(expect, 0.4))
   fmt.Printf("%+v %v\n", actual, actual.tolerance(expect, 0.1))
}

var representations = []representation{
   {
      bandwidth: 5_096_445,
      codecs:    "avc1.640028",
   },
   {
      bandwidth: 2_748_690,
      codecs:    "avc1.64001f",
   },
   {
      bandwidth: 1_867_586,
      codecs:    "avc1.64001f",
   },
   {
      bandwidth: 1278765,
      codecs:    "avc1.64001f",
   },
   {
      bandwidth: 772927,
      codecs:    "avc1.64001f",
   },
   {
      bandwidth: 402389,
      codecs:    "avc1.64001f",
   },
   {
      bandwidth: 102803,
      codecs:    "mp4a.40.2",
   },
   {
      bandwidth: 1216,
      codecs:    "wvtt",
   },
}
