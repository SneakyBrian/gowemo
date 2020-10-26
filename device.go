package gowemo

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

// Device represents a Wemo Device
type Device struct {
	Type     string
	USN      string
	Location *url.URL
	Server   string

	Data DeviceData
}

// BinaryStateData is the request/response to set the binary state
type BinaryStateData struct {
	XMLName     xml.Name `xml:"BinaryState"`
	BinaryState string   `xml:",chardata"`
}

// NewDevice creates a new device
func NewDevice(serviceType string, usn string, location string, server string) (*Device, error) {

	locationURL, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("error parsing location url: %v", err)
	}

	device := Device{
		Type:     serviceType,
		USN:      usn,
		Location: locationURL,
		Server:   server,
	}

	err = getDeviceData(location, &device.Data)
	if err != nil {
		return nil, fmt.Errorf("error getting device data: %v", err)
	}

	return &device, nil
}

// UpdateDeviceData Updates the Device Data
func (device *Device) UpdateDeviceData() error {

	err := getDeviceData(device.Location.String(), &device.Data)
	if err != nil {
		return fmt.Errorf("error getting device data: %v", err)
	}

	return nil
}

// GetBinaryState gets the binary state (1=on/0=off for switches)
func (device *Device) GetBinaryState() (string, error) {

	response := &BinaryStateData{}

	err := device.soapAction(ServiceBasicEvent, "GetBinaryState", nil, response)
	if err != nil {
		return "", fmt.Errorf("soap action GetBinaryState: %v", err)
	}

	device.Data.Device.BinaryState = response.BinaryState

	return response.BinaryState, nil
}

// SetBinaryState sets the binary state (1=on/0=off for switches)
func (device *Device) SetBinaryState(state string) (string, error) {

	response := &BinaryStateData{}

	err := device.soapAction(ServiceBasicEvent, "SetBinaryState", &BinaryStateData{BinaryState: state}, response)
	if err != nil {
		return "", fmt.Errorf("soap action SetBinaryState: %v", err)
	}

	device.Data.Device.BinaryState = response.BinaryState

	return response.BinaryState, nil
}

func (device *Device) controlURLFor(serviceType string) (*url.URL, error) {
	// get control url for service type
	controlURL := ""
	for _, service := range device.Data.Device.ServiceList.Service {
		if service.ServiceType == serviceType {
			controlURL = service.ControlURL
		}
	}

	if controlURL == "" {
		return nil, fmt.Errorf("serviceType '%s' not found", serviceType)
	}

	return &url.URL{
		Scheme: device.Location.Scheme,
		Host:   device.Location.Host,
		Path:   controlURL,
	}, nil
}

func (device *Device) soapAction(serviceType string, action string, request interface{}, response interface{}) error {

	// fmt.Println(serviceType)
	// fmt.Println(action)

	// get the control URL
	controlURL, err := device.controlURLFor(serviceType)
	if err != nil {
		return fmt.Errorf("controlURLFor: %v", err)
	}

	// fmt.Println(controlURL)

	// Build Request XML
	buf, err := soapBody(serviceType, action, request)
	if err != nil {
		return fmt.Errorf("soapBody: %v", err)
	}

	// send soap request
	data, err := sendSoapRequest(controlURL, buf, serviceType, action)
	if err != nil {
		return fmt.Errorf("sendSoapRequest: %v", err)
	}

	// Pull out Response XML
	if err := parseResponse(action, data, response); err != nil {
		return fmt.Errorf("parseResponse: %v", err)
	}

	return nil
}
