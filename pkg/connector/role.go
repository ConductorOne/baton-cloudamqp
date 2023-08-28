package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-cloudamqp/pkg/cloudamqp"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const (
	roleAdmin             = "admin"
	roleDevops            = "devops"
	roleMember            = "member"
	roleMonitor           = "monitor"
	roleBillingManager    = "billing manager"
	roleComplianceManager = "compliance manager"
)

var teamAccessRoles = []string{
	roleAdmin, roleDevops, roleMember, roleMonitor, roleBillingManager, roleComplianceManager,
}

type roleResourceType struct {
	resourceType *v2.ResourceType
	client       *cloudamqp.Client
}

func (r *roleResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return r.resourceType
}

// roleResource creates a new connector resource for a CloudAMQP Role.
func roleResource(role string) (*v2.Resource, error) {
	displayName := titleCase(role)
	profile := map[string]interface{}{
		"role_id":   role,
		"role_name": displayName,
	}

	resource, err := rs.NewRoleResource(
		displayName,
		resourceTypeRole,
		role,
		[]rs.RoleTraitOption{rs.WithRoleProfile(profile)},
	)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (r *roleResourceType) List(ctx context.Context, parentID *v2.ResourceId, pt *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	rv := make([]*v2.Resource, 0, len(teamAccessRoles))

	for _, role := range teamAccessRoles {
		rr, err := roleResource(role)
		if err != nil {
			return nil, "", nil, err
		}

		rv = append(rv, rr)
	}

	return rv, "", nil, nil
}

func (r *roleResourceType) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	entitlementOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(resourceTypeUser),
		ent.WithDisplayName(fmt.Sprintf("%s role", resource.DisplayName)),
		ent.WithDescription(fmt.Sprintf("%s CloudAMQP role", resource.DisplayName)),
	}

	rv = append(rv, ent.NewAssignmentEntitlement(resource, roleMember, entitlementOptions...))

	return rv, "", nil, nil
}

func (r *roleResourceType) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	users, err := r.client.GetUsers(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("cloudamqp-connector: failed to get users: %w", err)
	}

	var rv []*v2.Grant
	for _, user := range users {
		userCopy := user

		ur, err := userResource(ctx, &userCopy)
		if err != nil {
			return nil, "", nil, fmt.Errorf("cloudamqp-connector: failed to build user resource: %w", err)
		}

		for _, role := range user.Roles {
			if role == resource.Id.Resource {
				rv = append(rv, grant.NewGrant(
					resource,
					roleMember,
					ur.Id,
				))
			}
		}
	}

	return rv, "", nil, nil
}

func (r *roleResourceType) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	if principal.Id.ResourceType != resourceTypeUser.Id {
		l.Warn(
			"cloudamqp-connector: only users can be granted roles",
			zap.String("principal_id", principal.Id.String()),
			zap.String("principal_type", principal.Id.ResourceType),
		)

		return nil, fmt.Errorf("cloudamqp-connector: only users can be granted roles")
	}

	userId, roleId := principal.Id.Resource, entitlement.Resource.Id.Resource
	err := r.client.UpdateUserRole(ctx, userId, roleId)
	if err != nil {
		return nil, fmt.Errorf("cloudamqp-connector: failed to update user role: %w", err)
	}

	return nil, nil
}

// Since user always has a role, this function will assign the user to the default role - member.
func (r *roleResourceType) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	principal := grant.Principal

	if principal.Id.ResourceType != resourceTypeUser.Id {
		l.Warn(
			"cloudamqp-connector: only users can have roles revoked",
			zap.String("principal_id", principal.Id.String()),
			zap.String("principal_type", principal.Id.ResourceType),
		)

		return nil, fmt.Errorf("cloudamqp-connector: only users can have roles revoked")
	}

	userId, roleId := principal.Id.Resource, roleMember
	err := r.client.UpdateUserRole(ctx, userId, roleId)
	if err != nil {
		return nil, fmt.Errorf("cloudamqp-connector: failed to update user role: %w", err)
	}

	return nil, nil
}

func roleBuilder(client *cloudamqp.Client) *roleResourceType {
	return &roleResourceType{
		resourceType: resourceTypeRole,
		client:       client,
	}
}
