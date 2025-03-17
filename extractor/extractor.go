package extractor

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"go.source.hueristiq.com/url/schemes"
	"go.source.hueristiq.com/url/tlds"
	"go.source.hueristiq.com/url/unicodes"
)

// Extractor configures the URL extraction process.
// It provides options for controlling whether URL schemes and hosts are mandatory, and allows custom
// regular expression patterns to be specified for these components. This allows fine-grained control over
// the types of URLs that are extracted from text.
//
// Fields:
//   - withScheme: A boolean flag indicating if a URL scheme (e.g., "http") is required in extracted URLs.
//   - withSchemePattern: A custom regular expression pattern to match URL schemes.
//   - withHost: A boolean flag indicating if a URL host (e.g., domain) is required in extracted URLs.
//   - withHostPattern: A custom regular expression pattern to match URL hosts.
type Extractor struct {
	withScheme        bool
	withSchemePattern string
	withHost          bool
	withHostPattern   string
}

// CompileRegex constructs and compiles a regular expression pattern for URL extraction.
// It builds a composite regex pattern that can capture various URL forms, taking into account user
// configuration for required schemes or hosts, as well as any custom patterns provided.
// The generated pattern can match fully-qualified URLs (with schemes and hosts), email addresses,
// and relative URLs, depending on the configuration.
//
// Returns:
//   - regex (*regexp.Regexp): A compiled regular expression object that can be used to extract URLs from text.
//
// Example:
//
//	// Create an extractor that requires URL schemes.
//	extractor := New(WithScheme())
//	regex := extractor.CompileRegex()
//	urls := regex.FindAllString(text, -1) // Extracts URLs from the provided text.
func (e *Extractor) CompileRegex() (regex *regexp.Regexp) {
	schemePattern := ExtractorSchemePattern

	if e.withScheme && e.withSchemePattern != "" {
		schemePattern = e.withSchemePattern
	}

	// Separate ASCII TLDs (top-level domains) from Unicode TLDs for the regular expression.
	// We do this by checking for the first TLD that starts with a rune >= utf8.RuneSelf.
	var asciiTLDs, unicodeTLDs []string

	for i, tld := range tlds.Official {
		if tld[0] >= utf8.RuneSelf {
			asciiTLDs = tlds.Official[:i:i]
			unicodeTLDs = tlds.Official[i:]

			break
		}
	}

	// Define the punycode pattern used for internationalized domain names.
	punycode := `xn--[a-z0-9-]+`

	// knownTLDPattern combines punycode, known ASCII TLDs, pseudo TLDs (like .local, .onion),
	// and the set of unicode TLDs into a single pattern. We add a word boundary (\b) for ASCII TLDs
	// and keep the pattern case-insensitive with `(?i)`.
	knownTLDPattern := `(?:(?i)` + punycode + `|` + anyOf(append(asciiTLDs, tlds.Pseudo...)...) + `\b|` + anyOf(unicodeTLDs...) + `)`

	// domainPattern matches typical domains by combining subdomains, known TLDs, or 'localhost'.
	domainPattern := `(?:` + _subdomainPattern + knownTLDPattern + `|localhost)`

	// Host patterns:
	//   - hostWithoutPortPattern: matches domain, IPv6 (in brackets), or bare IPv4.
	//   - hostWithPortOptionalPattern: hostWithoutPort plus an optional port.
	hostWithoutPortPattern := `(?:` + domainPattern + `|\[` + ExtractorIPv6Pattern + `\]|\b` + ExtractorIPv4Pattern + `\b)`
	hostWithPortOptionalPattern := `(?:` + hostWithoutPortPattern + ExtractorPortOptionalPattern + `)`

	// If a custom host pattern was provided, use it.
	if e.withHost && e.withHostPattern != "" {
		hostWithPortOptionalPattern = e.withHostPattern
	}

	// _IAuthorityPattern: user info + host:port
	_IAuthorityPattern := `(?:` + _IUserInfoOptionalPattern + hostWithPortOptionalPattern + `)`

	// _IAuthorityOptionalPattern allows an optional authority.
	_IAuthorityOptionalPattern := _IAuthorityPattern + `?`

	// Patterns for the URL forms:
	//   - webURL:        has authority + optional path
	//   - email:         simplified email pattern, capturing user@host
	//   - URLsWithSchemePattern: fully supports an optional or mandatory authority, with path
	webURL := _IAuthorityPattern + `(?:/` + pathCont + `|/)?`
	email := `(?P<relaxedEmail>[a-zA-Z0-9._%\-+]+@` + hostWithPortOptionalPattern + `)`
	URLsWithSchemePattern := schemePattern + _IAuthorityOptionalPattern + pathCont

	// If a custom host pattern is provided, ensure we require the authority.
	if e.withHostPattern != "" {
		URLsWithSchemePattern = schemePattern + _IAuthorityPattern + `(?:/` + pathCont + `|/)?`
	}

	// URLsWithHostPattern tries to match either web-like URLs or emails
	// without a specific scheme.
	URLsWithHostPattern := webURL + `|` + email

	// RelativeURLsPattern captures relative paths (e.g. /path/to/resource or
	// something/like/this).
	RelativeURLsPattern := `(\/[\w\/?=&#.-]*)|([\w\/?=&#.-]+?(?:\/[\w\/?=&#.-]+)+)`

	// Build the final combined pattern depending on the configuration:
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

// WithScheme sets the Extractor to require URL schemes during extraction.
// It marks the withScheme flag as true.
func (e *Extractor) WithScheme() {
	e.withScheme = true
}

// WithSchemePattern sets a custom regular expression pattern for matching URL schemes.
// It also marks the Extractor as requiring a scheme.
//
// Parameter:
//   - pattern (string): A regex pattern to match URL schemes.
func (e *Extractor) WithSchemePattern(pattern string) {
	e.WithScheme()

	e.withSchemePattern = pattern
}

// WithHost sets the Extractor to require URL hosts during extraction.
// It marks the withHost flag as true.
func (e *Extractor) WithHost() {
	e.withHost = true
}

// WithHostPattern sets a custom regular expression pattern for matching URL hosts.
// It also marks the Extractor as requiring a host.
//
// Parameter:
//   - pattern (string): A regex pattern to match URL hosts.
func (e *Extractor) WithHostPattern(pattern string) {
	e.WithHost()

	e.withHostPattern = pattern
}

// Option defines a functional option for configuring an Extractor instance.
// It allows the caller to pass in configuration functions that modify the Extractor's
// settings (e.g., requiring a scheme or host, or providing custom regex patterns).
type Option func(extractor *Extractor)

// Interface defines the behavior that an Extractor must implement.
// It requires the implementation of a method to compile a regular expression for URL extraction.
type Interface interface {
	CompileRegex() (regex *regexp.Regexp)
}

const (
	// _alphaCharacterSet defines the set of alphabetical characters.
	_alphaCharacterSet = `a-zA-Z`

	// _digitCHaracterSet defines the set of digit characters.
	_digitCHaracterSet = `0-9`

	// _IUnreservedCharacterSet is the set of unreserved characters allowed in URL components,
	// including letters, digits, hyphens, dots, underscores, tildes, and additional allowed Unicode characters.
	_IUnreservedCharacterSet = _alphaCharacterSet + _digitCHaracterSet + `\-\._~` + unicodes.AllowedUcsChar

	// _IEndUnreservedCharacterSet is similar to _IUnreservedCharacterSet but excludes certain punctuation.
	_IEndUnreservedCharacterSet = _alphaCharacterSet + _digitCHaracterSet + `\-_~` + unicodes.AllowedUcsCharMinusPunc

	// _subDelimsCharacterSet defines the set of sub-delimiter characters allowed in URL components.
	_subDelimsCharacterSet = `!\$&'\(\)\*\+,;=`

	// _endSubDelimsCharacterSet defines a subset of sub-delimiters typically allowed at the end of URL components.
	_endSubDelimsCharacterSet = `\$&\+=`

	// _pctEncodingPattern defines the pattern for percent-encoded characters.
	_pctEncodingPattern = `%[0-9a-fA-F]{2}`

	// _IUserInfoPattern matches user information (e.g., username:password) in a URL.
	// It allows unreserved characters, sub-delimiters, colons, or percent-encoded characters, followed by an "@".
	_IUserInfoPattern = `(?:(?:[` + _IUnreservedCharacterSet + _subDelimsCharacterSet + `:]|` + _pctEncodingPattern + `)+@)`

	// _IUserInfoOptionalPattern makes the user information optional.
	_IUserInfoOptionalPattern = _IUserInfoPattern + `?`

	// midIPathSegmentChar defines the set of characters allowed in the middle of a URL path segment.
	midIPathSegmentChar = _IUnreservedCharacterSet + `%` + _subDelimsCharacterSet + `:@`

	// endIPathSegmentChar defines the set of characters allowed at the end of a URL path segment.
	endIPathSegmentChar = _IEndUnreservedCharacterSet + `%` + _endSubDelimsCharacterSet

	// _IPrivateCharacters defines Unicode private use areas that are allowed in URLs.
	_IPrivateCharacters = `\x{E000}-\x{F8FF}\x{F0000}-\x{FFFFD}\x{100000}-\x{10FFFD}`

	// midIChar is the set of characters allowed in the middle of a URL component.
	midIChar = `/?#\\` + midIPathSegmentChar + _IPrivateCharacters

	// endIChar is the set of characters allowed at the end of a URL component.
	endIChar = `/#` + endIPathSegmentChar + _IPrivateCharacters

	// wellParen, wellBrack, and wellBrace match well-formed nested parentheses, brackets, and braces,
	// respectively. wellAll combines these into a single pattern.
	wellParen = `\((?:[` + midIChar + `]|\([` + midIChar + `]*\))*\)`
	wellBrack = `\[(?:[` + midIChar + `]|\[[` + midIChar + `]*\])*\]`
	wellBrace = `\{(?:[` + midIChar + `]|\{[` + midIChar + `]*\})*\}`
	wellAll   = wellParen + `|` + wellBrack + `|` + wellBrace

	// pathCont defines the pattern for a URL path continuation. It matches segments composed of
	// allowed middle characters and ensures that the path does not abruptly end with disallowed characters.
	pathCont = `(?:[` + midIChar + `]*(?:` + wellAll + `|[` + endIChar + `]))+`

	// _letter, _mark, _number are Unicode property patterns for letters, marks, and numbers.
	_letter = `\p{L}`
	_mark   = `\p{M}`
	_number = `\p{N}`

	// _IRICharctersPattern defines a pattern for a valid IRI character sequence.
	// It matches a character with Unicode letter, mark, or number properties, possibly followed
	// by a hyphenated sequence.
	_IRICharctersPattern = `[` + _letter + _mark + _number + `](?:[` + _letter + _mark + _number + `\-]*[` + _letter + _mark + _number + `])?`

	// _subdomainPattern matches one or more IRI character sequences followed by a dot, forming a subdomain.
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

var _ Interface = (*Extractor)(nil)

// New creates a new Extractor instance with optional configuration.
// The options can be used to customize the behavior of the Extractor, such as requiring URL schemes
// or hosts, and specifying custom regex patterns for these components.
//
// Parameter:
//   - options (...Option): A variadic list of Option functions that modify the Extractor's settings.
//
// Returns:
//   - extractor (*Extractor): A pointer to the configured Extractor instance.
func New(options ...Option) (extractor *Extractor) {
	extractor = &Extractor{}

	for _, option := range options {
		option(extractor)
	}

	return
}

// WithScheme returns an Option function that configures the Extractor to require URL schemes in the extraction process.
//
// Returns:
//   - option (Option): A function that sets the withScheme flag to true.
func WithScheme() (option Option) {
	return func(e *Extractor) {
		e.WithScheme()
	}
}

// WithSchemePattern returns an Option function that sets a custom regex pattern for matching URL schemes.
//
// Parameter:
//   - pattern (string): A regex pattern to match URL schemes.
//
// Returns:
//   - option (Option): A function that sets the custom scheme pattern and enables the withScheme flag.
func WithSchemePattern(pattern string) (option Option) {
	return func(e *Extractor) {
		e.WithSchemePattern(pattern)
	}
}

// WithHost returns an Option function that configures the Extractor to require URL hosts in the extraction process.
//
// Returns:
//   - option (Option): A function that sets the withHost flag to true.
func WithHost() (option Option) {
	return func(e *Extractor) {
		e.WithHost()
	}
}

// WithHostPattern returns an Option function that sets a custom regex pattern for matching URL hosts.
//
// Parameter:
//   - pattern (string): A regex pattern to match URL hosts.
//
// Returns:
//   - option (Option): A function that sets the custom host pattern and enables the withHost flag.
func WithHostPattern(pattern string) (option Option) {
	return func(e *Extractor) {
		e.WithHostPattern(pattern)
	}
}

// anyOf is a helper function that constructs a non-capturing regex pattern from a list of strings.
// It joins the provided strings with a "|" (alternation) operator and escapes each string to ensure
// special regex characters are treated literally.
//
// Parameter:
//   - strs (...string): A variadic list of strings to be included in the pattern.
//
// Returns:
//   - string: A regex pattern that matches any one of the provided strings.
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
