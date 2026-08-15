package prettify

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/hermanschaaf/prettyprint"
	"go.yaml.in/yaml/v3"
)

func Format(kind, input string, spaces int) (string, error) {
	if spaces < 1 || spaces > 8 {
		return "", fmt.Errorf("indent must be between 1 and 8 spaces")
	}
	indent := strings.Repeat(" ", spaces)
	switch kind {
	case "json":
		var out bytes.Buffer
		if err := json.Indent(&out, []byte(input), "", indent); err != nil {
			return "", fmt.Errorf("parse JSON: %w", err)
		}
		return out.String() + "\n", nil
	case "xml":
		return formatXML(input, indent)
	case "yaml":
		var node yaml.Node
		if err := yaml.Unmarshal([]byte(input), &node); err != nil {
			return "", fmt.Errorf("parse YAML: %w", err)
		}
		var out bytes.Buffer
		encoder := yaml.NewEncoder(&out)
		encoder.SetIndent(spaces)
		if err := encoder.Encode(&node); err != nil {
			return "", err
		}
		return out.String(), nil
	case "html":
		out, err := prettyprint.Prettify(input, indent)
		if err != nil {
			return "", fmt.Errorf("parse HTML: %w", err)
		}
		return out, nil
	case "javascript":
		result := api.Transform(input, api.TransformOptions{Loader: api.LoaderJS})
		if len(result.Errors) > 0 {
			return "", fmt.Errorf("parse JavaScript: %s", result.Errors[0].Text)
		}
		return reindent(string(result.Code), spaces), nil
	}
	return "", fmt.Errorf("unknown formatter")
}

func formatXML(input, indent string) (string, error) {
	decoder := xml.NewDecoder(strings.NewReader(input))
	var out bytes.Buffer
	encoder := xml.NewEncoder(&out)
	encoder.Indent("", indent)
	for {
		token, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return "", fmt.Errorf("parse XML: %w", err)
		}
		if err := encoder.EncodeToken(token); err != nil {
			return "", err
		}
	}
	if err := encoder.Flush(); err != nil {
		return "", err
	}
	return out.String(), nil
}
func reindent(input string, spaces int) string {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		count := 0
		for strings.HasPrefix(line, "  ") {
			count++
			line = line[2:]
		}
		lines[i] = strings.Repeat(" ", spaces*count) + line
	}
	return strings.Join(lines, "\n")
}
