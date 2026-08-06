package connector

import (
	"context"
	"testing"

	"github.com/conductorone/baton-cloudamqp/pkg/cloudamqp"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// TestUserResourceTraitAndAttributes pins what a synced user carries after
// profile and status moved off UserTrait onto the resource itself: the user
// trait must still be present with the email, and the profile and status must
// now be resource-level attributes.
func TestUserResourceTraitAndAttributes(t *testing.T) {
	res, err := userResource(context.Background(), &cloudamqp.User{
		BaseResource: cloudamqp.BaseResource{Id: "42"},
		Email:        "user@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	trait, err := resource.GetUserTrait(res)
	if err != nil {
		t.Fatalf("user trait missing: %v", err)
	}
	if got := trait.GetEmails(); len(got) != 1 || got[0].GetAddress() != "user@example.com" || !got[0].GetIsPrimary() {
		t.Fatalf("expected one primary email user@example.com, got %v", got)
	}

	if got := res.GetStatus().GetStatus(); got != v2.Status_RESOURCE_STATUS_ENABLED {
		t.Fatalf("resource status = %v, want RESOURCE_STATUS_ENABLED", got)
	}

	profile := res.GetProfile().AsMap()
	if profile["login"] != "user@example.com" || profile["user_id"] != "42" {
		t.Fatalf("resource profile = %v, want login and user_id set", profile)
	}
}

func TestResolveInviteEmail(t *testing.T) {
	email := func(addr string, primary bool) *v2.AccountInfo_Email {
		return &v2.AccountInfo_Email{Address: addr, IsPrimary: primary}
	}
	withProfileEmail := func(info *v2.AccountInfo, email string) *v2.AccountInfo {
		profile, err := structpb.NewStruct(map[string]interface{}{"email": email})
		if err != nil {
			t.Fatalf("failed to build profile: %v", err)
		}
		info.Profile = profile
		return info
	}

	for _, tt := range []struct {
		name    string
		info    *v2.AccountInfo
		want    string
		wantErr bool
	}{
		{
			name: "username login falls back to profile email",
			info: withProfileEmail(&v2.AccountInfo{Login: "test"}, "user@example.com"),
			want: "user@example.com",
		},
		{
			name: "profile email wins over emails and login",
			info: withProfileEmail(&v2.AccountInfo{
				Login:  "someone@else.com",
				Emails: []*v2.AccountInfo_Email{email("primary@example.com", true)},
			}, "profile@example.com"),
			want: "profile@example.com",
		},
		{
			name: "primary email used when profile has none",
			info: &v2.AccountInfo{
				Login:  "test",
				Emails: []*v2.AccountInfo_Email{email("alt@example.com", false), email("primary@example.com", true)},
			},
			want: "primary@example.com",
		},
		{
			name: "non-primary email used when no primary is marked",
			info: &v2.AccountInfo{Emails: []*v2.AccountInfo_Email{email("alt@example.com", false)}},
			want: "alt@example.com",
		},
		{
			name: "email-shaped login is accepted",
			info: &v2.AccountInfo{Login: " user@example.com "},
			want: "user@example.com",
		},
		{
			name:    "display-name form is rejected",
			info:    &v2.AccountInfo{Login: "Example User <user@example.com>"},
			wantErr: true,
		},
		{
			name:    "no email anywhere",
			info:    &v2.AccountInfo{Login: "test"},
			wantErr: true,
		},
		{
			name: "invalid primary email falls through to a valid secondary",
			info: &v2.AccountInfo{Emails: []*v2.AccountInfo_Email{email("not-an-email", true), email("alt@example.com", false)}},
			want: "alt@example.com",
		},
		{
			name: "non-email profile value falls through to emails",
			info: withProfileEmail(&v2.AccountInfo{
				Emails: []*v2.AccountInfo_Email{email("primary@example.com", true)},
			}, "not-an-email"),
			want: "primary@example.com",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveInviteEmail(tt.info)
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
