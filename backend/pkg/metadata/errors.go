package metadata

import "github.com/go-kratos/kratos/v3/errors"

const (
	reason string = "UNAUTHORIZED"
)

var (
	ErrNoMetadata       = errors.Unauthorized(reason, "metadata: missing server metadata")
	ErrNoOperatorHeader = errors.Unauthorized(reason, "metadata: missing operator header")
	ErrInvalidOperator  = errors.Unauthorized(reason, "metadata: invalid operator header")
)
