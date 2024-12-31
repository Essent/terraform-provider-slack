// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/slack-go/slack"
)

var _ provider.Provider = &SlackProvider{}
var _ provider.ProviderWithFunctions = &SlackProvider{}

type SlackProvider struct {
	version string
}

type SlackProviderModel struct {
	SlackToken types.String `tfsdk:"slack_token"`
}

type SlackProviderData struct {
	Client *slack.Client
}

func (p *SlackProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "slack"
	resp.Version = p.version
}

func (p *SlackProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"slack_token": schema.StringAttribute{
				MarkdownDescription: "Slack token to authenticate API calls.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *SlackProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SlackProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.SlackToken.IsNull() || data.SlackToken.IsUnknown() {
		resp.Diagnostics.AddError(
			"Missing Slack Token",
			"The `slack_token` must be provided to authenticate API calls.",
		)
		return
	}
	tflog.Info(ctx, "Configuring slack client")
	client := slack.New(data.SlackToken.ValueString())
	_, err := client.AuthTest()
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Slack Token",
			fmt.Sprintf("Unable to authenticate with Slack: %s", err),
		)
		return
	}

	resp.DataSourceData = &SlackProviderData{Client: client}
	resp.ResourceData = &SlackProviderData{Client: client}
}

func (p *SlackProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Add your resources here if needed
	}
}

func (p *SlackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserDataSource,
	}
}

func (p *SlackProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// Add your provider-level functions here if needed
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SlackProvider{version: version}
	}
}
