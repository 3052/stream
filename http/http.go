package http

import (
   "bufio"
   "bytes"
   "io"
   "log"
   "net/http"
   "net/url"
   "os"
   "time"
)

func Write(name string, resp *http.Response) error {
   file, err := os.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   data, err := resp.Request.URL.MarshalBinary()
   if err != nil {
      return err
   }
   _, err = file.Write(data)
   if err != nil {
      return err
   }
   _, err = file.Write([]byte{0})
   if err != nil {
      return err
   }
   return resp.Write(file)
}

func Read(name string) (*http.Response, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   buf := bufio.NewReader(file)
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

// firefox does this
// 29s left - 19.9 of 33.5 MB (540 KB/sec)
func log_progress(d time.Duration) {
   log.Println(d.Truncate(time.Second), "left")
}

func (p *ProgressBytes) durationB() time.Duration {
   return p.durationA() * time.Duration(p.byteB) / time.Duration(p.byteA)
}

type ProgressBytes struct {
   byteA int64
   byteB int64
   read io.Reader
   timeA time.Time
   timeB int64
}

func (p *ProgressBytes) Set(resp *http.Response) {
   p.byteB = resp.ContentLength
   p.read = resp.Body
   p.timeA = time.Now()
   p.timeB = time.Now().Unix()
}

func (p *ProgressBytes) durationA() time.Duration {
   return time.Since(p.timeA)
}

func (p *ProgressParts) durationA() time.Duration {
   return time.Since(p.timeA)
}

func (p *ProgressParts) Set(partB int) {
   p.partB = partB
   p.timeA = time.Now()
   p.timeB = time.Now().Unix()
}

type ProgressParts struct {
   partA int64
   partB int
   timeA time.Time
   timeB int64
}

// keep last two terms separate
func (p *ProgressParts) durationB() time.Duration {
   return p.durationA() * time.Duration(p.partB) / time.Duration(p.partA)
}

func (p *ProgressBytes) Read(data []byte) (int, error) {
   n, err := p.read.Read(data)
   p.byteA += int64(n)
   p.byteB -= int64(n)
   timeB := time.Now().Unix()
   if timeB > p.timeB {
      log_progress(p.durationB())
      p.timeB = timeB
   }
   return n, err
}

func (p *ProgressParts) Next() {
   p.partA++
   p.partB--
   timeB := time.Now().Unix()
   if timeB > p.timeB {
      log_progress(p.durationB())
      p.timeB = timeB
   }
}
