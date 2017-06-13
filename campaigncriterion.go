package bingads

import (
	"encoding/xml"
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

	e.EncodeElement(s.Status, st("Status"))
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

func (c *CampaignService) AddCampaignCriterions(t string, cs []CampaignCriterion) ([]int64, error) {
	req := AddCampaignCriterionsRequest{
		NS:                 "https://bingads.microsoft.com/CampaignManagement/v11",
		CriterionType:      t,
		CampaignCriterions: cs,
	}

	resp, err := c.Client.SendRequest(req, c.Endpoint, "AddCampaignCriterions")

	if err != nil {
		return nil, err
	}

	ret := AddCampaignCriterionsResponse{}
	err = xml.Unmarshal(resp, &ret)
	return ret.CampaignCriterionIds, err
}

type AddCampaignCriterionsRequest struct {
	XMLName            xml.Name            `xml:"AddCampaignCriterionsRequest"`
	NS                 string              `xml:"xmlns,attr"`
	CampaignCriterions []CampaignCriterion `xml:"CampaignCriterions>CampaignCriterion"`
	CriterionType      string
}

type AddCampaignCriterionsResponse struct {
	CampaignCriterionIds Longs `xml:"CampaignCriterionIds>long"`
}

func (c *CampaignService) GetCampaignCriterionsByIds(campaign int64) ([]CampaignCriterion, error) {
	req := GetCampaignCriterionsByIdsRequest{
		NS:            "https://bingads.microsoft.com/CampaignManagement/v11",
		CampaignId:    campaign,
		CriterionType: "ProductScope",
	}

	resp, err := c.Client.SendRequest(req, c.Endpoint, "GetCampaignCriterionsByIds")

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
	CampaignCriterions []CampaignCriterion `xml:"https://bingads.microsoft.com/CampaignManagement/v11 CampaignCriterions>CampaignCriterion"`
}
