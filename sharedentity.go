package bingads

import (
	"encoding/xml"
)

type GetSharedEntitiesByAccountIdRequest struct {
	XMLName          xml.Name `xml:"GetSharedEntitiesByAccountIdRequest"`
	NS               string   `xml:"xmlns,attr"`
	SharedEntityType string
}

type GetSharedEntitiesByAccountIdResponse struct {
	SharedEntities []NegativeKeywordList `xml:"SharedEntities>SharedEntity"`
}

//inherits from sharedentity and sharedlist
type NegativeKeywordList struct {
	AssociationCount int
	Id               int64 `xml:",omitempty"`
	Name             string
	ItemCount        int               `xml:",omitempty"`
	ListItems        []NegativeKeyword `xml:"ListIems>SharedListItem"`
}

func (s *NegativeKeywordList) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = []xml.Attr{xml.Attr{Name: xml.Name{Local: "i:type"}, Value: "NegativeKeywordList"}}
	e.EncodeToken(start)
	e.EncodeElement(s.AssociationCount, st("AssociationCount"))
	if s.Id == 0 {
		if err := e.EncodeElement("", st("Id", "i:nil", "true")); err != nil {
			return err
		}
	} else {
		if err := e.EncodeElement(s.Id, st("Id")); err != nil {
			return err
		}
	}
	e.EncodeElement(s.Name, st("Name"))
	e.EncodeElement("NegativeKeywordList", st("Type"))
	e.EncodeElement(s.ItemCount, st("ItemCount"))
	if err := e.Encode(s.ListItems); err != nil {
		return err
	}

	e.EncodeToken(xml.EndElement{start.Name})

	return nil

}

func (c *CampaignService) GetSharedEntitiesByAccountId(entityType string) ([]NegativeKeywordList, error) {
	req := GetSharedEntitiesByAccountIdRequest{
		NS:               "https://bingads.microsoft.com/CampaignManagement/v11",
		SharedEntityType: "NegativeKeywordList",
	}
	resp, err := c.client.SendRequest(req, c.endpoint, "GetSharedEntitiesByAccountId")

	if err != nil {
		return nil, err
	}

	ret := GetSharedEntitiesByAccountIdResponse{}
	err = xml.Unmarshal(resp, &ret)

	if err != nil {
		return nil, err
	}

	return ret.SharedEntities, err

}

type GetListItemsBySharedListRequest struct {
	XMLName    xml.Name `xml:"GetListItemsBySharedListRequest"`
	NS         string   `xml:"xmlns,attr"`
	SharedList NegativeKeywordList
}

type GetListItemsBySharedListResponse struct {
	ListItems []NegativeKeyword `xml:"ListItems>SharedListItem"`
}

//Exact Phrase Broad Content
type NegativeKeyword struct {
	XMLName   xml.Name `xml:"NegativeKeyword"`
	Id        int64    `xml:",omitempty"`
	TypeAttr  string   `xml:"i:type,attr,omitempty"`
	Type      string
	MatchType string
	Text      string
}

func (c *CampaignService) GetListItemsBySharedList(list NegativeKeywordList) ([]NegativeKeyword, error) {
	req := GetListItemsBySharedListRequest{
		NS:         "https://bingads.microsoft.com/CampaignManagement/v11",
		SharedList: list,
	}
	resp, err := c.client.SendRequest(req, c.endpoint, "GetListItemsBySharedList")

	if err != nil {
		return nil, err
	}

	ret := GetListItemsBySharedListResponse{}
	err = xml.Unmarshal(resp, &ret)

	if err != nil {
		return nil, err
	}

	return ret.ListItems, err
}

type AddSharedEntityRequest struct {
	XMLName      xml.Name `xml:"AddSharedEntityRequest"`
	NS           string   `xml:"xmlns,attr"`
	ListItems    []NegativeKeyword
	SharedEntity *NegativeKeywordList
}

type AddSharedEntityResponse struct {
	ListItemIds    []int64 `xml:"ListItemIds>long"`
	SharedEntityId int64
}

func (c *CampaignService) AddSharedEntity(entity *NegativeKeywordList, items []NegativeKeyword) (*AddSharedEntityResponse, error) {
	req := AddSharedEntityRequest{
		NS:           "https://bingads.microsoft.com/CampaignManagement/v11",
		SharedEntity: entity,
		ListItems:    items,
	}
	resp, err := c.client.SendRequest(req, c.endpoint, "AddSharedEntity")

	if err != nil {
		return nil, err
	}

	ret := &AddSharedEntityResponse{}
	err = xml.Unmarshal(resp, ret)
	return ret, err
}
