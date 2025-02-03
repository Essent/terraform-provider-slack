# Terraform Provider Slack

This provider manages Slack usergroups as a resource and provides data sources for other Slack objects (e.g., channels, users).
Compared to the original [provider](https://github.com/pablovarela/terraform-provider-slack) this version contains various changes:
- Rebuilt around the terraform plugin framework
- Added rate limit handling for all Slack API endpoints
- Improved handling of conflicts with existing usergroups
- User datasources no longer return disabled/deleted users
- Removed the need for costly lookups by introducing new datasources that return lists of usergroups/users

Contributions are welcome!
Some ideas for future improvements:
- Implement a conversation resource to manage channels

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
