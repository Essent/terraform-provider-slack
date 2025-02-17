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
	"go.uber.org/mock/gomock"
)

func Test_DataSource_User_ByEmail(t *testing.T) {
	testConfig(t, tr.TestStep{
		// arrange
		PreConfig: func() {
			u := tb.NewUserBuilder().WithID("<ID>").WithName("<NAME>").WithEmail("<GIVEN_EMAIL>").Build()

			m := tb.MockSlackClient()
			m.EXPECT().GetUserByEmail(gomock.Any(), "<GIVEN_EMAIL>").Return(u, nil).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_user" "user_by_email" {
				email = "<GIVEN_EMAIL>"
			}
		`,
		// assert
		Check: tr.ComposeTestCheckFunc(
			tr.TestCheckResourceAttrSet("data.slack_user.user_by_email", "id"),
			tr.TestCheckResourceAttrSet("data.slack_user.user_by_email", "name"),
			tr.TestCheckResourceAttrSet("data.slack_user.user_by_email", "email"),

			tr.TestCheckResourceAttrWith("data.slack_user.user_by_email", "id", tb.ExpectString("<ID>")),
			tr.TestCheckResourceAttrWith("data.slack_user.user_by_email", "name", tb.ExpectString("<NAME>")),
			tr.TestCheckResourceAttrWith("data.slack_user.user_by_email", "email", tb.ExpectString("<GIVEN_EMAIL>")),
		),
	})
}

func Test_DataSource_User_ByID(t *testing.T) {
	testConfig(t, tr.TestStep{
		// arrange
		PreConfig: func() {
			u := tb.NewUserBuilder().WithID("<GIVEN_ID>").WithName("<NAME>").WithEmail("<EMAIL>").Build()

			m := tb.MockSlackClient()
			m.EXPECT().GetUserInfo(gomock.Any(), "<GIVEN_ID>").Return(u, nil).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_user" "user_by_id" {
				id = "<GIVEN_ID>"
			}
		`,
		// assert
		Check: tr.ComposeTestCheckFunc(
			tr.TestCheckResourceAttrSet("data.slack_user.user_by_id", "id"),
			tr.TestCheckResourceAttrSet("data.slack_user.user_by_id", "name"),
			tr.TestCheckResourceAttrSet("data.slack_user.user_by_id", "email"),

			tr.TestCheckResourceAttrWith("data.slack_user.user_by_id", "id", tb.ExpectString("<GIVEN_ID>")),
			tr.TestCheckResourceAttrWith("data.slack_user.user_by_id", "name", tb.ExpectString("<NAME>")),
			tr.TestCheckResourceAttrWith("data.slack_user.user_by_id", "email", tb.ExpectString("<EMAIL>")),
		),
	})
}

func Test_DataSource_User_Error_When_IdAndEmail_BothSpecified(t *testing.T) {
	testConfig(t, tr.TestStep{
		// arrange
		PreConfig: func() {
			u := tb.NewUserBuilder().WithID("<GIVEN_ID>").WithName("<NAME>").WithEmail("<GIVEN_EMAIL>").Build()

			m := tb.MockSlackClient()
			m.EXPECT().GetUserInfo(gomock.Any(), "<GIVEN_ID>").Return(u, nil).AnyTimes()
			m.EXPECT().GetUserByEmail(gomock.Any(), "<GIVEN_EMAIL>").Return(u, nil).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_user" "user" {
				id = "<GIVEN_ID>"
				email = "<GIVEN_EMAIL>"
			}
		`,
		// assert
		ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
	})
}

func Test_DataSource_User_Error_When_RetrievalFailed_ByID(t *testing.T) {
	testConfig(t, tr.TestStep{
		// arrange
		PreConfig: func() {
			m := tb.MockSlackClient()
			m.EXPECT().GetUserInfo(gomock.Any(), "<GIVEN_ID>").Return(nil, errors.New("<SLACK_ERR>")).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_user" "user_by_id" {
				id = "<GIVEN_ID>"
			}
		`,
		// assert
		ExpectError: regexp.MustCompile("<SLACK_ERR>"),
	})
}

func Test_DataSource_User_Error_When_RetrievalFailed_ByEmail(t *testing.T) {
	testConfig(t, tr.TestStep{
		// arrange
		PreConfig: func() {
			m := tb.MockSlackClient()
			m.EXPECT().GetUserByEmail(gomock.Any(), "<GIVEN_EMAIL>").Return(nil, errors.New("<SLACK_ERR>")).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_user" "user_by_email" {
				email = "<GIVEN_EMAIL>"
			}
		`,
		// assert
		ExpectError: regexp.MustCompile("<SLACK_ERR>"),
	})
}

func Test_DataSource_User_Error_When_Deleted_ByID(t *testing.T) {
	testConfig(t, tr.TestStep{
		// arrange
		PreConfig: func() {
			u := tb.NewUserBuilder().WithID("<GIVEN_ID>").WithName("<NAME>").WithEmail("<EMAIL>").WithDeleted(true).Build()

			m := tb.MockSlackClient()
			m.EXPECT().GetUserInfo(gomock.Any(), "<GIVEN_ID>").Return(u, nil).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_user" "user_by_id" {
				id = "<GIVEN_ID>"
			}
		`,
		// assert
		ExpectError: regexp.MustCompile("User is deactivated"),
	})
}

func Test_DataSource_User_Error_When_Deleted_ByEmail(t *testing.T) {
	testConfig(t, tr.TestStep{
		// arrange
		PreConfig: func() {
			u := tb.NewUserBuilder().WithID("<ID>").WithName("<NAME>").WithEmail("<GIVEN_EMAIL>").WithDeleted(true).Build()

			m := tb.MockSlackClient()
			m.EXPECT().GetUserByEmail(gomock.Any(), "<GIVEN_EMAIL>").Return(u, nil).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_user" "user_by_email" {
				email = "<GIVEN_EMAIL>"
			}
		`,
		// assert
		ExpectError: regexp.MustCompile("User is deactivated"),
	})
}

func Test_DataSource_User_Error_WhenSlackClientNil(t *testing.T) {
	// arrange
	res := &datasource.ConfigureResponse{}
	req := datasource.ConfigureRequest{
		ProviderData: &SlackProviderData{
			Client: nil,
		},
	}

	test_instance := UserDataSource{}

	// act
	test_instance.Configure(context.Background(), req, res)

	// assert
	if res.Diagnostics.Errors()[0].Summary() != "Invalid Provider Data" {
		t.Errorf("Expected error summary to be 'Invalid Provider Data', got: %s", res.Diagnostics.Errors()[0].Summary())
	}
}
