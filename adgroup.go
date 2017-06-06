package bingads

import (
	"encoding/xml"
)

type Bid struct {
	Amount float64
}

type Date struct {
	Month int
	Date  int
	Year  int
}

//OwnedAndOperatedOnly
//apparently required fields needs to be first or we get validation errors
type AdGroup struct {
	Id                          int64  `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Id,omitempty"`
	AdDistribution              string `xml:",omitempty"`
	Language                    string `xml:",omitempty"`
	Name                        string `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Name"`
	Status                      string `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Status,omitempty"`
	StartDate                   *Date
	EndDate                     *Date
	BiddingScheme               *BiddingScheme
	NativeBidAdjustment         int    `xml:",omitempty"`
	Network                     string `xml:",omitempty"`
	PricingModel                string `xml:",omitempty"`
	RemarketingTargetingSetting string `xml:",omitempty"`
	SearchBid                   *Bid
	TrackingUrlTemplate         string `xml:",omitempty"`
}

type GetAdGroupsByCampaignIdRequest struct {
	XMLName    xml.Name `xml:"GetAdGroupsByCampaignIdRequest"`
	CampaignId int64    `xml:"CampaignId"`
	NS         string   `xml:"xmlns,attr"`
}

type GetAdGroupsByCampaignIdResponse struct {
	AdGroups []AdGroup `xml:"https://bingads.microsoft.com/CampaignManagement/v11 AdGroups>AdGroup"`
}

func (c *CampaignService) GetAdgroupsByCampaign(campaign int64) ([]AdGroup, error) {
	req := GetAdGroupsByCampaignIdRequest{
		NS:         "https://bingads.microsoft.com/CampaignManagement/v11",
		CampaignId: campaign,
	}

	resp, err := c.client.SendRequest(req, c.endpoint, "GetAdGroupsByCampaignId")

	if err != nil {
		return nil, err
	}

	ret := GetAdGroupsByCampaignIdResponse{}
	err = xml.Unmarshal(resp, &ret)
	return ret.AdGroups, err
}

func (c *CampaignService) AddAdGroups(campaign int64, adgroups []AdGroup) ([]int64, error) {
	req := AddAdGroupsRequest{
		NS:         "https://bingads.microsoft.com/CampaignManagement/v11",
		CampaignId: campaign,
		AdGroups:   adgroups,
	}

	resp, err := c.client.SendRequest(req, c.endpoint, "AddAdGroups")

	if err != nil {
		return nil, err
	}

	ret := AddAdGroupsResponse{}
	err = xml.Unmarshal(resp, &ret)
	return ret.AdGroupIds, err

}

type AddAdGroupsRequest struct {
	XMLName    xml.Name `xml:"AddAdGroupsRequest"`
	NS         string   `xml:"xmlns,attr"`
	CampaignId int64
	AdGroups   []AdGroup `xml:"AdGroups>AdGroup"`
}

type AddAdGroupsResponse struct {
	AdGroupIds []int64 `xml:"AdGroupIds>long"`
}
