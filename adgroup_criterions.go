package bingads

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

type ProductCondition struct {
	Attribute string
	Operand   string
}

//TODO: derived types
//https://msdn.microsoft.com/en-us/library/bing-ads-campaign-management-getcampaigncriterionsbyids.aspx
type Criterion struct {
	Type      string
	TypeAttr  string            `xml:"i:type,attr"`
	Condition *ProductCondition //`xml:"Conditions>ProductCondition"`
	//should be nullable int64
	ParentCriterionId string               `xml:",omitempty"`
	PartitionType     ProductPartitionType `xml:",omitempty"`
}

type CriterionType string

const (
	ProductPartition CriterionType = "ProductPartition"
	ProductScope                   = "ProductScope"
)

type GetAdGroupCriterionsByIdsRequest struct {
	XMLName xml.Name `xml:"GetAdGroupCriterionsByIdsRequest"`

	NS            string        `xml:"xmlns,attr"`
	AdGroupId     int64         `xml:"AdGroupId"`
	CriterionType CriterionType `xml:"CriterionType"`
}

type GetAdGroupCriterionsByIdsResponse struct {
	AdGroupCriterions []BiddableAdGroupCriterion `xml:"https://bingads.microsoft.com/CampaignManagement/v11 AdGroupCriterions>AdGroupCriterion"`
}

type AdGroupCriterion struct {
	Id           int64 `xml:",omitempty"`
	AdGroupId    int64
	Criterion    Criterion
	Status       CriterionStatus
	Type         string
	CriterionBid CriterionBid
}
type BiddableAdGroupCriterion struct {
	TypeAttr     string `xml:"i:type,attr"`
	AdGroupId    int64
	Criterion    Criterion
	Id           int64 `xml:",omitempty"`
	Status       CriterionStatus
	Type         string
	CriterionBid CriterionBid
}

type NegativeAdGroupCriterion struct {
	Id        int64
	AdGroupId int64
	Criterion Criterion
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
		CriterionType: ProductPartition,
		NS:            "https://bingads.microsoft.com/CampaignManagement/v11",
	}
	resp, err := c.client.SendRequest(req, c.endpoint, "GetAdGroupCriterionsByIds")

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

type AdGroupCriterionAction struct {
	Action           string
	AdGroupCriterion BiddableAdGroupCriterion
}

type ApplyProductPartitionActionsRequest struct {
	XMLName          xml.Name                 `xml:"ApplyProductPartitionActionsRequest"`
	NS               string                   `xml:"xmlns,attr"`
	CriterionActions []AdGroupCriterionAction `xml:"CriterionActions>AdGroupCriterionAction"`
}

type AdGroupCriterionIds []int64

type ApplyProductPartitionActionsResponse struct {
	AdGroupCriterionIds AdGroupCriterionIds `xml:"AdGroupCriterionIds>long"`
}

func (s *AdGroupCriterionIds) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
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
		}
	}

	return nil
}

func (c *CampaignService) ApplyProductPartitionActions(actions []AdGroupCriterionAction) ([]int64, error) {
	req := ApplyProductPartitionActionsRequest{
		NS:               "https://bingads.microsoft.com/CampaignManagement/v11",
		CriterionActions: actions,
	}
	resp, err := c.client.SendRequest(req, c.endpoint, "ApplyProductPartitionActions")

	if err != nil {
		return nil, err
	}

	ret := ApplyProductPartitionActionsResponse{}
	err = xml.Unmarshal(resp, &ret)
	return ret.AdGroupCriterionIds, err

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
