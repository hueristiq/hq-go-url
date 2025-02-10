package parser

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// URL extends the standard net/url.URL struct by embedding it and adding additional domain-related
// information. The Domain field holds a pointer to a Domain struct which represents the parsed
// domain broken down into subdomain, SLD, and TLD components.
//
// By extending net/url.URL, URL can be used seamlessly with existing HTTP libraries while
// providing extra domain parsing functionality.
type URL struct {
	*url.URL

	Domain *Domain
}

// URLParser is responsible for parsing raw URL strings into a custom URL struct that includes both the
// standard URL components (scheme, host, path, etc.) and domain-specific details obtained via a DomainParser.
// It also supports the addition of a default scheme if the input URL is missing one.
//
// Fields:
//   - dp ( *DomainParser ): A pointer to a DomainParser that extracts domain components from the host.
//   - scheme (string): The default scheme to apply if the raw URL does not include one.
type URLParser struct {
	dp *DomainParser

	scheme string
}

// Parse takes a raw URL string and converts it into a URL struct that encapsulates both the standard
// URL components and the parsed domain information. If a default scheme has been set via WithDefaultScheme,
// it will be added to the raw URL string if missing. The host portion of the URL is further processed by the
// DomainParser to split it into subdomain, SLD, and TLD (if the host is not an IP address).
//
// Parameters:
//   - unparsed (string): The raw URL string to be parsed.
//
// Returns:
//   - parsed (*URL): A pointer to the resulting URL struct containing both net/url components and domain details.
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

// WithDefaultScheme sets the default scheme for the URLParser. This scheme will be prepended to any
// URL strings that do not already include a scheme.
//
// Parameters:
//   - scheme (string): The default scheme to use (e.g., "http", "https").
func (p *URLParser) WithDefaultScheme(scheme string) {
	p.scheme = scheme
}

// URLParserOption defines a function type for configuring a URLParser instance.
// Options can be used to set the default scheme or any other parser-specific configurations.
//
// Example:
//
//	parser := NewURLParser(URLParserWithDefaultScheme("https"))
type URLParserOption func(parser *URLParser)

// URLParserInterface defines the standard interface for URL parsing functionality.
// Any type implementing this interface must provide a Parse method that converts a raw URL string
// into a parsed URL struct.
type URLParserInterface interface {
	Parse(unparsed string) (parsed *URL, err error)
}

// Ensure that Parser implements the Interface.
var _ URLParserInterface = (*URLParser)(nil)

// NewURLParser creates and initializes a new URLParser instance with default settings.
// A new DomainParser is automatically created and associated with the URLParser to handle
// domain extraction from URL hostnames. Additional configuration can be applied via the
// provided URLParserOption functions.
//
// Parameters:
//   - options (...URLParserOption): A variadic list of functions to configure the URLParser.
//
// Returns:
//   - parser (*URLParser): A pointer to the newly created and configured URLParser instance.
func NewURLParser(options ...URLParserOption) (parser *URLParser) {
	parser = &URLParser{
		dp: NewDomainParser(),
	}

	for _, option := range options {
		option(parser)
	}

	return
}

// URLParserWithDefaultScheme returns a URLParserOption that sets the default scheme for the URLParser.
// This option is useful when you want to ensure that URLs missing a scheme are treated as absolute URLs
// with the specified scheme.
//
// Parameters:
//   - scheme (string): The default scheme to set (e.g., "http", "https").
//
// Returns:
//   - option (URLParserOption): A function that applies the default scheme to a URLParser instance.
func URLParserWithDefaultScheme(scheme string) (option URLParserOption) {
	return func(p *URLParser) {
		p.WithDefaultScheme(scheme)
	}
}

// addScheme is a helper function that adds a scheme to a raw URL string if it is missing one.
// It checks for common URL patterns and prepends the specified scheme to ensure the URL is absolute.
//
// Parameters:
//   - inURL (string): The raw URL string, which may lack a scheme.
//   - scheme (string): The scheme to add (e.g., "https") if missing.
//
// Returns:
//   - outURL (string): The resulting URL string with the scheme added, if necessary.
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
