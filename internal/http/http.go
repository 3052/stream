package http

import (
   "bufio"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
)

func Write(dir string, resp *http.Response) error {
   data, err := resp.Request.URL.MarshalBinary()
   if err != nil {
      return err
   }
   err = os.WriteFile(request(dir), data, os.ModePerm)
   if err != nil {
      return err
   }
   file, err := os.Create(response(dir))
   if err != nil {
      return err
   }
   defer file.Close()
   return resp.Write(file)
}

func Read(dir string) (*http.Response, error) {
   data, err := os.ReadFile(request(dir))
   if err != nil {
      return nil, err
   }
   var req url.URL
   err = req.UnmarshalBinary(data)
   if err != nil {
      return nil, err
   }
   file, err := os.Open(response(dir))
   if err != nil {
      return nil, err
   }
   defer file.Close()
   return http.ReadResponse(bufio.NewReader(file), &http.Request{URL: &req})
}

func request(dir string) string {
   return filepath.Join(dir, "Request URL")
}

func response(dir string) string {
   return filepath.Join(dir, "Response")
}
