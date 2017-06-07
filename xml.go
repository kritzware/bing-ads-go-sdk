package bingads

import (
	"encoding/xml"
	"fmt"
)

func findAttr(xs []xml.Attr, name string) (string, error) {
	for _, x := range xs {
		if x.Name.Local == name {
			return x.Value, nil
		}
	}

	return "", fmt.Errorf("attribute %s not found", name)
}
