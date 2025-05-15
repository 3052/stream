package dash

import (
   "fmt"
   "testing"
)

func Test(t *testing.T) {
   fmt.Println(index(representations, 3_300_000, 0.4))
   fmt.Println(index(representations, 3_300_000, 0.1))
}

var representations = []representation{
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
