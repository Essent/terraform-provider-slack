package provider

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/essent/terraform-provider-slack/internal/tb"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"go.uber.org/mock/gomock"
)

func Test_DataSource_User_ByEmail(t *testing.T) {
	testConfig(t, resource.TestStep{
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
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet("data.slack_user.user_by_email", "id"),
			resource.TestCheckResourceAttrSet("data.slack_user.user_by_email", "name"),
			resource.TestCheckResourceAttrSet("data.slack_user.user_by_email", "email"),

			resource.TestCheckResourceAttrWith("data.slack_user.user_by_email", "id", tb.ExpectString("<ID>")),
			resource.TestCheckResourceAttrWith("data.slack_user.user_by_email", "name", tb.ExpectString("<NAME>")),
			resource.TestCheckResourceAttrWith("data.slack_user.user_by_email", "email", tb.ExpectString("<GIVEN_EMAIL>")),
		),
	})
}

func Test_DataSource_User_ByID(t *testing.T) {
	testConfig(t, resource.TestStep{
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
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet("data.slack_user.user_by_id", "id"),
			resource.TestCheckResourceAttrSet("data.slack_user.user_by_id", "name"),
			resource.TestCheckResourceAttrSet("data.slack_user.user_by_id", "email"),

			resource.TestCheckResourceAttrWith("data.slack_user.user_by_id", "id", tb.ExpectString("<GIVEN_ID>")),
			resource.TestCheckResourceAttrWith("data.slack_user.user_by_id", "name", tb.ExpectString("<NAME>")),
			resource.TestCheckResourceAttrWith("data.slack_user.user_by_id", "email", tb.ExpectString("<EMAIL>")),
		),
	})
}

func Test_DataSource_User_Error_When_IdAndEmail_BothSpecified(t *testing.T) {
	testConfig(t, resource.TestStep{
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
	testConfig(t, resource.TestStep{
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
	testConfig(t, resource.TestStep{
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
	testConfig(t, resource.TestStep{
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
	testConfig(t, resource.TestStep{
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
