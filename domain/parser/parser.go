package parser

import (
	"index/suffixarray"
	"strings"

	"go.source.hueristiq.com/url/domain"
	"go.source.hueristiq.com/url/tlds"
)

// Parser is responsible for parsing domain names into their constituent parts: subdomain,
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
type Parser struct {
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
func (p *Parser) Parse(unparsed string) (parsed *domain.Domain) {
	parsed = &domain.Domain{}

	parts := strings.Split(unparsed, ".")

	if len(parts) <= 1 {
		parsed.SLD = unparsed

		return
	}

	TLDOffset := p.findTLDOffset(parts)

	if TLDOffset < 0 {
		parsed.SLD = unparsed

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
func (p *Parser) findTLDOffset(parts []string) (offset int) {
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

// Interface defines the interface for domain parsing functionality.
type Interface interface {
	Parse(unparsed string) (parsed *domain.Domain)

	findTLDOffset(parts []string) (offset int)
}

// OptionFunc defines a function type for configuring a Parser instance.
// This allows customization options like specifying custom TLDs.
//
// Example:
//
//	parser := New(DomainWithTLDs("custom", "tld"))
type OptionFunc func(*Parser)

// Ensure type compatibility with the Interface.
var _ Interface = &Parser{}

// New creates a new Parser instance and initializes it with a comprehensive list
// of TLDs, including both standard TLDs and pseudo-TLDs. Additional options can be passed to customize
// the parser, such as using a custom set of TLDs.
//
// Parameters:
//   - opts (variadic OptionFunc): Optional configuration options.
//
// Returns:
//   - parser (*Parser): A pointer to the initialized Parser.
func New(opts ...OptionFunc) (parser *Parser) {
	parser = &Parser{}

	TLDs := []string{}

	TLDs = append(TLDs, tlds.Official...)
	TLDs = append(TLDs, tlds.Pseudo...)

	parser.sa = suffixarray.New([]byte("\x00" + strings.Join(TLDs, "\x00") + "\x00"))

	for _, opt := range opts {
		opt(parser)
	}

	return
}

// WithTLDs allows the Parser to be initialized with a custom set of TLDs.
// This option is useful for handling non-standard or niche TLDs that may not be included
// in the default set.
//
// Parameters:
//   - TLDs ([]string): A slice of custom TLDs to be used by the Parser.
//
// Returns:
//   - A OptionFunc that applies the custom TLDs to the parser.
func WithTLDs(TLDs ...string) OptionFunc {
	return func(p *Parser) {
		p.sa = suffixarray.New([]byte("\x00" + strings.Join(TLDs, "\x00") + "\x00"))
	}
}
