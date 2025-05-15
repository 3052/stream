package dash

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

func value(values []representation, expect int) representation {
   value1 := values[0]
   for _, value2 := range values[1:] {
      if value2.variation(expect) < value1.variation(expect) {
         value1 = value2
      }
   }
   return value1
}
