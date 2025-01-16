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

func (c *clientImpl) AuthTest(ctx context.Context) (*slack.AuthTestResponse, error) {
	return c.base.AuthTestContext(ctx)
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

func (c *clientImpl) CreateUserGroup(ctx context.Context, userGroup slack.UserGroup) (slack.UserGroup, error) {
	return c.base.CreateUserGroupContext(ctx, userGroup)
}

func (c *clientImpl) DisableUserGroup(ctx context.Context, userGroup string) (slack.UserGroup, error) {
	return c.base.DisableUserGroupContext(ctx, userGroup)
}

func (c *clientImpl) EnableUserGroup(ctx context.Context, userGroup string) (slack.UserGroup, error) {
	return c.base.EnableUserGroupContext(ctx, userGroup)
}

func (c *clientImpl) UpdateUserGroup(ctx context.Context, userGroupID string, options ...slack.UpdateUserGroupsOption) (slack.UserGroup, error) {
	return c.base.UpdateUserGroupContext(ctx, userGroupID, options...)
}

func (c *clientImpl) UpdateUserGroupMembers(ctx context.Context, userGroup string, members string) (slack.UserGroup, error) {
	return c.base.UpdateUserGroupMembersContext(ctx, userGroup, members)
}
