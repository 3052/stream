package http

import (
   "bufio"
   "bytes"
   "io"
   "log"
   "net/http"
   "net/url"
   "os"
)

func Write(w io.Writer, resp *http.Response) error {
   data, err := resp.Request.URL.MarshalBinary()
   if err != nil {
      return err
   }
   _, err = w.Write(data)
   if err != nil {
      return err
   }
   _, err = w.Write([]byte{0})
   if err != nil {
      return err
   }
   return resp.Write(w)
}

func WriteFile(name string, resp *http.Response) error {
   log.Println("Create", name)
   file, err := os.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   return Write(file, resp)
}

func Read(r *bufio.Reader) (*http.Response, error) {
   data, err := r.ReadSlice(0)
   if err != nil {
      return nil, err
   }
   data = bytes.TrimSuffix(data, []byte{0})
   var u url.URL
   err = u.UnmarshalBinary(data)
   if err != nil {
      return nil, err
   }
   return http.ReadResponse(r, &http.Request{URL: &u})
}

func ReadFile(name string) (*http.Response, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   return Read(bufio.NewReader(file))
}
