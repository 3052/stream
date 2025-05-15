package main

import (
   "fmt"
   "slices"
)

var representations = []*representation{
   {
      bandwidth: 5_096_445,
      codecs:    "avc1.640028",
   },
   {
      bandwidth: 2_748_690,
      codecs:    "avc1.64001f",
   },
   {
      bandwidth: 1_867_586,
      codecs:    "avc1.64001f",
   },
   {
      bandwidth: 1278765,
      codecs:    "avc1.64001f",
   },
   {
      bandwidth: 772927,
      codecs:    "avc1.64001f",
   },
   {
      bandwidth: 402389,
      codecs:    "avc1.64001f",
   },
   {
      bandwidth: 102803,
      codecs:    "mp4a.40.2",
   },
   {
      bandwidth: 1216,
      codecs:    "wvtt",
   },
}

type representation struct {
   bandwidth int
   codecs    string
}

func (r *representation) variation(expect int) int {
   variation := r.bandwidth - expect
   if variation < 0 {
      return -variation
   }
   return variation
}

// github.com/golang/go/blob/go1.24.3/src/math/all_test.go#L2146
func (r *representation) tolerance(expect int, percent float64) bool {
   return float64(r.variation(expect)) <= float64(expect)*percent
}

func (a *representation) compare(b *representation, expect int) int {
   return a.variation(expect) - b.variation(expect)
}

func main() {
   const expect = 3_300_000
   rep := slices.MinFunc(representations, func(a, b *representation) int {
      return a.compare(b, expect)
   })
   fmt.Printf("%+v %v\n", rep, rep.tolerance(expect, 0.4))
   fmt.Printf("%+v %v\n", rep, rep.tolerance(expect, 0.1))
}
