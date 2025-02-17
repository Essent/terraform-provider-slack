// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package slackExt

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
)

type queriesImpl struct {
	client Client
}

func (q *queriesImpl) FindUserGroupByField(ctx context.Context, field, value string, includeDisabled bool) (slack.UserGroup, error) {
	groups, err := q.client.GetUserGroups(ctx,
		slack.GetUserGroupsOptionIncludeDisabled(includeDisabled),
		slack.GetUserGroupsOptionIncludeUsers(true),
	)
	if err != nil {
		return slack.UserGroup{}, err
	}

	for _, g := range groups {
		var matches bool
		switch field {
		case "name":
			matches = (g.Name == value)
		case "handle":
			matches = (g.Handle == value)
		case "id":
			matches = (g.ID == value)
		default:
			continue
		}

		if matches {
			if !includeDisabled && g.DateDelete == 0 {
				return g, nil
			} else if includeDisabled {
				return g, nil
			}
		}
	}

	return slack.UserGroup{}, fmt.Errorf("no usergroup with %s %q found", field, value)
}
