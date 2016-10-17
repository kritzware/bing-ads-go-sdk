package base

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	EnvelopeNamespace = "http://schemas.xmlsoap.org/soap/envelope/"
	BingNamespace     = "https://bingads.microsoft.com/CampaignManagement/v10"
)

type RequestEnvelope struct {
	XMLName xml.Name `xml:"SOAP-ENV:Envelope"`
	EnvNS   string   `xml:"xmlns:SOAP-ENV,attr"`
	BingNS  string   `xml:"xmlns:ns1,attr"`
	Header  RequestHeader
	Body    RequestBody
}

type ResponseEnvelope struct {
	XMLName xml.Name    `xml:"s:Envelope"`
	Body    interface{} `xml:"s:Body"`
}

type RequestBody struct {
	XMLName xml.Name `xml:"SOAP-ENV:Body"`
	Body    interface{}
}

type RequestHeader struct {
	XMLName             xml.Name `xml:"SOAP-ENV:Header"`
	CustomerAccountId   string   `xml:"ns1:CustomerAccountId"`
	CustomerId          string   `xml:"ns1:CustomerId"`
	DeveloperToken      string   `xml:"ns1:DeveloperToken"`
	AuthenticationToken string   `xml:"ns1:AuthenticationToken"`
	Username            string   `xml:"ns1:UserName"`
	Password            string   `xml:"ns1:Password"`
}

type BingClient struct {
	CustomerAccountId string
	customerId        string
	developerToken    string
	authToken         string
	username          string
	password          string
}

func (b *BingClient) MakeRequest(body interface{}) RequestEnvelope {
	return RequestEnvelope{
		EnvNS:  "http://schemas.xmlsoap.org/soap/envelope/",
		BingNS: "https://bingads.microsoft.com/CampaignManagement/v10",
		Header: RequestHeader{
			CustomerAccountId:   b.CustomerAccountId,
			AuthenticationToken: b.authToken,
			DeveloperToken:      b.developerToken,
			Username:            b.username,
			Password:            b.password,
		},
		Body: RequestBody{
			Body: body,
		},
	}
}

func (b *BingClient) SendRequest(r RequestEnvelope, endpoint string, soapAction string) (resp []byte, err error) {
	req, err := xml.Marshal(r)

	if err != nil {
		return
	}

	fmt.Printf("\n\n\n%s\n\n\n\n", string(req))

	httpRequest, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(req))

	if err != nil {
		return
	}

	httpRequest.Header.Add("Content-Type", "text/xml; charset=utf-8")
	httpRequest.Header.Add("SOAPAction", soapAction)

	c := http.Client{}

	response, err := c.Do(httpRequest)

	if err != nil {
		return
	}

	return ioutil.ReadAll(response.Body)
}

func NewBingClient(customerAccountId string, customerId string, developerToken string, authToken string, username string, password string) *BingClient {
	return &BingClient{
		CustomerAccountId: customerAccountId,
		customerId:        customerId,
		developerToken:    developerToken,
		authToken:         authToken,
		username:          username,
		password:          password,
	}
}