package schemes

// NoAuthority is a sorted list of some well-known url schemes that are
// followed by ":" instead of "://". The list includes both officially
// registered and unofficial schemes.
var NoAuthority = []string{
	`bitcoin`, // Bitcoin
	`cid`,     // Content-ID
	`file`,    // Files
	`magnet`,  // Torrent magnets
	`mailto`,  // Mail
	`mid`,     // Message-ID
	`sms`,     // SMS
	`tel`,     // Telephone
	`xmpp`,    // XMPP
}
