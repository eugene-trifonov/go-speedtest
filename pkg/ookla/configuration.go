package ookla

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"time"

	geo "github.com/kellydunn/golang-geo"
)

// Configuration represents the xml configuration structure from ookla's server
type Configuration struct {
	Client       Client       `xml:"client"`
	ServerConfig ServerConfig `xml:"server-config"`
	Times        Times        `xml:"times"`
	Download     Download     `xml:"socket-download"`
	Upload       Upload       `xml:"socket-upload"`
	Latency      Latency      `xml:"socket-latency"`
}

// Client represents the xml configuration for client
type Client struct {
	IP        string  `xml:"ip,attr"`
	ISP       string  `xml:"isp,attr"`
	Latitude  float64 `xml:"lat,attr"`
	Longitude float64 `xml:"lon,attr"`
}

// ServerConfig represents the xml configuration for server config
type ServerConfig struct {
	IgnoreIDs   string `xml:"ignoreids,attr"`
	ThreadCount string `xml:"threadcount,attr"`
}

// Times represents the xml configuration for how many times transfer can be tested
type Times struct {
	DownloadOne   int `xml:"dl1,attr"`
	DownloadTwo   int `xml:"dl2,attr"`
	DownloadThree int `xml:"dl3,attr"`
	UploadOne     int `xml:"ul1,attr"`
	UploadTwo     int `xml:"ul2,attr"`
	UploadThree   int `xml:"ul3,attr"`
}

// Download represents the xml configuration for a download timeout
type Download struct {
	Timeout      float64 `xml:"testlength,attr"`
	PacketLength int     `xml:"packetlength,attr"`
}

// Upload represents the xml configuration for a upload timeout
type Upload struct {
	Timeout      float64 `xml:"testlength,attr"`
	PacketLength int     `xml:"packetlength,attr"`
}

// Latency represents the xml configurtion of latency
type Latency struct {
	Length float64 `xml:"testlength,attr"`
}

// Servers is the root element in xml configuration
type Servers struct {
	Servers []Server `xml:"servers>server"`
}

// Server represents the configuration for a server
type Server struct {
	CC        string        `xml:"cc,attr" json:"cc"`
	Country   string        `xml:"country,attr" json:"country"`
	ID        int           `xml:"id,attr" json:"id"`
	Latitude  float64       `xml:"lat,attr" json:"lat"`
	Longitude float64       `xml:"lon,attr" json:"lon"`
	Name      string        `xml:"name,attr" json:"name"`
	Sponsor   string        `xml:"sponsor,attr" json:"sponsor"`
	URL       string        `xml:"url,attr" json:"url"`
	URL2      string        `xml:"url2,attr" json:"url2"`
	Host      string        `xml:"host,attr" json:"host"`
	Distance  float64       `xml:"distance,attr" json:"distance"`
	Latency   time.Duration `xml:"latency,attr" json:"latency"`
	tcpAddr   *net.TCPAddr
}

func (s Server) dial(ctx context.Context) (net.Conn, error) {
	var dialer net.Dialer
	return dialer.DialContext(ctx, "tcp", s.tcpAddr.String())
}

// GetConfiguration fetchs ookla's Configuration
func GetConfiguration() (config Configuration, err error) {
	res, err := http.Get("https://www.speedtest.net/speedtest-config.php")
	if err != nil {
		return config, fmt.Errorf("failed to load configuration: %w", err)
	}
	defer res.Body.Close()

	settingsBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return config, err
	}

	err = xml.Unmarshal(settingsBody, &config)
	return config, nil
}

// GetServers gets near located servers for a client
func GetServers(config Configuration) ([]Server, error) {
	res, err := http.Get("https://www.speedtest.net/speedtest-servers.php")
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Speedtest.net servers: %w", err)
	}
	defer res.Body.Close()
	serversBody, _ := ioutil.ReadAll(res.Body)
	var allServers Servers
	xml.Unmarshal(serversBody, &allServers)

	if len(allServers.Servers) == 0 {
		return nil, errors.New("no servers found")
	}

	servers := allServers.Servers

	for i := range servers {
		addr, err := net.ResolveTCPAddr("tcp", servers[i].Host)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve address %q: %w", servers[i].Host, err)
		}
		servers[i].tcpAddr = addr

		servers[i].setDistances(config.Client.Latitude, config.Client.Longitude)
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Distance > servers[j].Distance
	})

	return servers[:2], nil
}

func (s *Server) setDistances(latitude, longitude float64) {
	me := geo.NewPoint(latitude, longitude)
	serverPoint := geo.NewPoint(s.Latitude, s.Longitude)
	distance := me.GreatCircleDistance(serverPoint)
	s.Distance = distance
}
