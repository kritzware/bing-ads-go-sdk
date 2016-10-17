package error

import "encoding/xml"

type SoapFault struct {
	Name xml.Name `xml:"s:Fault"`
	FaultCode string `xml:"faultcode"`
	FaultString string `xml:"faultstring"`
	Errors []ApiError `xml:"detail>AdApiFaultDetail>Errors"`
}

type ApiError struct {
	Name xml.Name `xml:"AdApiError"`
	Code int `xml:"Code"`
	ErrorCode string `xml:"ErrorCode"`
	Message string `xml:"Message"`
}