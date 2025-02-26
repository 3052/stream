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

func Read(r io.Reader) (*http.Response, error) {
   buf := bufio.NewReader(r)
   data, err := buf.ReadSlice(0)
   if err != nil {
      return nil, err
   }
   data = bytes.TrimSuffix(data, []byte{0})
   var u url.URL
   err = u.UnmarshalBinary(data)
   if err != nil {
      return nil, err
   }
   return http.ReadResponse(buf, &http.Request{URL: &u})
}

// we cannot use os.Open because the full file is not read until the response
// body is read
func ReadFile(name string) (*http.Response, error) {
   data, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   return Read(bytes.NewReader(data))
}

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
