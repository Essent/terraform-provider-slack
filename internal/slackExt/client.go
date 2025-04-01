// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package slackExt

import (
	"context"

	"github.com/slack-go/slack"
)

type Client interface {
	AuthTest(ctx context.Context) (*slack.AuthTestResponse, error)
	GetUserInfo(ctx context.Context, user string) (*slack.User, error)
	GetUserByEmail(ctx context.Context, email string) (*slack.User, error)
	GetUsersContext(ctx context.Context) ([]slack.User, error)
	GetUserGroups(ctx context.Context, options ...slack.GetUserGroupsOption) ([]slack.UserGroup, error)
	GetConversationInfo(ctx context.Context, input *slack.GetConversationInfoInput) (*slack.Channel, error)

	CreateUserGroup(ctx context.Context, userGroup slack.UserGroup) (slack.UserGroup, error)
	DisableUserGroup(ctx context.Context, userGroup string) (slack.UserGroup, error)
	EnableUserGroup(ctx context.Context, userGroup string) (slack.UserGroup, error)
	UpdateUserGroup(ctx context.Context, userGroupID string, options ...slack.UpdateUserGroupsOption) (slack.UserGroup, error)
	UpdateUserGroupMembers(ctx context.Context, userGroup string, members string) (slack.UserGroup, error)
}

func New(base *slack.Client) Client {
	return &clientRateLimit{&clientImpl{base}}
}
