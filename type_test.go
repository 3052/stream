package net

import (
   "41.neocities.org/dash"
   "fmt"
   "io"
   "log"
   "net/http"
   "testing"
)

func TestExpect(t *testing.T) {
   const expect = 3_300_000
   actual := expected(representations, expect)
   fmt.Printf("%+v %v\n", actual, tolerance(&actual, expect, 0.4))
   fmt.Printf("%+v %v\n", actual, tolerance(&actual, expect, 0.1))
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
