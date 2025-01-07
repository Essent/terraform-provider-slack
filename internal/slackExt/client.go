package slackExt

import (
	"context"

	"github.com/slack-go/slack"
)

type Client interface {
	GetUserInfo(user string) (*slack.User, error)
	GetUserByEmail(email string) (*slack.User, error)
	GetUsersContext(ctx context.Context) ([]slack.User, error)
	GetUserGroups(options ...slack.GetUserGroupsOption) ([]slack.UserGroup, error)
}

func New(base *slack.Client) Client {
	return &clientRateLimit{&clientImpl{base}}
}
