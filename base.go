package bingads

import (
	"encoding/xml"
	"net/http"
)

var (
	EnvelopeNamespace = "http://schemas.xmlsoap.org/soap/envelope/"
	BingNamespace     = "https://bingads.microsoft.com/CampaignManagement/v10"
)

type RequestEnvelope struct {
	XMLName xml.Name `xml:"s:Envelope"`
	EnvNS   string   `xml:"xmlns:i,attr"`
	EnvSS   string   `xml:"xmlns:s,attr"`
	Header  RequestHeader
	Body    RequestBody
}

type SoapResponseBody struct {
	OperationResponse []byte `xml:",innerxml"`
}

type SoapResponseEnvelope struct {
	XMLName xml.Name         `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  TrackingId       `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`
	Body    SoapResponseBody `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}

type TrackingId struct {
	Nil        bool   `xml:"http://www.w3.org/2001/XMLSchema-instance nil,attr"`
	TrackingId string `xml:"https://adcenter.microsoft.com/v8 TrackingId"`
}

type RequestBody struct {
	XMLName xml.Name `xml:"s:Body"`
	Body    interface{}
}

type RequestHeader struct {
	XMLName             xml.Name `xml:"s:Header"`
	BingNS              string   `xml:"xmlns,attr"`
	Action              string
	AuthenticationToken string `xml:"AuthenticationToken,omitempty"`
	CustomerAccountId   string `xml:"CustomerAccountId"`
	CustomerId          string `xml:"CustomerId"`
	DeveloperToken      string `xml:"DeveloperToken"`
	Password            string `xml:"Password,omitempty"`
	Username            string `xml:"UserName,omitempty"`
}

type Session struct {
	AccountId      string
	CustomerId     string
	DeveloperToken string
	AuthToken      string
	Username       string
	Password       string
	HTTPClient     HttpClient
	ClientId       string
	ClientSecret   string
	RefreshToken   string
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Config struct {
	HTTPClient     *http.Client
	AccountId      string
	CustomerId     string
	DeveloperToken string
	AuthToken      string
	Username       string
	Password       string
}

func NewSession(config *Config) *Session {
	s := &Session{
		AccountId:      config.AccountId,
		CustomerId:     config.CustomerId,
		DeveloperToken: config.DeveloperToken,
		AuthToken:      config.AuthToken,
		Username:       config.Username,
		Password:       config.Password,
		HTTPClient:     config.HTTPClient,
	}

	if s.HTTPClient == nil {
		s.HTTPClient = &http.Client{}
	}

	return s
}
