package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type StringClient string

func (s StringClient) SendRequest(_ interface{}, _, _ string) ([]byte, error) {
	return []byte(s), nil
}

func getTestClient() *CampaignService {
	client := &BingClient{
		CustomerAccountId: os.Getenv("BING_ACCOUNT_ID"),
		customerId:        os.Getenv("BING_CUSTOMER_ID"),
		username:          os.Getenv("BING_USERNAME"),
		password:          os.Getenv("BING_PASSWORD"),
		developerToken:    os.Getenv("BING_DEV_TOKEN"),
	}

	return &CampaignService{
		endpoint: "https://campaign.api.sandbox.bingads.microsoft.com/Api/Advertiser/CampaignManagement/v11/CampaignManagementService.svc",
		client:   client,
	}
}

func TestGetSandboxCampaigns(t *testing.T) {
	client := StringClient(`<GetCampaignsByAccountIdResponse xmlns="https://bingads.microsoft.com/CampaignManagement/v11"><Campaigns xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><Campaign><BiddingScheme i:type="ManualCpcBiddingScheme"><Type>ManualCpc</Type></BiddingScheme><BudgetType>DailyBudgetStandard</BudgetType><DailyBudget>25</DailyBudget><Description>dota2</Description><ForwardCompatibilityMap xmlns:a="http://schemas.datacontract.org/2004/07/System.Collections.Generic"/><Id>804002264</Id><Name>dota2</Name><NativeBidAdjustment i:nil="true"/><Status>Active</Status><TimeZone>EasternTimeUSCanada</TimeZone><TrackingUrlTemplate i:nil="true"/><UrlCustomParameters i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/><CampaignType>Shopping</CampaignType><Settings><Setting i:type="ShoppingSetting"><Type>ShoppingSetting</Type><LocalInventoryAdsEnabled i:nil="true"/><Priority>0</Priority><SalesCountryCode>US</SalesCountryCode><StoreId>1387210</StoreId></Setting></Settings><BudgetId i:nil="true"/><Languages i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/></Campaign></Campaigns></GetCampaignsByAccountIdResponse>`)

	svc := &CampaignService{
		endpoint: "https://campaign.api.sandbox.bingads.microsoft.com/Api/Advertiser/CampaignManagement/v11/CampaignManagementService.svc",
		client:   client,
	}
	//svc = getTestClient()
	res, err := svc.GetCampaignsByAccountId(os.Getenv("BING_ACCOUNT_ID"), Shopping)

	if err != nil {
		t.Fatal(err)
	}

	expected := []Campaign{{
		BiddingScheme: BiddingScheme{"ManualCpc"},
		BudgetType:    "DailyBudgetStandard",
		DailyBudget:   25,
		Description:   "dota2",
		Id:            804002264,
		Name:          "dota2",
		Status:        "Active",
		TimeZone:      "EasternTimeUSCanada",
	}}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

func TestGetSandBoxAdGroups(t *testing.T) {
	client := getTestClient()
	res, err := client.GetAdgroupsByCampaign(804002264)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}

func TestGetSandboxCriterion(t *testing.T) {
	client := getTestClient()
	res, err := client.GetAdGroupCriterionsByIds(1159984767305062)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}
