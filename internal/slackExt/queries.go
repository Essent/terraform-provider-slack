package slackExt

import (
	"context"

	"github.com/slack-go/slack"
)

type Queries interface {
	FindUserGroupByField(ctx context.Context, field, value string, includeDisabled bool) (slack.UserGroup, error)
}

func NewQueries(client Client) Queries {
	return &queriesImpl{client}
}
