// Package xmlutils validates pasted XML without reading external resources.
package xmlutils

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// DocumentInfo describes the one root element in a validated XML document.
type DocumentInfo struct {
	Root       string
	Attributes int
}

// Validate checks for one well-formed XML root element and returns its details.
func Validate(input string) (DocumentInfo, error) {
	if strings.TrimSpace(input) == "" {
		return DocumentInfo{}, fmt.Errorf("XML input is required")
	}
	decoder := xml.NewDecoder(strings.NewReader(input))
	depth, roots := 0, 0
	var result DocumentInfo
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			if syntax, ok := err.(*xml.SyntaxError); ok {
				return DocumentInfo{}, fmt.Errorf("line %d: %s", syntax.Line, syntax.Msg)
			}
			return DocumentInfo{}, err
		}
		switch value := token.(type) {
		case xml.StartElement:
			if depth == 0 {
				roots++
				if roots > 1 {
					return DocumentInfo{}, fmt.Errorf("XML contains more than one root element")
				}
				result.Root = qualifiedName(value.Name)
				result.Attributes = len(value.Attr)
			}
			depth++
		case xml.EndElement:
			depth--
		}
	}
	if roots == 0 {
		return DocumentInfo{}, fmt.Errorf("XML document has no root element")
	}
	return result, nil
}

func qualifiedName(name xml.Name) string {
	if name.Space == "" {
		return name.Local
	}
	return name.Space + ":" + name.Local
}
