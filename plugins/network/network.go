// Package network parses concise interface and route command output.
package network

import (
	"bufio"
	"strings"
)

type Interface struct {
	Name, Status string
	Addresses    []string
}

func (i Interface) Row() map[string]any {
	status := i.Status
	if status == "" {
		status = "unknown"
	}
	return map[string]any{"interface": i.Name, "status": status, "addresses": strings.Join(i.Addresses, ", ")}
}

// ParseDarwinIfconfig parses macOS ifconfig blocks.
func ParseDarwinIfconfig(input string) []Interface {
	items := make([]Interface, 0)
	var current *Interface
	flush := func() {
		if current != nil {
			items = append(items, *current)
		}
	}
	for _, line := range strings.Split(input, "\n") {
		if line == "" || (!strings.HasPrefix(line, "\t") && strings.Contains(line, ": flags=")) {
			flush()
			if colon := strings.IndexByte(line, ':'); colon > 0 {
				current = &Interface{Name: line[:colon]}
			} else {
				current = nil
			}
			continue
		}
		if current == nil {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 && (fields[0] == "inet" || fields[0] == "inet6") {
			address := strings.Split(fields[1], "%")[0]
			current.Addresses = append(current.Addresses, address)
		}
		if len(fields) >= 2 && fields[0] == "status:" {
			current.Status = fields[1]
		}
	}
	flush()
	return items
}

// ParseLinuxAddresses parses `ip -o addr show` output, adding each address to
// its interface. Link state is supplied separately by ParseLinuxLinks.
func ParseLinuxAddresses(input string) []Interface {
	byName := map[string]*Interface{}
	order := make([]string, 0)
	for _, line := range strings.Split(input, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 4 || (fields[2] != "inet" && fields[2] != "inet6") {
			continue
		}
		name := strings.Split(fields[1], "@")[0]
		item := byName[name]
		if item == nil {
			item = &Interface{Name: name}
			byName[name] = item
			order = append(order, name)
		}
		item.Addresses = append(item.Addresses, strings.Split(fields[3], "/")[0])
	}
	items := make([]Interface, 0, len(order))
	for _, name := range order {
		items = append(items, *byName[name])
	}
	return items
}

func ParseGateway(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		for index, field := range fields {
			if (field == "gateway:" || field == "via") && index+1 < len(fields) {
				return fields[index+1]
			}
		}
	}
	return ""
}
