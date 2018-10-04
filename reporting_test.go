package bingads

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func reportingService() *ReportingService {
	config := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.live-int.com/oauth20_token.srf",
			TokenURL: "https://login.live-int.com/oauth20_token.srf",
		},
		RedirectURL: "https://localhost",
	}

	ts := config.TokenSource(context.TODO(), &oauth2.Token{
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
	})
	session := &Session{
		AccountId:      os.Getenv("BING_ACCOUNT_ID"),
		CustomerId:     os.Getenv("BING_CUSTOMER_ID"),
		DeveloperToken: os.Getenv("BING_DEV_TOKEN"),
		HTTPClient:     &http.Client{},
		TokenSource:    ts,
	}

	svc := NewReportingService(session)
	svc.Endpoint = strings.Replace(svc.Endpoint, "bingads", "sandbox.bingads", 1)
	return svc
}

func TestReportProductDimension(t *testing.T) {
	svc := reportingService()
	accountId, _ := strconv.ParseInt(svc.Session.AccountId, 10, 64)
	rr := &ProductDimensionPerformanceReportRequest{
		Scope: ReportScope{
			AccountIds: Longs{
				accountId,
			},
		},
		Aggregation: "Daily",
		Columns: []string{
			"TimePeriod",
			"AccountName",
			"AccountNumber",
			"AdGroupId",
			"AdGroupName",
			"CampaignName",
			"DeviceType",
			"Language",
			"MerchantProductId",
			"Title",
			"Condition",
			"Brand",
			"Impressions",
			"Clicks",
			"Ctr",
			"AverageCpc",
			"Spend",
			"Conversions",
			"Revenue",
		},
		Time: ReportTime{PredefinedTime: "LastMonth"},
	}
	id, err := svc.SubmitReportRequest(rr)
	if err != nil {
		t.Fatal(err)
	}

	pollReport(t, svc, id)

}

func TestReportProductDimensionWithCustomPeriod(t *testing.T) {
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	svc := reportingService()
	accountId, _ := strconv.ParseInt(svc.Session.AccountId, 10, 64)
	rr := &ProductDimensionPerformanceReportRequest{
		Scope: ReportScope{
			AccountIds: Longs{accountId},
		},
		Aggregation: "Daily",
		Columns: []string{
			"TimePeriod",
			"AccountName",
			"AccountNumber",
			"AdGroupId",
			"AdGroupName",
			"CampaignName",
			"DeviceType",
			"Language",
			"MerchantProductId",
			"Title",
			"Condition",
			"Brand",
			"Impressions",
			"Clicks",
			"Ctr",
			"AverageCpc",
			"Spend",
			"Conversions",
			"Revenue",
		},
		Time: ReportTime{
			CustomDateRangeEnd:   Date{Year: now.Year(), Month: int(now.Month()), Day: now.Day()},
			CustomDateRangeStart: Date{Year: lastMonth.Year(), Month: int(lastMonth.Month()), Day: lastMonth.Day()},
		},
	}
	id, err := svc.SubmitReportRequest(rr)
	if err != nil {
		t.Fatal(err)
	}

	pollReport(t, svc, id)
}

func TestReportProductPartition(t *testing.T) {
	svc := reportingService()
	accountId, _ := strconv.ParseInt(svc.Session.AccountId, 10, 64)
	rr := &ProductPartitionPerformanceReportRequest{
		Scope: ReportScope{
			AccountIds: Longs{
				accountId,
			},
		},
		Aggregation: "Daily",
		Columns: []string{
			"TimePeriod",
			"AccountName",
			"AccountNumber",
			"AdGroupCriterionId",
			"AdGroupId",
			"AdGroupName",
			"CampaignId",
			"CampaignName",
			"DeviceType",
			"Impressions",
			"ImpressionSharePercent",
			"Clicks",
			"Ctr",
			"AverageCpc",
			"Spend",
			"Conversions",
			"Revenue",
			"PartitionType",
			"ProductGroup",
		},
		Time: ReportTime{PredefinedTime: "LastMonth"},
	}
	id, err := svc.SubmitReportRequest(rr)
	if err != nil {
		t.Fatal(err)
	}

	svc = reportingService()
	pollReport(t, svc, id)
}

func TestReportProductPartitionWithCustomPeriod(t *testing.T) {
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	svc := reportingService()
	accountId, _ := strconv.ParseInt(svc.Session.AccountId, 10, 64)
	rr := &ProductPartitionPerformanceReportRequest{
		Scope: ReportScope{
			AccountIds: Longs{accountId},
		},
		Aggregation: "Daily",
		Columns: []string{
			"TimePeriod",
			"AccountName",
			"AccountNumber",
			"AdGroupCriterionId",
			"AdGroupId",
			"AdGroupName",
			"CampaignId",
			"CampaignName",
			"DeviceType",
			"Impressions",
			"ImpressionSharePercent",
			"Clicks",
			"Ctr",
			"AverageCpc",
			"Spend",
			"Conversions",
			"Revenue",
			"PartitionType",
			"ProductGroup",
		},
		Time: ReportTime{
			CustomDateRangeEnd:   Date{Year: now.Year(), Month: int(now.Month()), Day: now.Day()},
			CustomDateRangeStart: Date{Year: lastMonth.Year(), Month: int(lastMonth.Month()), Day: lastMonth.Day()},
		},
	}
	id, err := svc.SubmitReportRequest(rr)
	if err != nil {
		t.Fatal(err)
	}

	svc = reportingService()
	pollReport(t, svc, id)
}

func TestReportAdGroup(t *testing.T) {
	svc := reportingService()
	accountId, _ := strconv.ParseInt(svc.Session.AccountId, 10, 64)
	rr := &AdGroupPerformanceReportRequest{
		Scope: ReportScope{
			AccountIds: Longs{
				accountId,
			},
		},
		Aggregation: "Daily",
		Columns: []string{
			"TimePeriod",
			"AccountId",
			"AccountName",
			"AccountNumber",
			"AdGroupId",
			"AdGroupName",
			"DeviceType",
			"Impressions",
			"Clicks",
			"Ctr",
			"AverageCpc",
			"Spend",
			"Conversions",
			"Revenue",
		},
		Time: ReportTime{PredefinedTime: "LastMonth"},
	}
	id, err := svc.SubmitReportRequest(rr)
	if err != nil {
		t.Fatal(err)
	}

	svc = reportingService()
	pollReport(t, svc, id)
}

func TestReportAdGroupWithCustomPeriod(t *testing.T) {
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	svc := reportingService()
	accountId, _ := strconv.ParseInt(svc.Session.AccountId, 10, 64)
	rr := &AdGroupPerformanceReportRequest{
		Scope: ReportScope{
			AccountIds: Longs{accountId},
		},
		Aggregation: "Daily",
		Columns: []string{
			"TimePeriod",
			"AccountId",
			"AccountName",
			"AccountNumber",
			"AdGroupId",
			"AdGroupName",
			"DeviceType",
			"Impressions",
			"Clicks",
			"Ctr",
			"AverageCpc",
			"Spend",
			"Conversions",
			"Revenue",
		},
		Time: ReportTime{
			CustomDateRangeEnd:   Date{Year: now.Year(), Month: int(now.Month()), Day: now.Day()},
			CustomDateRangeStart: Date{Year: lastMonth.Year(), Month: int(lastMonth.Month()), Day: lastMonth.Day()},
		},
	}
	id, err := svc.SubmitReportRequest(rr)
	if err != nil {
		t.Fatal(err)
	}

	svc = reportingService()
	pollReport(t, svc, id)
}

func pollReport(t *testing.T, svc *ReportingService, id string) {
	fmt.Printf("polling for report %s...\n", id)
	for {
		time.Sleep(10 * time.Second)
		status, err := svc.PollGenerateReport(id)
		if err != nil {
			t.Fatal(err)
		}
		switch status.Status {
		case "Success":
			fmt.Println(status.ReportDownloadUrl)
			return
		default:
			fmt.Println(status.Status)
		}
	}
}
