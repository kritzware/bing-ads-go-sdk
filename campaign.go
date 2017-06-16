package bingads

import (
	"encoding/xml"
)

const (
	ManualCpc BiddingScheme = "ManualCpcBiddingScheme"
)

type BiddingScheme string

func (s BiddingScheme) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = ats("i:type", string(s))
	e.EncodeElement(string(s), st("Type"))
	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

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

func (s CampaignSettings) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = ats("i:type", s.Type)
	return e.EncodeElement(s, start)
}

//TODO: split into shoppingsetting + dynamicsearachadssetting
//SalesCountryCode: US
type CampaignSettings struct {
	Type string

	LocalInventoryAdsEnabled string `xml:",omitempty"`
	Priority                 int
	SalesCountryCode         string
	StoreId                  int64

	DomainName string `xml:",omitempty"`
	Language   string `xml:",omitempty"`
}

//TODO: maybe leave as string
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

type CampaignService struct {
	Endpoint string
	Client   Client
}

func NewCampaignService(client *Session) *CampaignService {
	return &CampaignService{
		Endpoint: "https://campaign.api.bingads.microsoft.com/Api/Advertiser/CampaignManagement/v11/CampaignManagementService.svc",
		Client:   client,
	}
}

type GetCampaignsByAccountIdRequest struct {
	XMLName      xml.Name     `xml:"GetCampaignsByAccountIdRequest"`
	NS           string       `xml:"xmlns,attr"`
	AccountId    string        `xml:"AccountId"`
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

	resp, err := c.Client.SendRequest(req, c.Endpoint, "GetCampaignsByAccountId")

	if err != nil {
		return nil, err
	}

	campaignResponse := GetCampaignsByAccountIdResponse{}

	err = xml.Unmarshal(resp, &campaignResponse)
	return campaignResponse.Campaigns, err
}

func (c *CampaignService) AddCampaigns(account string, campaigns []Campaign) (*AddCampaignsResponse, error) {
	req := AddCampaignsRequest{
		NS:        "https://bingads.microsoft.com/CampaignManagement/v11",
		Campaigns: campaigns,
		AccountId: account,
	}

	resp, err := c.Client.SendRequest(req, c.Endpoint, "AddCampaigns")

	if err != nil {
		return nil, err
	}

	ret := &AddCampaignsResponse{}
	err = xml.Unmarshal(resp, ret)
	return ret, err
}

type AddCampaignsRequest struct {
	XMLName   xml.Name   `xml:"AddCampaignsRequest"`
	NS        string     `xml:"xmlns,attr"`
	AccountId string      `xml:"AccountId"`
	Campaigns []Campaign `xml:"Campaigns>Campaign"`
}
type AddCampaignsResponse struct {
	CampaignIds   []int64      `xml:"CampaignIds>long"`
	PartialErrors []BatchError `xml:"PartialErrors>BatchError"`
}
