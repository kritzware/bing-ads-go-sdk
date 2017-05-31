package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client interface {
	SendRequest(interface{}, string, string) ([]byte, error)
}
type AuthHeader struct {
	Action              string `xml:"https://adcenter.microsoft.com/v8 Action"`
	ApplicationToken    string `xml:"https://adcenter.microsoft.com/v8 ApplicationToken"`
	AuthenticationToken string `xml:"https://adcenter.microsoft.com/v8 AuthenticationToken"`
	CustomerAccountId   int64  `xml:"https://adcenter.microsoft.com/v8 CustomerAccountId"`
	CustomerId          int64  `xml:"https://adcenter.microsoft.com/v8 CustomerId"`
	DeveloperToken      string `xml:"https://adcenter.microsoft.com/v8 DeveloperToken"`
	UserName            string `xml:"https://adcenter.microsoft.com/v8 UserName"`
	Password            string `xml:"https://adcenter.microsoft.com/v8 Password"`
}

type AdApiError struct {
	Code      int64  `xml:"AdApiError>Code"`
	Details   string `xml:"AdApiError>Details"`
	ErrorCode string `xml:"AdApiError>ErrorCode"`
	Message   string `xml:"AdApiError>Message"`
}

type BatchError struct {
	Code      int64  `xml:"BatchError>Code"`
	Details   string `xml:"BatchError>Details"`
	ErrorCode string `xml:"BatchError>ErrorCode"`
	Index     int64  `xml:"BatchError>Index"`
	Message   string `xml:"BatchError>Message"`
}

type EditorialError struct {
	Appealable       bool   `xml:"EditorialError>Appealable"`
	Code             int64  `xml:"EditorialError>Code"`
	DisapprovedText  string `xml:"EditorialError>DisapprovedText"`
	ErrorCode        string `xml:"EditorialError>ErrorCode"`
	Index            int64  `xml:"EditorialError>Index"`
	Message          string `xml:"EditorialError>Message"`
	PublisherCountry string `xml:"EditorialError>PublisherCountry"`
}

type GoalError struct {
	BatchErrors []BatchError `xml:"GoalError>BatchErrors"`
	Index       int64        `xml:"GoalError>Index"`
	StepErrors  []BatchError `xml:"GoalError>StepErrors"`
}

type OperationError struct {
	Code      int64  `xml:"OperationError>Code"`
	Details   string `xml:"OperationError>Details"`
	ErrorCode string `xml:"OperationError>ErrorCode"`
	Message   string `xml:"OperationError>Message"`
}
type Fault struct {
	FaultCode   string `xml:"faultcode"`
	FaultString string `xml:"faultstring"`
	Detail      struct {
		XMLName xml.Name   `xml:"detail"`
		Errors  ErrorsType `xml:",any"`
	}
}

type ErrorsType struct {
	TrackingId      string           `xml:"TrackingId"`
	AdApiErrors     []AdApiError     `xml:"Errors"`
	BatchErrors     []BatchError     `xml:"BatchErrors"`
	EditorialErrors []EditorialError `xml:"EditorialErrors"`
	GoalErrors      []GoalError      `xml:"GoalErrors"`
	OperationErrors []OperationError `xml:"OperationErrors"`
}

func (f *ErrorsType) Error() string {
	errors := []string{}
	for _, e := range f.AdApiErrors {
		errors = append(errors, fmt.Sprintf("%s", e.Message))
	}
	for _, e := range f.BatchErrors {
		errors = append(errors, fmt.Sprintf("%s", e.Message))
	}
	for _, e := range f.EditorialErrors {
		errors = append(errors, fmt.Sprintf("%s", e.Message))
	}
	for _, e := range f.GoalErrors {
		for _, be := range e.BatchErrors {
			errors = append(errors, fmt.Sprintf("%s", be.Message))
		}
		for _, be := range e.StepErrors {
			errors = append(errors, fmt.Sprintf("%s", be.Message))
		}
	}
	for _, e := range f.OperationErrors {
		errors = append(errors, fmt.Sprintf("%s", e.Message))
	}
	return strings.Join(errors, "\n")
}

func (b *BingClient) SendRequest(body interface{}, endpoint string, soapAction string) (resp []byte, err error) {
	envelope := RequestEnvelope{
		EnvNS:  "http://schemas.xmlsoap.org/soap/envelope/",
		BingNS: "https://bingads.microsoft.com/CampaignManagement/v11",
		Header: RequestHeader{
			CustomerAccountId:   b.accountId,
			AuthenticationToken: b.authToken,
			DeveloperToken:      b.developerToken,
			Username:            b.username,
			Password:            b.password,
		},
		Body: RequestBody{
			Body: body,
		},
	}

	req, err := xml.Marshal(envelope)

	if err != nil {
		return
	}

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

	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	res := SoapResponseEnvelope{}

	err = xml.Unmarshal(raw, &res)
	if err != nil {
		return
	}

	switch response.StatusCode {
	case 400, 401, 403, 405, 500:
		fault := Fault{}
		err = xml.Unmarshal(res.Body.OperationResponse, &fault)
		if err != nil {
			return res.Body.OperationResponse, err
		}
		return res.Body.OperationResponse, &fault.Detail.Errors //errors
	}

	fmt.Println(string(raw))

	return res.Body.OperationResponse, err
}

func New(customerAccountId string, customerId string, developerToken string, authToken string, username string, password string) *BingClient {
	return &BingClient{
		accountId:      customerAccountId,
		customerId:     customerId,
		developerToken: developerToken,
		authToken:      authToken,
		username:       username,
		password:       password,
	}
}
