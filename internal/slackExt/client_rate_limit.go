package slackExt

import (
	"context"
	"time"

	"github.com/slack-go/slack"
)

type clientRateLimit struct {
	base Client
}

func rateLimit[R any](ctx context.Context, f func() (R, error), getZeroValue func() R) (result R, err error) {
	for {
		result, err = f()

		if err == nil {
			return result, nil
		}

		if rateLimitedError, ok := err.(*slack.RateLimitedError); ok {
			select {
			case <-time.After(rateLimitedError.RetryAfter):
			case <-ctx.Done():
				return getZeroValue(), ctx.Err()
			}
		} else {
			return getZeroValue(), err
		}
	}
}

func (c *clientRateLimit) GetUserInfo(ctx context.Context, user string) (result *slack.User, err error) {
	return rateLimit(ctx, func() (*slack.User, error) {
		return c.base.GetUserInfo(ctx, user)
	}, func() *slack.User { return nil })
}

func (c *clientRateLimit) GetUserByEmail(ctx context.Context, email string) (*slack.User, error) {
	return rateLimit(ctx, func() (*slack.User, error) {
		return c.base.GetUserByEmail(ctx, email)
	}, func() *slack.User { return nil })
}

func (c *clientRateLimit) GetUsersContext(ctx context.Context) ([]slack.User, error) {
	return c.base.GetUsersContext(ctx)
}

func (c *clientRateLimit) GetUserGroups(ctx context.Context, options ...slack.GetUserGroupsOption) ([]slack.UserGroup, error) {
	return rateLimit(ctx, func() ([]slack.UserGroup, error) {
		return c.base.GetUserGroups(ctx, options...)
	}, func() []slack.UserGroup { return []slack.UserGroup{} })
}
