package gowemo

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DeviceData describes the data returned for the device
type DeviceData struct {
	XMLName     xml.Name `xml:"root"`
	Text        string   `xml:",chardata"`
	Xmlns       string   `xml:"xmlns,attr"`
	SpecVersion struct {
		Text  string `xml:",chardata"`
		Major string `xml:"major"`
		Minor string `xml:"minor"`
	} `xml:"specVersion"`
	Device struct {
		Text             string `xml:",chardata"`
		DeviceType       string `xml:"deviceType"`
		FriendlyName     string `xml:"friendlyName"`
		Manufacturer     string `xml:"manufacturer"`
		ManufacturerURL  string `xml:"manufacturerURL"`
		ModelDescription string `xml:"modelDescription"`
		ModelName        string `xml:"modelName"`
		ModelNumber      string `xml:"modelNumber"`
		ModelURL         string `xml:"modelURL"`
		SerialNumber     string `xml:"serialNumber"`
		UDN              string `xml:"UDN"`
		UPC              string `xml:"UPC"`
		MacAddress       string `xml:"macAddress"`
		FirmwareVersion  string `xml:"firmwareVersion"`
		IconVersion      string `xml:"iconVersion"`
		BinaryState      string `xml:"binaryState"`
		NewAlgo          string `xml:"new_algo"`
		IconList         struct {
			Text string `xml:",chardata"`
			Icon struct {
				Text     string `xml:",chardata"`
				Mimetype string `xml:"mimetype"`
				Width    string `xml:"width"`
				Height   string `xml:"height"`
				Depth    string `xml:"depth"`
				URL      string `xml:"url"`
			} `xml:"icon"`
		} `xml:"iconList"`
		ServiceList struct {
			Text    string `xml:",chardata"`
			Service []struct {
				Text        string `xml:",chardata"`
				ServiceType string `xml:"serviceType"`
				ServiceID   string `xml:"serviceId"`
				ControlURL  string `xml:"controlURL"`
				EventSubURL string `xml:"eventSubURL"`
				SCPDURL     string `xml:"SCPDURL"`
			} `xml:"service"`
		} `xml:"serviceList"`
		PresentationURL string `xml:"presentationURL"`
	} `xml:"device"`
}

func getDeviceData(location string, deviceData *DeviceData) error {

	if deviceData == nil {
		return errors.New("deviceData is nil")
	}

	resp, err := http.Get(location)
	if err != nil {
		return fmt.Errorf("http get: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response status: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read http response body: %v", err)
	}

	err = xml.Unmarshal(data, deviceData)
	if err != nil {
		return fmt.Errorf("device data xml unmarshal: %v", err)
	}

	return nil
}
