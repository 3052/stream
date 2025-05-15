package dash

type representation struct {
   bandwidth int64
   codecs    string
}

// github.com/golang/go/blob/go1.24.3/src/math/all_test.go#L2146
func (r *representation) variation(expected int64) int64 {
   variation := r.bandwidth - expected
   if variation < 0 {
      return -variation
   }
   return variation
}

func index(values []representation, expected int64, percent float64) int {
   i := -1
   for j, value := range values {
      variation := value.variation(expected)
      if float64(variation) <= float64(expected)*percent {
         if i == -1 {
            i = j
         } else if variation < values[i].variation(expected) {
            i = j
         }
      }
   }
   return i
}
