package net

import (
   "41.neocities.org/dash"
   "fmt"
   "log"
   "testing"
)

var representations = []dash.Representation{
   {Bandwidth: 5_096_445},
   {Bandwidth: 2_748_690},
   {Bandwidth: 1_867_586},
   {Bandwidth: 1278765},
   {Bandwidth: 772927},
   {Bandwidth: 402389},
   {Bandwidth: 102803},
   {Bandwidth: 1216},
}

func TestTolerance(t *testing.T) {
   correct := []int64{3_000_000, 100_000}
   for _, correct1 := range correct {
      variation := float64(correct1) * 0.1
      log.Println(
         "tolerance", correct1-int64(variation), correct1+int64(variation),
      )
   }
   for _, actual := range representations {
      fmt.Println(
         actual.Bandwidth, tolerance(&actual, correct, 0.3),
      )
   }
}
