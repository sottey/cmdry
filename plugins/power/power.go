package power

import (
	"regexp"
	"strconv"
	"strings"
)

type Info struct {
	Present                  bool
	Source, State, Remaining string
	Percent                  int
}

var percentPattern = regexp.MustCompile(`([0-9]+)%`)

func ParsePMSet(input string) Info {
	info := Info{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Now drawing from '") {
			info.Source = strings.TrimSuffix(strings.TrimPrefix(line, "Now drawing from '"), "'")
		}
		if strings.Contains(line, "InternalBattery") {
			info.Present = true
			if match := percentPattern.FindStringSubmatch(line); len(match) == 2 {
				info.Percent, _ = strconv.Atoi(match[1])
			}
			parts := strings.Split(line, ";")
			if len(parts) > 1 {
				info.State = strings.TrimSpace(parts[1])
			}
			if len(parts) > 2 {
				info.Remaining = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(parts[2]), "present: true"))
			}
		}
	}
	return info
}

func (i Info) Rows() []map[string]any {
	return []map[string]any{{"label": "Power source", "value": value(i.Source)}, {"label": "Battery state", "value": value(i.State)}, {"label": "Estimated time", "value": value(i.Remaining)}}
}
func value(s string) string {
	if s == "" {
		return "Unavailable"
	}
	return s
}
