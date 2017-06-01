package bingads

import (
	"encoding/xml"
	"fmt"
)

type BiddingScheme struct {
	Type string
}

type Campaign struct {
	BiddingScheme BiddingScheme `xml:"https://bingads.microsoft.com/CampaignManagement/v11 BiddingScheme"`
	BudgetType    string        `xml:"https://bingads.microsoft.com/CampaignManagement/v11 BudgetType"`
	DailyBudget   float64       `xml:"https://bingads.microsoft.com/CampaignManagement/v11 DailyBudget"`
	Description   string        `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Description"`
	Id            int64         `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Id"`
	Name          string        `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Name"`
	//maybe parse into sql nullable int?
	//NativeBidAdjustment int     `xml:"https://bingads.microsoft.com/CampaignManagement/v11 NativeBidAdjustment"`
	Status   string `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Status"`
	TimeZone string `xml:"https://bingads.microsoft.com/CampaignManagement/v11 TimeZone"`
}

type CampaignType string

const (
	SearchAndContent CampaignType = "SearchAndContent"
	Shopping         CampaignType = "Shopping"
	DynamicSearchAds CampaignType = "DynamicSearchAds"
)

type GetCampaignsByIdsResponse struct {
}

type CampaignService struct {
	endpoint string
	client   Client
}

func NewCampaignService(client *BingClient) *CampaignService {
	return &CampaignService{
		endpoint: "https://campaign.api.bingads.microsoft.com/Api/Advertiser/CampaignManagement/v11/CampaignManagementService.svc",
		client:   client,
	}
}

type GetCampaignsByAccountIdRequest struct {
	XMLName      xml.Name     `xml:"GetCampaignsByAccountIdRequest"`
	AccountId    string       `xml:"AccountId"`
	CampaignType CampaignType `xml:"CampaignType"`
}

type GetCampaignsByAccountIdResponse struct {
	Campaigns []Campaign `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Campaigns>Campaign"`
}

func (c *CampaignService) GetCampaignsByAccountId(account string, campaignType CampaignType) ([]Campaign, error) {
	req := GetCampaignsByAccountIdRequest{
		CampaignType: campaignType,
		AccountId:    account,
	}

	resp, err := c.client.SendRequest(req, c.endpoint, "GetCampaignsByAccountId")

	if err != nil {
		return nil, err
	}

	campaignResponse := GetCampaignsByAccountIdResponse{}

	err = xml.Unmarshal(resp, &campaignResponse)
	return campaignResponse.Campaigns, err
}

type AdGroup struct {
	Id     int64  `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Id"`
	Name   string `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Name"`
	Status string `xml:"https://bingads.microsoft.com/CampaignManagement/v11 Status"`
}

type GetAdGroupsByCampaignIdRequest struct {
	XMLName    xml.Name `xml:"GetAdGroupsByCampaignIdRequest"`
	CampaignId int64    `xml:"CampaignId"`
}

type GetAdGroupsByCampaignIdResponse struct {
	AdGroups []AdGroup `xml:"https://bingads.microsoft.com/CampaignManagement/v11 AdGroups>AdGroup"`
}

func (c *CampaignService) GetAdgroupsByCampaign(campaign int64) ([]AdGroup, error) {

	req := GetAdGroupsByCampaignIdRequest{
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

func findAttr(xs []xml.Attr, name string) (string, error) {
	fmt.Println(xs)
	for _, x := range xs {
		fmt.Println(x.Name.Local)
		if x.Name.Local == name {
			return x.Value, nil
		}
	}

	return "", fmt.Errorf("attribute %s not found", name)

}
