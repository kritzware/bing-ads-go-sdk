package bingads

import (
	"encoding/xml"
)

type AdGroup struct {
	Id     int64  `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Id"`
	Name   string `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Name"`
	Status string `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Status"`
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
