package connector

import (
	"context"
	"fmt"
	"strings"

	"github.com/ConductorOne/baton-cloudamqp/pkg/cloudamqp"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

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
		"login": user.Email,
		"roles": strings.Join(user.Roles, ","),
	}

	ret, err := resource.NewUserResource(
		user.Email,
		resourceTypeUser,
		user.Id,
		[]resource.UserTraitOption{
			resource.WithEmail(user.Email, true),
			resource.WithUserProfile(profile),
		},
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (u *userResourceType) List(ctx context.Context, parentID *v2.ResourceId, pt *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	users, err := u.client.GetUsers(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("cloudamqp-connector: failed to list users: %w", err)
	}

	rv := make([]*v2.Resource, 0, len(users))
	for _, user := range users {
		userCopy := user

		ur, err := userResource(ctx, &userCopy)
		if err != nil {
			return nil, "", nil, err
		}

		rv = append(rv, ur)
	}

	return rv, "", nil, nil
}

func (u *userResourceType) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (u *userResourceType) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func userBuilder(client *cloudamqp.Client) *userResourceType {
	return &userResourceType{
		resourceType: resourceTypeUser,
		client:       client,
	}
}
