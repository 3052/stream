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

func (a *representation) expect(values []representation, expect int) {
   *a = values[0]
   for _, b := range values[1:] {
      if b.variation(expect) < a.variation(expect) {
         *a = b
      }
   }
}

// github.com/golang/go/blob/go1.24.3/src/math/all_test.go#L2146
func (r *representation) tolerance(expect int, percent float64) bool {
   return float64(r.variation(expect)) <= float64(expect)*percent
}
