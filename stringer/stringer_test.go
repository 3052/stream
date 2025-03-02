package stringer

import (
   "fmt"
   "testing"
)

var cardinal_tests = []struct{
   in Cardinal
   out string
}{
   {123.45, "123"},
   {123.45*1000, "123.45 thousand"},
   {123.45*1000*1000, "123.45 million"},
   {123.45*1000*1000*1000, "123.45 billion"},
}

var percent_tests = []struct{
   in Percent
   out string
}{
   {0.0123, "1.23 %"},
   {0.1234, "12.34 %"},
}

var rate_tests = []struct{
   in Rate
   out string
}{
   {123.45, "123 byte/s"},
   {123.45*1000, "123.45 kilobyte/s"},
   {123.45*1000*1000, "123.45 megabyte/s"},
   {123.45*1000*1000*1000, "123.45 gigabyte/s"},
}

var size_tests = []struct{
   in Size
   out string
}{
   {123.45, "123 byte"},
   {123.45*1000, "123.45 kilobyte"},
   {123.45*1000*1000, "123.45 megabyte"},
   {123.45*1000*1000*1000, "123.45 gigabyte"},
}

func TestCardinal(t *testing.T) {
   for _, test := range cardinal_tests {
      if fmt.Sprint(test.in) != test.out {
         t.Fatal(test)
      }
   }
}

func TestPercent(t *testing.T) {
   for _, test := range percent_tests {
      if fmt.Sprint(test.in) != test.out {
         t.Fatal(test)
      }
   }
}

func TestRate(t *testing.T) {
   for _, test := range rate_tests {
      if fmt.Sprint(test.in) != test.out {
         t.Fatal(test)
      }
   }
}

func TestSize(t *testing.T) {
   for _, test := range size_tests {
      if fmt.Sprint(test.in) != test.out {
         t.Fatal(test)
      }
   }
}
