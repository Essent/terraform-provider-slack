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

func Test_DataSource_Conversation(t *testing.T) {
	// arrange
	cb := tb.NewChannelBuilder().WithID("<GIVEN_ID>").WithTopic("<TOPIC>").WithPurpose("<PURPOSE>")
	cb.WithCreated(1234567890).WithCreator("<CREATOR>").WithIsArchived(tb.RandBool()).WithIsShared(tb.RandBool())
	cb.WithIsExtShared(tb.RandBool()).WithIsOrgShared(tb.RandBool()).WithIsGeneral(tb.RandBool())

	c := cb.Build()

	testConfig(t, resource.TestStep{
		PreConfig: func() {
			expected_conversation_info_input := &slack.GetConversationInfoInput{
				ChannelID:         "<GIVEN_ID>",
				IncludeLocale:     false,
				IncludeNumMembers: false,
			}

			m := tb.MockSlackClient()
			m.EXPECT().GetConversationInfo(gomock.Any(), expected_conversation_info_input).Return(c, nil).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_conversation" "channel" {
				channel_id = "<GIVEN_ID>"
			}
		`,
		// assert
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet("data.slack_conversation.channel", "channel_id"),
			resource.TestCheckResourceAttrSet("data.slack_conversation.channel", "topic"),
			resource.TestCheckResourceAttrSet("data.slack_conversation.channel", "purpose"),
			resource.TestCheckResourceAttrSet("data.slack_conversation.channel", "created"),
			resource.TestCheckResourceAttrSet("data.slack_conversation.channel", "creator"),
			resource.TestCheckResourceAttrSet("data.slack_conversation.channel", "is_archived"),
			resource.TestCheckResourceAttrSet("data.slack_conversation.channel", "is_shared"),
			resource.TestCheckResourceAttrSet("data.slack_conversation.channel", "is_ext_shared"),
			resource.TestCheckResourceAttrSet("data.slack_conversation.channel", "is_org_shared"),
			resource.TestCheckResourceAttrSet("data.slack_conversation.channel", "is_general"),

			resource.TestCheckResourceAttrWith("data.slack_conversation.channel", "channel_id", tb.ExpectString("<GIVEN_ID>")),
			resource.TestCheckResourceAttrWith("data.slack_conversation.channel", "topic", tb.ExpectString("<TOPIC>")),
			resource.TestCheckResourceAttrWith("data.slack_conversation.channel", "purpose", tb.ExpectString("<PURPOSE>")),
			resource.TestCheckResourceAttrWith("data.slack_conversation.channel", "created", tb.ExpectString("1234567890")),
			resource.TestCheckResourceAttrWith("data.slack_conversation.channel", "creator", tb.ExpectString("<CREATOR>")),
			resource.TestCheckResourceAttrWith("data.slack_conversation.channel", "is_archived", tb.ExpectBool(c.IsArchived)),
			resource.TestCheckResourceAttrWith("data.slack_conversation.channel", "is_shared", tb.ExpectBool(c.IsShared)),
			resource.TestCheckResourceAttrWith("data.slack_conversation.channel", "is_ext_shared", tb.ExpectBool(c.IsExtShared)),
			resource.TestCheckResourceAttrWith("data.slack_conversation.channel", "is_org_shared", tb.ExpectBool(c.IsOrgShared)),
			resource.TestCheckResourceAttrWith("data.slack_conversation.channel", "is_general", tb.ExpectBool(c.IsGeneral)),
		),
	})
}

func Test_DataSource_Conversation_Error_When_RetrievalFailed(t *testing.T) {
	testConfig(t, resource.TestStep{
		// arrange
		PreConfig: func() {
			m := tb.MockSlackClient()
			m.EXPECT().GetConversationInfo(gomock.Any(), gomock.Any()).Return(nil, errors.New("<SLACK_ERROR>")).AnyTimes()
		},
		Config: `
			provider slack {
				slack_token = "<SLACK_TOKEN>"
			}

			data "slack_conversation" "channel" {
				channel_id = "<GIVEN_ID>"
			}
		`,
		// assert
		ExpectError: regexp.MustCompile("<SLACK_ERROR>"),
	})
}

func Test_DataSource_Conversation_Error_WhenSlackClientNil(t *testing.T) {
	// arrange
	res := &datasource.ConfigureResponse{}
	req := datasource.ConfigureRequest{
		ProviderData: &SlackProviderData{
			Client: nil,
		},
	}

	test_instance := ConversationDataSource{}

	// act
	test_instance.Configure(context.Background(), req, res)

	// assert
	if res.Diagnostics.Errors()[0].Summary() != "Invalid Provider Data" {
		t.Errorf("Expected error summary to be 'Invalid Provider Data', got: %s", res.Diagnostics.Errors()[0].Summary())
	}
}
