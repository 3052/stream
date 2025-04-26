package progress

import (
   "io"
   "log"
   "net/http"
   "time"
)

func (b *Byte) Read(data []byte) (int, error) {
   n, err := b.read.Read(data)
   b.byteA += int64(n)
   b.byteB -= int64(n)
   timeB := time.Now().Unix()
   if timeB > b.timeB {
      log.Println(b.durationB().Truncate(time.Second), "left")
      b.timeB = timeB
   }
   return n, err
}

func (b *Byte) durationB() time.Duration {
   return b.durationA() * time.Duration(b.byteB) / time.Duration(b.byteA)
}

type Byte struct {
   byteA int64
   byteB int64
   read  io.Reader
   timeA time.Time
   timeB int64
}

func (b *Byte) Set(resp *http.Response) {
   b.byteB = resp.ContentLength
   b.read = resp.Body
   b.timeA = time.Now()
   b.timeB = time.Now().Unix()
}

func (b *Byte) durationA() time.Duration {
   return time.Since(b.timeA)
}

// firefox does this
// 29s left - 19.9 of 33.5 MB (540 KB/sec)

func (s *Segment) Next() {
   s.segmentA++
   s.segmentB--
   timeB := time.Now().Unix()
   if timeB > s.timeB {
      log.Println(
         s.segmentB, "segment",
         s.durationB().Truncate(time.Second),
         "left",
      )
      s.timeB = timeB
   }
}

func (s *Segment) Set(segmentB int) {
   s.segmentB = segmentB
   s.timeA = time.Now()
   s.timeB = time.Now().Unix()
}

type Segment struct {
   segmentA int64
   segmentB int
   timeA    time.Time
   timeB    int64
}

///

func (s *Segment) durationA() time.Duration {
   return time.Since(s.timeA)
}

// keep last two terms separate
func (s *Segment) durationB() time.Duration {
   return s.durationA() * time.Duration(s.segmentB) / time.Duration(s.segmentA)
}
