package nord

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func Proxy(username, password, hostname string) string {
	var b strings.Builder
	b.WriteString("https://")
	b.WriteString(username)
	b.WriteByte(':')
	b.WriteString(password)
	b.WriteByte('@')
	b.WriteString(hostname)
	b.WriteString(":89")
	return b.String()
}

type ServerLoad struct {
	Count    int
	Country  string
	City     string
	Hostname string
}

type Server struct {
	Hostname  string
	Locations []struct {
		Country struct {
			City struct {
				DnsName string `json:"dns_name"`
			}
			Code string
		}
	}
	Status       string
	Technologies []struct {
		Identifier string
	}
}

// limit <= -1 for default
// limit == 0 for all
func GetServers(limit int) ([]Server, error) {
	req, _ := http.NewRequest("", "https://api.nordvpn.com/v2/servers", nil)
	if limit >= 0 {
		req.URL.RawQuery = "limit=" + strconv.Itoa(limit)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var servers []Server
	err = json.NewDecoder(resp.Body).Decode(&servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

type ServerLoads []*ServerLoad

func (s ServerLoads) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ServerLoads) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}

func GetServerLoads(servers []Server) ServerLoads {
	loads := make(ServerLoads, 0, len(servers))
	for _, server1 := range servers {
		if server1.proxy_ssl() {
			var load ServerLoad
			load.Hostname = server1.Hostname
			for _, location := range server1.Locations {
				load.Country = location.Country.Code
				load.City = location.Country.City.DnsName
			}
			loads = append(loads, &load)
		}
	}
	return loads
}

func (s ServerLoads) Country(code string) string {
	var load *ServerLoad
	for _, load1 := range s {
		if load1.Country == code {
			if load != nil {
				if load1.Count >= load.Count {
					continue
				}
			}
			load = load1
		}
	}
	load.Count++
	return load.Hostname
}

func (s *Server) proxy_ssl() bool {
	for _, technology := range s.Technologies {
		if technology.Identifier == "proxy_ssl" {
			return true
		}
	}
	return false
}
