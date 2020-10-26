package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SneakyBrian/gowemo"
)

func main() {

	deviceMap := make(map[string]*gowemo.Device)

	deviceChan := make(chan *gowemo.Device)
	removeDeviceChan := make(chan string)

	// run the discovery process every 1 minute
	go func() {
		for {
			gowemo.Discover(gowemo.ServiceBasicEvent, deviceChan, 10*time.Second)
			time.Sleep(1 * time.Minute)
		}
	}()

	// and then monitor for new and removed devices
	quitMonitorChan := gowemo.Monitor(gowemo.ServiceBasicEvent, deviceChan, removeDeviceChan)

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case device := <-deviceChan:
			// check if we've got it in the map of discovered devices
			_, ok := deviceMap[device.USN]
			if !ok {
				// if not, then it's a new device
				fmt.Printf("New Device Found! %s - %s\n", device.Data.Device.FriendlyName, device.Data.Device.BinaryState)
				deviceMap[device.USN] = device

				state, _ := device.GetBinaryState()
				// try toggling the switch state...
				fmt.Printf("%s - %s\n", device.Data.Device.FriendlyName, device.Data.Device.BinaryState)
				if state == "1" {
					device.SetBinaryState("0")
				} else {
					device.SetBinaryState("1")
				}
				fmt.Printf("%s - %s\n", device.Data.Device.FriendlyName, device.Data.Device.BinaryState)
			}
		case removed := <-removeDeviceChan:
			// check if we've got it in the map of discovered devices
			device, ok := deviceMap[removed]
			if ok {
				// if we have, then we should remove it
				fmt.Printf("Device Removed! %s - %s\n", device.Data.Device.FriendlyName, device.Data.Device.BinaryState)
				delete(deviceMap, removed)
			}
		case sig := <-termChan:
			fmt.Printf("We get signal! %v\n", sig)
			close(quitMonitorChan)
			fmt.Println("exiting...")
			return
		}
	}
}
