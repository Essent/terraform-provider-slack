// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/essent/terraform-provider-slack/internal/tb"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/slack-go/slack"
	"go.uber.org/mock/gomock"
)

func Test_DataSource_AllUsergroups(t *testing.T) {
	testConfig(t, resource.TestStep{
		// arrange
		PreConfig: func() {
			ugbA := tb.NewUsergroupBuilder().WithChannels([]string{"<CHANNEL_A_A>", "<CHANNEL_A_B>"}).WithUsers([]string{"<USER_A_A>", "<USER_A_B>"})
			ugbA.WithID("<ID_A>").WithName("<NAME_A>").WithDescription("<DESC_A>").WithHandle("<HANDLE_A>")
			ugA := ugbA.Build()

			ugbB := tb.NewUsergroupBuilder().WithChannels([]string{"<CHANNEL_B_A>", "<CHANNEL_B_B>"}).WithUsers([]string{"<USER_B_A>", "<USER_B_B>"})
			ugbB.WithID("<ID_B>").WithName("<NAME_B>").WithDescription("<DESC_B>").WithHandle("<HANDLE_B>")
			ugB := ugbB.Build()

			m := tb.MockSlackClient()
			m.EXPECT().GetUserGroups(gomock.Any(), gomock.Any()).Return([]slack.UserGroup{*ugA, *ugB}, nil).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_all_usergroups" "all_usersgroups" {}
		`,
		// assert
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "total_usergroups"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.0.id"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.0.name"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.0.description"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.0.handle"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.0.channels.0"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.0.channels.1"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.0.users.0"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.0.users.1"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.1.id"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.1.name"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.1.description"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.1.handle"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.1.channels.0"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.1.channels.1"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.1.users.0"),
			resource.TestCheckResourceAttrSet("data.slack_all_usergroups.all_usersgroups", "usergroups.1.users.1"),

			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "total_usergroups", tb.ExpectString("2")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.0.id", tb.ExpectString("<ID_A>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.0.name", tb.ExpectString("<NAME_A>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.0.description", tb.ExpectString("<DESC_A>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.0.handle", tb.ExpectString("<HANDLE_A>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.0.channels.0", tb.ExpectString("<CHANNEL_A_A>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.0.channels.1", tb.ExpectString("<CHANNEL_A_B>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.0.users.0", tb.ExpectString("<USER_A_A>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.0.users.1", tb.ExpectString("<USER_A_B>")),

			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.1.id", tb.ExpectString("<ID_B>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.1.name", tb.ExpectString("<NAME_B>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.1.description", tb.ExpectString("<DESC_B>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.1.handle", tb.ExpectString("<HANDLE_B>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.1.channels.0", tb.ExpectString("<CHANNEL_B_A>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.1.channels.1", tb.ExpectString("<CHANNEL_B_B>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.1.users.0", tb.ExpectString("<USER_B_A>")),
			resource.TestCheckResourceAttrWith("data.slack_all_usergroups.all_usersgroups", "usergroups.1.users.1", tb.ExpectString("<USER_B_B>")),
		),
	})
}

func Test_DataSource_AllUsergroups_Error_When_RetrievalFailed(t *testing.T) {
	testConfig(t, resource.TestStep{
		// arrange
		PreConfig: func() {
			m := tb.MockSlackClient()
			m.EXPECT().GetUserGroups(gomock.Any(), gomock.Any()).Return(nil, errors.New("<SLACK_ERROR>")).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_all_usergroups" "all_usergroups" {}
		`,
		// assert
		ExpectError: regexp.MustCompile("<SLACK_ERROR>"),
	})
}

func Test_DataSource_AllUsergroups_Error_WhenSlackClientNil(t *testing.T) {
	// arrange
	res := &datasource.ConfigureResponse{}
	req := datasource.ConfigureRequest{
		ProviderData: &SlackProviderData{
			Client: nil,
		},
	}

	test_instance := AllUserGroupsDataSource{}

	// act
	test_instance.Configure(context.Background(), req, res)

	// assert
	if res.Diagnostics.Errors()[0].Summary() != "Invalid Provider Data" {
		t.Errorf("Expected error summary to be 'Invalid Provider Data', got: %s", res.Diagnostics.Errors()[0].Summary())
	}
}
