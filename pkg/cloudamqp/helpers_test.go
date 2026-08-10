package cloudamqp

import (
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestNewAPIError pins that the vendor's JSON error body becomes the status
// message, since IsAlreadyInvitedError depends on that message to distinguish
// a duplicate invite from other HTTP 400s.
func TestNewAPIError(t *testing.T) {
	for _, tt := range []struct {
		name    string
		code    int
		body    string
		wantMsg string
	}{
		{name: "already invited", code: 400, body: `{"error":"User is already invited"}`, wantMsg: "User is already invited"},
		{name: "invalid email", code: 400, body: `{"error":"Invalid email"}`, wantMsg: "Invalid email"},
		{name: "empty body falls back", code: 500, body: "", wantMsg: "Request failed"},
		{name: "unparseable body falls back", code: 500, body: "not json", wantMsg: "Request failed"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := newAPIError(tt.code, strings.NewReader(tt.body))
			if status.Code(err) != codes.Code(uint32(tt.code)) { //nolint:gosec
				t.Fatalf("code = %v, want %d", status.Code(err), tt.code)
			}
			if status.Convert(err).Message() != tt.wantMsg {
				t.Fatalf("message = %q, want %q", status.Convert(err).Message(), tt.wantMsg)
			}
		})
	}
}

func TestIsAlreadyInvitedError(t *testing.T) {
	for _, tt := range []struct {
		name string
		err  error
		want bool
	}{
		{name: "duplicate invite", err: newAPIError(400, strings.NewReader(`{"error":"User is already invited"}`)), want: true},
		{name: "invalid email is not a duplicate", err: newAPIError(400, strings.NewReader(`{"error":"Invalid email"}`)), want: false},
		{name: "409 conflict is not this error", err: newAPIError(409, strings.NewReader(`{"error":"User is already invited"}`)), want: false},
		{name: "nil error", err: nil, want: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlreadyInvitedError(tt.err); got != tt.want {
				t.Fatalf("IsAlreadyInvitedError() = %v, want %v", got, tt.want)
			}
		})
	}
}
