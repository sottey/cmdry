// Package networktools provides bounded, local-only networking utilities.
package networktools

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/netip"
	"regexp"
	"sort"
	"strings"
)

const MaxExpandedAddresses = 1024

// Subnet describes the address space for an IPv4 prefix.
type Subnet struct {
	Address, Network, Broadcast, FirstHost, LastHost string
	Prefix                                           int
	Total, Usable                                    uint64
}

// CalculateSubnet returns the network details for one IPv4 address and prefix.
func CalculateSubnet(input string, prefix int) (Subnet, error) {
	address, err := parseIPv4(input)
	if err != nil {
		return Subnet{}, err
	}
	if prefix < 0 || prefix > 32 {
		return Subnet{}, fmt.Errorf("prefix length must be between 0 and 32")
	}
	value := ipv4Value(address)
	mask := ipv4Mask(prefix)
	network, broadcast := value&mask, value|^mask
	total := uint64(1) << (32 - prefix)
	first, last, usable := network, broadcast, total
	if prefix <= 30 {
		first, last, usable = network+1, broadcast-1, total-2
	}
	return Subnet{
		Address: address.String(), Prefix: prefix,
		Network: ipv4FromValue(network).String(), Broadcast: ipv4FromValue(broadcast).String(),
		FirstHost: ipv4FromValue(first).String(), LastHost: ipv4FromValue(last).String(),
		Total: total, Usable: usable,
	}, nil
}

// IPv4Representations returns familiar representations of a dotted IPv4 address.
func IPv4Representations(input string) (map[string]string, error) {
	address, err := parseIPv4(input)
	if err != nil {
		return nil, err
	}
	bytes := address.As4()
	parts := make([]string, len(bytes))
	for index, value := range bytes {
		parts[index] = fmt.Sprintf("%08b", value)
	}
	return map[string]string{
		"dotted":  address.String(),
		"decimal": fmt.Sprintf("%d", ipv4Value(address)),
		"binary":  strings.Join(parts, "."),
		"hex":     fmt.Sprintf("0x%02X%02X%02X%02X", bytes[0], bytes[1], bytes[2], bytes[3]),
	}, nil
}

// ExpandIPv4Range expands an IPv4 CIDR or inclusive start-end range.
func ExpandIPv4Range(input string) ([]netip.Addr, error) {
	input = strings.TrimSpace(input)
	var first, last uint32
	if prefix, err := netip.ParsePrefix(input); err == nil {
		if !prefix.Addr().Is4() {
			return nil, fmt.Errorf("only IPv4 CIDRs are supported")
		}
		prefix = prefix.Masked()
		first = ipv4Value(prefix.Addr())
		last = first | ^ipv4Mask(prefix.Bits())
	} else {
		parts := strings.Split(input, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("enter an IPv4 CIDR or an inclusive start-end range")
		}
		start, startErr := parseIPv4(parts[0])
		end, endErr := parseIPv4(parts[1])
		if startErr != nil || endErr != nil {
			return nil, fmt.Errorf("range endpoints must be valid IPv4 addresses")
		}
		first, last = ipv4Value(start), ipv4Value(end)
		if last < first {
			return nil, fmt.Errorf("range end must not be before its start")
		}
	}
	count := uint64(last) - uint64(first) + 1
	if count > MaxExpandedAddresses {
		return nil, fmt.Errorf("range contains %d addresses; the limit is %d", count, MaxExpandedAddresses)
	}
	addresses := make([]netip.Addr, 0, count)
	for value := first; ; value++ {
		addresses = append(addresses, ipv4FromValue(value))
		if value == last {
			break
		}
	}
	return addresses, nil
}

// NormalizeMAC validates a MAC address and returns its canonical colon-separated form.
func NormalizeMAC(input string) (string, [6]byte, error) {
	clean := strings.NewReplacer(":", "", "-", "", ".", "", " ", "").Replace(strings.TrimSpace(input))
	if len(clean) != 12 {
		return "", [6]byte{}, fmt.Errorf("MAC address must contain exactly 12 hexadecimal characters")
	}
	decoded, err := hex.DecodeString(clean)
	if err != nil || len(decoded) != 6 {
		return "", [6]byte{}, fmt.Errorf("MAC address contains invalid hexadecimal characters")
	}
	var mac [6]byte
	copy(mac[:], decoded)
	return formatMAC(mac), mac, nil
}

// LookupMACVendor searches a bundled offline table by OUI prefix.
func LookupMACVendor(input string) (canonical, vendor string, err error) {
	canonical, mac, err := NormalizeMAC(input)
	if err != nil {
		return "", "", err
	}
	if mac[0]&0x02 != 0 {
		return canonical, "Locally administered (no assigned vendor)", nil
	}
	vendor = offlineVendors[[3]byte{mac[0], mac[1], mac[2]}]
	if vendor == "" {
		vendor = "Not found in bundled vendor table"
	}
	return canonical, vendor, nil
}

// GenerateMAC returns a cryptographically random, locally administered unicast MAC address.
func GenerateMAC() (string, error) {
	var mac [6]byte
	if _, err := rand.Read(mac[:]); err != nil {
		return "", fmt.Errorf("generate random MAC address: %w", err)
	}
	mac[0] = mac[0]&^byte(0x01) | 0x02
	return formatMAC(mac), nil
}

// GenerateULAPrefix returns an RFC 4193 locally assigned /48 prefix.
func GenerateULAPrefix() (prefix, globalID string, err error) {
	var id [5]byte
	if _, err := rand.Read(id[:]); err != nil {
		return "", "", fmt.Errorf("generate ULA prefix: %w", err)
	}
	globalID = strings.ToUpper(hex.EncodeToString(id[:]))
	prefix = fmt.Sprintf("fd%02x:%02x%02x:%02x%02x::/48", id[0], id[1], id[2], id[3], id[4])
	return prefix, globalID, nil
}

// ExtractIPAddresses finds valid IPv4 and IPv6 literals, normalizes them, and
// preserves first-seen order. It deliberately does not resolve hostnames.
func ExtractIPAddresses(input string) []string {
	type foundAddress struct {
		index   int
		address netip.Addr
	}
	found := make([]foundAddress, 0)
	for _, match := range ipv4AddressPattern.FindAllStringSubmatchIndex(input, -1) {
		address, err := netip.ParseAddr(input[match[2]:match[3]])
		if err == nil && address.Is4() {
			found = append(found, foundAddress{index: match[2], address: address})
		}
	}
	for start, end := 0, 0; start < len(input); start = end {
		for start < len(input) && !isIPCharacter(input[start]) {
			start++
		}
		end = start
		for end < len(input) && isIPCharacter(input[end]) {
			end++
		}
		if start == end {
			continue
		}
		candidate := strings.Trim(input[start:end], ".")
		address, err := netip.ParseAddr(candidate)
		if err == nil && address.Is6() {
			found = append(found, foundAddress{index: start, address: address})
		}
	}
	sort.Slice(found, func(left, right int) bool { return found[left].index < found[right].index })
	seen := make(map[netip.Addr]struct{})
	addresses := make([]string, 0, len(found))
	for _, value := range found {
		address := value.address.Unmap()
		if _, exists := seen[address]; exists {
			continue
		}
		seen[address] = struct{}{}
		addresses = append(addresses, address.String())
	}
	return addresses
}

var ipv4AddressPattern = regexp.MustCompile(`(?:^|[^0-9A-Za-z])((?:[0-9]{1,3}\.){3}[0-9]{1,3})(?:$|[^0-9A-Za-z])`)
var macAddressPattern = regexp.MustCompile(`(?i)(?:[0-9a-f]{2}[:-]){5}[0-9a-f]{2}|(?:[0-9a-f]{4}\.){2}[0-9a-f]{4}`)

func isIPCharacter(value byte) bool {
	return (value >= '0' && value <= '9') || (value >= 'a' && value <= 'f') || (value >= 'A' && value <= 'F') || value == ':' || value == '.'
}

// ExtractMACAddresses finds common MAC-address spellings, normalizes them to
// uppercase colon-separated form, and preserves first-seen order.
func ExtractMACAddresses(input string) []string {
	matches := macAddressPattern.FindAllStringIndex(input, -1)
	seen := make(map[string]struct{}, len(matches))
	addresses := make([]string, 0, len(matches))
	for _, match := range matches {
		if (match[0] > 0 && isHex(input[match[0]-1])) || (match[1] < len(input) && isHex(input[match[1]])) {
			continue
		}
		canonical, _, err := NormalizeMAC(input[match[0]:match[1]])
		if err != nil {
			continue
		}
		if _, exists := seen[canonical]; exists {
			continue
		}
		seen[canonical] = struct{}{}
		addresses = append(addresses, canonical)
	}
	return addresses
}

func isHex(value byte) bool {
	return (value >= '0' && value <= '9') || (value >= 'a' && value <= 'f') || (value >= 'A' && value <= 'F')
}

func parseIPv4(input string) (netip.Addr, error) {
	address, err := netip.ParseAddr(strings.TrimSpace(input))
	if err != nil || !address.Is4() {
		return netip.Addr{}, fmt.Errorf("enter a valid IPv4 address")
	}
	return address, nil
}

func ipv4Value(address netip.Addr) uint32 {
	bytes := address.As4()
	return uint32(bytes[0])<<24 | uint32(bytes[1])<<16 | uint32(bytes[2])<<8 | uint32(bytes[3])
}

func ipv4FromValue(value uint32) netip.Addr {
	return netip.AddrFrom4([4]byte{byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
}

func ipv4Mask(prefix int) uint32 {
	if prefix == 0 {
		return 0
	}
	return ^uint32(0) << (32 - prefix)
}

func formatMAC(mac [6]byte) string {
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}
