package net

import (
   "41.neocities.org/dash"
   "fmt"
   "testing"
)

func TestTolerance(t *testing.T) {
   for _, actual := range representations {
      fmt.Println(actual.Bandwidth, tolerance(&actual, 3_300_000, 0.4))
   }
}

var representations = []dash.Representation{
   { Bandwidth: 5_096_445 },
   { Bandwidth: 2_748_690 },
   { Bandwidth: 1_867_586 },
   { Bandwidth: 1278765 },
   { Bandwidth: 772927 },
   { Bandwidth: 402389 },
   { Bandwidth: 102803 },
   { Bandwidth: 1216 },
}
