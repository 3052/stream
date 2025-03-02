package progress

import (
   "io"
   "log"
   "net/http"
   "time"
)

// firefox does this
// 29s left - 19.9 of 33.5 MB (540 KB/sec)
func log_progress(d time.Duration) {
   log.Println(d.Truncate(time.Second), "left")
}

func (b *Bytes) durationB() time.Duration {
   return b.durationA() * time.Duration(b.byteB) / time.Duration(b.byteA)
}

type Bytes struct {
   byteA int64
   byteB int64
   read io.Reader
   timeA time.Time
   timeB int64
}

func (b *Bytes) Set(resp *http.Response) {
   b.byteB = resp.ContentLength
   b.read = resp.Body
   b.timeA = time.Now()
   b.timeB = time.Now().Unix()
}

func (b *Bytes) durationA() time.Duration {
   return time.Since(b.timeA)
}

func (p *Parts) durationA() time.Duration {
   return time.Since(p.timeA)
}

func (p *Parts) Set(partB int) {
   p.partB = partB
   p.timeA = time.Now()
   p.timeB = time.Now().Unix()
}

type Parts struct {
   partA int64
   partB int
   timeA time.Time
   timeB int64
}

// keep last two terms separate
func (p *Parts) durationB() time.Duration {
   return p.durationA() * time.Duration(p.partB) / time.Duration(p.partA)
}

func (b *Bytes) Read(data []byte) (int, error) {
   n, err := b.read.Read(data)
   b.byteA += int64(n)
   b.byteB -= int64(n)
   timeB := time.Now().Unix()
   if timeB > b.timeB {
      log_progress(b.durationB())
      b.timeB = timeB
   }
   return n, err
}

func (p *Parts) Next() {
   p.partA++
   p.partB--
   timeB := time.Now().Unix()
   if timeB > p.timeB {
      log_progress(p.durationB())
      p.timeB = timeB
   }
}
