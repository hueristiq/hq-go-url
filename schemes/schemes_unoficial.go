package schemes

// Unofficial is a sorted list of some well-known URL schemes that are not yet officially registered.
// These schemes are commonly associated with specific software or services, and while they may not be
// part of the official URI scheme registry, they are widely recognized and used.
//
// This list primarily gathers schemes from:
//   - https://en.wikipedia.org/wiki/List_of_URI_schemes#Unofficial_but_common_URI_schemes
//
// The schemes in this list are useful in specific contexts, such as handling custom protocols for
// software applications or services.
var Unofficial = []string{
	`gemini`,        // Gemini - a lightweight internet protocol for navigating and publishing on the web.
	`jdbc`,          // Java Database Connectivity (JDBC) - for connecting to databases from Java applications.
	`moz-extension`, // Firefox extension - used for accessing Firefox extensions.
	`postgres`,      // PostgreSQL (short form) - a short form of the PostgreSQL database scheme.
	`postgresql`,    // PostgreSQL - full form for PostgreSQL database connections.
	`slack`,         // Slack - used for handling Slack URIs.
	`zoommtg`,       // Zoom (desktop) - used by the Zoom desktop application.
	`zoomus`,        // Zoom (mobile) - used by the Zoom mobile application.
}
