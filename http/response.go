package http

import (
   "bufio"
   "bytes"
   "net/http"
   "net/url"
)

func Marshal(resp *http.Response) ([]byte, error) {
   data, err := resp.Request.URL.MarshalBinary()
   if err != nil {
      return nil, err
   }
   buf := bytes.NewBuffer(append(data, 0))
   err = resp.Write(buf)
   if err != nil {
      return nil, err
   }
   return buf.Bytes(), nil
}

func Unmarshal(data []byte) (*http.Response, error) {
   buf := bufio.NewReader(bytes.NewReader(data))
   data, err := buf.ReadSlice(0)
   if err != nil {
      return nil, err
   }
   var u url.URL
   err = u.UnmarshalBinary(bytes.TrimSuffix(data, []byte{0}))
   if err != nil {
      return nil, err
   }
   return http.ReadResponse(buf, &http.Request{URL: &u})
}
