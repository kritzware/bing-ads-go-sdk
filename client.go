package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (b *BingClient) MakeRequest(body interface{}) RequestEnvelope {
	return RequestEnvelope{
		EnvNS:  "http://schemas.xmlsoap.org/soap/envelope/",
		BingNS: "https://bingads.microsoft.com/CampaignManagement/v11",
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

	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	res := SoapResponseEnvelope{}

	err = xml.Unmarshal(raw, &res)
	if err != nil {
		return
	}

	fmt.Println(string(res.Body.OperationResponse))

	return res.Body.OperationResponse, err
}

func New(customerAccountId string, customerId string, developerToken string, authToken string, username string, password string) *BingClient {
	return &BingClient{
		CustomerAccountId: customerAccountId,
		customerId:        customerId,
		developerToken:    developerToken,
		authToken:         authToken,
		username:          username,
		password:          password,
	}
}
