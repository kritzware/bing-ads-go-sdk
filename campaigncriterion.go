package bingads

import (
	"encoding/xml"
)

type NilId struct {
	Value int64  `xml:",chardata"`
	Attr  string `xml:",attr"`
}

func st(name string, attrs ...string) xml.StartElement {
	ret := xml.StartElement{
		Name: xml.Name{Local: name},
	}

	for i := 0; i < len(attrs); i += 2 {
		ret.Attr = append(ret.Attr, xml.Attr{xml.Name{Local: attrs[i]}, attrs[i+1]})
	}

	return ret
}

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
	e.EncodeElement("", st("Id", "i:nil", "true"))

	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

//FIXME: make Id nullable int?
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
	Type     string
	TypeAttr string `xml:"i:type,attr"`
	Amount   float64
}

func (c *CampaignService) AddCampaignCriterions(t CriterionType, cs []CampaignCriterion) ([]int64, error) {
	req := AddCampaignCriterionsRequest{
		NS:                 "https://bingads.microsoft.com/CampaignManagement/v11",
		CriterionType:      t,
		CampaignCriterions: cs,
	}

	resp, err := c.client.SendRequest(req, c.endpoint, "AddCampaignCriterions")

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
	CampaignCriterionIds []int64 `xml:"CampaignCriterionIds>a1:long"`
}

func (c *CampaignService) GetCampaignCriterionsByIds(campaign int64) ([]CampaignCriterion, error) {
	req := GetCampaignCriterionsByIdsRequest{
		NS:            "https://bingads.microsoft.com/CampaignManagement/v11",
		CampaignId:    campaign,
		CriterionType: ProductScope,
	}

	resp, err := c.client.SendRequest(req, c.endpoint, "GetCampaignCriterionsByIds")

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
