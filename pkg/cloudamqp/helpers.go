package cloudamqp

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// doRequest encodes a failed HTTP response as a gRPC status whose code is the
// raw HTTP status code (see Client.doRequest). These helpers are the single
// place that classification of those errors lives, used by the connector layer
// for idempotency decisions (already-exists on create, not-found on delete).

// IsNotFoundError reports whether err represents an HTTP 404 (or a NotFound
// status produced by the client's own lookups).
func IsNotFoundError(err error) bool {
	switch status.Code(err) {
	case codes.NotFound, codes.Code(http.StatusNotFound):
		return true
	default:
		return false
	}
}

// IsAlreadyExistsError reports whether err represents an HTTP 409 Conflict, which
// CloudAMQP returns when the team member or invitation already exists. Only 409 is
// treated as already-exists: a 422 Unprocessable Entity is a genuine validation
// failure (e.g. a malformed email) and must surface as an error rather than be
// masked as "already exists".
func IsAlreadyExistsError(err error) bool {
	return status.Code(err) == codes.Code(http.StatusConflict)
}
