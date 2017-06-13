package bingads

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"testing"
)

type StringClient string

func (s StringClient) SendRequest(_ interface{}, _, _ string) ([]byte, error) {
	return []byte(s), nil
}

func getIds(xs []NegativeKeywordList) []int64 {
	r := make([]int64, len(xs))

	for i := 0; i < len(xs); i++ {
		r[i] = xs[i].Id
	}
	return r
}

func TestSandboxGetSharedEntities(t *testing.T) {
	svc := getTestClient()

	existing, err := svc.GetSharedEntitiesByAccountId("NegativeKeywordList")
	if err != nil {
		t.Fatal(existing)
	}

	a, err := svc.GetSharedEntityAssociationsByEntityIds([]int64{804004280})

	if err != nil {
		t.Fatal(err)
	}

	if len(a.Associations) > 0 {
		partials, err := svc.DeleteSharedEntityAssociations(a.Associations)
		if err != nil {
			t.Fatal(err)
		}

		if len(partials) > 0 {
			t.Fatal(partials)
		}
	}

	partials, err := svc.DeleteSharedEntities(existing)
	if err != nil {
		t.Fatal(err)
	}

	if len(partials) > 0 {
		t.Fatal(partials)
	}

	items := []NegativeKeyword{{
		//Id:        63001000817,
		MatchType: "Phrase",
		Text:      "asdf-1",
	}}

	added, err := svc.AddSharedEntity(&NegativeKeywordList{
		Name: "asdf negative keyword list",
	}, items)
	if err != nil {
		t.Fatal(err)
	}

	items2 := []NegativeKeyword{{
		//Id:        63001000817,
		MatchType: "Phrase",
		Text:      "asdf-2",
	}}

	sa, err := svc.AddListItemsToSharedList(&NegativeKeywordList{Id: added.SharedEntityId}, items2)
	if err != nil {
		t.Fatal(err)
	}

	if len(sa.PartialErrors) > 0 {
		t.Fatal(sa.PartialErrors)
	}

	err = svc.SetSharedEntityAssociations([]SharedEntityAssociation{
		{
			EntityId:         804004280,
			EntityType:       "Campaign",
			SharedEntityId:   added.SharedEntityId,
			SharedEntityType: "NegativeKeywordList",
		},
	})

	if err != nil {
		t.Error(err)
	}
}

func TestUnmarshalResponse(t *testing.T) {

	s := StringClient(`<ApplyProductPartitionActionsResponse xmlns="https://bingads.microsoft.com/CampaignManagement/v11"><AdGroupCriterionIds xmlns:a="http://schemas.datacontract.org/2004/07/System" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><a:long>1</a:long><a:long>2</a:long><a:long>3</a:long></AdGroupCriterionIds><PartialErrors xmlns:i="http://www.w3.org/2001/XMLSchema-instance"/></ApplyProductPartitionActionsResponse>`)
	svc := &CampaignService{
		Endpoint: "https://campaign.api.sandbox.bingads.microsoft.com/Api/Advertiser/CampaignManagement/v11/CampaignManagementService.svc",
		Client:   s,
	}

	res, err := svc.ApplyProductPartitionActions(nil)

	if err != nil {
		t.Error(err)
	}
	if len(res.AdGroupCriterionIds) != 3 {
		t.Errorf("expected 3 ids, got %d", len(res.AdGroupCriterionIds))
	}

	s = StringClient(`<ApplyProductPartitionActionsResponse xmlns="https://bingads.microsoft.com/CampaignManagement/v11"><AdGroupCriterionIds xmlns:a="http://schemas.datacontract.org/2004/07/System" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><a:long i:nil="true"/></AdGroupCriterionIds><PartialErrors xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><BatchError><Code>4129</Code><Details i:nil="true"/><ErrorCode>CampaignServiceDuplicateProductConditions</ErrorCode><FieldPath>ProductConditions</FieldPath><ForwardCompatibilityMap i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/System.Collections.Generic"/><Index>0</Index><Message>Children of a product partition node cannot contain duplicate product conditions.</Message><Type>BatchError</Type></BatchError></PartialErrors></ApplyProductPartitionActionsResponse>`)
	svc = &CampaignService{
		Endpoint: "https://campaign.api.sandbox.bingads.microsoft.com/Api/Advertiser/CampaignManagement/v11/CampaignManagementService.svc",
		Client:   s,
	}

	res, err = svc.ApplyProductPartitionActions(nil)
	if err != nil {
		t.Error(err)
	}

	if len(res.AdGroupCriterionIds) != 0 {
		t.Errorf("expected no ids")
	}
}

func TestSandboxApplyProductPartitionActions(t *testing.T) {
	svc := getTestClient()

	existing, err := svc.GetAdGroupCriterionsByIds(1167681348701053)
	if err != nil {
		t.Fatal(err)
	}

	if len(existing) == 0 {
		t.Fatalf("expected all products root group to exist")
	}

	fmt.Println(existing)

	parentid := fmt.Sprintf("%d", existing[0].Id)

	tomodifyid, err := func() (int64, error) {
		for _, e := range existing {
			if e.Criterion.Condition.Attribute == "str" {
				return e.Id, nil
			}
		}
		return 0, fmt.Errorf("expected existing str partition")
	}()

	if err != nil {
		t.Fatal(err)
	}

	cleanup := func(id int64) {
		toremove := BiddableAdGroupCriterion{
			AdGroupId: 1167681348701053,
			Id:        id,
		}

		actions := []AdGroupCriterionAction{{"Delete", toremove}}
		res, err := svc.ApplyProductPartitionActions(actions)

		if err != nil {
			t.Fatal(err)
		}

		if len(res.AdGroupCriterionIds) != 1 {
			t.Fatalf("expected 1 delete, got %d", len(res.AdGroupCriterionIds))
		}
	}

	a := BiddableAdGroupCriterion{
		AdGroupId: 1167681348701053,
		Criterion: Criterion{
			Condition:         ProductCondition{"int", "ProductType1"},
			ParentCriterionId: parentid,
			PartitionType:     "Unit",
		},
		Status: "Active",
		CriterionBid: CriterionBid{
			Amount: 0.35,
		},
	}

	b := BiddableAdGroupCriterion{
		AdGroupId: 1159984767306214,
		Criterion: Criterion{
			Condition:         ProductCondition{"str", "ProductType1"},
			ParentCriterionId: parentid,
			PartitionType:     "Unit",
		},
		Status: "Active",
		CriterionBid: CriterionBid{
			Amount: 0.35,
		},
	}

	c := BiddableAdGroupCriterion{
		AdGroupId: 1167681348701053,
		Id:        tomodifyid,
		Criterion: Criterion{
			Condition:         ProductCondition{"str", "ProductType1"},
			ParentCriterionId: parentid,
			PartitionType:     "Unit",
		},
		Status: "Active",
		CriterionBid: CriterionBid{
			Amount: 0.35,
		},
	}

	expectedPartialFailures := map[int]bool{
		0: true,
		2: true,
	}

	actions := []AdGroupCriterionAction{
		{"Add", b},
		{"Add", a},
		{"Add", b},
		{"Update", c},
	}

	applied, err := svc.ApplyProductPartitionActions(actions)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(applied)

	expectedSuccessful := len(actions) - len(expectedPartialFailures)

	if len(applied.AdGroupCriterionIds) != expectedSuccessful {
		t.Errorf("expected %d created, got %d", expectedSuccessful, len(applied.AdGroupCriterionIds))
	}

	for _, failure := range applied.PartialErrors {
		if !expectedPartialFailures[failure.Index] {
			t.Errorf("unexpected partial failure at index %d", failure.Index)
		}
	}

	cleanup(applied.AdGroupCriterionIds[0])
}

func TestUnmarshalCampaignScope(t *testing.T) {
	s := &CampaignService{
		Client: StringClient(`<GetCampaignCriterionsByIdsResponse xmlns="https://bingads.microsoft.com/CampaignManagement/v11"><CampaignCriterions xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><CampaignCriterion i:type="BiddableCampaignCriterion"><CampaignId>283025743</CampaignId><Criterion i:type="ProductScope"><Type>ProductScope</Type><Conditions><ProductCondition><Attribute>top_brand</Attribute><Operand>CustomLabel0</Operand></ProductCondition></Conditions></Criterion><ForwardCompatibilityMap i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/System.Collections.Generic"/><Id>840001008</Id><Status i:nil="true"/><Type>BiddableCampaignCriterion</Type><CriterionBid i:nil="true"/></CampaignCriterion></CampaignCriterions><PartialErrors i:nil="true" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"/></GetCampaignCriterionsByIdsResponse>`),
	}

	res, err := s.GetCampaignCriterionsByIds(123)

	if err != nil {
		t.Fatal(err)
	}

	//	{283025743 {ProductScope [{top_brand CustomLabel0}]  } 840001008  BiddableCampaignCriterion}]
	expected := []CampaignCriterion{{
		CampaignId: 283025743,
		Type:       "BiddableCampaignCriterion",
		Id:         840001008,
		Criterion: Criterion{
			Type:      ProductScope,
			Condition: ProductCondition{"top_brand", "CustomLabel0"},
		},
	}}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

func getTestClient() *CampaignService {
	client := &Session{
		AccountId:      os.Getenv("BING_ACCOUNT_ID"),
		CustomerId:     os.Getenv("BING_CUSTOMER_ID"),
		Username:       os.Getenv("BING_USERNAME"),
		Password:       os.Getenv("BING_PASSWORD"),
		DeveloperToken: os.Getenv("BING_DEV_TOKEN"),
		HTTPClient:     &http.Client{},
	}

	return &CampaignService{
		Endpoint: "https://campaign.api.sandbox.bingads.microsoft.com/Api/Advertiser/CampaignManagement/v11/CampaignManagementService.svc",
		Client:   client,
	}
}

func TestAddAdGroupSandbox(t *testing.T) {

	svc := getTestClient()

	toadd := []AdGroup{
		{
			Name:           "new adgroup1",
			Language:       "English",
			Network:        "OwnedAndOperatedOnly",
			AdDistribution: "Search",
			BiddingScheme:  BiddingScheme("ManualCpcBiddingScheme"),
		},
	}

	res, err := svc.AddAdGroups(804004280, toadd)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

func TestAddCampaignCriterions(t *testing.T) {
	svc := getTestClient()

	/*
		camps, err := svc.GetCampaignsByAccountId(os.Getenv("BING_ACCOUNT_ID"), Shopping)

		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(camps)

		crits, err := svc.GetCampaignCriterionsByIds(804004280)

		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(crits)
	*/

	cs := []CampaignCriterion{
		{
			Type:       "BiddableCampaignCriterion",
			CampaignId: 804004280,
			Nil:        "true",
			Criterion: Criterion{
				Type:      "ProductScope",
				Condition: ProductCondition{"top_brand", "CustomLabel0"},
			},
		},
	}

	res, err := svc.AddCampaignCriterions(ProductScope, cs)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

func TestAddCampaigns(t *testing.T) {

	svc := getTestClient()

	toadd := []Campaign{{
		BiddingScheme: ManualCpc,
		BudgetType:    "DailyBudgetStandard",
		DailyBudget:   25,
		Description:   "a new campaign",
		Name:          "newcapaign",
		Status:        "Active",
		TimeZone:      "EasternTimeUSCanada",
		CampaignType:  Shopping,
		Settings: []CampaignSettings{{
			Type:             "ShoppingSetting",
			TypeAttr:         "ShoppingSetting",
			Priority:         0,
			SalesCountryCode: "US",
			StoreId:          1397151,
		}},
	}}

	accountid, err := strconv.ParseInt(os.Getenv("BING_ACCOUNT_ID"), 10, 64)

	if err != nil {
		t.Fatal(err)
	}
	ids, err := svc.AddCampaigns(accountid, toadd)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(ids)
}

//addcampaignerror
//<AddCampaignsResponse xmlns="https://bingads.microsoft.com/CampaignManagement/v11"><CampaignIds xmlns:a="http://schemas.datacontract.org/2004/07/System" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><a:long i:nil="true"/></CampaignIds><PartialErrors xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><BatchError><Code>1154</Code><Details i:nil="true"/><ErrorCode>CampaignServiceCampaignShoppingCampaignStoreIdInvalid</ErrorCode><FieldPath i:nil="true"/><ForwardCompatibilityMap i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/System.Collections.Generic"/><Index>0</Index><Message>The store ID of the shopping campaign is invalid.</Message><Type>BatchError</Type></BatchError></PartialErrors></AddCampaignsResponse>

func TestUnmarshalCampaigns(t *testing.T) {
	client := StringClient(`<GetCampaignsByAccountIdResponse xmlns="https://bingads.microsoft.com/CampaignManagement/v11"><Campaigns xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><Campaign><BiddingScheme i:type="ManualCpcBiddingScheme"><Type>ManualCpc</Type></BiddingScheme><BudgetType>DailyBudgetStandard</BudgetType><DailyBudget>25</DailyBudget><Description>dota2</Description><ForwardCompatibilityMap xmlns:a="http://schemas.datacontract.org/2004/07/System.Collections.Generic"/><Id>804002264</Id><Name>dota2</Name><NativeBidAdjustment i:nil="true"/><Status>Active</Status><TimeZone>EasternTimeUSCanada</TimeZone><TrackingUrlTemplate i:nil="true"/><UrlCustomParameters i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/><CampaignType>Shopping</CampaignType><Settings><Setting i:type="ShoppingSetting"><Type>ShoppingSetting</Type><LocalInventoryAdsEnabled i:nil="true"/><Priority>0</Priority><SalesCountryCode>US</SalesCountryCode><StoreId>1387210</StoreId></Setting></Settings><BudgetId i:nil="true"/><Languages i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/></Campaign></Campaigns></GetCampaignsByAccountIdResponse>`)

	svc := &CampaignService{
		Endpoint: "https://campaign.api.sandbox.bingads.microsoft.com/Api/Advertiser/CampaignManagement/v11/CampaignManagementService.svc",
		Client:   client,
	}
	//svc = getTestClient()

	accountid, err := strconv.ParseInt(os.Getenv("BING_ACCOUNT_ID"), 10, 64)

	if err != nil {
		t.Fatal(err)
	}
	res, err := svc.GetCampaignsByAccountId(accountid, Shopping)

	if err != nil {
		t.Fatal(err)
	}

	expected := []Campaign{{
		BiddingScheme: BiddingScheme("ManualCpc"),
		BudgetType:    "DailyBudgetStandard",
		DailyBudget:   25,
		Description:   "dota2",
		Id:            804002264,
		Name:          "dota2",
		Status:        "Active",
		TimeZone:      "EasternTimeUSCanada",
		CampaignType:  Shopping,
		Settings: []CampaignSettings{{
			Type:             "ShoppingSetting",
			Priority:         0,
			SalesCountryCode: "US",
			StoreId:          1387210,
		}},
	}}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

func TestGetSandboxAdGroups(t *testing.T) {
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

func TestUnmarshalAdgroupCriterions(t *testing.T) {
	/*
		s := `
		<AdGroupCriterion i:type="BiddableAdGroupCriterion">
			<AdGroupId>1159984767305062</AdGroupId>
			<Criterion i:type="ProductPartition">
				<Type>ProductPartition</Type>
				<Condition>
					<Attribute>agi</Attribute>
					<Operand>ProductType1</Operand>
				</Condition>
				<ParentCriterionId>4576098675327012</ParentCriterionId>
				<PartitionType>Subdivision</PartitionType>
			</Criterion>
			<Id>4576098675327036</Id>
			<Status>Active</Status>
			<Type>BiddableAdGroupCriterion</Type>
			<CriterionBid i:type="FixedBid">
				<Type>FixedBid</Type>
				<Amount>0.35</Amount>
			</CriterionBid>
			<DestinationUrl i:nil="true"/>
			<EditorialStatus i:nil="true"/>
			<FinalAppUrls i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/>
			<FinalMobileUrls i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/>
			<FinalUrls i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/>
			<TrackingUrlTemplate i:nil="true"/><UrlCustomParameters i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/>
		</AdGroupCriterion>`
		a := BiddableAdGroupCriterion{}
	*/
	s := StringClient(`<GetAdGroupCriterionsByIdsResponse xmlns="https://bingads.microsoft.com/CampaignManagement/v11"><AdGroupCriterions xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><AdGroupCriterion i:type="BiddableAdGroupCriterion"><AdGroupId>1159984767305062</AdGroupId><Criterion i:type="ProductPartition"><Type>ProductPartition</Type><Condition><Attribute/><Operand>All</Operand></Condition><ParentCriterionId i:nil="true"/><PartitionType>Subdivision</PartitionType></Criterion><Id>4576098675327012</Id><Status>Active</Status><Type>BiddableAdGroupCriterion</Type><CriterionBid i:nil="true"/><DestinationUrl i:nil="true"/><EditorialStatus i:nil="true"/><FinalAppUrls i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/><FinalMobileUrls i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/><FinalUrls i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/><TrackingUrlTemplate i:nil="true"/><UrlCustomParameters i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/></AdGroupCriterion><AdGroupCriterion i:type="BiddableAdGroupCriterion"><AdGroupId>1159984767305062</AdGroupId><Criterion i:type="ProductPartition"><Type>ProductPartition</Type><Condition><Attribute>agi</Attribute><Operand>ProductType1</Operand></Condition><ParentCriterionId>4576098675327012</ParentCriterionId><PartitionType>Unit</PartitionType></Criterion><Id>4576098675327013</Id><Status>Active</Status><Type>BiddableAdGroupCriterion</Type><CriterionBid i:type="FixedBid"><Type>FixedBid</Type><Amount>1</Amount></CriterionBid><DestinationUrl i:nil="true"/><EditorialStatus i:nil="true"/><FinalAppUrls i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/><FinalMobileUrls i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/><FinalUrls i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/><TrackingUrlTemplate i:nil="true"/><UrlCustomParameters i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/></AdGroupCriterion><AdGroupCriterion i:type="BiddableAdGroupCriterion"><AdGroupId>1159984767305062</AdGroupId><Criterion i:type="ProductPartition"><Type>ProductPartition</Type><Condition><Attribute>int</Attribute><Operand>ProductType1</Operand></Condition><ParentCriterionId>4576098675327012</ParentCriterionId><PartitionType>Unit</PartitionType></Criterion><Id>4576098675327014</Id><Status>Active</Status><Type>BiddableAdGroupCriterion</Type><CriterionBid i:type="FixedBid"><Type>FixedBid</Type><Amount>1</Amount></CriterionBid><DestinationUrl i:nil="true"/><EditorialStatus i:nil="true"/><FinalAppUrls i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/><FinalMobileUrls i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/><FinalUrls i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/><TrackingUrlTemplate i:nil="true"/><UrlCustomParameters i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/></AdGroupCriterion><AdGroupCriterion i:type="BiddableAdGroupCriterion"><AdGroupId>1159984767305062</AdGroupId><Criterion i:type="ProductPartition"><Type>ProductPartition</Type><Condition><Attribute>str</Attribute><Operand>ProductType1</Operand></Condition><ParentCriterionId>4576098675327012</ParentCriterionId><PartitionType>Unit</PartitionType></Criterion><Id>4576098675327015</Id><Status>Active</Status><Type>BiddableAdGroupCriterion</Type><CriterionBid i:type="FixedBid"><Type>FixedBid</Type><Amount>1</Amount></CriterionBid><DestinationUrl i:nil="true"/><EditorialStatus i:nil="true"/><FinalAppUrls i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/><FinalMobileUrls i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/><FinalUrls i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/><TrackingUrlTemplate i:nil="true"/><UrlCustomParameters i:nil="true" xmlns:a="http:/  chemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/></AdGroupCriterion><AdGroupCriterion i:type="BiddableAdGroupCriterion"><AdGroupId>1159984767305062</AdGroupId><Criterion i:type="ProductPartition"><Type>ProductPartition</Type><Condition><Attribute/><Operand>ProductType1</Operand></Condition><ParentCriterionId>4576098675327012</ParentCriterionId><PartitionType>Unit</PartitionType></Criterion><Id>4576098675327016</Id><Status>Active</Status><Type>BiddableAdGroupCriterion</Type><CriterionBid i:type="FixedBid"><Type>FixedBid</Type><Amount>1</Amount></CriterionBid><DestinationUrl i:nil="true"/><EditorialStatus i:nil="true"/><FinalAppUrls i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/><FinalMobileUrls i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/><FinalUrls i:nil="true" xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/><TrackingUrlTemplate i:nil="true"/><UrlCustomParameters i:nil="true" xmlns:a="http://schemas.datacontract.org/2004/07/Microsoft.AdCenter.Advertiser.CampaignManagement.Api.DataContracts.V11"/></AdGroupCriterion></AdGroupCriterions></GetAdGroupCriterionsByIdsResponse>`)

	svc := &CampaignService{
		Endpoint: "https://campaign.api.sandbox.bingads.microsoft.com/Api/Advertiser/CampaignManagement/v11/CampaignManagementService.svc",
		Client:   s,
	}

	res, err := svc.GetAdGroupCriterionsByIds(1)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)

}
