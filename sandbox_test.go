package main

import (
	"fmt"
	"os"
	"testing"
)

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
	client := getTestClient()

	res, err := client.GetCampaignsByAccountId(Shopping)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
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
