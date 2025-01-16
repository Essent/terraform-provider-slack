// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"regexp"
	"testing"

	tr "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/essent/terraform-provider-slack/internal/tb"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/slack-go/slack"
	"go.uber.org/mock/gomock"
)

func Test_DataSource_AllUsers(t *testing.T) {
	testConfig(t, tr.TestStep{
		PreConfig: func() {
			uA := tb.NewUserBuilder().WithID("<ID_A>").WithName("<NAME_A>").WithEmail("<EMAIL_A>").Build()
			uB := tb.NewUserBuilder().WithID("<ID_B>").WithName("<NAME_B>").WithEmail("<EMAIL_B>").Build()

			m := tb.MockSlackClient()
			m.EXPECT().GetUsersContext(gomock.Any()).Return([]slack.User{*uA, *uB}, nil).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_all_users" "all_users" {}
		`,
		Check: tr.ComposeTestCheckFunc(
			tr.TestCheckResourceAttrSet("data.slack_all_users.all_users", "total_users"),
			tr.TestCheckResourceAttrSet("data.slack_all_users.all_users", "users.0.id"),
			tr.TestCheckResourceAttrSet("data.slack_all_users.all_users", "users.0.name"),
			tr.TestCheckResourceAttrSet("data.slack_all_users.all_users", "users.0.email"),

			tr.TestCheckResourceAttrWith("data.slack_all_users.all_users", "total_users", tb.ExpectString("2")),
			tr.TestCheckResourceAttrWith("data.slack_all_users.all_users", "users.0.id", tb.ExpectString("<ID_A>")),
			tr.TestCheckResourceAttrWith("data.slack_all_users.all_users", "users.0.name", tb.ExpectString("<NAME_A>")),
			tr.TestCheckResourceAttrWith("data.slack_all_users.all_users", "users.0.email", tb.ExpectString("<EMAIL_A>")),
			tr.TestCheckResourceAttrWith("data.slack_all_users.all_users", "users.1.id", tb.ExpectString("<ID_B>")),
			tr.TestCheckResourceAttrWith("data.slack_all_users.all_users", "users.1.name", tb.ExpectString("<NAME_B>")),
			tr.TestCheckResourceAttrWith("data.slack_all_users.all_users", "users.1.email", tb.ExpectString("<EMAIL_B>")),
		),
	})
}

func Test_DataSource_AllUsers_Error_When_RetrievalFailed(t *testing.T) {
	testConfig(t, tr.TestStep{
		// arrange
		PreConfig: func() {
			m := tb.MockSlackClient()
			m.EXPECT().GetUsersContext(gomock.Any()).Return(nil, errors.New("<SLACK_ERROR>")).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_all_users" "all_users" {}
		`,
		// assert
		ExpectError: regexp.MustCompile(`<SLACK_ERROR>`),
	})
}

func Test_DataSource_AllUsers_Error_WhenSlackClientNil(t *testing.T) {
	// arrange
	res := &datasource.ConfigureResponse{}
	req := datasource.ConfigureRequest{
		ProviderData: &SlackProviderData{
			Client: nil,
		},
	}

	test_instance := AllUsersDataSource{}

	// act
	test_instance.Configure(context.Background(), req, res)

	// assert
	if res.Diagnostics.Errors()[0].Summary() != "Invalid Provider Data" {
		t.Errorf("Expected error summary to be 'Invalid Provider Data', got: %s", res.Diagnostics.Errors()[0].Summary())
	}
}
