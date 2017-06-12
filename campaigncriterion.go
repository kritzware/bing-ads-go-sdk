package bingads

import (
	"encoding/xml"
)

func (s *CampaignCriterion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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
	e.Encode(s.Type)

	if s.Id == 0 {
		e.EncodeElement("", st("Id", "i:nil", "true"))
	} else {
		e.EncodeElement(s.Id, st("Id"))
	}

	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

//TODO: write marshal fns
type CampaignCriterion struct {
	TypeAttr                string `xml:"i:type,attr"`
	CampaignId              int64
	Criterion               Criterion
	Id                      int64
	Nil                     string          `xml:"i:nil,attr"`
	Status                  CriterionStatus `xml:",omitempty"`
	Type                    string          `xml:",omitempty"`
	ForwardCompatibilityMap []string

	//	CriterionBid
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

func (c *CampaignService) AddCampaignCriterions(t CriterionType, cs []CampaignCriterion) ([]int64, error) {
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
	CriterionType      CriterionType
}

type AddCampaignCriterionsResponse struct {
	CampaignCriterionIds []int64 `xml:"CampaignCriterionIds>a:long"`
}

func (c *CampaignService) GetCampaignCriterionsByIds(campaign int64) ([]CampaignCriterion, error) {
	req := GetCampaignCriterionsByIdsRequest{
		NS:            "https://bingads.microsoft.com/CampaignManagement/v11",
		CampaignId:    campaign,
		CriterionType: ProductScope,
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
	XMLName       xml.Name      `xml:"GetCampaignCriterionsByIdsRequest"`
	CampaignId    int64         `xml:"CampaignId"`
	CriterionType CriterionType `xml:"CriterionType"`
	NS            string        `xml:"xmlns,attr"`
}

type GetCampaignCriterionsByIdsResponse struct {
	CampaignCriterions []CampaignCriterion `xml:"https://bingads.microsoft.com/CampaignManagement/v11 CampaignCriterions>CampaignCriterion"`
}
