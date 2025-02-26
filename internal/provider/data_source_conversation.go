// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/essent/terraform-provider-slack/internal/slackExt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/slack-go/slack"
)

var _ datasource.DataSource = &ConversationDataSource{}

func NewConversationDataSource() datasource.DataSource {
	return &ConversationDataSource{}
}

type ConversationDataSource struct {
	client slackExt.Client
}

type ConversationDataSourceModel struct {
	ChannelID   types.String `tfsdk:"channel_id"`
	Topic       types.String `tfsdk:"topic"`
	Purpose     types.String `tfsdk:"purpose"`
	Created     types.Int64  `tfsdk:"created"`
	Creator     types.String `tfsdk:"creator"`
	IsArchived  types.Bool   `tfsdk:"is_archived"`
	IsShared    types.Bool   `tfsdk:"is_shared"`
	IsExtShared types.Bool   `tfsdk:"is_ext_shared"`
	IsOrgShared types.Bool   `tfsdk:"is_org_shared"`
	IsGeneral   types.Bool   `tfsdk:"is_general"`
}

func (d *ConversationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_conversation"
}

func (d *ConversationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Retrieve information about a Slack conversation by its ID.

This datasource requires the following scopes:

- channels:read (public channels)
- groups:read (private channels)
- im:read (optional)
- mpim:read (optional)`,
		Attributes: map[string]schema.Attribute{
			"channel_id": schema.StringAttribute{
				MarkdownDescription: "The Slack channel ID to look up.",
				Required:            true,
			},
			"topic": schema.StringAttribute{
				MarkdownDescription: "The channel topic.",
				Computed:            true,
			},
			"purpose": schema.StringAttribute{
				MarkdownDescription: "The channel purpose.",
				Computed:            true,
			},
			"created": schema.Int64Attribute{
				MarkdownDescription: "UNIX timestamp when the channel was created.",
				Computed:            true,
			},
			"creator": schema.StringAttribute{
				MarkdownDescription: "User ID of the channel creator.",
				Computed:            true,
			},
			"is_archived": schema.BoolAttribute{
				MarkdownDescription: "True if the channel is archived.",
				Computed:            true,
			},
			"is_shared": schema.BoolAttribute{
				MarkdownDescription: "True if the channel is shared.",
				Computed:            true,
			},
			"is_ext_shared": schema.BoolAttribute{
				MarkdownDescription: "True if the channel is externally shared.",
				Computed:            true,
			},
			"is_org_shared": schema.BoolAttribute{
				MarkdownDescription: "True if the channel is shared across the org.",
				Computed:            true,
			},
			"is_general": schema.BoolAttribute{
				MarkdownDescription: "True if this is the #general channel.",
				Computed:            true,
			},
		},
	}
}

func (d *ConversationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*SlackProviderData)
	if !ok || providerData.Client == nil {
		resp.Diagnostics.AddError(
			"Invalid Provider Data",
			fmt.Sprintf("Expected *SlackProviderData with initialized client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = providerData.Client
}

func (d *ConversationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConversationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channel, err := d.client.GetConversationInfo(ctx, &slack.GetConversationInfoInput{
		ChannelID:         data.ChannelID.ValueString(),
		IncludeLocale:     false,
		IncludeNumMembers: false,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Could not get conversation info for channel %s: %s", data.ChannelID.ValueString(), err),
		)
		return
	}

	data.ChannelID = types.StringValue(channel.ID)
	data.Topic = types.StringValue(channel.Topic.Value)
	data.Purpose = types.StringValue(channel.Purpose.Value)
	data.Created = types.Int64Value(int64(channel.Created))
	data.Creator = types.StringValue(channel.Creator)
	data.IsArchived = types.BoolValue(channel.IsArchived)
	data.IsShared = types.BoolValue(channel.IsShared)
	data.IsExtShared = types.BoolValue(channel.IsExtShared)
	data.IsOrgShared = types.BoolValue(channel.IsOrgShared)
	data.IsGeneral = types.BoolValue(channel.IsGeneral)

	tflog.Trace(ctx, "Fetched Slack channel data", map[string]any{"channel_id": channel.ID})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
