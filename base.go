package main

import (
	"encoding/xml"
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
	accountId      string
	customerId     string
	developerToken string
	authToken      string
	username       string
	password       string
}

func NewBingClient(customerAccountId string, customerId string, developerToken string, authToken string, username string, password string) *BingClient {
	return &BingClient{
		accountId:      customerAccountId,
		customerId:     customerId,
		developerToken: developerToken,
		authToken:      authToken,
		username:       username,
		password:       password,
	}
}
