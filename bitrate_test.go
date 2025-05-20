package net

import (
   "41.neocities.org/dash"
   "fmt"
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

func TestBitrate(t *testing.T) {
   correct := Bitrate{
      Value: [][2]int{{2_000_000, 3_000_000}},
   }
   for _, actual := range representations {
      fmt.Println(actual.Bandwidth, correct.contains(actual.Bandwidth))
   }
}
