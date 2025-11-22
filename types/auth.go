package types

// Captcha specifies the captcha mode used during authentication.
// API clients must set this to "api" when performing programmatic login.
type Captcha string

const API Captcha = "api"

// Remember indicates whether a long-lived or short-lived authentication
// token should be issued during login.
type Remember string

const (
	RememberYes Remember = "yes" // long-lived key
	RememberNo  Remember = "no"  // short-lived key
)

// AuthenticationParams defines the request payload and required headers
// for initiating a Nobitex authentication session.
type AuthenticationParams struct {
	// Username is the account email used to log in.
	Username string `json:"username"`

	// Password is the account password associated with the username.
	Password string `json:"password"`

	// Remember determines whether the server returns a long-lived token.
	Remember Remember `json:"remember"`

	// Captcha must be set to "api" for bot authentication flows.
	Captcha Captcha `json:"captcha"`

	// XTOTP is the 6-digit TOTP value sent through the X-TOTP header.
	// This is required when Captcha=="api".
	XTOTP string `json:"X-TOTP,omitempty"`

	// UserAgent identifies the application performing authentication.
	// Example: "TraderBot/1.2.0"
	UserAgent string `json:"User-Agent"`
}

// AuthenticationResponse contains the issued session token and metadata
// returned after a successful authentication request.
type AuthenticationResponse struct {
	// Status indicates the result of the login attempt (commonly "ok").
	Status string `json:"status"`

	// Key is the session token returned by the server. This must be provided
	// as `Authorization: Token <key>` for authenticated requests.
	Key string `json:"key"`

	// Device represents a unique ID assigned to the authenticated session.
	Device string `json:"device"`
}
