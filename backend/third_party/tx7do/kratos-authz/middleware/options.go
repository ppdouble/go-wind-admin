package middleware

import (
	"go-wind-admin/pkg/v2logcompat"
)

type Option func(*options)

type options struct {
	log *log.Helper
}

func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.log = log.NewHelper(log.With(logger, "module", "authz.middleware"))
	}
}
