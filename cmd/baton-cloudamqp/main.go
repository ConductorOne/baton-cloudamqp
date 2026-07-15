package main

import (
	"context"

	cfg "github.com/conductorone/baton-cloudamqp/pkg/config"
	"github.com/conductorone/baton-cloudamqp/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.Background()

	config.RunConnector(
		ctx,
		"baton-cloudamqp",
		version,
		cfg.Config,
		getConnector,
		connectorrunner.WithDefaultCapabilitiesConnectorBuilderV2(&connector.CloudAMQP{}),
	)
}

func getConnector(ctx context.Context, c *cfg.Cloudamqp, _ *cli.ConnectorOpts) (connectorbuilder.ConnectorBuilderV2, []connectorbuilder.Opt, error) {
	l := ctxzap.Extract(ctx)

	cloudamqpConnector, err := connector.New(ctx, c.Token, c.BaseUrl)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, nil, err
	}

	return cloudamqpConnector, nil, nil
}
