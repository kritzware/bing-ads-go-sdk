package bingads

import (
	"encoding/xml"
	"fmt"
)

type ReportingService struct {
	Endpoint string
	Session  *Session
}

func NewReportingService(session *Session) *ReportingService {
	return &ReportingService{
		Endpoint: "https://api.bingads.microsoft.com/Api/Advertiser/Reporting/v12/ReportingService.svc",
		Session:  session,
	}
}

/*
Aggregation ::
Summary
Daily
Weekly
Monthly
Yearly
*/

//TODO: Filter
type PerformanceReportRequest struct {
	XMLName     xml.Name `xml:"ReportRequest"`
	Type        string   `xml:"i:type,attr"`
	Aggregation string
	Columns     []string `xml:"Columns>AdGroupPerformanceReportColumn"`
	Scope       ReportScope
	Time        ReportTime
}

type ProductDimensionPerformanceReportRequest PerformanceReportRequest
type ProductPartitionPerformanceReportRequest PerformanceReportRequest
type AdGroupPerformanceReportRequest PerformanceReportRequest

type ReportScope struct {
	XMLName    xml.Name              `xml:"Scope"`
	AccountIds Longs                 `xml:"AccountIds>long,omitempty"`
	AdGroups   []AdGroupReportScope  `xml:"AdGroups>AdGroupReportScope,omitempty"`
	Campaigns  []CampaignReportScope `xml:"Campaigns>CampaignReportScope,omitempty"`
}

func (s ReportScope) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Scope"
	e.EncodeToken(start)
	if len(s.AccountIds) > 0 {
		acts := st("AccountIds", "xmlns:a1", "http://schemas.microsoft.com/2003/10/Serialization/Arrays")
		e.EncodeToken(acts)
		for i := 0; i < len(s.AccountIds); i++ {
			e.EncodeElement(s.AccountIds[i], st("a1:long"))
		}
		e.EncodeToken(acts.End())
	}

	if len(s.AdGroups) > 0 {
		acts := st("AdGroups")
		e.EncodeToken(acts)
		e.Encode(s.AdGroups)
		e.EncodeToken(acts.End())
	}

	if len(s.Campaigns) > 0 {
		acts := st("Campaigns")
		e.EncodeToken(acts)
		e.Encode(s.Campaigns)
		e.EncodeToken(acts.End())
	}
	e.EncodeToken(start.End())
	return nil
}

type AdGroupReportScope struct {
	AccountId  int64
	CampaignId int64
	AdGroupId  int64
}

type CampaignReportScope struct {
	AccountId  int64
	CampaignId int64
}

/*
Today | Yesterday | LastSevenDays | ThisWeek | LastWeek | LastFourWeeks | ThisMonth | LastMonth | LastThreeMonths | LastSixMonths | ThisYear | LastYear
*/

type ReportTime struct {
	XMLName              xml.Name `xml:"Time"`
	CustomDateRangeEnd   Date     `xml:",omitempty"`
	CustomDateRangeStart Date     `xml:",omitempty"`
	PredefinedTime       string   `xml:",omitempty"`
}

func (s ReportTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Time"
	e.EncodeToken(start)
	if s.PredefinedTime != "" {
		e.EncodeElement(s.PredefinedTime, st("PredefinedTime"))
	} else {
		e.EncodeElement(s.CustomDateRangeEnd, st("CustomDateRangeEnd"))
		e.EncodeElement(s.CustomDateRangeStart, st("CustomDateRangeStart"))
	}
	e.EncodeToken(start.End())
	return nil
}

func (s *ProductDimensionPerformanceReportRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	req := PerformanceReportRequest(*s)
	return marshallPerformanceReportRequest(e, req, "ProductDimensionPerformanceReport")
}

func (s *ProductPartitionPerformanceReportRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	req := PerformanceReportRequest(*s)
	return marshallPerformanceReportRequest(e, req, "ProductPartitionPerformanceReport")
}

func (s *AdGroupPerformanceReportRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	req := PerformanceReportRequest(*s)
	return marshallPerformanceReportRequest(e, req, "AdGroupPerformanceReport")
}

func (s PerformanceReportRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if s.Type == "" {
		return fmt.Errorf("missing report type")
	}
	return marshallPerformanceReportRequest(e, s, s.Type)
}

func marshallPerformanceReportRequest(e *xml.Encoder, s PerformanceReportRequest, t string) error {
	start := st("ReportRequest", "i:type", t+"Request")
	e.EncodeToken(start)
	excludes := []string{"ExcludeReportFooter", "ExcludeReportHeader"}
	for i := 0; i < len(excludes); i++ {
		e.EncodeElement(true, st(excludes[i]))
	}
	if s.Aggregation != "" {
		e.EncodeElement(s.Aggregation, st("Aggregation"))
	}
	if s.Columns == nil || len(s.Columns) == 0 {
		return fmt.Errorf("no columns selected")
	}
	cols := st("Columns", "i:nil", "false")
	e.EncodeToken(cols)
	for i := 0; i < len(s.Columns); i++ {
		e.EncodeElement(s.Columns[i], st(t+"Column"))
	}
	e.EncodeToken(cols.End())
	e.Encode(s.Scope)
	e.Encode(s.Time)
	e.EncodeToken(start.End())
	return nil
}

type SubmitGenerateReportRequest struct {
	XMLName xml.Name `xml:"SubmitGenerateReportRequest"`

	NS            string `xml:"xmlns,attr"`
	ReportRequest interface{}
}

type SubmitGenerateReportResponse struct {
	ReportRequestId string
}

type PollGenerateReportRequest struct {
	XMLName         xml.Name `xml:"PollGenerateReportRequest"`
	NS              string   `xml:"xmlns,attr"`
	ReportRequestId string
}

type PollGenerateReportResponse struct {
	ReportRequestStatus ReportRequestStatus
}

//Status :: Error | Success | Pending
type ReportRequestStatus struct {
	ReportDownloadUrl string
	Status            string
}

func (c *ReportingService) PollGenerateReport(id string) (*ReportRequestStatus, error) {
	req := PollGenerateReportRequest{
		ReportRequestId: id,
		NS:              "https://bingads.microsoft.com/Reporting/v12",
	}
	resp, err := c.Session.reportRequest(req, c.Endpoint, "PollGenerateReport")
	if err != nil {
		return nil, err
	}

	ret := PollGenerateReportResponse{}
	err = xml.Unmarshal(resp, &ret)
	if err != nil {
		return nil, err
	}
	return &ret.ReportRequestStatus, nil
}

func (c *ReportingService) SubmitReportRequest(rr interface{}) (string, error) {
	req := SubmitGenerateReportRequest{
		ReportRequest: rr,
		NS:            "https://bingads.microsoft.com/Reporting/v12",
	}
	resp, err := c.Session.reportRequest(req, c.Endpoint, "SubmitGenerateReport")

	if err != nil {
		return "", err
	}

	ret := SubmitGenerateReportResponse{}
	err = xml.Unmarshal(resp, &ret)

	if err != nil {
		return "", err
	}

	return ret.ReportRequestId, nil

}
