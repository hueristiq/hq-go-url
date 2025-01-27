package url

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"go.source.hueristiq.com/url/schemes"
	"go.source.hueristiq.com/url/tlds"
	"go.source.hueristiq.com/url/unicodes"
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

// ExtractorOptionFunc defines a function type for configuring Extractor instances.
// It allows users to pass options that modify the behavior of the Extractor, such as whether
// to include schemes or hosts in URL extraction.
type ExtractorOptionFunc func(*Extractor)

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
	// It matches any URL scheme that starts with alphabetical characters (a-z, A-Z), followed by
	// any combination of alphabets, dots (.), hyphens (-), or plus signs (+), and ends with "://".
	// Additionally, it matches schemes from a predefined list that do not require an authority (host),
	// ending with just a colon (":"). These are known as "no-authority" schemes (e.g., "mailto:").
	//
	// This pattern covers a broad range of schemes, making it versatile for extracting different types
	// of URLs, whether they require an authority component or not.
	ExtractorSchemePattern = `(?:[a-zA-Z][a-zA-Z.\-+]*://|` + anyOf(schemes.NoAuthority...) + `:)`

	// ExtractorKnownOfficialSchemePattern defines a pattern for matching officially recognized
	// URL schemes. These include well-known schemes such as "http", "https", "ftp", etc., as registered
	// with IANA. The pattern ensures that the scheme is followed by "://".
	//
	// This pattern ensures that only officially recognized schemes are matched.
	ExtractorKnownOfficialSchemePattern = `(?:` + anyOf(schemes.Official...) + `://)`

	// ExtractorKnownUnofficialSchemePattern defines a pattern for matching unofficial or less commonly
	// used URL schemes. These schemes may not be registered with IANA but are still valid in specific contexts,
	// such as application-specific schemes (e.g., "slack://", "zoommtg://").
	// The pattern ensures that the scheme is followed by "://".
	//
	// This pattern is useful for applications that work with unofficial or niche schemes.
	ExtractorKnownUnofficialSchemePattern = `(?:` + anyOf(schemes.Unofficial...) + `://)`

	// ExtractorKnownNoAuthoritySchemePattern defines a pattern for matching URL schemes that
	// do not require an authority component (host). These schemes are followed by a colon (":") rather than "://".
	// Examples include "mailto:", "tel:", and "sms:".
	//
	// This pattern is used for schemes where a host is not applicable, making it suitable for schemes
	// that involve direct communication (e.g., email or telephone).
	ExtractorKnownNoAuthoritySchemePattern = `(?:` + anyOf(schemes.NoAuthority...) + `:)`

	// ExtractorKnownSchemePattern combines the patterns for officially recognized, unofficial,
	// and no-authority-required schemes into a single comprehensive pattern.
	// It is case-insensitive (denoted by "(?i)") and matches the broadest possible range of URLs.
	//
	// This pattern is suitable for extracting any known scheme, regardless of its official status
	// or whether it requires an authority component.
	ExtractorKnownSchemePattern = `(?:(?i)(?:` + anyOf(schemes.Official...) + `|` + anyOf(schemes.Unofficial...) + `)://|` + anyOf(schemes.NoAuthority...) + `:)`

	// ExtractorIPv4Pattern defines a pattern for matching valid IPv4 addresses.
	// It matches four groups of 1 to 3 digits (0-255) separated by periods (e.g., "192.168.0.1").
	//
	// This pattern is essential for extracting or validating IPv4 addresses in URLs or hostnames.
	ExtractorIPv4Pattern = `(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])`

	// ExtractorNonEmptyIPv6Pattern defines a detailed pattern for matching valid, non-empty IPv6 addresses.
	// It accounts for various valid formats of IPv6 addresses, including those with elisions ("::") and IPv4
	// address representations.
	//
	// This pattern supports matching fully expanded IPv6 addresses, elided sections, and IPv4-mapped IPv6 addresses.
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

	// ExtractorIPv6Pattern is a comprehensive pattern that matches both fully expanded and compressed IPv6 addresses.
	// It also handles "::" elision and optional IPv4-mapped sections.
	ExtractorIPv6Pattern = `(?:` + ExtractorNonEmptyIPv6Pattern + `|::)`

	// ExtractorPortPattern defines a pattern for matching port numbers in URLs.
	// It matches valid port numbers (1 to 65535) that are typically found in network addresses.
	// The port number is preceded by a colon (":").
	ExtractorPortPattern = `(?::[0-9]{1,4}|[1-5][0-9]{4}|6[0-5][0-9]{3}\b)`

	// ExtractorPortOptionalPattern is similar to ExtractorPortPattern but makes the port number optional.
	// This is useful for matching URLs where the port may or may not be specified.
	ExtractorPortOptionalPattern = ExtractorPortPattern + `?`
)

// Ensure that Extractor implements the ExtractorInterface.
var _ ExtractorInterface = &Extractor{}

// NewExtractor creates a new Extractor instance with optional configuration.
// The options can be used to customize how URLs are extracted, such as whether
// to include URL schemes or hosts.
func NewExtractor(opts ...ExtractorOptionFunc) (extractor *Extractor) {
	extractor = &Extractor{}

	for _, opt := range opts {
		opt(extractor)
	}

	return
}

// ExtractorWithScheme returns an option function that configures the Extractor
// to require URL schemes in the extraction process.
func ExtractorWithScheme() ExtractorOptionFunc {
	return func(e *Extractor) {
		e.withScheme = true
	}
}

// ExtractorWithSchemePattern returns an option function that allows specifying
// a custom regex pattern for matching URL schemes.
func ExtractorWithSchemePattern(pattern string) ExtractorOptionFunc {
	return func(e *Extractor) {
		e.withScheme = true
		e.withSchemePattern = pattern
	}
}

// ExtractorWithHost returns an option function that configures the Extractor
// to require URL hosts in the extraction process.
func ExtractorWithHost() ExtractorOptionFunc {
	return func(e *Extractor) {
		e.withHost = true
	}
}

// ExtractorWithHostPattern returns an option function that allows specifying
// a custom regex pattern for matching URL hosts.
func ExtractorWithHostPattern(pattern string) ExtractorOptionFunc {
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
