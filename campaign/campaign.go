package campaign

import (
	"encoding/xml"

	"github.com/getsidecar/bing-ads-go-sdk/base"
)

type GetCampaignsByIdsRequest struct {
	XMLName      xml.Name `xml:"ns1:GetCampaignsByAccountIdRequest"`
	AccountId    string   `xml:"ns1:AccountId"`
	CampaignType string   `xml:"ns1:CampaignType"`
}

type GetCampaignsByIdsResponse struct {
}

type CampaignService struct {
	endpoint string
	client   *base.BingClient
}

func NewCampaignService(client *base.BingClient) *CampaignService {
	return &CampaignService{
		endpoint: "https://campaign.api.bingads.microsoft.com/Api/Advertiser/CampaignManagement/V10/CampaignManagementService.svc?singleWsdl",
		client:   client,
	}
}

func (c *CampaignService) GetCampaignsByAccountIds(campaignType string) (*GetCampaignsByIdsResponse, error) {
	req := c.client.MakeRequest(GetCampaignsByIdsRequest{
		CampaignType: campaignType,
		AccountId:    c.client.CustomerAccountId,
	})

	resp, err := c.client.SendRequest(req, c.endpoint, "GetCampaignsByAccountIds")

	if err != nil {
		return nil, err
	}

	var campaignResponse GetCampaignsByIdsResponse
	err = xml.Unmarshal(resp, &campaignResponse)

	return &campaignResponse, err
}
