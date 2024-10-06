package url

import "strings"

// Domain struct represents the structure of a parsed domain name, including its subdomain, root domain, and top-level domain (TLD).
type Domain struct {
	Sub      string
	Root     string
	TopLevel string
}

// String assembles the domain components back into a full domain string.
func (d *Domain) String() (domain string) {
	var parts []string

	if d.Sub != "" {
		parts = append(parts, d.Sub)
	}

	if d.Root != "" {
		parts = append(parts, d.Root)
	}

	if d.TopLevel != "" {
		parts = append(parts, d.TopLevel)
	}

	domain = strings.Join(parts, ".")

	return
}

// DomainInterface defines a standard interface for any domain representation.
type DomainInterface interface {
	String() (domain string)
}
