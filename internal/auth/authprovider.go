package auth

import (
	"github.com/qiangxue/fasthttp-routing"
)

// Provider describes which functions
// an Authentication provider must provide.
type Provider interface {
	// Authenticate takes the context of a HTTP
	// request and returns an object which should
	// describe the authenticated user/account and
	// and error if the authentication failes.
	Authenticate(ctx *routing.Context) (interface{}, error)
}
