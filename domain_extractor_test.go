package url_test

import (
	"testing"

	hqgourl "github.com/hueristiq/hq-go-url"
	"github.com/hueristiq/hq-go-url/tlds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainExtractor_CompileRegex_Default(t *testing.T) {
	t.Parallel()

	// Initialize DomainExtractor with default settings.
	extractor := hqgourl.NewDomainExtractor()

	// Compile the regex.
	regex := extractor.CompileRegex()

	// Ensure the regex is not nil.
	require.NotNil(t, regex)

	// Test that the regex matches valid domain patterns.
	tests := []struct {
		input    string
		expected bool
	}{
		{"example.com", true},
		{"www.example.com", true},
		{"http://www.example.com", true},
		{"localhost", true},
		{"http://localhost", true},
		{"example.localhost", true},
		{"http://example.localhost", true},
		{"xn--example-q9a.com", true}, // IDN with punycode.
		{"example.co.uk", true},
		{"invalid_domain", false},
		{"just_text", false},
		{"ftp://example.com", true},
	}

	for _, tt := range tests {
		assert.Equalf(t, tt.expected, regex.MatchString(tt.input), "failed on input: %s", tt.input)
	}
}

func TestDomainExtractor_CompileRegex_CustomRootDomainPattern(t *testing.T) {
	t.Parallel()

	// Initialize DomainExtractor with a custom root domain pattern.
	extractor := hqgourl.NewDomainExtractor(
		hqgourl.DomainExtractorWithRootDomainPattern(`(?:example|rootdomain)`), // Custom root domain pattern
	)

	// Compile the regex.
	regex := extractor.CompileRegex()

	// Ensure the regex is not nil.
	require.NotNil(t, regex)

	// Test cases for custom root domain pattern.
	tests := []struct {
		input    string
		expected bool
	}{
		{"rootdomain.com", true},
		{"my-root-domain.org", false},
		{"not_valid_domain", false},
		{"example.com", true},
		{"www.example.com", true},
		{"localhost", false},
	}

	for _, tt := range tests {
		assert.Equalf(t, tt.expected, regex.MatchString(tt.input), "failed on input: %s", tt.input)
	}
}

func TestDomainExtractor_CompileRegex_CustomTLDPattern(t *testing.T) {
	t.Parallel()

	// Initialize DomainExtractor with a custom TLD pattern.
	extractor := hqgourl.NewDomainExtractor(
		hqgourl.DomainExtractorWithTLDPattern(`(?:com|net|org)`), // Custom TLD pattern
	)

	// Compile the regex.
	regex := extractor.CompileRegex()

	// Ensure the regex is not nil.
	require.NotNil(t, regex)

	// Test cases for custom TLD pattern.
	tests := []struct {
		input    string
		expected bool
	}{
		{"example.com", true},
		{"example.org", true},
		{"example.net", true},
		{"example.co.uk", false},
		{"localhost", false},
	}

	for _, tt := range tests {
		assert.Equalf(t, tt.expected, regex.MatchString(tt.input), "failed on input: %s", tt.input)
	}
}

func TestDomainExtractor_CompileRegex_CustomRootDomainAndTLDPattern(t *testing.T) {
	t.Parallel()

	// Initialize DomainExtractor with custom root domain and TLD patterns.
	extractor := hqgourl.NewDomainExtractor(
		hqgourl.DomainExtractorWithRootDomainPattern(`[a-zA-Z0-9-]+`),
		hqgourl.DomainExtractorWithTLDPattern(`(?:com|net)`),
	)

	// Compile the regex.
	regex := extractor.CompileRegex()

	// Ensure the regex is not nil.
	require.NotNil(t, regex)

	// Test cases for custom root domain and TLD pattern.
	tests := []struct {
		input    string
		expected bool
	}{
		{"example.com", true},
		{"example.net", true},
		{"example.org", false}, // TLD pattern restricts to com/net.
		{"localhost", false},
		{"subdomain.example.com", true},
	}

	for _, tt := range tests {
		assert.Equalf(t, tt.expected, regex.MatchString(tt.input), "failed on input: %s", tt.input)
	}
}

func TestDomainExtractor_CompileRegex_TLDSeparation(t *testing.T) {
	t.Parallel()

	// Simulate a scenario where the TLDs include both ASCII and Unicode values.
	originalTLDs := tlds.Official
	tlds.Official = []string{"com", "org", "xn--unicode", "测试"}

	// Initialize the DomainExtractor.
	extractor := hqgourl.NewDomainExtractor()

	// Compile the regex.
	regex := extractor.CompileRegex()

	// Restore the original TLD list.
	t.Cleanup(func() { tlds.Official = originalTLDs })

	// Ensure the regex is not nil.
	require.NotNil(t, regex)

	// Test cases for ASCII and Unicode TLDs.
	tests := []struct {
		input    string
		expected bool
	}{
		{"example.com", true},
		{"example.org", true},
		{"example.测试", true},          // Unicode TLD.
		{"example.xn--unicode", true}, // Punycode.
		{"example.co.uk", false},      // TLD not in the list.
		{"localhost", true},
	}

	for _, tt := range tests {
		assert.Equalf(t, tt.expected, regex.MatchString(tt.input), "failed on input: %s", tt.input)
	}
}

func TestDomainExtractor_CustomPatterns_Failures(t *testing.T) {
	t.Parallel()

	// Test with invalid root domain and TLD patterns.
	extractor := hqgourl.NewDomainExtractor(
		hqgourl.DomainExtractorWithRootDomainPattern(`[invalid`), // Invalid regex pattern
		hqgourl.DomainExtractorWithTLDPattern(`(`),               // Invalid regex pattern
	)

	// Expecting a panic due to invalid regex patterns.
	assert.Panics(t, func() {
		extractor.CompileRegex()
	}, "Expected panic with invalid regex patterns")
}

func TestDomainExtractor_CustomPatterns_Empty(t *testing.T) {
	t.Parallel()

	// Test with empty custom root domain and TLD patterns.
	extractor := hqgourl.NewDomainExtractor(
		hqgourl.DomainExtractorWithRootDomainPattern(""),
		hqgourl.DomainExtractorWithTLDPattern(""),
	)

	// Compile the regex.
	regex := extractor.CompileRegex()

	// Ensure the regex is not nil.
	require.NotNil(t, regex)

	// The regex should fall back to default behavior.
	tests := []struct {
		input    string
		expected bool
	}{
		{"example.com", true},
		{"localhost", true},
		{"invalid_domain", false},
	}

	for _, tt := range tests {
		assert.Equalf(t, tt.expected, regex.MatchString(tt.input), "failed on input: %s", tt.input)
	}
}
