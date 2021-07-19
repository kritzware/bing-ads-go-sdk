package bingads

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

//ProductType1-5 | CategoryL1-5 | Id | Condition | Brand | CustomLabel0-4
type ProductCondition struct {
	Attribute string
	Operand   string
}

//TODO: derived types
//https://msdn.microsoft.com/en-us/library/bing-ads-campaign-management-getcampaigncriterionsbyids.aspx
type ProductPartition struct {
	Type              string
	Condition         ProductCondition
	ParentCriterionId int64
	PartitionType     string `xml:",omitempty"`
}

type ProductScope struct {
	Type       string
	Conditions []ProductCondition `xml:"Conditions>ProductCondition"`
}

func (s ProductScope) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Criterion"
	start.Attr = ats("i:type", "ProductScope")
	e.EncodeToken(start)
	e.EncodeElement("ProductScope", st("Type"))

	e.EncodeToken(st("Conditions"))

	for _, x := range s.Conditions {
		e.EncodeElement(x, st("ProductCondition"))
	}

	e.EncodeToken(xml.EndElement{xml.Name{Local: "Conditions"}})

	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

func (s ProductPartition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Criterion"
	start.Attr = ats("i:type", "ProductPartition")
	e.EncodeToken(start)
	e.EncodeElement("ProductPartition", st("Type"))
	//e.Encode(s.Condition)
	e.EncodeElement(s.Condition, st("Condition"))

	if s.ParentCriterionId != 0 {
		e.EncodeElement(s.ParentCriterionId, st("ParentCriterionId"))
	} else {
		e.EncodeElement("", st("ParentCriterionId", "i:nil", "true"))
	}

	if s.PartitionType != "" {
		e.EncodeElement(s.PartitionType, st("PartitionType"))
	}

	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

type GetAdGroupCriterionsByIdsRequest struct {
	XMLName xml.Name `xml:"GetAdGroupCriterionsByIdsRequest"`

	NS            string `xml:"xmlns,attr"`
	AdGroupId     int64  `xml:"AdGroupId"`
	CriterionType string `xml:"CriterionType"`
}

type GetAdGroupCriterionsByIdsResponse struct {
	AdGroupCriterions []BiddableAdGroupCriterion `xml:"https://bingads.microsoft.com/CampaignManagement/v13 AdGroupCriterions>AdGroupCriterion"`
}

/*
type AdGroupCriterion struct {
	Id           int64 `xml:",omitempty"`
	AdGroupId    int64
	Criterion    ProductPartition
	Status       CriterionStatus
	Type         string
	CriterionBid CriterionBid
}
*/
/*
func marshalCriterion(c Criterion, e *xml.Encoder) error {
	switch c.(type) {
	case ProductPartition:
		return e.Encode(c)
	case ProductScope:
		return e.Encode(c)
	}

	return fmt.Errorf("unknown criterion type")
}
*/

type BiddableAdGroupCriterion struct {
	AdGroupId    int64
	Criterion    ProductPartition
	Id           int64 `xml:",omitempty"`
	Status       string
	Type         string
	CriterionBid CriterionBid
}

//FIXME: maybe switch to pointers
func (s BiddableAdGroupCriterion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = []xml.Attr{xml.Attr{Name: xml.Name{Local: "i:type"}, Value: "BiddableAdGroupCriterion"}}
	e.EncodeToken(start)
	e.EncodeElement(s.AdGroupId, st("AdGroupId"))

	if s.Criterion != (ProductPartition{}) {
		e.Encode(s.Criterion)
		//marshalCriterion(s.Criterion, e)
	}

	if s.Id != 0 {
		e.EncodeElement(s.Id, st("Id"))
	}

	if s.Status != "" {
		e.EncodeElement(s.Status, st("Status"))
	}
	e.EncodeElement("BiddableAdGroupCriterion", st("Type"))
	if s.CriterionBid != (CriterionBid{}) {
		e.Encode(s.CriterionBid)
	}
	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

//func (s BiddableAdGroupCriterion) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error { }

type NegativeAdGroupCriterion struct {
	Id        int64
	AdGroupId int64
	Criterion ProductPartition
	Status    CriterionStatus
	Type      string
}

type AdGroupCriterions []interface{}

type ProductPartitionType string

const (
	Subdivision ProductPartitionType = "Subdivision"
	Unit        ProductPartitionType = "Unit"
)

func (c *CampaignService) GetAdGroupCriterionsByIds(adgroup int64) ([]BiddableAdGroupCriterion, error) {
	req := GetAdGroupCriterionsByIdsRequest{
		AdGroupId:     adgroup,
		CriterionType: "ProductPartition",
		NS:            "https://bingads.microsoft.com/CampaignManagement/v13",
	}
	resp, err := c.Session.SendRequest(req, c.Endpoint, "GetAdGroupCriterionsByIds")

	if err != nil {
		return nil, err
	}

	ret := GetAdGroupCriterionsByIdsResponse{}
	err = xml.Unmarshal(resp, &ret)

	if err != nil {
		return nil, err
	}

	return ret.AdGroupCriterions, err

}

//Action :: Add | Delete | Update
type AdGroupCriterionAction struct {
	Action           string
	AdGroupCriterion BiddableAdGroupCriterion
}

type ApplyProductPartitionActionsRequest struct {
	XMLName          xml.Name                 `xml:"ApplyProductPartitionActionsRequest"`
	NS               string                   `xml:"xmlns,attr"`
	CriterionActions []AdGroupCriterionAction `xml:"CriterionActions>AdGroupCriterionAction"`
}

//type AdGroupCriterionIds []int64
type Longs []int64

type ApplyProductPartitionActionsResponse struct {
	AdGroupCriterionIds Longs        `xml:"AdGroupCriterionIds>long"`
	PartialErrors       []BatchError `xml:"PartialErrors>BatchError"`
}

type productPartition struct {
	Type              string
	Condition         ProductCondition
	ParentCriterionId string
	PartitionType     string `xml:",omitempty"`
}

func (s *ProductPartition) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	p := productPartition{}
	dec.DecodeElement(&p, &start)
	s.Condition = p.Condition
	s.PartitionType = p.PartitionType
	s.Type = p.Type
	if p.ParentCriterionId != "" {
		n, _ := strconv.ParseInt(p.ParentCriterionId, 10, 64)
		s.ParentCriterionId = n
	}
	return nil
}

func (s *Longs) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	open := false
	skip := false
	if isNil, _ := findAttr(start.Attr, "nil"); isNil == "true" {
		skip = true
	}
	for token, err := dec.Token(); err == nil; token, err = dec.Token() {
		if err != nil {
			return err
		}

		switch next := token.(type) {
		case xml.CharData:
			i, err := strconv.ParseInt(string(next), 10, 64)
			if err != nil {
				return err
			}
			*s = append(*s, i)
			open = true
		case xml.EndElement:
			//instead of skipping nils, replace with 0s, so we can map partial success ids
			if open == false && !skip {
				*s = append(*s, 0)
			}
			open = false
		}
	}

	return nil
}

//TODO: should we handle mapping successful ids to actions
func (c *CampaignService) ApplyProductPartitionActions(actions []AdGroupCriterionAction) (*ApplyProductPartitionActionsResponse, error) {
	req := ApplyProductPartitionActionsRequest{
		NS:               "https://bingads.microsoft.com/CampaignManagement/v13",
		CriterionActions: actions,
	}
	resp, err := c.Session.SendRequest(req, c.Endpoint, "ApplyProductPartitionActions")

	if err != nil {
		return nil, err
	}

	ret := &ApplyProductPartitionActionsResponse{}
	err = xml.Unmarshal(resp, ret)
	return ret, err

}

func (agcs *AdGroupCriterions) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	adGroupCriterionType, err := findAttr(start.Attr, "type")
	if err != nil {
		return err
	}
	switch adGroupCriterionType {
	case "BiddableAdGroupCriterion":
		bagc := BiddableAdGroupCriterion{}
		err := dec.DecodeElement(&bagc, &start)
		if err != nil {
			return err
		}
		*agcs = append(*agcs, bagc)
	case "NegativeAdGroupCriterion":
		nagc := NegativeAdGroupCriterion{}
		err := dec.DecodeElement(&nagc, &start)
		if err != nil {
			return err
		}
		*agcs = append(*agcs, nagc)
	default:
		return fmt.Errorf("unknown AdGroupCriterion -> %#v", adGroupCriterionType)
	}
	return nil
}
