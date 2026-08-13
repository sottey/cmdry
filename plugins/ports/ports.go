package ports

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
)

type Port struct {
	Port                       int
	Protocol, Address, Process string
	PID                        *int
}

var processPattern = regexp.MustCompile(`\(\("([^\"]+)"(?:,pid=([0-9]+))?`)

// ParseSS parses the stable, numeric, headerless output of ss -H -l -n -t -u -p.
func ParseSS(input string) []Port {
	var result []Port
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		protocol := strings.ToUpper(fields[0])
		if protocol != "TCP" && protocol != "UDP" {
			continue
		}
		address, port, ok := splitAddress(fields[4])
		if !ok {
			continue
		}
		item := Port{Port: port, Protocol: protocol, Address: address}
		if m := processPattern.FindStringSubmatch(line); len(m) > 0 {
			item.Process = m[1]
			if len(m) > 2 && m[2] != "" {
				pid, _ := strconv.Atoi(m[2])
				item.PID = &pid
			}
		}
		result = append(result, item)
	}
	return result
}

// ParseLsof parses the machine-readable fields emitted by:
// lsof -nP -FpcPnT -iTCP -sTCP:LISTEN
// lsof -nP -FpcPnT -iUDP
// Connected UDP sockets are excluded because they are clients, not listeners.
func ParseLsof(input string) []Port {
	var result []Port
	var pid *int
	var process, protocol, endpoint string
	haveFile := false
	appendCurrent := func() {
		if !haveFile || (protocol != "TCP" && protocol != "UDP") || strings.Contains(endpoint, "->") {
			return
		}
		address, port, ok := splitAddress(endpoint)
		if !ok {
			return
		}
		result = append(result, Port{Port: port, Protocol: protocol, Address: address, Process: process, PID: pid})
	}
	for _, line := range strings.Split(input, "\n") {
		if len(line) < 2 {
			continue
		}
		switch line[0] {
		case 'p':
			appendCurrent()
			haveFile = false
			process, protocol, endpoint = "", "", ""
			value, err := strconv.Atoi(line[1:])
			if err != nil {
				pid = nil
			} else {
				pid = &value
			}
		case 'c':
			process = line[1:]
		case 'f':
			appendCurrent()
			haveFile = true
			protocol, endpoint = "", ""
		case 'P':
			protocol = strings.ToUpper(line[1:])
		case 'n':
			endpoint = line[1:]
		}
	}
	appendCurrent()
	return result
}
func splitAddress(value string) (string, int, bool) {
	idx := strings.LastIndex(value, ":")
	if idx < 0 {
		return "", 0, false
	}
	port, err := strconv.Atoi(value[idx+1:])
	if err != nil {
		return "", 0, false
	}
	address := strings.Trim(value[:idx], "[]")
	if address == "" {
		address = "*"
	}
	return address, port, true
}
func (p Port) Row() map[string]any {
	row := map[string]any{"port": p.Port, "protocol": p.Protocol, "address": p.Address, "process": nil, "pid": nil}
	if p.Process != "" {
		row["process"] = p.Process
	}
	if p.PID != nil {
		row["pid"] = *p.PID
	}
	return row
}
func Summary(ports []Port) (int, int) {
	tcp, udp := 0, 0
	for _, p := range ports {
		if p.Protocol == "TCP" {
			tcp++
		} else if p.Protocol == "UDP" {
			udp++
		}
	}
	return tcp, udp
}
