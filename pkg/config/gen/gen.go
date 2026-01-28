package main

import (
	cfg "github.com/conductorone/baton-cloudamqp/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("cloudamqp", cfg.Config)
}
