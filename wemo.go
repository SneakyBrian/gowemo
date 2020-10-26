package gowemo

import (
	"log"
	"time"

	"github.com/koron/go-ssdp"
)

// Discover runs discovery process
func Discover(serviceType string, newDeviceChan chan<- *Device, duration time.Duration) {

	searchSeconds := int(duration.Seconds())

	services, err := ssdp.Search(serviceType, searchSeconds, "")
	if err != nil {
		log.Fatal(err)
	}
	for _, service := range services {
		device, err := NewDevice(service.Type, service.USN, service.Location, service.Server)
		if err == nil {
			newDeviceChan <- device
		}
	}
}

// Monitor monitors for new devices
// returns a channel which when signalled will close the monitor
func Monitor(serviceType string, newDeviceChan chan<- *Device, removeDeviceChan chan<- string) chan<- struct{} {

	quitChan := make(chan struct{})

	m := &ssdp.Monitor{
		Alive: func(m *ssdp.AliveMessage) {
			if m.Type == serviceType {
				// fmt.Printf("*** Alive: From=%s Type=%s USN=%s Location=%s Server=%s MaxAge=%d\n",
				// 	m.From.String(), m.Type, m.USN, m.Location, m.Server, m.MaxAge())

				device, err := NewDevice(m.Type, m.USN, m.Location, m.Server)
				if err == nil {
					newDeviceChan <- device
				}
			}
		},
		Bye: func(m *ssdp.ByeMessage) {
			if m.Type == serviceType {
				// fmt.Printf("*** Bye: From=%s Type=%s USN=%s\n", m.From.String(), m.Type, m.USN)
			}
		},
		Search: func(m *ssdp.SearchMessage) {
			if m.Type == serviceType {
				// fmt.Printf("*** Search: From=%s Type=%s\n", m.From.String(), m.Type)
			}
		},
	}

	m.Start()

	go func() {
		// when the quit chan is signalled
		// then we close the monitor
		<-quitChan
		m.Close()
	}()

	return quitChan
}
