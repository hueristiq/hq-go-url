package url

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/hueristiq/hq-go-url/schemes"
	"github.com/hueristiq/hq-go-url/tlds"
	"github.com/hueristiq/hq-go-url/unicodes"
)

// Extractor is a struct that configures the URL extraction process.
// It provides options for controlling whether URL schemes and hosts are mandatory,
// and allows custom regular expression patterns to be specified for these components.
// This allows fine-grained control over the types of URLs that are extracted from text.
type Extractor struct {
	withScheme        bool   // Specifies if a scheme (e.g., http) is mandatory in extracted URLs.
	withSchemePattern string // A custom regex pattern for matching URL schemes (optional).
	withHost          bool   // Specifies if a host (e.g., domain) is mandatory in extracted URLs.
	withHostPattern   string // A custom regex pattern for matching URL hosts (optional).
}

// CompileRegex constructs and compiles a regular expression based on the Extractor configuration.
// It builds a regex pattern that can capture various forms of URLs, including those with or without
// schemes and hosts. The method also supports custom patterns provided by the user, ensuring that the
// longest possible match for a URL is found, improving accuracy in URL extraction.
func (e *Extractor) CompileRegex() (regex *regexp.Regexp) {
	// Set the default scheme pattern or use the user-specified one.
	schemePattern := ExtractorSchemePattern

	if e.withScheme && e.withSchemePattern != "" {
		schemePattern = e.withSchemePattern
	}

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
	knownTLDPattern := `(?:(?i)` + punycode + `|` + anyOf(append(asciiTLDs, tlds.Pseudo...)...) + `\b|` + anyOf(unicodeTLDs...) + `)`
	domainPattern := `(?:` + _subdomainPattern + knownTLDPattern + `|localhost)`

	// Host and authority patterns for matching URLs with optional ports.
	hostWithoutPortPattern := `(?:` + domainPattern + `|\[` + ExtractorIPv6Pattern + `\]|\b` + ExtractorIPv4Pattern + `\b)`
	hostWithPortOptionalPattern := `(?:` + hostWithoutPortPattern + ExtractorPortOptionalPattern + `)`

	if e.withHost && e.withHostPattern != "" {
		hostWithPortOptionalPattern = e.withHostPattern
	}

	// Authority patterns for matching URLs with optional user info and host.
	_IAuthorityPattern := `(?:` + _IUserInfoOptionalPattern + hostWithPortOptionalPattern + `)`
	_IAuthorityOptionalPattern := _IAuthorityPattern + `?`

	// Define patterns for different types of URLs.
	webURL := _IAuthorityPattern + `(?:/` + pathCont + `|/)?`

	// Emails pattern.
	email := `(?P<relaxedEmail>[a-zA-Z0-9._%\-+]+@` + hostWithPortOptionalPattern + `)`

	URLsWithSchemePattern := schemePattern + _IAuthorityOptionalPattern + pathCont

	// If custom host pattern is provided, prioritize it.
	if e.withHostPattern != "" {
		URLsWithSchemePattern = schemePattern + _IAuthorityPattern + `(?:/` + pathCont + `|/)?`
	}

	// Combine various URL matching patterns for full URL extraction.
	URLsWithHostPattern := webURL + `|` + email

	RelativeURLsPattern := `(\/[\w\/?=&#.-]*)|([\w\/?=&#.-]+?(?:\/[\w\/?=&#.-]+)+)`

	// Select the final pattern based on the configuration.
	var pattern string

	switch {
	case e.withScheme:
		pattern = URLsWithSchemePattern
	case e.withHost:
		pattern = URLsWithSchemePattern + `|` + URLsWithHostPattern
	default:
		pattern = URLsWithSchemePattern + `|` + URLsWithHostPattern + `|` + RelativeURLsPattern
	}

	// Compiling the final regex pattern.
	regex = regexp.MustCompile(pattern)
	// Ensures the longest possible match is found.
	regex.Longest()

	return
}

// ExtractorOptionsFunc defines a function type for configuring Extractor instances.
// It allows users to pass options that modify the behavior of the Extractor, such as whether
// to include schemes or hosts in URL extraction.
type ExtractorOptionsFunc func(*Extractor)

// ExtractorInterface defines the interface that Extractor should implement.
// It ensures that Extractor has the ability to compile regex patterns for URL extraction.
type ExtractorInterface interface {
	CompileRegex() (regex *regexp.Regexp)
}

const (
	_alphaCharacterSet          = `a-zA-Z`
	_digitCHaracterSet          = `0-9`
	_IUnreservedCharacterSet    = _alphaCharacterSet + _digitCHaracterSet + `\-\._~` + unicodes.AllowedUcsChar
	_IEndUnreservedCharacterSet = _alphaCharacterSet + _digitCHaracterSet + `\-_~` + unicodes.AllowedUcsCharMinusPunc
	_subDelimsCharacterSet      = `!\$&'\(\)\*\+,;=`
	_endSubDelimsCharacterSet   = `\$&\+=`
	_pctEncodingPattern         = `%[0-9a-fA-F]{2}`

	_IUserInfoPattern         = `(?:(?:[` + _IUnreservedCharacterSet + _subDelimsCharacterSet + `:]|` + _pctEncodingPattern + `)+@)`
	_IUserInfoOptionalPattern = _IUserInfoPattern + `?`

	ExtractorIPv4Pattern         = `(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])`
	ExtractorNonEmptyIPv6Pattern = `(?:` +
		// 7 colon-terminated chomps, followed by a final chomp or the rest of an elision.
		`(?:[0-9a-fA-F]{1,4}:){7}(?:[0-9a-fA-F]{1,4}|:)|` +
		// 6 chomps, followed by an IPv4 address or elision with final chomp or final elision.
		`(?:[0-9a-fA-F]{1,4}:){6}(?:` + ExtractorIPv4Pattern + `|:[0-9a-fA-F]{1,4}|:)|` +
		// 5 chomps, followed by an elision with optional IPv4 or up to 2 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){5}(?::` + ExtractorIPv4Pattern + `|(?::[0-9a-fA-F]{1,4}){1,2}|:)|` +
		// 4 chomps, followed by an elision with optional IPv4 (optionally preceded by a chomp) or
		// up to 3 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){4}(?:(?::[0-9a-fA-F]{1,4}){0,1}:` + ExtractorIPv4Pattern + `|(?::[0-9a-fA-F]{1,4}){1,3}|:)|` +
		// 3 chomps, followed by an elision with optional IPv4 (preceded by up to 2 chomps) or
		// up to 4 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){3}(?:(?::[0-9a-fA-F]{1,4}){0,2}:` + ExtractorIPv4Pattern + `|(?::[0-9a-fA-F]{1,4}){1,4}|:)|` +
		// 2 chomps, followed by an elision with optional IPv4 (preceded by up to 3 chomps) or
		// up to 5 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){2}(?:(?::[0-9a-fA-F]{1,4}){0,3}:` + ExtractorIPv4Pattern + `|(?::[0-9a-fA-F]{1,4}){1,5}|:)|` +
		// 1 chomp, followed by an elision with optional IPv4 (preceded by up to 4 chomps) or
		// up to 6 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){1}(?:(?::[0-9a-fA-F]{1,4}){0,4}:` + ExtractorIPv4Pattern + `|(?::[0-9a-fA-F]{1,4}){1,6}|:)|` +
		// elision, followed by optional IPv4 (preceded by up to 5 chomps) or up to 7 final chomps.
		// `:` is an intentionally omitted alternative, to avoid matching `::`.
		`:(?:(?::[0-9a-fA-F]{1,4}){0,5}:` + ExtractorIPv4Pattern + `|(?::[0-9a-fA-F]{1,4}){1,7})` +
		`)`
	ExtractorIPv6Pattern = `(?:` + ExtractorNonEmptyIPv6Pattern + `|::)`

	ExtractorPortPattern         = `(?::[0-9]{1,4}|[1-5][0-9]{4}|6[0-5][0-9]{3}\b)`
	ExtractorPortOptionalPattern = ExtractorPortPattern + `?`

	midIPathSegmentChar = _IUnreservedCharacterSet + `%` + _subDelimsCharacterSet + `:@`
	endIPathSegmentChar = _IEndUnreservedCharacterSet + `%` + _endSubDelimsCharacterSet

	_IPrivateCharacters = `\x{E000}-\x{F8FF}\x{F0000}-\x{FFFFD}\x{100000}-\x{10FFFD}`

	midIChar  = `/?#\\` + midIPathSegmentChar + _IPrivateCharacters
	endIChar  = `/#` + endIPathSegmentChar + _IPrivateCharacters
	wellParen = `\((?:[` + midIChar + `]|\([` + midIChar + `]*\))*\)`
	wellBrack = `\[(?:[` + midIChar + `]|\[[` + midIChar + `]*\])*\]`
	wellBrace = `\{(?:[` + midIChar + `]|\{[` + midIChar + `]*\})*\}`
	wellAll   = wellParen + `|` + wellBrack + `|` + wellBrace
	pathCont  = `(?:[` + midIChar + `]*(?:` + wellAll + `|[` + endIChar + `]))+`

	_letter              = `\p{L}`
	_mark                = `\p{M}`
	_number              = `\p{N}`
	_IRICharctersPattern = `[` + _letter + _mark + _number + `](?:[` + _letter + _mark + _number + `\-]*[` + _letter + _mark + _number + `])?`

	_subdomainPattern = `(?:` + _IRICharctersPattern + `\.)+`
)

var (
	// ExtractorSchemePattern defines a general pattern for matching URL schemes.
	// It matches any scheme that starts with alphabetical characters followed by any combination
	// of alphabets, dots, hyphens, or pluses, and ends with "://". It also matches any scheme
	// from a predefined list that does not require authority (host), ending with ":".
	ExtractorSchemePattern = `(?:[a-zA-Z][a-zA-Z.\-+]*://|` + anyOf(schemes.NoAuthority...) + `:)`
	// ExtractorKnownOfficialSchemePattern defines a pattern for matching officially recognized
	// URL schemes. This includes schemes like "http", "https", "ftp", etc., and is strictly based
	// on the schemes defined in the schemes.Schemes slice, ensuring a match ends with "://".
	ExtractorKnownOfficialSchemePattern = `(?:` + anyOf(schemes.Official...) + `://)`
	// ExtractorKnownUnofficialSchemePattern defines a pattern for matching unofficial or
	// less commonly used URL schemes. Similar to the official pattern but based on the
	// schemes.SchemesUnofficial slice, it supports schemes that might not be universally recognized
	// but are valid in specific contexts, ending with "://".
	ExtractorKnownUnofficialSchemePattern = `(?:` + anyOf(schemes.Unofficial...) + `://)`
	// ExtractorKnownNoAuthoritySchemePattern defines a pattern for matching schemes that
	// do not require an authority (host) component. This is useful for schemes like "mailto:",
	// "tel:", and others where a host is not applicable, ending with ":".
	ExtractorKnownNoAuthoritySchemePattern = `(?:` + anyOf(schemes.NoAuthority...) + `:)`
	// ExtractorKnownSchemePattern combines the patterns for officially recognized,
	// unofficial, and no-authority-required schemes into one comprehensive pattern. It is
	// case-insensitive (noted by "(?i)") and designed to match a wide range of schemes, accommodating
	// the broadest possible set of URLs.
	ExtractorKnownSchemePattern = `(?:(?i)(?:` + anyOf(schemes.Official...) + `|` + anyOf(schemes.Unofficial...) + `)://|` + anyOf(schemes.NoAuthority...) + `:)`

	_ ExtractorInterface = &Extractor{}
)

// NewExtractor creates a new Extractor instance with optional configuration.
// The options can be used to customize how URLs are extracted, such as whether
// to include URL schemes or hosts.
func NewExtractor(opts ...ExtractorOptionsFunc) (extractor *Extractor) {
	extractor = &Extractor{}

	for _, opt := range opts {
		opt(extractor)
	}

	return
}

// ExtractorWithScheme returns an option function that configures the Extractor
// to require URL schemes in the extraction process.
func ExtractorWithScheme() ExtractorOptionsFunc {
	return func(e *Extractor) {
		e.withScheme = true
	}
}

// ExtractorWithSchemePattern returns an option function that allows specifying
// a custom regex pattern for matching URL schemes.
func ExtractorWithSchemePattern(pattern string) ExtractorOptionsFunc {
	return func(e *Extractor) {
		e.withScheme = true
		e.withSchemePattern = pattern
	}
}

// ExtractorWithHost returns an option function that configures the Extractor
// to require URL hosts in the extraction process.
func ExtractorWithHost() ExtractorOptionsFunc {
	return func(e *Extractor) {
		e.withHost = true
	}
}

// ExtractorWithHostPattern returns an option function that allows specifying
// a custom regex pattern for matching URL hosts.
func ExtractorWithHostPattern(pattern string) ExtractorOptionsFunc {
	return func(e *Extractor) {
		e.withHost = true
		e.withHostPattern = pattern
	}
}

// anyOf is a helper function that constructs a regex pattern from a list of strings.
// It joins the provided strings into a single regular expression, ensuring that
// each string is properly escaped for use in regex matching.
func anyOf(strs ...string) string {
	var b strings.Builder

	b.WriteString("(?:")

	for i, s := range strs {
		if i != 0 {
			b.WriteByte('|')
		}

		b.WriteString(regexp.QuoteMeta(s))
	}

	b.WriteByte(')')

	return b.String()
}
