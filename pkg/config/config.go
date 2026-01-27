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

	// FieldRelationships defines relationships between the fields listed in
	// Config that can be automatically validated.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run ./gen
var Config = field.NewConfiguration([]field.SchemaField{
	AccessToken,
})

