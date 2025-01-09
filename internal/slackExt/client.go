// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package slackExt

import (
	"context"

	"github.com/slack-go/slack"
)

type Client interface {
	GetUserInfo(ctx context.Context, user string) (*slack.User, error)
	GetUserByEmail(ctx context.Context, email string) (*slack.User, error)
	GetUsersContext(ctx context.Context) ([]slack.User, error)
	GetUserGroups(ctx context.Context, options ...slack.GetUserGroupsOption) ([]slack.UserGroup, error)
	GetConversationInfo(ctx context.Context, input *slack.GetConversationInfoInput) (*slack.Channel, error)
}

func New(base *slack.Client) Client {
	return &clientRateLimit{&clientImpl{base}}
}
