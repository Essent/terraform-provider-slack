package slackExt

import (
	"context"

	"github.com/slack-go/slack"
)

type clientImpl struct {
	base *slack.Client
}

func (c *clientImpl) GetUserInfo(user string) (*slack.User, error) {
	return c.base.GetUserInfo(user)
}

func (c *clientImpl) GetUserByEmail(email string) (*slack.User, error) {
	return c.base.GetUserByEmail(email)
}

func (c *clientImpl) GetUsersContext(ctx context.Context) ([]slack.User, error) {
	return c.base.GetUsersContext(ctx)
}

func (c *clientImpl) GetUserGroups(options ...slack.GetUserGroupsOption) ([]slack.UserGroup, error) {
	return c.base.GetUserGroups(options...)
}
