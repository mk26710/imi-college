package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

const DEFAULT_TOKEN_SIZE int = 72

func NewToken(size int) (string, error) {
	bytes := make([]byte, size)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

var ErrTokenNotFound error = errors.New("user token is empty or not found")

// attempts to read auth token from the request's cookies
// and from Authorzation header in case there's no cookie
//
// returned token string will be only the token itself
// without any prefix, so "Bearer " part gets removed
//
// if there is no token or user provided invalid string
// the error value will be set to ErrTokenNotFound
func ExtractToken(r *http.Request) (string, error) {
	var inputToken string

	// attempt to read token from cookie first
	if cookie, err := r.Cookie("token"); err == nil {
		inputToken = cookie.Value
	}

	// if cookie wasn't read or empty attempt reading header
	if header := r.Header.Get("Authorization"); len(header) > 0 && len(inputToken) == 0 {
		inputToken = header
	}

	// if token is empty then there's no token
	if len(inputToken) == 0 {
		return "", ErrTokenNotFound
	}

	// make sure token has prefix and cut it
	rawToken, isCut := strings.CutPrefix(inputToken, "Bearer ")
	if !isCut {
		return "", ErrTokenNotFound
	}

	return rawToken, nil
}
