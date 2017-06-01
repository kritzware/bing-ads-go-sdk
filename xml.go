package bingads

import (
	"encoding/xml"
	"fmt"
)

func findAttr(xs []xml.Attr, name string) (string, error) {
	fmt.Println(xs)
	for _, x := range xs {
		fmt.Println(x.Name.Local)
		if x.Name.Local == name {
			return x.Value, nil
		}
	}

	return "", fmt.Errorf("attribute %s not found", name)

}
