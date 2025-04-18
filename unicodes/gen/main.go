package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
	"unicode"
)

var (
	unicodesTemplate = template.Must(template.New("schemes").Parse(`// This file is autogenerated by ./gen/main.go. Please do not edit manually.
package unicodes

// AllowedUcsChar defines a range of allowed Unicode characters.
//
// This set includes various characters spanning multiple Unicode blocks.
// It supports a wide range of characters, including those from different languages,
// symbols, and select punctuation marks.
const AllowedUcsChar = {{.withPunc}}

// AllowedUcsCharMinusPunc defines a range of allowed Unicode characters,
// excluding certain punctuation marks.
//
// This set is used in contexts where punctuation is restricted, but other characters
// from AllowedUcsChar are allowed. This is useful for filtering input in usernames,
// identifiers, or text fields that should not contain punctuation.
const AllowedUcsCharMinusPunc = {{.withoutPunc}}
`))
	unicodesFile string
)

func init() {
	flag.StringVar(&unicodesFile, "output", "", "Specify the output file path for the generated Go source file.")

	flag.Usage = func() {
		h := "USAGE:\n"
		h += "  schemes [OPTIONS]\n"

		h += "\nOPTIONS:\n"
		h += " -output string    Specify the output file path for the generated Go source file.\n"

		fmt.Fprintln(os.Stderr, h)
	}

	flag.Parse()
}

func main() {
	if unicodesFile == "" {
		log.Fatalln("Output file path is required. Use -output to specify the output file path.")
	}

	log.Printf("Generating %s...\n", unicodesFile)

	if err := writeUnicode(); err != nil {
		log.Fatal(err)
	}

	log.Println("Unicodes file generated successfully.")
}

// visit iterates through the given Unicode range table and applies the given function `fn` to each rune.
func visit(rt *unicode.RangeTable, fn func(rune)) {
	for _, r16 := range rt.R16 {
		for r := rune(r16.Lo); r <= rune(r16.Hi); r += rune(r16.Stride) {
			fn(r)
		}
	}

	for _, r32 := range rt.R32 {
		for r := rune(r32.Lo); r <= rune(r32.Hi); r += rune(r32.Stride) {
			fn(r)
		}
	}
}

// writeUnicode generates the Unicode ranges and writes them to a Go file using the provided template.
func writeUnicode() error {
	// rfc3987Ranges contains the ranges of valid code points specified by RFC 3987.
	rfc3987Ranges := [][2]rune{
		{0xA0, 0xD7FF},
		{0xF900, 0xFDCF},
		{0xFDF0, 0xFFEF},
		{0x10000, 0x1FFFD},
		{0x20000, 0x2FFFD},
		{0x30000, 0x3FFFD},
		{0x40000, 0x4FFFD},
		{0x50000, 0x5FFFD},
		{0x60000, 0x6FFFD},
		{0x70000, 0x7FFFD},
		{0x80000, 0x8FFFD},
		{0x90000, 0x9FFFD},
		{0xA0000, 0xAFFFD},
		{0xB0000, 0xBFFFD},
		{0xC0000, 0xCFFFD},
		{0xD0000, 0xDFFFD},
		{0xE1000, 0xEFFFD},
	}

	// removeRune accepts a slice of inclusive code point ranges (in ascending order)
	// and returns a new slice that is equivalent except for excluding a specified rune
	// by removing/replacing/splitting any range containing it.
	// Its linear searches over the ranges (including those added by previous invocations)
	// are inefficient, but acceptable because this code runs only at build time.
	removeRune := func(ranges [][2]rune, cp rune) [][2]rune {
		for i, r := range ranges {
			// Ranges are in ascending order. Skip any that precede `cp`,
			// and bail out upon reaching one that follows `cp`.
			if r[1] < cp {
				continue
			} else if cp < r[0] {
				break
			}

			// Modify the range to exclude the specified rune.
			// `cp` is in this range and must be removed from it.
			if cp == r[0] && cp == r[1] {
				// Remove this single-element range.
				return append(ranges[0:i], ranges[i+1:]...)
			} else if cp == r[0] {
				// Remove the first element of this range.
				newRange := [2]rune{r[0] + 1, r[1]}
				newTail := append([][2]rune{newRange}, ranges[i+1:]...)

				return append(ranges[0:i], newTail...)
			} else if cp == r[1] {
				// Remove the last element of this range.
				newRange := [2]rune{r[0], r[1] - 1}
				newTail := append([][2]rune{newRange}, ranges[i+1:]...)

				return append(ranges[0:i], newTail...)
			} else {
				// Split this range.
				newTail := append(
					[][2]rune{
						{r[0], cp - 1},
						{cp + 1, r[1]},
					},
					ranges[i+1:]...)

				return append(ranges[0:i], newTail...)
			}
		}

		return ranges
	}

	// sepFreeRanges excludes separators from rfc3987Ranges.
	// sepFreeRanges contains the ranges excluding separator characters (from Unicode category Z).
	sepFreeRanges := append([][2]rune{}, rfc3987Ranges...)

	visit(unicode.Z, func(cp rune) {
		sepFreeRanges = removeRune(sepFreeRanges, cp)
	})

	// puncFreeRanges excludes punctuation from sepFreeRanges.
	// puncFreeRanges contains the ranges excluding punctuation characters (from Unicode category Po).
	puncFreeRanges := append([][2]rune{}, sepFreeRanges...)

	visit(unicode.Po, func(cp rune) {
		puncFreeRanges = removeRune(puncFreeRanges, cp)
	})

	// Build the corresponding regular expression character class contents.
	// characterClassContents builds a string representing the character ranges for regular expression purposes.
	characterClassContents := func(ranges [][2]rune) strings.Builder {
		var builder strings.Builder

		for _, r := range ranges {
			// regexp.QuoteMeta is not necessary because all metacharacters are ASCII.
			// cf. https://golang.org/s/re2syntax and
			// https://cs.opensource.google/go/go/+/refs/tags/go1.17.6:src/regexp/regexp.go;l=721
			builder.WriteRune(r[0])

			if r[0] == r[1] {
				continue
			}

			builder.WriteRune('-')
			builder.WriteRune(r[1])
		}

		return builder
	}

	// Generate the allowed character sets.
	allowedUcsChar := characterClassContents(sepFreeRanges)
	allowedUcsCharMinusPunc := characterClassContents(puncFreeRanges)

	f, err := os.Create(unicodesFile)
	if err != nil {
		return err
	}

	defer f.Close()

	return unicodesTemplate.Execute(f, map[string]string{
		"withPunc":    strconv.Quote(allowedUcsChar.String()),
		"withoutPunc": strconv.Quote(allowedUcsCharMinusPunc.String()),
	})
}
