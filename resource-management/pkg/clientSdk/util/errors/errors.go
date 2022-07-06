package errors

import (
	"global-resource-service/resource-management/pkg/common-lib/types"
)

// ErrorReporter converts generic errors into runtime.Object errors without
// requiring the caller to take a dependency on meta/v1 (where Status lives).
// This prevents circular dependencies in core watch code.
type ErrorReporter struct {
	code   int
	verb   string
	reason string
}

// NewClientErrorReporter will respond with valid v1.Status objects that report
// unexpected server responses. Primarily used by watch to report errors when
// we attempt to decode a response from the server and it is not in the form
// we expect. Because watch is a dependency of the core api, we can't return
// meta/v1.Status in that package and so much inject this interface to convert a
// generic error as appropriate. The reason is passed as a unique status cause
// on the returned status, otherwise the generic "ClientError" is returned.
func NewClientErrorReporter(code int, verb string, reason string) *ErrorReporter {
	return &ErrorReporter{
		code:   code,
		verb:   verb,
		reason: reason,
	}
}

// TODO: since we dont have a status oject, so should expand the report the error.
//
func (r *ErrorReporter) AsObject(err error) *types.LogicalNode {
	return &types.LogicalNode{}
}
