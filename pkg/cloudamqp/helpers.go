package cloudamqp

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// doRequest encodes a failed HTTP response as a gRPC status whose code is the
// raw HTTP status code (see Client.doRequest). These helpers are the single
// place that classification of those errors lives, used by the connector layer
// for idempotency decisions (already-exists on create, not-found on delete).

// maxAPIErrorBodyBytes bounds how much of a failed response body newAPIError
// reads. CloudAMQP's error bodies are a short {"error":"..."} object; this
// just keeps a misbehaving upstream from forcing an unbounded read.
const maxAPIErrorBodyBytes = 1 << 16

// alreadyInvitedMsgFragment is the substring CloudAMQP's error message carries
// for a duplicate POST /team/invite (see IsAlreadyInvitedError).
const alreadyInvitedMsgFragment = "already invited"

// apiErrorBody is the shape CloudAMQP uses for error responses, e.g.
// {"error":"User is already invited"}.
type apiErrorBody struct {
	Error string `json:"error"`
}

// newAPIError builds a gRPC status error from a failed HTTP response, using the
// raw status code and, when present, the vendor's JSON error message. The
// message matters here: CloudAMQP answers both a duplicate invite and a
// malformed email with the same HTTP 400, distinguishable only by this body
// (see IsAlreadyInvitedError). If the body is missing or unparseable, fall back
// to a generic message rather than fail the request over it.
func newAPIError(statusCode int, body io.Reader) error {
	msg := "Request failed"
	if raw, err := io.ReadAll(io.LimitReader(body, maxAPIErrorBodyBytes)); err == nil && len(raw) > 0 {
		var parsed apiErrorBody
		if json.Unmarshal(raw, &parsed) == nil && parsed.Error != "" {
			msg = parsed.Error
		}
	}
	//nolint:gosec // statusCode is always a valid HTTP code, so the int-to-uint32 conversion cannot overflow or wrap
	return status.Error(codes.Code(uint32(statusCode)), msg)
}

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

// IsAlreadyInvitedError reports whether err represents CloudAMQP's response to
// a duplicate POST /team/invite. Unlike member/invitation conflicts elsewhere
// in the API, the invite endpoint never returns 409 for this case: both a
// duplicate invite and a malformed email come back as HTTP 400, distinguished
// only by the JSON error body ("User is already invited" vs "Invalid email").
// Matching on the raw HTTP 400 alone would misclassify a genuine validation
// failure as already-invited, so the message must be checked too.
func IsAlreadyInvitedError(err error) bool {
	if status.Code(err) != codes.Code(http.StatusBadRequest) {
		return false
	}
	return strings.Contains(strings.ToLower(status.Convert(err).Message()), alreadyInvitedMsgFragment)
}
