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
	CheckConflicts(ctx context.Context, name, handle string, includeDisabled bool) error
	FindUserGroupByField(ctx context.Context, searchVal, searchField string, includeDisabled bool) (slack.UserGroup, error)
	UpdateUserGroupMembership(ctx context.Context, groupID string, users []string) error
}

type userGroupServiceImpl struct {
	client slackExt.Client
}

func NewUserGroupService(client slackExt.Client) UserGroupService {
	return &userGroupServiceImpl{
		client: client,
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
		Users:            listToStringSlice(m.Users),
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
			existingGroup, errLookup := s.FindUserGroupByField(ctx, lookupValue, lookupField, true)
			if errLookup != nil {
				return "", fmt.Errorf("Slack returned %q, and %q when trying to find group with %s : %s",
					errCreate.Error(), errLookup.Error(), lookupField, lookupValue)
			}

			if existingGroup.DateDelete == 0 {
				return "", fmt.Errorf("Conflict when creating group '%s' (conflicts with group ID: %s). Cannot reuse an enabled group.",
					createReq.Name, existingGroup.ID)
			}

			if errEnable := s.EnableAndUpdateUserGroup(ctx, existingGroup.ID, plan); errEnable != nil {
				return "", fmt.Errorf("Enable/Update Error: %s", errEnable.Error())
			}
			return existingGroup.ID, nil
		} else {
			return "", fmt.Errorf("Error creating user group: %q", errCreate.Error())
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
	return s.FindUserGroupByField(ctx, id, "id", false)
}

func (s *userGroupServiceImpl) UpdateGroup(ctx context.Context, id string, plan *UserGroupPlan) error {
	return s.EnableAndUpdateUserGroup(ctx, id, plan)
}

func (s *userGroupServiceImpl) DeleteGroup(ctx context.Context, id string) error {
	_, err := s.client.DisableUserGroup(ctx, id)
	if err != nil {
		return fmt.Errorf("Could not disable usergroup: %s", err)
	}
	return nil
}

func (s *userGroupServiceImpl) CheckConflicts(ctx context.Context, name, handle string, includeDisabled bool) error {
	existingByName, errNameLookup := s.FindUserGroupByField(ctx, name, "name", includeDisabled)
	if errNameLookup == nil {
		return fmt.Errorf("Conflict: Existing Enabled Group With Same Name\nAn enabled user group named %q already exists (ID: %s).",
			existingByName.Name, existingByName.ID)
	} else if !strings.Contains(errNameLookup.Error(), "no usergroup with name") {
		return fmt.Errorf("Error Checking Name Conflict\n%s", errNameLookup.Error())
	}

	existingByHandle, errHandleLookup := s.FindUserGroupByField(ctx, handle, "handle", includeDisabled)
	if errHandleLookup == nil {
		return fmt.Errorf("Conflict: Existing Enabled Group With Same Handle\nAn enabled user group with handle %q already exists (ID: %s).", existingByHandle.Handle, existingByHandle.ID)
	} else if !strings.Contains(errHandleLookup.Error(), "no usergroup with handle") {
		return fmt.Errorf("Error Checking Handle Conflict\n%s", errHandleLookup.Error())
	}

	return nil
}

func (s *userGroupServiceImpl) FindUserGroupByField(ctx context.Context, searchVal, searchField string, includeDisabled bool) (slack.UserGroup, error) {
	groups, err := s.client.GetUserGroups(ctx,
		slack.GetUserGroupsOptionIncludeDisabled(includeDisabled),
		slack.GetUserGroupsOptionIncludeUsers(true),
	)
	if err != nil {
		return slack.UserGroup{}, err
	}

	for _, g := range groups {
		var matches bool
		switch searchField {
		case "name":
			matches = strings.EqualFold(g.Name, searchVal)
		case "handle":
			matches = strings.EqualFold(g.Handle, searchVal)
		case "id":
			matches = (g.ID == searchVal)
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

	return slack.UserGroup{}, fmt.Errorf("no usergroup with %s %q found", searchField, searchVal)
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
