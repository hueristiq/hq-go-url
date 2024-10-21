package url

import "strings"

// Domain represents a parsed domain name. It breaks down the domain into three main components:
//   - Sub: The subdomain part of the domain (e.g., "www" in "www.example.com").
//   - Root: The root domain or second-level domain (e.g., "example" in "www.example.com").
//   - TopLevel: The top-level domain (TLD) (e.g., "com" in "www.example.com").
//
// This struct allows for easy manipulation of domain name parts and can be used in various domain parsing
// and validation scenarios.
type Domain struct {
	Sub      string // The subdomain part of the domain.
	Root     string // The root domain or second-level domain.
	TopLevel string // The top-level domain (TLD).
}

// String reassembles the components of the domain (subdomain, root domain, and TLD) back into a complete
// domain name string. It joins the non-empty components using a dot ("."). If any of the components are empty,
// they will be omitted from the final string.
//
// Example:
//
//	If Sub = "www", Root = "example", and TopLevel = "com", the output will be "www.example.com".
func (d *Domain) String() (domain string) {
	var parts []string

	// Add subdomain if it exists.
	if d.Sub != "" {
		parts = append(parts, d.Sub)
	}

	// Add root domain if it exists.
	if d.Root != "" {
		parts = append(parts, d.Root)
	}

	// Add top-level domain if it exists.
	if d.TopLevel != "" {
		parts = append(parts, d.TopLevel)
	}

	// Join the parts with dots to form the full domain.
	domain = strings.Join(parts, ".")

	return
}

// DomainInterface defines an interface for Domain representations. It ensures that any Domain struct
// can implement the String method to convert the domain structure back into a string form.
type DomainInterface interface {
	String() (domain string) // Converts the domain components back into a full domain name.
}

// Ensure type compatibility with the interface.
// This ensures that Domain structs correctly implement the required interface.
var _ DomainInterface = &Domain{}
