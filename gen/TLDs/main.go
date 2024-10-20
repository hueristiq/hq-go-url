package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

var (
	// Output file path for the generated Go source file.
	output string

	// Template for the autogenerated Go file containing the list of TLDs.
	tmpl = template.Must(template.New("schemes").Parse(`// This file is autogenerated by the TLDs generator. Please do not edit manually.
package tlds

// Official is a sorted list of public top-level domains (TLDs) and effective top-level domains (eTLDs).
// TLDs are the highest level in the hierarchical domain name system of the Internet. eTLDs include
// top-level domains and public suffixes, such as country code second-level domains (e.g., "co.uk" or "gov.in"),
// that are commonly used for websites.
//
// The list is curated from official sources:
//   - https://data.iana.org/TLD/tlds-alpha-by-domain.txt: Contains a list of all current IANA-approved TLDs.
//   - https://publicsuffix.org/list/public_suffix_list.dat: Contains a list of public suffixes managed by the Public Suffix List,
//     which identifies domain suffixes under which Internet users can register names.
//
// This list is automatically generated to ensure it stays up to date with the latest TLDs and public suffixes.
var Official = []string{
{{- range $_, $TLD := .TLDs}}
	"{{$TLD}}",
{{- end}}
}
`))
)

func init() {
	// Define the command-line flag for output file path
	flag.StringVar(&output, "output", "", "Specify the output file path for the generated Go source file.")

	// Custom usage message for the command-line flag
	flag.Usage = func() {
		h := "USAGE:\n"
		h += "  schemes [OPTIONS]\n"

		h += "\nOPTIONS:\n"
		h += " -output string    Specify the output file path for the generated Go source file.\n"

		fmt.Fprintln(os.Stderr, h)
	}

	// Parse command-line flags
	flag.Parse()
}

func main() {
	// Ensure that an output file path is specified
	if output == "" {
		log.Fatalln("Output file path is required. Use -output to specify the output file path.")
	}

	log.Printf("Generating %s...\n", output)

	// Fetch TLDs from IANA
	TLDs, err := getTLDsFromIANA()
	if err != nil {
		log.Fatalf("Failed to get TLDs from IANA: %v\n", err)
	}

	// Fetch effective TLDs from the Public Suffix list
	eTLDs, err := getEffectiveTLDsFromPublicSuffix()
	if err != nil {
		log.Fatalf("Failed to get effective TLDs from Public Suffix: %v\n", err)
	}

	// Combine both TLDs and eTLDs
	TLDs = append(TLDs, eTLDs...)

	// Sort the combined list of TLDs
	sort.Strings(TLDs)

	// Remove duplicate entries
	TLDs = removeDuplicates(TLDs)

	// Write the TLDs to the output file
	if err := writeTLDsToFile(TLDs, output); err != nil {
		log.Fatalf("Failed to write schemes to file: %v\n", err)
	}

	log.Println("TLDs file generated successfully.")
}

// getTLDsFromIANA fetches the list of TLDs from the IANA TLD list and returns them.
func getTLDsFromIANA() (TLDs []string, err error) {
	// Perform HTTP GET request to fetch the IANA TLD list
	var res *http.Response

	res, err = http.Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
	if err != nil {
		err = fmt.Errorf("failed to fetch IANA TLDs: %w", err)

		return
	}

	defer res.Body.Close()

	// Regular expression to match valid TLD entries (ignore comments)
	re := regexp.MustCompile(`^[^#]+$`)

	// Scan through the response body line by line
	scanner := bufio.NewScanner(res.Body)

	for scanner.Scan() {
		line := scanner.Text()

		line = strings.TrimSpace(line)
		line = strings.ToLower(line)

		// Extract valid TLDs (skip comments and entries starting with "xn--")
		TLD := re.FindString(line)

		if TLD == "" || strings.HasPrefix(TLD, "xn--") {
			continue
		}

		TLDs = append(TLDs, TLD)
	}

	// Check for errors during scanning
	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("scanner error: %w", err)

		return
	}

	return
}

func getEffectiveTLDsFromPublicSuffix() (eTLDs []string, err error) {
	// Perform HTTP GET request to fetch the Public Suffix list
	var res *http.Response

	res, err = http.Get("https://publicsuffix.org/list/effective_tld_names.dat")
	if err != nil {
		err = fmt.Errorf("failed to fetch Public Suffix TLDs: %w", err)

		return
	}

	defer res.Body.Close()

	// Scan through the response body line by line
	scanner := bufio.NewScanner(res.Body)

	for scanner.Scan() {
		line := scanner.Text()

		line = strings.TrimSpace(line)

		// Stop reading when encountering private domain section
		if strings.HasPrefix(line, "// ===BEGIN PRIVATE DOMAINS") {
			break
		}

		// Skip comments
		if strings.HasPrefix(line, "//") {
			continue
		}

		TLD := line

		// Remove special characters
		TLD = strings.ReplaceAll(TLD, "*.", "")
		TLD = strings.ReplaceAll(TLD, "!", "")

		if TLD == "" {
			continue
		}

		eTLDs = append(eTLDs, TLD)
	}

	// Check for errors during scanning
	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("scanner error: %w", err)

		return
	}

	return
}

// removeDuplicates
// removes duplicate elements from a slice of any type that satisfies the comparable constraint.
func removeDuplicates[T comparable](slice []T) []T {
	keys := make(map[T]bool)

	var list []T

	for _, entry := range slice {
		if _, exists := keys[entry]; !exists {
			keys[entry] = true

			list = append(list, entry)
		}
	}

	return list
}

// writeTLDsToFile writes the generated list of URI schemes to the specified file
// using a Go source file template.
func writeTLDsToFile(TLDs []string, output string) (err error) {
	// Create the output file
	file, err := os.Create(output)
	if err != nil {
		err = fmt.Errorf("failed to create output file: %w", err)

		return
	}

	defer file.Close()

	// Execute the template and write to the output file
	data := struct {
		TLDs []string
	}{
		TLDs: TLDs,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return
}
