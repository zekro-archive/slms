package auth

import (
	"errors"
	"strings"

	routing "github.com/qiangxue/fasthttp-routing"
)

var (
	// ErrUnauthorized indicates an unauthorized
	// request attempt
	ErrUnauthorized = errors.New("unauthorized")
)

// TokenAuthProvider provides an authentication
// method using basic header based authentication
// tokens.
type TokenAuthProvider struct {
	token string
}

// NewTokenAuthProvider creates a new instance
// of TokenAuthProvider passing the token to be
// used for authentication.
func NewTokenAuthProvider(token string) *TokenAuthProvider {
	return &TokenAuthProvider{
		token: token,
	}
}

// Authenticate checks the Authorization header
// for a Basic token and checks equality to the
// defined reference API token.
func (tap *TokenAuthProvider) Authenticate(ctx *routing.Context) (interface{}, error) {
	authVal := string(ctx.Request.Header.Peek("Authorization"))
	if authVal == "" || !strings.HasPrefix(strings.ToLower(authVal), "basic ") {
		return nil, ErrUnauthorized
	}

	authValSplit := strings.SplitN(authVal, " ", 2)
	if len(authValSplit) < 2 {
		return nil, ErrUnauthorized
	}

	if authValSplit[1] != tap.token {
		return nil, ErrUnauthorized
	}

	return "authorized", nil
}
