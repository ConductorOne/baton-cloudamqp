![Baton Logo](./docs/images/baton-logo.png)

# `baton-cloudamqp` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-cloudamqp.svg)](https://pkg.go.dev/github.com/conductorone/baton-cloudamqp) ![main ci](https://github.com/conductorone/baton-cloudamqp/actions/workflows/main.yaml/badge.svg)

`baton-cloudamqp` is a connector for CloudAMQP built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the CloudAMQP Customer API to sync data about users and their roles.

Check out [Baton](https://github.com/conductorone/baton) to learn more about the project in general.

# Prerequisites

To work with the connector, you need to obtain API access token from CloudAMQP. 

Create token by going to the top menu bar, selecting the dropdown under your team name  -> `Team Settings` and from there click on `API Access` item in left menu bar where you can create either full access keys (shared with other team admins) or Personal access keys.

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-cloudamqp

BATON_TOKEN=token baton-cloudamqp
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_TOKEN=token ghcr.io/conductorone/baton-cloudamqp:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-cloudamqp/cmd/baton-cloudamqp@main

BATON_TOKEN=token baton-cloudamqp
baton resources
```

# Data Model

`baton-cloudamqp` will pull down information about the following CloudAMQP resources:

- Users

By default, `baton-cloudamqp` will sync information only from account based on provided credential.

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-cloudamqp` Command Line Usage

```
baton-cloudamqp

Usage:
  baton-cloudamqp [flags]
  baton-cloudamqp [command]

Available Commands:
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
  -f, --file string                   The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                          help for baton-cloudamqp
      --log-format string             The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string              The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --token string                  The CloudAMQP access token used to connect to the CloudAMQP API. ($BATON_TOKEN)
  -v, --version                       version for baton-cloudamqp

Use "baton-cloudamqp [command] --help" for more information about a command.
```
