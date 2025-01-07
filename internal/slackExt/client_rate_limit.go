package slackExt

import (
	"context"
	"time"

	"github.com/slack-go/slack"
)

type clientRateLimit struct {
	base Client
}

func rateLimit[R any](f func() (R, error)) (R, error) {
	for {
		result, err := f()

		if err == nil {
			return result, nil
		}

		if rateLimitedError, ok := err.(*slack.RateLimitedError); ok {
			<-time.After(rateLimitedError.RetryAfter)
			err = nil
		}
	}
}

func (c *clientRateLimit) GetUserInfo(user string) (result *slack.User, err error) {
	return rateLimit(func() (*slack.User, error) {
		return c.base.GetUserInfo(user)
	})
}

func (c *clientRateLimit) GetUserByEmail(email string) (*slack.User, error) {
	return rateLimit(func() (*slack.User, error) {
		return c.base.GetUserByEmail(email)
	})
}

func (c *clientRateLimit) GetUsersContext(ctx context.Context) ([]slack.User, error) {
	return c.base.GetUsersContext(ctx)
}

func (c *clientRateLimit) GetUserGroups(options ...slack.GetUserGroupsOption) ([]slack.UserGroup, error) {
	return rateLimit(func() ([]slack.UserGroup, error) {
		return c.base.GetUserGroups(options...)
	})
}
