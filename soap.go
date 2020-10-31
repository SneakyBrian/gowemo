package gowemo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"text/template"
)

var (
	requestTmpl *template.Template
)

func init() {
	requestTmpl = template.Must(template.New("request").Parse(`<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
<s:Body>
<u:{{.Action}} xmlns:u="{{.ServiceType}}">
{{ .RequestBody }}
</u:{{.Action}}>
</s:Body>
</s:Envelope>`))
}

func soapBody(serviceType string, action string, request interface{}) (*bytes.Buffer, error) {

	reqXMLBytes, err := xml.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("xml.Marshal: %v", err)
	}

	reqXML := string(reqXMLBytes)

	var buf bytes.Buffer

	err = requestTmpl.Execute(&buf, &struct {
		Action      string
		RequestBody string
		ServiceType string
	}{
		action,
		reqXML,
		serviceType,
	})

	return &buf, nil
}

func sendSoapRequest(controlURL *url.URL, buf *bytes.Buffer, serviceType string, action string) ([]byte, error) {

	client := &http.Client{}

	// build a new request, but not doing the POST yet
	req, err := http.NewRequest("POST", controlURL.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}

	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	// have to set this header like this so it doesn't get changed to "Soapaction"
	req.Header["SOAPACTION"] = []string{"\"" + serviceType + "#" + action + "\""}

	// Save a copy of this request for debugging.
	// requestDump, err := httputil.DumpRequest(req, true)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(requestDump))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do(req): %v", err)
	}

	defer resp.Body.Close()

	// fmt.Println()

	// respDump, err := httputil.DumpResponse(resp, true)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(respDump))
	// fmt.Println(string(data))

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read http response body: %v", err)
	}

	return data, nil
}

func parseResponse(action string, data []byte, response interface{}) error {

	outbuf := bytes.NewBuffer(data)

	decoder := xml.NewDecoder(outbuf)
	var inElement string
	process := false

	for {
		// Read tokens from the XML document in a stream.
		t, err := decoder.Token()
		if err != nil && err != io.EOF {
			return fmt.Errorf("Token: %v", err)
		}
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// println(inElement)
			// ...and its name is action + "Response"
			if inElement == action+"Response" {
				//set the flag to process the next element
				process = true
				continue
			}
			if process {
				// decode a whole chunk of following XML into the response
				if err := decoder.DecodeElement(&response, &se); err != nil {
					// fmt.Println(err)
					return fmt.Errorf("DecodeElement: %v", err)
				}
				process = false
			}
		default:
		}
	}

	return nil
}
