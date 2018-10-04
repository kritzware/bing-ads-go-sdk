package bingads

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type BulkService struct {
	Endpoint string
	Session  *Session
}

func NewBulkService(session *Session) *BulkService {
	return &BulkService{
		Endpoint: "https://bulk.api.bingads.microsoft.com/Api/Advertiser/CampaignManagement/V12/BulkService.svc",
		Session:  session,
	}
}

type BulkAccountIds struct {
	XMLName xml.Name `xml:"AccountIds"`
	NS      string   `xml:"xmlns:a2,attr"`
	Ids     []int64  `xml:"a2:long"`
}

type GetBulkCampaignsByAccountIdRequest struct {
	XMLName          xml.Name `xml:"DownloadCampaignsByAccountIdsRequest"`
	NS               string   `xml:"xmlns,attr"`
	AccountIds       BulkAccountIds
	DownloadEntities []string `xml:"DownloadEntities>DownloadEntity"`
	FormatVersion    string
	CompressionType  string
}

type GetBulkDownloadStatusResponse struct {
	PercentComplete int
	RequestStatus   string
	ResultFileUrl   string
}

//GetBulkCampaignsByAccountId
func (c *BulkService) GetAdGroupCampaign(adgroup int64) (int64, error) {
	accountid, _ := strconv.ParseInt(c.Session.AccountId, 10, 64)
	req := GetBulkCampaignsByAccountIdRequest{
		NS: "https://bingads.microsoft.com/CampaignManagement/v12",
		AccountIds: BulkAccountIds{
			NS:  "http://schemas.microsoft.com/2003/10/Serialization/Arrays",
			Ids: []int64{accountid},
		},
		CompressionType:  "Zip",
		DownloadEntities: []string{"Campaigns", "AdGroups"},
		FormatVersion:    "6.0",
	}

	resp, err := c.Session.SendRequest(req, c.Endpoint, "DownloadCampaignsByAccountIds")
	if err != nil {
		return -1, err
	}

	type DownloadCampaignsByAccountIdsResponse struct {
		DownloadRequestId string
	}

	downloadresponse := DownloadCampaignsByAccountIdsResponse{}
	if err = xml.Unmarshal(resp, &downloadresponse); err != nil {
		return -1, err
	}

	type GetBulkDownloadStatusRequest struct {
		XMLName   xml.Name `xml:"GetBulkDownloadStatusRequest"`
		NS        string   `xml:"xmlns,attr"`
		RequestId string
	}

	req2 := GetBulkDownloadStatusRequest{
		RequestId: downloadresponse.DownloadRequestId,
		NS:        "https://bingads.microsoft.com/CampaignManagement/v12",
	}

	url, err := func() (string, error) {
		//2 minutes max wait
		for i := 0; i < 20; i++ {
			time.Sleep(time.Duration(1000+500*i) * time.Millisecond)
			resp, err = c.Session.SendRequest(req2, c.Endpoint, "GetBulkDownloadStatus")
			if err != nil {
				return "", err
			}
			status := GetBulkDownloadStatusResponse{}
			if err = xml.Unmarshal(resp, &status); err != nil {
				return "", err
			}

			switch status.RequestStatus {
			case "InProgress":
			case "Completed":
				return status.ResultFileUrl, nil
			default:
				return "", fmt.Errorf(status.RequestStatus)

			}
		}

		return "", fmt.Errorf("timed out waiting for bulk request")
	}()
	if err != nil {
		return -1, err
	}

	greq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	gres, err := c.Session.HTTPClient.Do(greq)
	if err != nil {
		return -1, err
	}

	defer gres.Body.Close()

	buffer := bytes.NewBuffer(nil)
	n, err := buffer.ReadFrom(gres.Body)
	if err != nil {
		return -1, err
	}

	reader, err := zip.NewReader(bytes.NewReader(buffer.Bytes()), n)
	if err != nil {
		return -1, err
	}
	b, _ := reader.File[0].Open()

	csvReader := csv.NewReader(b)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return -1, err
	}

	if len(rows) == 0 {
		return -1, fmt.Errorf("missing header")
	}

	//Type,Status,Id,Parent Id,
	header := rows[0]
	idcol, err := func() (int, error) {
		for i := range header {
			if header[i] == "Id" {
				return i, nil
			}
		}
		return -1, fmt.Errorf("missing Id col")
	}()
	if err != nil {
		return -1, err
	}

	parentcol, err := func() (int, error) {
		for i := range header {
			if header[i] == "Parent Id" {
				return i, nil
			}
		}
		return -1, fmt.Errorf("missing Parent Id col")
	}()
	if err != nil {
		return -1, err
	}

	target := strconv.FormatInt(adgroup, 10)

	for i := 1; i < len(rows); i++ {
		if rows[i][idcol] == target {
			return strconv.ParseInt(rows[i][parentcol], 10, 64)
		}
	}

	return -1, fmt.Errorf("unable to find matching adgroup")
}
