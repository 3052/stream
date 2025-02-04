package http

import (
   "io"
   "net/http"
   "net/url"
   "testing"
)

func TestBytes(t *testing.T) {
   port := Transport{DisableCompression: true}
   port.ProxyFromEnvironment()
   port.DefaultClient()
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
