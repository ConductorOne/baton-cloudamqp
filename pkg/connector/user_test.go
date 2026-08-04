package connector

import (
	"testing"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestResolveInviteEmail(t *testing.T) {
	email := func(addr string, primary bool) *v2.AccountInfo_Email {
		return &v2.AccountInfo_Email{Address: addr, IsPrimary: primary}
	}

	for _, tt := range []struct {
		name    string
		info    *v2.AccountInfo
		profile map[string]interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "username login falls back to profile email",
			info:    &v2.AccountInfo{Login: "test"},
			profile: map[string]interface{}{"email": "user@example.com"},
			want:    "user@example.com",
		},
		{
			name: "primary email wins over login and profile",
			info: &v2.AccountInfo{
				Login:  "someone@else.com",
				Emails: []*v2.AccountInfo_Email{email("alt@example.com", false), email("primary@example.com", true)},
			},
			profile: map[string]interface{}{"email": "profile@example.com"},
			want:    "primary@example.com",
		},
		{
			name:    "non-primary email used when no primary is marked",
			info:    &v2.AccountInfo{Emails: []*v2.AccountInfo_Email{email("alt@example.com", false)}},
			profile: map[string]interface{}{"email": "profile@example.com"},
			want:    "alt@example.com",
		},
		{
			name:    "email-shaped login is accepted",
			info:    &v2.AccountInfo{Login: " user@example.com "},
			profile: map[string]interface{}{},
			want:    "user@example.com",
		},
		{
			name:    "display-name form is rejected",
			info:    &v2.AccountInfo{Login: "Example User <user@example.com>"},
			profile: map[string]interface{}{},
			wantErr: true,
		},
		{
			name:    "no email anywhere",
			info:    &v2.AccountInfo{Login: "test"},
			profile: map[string]interface{}{},
			wantErr: true,
		},
		{
			name:    "invalid emails are skipped in favour of profile",
			info:    &v2.AccountInfo{Emails: []*v2.AccountInfo_Email{email("not-an-email", true)}},
			profile: map[string]interface{}{"email": "profile@example.com"},
			want:    "profile@example.com",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveInviteEmail(tt.info, tt.profile)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got email %q", got)
				}
				if code := status.Code(err); code != codes.InvalidArgument {
					t.Fatalf("expected InvalidArgument, got %s", code)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
