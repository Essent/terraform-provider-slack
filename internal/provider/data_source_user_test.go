package provider

import (
	"testing"

	"github.com/essent/terraform-provider-slack/internal/testBed"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/slack-go/slack"
)

func Test_WHATEVER(t *testing.T) {
	defer testBed.Finish()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

			m := testBed.MockSlackClient()
			m.EXPECT().AuthTest(gomock.Any()).Return(&slack.AuthTestResponse{}, nil).AnyTimes()
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					// arrange
					u := testBed.NewUserBuilder().WithID("<ID>").WithName("<NAME>").WithEmail("<GIVEN_EMAIL>").Build()

					m := testBed.MockSlackClient()
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

					resource.TestCheckResourceAttrWith("data.slack_user.user_by_email", "id", testBed.ExpectString("<ID>")),
					resource.TestCheckResourceAttrWith("data.slack_user.user_by_email", "name", testBed.ExpectString("<NAME>")),
					resource.TestCheckResourceAttrWith("data.slack_user.user_by_email", "email", testBed.ExpectString("<GIVEN_EMAIL>")),
				),
			},
		},
	})
}
