package main

import ()

/*
type AdGroupCriterion struct {
	XMLName   xml.Name  `xml:"AdGroupCriterion"`
	Type      string    `xml:"type,attr"`
	AdGroupId string    `xml:"AdGroupId"`
	Criterion Criterion `xml:"Criterion"`
	Id        string    `xml:"Id"`
	Status    string    `xml:"status"`
}

type Criterion struct {
	Type              string    `xml:"type,attr"`
	Condition         Condition `xml:"Condition"`
	ParentCriterionId string    `xml:"ParentCriterionId"`
	PartitionType     string    `xml:"PartitionType"`
}

type Condition struct {
	Attribute string `xml:"Attribute"`
	Operand   string `xml:"Operand"`
}

type AdgroupService struct {
	endpoint string
	client   *BingClient
}

func NewAdgroupService(client *BingClient) *AdgroupService {
	return &AdgroupService{
		endpoint: "https://campaign.api.bingads.microsoft.com/Api/Advertiser/CampaignManagement/V10/CampaignManagementService.svc?singleWsdl",
		client:   client,
	}
}

func (a *AdgroupService) GetAdGroupCriterionsByIdsRequest(adgroup string, criterionType string) {
	req := a.client.MakeRequest(GetAdGroupCriterionsByIdsRequest{
		AccountId:     a.client.CustomerAccountId,
		AdGroupId:     adgroup,
		CriterionType: criterionType,
	})
	resp, _ := a.client.SendRequest(req, a.endpoint, "GetAdGroupCriterionsByIds")

	fmt.Println(string(resp))

	var r GetAdGroupCriterionsByIdsResponse
	err := xml.Unmarshal(resp, &r)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", r)
}
*/
