package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	AccessToken = field.StringField(
		"token",
		field.WithDescription("The CloudAMQP access token used to connect to the CloudAMQP API."),
		field.WithRequired(true),
		field.WithIsSecret(true),
		field.WithDisplayName("Access token"),
	)
	BaseURLField = field.StringField(
		"base-url",
		field.WithDisplayName("Base URL"),
		field.WithDescription("Override the CloudAMQP API URL (for testing)"),
		field.WithHidden(true),
		field.WithExportTarget(field.ExportTargetCLIOnly),
	)

	// ConfigurationFields defines the external configuration required for the
	// connector to run.
	ConfigurationFields = []field.SchemaField{
		AccessToken,
		BaseURLField,
	}

	// FieldRelationships defines relationships between the fields listed in
	// ConfigurationFields that can be automatically validated.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(ConfigurationFields)

