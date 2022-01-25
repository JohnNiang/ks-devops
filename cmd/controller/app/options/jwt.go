package options

import (
	authoptions "kubesphere.io/devops/pkg/apiserver/authentication/options"
	"time"
)

// JWTOptions contain some options of JWT, such as secret and clock skew.
type JWTOptions struct {
	// MaximumClockSkew indicates token verification maximum time difference.
	MaximumClockSkew time.Duration
	// Secret is used to sign JWT token.
	Secret string
}

func NewJWTOptions(options *authoptions.AuthenticationOptions) *JWTOptions {
	if options == nil {
		return nil
	}
	return &JWTOptions{
		Secret:           options.JwtSecret,
		MaximumClockSkew: options.MaximumClockSkew,
	}
}
