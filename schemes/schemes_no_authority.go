package schemes

// NoAuthority is a sorted list of well-known URL schemes that do not require an authority component.
// Instead of being followed by "://", these schemes are followed directly by a colon (":").
// This format is typically used when the scheme does not require a host or authority part in the URL.
//
// The list includes both officially registered schemes and commonly used unofficial schemes.
//
// These schemes are used in various contexts, such as addressing specific resources or services
// (e.g., email, telephone, file access, etc.).
var NoAuthority = []string{
	`bitcoin`, // Bitcoin - Used to create Bitcoin payment URIs.
	`cid`,     // Content-ID - Identifies a specific piece of content.
	`file`,    // Files - Refers to files on the local filesystem.
	`magnet`,  // Torrent magnets - Used for referencing torrents without relying on a tracker.
	`mailto`,  // Mail - Refers to email addresses for sending mail.
	`mid`,     // Message-ID - Refers to email message IDs.
	`sms`,     // SMS - Used for sending SMS messages.
	`tel`,     // Telephone - Refers to telephone numbers.
	`xmpp`,    // XMPP - Used for addressing XMPP (Jabber) communication services.
}
