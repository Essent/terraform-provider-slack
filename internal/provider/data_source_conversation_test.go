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

func Test_DataSource_Conversation(t *testing.T) {
	// arrange
	cb := tb.NewChannelBuilder().WithID("<GIVEN_ID>").WithTopic("<TOPIC>").WithPurpose("<PURPOSE>")
	cb.WithCreated(1234567890).WithCreator("<CREATOR>").WithIsArchived(tb.RandBool()).WithIsShared(tb.RandBool())
	cb.WithIsExtShared(tb.RandBool()).WithIsOrgShared(tb.RandBool()).WithIsGeneral(tb.RandBool())

	c := cb.Build()

	testConfig(t, tr.TestStep{
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
		Check: tr.ComposeTestCheckFunc(
			tr.TestCheckResourceAttrSet("data.slack_conversation.channel", "channel_id"),
			tr.TestCheckResourceAttrSet("data.slack_conversation.channel", "topic"),
			tr.TestCheckResourceAttrSet("data.slack_conversation.channel", "purpose"),
			tr.TestCheckResourceAttrSet("data.slack_conversation.channel", "created"),
			tr.TestCheckResourceAttrSet("data.slack_conversation.channel", "creator"),
			tr.TestCheckResourceAttrSet("data.slack_conversation.channel", "is_archived"),
			tr.TestCheckResourceAttrSet("data.slack_conversation.channel", "is_shared"),
			tr.TestCheckResourceAttrSet("data.slack_conversation.channel", "is_ext_shared"),
			tr.TestCheckResourceAttrSet("data.slack_conversation.channel", "is_org_shared"),
			tr.TestCheckResourceAttrSet("data.slack_conversation.channel", "is_general"),

			tr.TestCheckResourceAttrWith("data.slack_conversation.channel", "channel_id", tb.ExpectString("<GIVEN_ID>")),
			tr.TestCheckResourceAttrWith("data.slack_conversation.channel", "topic", tb.ExpectString("<TOPIC>")),
			tr.TestCheckResourceAttrWith("data.slack_conversation.channel", "purpose", tb.ExpectString("<PURPOSE>")),
			tr.TestCheckResourceAttrWith("data.slack_conversation.channel", "created", tb.ExpectString("1234567890")),
			tr.TestCheckResourceAttrWith("data.slack_conversation.channel", "creator", tb.ExpectString("<CREATOR>")),
			tr.TestCheckResourceAttrWith("data.slack_conversation.channel", "is_archived", tb.ExpectBool(c.IsArchived)),
			tr.TestCheckResourceAttrWith("data.slack_conversation.channel", "is_shared", tb.ExpectBool(c.IsShared)),
			tr.TestCheckResourceAttrWith("data.slack_conversation.channel", "is_ext_shared", tb.ExpectBool(c.IsExtShared)),
			tr.TestCheckResourceAttrWith("data.slack_conversation.channel", "is_org_shared", tb.ExpectBool(c.IsOrgShared)),
			tr.TestCheckResourceAttrWith("data.slack_conversation.channel", "is_general", tb.ExpectBool(c.IsGeneral)),
		),
	})
}

func Test_DataSource_Conversation_Error_When_RetrievalFailed(t *testing.T) {
	testConfig(t, tr.TestStep{
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
