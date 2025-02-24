package http

import (
   "bufio"
   "net/http"
   "net/url"
   "os"
)

func read() (*http.Response, error) {
   data, err := os.ReadFile("Request")
   if err != nil {
      return nil, err
   }
   var req url.URL
   err = req.UnmarshalBinary(data)
   if err != nil {
      return nil, err
   }
   file, err := os.Open("Response")
   if err != nil {
      return nil, err
   }
   defer file.Close()
   return http.ReadResponse(bufio.NewReader(file), &http.Request{URL: &req})
}

func write(req *http.Request) error {
   data, err := req.URL.MarshalBinary()
   if err != nil {
      return err
   }
   err = os.WriteFile("Request", data, os.ModePerm)
   if err != nil {
      return err
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   file, err := os.Create("Response")
   if err != nil {
      return err
   }
   defer file.Close()
   return resp.Write(file)
}
