package http

import (
   "bytes"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "os"
)

type Response struct {
   Url *url.URL
   Body []byte
}

func (r *Response) New(resp *http.Response) error {
   defer resp.Body.Close()
   var err error
   r.Body, err = io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   r.Url = resp.Request.URL
   return nil
}

func (r *Response) Read(name string) error {
   data, err := os.ReadFile(name)
   if err != nil {
      return err
   }
   data, r.Body, _ = bytes.Cut(data, []byte{'\n'})
   r.Url = &url.URL{}
   return r.Url.UnmarshalBinary(data)
}

func (r *Response) Write(name string) error {
   log.Println("Create", name)
   file, err := os.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   _, err = fmt.Fprintln(file, r.Url)
   if err != nil {
      return err
   }
   _, err = file.Write(r.Body)
   if err != nil {
      return err
   }
   return nil
}
