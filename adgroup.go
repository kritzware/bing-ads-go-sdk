package bingads

import (
	"encoding/xml"
)

type Bid struct {
	Amount float64
}

type Date struct {
	Day   int
	Month int
	Year  int
}

//OwnedAndOperatedOnly
//apparently required fields needs to be first or we get validation errors
type AdGroup struct {
	Id                          int64  `xml:",omitempty"`
	AdDistribution              string `xml:",omitempty"`
	Language                    string `xml:",omitempty"`
	Name                        string
	Status                      string `xml:",omitempty"`
	StartDate                   *Date
	EndDate                     *Date
	BiddingScheme               BiddingScheme
	NativeBidAdjustment         string `xml:",omitempty"`
	Network                     string `xml:",omitempty"`
	PricingModel                string `xml:",omitempty"`
	RemarketingTargetingSetting string `xml:",omitempty"`
	SearchBid                   *Bid
	TrackingUrlTemplate         string `xml:",omitempty"`
}

type GetAdGroupsByCampaignIdRequest struct {
	XMLName    xml.Name `xml:"GetAdGroupsByCampaignIdRequest"`
	NS         string   `xml:"xmlns,attr"`
	CampaignId int64
}

type GetAdGroupsByCampaignIdResponse struct {
	AdGroups []AdGroup `xml:"https://bingads.microsoft.com/CampaignManagement/v12 AdGroups>AdGroup"`
}

func (c *CampaignService) GetAdgroupsByCampaign(campaign int64) ([]AdGroup, error) {
	req := GetAdGroupsByCampaignIdRequest{
		NS:         "https://bingads.microsoft.com/CampaignManagement/v12",
		CampaignId: campaign,
	}

	resp, err := c.Session.SendRequest(req, c.Endpoint, "GetAdGroupsByCampaignId")

	if err != nil {
		return nil, err
	}

	ret := GetAdGroupsByCampaignIdResponse{}
	err = xml.Unmarshal(resp, &ret)
	return ret.AdGroups, err
}

func (c *CampaignService) AddAdGroups(campaign int64, adgroups []AdGroup) (*AddAdGroupsResponse, error) {
	req := AddAdGroupsRequest{
		NS:         "https://bingads.microsoft.com/CampaignManagement/v12",
		CampaignId: campaign,
		AdGroups:   adgroups,
	}

	resp, err := c.Session.SendRequest(req, c.Endpoint, "AddAdGroups")

	if err != nil {
		return nil, err
	}

	ret := &AddAdGroupsResponse{}
	err = xml.Unmarshal(resp, ret)
	return ret, err

}

type AddAdGroupsRequest struct {
	XMLName    xml.Name `xml:"AddAdGroupsRequest"`
	NS         string   `xml:"xmlns,attr"`
	CampaignId int64
	AdGroups   []AdGroup `xml:"AdGroups>AdGroup"`
}

type AddAdGroupsResponse struct {
	AdGroupIds    Longs        `xml:"AdGroupIds>long"`
	PartialErrors []BatchError `xml:"PartialErrors>BatchError"`
}
