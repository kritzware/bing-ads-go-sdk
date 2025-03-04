package bingads

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

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

var debug = os.Getenv("BING_SDK_DEBUG")

func (b *Session) SendRequest(body interface{}, endpoint string, soapAction string) ([]byte, error) {
	return b.sendRequest(body, endpoint, soapAction, campaignns)
}

var campaignns = "https://bingads.microsoft.com/CampaignManagement/v13"
var reportingns = "https://bingads.microsoft.com/Reporting/v13"

func (b *Session) reportRequest(body interface{}, endpoint string, soapAction string) ([]byte, error) {
	return b.sendRequest(body, endpoint, soapAction, reportingns)
}

//FIXME: namespace
func (b *Session) sendRequest(body interface{}, endpoint string, soapAction string, ns string) ([]byte, error) {
	header := RequestHeader{
		BingNS:            ns,
		Action:            soapAction,
		CustomerAccountId: b.AccountId,
		CustomerId:        b.CustomerId,
		DeveloperToken:    b.DeveloperToken,
	}
	if b.TokenSource != nil {
		token, err := b.TokenSource.Token()
		if err != nil {
			return nil, err
		}
		header.AuthenticationToken = token.AccessToken
	} else {
		header.Username = b.Username
		header.Password = b.Password
	}

	envelope := RequestEnvelope{
		EnvNS:  "http://www.w3.org/2001/XMLSchema-instance",
		EnvSS:  "http://schemas.xmlsoap.org/soap/envelope/",
		Header: header,
		Body: RequestBody{
			Body: body,
		},
	}

	req, err := xml.MarshalIndent(envelope, "", "  ")

	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(req))

	if err != nil {
		return nil, err
	}

	httpRequest.Header.Add("Content-Type", "text/xml; charset=utf-8")
	httpRequest.Header.Add("SOAPAction", soapAction)

	response, err := b.HTTPClient.Do(httpRequest)

	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if debug != "" {
		fmt.Println(string(req))
	}
	//fmt.Println(string(raw))

	res := SoapResponseEnvelope{}

	err = xml.Unmarshal(raw, &res)
	if err != nil {
		return nil, err
	}

	if debug != "" {
		fmt.Println(">>>")
		fmt.Println(string(res.Body.OperationResponse))
		fmt.Println(">>>")
	}

	switch response.StatusCode {
	case 400, 401, 403, 405, 500:
		fault := Fault{}
		err = xml.Unmarshal(res.Body.OperationResponse, &fault)
		if err != nil {
			return res.Body.OperationResponse, err
		}
		for _, e := range fault.Detail.Errors.AdApiErrors {
			switch e.ErrorCode {
			case "AuthenticationTokenExpired", "InvalidCredentials", "InternalError", "CallRateExceeded":
				return res.Body.OperationResponse, &baseError{
					code:    e.ErrorCode,
					origErr: &fault.Detail.Errors,
				}
			}
		}
		return res.Body.OperationResponse, &fault.Detail.Errors
	}

	return res.Body.OperationResponse, err
}

type SessionConfig struct {
	OAuth2Config   *oauth2.Config
	OAuth2Token    *oauth2.Token
	AccountId      string
	CustomerId     string
	DeveloperToken string
	HTTPClient     HttpClient
}

func NewSession(config SessionConfig) *Session {
	tokenSource := config.OAuth2Config.TokenSource(context.TODO(), config.OAuth2Token)

	return &Session{
		AccountId:      config.AccountId,
		CustomerId:     config.CustomerId,
		DeveloperToken: config.DeveloperToken,
		HTTPClient:     config.HTTPClient,
		TokenSource:    tokenSource,
	}
}
