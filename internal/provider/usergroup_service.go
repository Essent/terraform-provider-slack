// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/essent/terraform-provider-slack/internal/slackExt"
	"github.com/slack-go/slack"
)

type UserGroupService interface {
	CreateGroup(ctx context.Context, plan *UserGroupPlan) (string, error)
	EnableAndUpdateUserGroup(ctx context.Context, groupID string, plan *UserGroupPlan) error
	ReadGroup(ctx context.Context, id string) (slack.UserGroup, error)
	UpdateGroup(ctx context.Context, id string, plan *UserGroupPlan) error
	DeleteGroup(ctx context.Context, id string) error
	CheckConflicts(ctx context.Context, id string, name string, handle string, includeDisabled bool) error
	UpdateUserGroupMembership(ctx context.Context, groupID string, users []string) error
}

type userGroupServiceImpl struct {
	client  slackExt.Client
	queries slackExt.Queries
}

func NewUserGroupService(client slackExt.Client) UserGroupService {
	return &userGroupServiceImpl{
		client:  client,
		queries: slackExt.NewQueries(client),
	}
}

type UserGroupPlan struct {
	ID               string
	Name             string
	Description      string
	Handle           string
	Channels         []string
	Users            []string
	PreventConflicts bool
}

func toPlan(m *UserGroupResourceModel) *UserGroupPlan {
	return &UserGroupPlan{
		ID:               m.ID.ValueString(),
		Name:             m.Name.ValueString(),
		Description:      m.Description.ValueString(),
		Handle:           m.Handle.ValueString(),
		Channels:         listToStringSlice(m.Channels),
		Users:            setToStringSlice(m.Users),
		PreventConflicts: m.PreventConflicts.ValueBool(),
	}
}

func (s *userGroupServiceImpl) CreateGroup(ctx context.Context, plan *UserGroupPlan) (string, error) {
	createReq := slack.UserGroup{
		Name:        plan.Name,
		Description: plan.Description,
		Handle:      plan.Handle,
		Prefs: slack.UserGroupPrefs{
			Channels: plan.Channels,
		},
	}

	created, errCreate := s.client.CreateUserGroup(ctx, createReq)
	if errCreate != nil {
		var lookupField, lookupValue string
		switch errCreate.Error() {
		case "name_already_exists":
			lookupField = "name"
			lookupValue = createReq.Name
		case "handle_already_exists":
			lookupField = "handle"
			lookupValue = createReq.Handle
		}

		if lookupField != "" {
			existingGroup, errLookup := s.queries.FindUserGroupByField(ctx, lookupField, lookupValue, true)
			if errLookup != nil {
				return "", fmt.Errorf("slack returned %q, and %q when trying to find group with %s : %s",
					errCreate.Error(), errLookup.Error(), lookupField, lookupValue)
			}

			if existingGroup.DateDelete == 0 {
				return "", fmt.Errorf("conflict when creating group '%s' (conflicts with group ID: %s): cannot reuse an enabled group",
					createReq.Name, existingGroup.ID)
			}

			if errEnable := s.EnableAndUpdateUserGroup(ctx, existingGroup.ID, plan); errEnable != nil {
				return "", fmt.Errorf("Enable/Update Error: %s", errEnable.Error())
			}
			return existingGroup.ID, nil
		} else {
			return "", fmt.Errorf("error creating user group: %q", errCreate.Error())
		}
	}

	if err := s.UpdateUserGroupMembership(ctx, created.ID, plan.Users); err != nil {
		return "", err
	}

	return created.ID, nil
}

func (s *userGroupServiceImpl) EnableAndUpdateUserGroup(ctx context.Context, groupID string, plan *UserGroupPlan) error {
	_, err := s.client.EnableUserGroup(ctx, groupID)
	if err != nil && err.Error() != "already_enabled" {
		return fmt.Errorf("could not enable usergroup %s: %w", groupID, err)
	}

	opts := []slack.UpdateUserGroupsOption{
		slack.UpdateUserGroupsOptionName(plan.Name),
		slack.UpdateUserGroupsOptionHandle(plan.Handle),
		slack.UpdateUserGroupsOptionDescription(&[]string{plan.Description}[0]),
		slack.UpdateUserGroupsOptionChannels(plan.Channels),
	}

	if _, err := s.client.UpdateUserGroup(ctx, groupID, opts...); err != nil {
		return fmt.Errorf("could not update usergroup %s: %w", groupID, err)
	}

	return s.UpdateUserGroupMembership(ctx, groupID, plan.Users)
}

func (s *userGroupServiceImpl) ReadGroup(ctx context.Context, id string) (slack.UserGroup, error) {
	return s.queries.FindUserGroupByField(ctx, "id", id, false)
}

func (s *userGroupServiceImpl) UpdateGroup(ctx context.Context, id string, plan *UserGroupPlan) error {
	return s.EnableAndUpdateUserGroup(ctx, id, plan)
}

func (s *userGroupServiceImpl) DeleteGroup(ctx context.Context, id string) error {
	_, err := s.client.DisableUserGroup(ctx, id)
	if err != nil {
		return fmt.Errorf("could not disable usergroup: %s", err)
	}
	return nil
}

func (s *userGroupServiceImpl) CheckConflicts(ctx context.Context, resourceID, name, handle string, includeDisabled bool) error {
	conflict := func(field, val string) error {
		existing, errLookup := s.queries.FindUserGroupByField(ctx, field, val, includeDisabled)
		if errLookup != nil {
			if !strings.Contains(errLookup.Error(), "no usergroup with "+field) {
				return fmt.Errorf("error checking %s conflict\n%s", field, errLookup.Error())
			}
			return nil
		}
		if existing.ID != resourceID {
			return fmt.Errorf("conflict: existing usergroup with same %s\nA usergroup with %s %q already exists (ID: %s)", field, field, val, existing.ID)
		}
		return nil
	}

	var errs []string
	if err := conflict("name", name); err != nil {
		errs = append(errs, err.Error())
	}
	if err := conflict("handle", handle); err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errs, "\n\n"))
}

func (s *userGroupServiceImpl) UpdateUserGroupMembership(ctx context.Context, groupID string, users []string) error {
	joined := "[]"
	if len(users) > 0 {
		joined = strings.Join(users, ",")
	}

	_, err := s.client.UpdateUserGroupMembers(ctx, groupID, joined)
	if err != nil {
		return fmt.Errorf("could not update usergroup members: %s", err)
	}
	return nil
}
