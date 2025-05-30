package nord

import (
   "fmt"
   "os"
   "testing"
)

func TestServerLoad(t *testing.T) {
   loads := ServerLoads{
      {Count: 2, Country: "CA", City: "toronto", Hostname: "ca1788.nordvpn.com"},
      {Count: 3, Country: "CA", City: "toronto", Hostname: "ca1789.nordvpn.com"},
   }
   for range 9 {
      fmt.Println(loads.Country("CA"))
   }
}

func TestCountry(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/platform/nord/ServerLoads")
   if err != nil {
      t.Fatal(err)
   }
   var loads ServerLoads
   err = loads.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   for _, load := range loads {
      if load.Country == "CA" {
         fmt.Printf("%+v\n\n", load)
      }
   }
}

func TestWrite(t *testing.T) {
   servers, err := GetServers(0)
   if err != nil {
      t.Fatal(err)
   }
   data, err := GetServerLoads(servers).Marshal()
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home+"/platform/nord/ServerLoads", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
