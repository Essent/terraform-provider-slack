// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package slackExt

import (
	"context"

	"github.com/slack-go/slack"
)

type clientImpl struct {
	base *slack.Client
}

func (c *clientImpl) GetUserInfo(ctx context.Context, user string) (*slack.User, error) {
	return c.base.GetUserInfoContext(ctx, user)
}

func (c *clientImpl) GetUserByEmail(ctx context.Context, email string) (*slack.User, error) {
	return c.base.GetUserByEmailContext(ctx, email)
}

func (c *clientImpl) GetUsersContext(ctx context.Context) ([]slack.User, error) {
	return c.base.GetUsersContext(ctx)
}

func (c *clientImpl) GetUserGroups(ctx context.Context, options ...slack.GetUserGroupsOption) ([]slack.UserGroup, error) {
	return c.base.GetUserGroupsContext(ctx, options...)
}

func (c *clientImpl) GetConversationInfo(ctx context.Context, input *slack.GetConversationInfoInput) (*slack.Channel, error) {
	return c.base.GetConversationInfoContext(ctx, input)
}
