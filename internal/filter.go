package main

import (
   "flag"
   "fmt"
   "net/http"
   "strings"
)

type filters []filter

func (f *filter) String() string {
   var b []byte
   if f.bandwidth.start >= 1 {
      b = fmt.Append(b, "bs=", f.bandwidth.start)
   }
   if f.bandwidth.end >= 1 {
      if b != nil {
         b = append(b, ';')
      }
      b = fmt.Append(b, "be=", f.bandwidth.end)
   }
   if f.language != "" {
      if b != nil {
         b = append(b, ';')
      }
      b = fmt.Append(b, "l=", f.language)
   }
   return string(b)
}

func (f *filter) Set(data string) error {
   cookies, err := http.ParseCookie(data)
   if err != nil {
      return err
   }
   for _, cookie := range cookies {
      switch cookie.Name {
      case "bs":
         _, err = fmt.Sscan(cookie.Value, &f.bandwidth.start)
      case "be":
         _, err = fmt.Sscan(cookie.Value, &f.bandwidth.end)
      case "l":
         f.language = cookie.Value
      }
      if err != nil {
         return err
      }
   }
   return nil
}

type bandwidth struct {
   start int
   end   int
}

type filter struct {
   bandwidth bandwidth
   language  string
}

func (f filters) String() string {
   var b []byte
   for i, value := range f {
      if i >= 1 {
         b = append(b, ',')
      }
      b = fmt.Append(b, &value)
   }
   return string(b)
}

func (f *filters) Set(data string) error {
   *f = nil
   for _, data := range strings.Split(data, ",") {
      var filter1 filter
      err := filter1.Set(data)
      if err != nil {
         return err
      }
      *f = append(*f, filter1)
   }
   return nil
}

const usage = `bs = bandwidth start
be = bandwidth end
l = language
`

func main() {
   filters1 := filters{
      {bandwidth: bandwidth{100_000, 200_000}, language: "english"},
      {bandwidth: bandwidth{3_000_000, 4_000_000}},
   }
   flag.Var(&filters1, "f", usage)
   flag.Usage()
   flag.Parse()
   fmt.Printf("%#v\n", filters1)
}
