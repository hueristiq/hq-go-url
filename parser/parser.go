package parser

import (
	"fmt"
	"index/suffixarray"
	"net"
	"net/url"
	"strings"

	"github.com/hueristiq/hq-go-url/tlds"
)

// URL extends the standard net/url.URL struct by embedding it and adding additional domain-related
// information. The Domain field holds a pointer to a Domain struct representing the parsed
// domain broken down into subdomain, second-level domain (SLD), and top-level domain (TLD).
//
// This design enables seamless integration with existing HTTP libraries while providing enhanced
// domain parsing functionality.
type URL struct {
	*url.URL

	Domain *Domain
}

// Domain represents a parsed domain name, broken down into three main components:
//   - Subdomain: The subdomain part of the domain (e.g., "www" in "www.example.com").
//   - SecondLevelDomain (SLD): The core part of the domain (e.g., "example" in "www.example.com").
//   - TopLevelDomain (TLD): The domain extension (e.g., "com" in "www.example.com").
//
// This structure is useful for analysis or manipulation of domain names.
type Domain struct {
	TopLevelDomain    string
	SecondLevelDomain string
	Subdomain         string
}

// String reconstructs a full domain name from its individual components. It joins the non-empty
// components using a dot ("."). Empty components are omitted from the final output.
//
// Examples:
//   - "www.example.com" when all components are provided.
//   - "example.com" when the subdomain is empty.
//   - "example" when only the SLD is provided.
//
// Returns:
//   - A string representing the reconstructed domain.
func (d *Domain) String() (domain string) {
	var parts []string

	if d.Subdomain != "" {
		parts = append(parts, d.Subdomain)
	}

	if d.SecondLevelDomain != "" {
		parts = append(parts, d.SecondLevelDomain)
	}

	if d.TopLevelDomain != "" {
		parts = append(parts, d.TopLevelDomain)
	}

	domain = strings.Join(parts, ".")

	return
}

// Parser is responsible for converting raw URL strings into the custom URL struct that includes
// both the standard URL components and additional domain details. It supports adding a default
// scheme if the URL is missing one, and it uses a suffix array for efficient TLD lookup.
//
// Fields:
//   - scheme (string): The default scheme (e.g., "http", "https") to apply when missing.
//   - sa (*suffixarray.Index): A suffix array used for fast lookup of TLD strings.
type Parser struct {
	scheme string

	sa *suffixarray.Index
}

// SetDefaultScheme sets the default scheme for the Parser. This scheme is applied to URL strings
// that do not already include a scheme.
//
// Parameters:
//   - scheme (string): The default scheme to use (e.g., "http", "https").
func (p *Parser) SetDefaultScheme(scheme string) {
	p.scheme = scheme
}

// SetTLDs configures the Parser to use a custom set of TLDs by building a new suffix array.
// It concatenates the provided TLD strings with a delimiter and initializes the suffix array for lookups.
//
// Parameters:
//   - TLDs (...string): A slice of custom TLD strings to be used by the Parser.
func (p *Parser) SetTLDs(TLDs ...string) {
	p.sa = suffixarray.New([]byte("\x00" + strings.Join(TLDs, "\x00") + "\x00"))
}

// Parse converts a raw URL string into a URL struct that encapsulates both the standard URL
// components (parsed by net/url) and the extracted domain components. If a default scheme has
// been set (via SetDefaultScheme), it is applied to the raw URL string if missing. The host part
// of the URL is further processed to extract the subdomain, SLD, and TLD using a suffix array.
//
// Parameters:
//   - raw (string): The raw URL string to parse.
//
// Returns:
//   - parsed (*URL): A pointer to the resulting URL struct with embedded net/url.URL and domain details.
//   - err (error): An error if the URL cannot be parsed.
func (p *Parser) Parse(raw string) (parsed *URL, err error) {
	parsed = &URL{}

	if p.scheme != "" {
		raw = p.addScheme(raw)
	}

	parsed.URL, err = url.Parse(raw)
	if err != nil {
		err = fmt.Errorf("failed to parse URL: %w", err)

		return
	}

	hostname := parsed.Hostname()

	if net.ParseIP(hostname) == nil {
		parsed.Domain = &Domain{}

		parts := strings.Split(hostname, ".")

		if len(parts) <= 1 {
			parsed.Domain.SecondLevelDomain = hostname

			return
		}

		TLDOffset := p.findTLDOffset(parts)

		if TLDOffset < 0 {
			parsed.Domain.SecondLevelDomain = hostname

			return
		}

		parsed.Domain.Subdomain = strings.Join(parts[:TLDOffset], ".")
		parsed.Domain.SecondLevelDomain = parts[TLDOffset]
		parsed.Domain.TopLevelDomain = strings.Join(parts[TLDOffset+1:], ".")
	}

	return
}

// addScheme is a helper method that ensures the URL string is absolute by adding a default scheme
// if one is missing. It checks for various URL patterns and prepends the scheme accordingly.
//
// Parameters:
//   - inURL (string): The input URL string which may be missing a scheme.
//
// Returns:
//   - outURL (string): The URL string with the scheme added if necessary.
func (p *Parser) addScheme(inURL string) (outURL string) {
	switch {
	case strings.HasPrefix(inURL, "//"):
		outURL = p.scheme + ":" + inURL
	case strings.HasPrefix(inURL, "://"):
		outURL = p.scheme + inURL
	case !strings.Contains(inURL, "//"):
		outURL = p.scheme + "://" + inURL
	default:
		outURL = inURL
	}

	return
}

// findTLDOffset examines the domain components (split by ".") in reverse order to determine the
// starting index of the TLD. It uses a suffix array to quickly verify if a segment of the domain
// matches a known TLD.
//
// Parameters:
//   - parts ([]string): Slice of domain components (e.g., ["www", "example", "com"]).
//
// Returns:
//   - offset (int): The index of the SLD (one position before the TLD begins), or -1 if no valid TLD is found.
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

// OptionFunc defines a function type used for configuring a Parser instance. Options allow customization
// of the Parser (e.g., setting a default scheme or custom TLDs) during initialization.
type OptionFunc func(parser *Parser)

// Interface defines the standard interface for URL parsing functionality. Any type that implements
// this interface must provide a Parse method to convert a raw URL string into a parsed URL struct.
type Interface interface {
	Parse(raw string) (parsed *URL, err error)
}

// Ensure that Parser implements the Interface.
var _ Interface = (*Parser)(nil)

// New creates and initializes a new Parser instance with default settings. It builds a suffix
// array using a default set of TLDs from the imported tlds package. Additional configurations can
// be applied via the provided OptionFunc functions.
//
// Parameters:
//   - ofs (...OptionFunc): A variadic list of OptionFunc functions to configure the Parser.
//
// Returns:
//   - parser (*Parser): A pointer to the newly created Parser instance.
func New(ofs ...OptionFunc) (parser *Parser) {
	parser = &Parser{}

	TLDs := []string{}

	TLDs = append(TLDs, tlds.Official...)
	TLDs = append(TLDs, tlds.Pseudo...)

	parser.sa = suffixarray.New([]byte("\x00" + strings.Join(TLDs, "\x00") + "\x00"))

	for _, f := range ofs {
		f(parser)
	}

	return
}

// WithDefaultScheme returns an OptionFunc that sets the default scheme for the Parser. This option
// ensures that URLs missing a scheme are treated as absolute URLs with the specified scheme.
//
// Parameters:
//   - scheme (string): The default scheme to set (e.g., "http", "https").
//
// Returns:
//   - (OptionFunc): An OptionFunc function that applies the default scheme to a Parser instance.
func WithDefaultScheme(scheme string) OptionFunc {
	return func(parser *Parser) {
		parser.SetDefaultScheme(scheme)
	}
}

// WithTLDs returns an OptionFunc that configures the Parser with a custom set of TLDs. This is useful
// for handling non-standard or niche TLDs that may not be included in the default set.
//
// Parameters:
//   - TLDs (...string): A slice of custom TLD strings to be used by the Parser.
//
// Returns:
//   - (OptionFunc): An OptionFunc function that applies the custom TLDs to the Parser.
func WithTLDs(TLDs ...string) OptionFunc {
	return func(parser *Parser) {
		parser.SetTLDs(TLDs...)
	}
}
