package url

import (
	"regexp"
	"unicode/utf8"

	"github.com/hueristiq/hq-go-url/tlds"
)

// DomainExtractor is responsible for extracting domain names, including both root domains
// and top-level domains (TLDs), using regular expressions. It provides flexibility in the
// domain extraction process by allowing custom patterns for both root domains and TLDs.
type DomainExtractor struct {
	RootDomainPattern     string // Custom regex pattern for matching the root domain (e.g., "example").
	TopLevelDomainPattern string // Custom regex pattern for matching the TLD (e.g., "com").
}

// CompileRegex compiles a regular expression based on the configured DomainExtractor.
// It builds a regex that can match domains, combining the root domain pattern with the top-level domain (TLD) pattern.
// The method separates ASCII and Unicode TLDs and includes a punycode pattern to handle internationalized domain names (IDNs).
// It also ensures that the regex captures the longest possible domain match.
//
// Returns:
//   - regex: The compiled regular expression for matching domain names.
func (e *DomainExtractor) CompileRegex() (regex *regexp.Regexp) {
	// Default root domain pattern or use a user-specified one.
	RootDomainPattern := _subdomainPattern

	if e.RootDomainPattern != "" {
		RootDomainPattern = `(?:\w+[.])*` + e.RootDomainPattern + `\.`
	}

	// Define a pattern for known TLDs, including punycode, ASCII TLDs, and Unicode TLDs.
	// Separate ASCII TLDs from Unicode TLDs for the regular expression.
	var asciiTLDs, unicodeTLDs []string

	for i, tld := range tlds.Official {
		if tld[0] >= utf8.RuneSelf {
			asciiTLDs = tlds.Official[:i:i]
			unicodeTLDs = tlds.Official[i:]

			break
		}
	}

	// Define regular expression components for known TLDs and domains.
	punycode := `xn--[a-z0-9-]+`
	TopLevelDomainPattern := `(?:(?i)` + punycode + `|` + anyOf(append(asciiTLDs, tlds.Pseudo...)...) + `\b|` + anyOf(unicodeTLDs...) + `)`

	if e.TopLevelDomainPattern != "" {
		TopLevelDomainPattern = e.TopLevelDomainPattern
	}

	// Combine the root domain and TLD patterns to form the complete domain pattern.
	pattern := `(?:` + RootDomainPattern + TopLevelDomainPattern + `)`

	if e.RootDomainPattern == "" && e.TopLevelDomainPattern == "" {
		pattern = `(?:` + RootDomainPattern + TopLevelDomainPattern + `|localhost)`
	}

	// Compile the regex and set it to find the longest possible match.
	regex = regexp.MustCompile(pattern)

	regex.Longest()

	return
}

// DomainExtractorOptionFunc defines a function type for configuring a DomainExtractor.
// It allows setting options like custom patterns for root domains and TLDs.
type DomainExtractorOptionFunc func(*DomainExtractor)

// DomainExtractorInterface defines the interface for domain extraction functionality.
// It ensures that any domain extractor can compile regular expressions to match domain names.
type DomainExtractorInterface interface {
	CompileRegex() (regex *regexp.Regexp)
}

// Ensure that DomainExtractor implements the DomainExtractorInterface.
var _ DomainExtractorInterface = &DomainExtractor{}

// NewDomainExtractor creates and initializes a DomainExtractor with optional configurations.
// By default, it uses pre-defined patterns for extracting root domains and TLDs, but custom
// patterns can be applied using the provided options.
//
// Returns:
//   - extractor: A pointer to the initialized DomainExtractor.
func NewDomainExtractor(opts ...DomainExtractorOptionFunc) (extractor *DomainExtractor) {
	extractor = &DomainExtractor{}

	// Apply any provided options to customize the extractor.
	for _, opt := range opts {
		opt(extractor)
	}

	return
}

// DomainExtractorWithRootDomainPattern returns an option function to configure the DomainExtractor
// with a custom regex pattern for matching root domains (e.g., "example" in "example.com").
//
// Parameters:
//   - pattern: The custom root domain regex pattern.
//
// Returns:
//   - A function that applies the custom root domain pattern to the DomainExtractor.
func DomainExtractorWithRootDomainPattern(pattern string) DomainExtractorOptionFunc {
	return func(e *DomainExtractor) {
		e.RootDomainPattern = pattern
	}
}

// DomainExtractorWithTLDPattern returns an option function to configure the DomainExtractor
// with a custom regex pattern for matching top-level domains (TLDs) (e.g., "com" in "example.com").
//
// Parameters:
//   - pattern: The custom TLD regex pattern.
//
// Returns:
//   - A function that applies the custom TLD pattern to the DomainExtractor.
func DomainExtractorWithTLDPattern(pattern string) DomainExtractorOptionFunc {
	return func(e *DomainExtractor) {
		e.TopLevelDomainPattern = pattern
	}
}
