package mullvad

import (
   "errors"
   "log"
   "net/http"
   "os/exec"
   "strings"
   "time"
)

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
   var err error
   if req.Header.Get("proxy") != "" {
      if !connected {
         err = Connect()
         connected = true
      }
   } else if connected {
      err = Disconnect()
      connected = false
   }
   if err != nil {
      return nil, err
   }
   if req.Header.Get("silent") == "" {
      log.Println(req.Method, req.URL)
   }
   return (*http.Transport)(t).RoundTrip(req)
}

func command(name string, arg ...string) *exec.Cmd {
   cmd := exec.Command(name, arg...)
   log.Print(cmd.Args)
   return cmd
}

type Transport http.Transport

var connected bool

func status(prefix string) error {
   after := time.After(30 * time.Second)
   for {
      select {
      case <-time.After(time.Second):
         data, err := command("mullvad", "status").Output()
         if err != nil {
            return err
         }
         data1 := string(data)
         if strings.HasPrefix(data1, prefix) {
            if strings.Contains(data1, " IPv4:") {
               return nil
            }
         }
      case <-after:
         return errors.New("timeout")
      }
   }
}

func Connect() error {
   err := command("mullvad", "connect").Run()
   if err != nil {
      return err
   }
   return status("Connected")
}

func Disconnect() error {
   err := command("mullvad", "disconnect").Run()
   if err != nil {
      return err
   }
   return status("Disconnected")
}
