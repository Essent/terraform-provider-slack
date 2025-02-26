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

var _ datasource.DataSource = &AllUserGroupsDataSource{}

func NewAllUserGroupsDataSource() datasource.DataSource {
	return &AllUserGroupsDataSource{}
}

type AllUserGroupsDataSource struct {
	client slackExt.Client
}

type AllUserGroupsDataSourceModel struct {
	TotalUserGroups types.Int64                        `tfsdk:"total_usergroups"`
	UserGroups      []AllUserGroupsDataSourceGroupItem `tfsdk:"usergroups"`
}

type AllUserGroupsDataSourceGroupItem struct {
	ID          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Handle      types.String   `tfsdk:"handle"`
	Channels    []types.String `tfsdk:"channels"`
	Users       []types.String `tfsdk:"users"`
}

func (d *AllUserGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_all_usergroups"
}

func (d *AllUserGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Retrieve a list of all Slack user groups.

This datasource requires the following scopes:

- usergroups:read`,
		Description: "Retrieve all Slack user groups.",
		Attributes: map[string]schema.Attribute{
			"total_usergroups": schema.Int64Attribute{
				Description: "Total number of user groups retrieved.",
				Computed:    true,
			},
			"usergroups": schema.ListNestedAttribute{
				Description: "List of Slack user groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "User group's Slack ID.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the user group.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the user group.",
							Computed:    true,
						},
						"handle": schema.StringAttribute{
							Description: "Handle of the user group (unique identifier).",
							Computed:    true,
						},
						"channels": schema.ListAttribute{
							Description: "Channels shared by the user group.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"users": schema.ListAttribute{
							Description: "List of user IDs in the user group.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *AllUserGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AllUserGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AllUserGroupsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userGroups, err := d.client.GetUserGroups(ctx, slack.GetUserGroupsOptionIncludeUsers(true))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to fetch Slack user groups: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "Fetched Slack user groups", map[string]any{"total_usergroups": len(userGroups)})

	var resultingList []AllUserGroupsDataSourceGroupItem
	for _, group := range userGroups {
		groupItem := AllUserGroupsDataSourceGroupItem{
			ID:          types.StringValue(group.ID),
			Name:        types.StringValue(group.Name),
			Description: types.StringValue(group.Description),
			Handle:      types.StringValue(group.Handle),
		}

		channels := make([]types.String, len(group.Prefs.Channels))
		for i, ch := range group.Prefs.Channels {
			channels[i] = types.StringValue(ch)
		}
		groupItem.Channels = channels

		users := make([]types.String, len(group.Users))
		for i, u := range group.Users {
			users[i] = types.StringValue(u)
		}
		groupItem.Users = users

		resultingList = append(resultingList, groupItem)
	}

	data.UserGroups = resultingList
	data.TotalUserGroups = types.Int64Value(int64(len(resultingList)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
