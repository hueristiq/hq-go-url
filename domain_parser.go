package url

import (
	"index/suffixarray"
	"strings"

	"github.com/hueristiq/hq-go-url/tlds"
)

// DomainParser is responsible for parsing domain names into their constituent parts: subdomain,
// root domain, and top-level domain (TLD). It uses a suffix array to efficiently identify the TLDs
// and extract the subdomain and root domain from the full domain string. The suffix array allows
// for fast lookups, even for large sets of known TLDs.
type DomainParser struct {
	sa *suffixarray.Index // Suffix array index for efficiently searching TLDs.
}

// Parse takes a full domain string (e.g., "www.example.com") and splits it into three main components:
// subdomain, root domain, and TLD. The method leverages the suffix array to identify the TLD and then
// extracts the other parts of the domain accordingly.
//
// Returns:
//   - parsed: A pointer to a Domain struct containing the subdomain, root domain, and TLD.
func (p *DomainParser) Parse(domain string) (parsed *Domain) {
	parsed = &Domain{}

	// Split the domain into parts based on the '.' character.
	parts := strings.Split(domain, ".")

	// If the domain has no dots or consists of a single part, treat the entire domain as the root.
	if len(parts) <= 1 {
		parsed.Root = domain

		return
	}

	// Use the findTLDOffset method to determine where the TLD begins in the domain parts.
	TLDOffset := p.findTLDOffset(parts)

	// If no valid TLD is found, treat the entire domain as the root.
	if TLDOffset < 0 {
		parsed.Root = domain

		return
	}

	// Separate the domain components based on the identified TLD offset.
	parsed.Sub = strings.Join(parts[:TLDOffset], ".")
	parsed.Root = parts[TLDOffset]
	parsed.TopLevel = strings.Join(parts[TLDOffset+1:], ".")

	return
}

// findTLDOffset searches the domain parts to find the position where the TLD starts.
// It works backwards through the domain parts to handle complex cases where subdomains may
// appear similar to TLDs. The suffix array is used to efficiently find known TLDs.
//
// Parameters:
//   - parts: A slice of domain components split by '.' (e.g., ["www", "example", "com"]).
//
// Returns:
//   - offset: The index of the last part that is the root domain, or -1 if no valid TLD is found.
func (p *DomainParser) findTLDOffset(parts []string) (offset int) {
	offset = -1

	partsLength := len(parts)
	partsLastIndex := partsLength - 1

	// Iterate through the parts starting from the right (TLD) to the left (subdomain).
	for i := partsLastIndex; i >= 0; i-- {
		// Reconstruct the potential TLD by joining parts from the current index to the end.
		TLD := strings.Join(parts[i:], ".")

		// Search for the TLD in the suffix array.
		indices := p.sa.Lookup([]byte(TLD), -1)

		// If a valid TLD is found, update the offset to the position before the TLD starts.
		if len(indices) > 0 {
			offset = i - 1
		} else {
			break
		}
	}

	return
}

// DomainParserInterface defines the interface for domain parsing functionality.
// This interface ensures that any struct implementing it can parse domain strings
// and find the TLD offset within domain components.
type DomainParserInterface interface {
	Parse(domain string) (parsed *Domain) // Parse the full domain into its components.

	findTLDOffset(parts []string) (offset int) // Find the TLD offset in the domain parts.
}

// DomainParserOptionFunc defines a function type for configuring a DomainParser instance.
// It is used to apply customization options, such as specifying custom TLDs.
type DomainParserOptionFunc func(*DomainParser)

// Ensure type compatibility with the interface.
// This ensures that DomainParser structs correctly implement the required interface.
var _ DomainParserInterface = &DomainParser{}

// NewDomainParser creates a new DomainParser instance, initializing it with a comprehensive list
// of top-level domains (TLDs), including both standard and pseudo-TLDs. Additional options can be
// passed to customize the parser, such as using a custom set of TLDs.
//
// Returns:
//   - parser: A pointer to the initialized DomainParser.
func NewDomainParser(opts ...DomainParserOptionFunc) (parser *DomainParser) {
	parser = &DomainParser{}

	// Combine standard and pseudo-TLDs for comprehensive coverage.
	TLDs := []string{}

	TLDs = append(TLDs, tlds.Official...)
	TLDs = append(TLDs, tlds.Pseudo...)

	// Initialize the suffix array with TLD data.
	parser.sa = suffixarray.New([]byte("\x00" + strings.Join(TLDs, "\x00") + "\x00"))

	// Apply any additional options
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
//   - TLDs: A slice of custom TLDs to be used by the DomainParser.
//
// Returns:
//   - A DomainParserOptionFunc that applies the custom TLDs to the parser.
func DomainParserWithTLDs(TLDs ...string) DomainParserOptionFunc {
	return func(p *DomainParser) {
		p.sa = suffixarray.New([]byte("\x00" + strings.Join(TLDs, "\x00") + "\x00"))
	}
}
