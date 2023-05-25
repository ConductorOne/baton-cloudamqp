package connector

import (
	"context"

	"github.com/ConductorOne/baton-cloudamqp/pkg/cloudamqp"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	resourceTypeUser = &v2.ResourceType{
		Id:          "user",
		DisplayName: "User",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_USER,
		},
	}
	resourceTypeRole = &v2.ResourceType{
		Id:          "role",
		DisplayName: "Role",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_ROLE,
		},
	}
)

type CloudAMQP struct {
	client *cloudamqp.Client
}

func (pd *CloudAMQP) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		userBuilder(pd.client),
		roleBuilder(pd.client),
	}
}

// Metadata returns metadata about the connector.
func (pd *CloudAMQP) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "CloudAMQP",
	}, nil
}

// Validate hits the CloudAMQP API to validate that the configured credentials are valid and compatible.
func (pd *CloudAMQP) Validate(ctx context.Context) (annotations.Annotations, error) {
	// should be able to list users
	_, err := pd.client.GetUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Provided Access Token is invalid")
	}

	return nil, nil
}

// New returns the CloudAMQP connector.
func New(ctx context.Context, password string) (*CloudAMQP, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	return &CloudAMQP{
		client: cloudamqp.NewClient(httpClient, password),
	}, nil
}
