// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/slack-go/slack"
)

var _ datasource.DataSource = &UserDataSource{}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

type UserDataSource struct {
	client *slack.Client
}

type UserDataSourceModel struct {
	UserID      types.String `tfsdk:"user_id"`
	Email       types.String `tfsdk:"email"`
	RealName    types.String `tfsdk:"real_name"`
	DisplayName types.String `tfsdk:"display_name"`
	ID          types.String `tfsdk:"id"`
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_data"
}

func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve Slack user information. Either `user_id` or `email` must be specified, but not both.",
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				MarkdownDescription: "Slack user ID to look up.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("email")),
				},
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email of the user to look up.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("user_id")),
				},
			},
			"real_name": schema.StringAttribute{
				MarkdownDescription: "User's real name.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "User's display name.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for Terraform state.",
				Computed:            true,
			},
		},
	}
}

func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// https://stackoverflow.com/questions/78623763/terraform-provider-method-configure-not-getting-called
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		user *slack.User
		err  error
	)

	if !data.UserID.IsNull() {
		user, err = d.client.GetUserInfo(data.UserID.ValueString())
	} else {
		user, err = d.client.GetUserByEmail(data.Email.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to fetch user info: %s", err),
		)
		return
	}

	if user.Deleted {
		resp.Diagnostics.AddError(
			"User is deactivated",
			"User is deactivated in Slack",
		)
		return
	}

	data.Email = types.StringValue(user.Profile.Email)
	data.RealName = types.StringValue(user.RealName)
	data.DisplayName = types.StringValue(user.Profile.DisplayName)
	data.ID = types.StringValue(user.ID)

	tflog.Trace(ctx, "Fetched Slack user data", map[string]any{"user_id": user.ID})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
