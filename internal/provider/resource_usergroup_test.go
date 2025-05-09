// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/essent/terraform-provider-slack/internal/tb"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tr "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/slack-go/slack"
	"go.uber.org/mock/gomock"
)

func Test_Resource_UserGroup(t *testing.T) {
	testConfig(t, tr.TestStep{
		// arrange
		PreConfig: func() {
			ub := tb.NewUsergroupBuilder().WithName("<NAME>").WithHandle("<HANDLE>")
			ub.WithID("<ID>").WithDescription("<DESCRIPTION>")
			u := ub.Build()

			m := tb.MockSlackClient()
			m.EXPECT().CreateUserGroup(gomock.Any(), gomock.Any()).Return(*u, nil).AnyTimes()
			m.EXPECT().UpdateUserGroupMembers(gomock.Any(), gomock.Any(), gomock.Any()).Return(*u, nil).AnyTimes()
			m.EXPECT().GetUserGroups(gomock.Any(), gomock.Any()).Return([]slack.UserGroup{*u}, nil).AnyTimes()
			m.EXPECT().DisableUserGroup(gomock.Any(), gomock.Any()).Return(*u, nil).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}
				
			resource "slack_usergroup" "group" {
				name = "<NAME>"
				handle = "<HANDLE>"
				description = "<DESCRIPTION>"
			}
		`,
		// assert
		Check: tr.ComposeTestCheckFunc(
			tr.TestCheckResourceAttrSet("slack_usergroup.group", "id"),
			tr.TestCheckResourceAttrSet("slack_usergroup.group", "name"),
			tr.TestCheckResourceAttrSet("slack_usergroup.group", "handle"),
			tr.TestCheckResourceAttrSet("slack_usergroup.group", "description"),

			tr.TestCheckResourceAttrWith("slack_usergroup.group", "id", tb.ExpectString("<ID>")),
			tr.TestCheckResourceAttrWith("slack_usergroup.group", "name", tb.ExpectString("<NAME>")),
			tr.TestCheckResourceAttrWith("slack_usergroup.group", "handle", tb.ExpectString("<HANDLE>")),
			tr.TestCheckResourceAttrWith("slack_usergroup.group", "description", tb.ExpectString("<DESCRIPTION>")),
		),
	})
}

func Test_Resource_UserGroup_Error_WhenSlackClientNil(t *testing.T) {
	// arrange
	res := &resource.ConfigureResponse{}
	req := resource.ConfigureRequest{
		ProviderData: &SlackProviderData{
			Client: nil,
		},
	}

	test_instance := UserGroupResource{}

	// act
	test_instance.Configure(context.Background(), req, res)

	// assert
	if !res.Diagnostics.HasError() {
		t.Error("Expected diagnostics to have errors, but got none")
		return
	}
	if res.Diagnostics.Errors()[0].Summary() != "Invalid Provider Data" {
		t.Errorf("Expected error summary to be 'Invalid Provider Data', got: %s", res.Diagnostics.Errors()[0].Summary())
	}
}

func Test_Resource_UserGroup_Error_WhenSlackQueriesNil(t *testing.T) {
	// arrange
	tb.Init(t)
	defer tb.Finish()

	res := &resource.ConfigureResponse{}
	req := resource.ConfigureRequest{
		ProviderData: &SlackProviderData{
			Client: tb.MockSlackClient(),
		},
	}

	test_instance := UserGroupResource{}

	// act
	test_instance.Configure(context.Background(), req, res)

	// assert
	if !res.Diagnostics.HasError() {
		t.Error("Expected diagnostics to have errors, but got none")
		return
	}
	if res.Diagnostics.Errors()[0].Summary() != "Invalid Provider Data" {
		t.Errorf("Expected error summary to be 'Invalid Provider Data', got: %s", res.Diagnostics.Errors()[0].Summary())
	}
}
