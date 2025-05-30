// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/essent/terraform-provider-slack/internal/provider/dependencies"
	"github.com/essent/terraform-provider-slack/internal/slackExt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ provider.Provider = &SlackProvider{}
var _ provider.ProviderWithFunctions = &SlackProvider{}

type SlackProvider struct {
	version      string
	dependencies dependencies.Dependencies
}

type SlackProviderModel struct {
	SlackToken types.String `tfsdk:"slack_token"`
}

type SlackProviderData struct {
	Client           slackExt.Client
	UserGroupService UserGroupService
}

func (p *SlackProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "slack"
	resp.Version = p.version
}

func (p *SlackProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Slack provider allows you to manage Slack resources. Use it to create, read, update, and delete Slack entities.",
		Attributes: map[string]schema.Attribute{
			"slack_token": schema.StringAttribute{
				MarkdownDescription: "The Slack token used for API authentication. It can be provided in the provider block or via the `SLACK_TOKEN` environment variable.",
				Optional:            true,
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

	slackToken := data.SlackToken.ValueString()

	// If slack_token was not set in the provider block, check the environment variable.
	if slackToken == "" {
		envToken, ok := os.LookupEnv("SLACK_TOKEN")
		if !ok || envToken == "" {
			resp.Diagnostics.AddError(
				"Missing Slack Token",
				"`slack_token` was not set in the provider block, and `SLACK_TOKEN` is not set in the environment.",
			)
			return
		}
		slackToken = envToken
	}

	tflog.Info(ctx, "Configuring slack client")
	client := p.dependencies.CreateSlackClient(slackToken)
	_, err := client.AuthTest(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Slack Token",
			fmt.Sprintf("Unable to authenticate with Slack: %s", err),
		)
		return
	}

	providerData := &SlackProviderData{
		Client:           client,
		UserGroupService: NewUserGroupService(client),
	}
	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

func (p *SlackProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserGroupResource,
	}
}

func (p *SlackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserDataSource,
		NewAllUsersDataSource,
		NewAllUserGroupsDataSource,
		NewConversationDataSource,
	}
}

func (p *SlackProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// Add your provider-level functions here if needed
	}
}

func New(version string, dependencies dependencies.Dependencies) func() provider.Provider {
	return func() provider.Provider {
		return &SlackProvider{version: version, dependencies: dependencies}
	}
}
