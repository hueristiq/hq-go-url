package parser

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	hqgourl "go.source.hueristiq.com/url"
	dparser "go.source.hueristiq.com/url/domain/parser"
)

// Parser is responsible for parsing URLs while also handling domain-related parsing through
// the use of a DomainParser. It extends basic URL parsing functionality by providing support
// for handling custom schemes and extracting domain components such as subdomains, root domains,
// and TLDs.
//
// Fields:
//   - dp (*DomainParser):
//   - A reference to a `DomainParser` used for extracting subdomain, root domain, and TLD information
//     from the host part of the URL.
//   - scheme (string):
//   - The default scheme to use when parsing URLs without a specified scheme. For example,
//     if a URL is missing a scheme (e.g., "www.example.com"), the `scheme` field will prepend a
//     default scheme like "https", resulting in "https://www.example.com".
//
// Methods:
//
//   - Parse(unparsed string) (parsed *URL, err error):
//   - Takes a raw URL string and parses it into a custom `URL` struct that includes both the
//     standard URL components (via the embedded `net/url.URL`) and domain-specific details.
//   - If the URL does not include a scheme, the default scheme is added (if specified).
//   - Additionally, the method uses the DomainParser to break down the domain into subdomain, root domain,
//     and TLD components.
//
// Example Usage:
//
//	parser := New(WithDefaultScheme("https"))
//	parsedURL, err := parser.Parse("example.com/path")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(parsedURL.Scheme)     // Output: https
//	fmt.Println(parsedURL.Hostname()) // Output: example.com
//	fmt.Println(parsedURL.Domain.Root) // Output: example
type Parser struct {
	dp *dparser.Parser

	scheme string
}

// Parse takes a raw URL string and parses it into a custom URL struct that includes:
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
func (p *Parser) Parse(unparsed string) (parsed *hqgourl.URL, err error) {
	parsed = &hqgourl.URL{}

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

// OptionFunc defines a function type for configuring a Parser instance.
// It is used to apply various options such as setting the default scheme.
//
// Example:
//
//	parser := New(WithDefaultScheme("https"))
type OptionFunc func(*Parser)

// Interface defines the interface that all Parser implementations must adhere to.
type Interface interface {
	Parse(unparsed string) (parsed *hqgourl.URL, err error)
}

// Ensure that Parser implements the Interface.
var _ Interface = &Parser{}

// New creates and initializes a new Parser with the given options. The Parser is also
// initialized with a DomainParser for extracting domain-specific details such as subdomain,
// root domain, and TLD. Additional configuration options can be applied using the variadic
// `opts` parameter.
//
// Parameters:
//   - opts: A variadic list of `OptionFunc` functions that can configure the Parser.
//
// Returns:
//   - parser (*Parser): A pointer to the initialized Parser instance.
func New(opts ...OptionFunc) (parser *Parser) {
	parser = &Parser{
		dp: dparser.New(),
	}

	for _, opt := range opts {
		opt(parser)
	}

	return
}

// WithDefaultScheme returns a `OptionFunc` that sets the default scheme for the Parser.
// This function allows you to specify a default scheme (e.g., "http" or "https") that will be added
// to URLs that don't provide one.
//
// Parameters:
//   - scheme (string): The default scheme to set (e.g., "http" or "https").
//
// Returns:
//   - A `OptionFunc` that applies the default scheme to the Parser.
func WithDefaultScheme(scheme string) OptionFunc {
	return func(p *Parser) {
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
