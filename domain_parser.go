package url

import (
	"index/suffixarray"
	"strings"

	"go.source.hueristiq.com/url/tlds"
)

// DomainParser is responsible for parsing domain names into their constituent parts: subdomain,
// root domain (SLD), and top-level domain (TLD). It utilizes a suffix array to efficiently identify TLDs
// from a comprehensive list of known TLDs (both standard and pseudo-TLDs). This allows the parser to split
// the domain into subdomain, root domain, and TLD components quickly and accurately.
//
// The suffix array helps in handling a large number of known TLDs and enables fast lookups, even for complex
// domain structures where subdomains might be mistaken for TLDs.
//
// Fields:
//   - sa (*suffixarray.Index):
//   - The suffix array index used for efficiently searching through known TLDs.
//   - This allows for rapid identification of the TLD in the domain string.
//
// Example Usage:
//
//	parser := NewDomainParser()
//	domain := "www.example.com"
//	parsedDomain := parser.Parse(domain)
//	fmt.Println(parsedDomain.Subdomain)  // Output: "www"
//	fmt.Println(parsedDomain.SLD)        // Output: "example"
//	fmt.Println(parsedDomain.TLD)        // Output: "com"
type DomainParser struct {
	sa *suffixarray.Index
}

// Parse takes a full domain string (e.g., "www.example.com") and splits it into three main components:
// subdomain, root domain (SLD), and TLD. The method uses the suffix array to identify the TLD and then
// extracts the subdomain and root domain from the rest of the domain string.
//
// Parameters:
//   - domain (string): The full domain string to be parsed.
//
// Returns:
//   - parsed (*Domain): A pointer to a Domain struct containing the subdomain, root domain (SLD), and TLD.
func (p *DomainParser) Parse(domain string) (parsed *Domain) {
	parsed = &Domain{}

	parts := strings.Split(domain, ".")

	if len(parts) <= 1 {
		parsed.SLD = domain

		return
	}

	TLDOffset := p.findTLDOffset(parts)

	if TLDOffset < 0 {
		parsed.SLD = domain

		return
	}

	parsed.Subdomain = strings.Join(parts[:TLDOffset], ".")
	parsed.SLD = parts[TLDOffset]
	parsed.TLD = strings.Join(parts[TLDOffset+1:], ".")

	return
}

// findTLDOffset searches the domain parts to find the position where the TLD starts.
// It works backward through the domain parts, from right (TLD) to left (subdomain),
// to handle complex cases where subdomains might appear similar to TLDs.
//
// This method uses the suffix array to efficiently identify known TLDs.
//
// Parameters:
//   - parts ([]string): A slice of domain components split by '.' (e.g., ["www", "example", "com"]).
//
// Returns:
//   - offset (int): The index of the root domain (SLD) or -1 if no valid TLD is found.
func (p *DomainParser) findTLDOffset(parts []string) (offset int) {
	offset = -1

	partsLength := len(parts)
	partsLastIndex := partsLength - 1

	for i := partsLastIndex; i >= 0; i-- {
		TLD := strings.Join(parts[i:], ".")

		indices := p.sa.Lookup([]byte(TLD), -1)

		if len(indices) > 0 {
			offset = i - 1
		} else {
			break
		}
	}

	return
}

// DomainParserInterface defines the interface for domain parsing functionality.
type DomainParserInterface interface {
	Parse(domain string) (parsed *Domain)

	findTLDOffset(parts []string) (offset int)
}

// DomainParserOptionFunc defines a function type for configuring a DomainParser instance.
// This allows customization options like specifying custom TLDs.
//
// Example:
//
//	parser := NewDomainParser(DomainParserWithTLDs("custom", "tld"))
type DomainParserOptionFunc func(*DomainParser)

// Ensure type compatibility with the DomainParserInterface.
var _ DomainParserInterface = &DomainParser{}

// NewDomainParser creates a new DomainParser instance and initializes it with a comprehensive list
// of TLDs, including both standard TLDs and pseudo-TLDs. Additional options can be passed to customize
// the parser, such as using a custom set of TLDs.
//
// Parameters:
//   - opts (variadic DomainParserOptionFunc): Optional configuration options.
//
// Returns:
//   - parser (*DomainParser): A pointer to the initialized DomainParser.
func NewDomainParser(opts ...DomainParserOptionFunc) (parser *DomainParser) {
	parser = &DomainParser{}

	TLDs := []string{}

	TLDs = append(TLDs, tlds.Official...)
	TLDs = append(TLDs, tlds.Pseudo...)

	parser.sa = suffixarray.New([]byte("\x00" + strings.Join(TLDs, "\x00") + "\x00"))

	for _, opt := range opts {
		opt(parser)
	}

	return
}

// DomainParserWithTLDs allows the DomainParser to be initialized with a custom set of TLDs.
// This option is useful for handling non-standard or niche TLDs that may not be included
// in the default set.
//
// Parameters:
//   - TLDs ([]string): A slice of custom TLDs to be used by the DomainParser.
//
// Returns:
//   - A DomainParserOptionFunc that applies the custom TLDs to the parser.
func DomainParserWithTLDs(TLDs ...string) DomainParserOptionFunc {
	return func(p *DomainParser) {
		p.sa = suffixarray.New([]byte("\x00" + strings.Join(TLDs, "\x00") + "\x00"))
	}
}
