package main

import (
   "41.neocities.org/net/nord"
   "flag"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "os"
   "os/exec"
   "path/filepath"
   "strings"
   "time"
)

func main() {
   log.SetFlags(log.Ltime)
   http.DefaultClient.Transport = transport{}
   write := flag.Bool("w", false, "write")
   country_code := flag.String("c", "", "country code")
   flag.Parse()
   name, err := os.UserHomeDir()
   if err != nil {
      panic(err)
   }
   name = filepath.ToSlash(name) + "/net/nord/ServerLoads"
   switch {
   case *country_code != "":
      err := do_country(name, *country_code)
      if err != nil {
         panic(err)
      }
   case *write:
      err := do_write(name)
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}

func (transport) RoundTrip(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultTransport.RoundTrip(req)
}

type transport struct{}

func do_write(name string) error {
   servers, err := nord.GetServers(0)
   if err != nil {
      return err
   }
   data, err := nord.GetServerLoads(servers).Marshal()
   if err != nil {
      return err
   }
   return write_file(name, data)
}
func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func command(name string, arg ...string) ([]byte, error) {
   c := exec.Command(name, arg...)
   log.Println("Output", c.Args)
   return c.Output()
}

func do_country(name, code string) error {
   data, err := command("password", "-i", "nordvpn.com")
   if err != nil {
      return err
   }
   username, password, _ := strings.Cut(string(data), ":")
   data, err = read_file(name)
   if err != nil {
      return err
   }
   var loads nord.ServerLoads
   err = loads.Unmarshal(data)
   if err != nil {
      return err
   }
   proxy := nord.Proxy(username, password, loads.Country(code))
   fmt.Printf("$env:http_proxy = %q\n", proxy)
   fmt.Printf("$env:https_proxy = %q\n", proxy)
   data, err = loads.Marshal()
   if err != nil {
      return err
   }
   return write_file(name, data)
}

const month = 30 * 24 *time.Hour

func read_file(name string) ([]byte, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   info, err := file.Stat()
   if err != nil {
      return nil, err
   }
   if time.Since(info.ModTime()) >= month {
      return nil, errors.New("ModTime")
   }
   return io.ReadAll(file)
}
