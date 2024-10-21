package url

import "strings"

// Domain represents a parsed domain name, broken down into three main components:
//   - Subdomain: The subdomain part of the domain (e.g., "www" in "www.example.com").
//   - SLD: The root domain, also known as the second-level domain (SLD), which is the core part of the domain
//     (e.g., "example" in "www.example.com").
//   - TLD: The top-level domain (TLD), which is the domain suffix or extension (e.g., "com" in "www.example.com").
//
// This struct is useful in scenarios where you need to manipulate and analyze domain names. It can be applied
// in tasks such as:
//   - Domain validation (e.g., ensuring that domains conform to expected formats).
//   - URL parsing (e.g., breaking down a URL into its domain components).
//   - Domain classification (e.g., identifying and grouping URLs by subdomain, root domain, or TLD).
//
// By splitting a domain into its components, you can easily identify domain hierarchies, manipulate specific parts of
// a domain, or analyze domain names for SEO, security, or categorization purposes.
//
// Example:
//
//	domain := Domain{
//	    Subdomain: "www",  // Subdomain part ("www")
//	    SLD:       "example",  // Second-level domain part ("example")
//	    TLD:       "com",  // Top-level domain part ("com")
//	}
//
//	// Output: "www.example.com"
//	fmt.Println(domain.String())
type Domain struct {
	Subdomain string
	SLD       string
	TLD       string
}

// String reassembles the components of the domain (Subdomain, SLD, and TLD) back into a complete
// domain name string. Non-empty components are joined with a dot ("."). If any component is missing,
// it is omitted from the final output. This method is useful for reconstructing domain names after parsing.
//
// Example:
//   - If Subdomain = "www", SLD = "example", and TLD = "com", the output will be "www.example.com".
//   - If Subdomain is empty, the output will be "example.com".
//   - If both Subdomain and TLD are empty, the output will be just the SLD "example".
//
// Returns:
//   - domain (string): The reconstructed domain name string.
func (d *Domain) String() (domain string) {
	var parts []string

	if d.Subdomain != "" {
		parts = append(parts, d.Subdomain)
	}

	if d.SLD != "" {
		parts = append(parts, d.SLD)
	}

	if d.TLD != "" {
		parts = append(parts, d.TLD)
	}

	domain = strings.Join(parts, ".")

	return
}

// DomainInterface defines an interface for domain representations.
type DomainInterface interface {
	String() (domain string)
}

// Ensure type compatibility with the DomainInterface.
var _ DomainInterface = &Domain{}
