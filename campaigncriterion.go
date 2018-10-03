package bingads

import (
	"encoding/xml"
	"errors"
)

func (s CampaignCriterion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = []xml.Attr{xml.Attr{Name: xml.Name{Local: "i:type"}, Value: s.Type}}
	e.EncodeToken(start)
	/*
		err := e.EncodeElement("", xml.StartElement{
			Name: xml.Name{Local: "Id"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "i:nil"}, Value: "true"},
			},
		})
	*/

	e.EncodeElement(s.CampaignId, st("CampaignId"))
	e.Encode(s.Criterion)

	if s.Id == 0 {
		e.EncodeElement("", st("Id", "i:nil", "true"))
	} else {
		e.EncodeElement(s.Id, st("Id"))
	}

	if s.Status != "" {
		e.EncodeElement(s.Status, st("Status"))
	}
	e.EncodeElement(s.Type, st("Type"))
	e.Encode(s.CriterionBid)
	//marshalCriterion(s.Criterion, e)

	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

type CampaignCriterion struct {
	CampaignId   int64
	Criterion    ProductScope
	Id           int64
	Status       CriterionStatus `xml:",omitempty"`
	Type         string          `xml:",omitempty"`
	CriterionBid CriterionBid
}

type CriterionStatus string

const (
	Active  CriterionStatus = "Active"
	Paused                  = "Paused"
	Deleted                 = "Deleted"
)

type FixedBid struct {
	Amount float64
}
type Multiplier struct {
	Multiplier float64
}
type CriterionBid struct {
	Type   string
	Amount float64
}

func (s CriterionBid) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = ats("i:type", "FixedBid")
	e.EncodeToken(start)

	e.EncodeElement("FixedBid", st("Type"))
	e.EncodeElement(s.Amount, st("Amount"))

	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

var CampaignCriterionAlreadyExists = errors.New("CampaignCriterionAlreadyExists")

func (c *CampaignService) AddCampaignCriterions(t string, cs []CampaignCriterion) ([]int64, error) {
	req := AddCampaignCriterionsRequest{
		NS:                 "https://bingads.microsoft.com/CampaignManagement/v12",
		CriterionType:      t,
		CampaignCriterions: cs,
	}

	resp, err := c.Session.SendRequest(req, c.Endpoint, "AddCampaignCriterions")

	if err != nil {
		return nil, err
	}

	ret := &AddCampaignCriterionsResponse{}
	if err := xml.Unmarshal(resp, ret); err != nil {
		return nil, err
	}

	if ret.NestedPartialErrors != nil {
		switch ret.NestedPartialErrors.ErrorCode {
		case "CampaignCriterionAlreadyExists":
			return nil, CampaignCriterionAlreadyExists
		default:
			return nil, ret.NestedPartialErrors
		}
	}

	return ret.CampaignCriterionIds, nil
}

type AddCampaignCriterionsRequest struct {
	XMLName            xml.Name            `xml:"AddCampaignCriterionsRequest"`
	NS                 string              `xml:"xmlns,attr"`
	CampaignCriterions []CampaignCriterion `xml:"CampaignCriterions>CampaignCriterion"`
	CriterionType      string
}

type BatchErrorCollection struct {
	Code        int
	Details     string
	ErrorCode   string
	Index       int
	Message     string
	Type        string
	BatchErrors []BatchError
}

func (s BatchErrorCollection) Error() string {
	return s.Message
}

type AddCampaignCriterionsResponse struct {
	CampaignCriterionIds Longs                 `xml:"CampaignCriterionIds>long"`
	NestedPartialErrors  *BatchErrorCollection `xml:"NestedPartialErrors>BatchErrorCollection"`
}

type UpdateCampaignCriterionsRequest struct {
	XMLName            xml.Name            `xml:"UpdateCampaignCriterionsRequest"`
	NS                 string              `xml:"xmlns,attr"`
	CampaignCriterions []CampaignCriterion `xml:"CampaignCriterions>CampaignCriterion"`
	CriterionType      string
}

type UpdateCampaignCriterionsResponse struct {
	NestedPartialErrors *BatchErrorCollection `xml:"NestedPartialErrors>BatchErrorCollection"`
}

func (c *CampaignService) UpdateCampaignCriterions(t string, cs []CampaignCriterion) error {
	req := UpdateCampaignCriterionsRequest{
		NS:                 "https://bingads.microsoft.com/CampaignManagement/v12",
		CampaignCriterions: cs,
		CriterionType:      t,
	}
	resp, err := c.Session.SendRequest(req, c.Endpoint, "UpdateCampaignCriterions")

	if err != nil {
		return err
	}

	ret := &UpdateCampaignCriterionsResponse{}
	if err := xml.Unmarshal(resp, ret); err != nil {
		return err
	}

	if ret.NestedPartialErrors != nil {
		return ret.NestedPartialErrors
	}

	return nil
}

func (c *CampaignService) GetCampaignCriterionsByIds(campaign int64) ([]CampaignCriterion, error) {
	req := GetCampaignCriterionsByIdsRequest{
		NS:            "https://bingads.microsoft.com/CampaignManagement/v12",
		CampaignId:    campaign,
		CriterionType: "ProductScope",
	}

	resp, err := c.Session.SendRequest(req, c.Endpoint, "GetCampaignCriterionsByIds")

	if err != nil {
		return nil, err
	}

	ret := GetCampaignCriterionsByIdsResponse{}
	err = xml.Unmarshal(resp, &ret)
	return ret.CampaignCriterions, err
}

type GetCampaignCriterionsByIdsRequest struct {
	XMLName       xml.Name `xml:"GetCampaignCriterionsByIdsRequest"`
	CampaignId    int64    `xml:"CampaignId"`
	CriterionType string   `xml:"CriterionType"`
	NS            string   `xml:"xmlns,attr"`
}

type GetCampaignCriterionsByIdsResponse struct {
	CampaignCriterions []CampaignCriterion `xml:"https://bingads.microsoft.com/CampaignManagement/v12 CampaignCriterions>CampaignCriterion"`
}
