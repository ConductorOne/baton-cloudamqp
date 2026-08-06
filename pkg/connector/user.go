package connector

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	"github.com/conductorone/baton-cloudamqp/pkg/cloudamqp"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ connectorbuilder.ResourceSyncerV2  = (*userResourceType)(nil)
	_ connectorbuilder.AccountManagerV2  = (*userResourceType)(nil)
	_ connectorbuilder.ResourceDeleterV2 = (*userResourceType)(nil)
)

// defaultInviteRole is used when CreateAccount is called without an explicit
// role. member is the least-privileged team role.
const defaultInviteRole = "member"

// knownInviteRoles are the team-member roles CloudAMQP accepts on POST /team/invite
// (per the CloudAMQP CLI). This is intentionally separate from the role-resource
// grant set in role.go, which models instance-access roles.
var knownInviteRoles = []string{"owner", "admin", "developer", "devops", "member"}

func isKnownInviteRole(role string) bool {
	for _, r := range knownInviteRoles {
		if r == role {
			return true
		}
	}
	return false
}

type userResourceType struct {
	resourceType *v2.ResourceType
	client       *cloudamqp.Client
}

func (u *userResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return u.resourceType
}

// Create a new connector resource for a CloudAMQP User.
func userResource(ctx context.Context, user *cloudamqp.User) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"login":   user.Email,
		"user_id": user.Id,
	}

	ret, err := resource.NewUserResource(
		user.Email,
		resourceTypeUser,
		user.Id,
		[]resource.UserTraitOption{
			resource.WithEmail(user.Email, true),
		},
		resource.WithResourceProfile(profile),
		resource.WithResourceStatus(v2.Status_RESOURCE_STATUS_ENABLED, ""),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (u *userResourceType) List(ctx context.Context, parentID *v2.ResourceId, opts resource.SyncOpAttrs) ([]*v2.Resource, *resource.SyncOpResults, error) {
	users, err := u.client.GetUsers(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("baton-cloudamqp: failed to list users: %w", err)
	}

	rv := make([]*v2.Resource, 0, len(users))
	for _, user := range users {
		userCopy := user

		ur, err := userResource(ctx, &userCopy)
		if err != nil {
			return nil, nil, err
		}

		rv = append(rv, ur)
	}

	return rv, nil, nil
}

func (u *userResourceType) Entitlements(_ context.Context, _ *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Entitlement, *resource.SyncOpResults, error) {
	return nil, nil, nil
}

func (u *userResourceType) Grants(_ context.Context, _ *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Grant, *resource.SyncOpResults, error) {
	return nil, nil, nil
}

// CreateAccountCapabilityDetails declares the credential options for account
// creation. CloudAMQP invitations are email-based and carry no password, so the
// connector advertises the no-password option. This is required for the SDK to
// detect the account-manager capability and wire up CreateAccount/Delete.
func (u *userResourceType) CreateAccountCapabilityDetails(_ context.Context) (*v2.CredentialDetailsAccountProvisioning, annotations.Annotations, error) {
	return &v2.CredentialDetailsAccountProvisioning{
		SupportedCredentialOptions: []v2.CapabilityDetailCredentialOption{
			v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
		},
		PreferredCredentialOption: v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
	}, nil, nil
}

// isEmailAddress reports whether s is a bare email address. CloudAMQP keys
// invitations on the email, so a value that is not one (a username, say) is
// rejected downstream with an opaque 400.
func isEmailAddress(s string) bool {
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return false
	}
	// Reject display-name forms ("Example User <user@example.com>"): CloudAMQP
	// wants the bare address.
	return addr.Address == s
}

// resolveInviteEmail picks the email address to invite from the account info C1
// supplies. Resolution order:
//
//  1. The schema-declared profile "email" field — the value the admin mapped or
//     typed. It WINS: resolving anything ahead of it would silently override an
//     explicit choice. If it's present but not email-shaped, that's a mapping
//     mistake, and we fail loud rather than silently inviting a different
//     address than the one the admin mapped.
//  2. AccountInfo.Emails — the primary entry, else any other valid address.
//  3. Login, but only when it is itself an email address.
//
// Login must never be preferred blindly. It is documented as "the user's login",
// which for many C1 directories is a username, and sending a username to
// POST /team/invite makes CloudAMQP reject the request with an opaque 400.
func resolveInviteEmail(accountInfo *v2.AccountInfo) (string, error) {
	// Key must match the "email" field of the account-creation schema.
	profileEmail, _ := resource.GetProfileStringValue(accountInfo.GetProfile(), "email")
	profileEmail = strings.TrimSpace(profileEmail)
	if profileEmail != "" {
		if !isEmailAddress(profileEmail) {
			return "", status.Errorf(codes.InvalidArgument,
				"baton-cloudamqp: create account: mapped profile \"email\" field (%q) is not a valid email address",
				profileEmail)
		}
		return profileEmail, nil
	}

	var secondary string
	for _, e := range accountInfo.GetEmails() {
		addr := strings.TrimSpace(e.GetAddress())
		if !isEmailAddress(addr) {
			continue
		}
		if e.GetIsPrimary() {
			return addr, nil
		}
		if secondary == "" {
			secondary = addr
		}
	}
	if secondary != "" {
		return secondary, nil
	}

	if login := strings.TrimSpace(accountInfo.GetLogin()); isEmailAddress(login) {
		return login, nil
	}

	return "", status.Errorf(codes.InvalidArgument,
		"baton-cloudamqp: create account: a valid email address is required; CloudAMQP invitations are keyed on email, "+
			"and none of the mapped profile \"email\" field, the account's emails, or the login (%q) is one",
		accountInfo.GetLogin())
}

// CreateAccount invites a new team member to CloudAMQP via POST /team/invite.
//
// Flow:
//  1. If a member with the email already exists, return AlreadyExistsResult with
//     the real id-keyed resource (idempotent; the account-provisioning CI re-runs
//     create).
//  2. Otherwise send the invitation. If CloudAMQP reports the email already
//     exists, re-resolve it: if now a member, AlreadyExistsResult; otherwise the
//     invite is still pending.
//  3. A freshly sent invitation has no stable user id until the invitee accepts
//     the email, so return ActionRequiredResult with no Resource — never fabricate
//     a resource keyed on the email, which could not reconcile with the real
//     id-keyed resource a later sync emits.
func (u *userResourceType) CreateAccount(
	ctx context.Context, accountInfo *v2.AccountInfo, _ *v2.LocalCredentialOptions,
) (connectorbuilder.CreateAccountResponse, []*v2.PlaintextData, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	email, err := resolveInviteEmail(accountInfo)
	if err != nil {
		return nil, nil, nil, err
	}

	role, _ := resource.GetProfileStringValue(accountInfo.GetProfile(), "role")
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		role = defaultInviteRole
	}
	if !isKnownInviteRole(role) {
		return nil, nil, nil, status.Errorf(codes.InvalidArgument, "baton-cloudamqp: create account: unknown role %q (valid: %v)", role, knownInviteRoles)
	}

	tagsRaw, _ := accountInfo.GetProfile().AsMap()["tags"].([]interface{})
	var tags []string
	for _, t := range tagsRaw {
		if s, ok := t.(string); ok {
			s = strings.TrimSpace(s)
			if s != "" {
				tags = append(tags, s)
			}
		}
	}

	// Step 1: short-circuit if the member already exists in this team.
	if existing, err := u.client.GetUserByEmail(ctx, email); err == nil {
		res, resErr := userResource(ctx, existing)
		if resErr != nil {
			return nil, nil, nil, resErr
		}
		l.Debug("baton-cloudamqp: member already exists, returning AlreadyExistsResult", zap.String("email", email))
		return &v2.CreateAccountResponse_AlreadyExistsResult{Resource: res}, nil, nil, nil
	} else if !cloudamqp.IsNotFoundError(err) {
		return nil, nil, nil, fmt.Errorf("baton-cloudamqp: create account: failed to look up member %s: %w", email, err)
	}

	// Step 2: send the invitation.
	if err := u.client.InviteUser(ctx, email, role, tags); err != nil {
		if !cloudamqp.IsAlreadyExistsError(err) {
			return nil, nil, nil, fmt.Errorf("baton-cloudamqp: create account: failed to invite %s: %w", email, err)
		}
		// Invite reported already-exists: re-resolve to decide member vs pending.
		if existing, ferr := u.client.GetUserByEmail(ctx, email); ferr == nil {
			res, resErr := userResource(ctx, existing)
			if resErr != nil {
				return nil, nil, nil, resErr
			}
			return &v2.CreateAccountResponse_AlreadyExistsResult{Resource: res}, nil, nil, nil
		} else if !cloudamqp.IsNotFoundError(ferr) {
			return nil, nil, nil, fmt.Errorf("baton-cloudamqp: create account: failed to re-resolve member %s after 409: %w", email, ferr)
		}
		return &v2.CreateAccountResponse_ActionRequiredResult{
			Message:               fmt.Sprintf("Invitation already pending for %s. The user must accept the email invitation to complete account creation.", email),
			IsCreateAccountResult: true,
		}, nil, nil, nil
	}

	// Step 3: invitation sent — no stable user id until the invitee accepts.
	return &v2.CreateAccountResponse_ActionRequiredResult{
		Message:               fmt.Sprintf("Invitation sent to %s. The user must accept the email invitation to complete account creation.", email),
		IsCreateAccountResult: true,
	}, nil, nil, nil
}

// Delete removes a team member from CloudAMQP. The platform delivers the user's
// id, but removal is keyed on email (POST /team/remove), so the email is resolved
// from the team list first. CloudAMQP has no soft-disable, so this is a hard
// removal. A not-found result at either step is treated as success because the
// platform retries deletes and the member may already be gone.
func (u *userResourceType) Delete(ctx context.Context, resourceID *v2.ResourceId, _ *v2.ResourceId) (annotations.Annotations, error) {
	user, err := u.client.GetUserByID(ctx, resourceID.Resource)
	if err != nil {
		if cloudamqp.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("baton-cloudamqp: delete: failed to resolve member %s: %w", resourceID.Resource, err)
	}

	if err := u.client.RemoveUser(ctx, user.Email); err != nil {
		if cloudamqp.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("baton-cloudamqp: delete: failed to remove member %s: %w", resourceID.Resource, err)
	}

	return nil, nil
}

func userBuilder(client *cloudamqp.Client) *userResourceType {
	return &userResourceType{
		resourceType: resourceTypeUser,
		client:       client,
	}
}
