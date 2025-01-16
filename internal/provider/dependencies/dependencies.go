// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dependencies

import (
	"github.com/essent/terraform-provider-slack/internal/slackExt"
	"github.com/slack-go/slack"
)

type Dependencies interface {
	CreateSlackClient(token string) slackExt.Client
	CreateSlackQueries(client slackExt.Client) slackExt.Queries
}

type dependenciesImpl struct {
}

func (d *dependenciesImpl) CreateSlackClient(token string) slackExt.Client {
	return slackExt.New(slack.New(token))
}

func (d *dependenciesImpl) CreateSlackQueries(client slackExt.Client) slackExt.Queries {
	return slackExt.NewQueries(client)
}

func New() Dependencies {
	return &dependenciesImpl{}
}
