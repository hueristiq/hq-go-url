package parser

import (
	"index/suffixarray"
	"strings"

	"go.source.hueristiq.com/url/tlds"
)

// Domain represents a parsed domain name, broken down into three main components:
//   - Subdomain: The subdomain part of the domain (e.g., "www" in "www.example.com").
//   - SLD: The root domain, also known as the second-level domain (SLD), which is the core part
//     of the domain (e.g., "example" in "www.example.com").
//   - TLD: The top-level domain (TLD), which is the domain suffix or extension (e.g., "com" in "www.example.com").
//
// This struct is useful in scenarios where the domain name must be analyzed or manipulated,
// such as in URL parsing, domain validation, or domain classification tasks.
type Domain struct {
	Subdomain string
	SLD       string
	TLD       string
}

// String reconstructs a full domain name from its components.
//
// This method joins non-empty components with a dot ("."). If any component is missing,
// it is omitted from the final output.
//
// Example Outputs:
// - "www.example.com" (when all components are present)
// - "example.com" (when the subdomain is empty)
// - "example" (when only the SLD is present)
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

// DomainParser is responsible for parsing domain names into their constituent parts: subdomain,
// root domain (SLD), and top-level domain (TLD). It utilizes a suffix array to efficiently identify TLDs
// from a comprehensive list of known TLDs (both standard and pseudo-TLDs). This allows the parser to split
// the domain into subdomain, root domain, and TLD components quickly and accurately.
//
// The suffix array helps in handling a large number of known TLDs and enables fast lookups, even for complex
// domain structures where subdomains might be mistaken for TLDs.
//
// Fields:
//   - sa (*suffixarray.Index): The suffix array index used for efficiently searching through known TLDs.
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
func (p *DomainParser) Parse(unparsed string) (parsed *Domain) {
	parsed = &Domain{}

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

// WithTLDs configures the DomainParser to use a custom set of TLDs by building a new suffix array.
// It takes a list of TLD strings, concatenates them with a separator, and builds the suffix array.
//
// Parameters:
//   - TLDs (...string): A slice of custom TLDs to be used by the Parser.
func (p *DomainParser) WithTLDs(TLDs ...string) {
	p.sa = suffixarray.New([]byte("\x00" + strings.Join(TLDs, "\x00") + "\x00"))
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

// DomainParserOption defines a function type for configuring a Parser instance.
// This allows customization options like specifying custom TLDs.
//
// Example:
//
//	parser := NewDomainParser(DomainDomainParserWithTLDs("custom", "tld"))
type DomainParserOption func(parser *DomainParser)

// DomainInterface defines a standard interface for domain name representation.
type DomainInterface interface {
	String() (domain string)
}

// DomainParserInterface defines the interface for domain parsing functionality.
type DomainParserInterface interface {
	Parse(unparsed string) (parsed *Domain)
	findTLDOffset(parts []string) (offset int)
}

// Ensuring Domain struct implements the DomainInterface.
var _ DomainInterface = (*Domain)(nil)

// Ensuring DomainParser struct implements the DomainParserInterface.
var _ DomainParserInterface = (*DomainParser)(nil)

// NewDomainParser creates a new Parser instance and initializes it with a comprehensive list
// of TLDs, including both standard TLDs and pseudo-TLDs. Additional options can be passed to customize
// the parser, such as using a custom set of TLDs.
//
// Parameters:
//   - options (...DomainParserOption): Optional configuration options.
//
// Returns:
//   - parser (*DomainParser): A pointer to the initialized Parser.
func NewDomainParser(options ...DomainParserOption) (parser *DomainParser) {
	parser = &DomainParser{}

	TLDs := []string{}

	TLDs = append(TLDs, tlds.Official...)
	TLDs = append(TLDs, tlds.Pseudo...)

	parser.sa = suffixarray.New([]byte("\x00" + strings.Join(TLDs, "\x00") + "\x00"))

	for _, opt := range options {
		opt(parser)
	}

	return
}

// DomainParserWithTLDs allows the Parser to be initialized with a custom set of TLDs.
// This option is useful for handling non-standard or niche TLDs that may not be included
// in the default set.
//
// Parameters:
//   - TLDs (...string): A slice of custom TLDs to be used by the Parser.
//
// Returns:
//   - option (DomainParserOption): A Option that applies the custom TLDs to the parser.
func DomainParserWithTLDs(TLDs ...string) (option DomainParserOption) {
	return func(p *DomainParser) {
		p.WithTLDs(TLDs...)
	}
}
