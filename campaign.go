package bingads

import (
	"encoding/xml"
)

type BiddingScheme struct {
	Type     string
	TypeAttr string `xml:"i:type,attr"`
}

var ManualCpc = BiddingScheme{Type: "ManualCpc", TypeAttr: "ManualCpcBiddingScheme"}

type Campaign struct {
	BiddingScheme BiddingScheme `xml:"BiddingScheme"`
	BudgetType    BudgetType    `xml:"BudgetType"`
	BudgetId      string        `xml:",omitempty"`
	DailyBudget   float64       `xml:"DailyBudget"`
	Description   string        `xml:"Description"`
	Id            int64         `xml:"Id,omitempty"`
	Name          string        `xml:"Name"`
	//maybe parse into sql nullable int?
	//NativeBidAdjustment int     `xml:"https://bingads.microsoft.com/CampaignManagement/v11 NativeBidAdjustment"`
	Status       string             `xml:"Status"`
	TimeZone     string             `xml:"TimeZone"`
	CampaignType CampaignType       `xml:"CampaignType"`
	Settings     []CampaignSettings `xml:"Settings>Setting"`
}

//TODO: split into shoppingsetting + dynamicsearachadssetting
type CampaignSettings struct {
	Type     string
	TypeAttr string `xml:"i:type,attr"`

	LocalInventoryAdsEnabled string `xml:",omitempty"`
	Priority                 int
	SalesCountryCode         string
	StoreId                  int64

	DomainName string `xml:",omitempty"`
	Language   string `xml:",omitempty"`
}

type BudgetType string

const (
	DailyBudgetAccelerated BudgetType = "DailyBudgetAccelerated"
	DailyBudgetStandard               = "DailyBudgetStandard"
)

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

func NewCampaignService(client *Session) *CampaignService {
	return &CampaignService{
		endpoint: "https://campaign.api.bingads.microsoft.com/Api/Advertiser/CampaignManagement/v11/CampaignManagementService.svc",
		client:   client,
	}
}

type GetCampaignsByAccountIdRequest struct {
	XMLName      xml.Name     `xml:"GetCampaignsByAccountIdRequest"`
	NS           string       `xml:"xmlns,attr"`
	AccountId    string       `xml:"AccountId"`
	CampaignType CampaignType `xml:"CampaignType"`
}

type GetCampaignsByAccountIdResponse struct {
	Campaigns []Campaign `xml:"Campaigns>Campaign"`
}

func (c *CampaignService) GetCampaignsByAccountId(account string, campaignType CampaignType) ([]Campaign, error) {
	req := GetCampaignsByAccountIdRequest{
		NS:           "https://bingads.microsoft.com/CampaignManagement/v11",
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

func (c *CampaignService) AddCampaigns(account string, campaigns []Campaign) ([]int64, error) {
	req := AddCampaignsRequest{
		NS:        "https://bingads.microsoft.com/CampaignManagement/v11",
		Campaigns: campaigns,
		AccountId: account,
	}

	resp, err := c.client.SendRequest(req, c.endpoint, "AddCampaigns")

	if err != nil {
		return nil, err
	}

	ret := AddCampaignsResponse{}
	err = xml.Unmarshal(resp, &ret)
	return ret.CampaignIds, err
}

type AddCampaignsRequest struct {
	XMLName   xml.Name   `xml:"AddCampaignsRequest"`
	NS        string     `xml:"xmlns,attr"`
	AccountId string     `xml:"AccountId"`
	Campaigns []Campaign `xml:"Campaigns>Campaign"`
}
type AddCampaignsResponse struct {
	CampaignIds   []int64 `xml:"CampaignIds>long"`
	PartialErrors []BatchError
}
