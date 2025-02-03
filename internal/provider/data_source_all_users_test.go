package provider

import (
	"errors"
	"regexp"
	"testing"

	"github.com/essent/terraform-provider-slack/internal/testBed"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/slack-go/slack"
	"go.uber.org/mock/gomock"
)

func Test_DataSource_AllUsers(t *testing.T) {
	testConfig(t, resource.TestStep{
		PreConfig: func() {
			uA := testBed.NewUserBuilder().WithID("<ID_A>").WithName("<NAME_A>").WithEmail("<EMAIL_A>").Build()
			uB := testBed.NewUserBuilder().WithID("<ID_B>").WithName("<NAME_B>").WithEmail("<EMAIL_B>").Build()

			m := testBed.MockSlackClient()
			m.EXPECT().GetUsersContext(gomock.Any()).Return([]slack.User{*uA, *uB}, nil).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_all_users" "all_users" {}
		`,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet("data.slack_all_users.all_users", "total_users"),
			resource.TestCheckResourceAttrSet("data.slack_all_users.all_users", "users.0.id"),
			resource.TestCheckResourceAttrSet("data.slack_all_users.all_users", "users.0.name"),
			resource.TestCheckResourceAttrSet("data.slack_all_users.all_users", "users.0.email"),

			resource.TestCheckResourceAttrWith("data.slack_all_users.all_users", "total_users", testBed.ExpectString("2")),
			resource.TestCheckResourceAttrWith("data.slack_all_users.all_users", "users.0.id", testBed.ExpectString("<ID_A>")),
			resource.TestCheckResourceAttrWith("data.slack_all_users.all_users", "users.0.name", testBed.ExpectString("<NAME_A>")),
			resource.TestCheckResourceAttrWith("data.slack_all_users.all_users", "users.0.email", testBed.ExpectString("<EMAIL_A>")),
			resource.TestCheckResourceAttrWith("data.slack_all_users.all_users", "users.1.id", testBed.ExpectString("<ID_B>")),
			resource.TestCheckResourceAttrWith("data.slack_all_users.all_users", "users.1.name", testBed.ExpectString("<NAME_B>")),
			resource.TestCheckResourceAttrWith("data.slack_all_users.all_users", "users.1.email", testBed.ExpectString("<EMAIL_B>")),
		),
	})
}

func Test_DataSource_Error_When_RetrievalFailed(t *testing.T) {
	testConfig(t, resource.TestStep{
		// arrange
		PreConfig: func() {
			m := testBed.MockSlackClient()
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
