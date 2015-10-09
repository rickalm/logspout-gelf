package gelf

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"time"

	"github.com/gliderlabs/logspout/router"
)

var hostname string

func init() {
	hostname, _ = os.Hostname()
	router.AdapterFactories.Register(NewGelfAdapter, "gelf")
}

// GelfAdapter is an adapter that streams UDP JSON to Graylog
type GelfAdapter struct {
	conn  net.Conn
	route *router.Route
}

// NewGelfAdapter creates a GelfAdapter with UDP as the default transport.
func NewGelfAdapter(route *router.Route) (router.LogAdapter, error) {
	transport, found := router.AdapterTransports.Lookup(route.AdapterTransport("udp"))
	if !found {
		return nil, errors.New("unable to find adapter: " + route.Adapter)
	}

	conn, err := transport.Dial(route.Address, route.Options)
	if err != nil {
		return nil, err
	}

	return &GelfAdapter{
		route: route,
		conn:  conn,
	}, nil
}

// Stream implements the router.LogAdapter interface.
func (a *GelfAdapter) Stream(logstream chan *router.Message) {
	for m := range logstream {

		msg := GelfMessage{
			Version:        "1.1",
      //Host:           os.hostname(),
			ShortMessage:   m.Data,
			Timestamp:      float64(m.Time.UnixNano()) / float64(time.Second),
			ContainerId:    m.Container.ID,
			ContainerName:  m.Container.Name,
			ContainerCmd:   strings.Join(m.Container.Config.Cmd," "),
			ImageId:        m.Container.Image,
			ImageName:      m.Container.Config.Image,
		}
		js, err := json.Marshal(msg)
		if err != nil {
			log.Println("Graylog:", err)
			continue
		}
		_, err = a.conn.Write(js)
		if err != nil {
			log.Println("Graylog:", err)
			continue
		}
	}
}

type GelfMessage struct {
	Version      string  `json:"version"`
	Host         string  `json:"host"`
	ShortMessage string  `json:"short_message"`
	FullMessage  string  `json:"message,omitempty"`
	Timestamp    float64 `json:"timestamp,omitempty"`
	Level        int     `json:"level,omitempty"`

	ImageId        string `json:"image_id,omitempty"`
	ImageName      string `json:"image_name,omitempty"`
	ContainerId    string `json:"container_id,omitempty"`
	ContainerName  string `json:"container_name,omitempty"`
	ContainerCmd   string `json:"command,omitempty"`
}

//version":"1.1","host":"devops-us-east-1a-20150923080634614","short_message":"Hello","level":6,"facility":"","@version":"1","@timestamp":"2015-10-08T19:15:43.089Z","source_host":"127.0.0.1","message":"","command":"echo Hello","container_name":"loving_bartik","created":"2015-10-08T19:15:42.885062465Z","image_id":"d7057cb020844f245031d27b76cb18af05db1cc3a96a29fa7777af75f5ac91a3","image_name":"busybox","tag":"develop"}

