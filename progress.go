package stream

import (
   "log"
   "time"
)

// firefox does this
// 29s left - 19.9 of 33.5 MB (540 KB/sec)

func (p *progress) Next() {
   p.segmentA++
   p.segmentB--
   timeB := time.Now().Unix()
   if timeB > p.timeB {
      log.Println(
         p.segmentB, "segment",
         p.durationB().Truncate(time.Second),
         "left",
      )
      p.timeB = timeB
   }
}

func (p *progress) Set(segmentB int) {
   p.segmentB = segmentB
   p.timeA = time.Now()
   p.timeB = time.Now().Unix()
}

type progress struct {
   segmentA int64
   segmentB int
   timeA    time.Time
   timeB    int64
}

func (p *progress) durationA() time.Duration {
   return time.Since(p.timeA)
}

// keep last two terms separate
func (p *progress) durationB() time.Duration {
   return p.durationA() * time.Duration(p.segmentB) / time.Duration(p.segmentA)
}
