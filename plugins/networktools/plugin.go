package networktools

import (
	"fmt"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

func RunIPv4SubnetCalculator() {
	cmdry.Run(plugin("ipv4-subnet-calculator", "IPv4 Subnet Calculator", "Calculate an IPv4 network, broadcast address, and usable host range locally.", []string{"ipv4", "subnet", "cidr", "mask", "network", "hosts"}, subnetForm, subnetRun))
}

func RunIPv4AddressConverter() {
	cmdry.Run(plugin("ipv4-address-converter", "IPv4 Address Converter", "Convert an IPv4 address among dotted decimal, integer, binary, and hexadecimal locally.", []string{"ipv4", "address", "decimal", "binary", "hexadecimal", "integer"}, addressForm, addressRun))
}

func RunIPv4RangeExpander() {
	cmdry.Run(plugin("ipv4-range-expander", "IPv4 Range Expander", "Expand a bounded IPv4 CIDR or inclusive address range locally.", []string{"ipv4", "cidr", "range", "addresses", "subnet"}, rangeForm, rangeRun))
}

func RunMACAddressLookup() {
	cmdry.Run(plugin("mac-address-lookup", "MAC Address Lookup", "Look up a MAC address prefix in Cmdry's bundled offline vendor table.", []string{"mac", "oui", "vendor", "ethernet", "network"}, lookupForm, lookupRun))
}

func RunMACAddressGenerator() {
	cmdry.Run(plugin("mac-address-generator", "MAC Address Generator", "Generate a random locally administered unicast MAC address.", []string{"mac", "ethernet", "random", "locally administered", "network"}, generatorForm, generatorRun))
}

func RunIPv6ULAGenerator() {
	cmdry.Run(plugin("ipv6-ula-generator", "IPv6 ULA Generator", "Generate a random RFC 4193 IPv6 unique-local /48 prefix locally.", []string{"ipv6", "ula", "unique local", "rfc 4193", "prefix"}, ulaForm, ulaRun))
}

func plugin(id, name, description string, terms []string, overview, action cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "network", Icon: "network", SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": action}}
}

func subnetForm(cmdry.Request) (cmdry.View, error) {
	return form("IPv4 Subnet Calculator", "Calculate IPv4 network details from an address and CIDR prefix.", "Calculate subnet", []cmdry.Field{{Name: "address", Label: "IPv4 address", Type: "text", Value: "192.168.7.233", Required: true}, {Name: "prefix", Label: "CIDR prefix", Type: "number", Value: "24", Min: "0", Max: "32", Required: true}}), nil
}

func subnetRun(request cmdry.Request) (cmdry.View, error) {
	prefix, err := strconv.Atoi(strings.TrimSpace(fmt.Sprint(request.Params["prefix"])))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("CIDR prefix must be a whole number")
	}
	result, err := CalculateSubnet(fmt.Sprint(request.Params["address"]), prefix)
	if err != nil {
		return cmdry.View{}, err
	}
	rows := []map[string]any{{"field": "Address", "value": result.Address}, {"field": "CIDR", "value": fmt.Sprintf("/%d", result.Prefix)}, {"field": "Network", "value": result.Network}, {"field": "Broadcast", "value": result.Broadcast}, {"field": "First usable host", "value": result.FirstHost}, {"field": "Last usable host", "value": result.LastHost}, {"field": "Total addresses", "value": strconv.FormatUint(result.Total, 10)}, {"field": "Usable hosts", "value": strconv.FormatUint(result.Usable, 10)}}
	return tableResult("IPv4 subnet", "Subnet details calculated locally.", "ipv4-subnet", rows), nil
}

func addressForm(cmdry.Request) (cmdry.View, error) {
	return form("IPv4 Address Converter", "Convert one dotted-decimal IPv4 address to common representations.", "Convert address", []cmdry.Field{{Name: "address", Label: "IPv4 address", Type: "text", Value: "192.168.7.233", Required: true}}), nil
}

func addressRun(request cmdry.Request) (cmdry.View, error) {
	values, err := IPv4Representations(fmt.Sprint(request.Params["address"]))
	if err != nil {
		return cmdry.View{}, err
	}
	rows := []map[string]any{{"field": "Dotted decimal", "value": values["dotted"]}, {"field": "Integer", "value": values["decimal"]}, {"field": "Binary", "value": values["binary"]}, {"field": "Hexadecimal", "value": values["hex"]}}
	return tableResult("IPv4 address conversion", "Equivalent values are shown in common IPv4 representations.", "ipv4-address", rows), nil
}

func rangeForm(cmdry.Request) (cmdry.View, error) {
	return form("IPv4 Range Expander", fmt.Sprintf("Enter an IPv4 CIDR or inclusive range. Expansions are limited to %d addresses.", MaxExpandedAddresses), "Expand range", []cmdry.Field{{Name: "range", Label: "CIDR or range", Type: "text", Value: "192.168.7.0/30", Required: true}}), nil
}

func rangeRun(request cmdry.Request) (cmdry.View, error) {
	addresses, err := ExpandIPv4Range(fmt.Sprint(request.Params["range"]))
	if err != nil {
		return cmdry.View{}, err
	}
	rows := make([]map[string]any, 0, len(addresses))
	for index, address := range addresses {
		rows = append(rows, map[string]any{"index": strconv.Itoa(index + 1), "address": address.String()})
	}
	return cmdry.View{Title: "IPv4 range expansion", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: fmt.Sprintf("%d addresses expanded", len(rows)), Message: "The address list was generated locally."}, {Type: "table", ID: "ipv4-range", Columns: []cmdry.Column{{Key: "index", Label: "#"}, {Key: "address", Label: "IPv4 address"}}, Rows: rows}, restart()}}, nil
}

func lookupForm(cmdry.Request) (cmdry.View, error) {
	return form("MAC Address Lookup", "Look up a vendor in the bundled offline OUI table. Unknown prefixes are not guessed.", "Look up MAC", []cmdry.Field{{Name: "mac", Label: "MAC address", Type: "text", Value: "00:50:56:AA:BB:CC", Required: true}}), nil
}

func lookupRun(request cmdry.Request) (cmdry.View, error) {
	mac, vendor, err := LookupMACVendor(fmt.Sprint(request.Params["mac"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return tableResult("MAC address lookup", "Vendor information comes from Cmdry's bundled offline table.", "mac-lookup", []map[string]any{{"field": "MAC address", "value": mac}, {"field": "Vendor", "value": vendor}}), nil
}

func generatorForm(cmdry.Request) (cmdry.View, error) {
	return oneClickView("MAC Address Generator", "Generate a cryptographically random, locally administered unicast MAC address.", "Generate MAC address"), nil
}

func generatorRun(cmdry.Request) (cmdry.View, error) {
	mac, err := GenerateMAC()
	if err != nil {
		return cmdry.View{}, err
	}
	return tableResult("Generated MAC address", "This address is locally administered and unicast.", "mac-generator", []map[string]any{{"field": "MAC address", "value": mac}}), nil
}

func ulaForm(cmdry.Request) (cmdry.View, error) {
	return oneClickView("IPv6 ULA Generator", "Generate a random locally assigned RFC 4193 unique-local IPv6 /48 prefix.", "Generate ULA prefix"), nil
}

func ulaRun(cmdry.Request) (cmdry.View, error) {
	prefix, globalID, err := GenerateULAPrefix()
	if err != nil {
		return cmdry.View{}, err
	}
	firstSubnet := strings.Replace(prefix, "::/48", ":0::/64", 1)
	return tableResult("Generated IPv6 ULA", "The random Global ID was generated locally.", "ipv6-ula", []map[string]any{{"field": "Unique-local prefix", "value": prefix}, {"field": "Global ID", "value": globalID}, {"field": "First subnet", "value": firstSubnet}}), nil
}

func form(title, description, submit string, fields []cmdry.Field) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "form", Title: title, Description: description + " Your data never leaves this device.", Action: "run", Submit: submit, Fields: fields}}}
}

func oneClickView(title, description, button string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "text", Title: title, Text: description + " Your data never leaves this device."}, {Type: "actions", Actions: []cmdry.Action{{ID: "run", Name: button, Method: "write"}}}}}
}

func tableResult(title, message, id string, rows []map[string]any) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: message}, {Type: "table", ID: id, Columns: []cmdry.Column{{Key: "field", Label: "Field"}, {Key: "value", Label: "Value"}}, Rows: rows}, restart()}}
}

func restart() cmdry.Component {
	return cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Start over"}}}
}
