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
	Id               int64
	Name             string
	ItemCount        int
	ListItems        []NegativeKeyword `xml:"ListIems>SharedListItem"`
}

//these nil ids
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
	e.EncodeElement(len(s.ListItems), st("ItemCount"))
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
	resp, err := c.Session.SendRequest(req, c.Endpoint, "GetSharedEntitiesByAccountId")

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
	SharedList *NegativeKeywordList
}

type GetListItemsBySharedListResponse struct {
	ListItems []NegativeKeyword `xml:"ListItems>SharedListItem"`
}

//MatchType: Exact Phrase Broad Content
type NegativeKeyword struct {
	XMLName   xml.Name `xml:"SharedListItem"`
	Id        int64
	MatchType string
	Text      string
}

func (s *NegativeKeyword) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = []xml.Attr{xml.Attr{Name: xml.Name{Local: "i:type"}, Value: "NegativeKeyword"}}
	e.EncodeToken(start)
	e.EncodeElement("NegativeKeyword", st("Type"))

	if s.Id == 0 {
		e.EncodeElement("", st("Id", "i:nil", "true"))
	} else {
		e.EncodeElement(s.Id, st("Id"))
	}

	e.EncodeElement(s.MatchType, st("MatchType"))
	e.EncodeElement(s.Text, st("Text"))
	e.EncodeToken(xml.EndElement{start.Name})

	return nil
}

func (c *CampaignService) GetListItemsBySharedList(list *NegativeKeywordList) ([]NegativeKeyword, error) {
	req := GetListItemsBySharedListRequest{
		NS:         "https://bingads.microsoft.com/CampaignManagement/v11",
		SharedList: list,
	}
	resp, err := c.Session.SendRequest(req, c.Endpoint, "GetListItemsBySharedList")

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
	SharedEntity *NegativeKeywordList
	ListItems    []NegativeKeyword `xml:"ListItems>SharedListItem"`
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
	resp, err := c.Session.SendRequest(req, c.Endpoint, "AddSharedEntity")

	if err != nil {
		return nil, err
	}

	ret := &AddSharedEntityResponse{}
	err = xml.Unmarshal(resp, ret)
	return ret, err
}

type DeleteSharedEntitiesRequest struct {
	XMLName        xml.Name              `xml:"DeleteSharedEntitiesRequest"`
	NS             string                `xml:"xmlns,attr"`
	SharedEntities []NegativeKeywordList `xml:"SharedEntities>SharedEntity"`
}

type BatchError struct {
	Code      int
	Details   string
	ErrorCode string
	FieldPath string
	Index     int
	Message   string
	Type      string
}

type DeleteSharedEntitiesResponse struct {
	PartialErrors []BatchError `xml:"PartialErrors>BatchError"`
}

func (c *CampaignService) DeleteSharedEntities(list []NegativeKeywordList) ([]BatchError, error) {
	req := DeleteSharedEntitiesRequest{
		NS:             "https://bingads.microsoft.com/CampaignManagement/v11",
		SharedEntities: list,
	}
	res, err := c.Session.SendRequest(req, c.Endpoint, "DeleteSharedEntities")
	if err != nil {
		return nil, err
	}

	ret := &DeleteSharedEntitiesResponse{}
	err = xml.Unmarshal(res, ret)
	return ret.PartialErrors, err
}

type SharedEntityAssociation struct {
	EntityId         int64
	EntityType       string
	SharedEntityId   int64
	SharedEntityType string
}

type SetSharedEntityAssociationsRequest struct {
	XMLName      xml.Name                  `xml:"SetSharedEntityAssociationsRequest"`
	NS           string                    `xml:"xmlns,attr"`
	Associations []SharedEntityAssociation `xml:"Associations>SharedEntityAssociation"`
}

type SetSharedEntityAssociationsResponse struct {
	ListItemIds    []int64 `xml:"ListItemIds>long"`
	SharedEntityId int64
}

func (c *CampaignService) SetSharedEntityAssociations(associations []SharedEntityAssociation) error {
	req := SetSharedEntityAssociationsRequest{
		NS:           "https://bingads.microsoft.com/CampaignManagement/v11",
		Associations: associations,
	}
	_, err := c.Session.SendRequest(req, c.Endpoint, "SetSharedEntityAssociations")

	return err
}

type GetSharedEntityAssociationsBySharedEntityIdsRequest struct {
	XMLName          xml.Name `xml:"GetSharedEntityAssociationsBySharedEntityIdsRequest"`
	NS               string   `xml:"xmlns,attr"`
	EntityType       string
	SharedEntityIds  SharedEntityIds // `xml:"SharedEntityIds>long"`
	SharedEntityType string
}

type SharedEntityIds []int64

func (s SharedEntityIds) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = []xml.Attr{xml.Attr{Name: xml.Name{Local: "i:nil"}, Value: "false"}, xml.Attr{Name: xml.Name{Local: "xmlns:a"}, Value: "http://schemas.microsoft.com/2003/10/Serialization/Arrays"}}
	e.EncodeToken(start)
	for _, id := range s {
		e.EncodeElement(id, st("a:long"))
	}
	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

type GetSharedEntityAssociationsBySharedEntityIdsResponse struct {
	Associations  []SharedEntityAssociation `xml:"Associations>SharedEntityAssociation"`
	PartialErrors []BatchError              `xml:"PartialErrors>BatchError"`
}

//gives us an internal service error in sandbox
func (c *CampaignService) GetSharedEntityAssociationsBySharedEntityIds(ids []int64) (*GetSharedEntityAssociationsBySharedEntityIdsResponse, error) {
	req := GetSharedEntityAssociationsBySharedEntityIdsRequest{
		NS:               "https://bingads.microsoft.com/CampaignManagement/v11",
		EntityType:       "Campaign",
		SharedEntityType: "NegativeKeywordList",
		SharedEntityIds:  SharedEntityIds(ids),
	}
	res, err := c.Session.SendRequest(req, c.Endpoint, "GetSharedEntityAssociationsBySharedEntityIds")
	if err != nil {
		return nil, err
	}

	ret := &GetSharedEntityAssociationsBySharedEntityIdsResponse{}
	err = xml.Unmarshal(res, ret)
	return ret, err
}

type EntityIds []int64

func (s EntityIds) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = []xml.Attr{xml.Attr{Name: xml.Name{Local: "i:nil"}, Value: "false"}, xml.Attr{Name: xml.Name{Local: "xmlns:a"}, Value: "http://schemas.microsoft.com/2003/10/Serialization/Arrays"}}
	e.EncodeToken(start)
	for _, id := range s {
		e.EncodeElement(id, st("a:long"))
	}
	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

type GetSharedEntityAssociationsByEntityIdsRequest struct {
	XMLName          xml.Name  `xml:"GetSharedEntityAssociationsByEntityIdsRequest"`
	NS               string    `xml:"xmlns,attr"`
	EntityIds        EntityIds // `xml:"EntityIds>long"`
	EntityType       string
	SharedEntityType string
}
type GetSharedEntityAssociationsByEntityIdsResponse struct {
	Associations  []SharedEntityAssociation `xml:"Associations>SharedEntityAssociation"`
	PartialErrors []BatchError              `xml:"PartialErrors>BatchError"`
}

func (c *CampaignService) GetSharedEntityAssociationsByEntityIds(ids []int64) (*GetSharedEntityAssociationsByEntityIdsResponse, error) {
	req := GetSharedEntityAssociationsByEntityIdsRequest{
		NS:               "https://bingads.microsoft.com/CampaignManagement/v11",
		EntityType:       "Campaign",
		SharedEntityType: "NegativeKeywordList",
		EntityIds:        ids,
	}
	res, err := c.Session.SendRequest(req, c.Endpoint, "GetSharedEntityAssociationsByEntityIds")
	if err != nil {
		return nil, err
	}

	ret := &GetSharedEntityAssociationsByEntityIdsResponse{}
	err = xml.Unmarshal(res, ret)
	return ret, err
}

type DeleteSharedEntityAssociationsRequest struct {
	XMLName      xml.Name                  `xml:"DeleteSharedEntityAssociationsRequest"`
	NS           string                    `xml:"xmlns,attr"`
	Associations []SharedEntityAssociation `xml:"Associations>SharedEntityAssociation"`
}

type DeleteSharedEntityAssociationsResponse struct {
	PartialErrors []BatchError `xml:"PartialErrors>BatchError"`
}

func (c *CampaignService) DeleteSharedEntityAssociations(list []SharedEntityAssociation) ([]BatchError, error) {
	req := DeleteSharedEntityAssociationsRequest{
		NS:           "https://bingads.microsoft.com/CampaignManagement/v11",
		Associations: list,
	}
	res, err := c.Session.SendRequest(req, c.Endpoint, "DeleteSharedEntityAssociations")
	if err != nil {
		return nil, err
	}

	ret := &DeleteSharedEntityAssociationsResponse{}
	err = xml.Unmarshal(res, ret)
	return ret.PartialErrors, err
}

type AddListItemsToSharedListRequest struct {
	XMLName    xml.Name          `xml:"AddListItemsToSharedListRequest"`
	NS         string            `xml:"xmlns,attr"`
	ListItems  []NegativeKeyword `xml:"ListItems>SharedListItem"`
	SharedList *NegativeKeywordList
}

type AddListItemsToSharedListResponse struct {
	ListItemIds   []int64      `xml:"ListItemIds>long"`
	PartialErrors []BatchError `xml:"PartialErrors>BatchError"`
}

func (c *CampaignService) AddListItemsToSharedList(list *NegativeKeywordList, items []NegativeKeyword) (*AddListItemsToSharedListResponse, error) {
	req := AddListItemsToSharedListRequest{
		NS:         "https://bingads.microsoft.com/CampaignManagement/v11",
		SharedList: list,
		ListItems:  items,
	}
	res, err := c.Session.SendRequest(req, c.Endpoint, "AddListItemsToSharedList")
	if err != nil {
		return nil, err
	}

	ret := &AddListItemsToSharedListResponse{}
	err = xml.Unmarshal(res, ret)
	return ret, err
}
