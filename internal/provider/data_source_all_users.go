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
)

var _ datasource.DataSource = &AllUsersDataSource{}

func NewAllUsersDataSource() datasource.DataSource {
	return &AllUsersDataSource{}
}

type AllUsersDataSource struct {
	client slackExt.Client
}

type AllUsersDataSourceModel struct {
	Totalusers types.Int64                       `tfsdk:"total_users"`
	Users      []AllUsersDataSourceModelUserItem `tfsdk:"users"`
}

type AllUsersDataSourceModelUserItem struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

func (d *AllUsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_all_users"
}

func (d *AllUsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Retrieve a list of all activated and non-bot Slack users.

This datasource requires the following scopes:

- users:read

If you get missing_scope errors while using this datasource, check the scopes against the documentation for the methods above.`,
		Attributes: map[string]schema.Attribute{
			"total_users": schema.Int64Attribute{
				Description: "Number of users returned.",
				Computed:    true,
			},
			"users": schema.ListNestedAttribute{
				Description: "List of activated and non-bot Slack users.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "User's Slack ID.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User's name.",
							Computed:    true,
						},
						"email": schema.StringAttribute{
							Description: "User's email address.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *AllUsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AllUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AllUsersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	users, err := d.client.GetUsersContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to fetch Slack users: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "Fetched Slack users", map[string]any{"total_users": len(users)})

	var resultingList []AllUsersDataSourceModelUserItem
	for _, user := range users {
		if !user.Deleted && !user.IsBot {
			resultingList = append(resultingList, AllUsersDataSourceModelUserItem{
				ID:    types.StringValue(user.ID),
				Name:  types.StringValue(user.Name),
				Email: types.StringValue(user.Profile.Email),
			})
		}
	}

	data.Users = resultingList
	data.Totalusers = types.Int64Value(int64(len(resultingList)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
