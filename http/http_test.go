package http

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "testing"
)

func TestWriteRead(t *testing.T) {
   resp, err := http.Get("http://httpbingo.org/get")
   if err != nil {
      t.Fatal(err)
   }
   err = Create("http.txt", resp)
   if err != nil {
      t.Fatal(err)
   }
   resp1, err := Open("http.txt")
   if err != nil {
      t.Fatal(err)
   }
   err = resp1.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%#v\n", resp1.Request.URL)
}

func TestParts(t *testing.T) {
   var (
      parts [9]struct{}
      progress ProgressParts
   )
   progress.Set(len(parts))
   for range parts {
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

func TestBytes(t *testing.T) {
   http.DefaultTransport = &http.Transport{DisableCompression: true}
   req := http.Request{URL: &url.URL{
      Scheme: "http",
      Host: "httpbingo.org",
      Path: "/drip",
      RawQuery: "delay=0&duration=9",
   }}
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {   
      t.Fatal(err)
   }
   defer resp.Body.Close()
   var progress ProgressBytes
   progress.Set(resp)
   _, err = io.Copy(io.Discard, &progress)
   if err != nil {   
      t.Fatal(err)
   }
}
