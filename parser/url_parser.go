package parser

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// URL extends the standard net/url URL struct by embedding it and adding additional fields
// for handling domain-related information. This extension provides a more detailed representation
// of the URL by including a separate `Domain` struct that breaks down the domain into Subdomain,
// second-level domain (SLD), and top-level domain (TLD).
type URL struct {
	*url.URL

	Domain *Domain
}

// URLParser is responsible for parsing URLs while also handling domain-related parsing through
// the use of a DomainParser. It extends basic URL parsing functionality by providing support
// for handling custom schemes and extracting domain components such as subdomains, root domains,
// and TLDs.
type URLParser struct {
	dp *DomainParser

	scheme string
}

// URLParser takes a raw URL string and parses it into a custom URL struct that includes:
//   - Standard URL components from `net/url` (scheme, host, path, etc.)
//   - Domain-specific details such as subdomain, root domain, and TLD.
//
// If the URL does not specify a scheme, the default scheme (if any) is added.
// The method also validates and parses the host and port (if specified).
//
// Parameters:
//   - unparsed (string): The raw URL string to parse.
//
// Returns:
//   - parsed (*URL): A pointer to the parsed URL struct containing both standard URL components
//     and domain-specific details.
//   - err (error): An error if the URL cannot be parsed.
func (p *URLParser) Parse(unparsed string) (parsed *URL, err error) {
	parsed = &URL{}

	if p.scheme != "" {
		unparsed = addScheme(unparsed, p.scheme)
	}

	parsed.URL, err = url.Parse(unparsed)
	if err != nil {
		err = fmt.Errorf("error parsing URL: %w", err)

		return
	}

	if net.ParseIP(parsed.Hostname()) == nil {
		parsed.Domain = p.dp.Parse(parsed.Hostname())
	}

	return
}

// URLParserOption defines a function type for configuring a Parser instance.
// It is used to apply various options such as setting the default scheme.
//
// Example:
//
//	parser := NewURLParserOption(URLParserWithDefaultScheme("https"))
type URLParserOption func(parser *URLParser)

// Ensuring URLParser implements the URLParserInterface.
type URLParserInterface interface {
	Parse(unparsed string) (parsed *URL, err error)
}

// Ensure that Parser implements the Interface.
var _ URLParserInterface = (*URLParser)(nil)

// New creates and initializes a new Parser with the given options. The Parser is also
// initialized with a DomainParser for extracting domain-specific details such as subdomain,
// root domain, and TLD. Additional configuration options can be applied using the variadic
// `options` parameter.
//
// Parameters:
//   - options (...URLParserOption): A variadic list of `Option` functions that can configure the Parser.
//
// Returns:
//   - parser (*URLParser): A pointer to the initialized Parser instance.
func NewURLParser(options ...URLParserOption) (parser *URLParser) {
	parser = &URLParser{
		dp: NewDomainParser(),
	}

	for _, option := range options {
		option(parser)
	}

	return
}

// URLParserWithDefaultScheme returns a `Option` that sets the default scheme for the Parser.
// This function allows you to specify a default scheme (e.g., "http" or "https") that will be added
// to URLs that don't provide one.
//
// Parameters:
//   - scheme (string): The default scheme to set (e.g., "http" or "https").
//
// Returns:
//   - option (URLParserOption): A `Option` that applies the default scheme to the Parser.
func URLParserWithDefaultScheme(scheme string) (option URLParserOption) {
	return func(p *URLParser) {
		p.scheme = scheme
	}
}

// addScheme is a helper function that adds a scheme to a URL string if it is missing.
// This ensures that URLs without schemes are treated as absolute URLs instead of relative paths.
//
// Parameters:
//   - inURL (string): The raw URL string, which may or may not have a scheme.
//   - scheme (string): The scheme to be added if one is missing (e.g., "https").
//
// Returns:
//   - outURL (string): The URL with the scheme added, if necessary.
func addScheme(inURL, scheme string) (outURL string) {
	switch {
	case strings.HasPrefix(inURL, "//"):
		outURL = scheme + ":" + inURL
	case strings.HasPrefix(inURL, "://"):
		outURL = scheme + inURL
	case !strings.Contains(inURL, "//"):
		outURL = scheme + "://" + inURL
	default:
		outURL = inURL
	}

	return
}
