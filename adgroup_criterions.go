package bingads

import (
	"encoding/xml"
	"fmt"
)

type ProductCondition struct {
	Attribute string
	Operand   string
}

type Criterion struct {
	Type      string
	Condition ProductCondition
	//should be nullable int64
	ParentCriterionId string
	PartitionType     ProductPartitionType
}

type CriterionType string

const (
	ProductPartition CriterionType = "ProductPartition"
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

type AdGroupCriterionStatus string

const (
	Active  AdGroupCriterionStatus = "Active"
	Paused                         = "Paused"
	Deleted                        = "Deleted"
)

//either fixed or multiplier
type CriterionBid struct {
	Type   string
	Amount float64
}

type AdGroupCriterion struct {
	Id           int64
	AdGroupId    int64
	Criterion    Criterion
	Status       AdGroupCriterionStatus
	Type         string
	CriterionBid CriterionBid
}
type BiddableAdGroupCriterion struct {
	Id           int64
	AdGroupId    int64
	Criterion    Criterion
	Status       AdGroupCriterionStatus
	Type         string
	CriterionBid CriterionBid
}

type NegativeAdGroupCriterion struct {
	Id        int64
	AdGroupId int64
	Criterion Criterion
	Status    AdGroupCriterionStatus
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
	return ret.AdGroupCriterions, err

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
