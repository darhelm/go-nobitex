package utils

import (
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
)

// GenerateOtpCode returns a time-based one-time password (TOTP) generated
// from the provided secret. It uses the current system time to compute
// the code following the standard TOTP algorithm.
//
// Parameters:
//
//	otpSecret — the Base32-encoded secret key used to generate the TOTP.
//
// Returns:
//
//	The generated 6-digit OTP code as a string.
//	An error if the code cannot be generated.
//
// This function is typically used when interacting with services that
// require TOTP-based two-factor authentication.
func GenerateOtpCode(otpSecret string) (string, error) {
	code, err := totp.GenerateCode(otpSecret, time.Now())
	if err != nil {
		return "", fmt.Errorf("could not generate otp code: %v", err)
	}

	return code, nil
}
