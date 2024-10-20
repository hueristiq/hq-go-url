package tlds

// Pseudo is a sorted list of widely used unofficial or pseudo top-level domains (TLDs).
// Pseudo-TLDs are domain name suffixes that function similarly to regular top-level domains
// (such as .com, .org), but are not part of the official Internet Assigned Numbers Authority (IANA) TLDs.
// These domains are often used in specific contexts such as private networks, experimental networks, or
// special-use domains.
//
// The list of pseudo-TLDs is curated from the following sources:
//   - https://en.wikipedia.org/wiki/Pseudo-top-level_domain
//   - https://en.wikipedia.org/wiki/Category:Pseudo-top-level_domains
//   - https://tools.ietf.org/html/draft-grothoff-iesg-special-use-p2p-names-00
//   - https://www.iana.org/assignments/special-use-domain-names/special-use-domain-names.xhtml
//
// Each pseudo-TLD in this list serves a specific purpose or has historical significance in its respective
// network or environment.
var Pseudo = []string{
	`bit`,       // Namecoin - a decentralized domain system based on the Namecoin blockchain.
	`example`,   // Example domain - reserved for use in documentation and examples.
	`exit`,      // Tor exit node - used for identifying Tor exit nodes in the Tor network.
	`gnu`,       // GNS by public key - GNU Name System, a decentralized name system.
	`i2p`,       // I2P network - Invisible Internet Project, an anonymous network layer.
	`invalid`,   // Invalid domain - reserved for invalid domain names.
	`local`,     // Local network - used in local networking environments.
	`localhost`, // Local network - refers to the local loopback interface (127.0.0.1).
	`test`,      // Test domain - reserved for use in testing environments.
	`zkey`,      // GNS domain name - used in the GNU Name System for public-key based domain names.
}
