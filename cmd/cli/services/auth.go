package services

import (
	"errors"
	"fmt"
	"os"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/clients"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/responses"

	"github.com/NorskHelsenett/ror/pkg/rlog"
)

var ErrCouldNotReachAuth = errors.New("failed to reach dex")
var ErrAccessTokenIsEmpty = errors.New("access token is empty")
var ErrAccessTokenExpired = errors.New("access token has expired")
var ErrAccesTokenFailedToVerify = errors.New("failed to verify accesstoken")

type Claims struct {
	Email           string   `json:"email"`
	IsEmailVerified bool     `json:"email_verified"`
	Name            string   `json:"name"`
	Groups          []string `json:"groups"`
	Audience        string   `json:"aud"`
	Issuer          string   `json:"iss"`
	ExpirationTime  int64    `json:"exp"`
}

// GetJWToken get a JWToken from dex
func GetJWToken(code *responses.CodeResponse) *responses.TokenResponse {
	JWToken, err := clients.FetchJWToken(code)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "there was an error authenticating")
		rlog.Fatal("fetching JWToken failed: ", err)
	}

	return JWToken
}

// GetDeviceCode NOTE: this is not just the device code, should find a better name
func GetDeviceCode() (*responses.CodeResponse, error) {
	err, code := clients.FetchDeviceCode()
	if err != nil {
		return nil, err
	}

	return code, nil
}
