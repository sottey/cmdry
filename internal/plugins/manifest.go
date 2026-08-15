package plugins

import (
	"fmt"
	"regexp"
	"strings"
)

// IDs are URL-path-safe lowercase identifiers. Periods allow reverse-domain
// names such as com.sottey.port-inspector without permitting path separators.
var idPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_.-]*$`)
var versionPattern = regexp.MustCompile(`^v?[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`)

func ValidateManifest(m Manifest) error {
	if m.ProtocolVersion != ProtocolVersion {
		return fmt.Errorf("unsupported protocol version %d", m.ProtocolVersion)
	}
	if !idPattern.MatchString(m.ID) {
		return fmt.Errorf("invalid plugin ID %q", m.ID)
	}
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("plugin name is required")
	}
	if len(m.SearchTerms) > 32 {
		return fmt.Errorf("too many search terms")
	}
	for _, term := range m.SearchTerms {
		if term = strings.TrimSpace(term); term == "" || len(term) > 80 {
			return fmt.Errorf("invalid search term")
		}
	}
	if !versionPattern.MatchString(m.Version) {
		return fmt.Errorf("invalid plugin version %q", m.Version)
	}
	pages, actions := map[string]bool{}, map[string]bool{}
	for _, p := range m.Pages {
		if !idPattern.MatchString(p.ID) || strings.TrimSpace(p.Name) == "" || pages[p.ID] {
			return fmt.Errorf("invalid or duplicate page ID %q", p.ID)
		}
		pages[p.ID] = true
	}
	for _, a := range m.Actions {
		if !idPattern.MatchString(a.ID) || strings.TrimSpace(a.Name) == "" || actions[a.ID] {
			return fmt.Errorf("invalid or duplicate action ID %q", a.ID)
		}
		actions[a.ID] = true
	}
	if len(m.Pages) == 0 || len(m.Actions) == 0 {
		return fmt.Errorf("plugin must declare at least one page and action")
	}
	return nil
}
func ValidID(id string) bool { return idPattern.MatchString(id) }
