# GoWemo

Go Library for interfacing with Wemo devices

Only handles Wemo Switches - discovery and turning on/off at present

Still in development, see the [TODO](TODO.md)

## Example

```go

deviceChan := make(chan *gowemo.Device)

// run the discover in a go routine
go gowemo.Discover(deviceChan, 10*time.Second)

// wait until we have a device
device := <-deviceChan

fmt.Printf("New Device Found! %s - %s\n", device.Data.Device.FriendlyName, device.Data.Device.BinaryState)

// try toggling the switch state...
state, _ := device.GetBinaryState()

fmt.Printf("%s - %s\n", device.Data.Device.FriendlyName, device.Data.Device.BinaryState)

if state == "1" {
    // turn switch off
    device.SetBinaryState("0")
} else {
    // turn switch on
    device.SetBinaryState("1")
}

fmt.Printf("%s - %s\n", device.Data.Device.FriendlyName, device.Data.Device.BinaryState)

```