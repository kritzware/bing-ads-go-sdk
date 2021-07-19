package bingads

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func reportingService() *ReportingService {
	// if os.Getenv("TEST_PROD") != "" {
	// 	session := getProdClient()
	// 	return NewReportingService(session)
	// }

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

	url := pollReport(t, svc, id)
	if err := verifyReportColumns(url, rr.Columns); err != nil {
		t.Error(err)
	}

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

	url := pollReport(t, svc, id)
	if err := verifyReportColumns(url, rr.Columns); err != nil {
		t.Error(err)
	}
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

	url := pollReport(t, svc, id)
	if err := verifyReportColumns(url, rr.Columns); err != nil {
		t.Error(err)
	}
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

	url := pollReport(t, svc, id)
	if err := verifyReportColumns(url, rr.Columns); err != nil {
		t.Error(err)
	}
}

func TestReportAdGroupDaily(t *testing.T) {
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

	url := pollReport(t, svc, id)
	if err := verifyReportColumns(url, rr.Columns); err != nil {
		t.Error(err)
	}
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

	url := pollReport(t, svc, id)
	if err := verifyReportColumns(url, rr.Columns); err != nil {
		t.Error(err)
	}
}

func compareColumns(xs, ys []string) bool {
	if len(xs) != len(ys) {
		return false
	}

	for i := range xs {
		if xs[i] != ys[i] {
			fmt.Println(xs[i], ys[i], len(xs[i]), len(ys[i]))
			return false
		}
	}

	return true
}

func verifyReportColumns(url string, expectedCols []string) error {
	if url == "" {
		fmt.Println("WARN: sandbox reports can't be downloaded")
		return nil
	}

	fmt.Println("downloading: " + url)
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	breader := bytes.NewReader(b)

	reader, err := zip.NewReader(breader, breader.Size())
	if err != nil {
		return err
	}

	if len(reader.File) != 1 {
		return fmt.Errorf("expected 1 file, got: %d", len(reader.File))
	}

	f, err := reader.File[0].Open()
	if err != nil {
		return err
	}

	defer f.Close()
	br := bufio.NewReader(f)

	line, _, err := br.ReadLine()
	if err != nil {
		return err
	}

	cols := strings.Split(strings.Replace(string(line), "\"", "", -1), ",")
	//removing the BOM
	cols[0] = cols[0][3:]

	if !compareColumns(cols, expectedCols) {
		return fmt.Errorf("expected cols: %v, got cols: %v", expectedCols, cols)
	}

	return nil
}

func pollReport(t *testing.T, svc *ReportingService, id string) string {
	for {
		fmt.Printf("polling for report %s...\n", id)
		time.Sleep(4 * time.Second)
		status, err := svc.PollGenerateReport(id)
		if err != nil {
			t.Fatal(err)
		}
		switch status.Status {
		case "Success":
			return status.ReportDownloadUrl
		case "Error":
			t.Fatalf("error polling report")
		default:
		}
	}
}
