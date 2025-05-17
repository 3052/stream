package net

import (
   "log"
   "net/http"
)

func init() {
   log.SetFlags(log.Ltime)
   http.DefaultClient.Transport = &transport{
      // x/net/http2: make Transport return nicer error when Amazon ALB hangs up
      // mid-response?
      // github.com/golang/go/issues/18639
      // x/net/http2: Transport ignores net/http.Transport.Proxy once connected
      // github.com/golang/go/issues/25793
      Protocols: &http.Protocols{},
      Proxy:     http.ProxyFromEnvironment,
   }
}

type transport http.Transport

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
   if req.Header.Get("silent") == "" {
      log.Println(req.Method, req.URL)
   }
   return (*http.Transport)(t).RoundTrip(req)
}
