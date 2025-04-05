// Package unicodes provides predefined sets of allowed Unicode character ranges.
// These ranges define characters that can be used in various contexts, including
// text processing, validation, and security filtering.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/hueristiq/hq-go-url/unicodes"
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
