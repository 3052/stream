package progress

import (
   "io"
   "log"
   "net/http"
   "net/url"
   "testing"
)

func TestSegment(t *testing.T) {
   log.SetFlags(log.Ltime)
   var (
      segment    [9]struct{}
      progress Segment
   )
   progress.Set(len(segment))
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
      progress.Next()
   }
}

func TestByte(t *testing.T) {
   http.DefaultClient.Transport = &http.Transport{DisableCompression: true}
   req := http.Request{URL: &url.URL{
      Scheme:   "http",
      Host:     "httpbingo.org",
      Path:     "/drip",
      RawQuery: "delay=0&duration=9",
   }}
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      t.Fatal(err)
   }
   defer resp.Body.Close()
   var progress Byte
   progress.Set(resp)
   _, err = io.Copy(io.Discard, &progress)
   if err != nil {
      t.Fatal(err)
   }
}
