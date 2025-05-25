package net

import (
   "errors"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

var Client = http.Client{
   Transport: &http.Transport{},
}

func get_segment(u *url.URL, head http.Header) ([]byte, error) {
   req := http.Request{Method: "GET", URL: u}
   if head != nil {
      req.Header = head
   } else {
      req.Header = http.Header{}
   }
   resp, err := Client.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   switch resp.StatusCode {
   case http.StatusOK, http.StatusPartialContent:
   default:
      var data strings.Builder
      resp.Write(&data)
      return nil, errors.New(data.String())
   }
   return io.ReadAll(resp.Body)
}

func init() {
   log.SetFlags(log.Ltime)
   http.DefaultClient.Transport = Proxy(true)
}

type Proxy bool

func (p Proxy) RoundTrip(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   if p {
      return http.DefaultTransport.RoundTrip(req)
   }
   return new(http.Transport).RoundTrip(req)
}
