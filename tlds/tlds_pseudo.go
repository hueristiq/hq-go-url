package tlds

// Pseudo is a sorted list of some widely used unofficial TLDs.
// The list is fetched from:
//   - https://en.wikipedia.org/wiki/Pseudo-top-level_domain
//   - https://en.wikipedia.org/wiki/Category:Pseudo-top-level_domains
//   - https://tools.ietf.org/html/draft-grothoff-iesg-special-use-p2p-names-00
//   - https://www.iana.org/assignments/special-use-domain-names/special-use-domain-names.xhtml
var Pseudo = []string{
	`bit`,       // Namecoin
	`example`,   // Example domain
	`exit`,      // Tor exit node
	`gnu`,       // GNS by public key
	`i2p`,       // I2P network
	`invalid`,   // Invalid domain
	`local`,     // Local network
	`localhost`, // Local network
	`test`,      // Test domain
	`zkey`,      // GNS domain name
}
