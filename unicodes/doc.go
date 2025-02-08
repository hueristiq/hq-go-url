// Package unicodes provides predefined sets of allowed Unicode character ranges.
// These ranges define characters that can be used in various contexts, including
// text processing, validation, and security filtering.
//
// This package is useful for applications that need to define and enforce
// character restrictions, such as usernames, identifiers, and text normalization.
//
// Features:
// - Defines a broad range of allowed Unicode characters for general text processing.
// - Provides a restricted subset that excludes certain punctuation marks.
// - Useful for applications involving security, data validation, and input sanitization.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "go.source.hueristiq.com/url/unicodes"
//	)
//
//	func main() {
//	    fmt.Println("Allowed Unicode Characters:", unicodes.AllowedUcsChar)
//	    fmt.Println("Allowed Unicode Characters (Minus Punctuation):", unicodes.AllowedUcsCharMinusPunc)
//	}
//
// References:
// - Unicode Character Database: https://www.unicode.org/ucd/
package unicodes
