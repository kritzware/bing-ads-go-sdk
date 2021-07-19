package bingads

import (
	"encoding/xml"
)

//BudgetType :: DailyBudgetStandard
type Campaign struct {
	BiddingScheme BiddingScheme `xml:"BiddingScheme"`
	BudgetType    string        `xml:"BudgetType"`
	BudgetId      string        `xml:",omitempty"`
	DailyBudget   float64       `xml:"DailyBudget"`
	Description   string        `xml:"Description"`
	Id            int64         `xml:"Id,omitempty"`
	Name          string        `xml:"Name"`
	//maybe parse into sql nullable int?
	//NativeBidAdjustment int     `xml:"https://bingads.microsoft.com/CampaignManagement/v13 NativeBidAdjustment"`
	Status       string             `xml:"Status"`
	TimeZone     string             `xml:"TimeZone"`
	CampaignType CampaignType       `xml:"CampaignType"`
	Settings     []CampaignSettings `xml:"Settings>Setting"`
}

//BiddingScheme :: ManualCpc
type BiddingScheme struct {
	Type string
}

func (s BiddingScheme) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = ats("i:type", s.Type+"BiddingScheme")
	e.EncodeToken(start)
	e.EncodeElement(s.Type, st("Type"))
	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

func (s CampaignSettings) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = ats("i:type", s.Type)
	e.EncodeToken(start)

	e.EncodeElement(s.Type, st("Type"))
	if s.LocalInventoryAdsEnabled != "" {
		e.EncodeElement(s.LocalInventoryAdsEnabled, st("LocalInventoryAdsEnabled"))
	}
	e.EncodeElement(s.Priority, st("Priority"))
	e.EncodeElement(s.SalesCountryCode, st("SalesCountryCode"))
	e.EncodeElement(s.StoreId, st("StoreId"))

	if s.Language != "" {
		e.EncodeElement(s.Language, st("Language"))
	}
	if s.DomainName != "" {
		e.EncodeElement(s.DomainName, st("DomainName"))
	}

	e.EncodeToken(xml.EndElement{start.Name})
	return nil
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
	Session  *Session
}

func NewCampaignService(session *Session) *CampaignService {
	return &CampaignService{
		Endpoint: "https://campaign.api.bingads.microsoft.com/Api/Advertiser/CampaignManagement/v13/CampaignManagementService.svc",
		Session:  session,
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

func (c *CampaignService) GetCampaignsByAccountId(campaignType CampaignType) ([]Campaign, error) {
	req := GetCampaignsByAccountIdRequest{
		NS:           "https://bingads.microsoft.com/CampaignManagement/v13",
		CampaignType: campaignType,
		AccountId:    c.Session.AccountId,
	}

	resp, err := c.Session.SendRequest(req, c.Endpoint, "GetCampaignsByAccountId")

	if err != nil {
		return nil, err
	}

	campaignResponse := GetCampaignsByAccountIdResponse{}

	err = xml.Unmarshal(resp, &campaignResponse)
	return campaignResponse.Campaigns, err
}

func (c *CampaignService) AddCampaigns(campaigns []Campaign) (*AddCampaignsResponse, error) {
	req := AddCampaignsRequest{
		NS:        "https://bingads.microsoft.com/CampaignManagement/v13",
		Campaigns: campaigns,
		AccountId: c.Session.AccountId,
	}

	resp, err := c.Session.SendRequest(req, c.Endpoint, "AddCampaigns")

	if err != nil {
		return nil, err
	}

	ret := &addCampaignsResponse{}
	if err := xml.Unmarshal(resp, ret); err != nil {
		return nil, err
	}

	return &AddCampaignsResponse{
		CampaignIds:   []int64(ret.CampaignIds),
		PartialErrors: ret.PartialErrors,
	}, nil
}

type AddCampaignsRequest struct {
	XMLName   xml.Name   `xml:"AddCampaignsRequest"`
	NS        string     `xml:"xmlns,attr"`
	AccountId string     `xml:"AccountId"`
	Campaigns []Campaign `xml:"Campaigns>Campaign"`
}
type addCampaignsResponse struct {
	CampaignIds   Longs        `xml:"CampaignIds>long"`
	PartialErrors []BatchError `xml:"PartialErrors>BatchError"`
}

type AddCampaignsResponse struct {
	CampaignIds   []int64      `xml:"CampaignIds>long"`
	PartialErrors []BatchError `xml:"PartialErrors>BatchError"`
}

func (c *CampaignService) DeleteCampaigns(campaigns []int64) (*DeleteCampaignsResponse, error) {
	req := DeleteCampaignsRequest{
		NS:          "https://bingads.microsoft.com/CampaignManagement/v13",
		CampaignIds: campaigns,
		AccountId:   c.Session.AccountId,
	}

	resp, err := c.Session.SendRequest(req, c.Endpoint, "DeleteCampaigns")

	if err != nil {
		return nil, err
	}

	ret := &DeleteCampaignsResponse{}
	if err := xml.Unmarshal(resp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

type DeleteCampaignsRequest struct {
	XMLName     xml.Name `xml:"DeleteCampaignsRequest"`
	NS          string   `xml:"xmlns,attr"`
	AccountId   string
	CampaignIds []int64 `xml:"CampaignIds>long"`
}

func (s DeleteCampaignsRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeToken(st("DeleteCampaignsRequest", "xmlns", s.NS))
	e.EncodeElement(s.AccountId, st("AccountId"))

	e.EncodeToken(st("CampaignIds", "xmlns:a1", "http://schemas.microsoft.com/2003/10/Serialization/Arrays"))
	for i := 0; i < len(s.CampaignIds); i++ {
		e.EncodeElement(s.CampaignIds[i], st("a1:long"))
	}
	e.EncodeToken(xml.EndElement{xml.Name{Local: "CampaignIds"}})

	e.EncodeToken(xml.EndElement{xml.Name{Local: "DeleteCampaignsRequest"}})
	return nil
}

type DeleteCampaignsResponse struct {
	PartialErrors []BatchError `xml:"PartialErrors>BatchError"`
}
