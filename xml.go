package bingads

import (
	"encoding/xml"
	"fmt"
)

func ats(attrs ...string) []xml.Attr {
	r := []xml.Attr{}

	for i := 0; i < len(attrs); i += 2 {
		r = append(r, xml.Attr{Name: xml.Name{Local: attrs[i]}, Value: attrs[i+1]})
	}

	return r
}

func st(name string, attrs ...string) xml.StartElement {
	ret := xml.StartElement{
		Name: xml.Name{Local: name},
	}

	for i := 0; i < len(attrs); i += 2 {
		ret.Attr = append(ret.Attr, xml.Attr{xml.Name{Local: attrs[i]}, attrs[i+1]})
	}

	return ret
}

func findAttr(xs []xml.Attr, name string) (string, error) {
	for _, x := range xs {
		if x.Name.Local == name {
			return x.Value, nil
		}
	}

	return "", fmt.Errorf("attribute %s not found", name)
}
