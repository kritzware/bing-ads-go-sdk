package main

import (
	"encoding/xml"

	"github.com/getsidecar/bing-ads-go-sdk/adgroup"
	"github.com/zak10/bing-ads-go-sdk/base"
)

type GetCampaignsByIdsRequest struct {
	XMLName      xml.Name `xml:"ns1:GetCampaignsByAccountIdRequest"`
	AccountId    string   `xml:"ns1:AccountId"`
	CampaignType string   `xml:"ns1:CampaignType"`
}

func main() {
	oauthToken := ""
	//wsdl := "https://campaign.api.bingads.microsoft.com/Api/Advertiser/CampaignManagement/V10/CampaignManagementService.svc?singleWsdl"
	CustomerAccountId := ""
	DeveloperToken := ""
	AdgroupId := ""
	client := base.NewBingClient(CustomerAccountId, "", DeveloperToken, oauthToken, "", "")
	adgroupService := adgroup.NewAdgroupService(client)
	adgroupService.GetAdGroupCriterionsByIdsRequest(AdgroupId, "ProductPartition")
}
