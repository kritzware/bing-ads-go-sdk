package main

import (
	"encoding/xml"
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
		endpoint: "https://campaign.api.bingads.microsoft.com/Api/Advertiser/CampaignManagement/V10/CampaignManagementService.svc?singleWsdl",
		client:   client,
	}
}

type GetCampaignsByAccountIdRequest struct {
	XMLName      xml.Name     `xml:"ns1:GetCampaignsByAccountIdRequest"`
	AccountId    string       `xml:"ns1:AccountId"`
	CampaignType CampaignType `xml:"ns1:CampaignType"`
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
	XMLName    xml.Name `xml:"ns1:GetAdGroupsByCampaignIdRequest"`
	CampaignId int64    `xml:"ns1:CampaignId"`
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

type CriterionType string

const (
	ProductPartition CriterionType = "ProductPartition"
)

type GetAdGroupCriterionsByIdsRequest struct {
	XMLName       xml.Name      `xml:"ns1:GetAdGroupCriterionsByIdsRequest"`
	AdGroupId     int64         `xml:"ns1:AdGroupId"`
	CriterionType CriterionType `xml:"ns1:CriterionType"`
}

type GetAdGroupCriterionsByIdsResponse struct {
	AdGroupCriterions []AdGroupCriterion `xml:"https://bingads.microsoft.com/CampaignManagement/v11 AdGroupCriterions>AdGroupCriterion"`
}

type AdGroupCriterionStatus string

const (
	Active  AdGroupCriterionStatus = "Active"
	Paused                         = "Paused"
	Deleted                        = "Deleted"
)

type AdGroupCriterion struct {
	Id        int64
	AdGroupId int64
	Criterion Criterion
	Status    AdGroupCriterionStatus
	Type      string
}

type ProductPartitionType string

const (
	Subdivision ProductPartitionType = "Subdivision"
	Unit        ProductPartitionType = "Unit"
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

func (c *CampaignService) GetAdGroupCriterionsByIds(adgroup int64) ([]AdGroupCriterion, error) {
	req := GetAdGroupCriterionsByIdsRequest{
		AdGroupId:     adgroup,
		CriterionType: ProductPartition,
	}
	resp, err := c.client.SendRequest(req, c.endpoint, "GetAdGroupCriterionsByIds")

	if err != nil {
		return nil, err
	}

	ret := GetAdGroupCriterionsByIdsResponse{}
	err = xml.Unmarshal(resp, &ret)
	return ret.AdGroupCriterions, err

}
