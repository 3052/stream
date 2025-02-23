package http

import (
   "io"
   "log"
   "net/http"
   "time"
)

type Transport struct{}

func (Transport) RoundTrip(req *http.Request) (*http.Response, error) {
   if req.Method == "" {
      req.Method = "GET"
   }
   log.Println(req.Method, req.URL)
   return http.DefaultTransport.RoundTrip(req)
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
