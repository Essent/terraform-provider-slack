// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/essent/terraform-provider-slack/internal/slackExt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/slack-go/slack"
)

var (
	_ resource.Resource                   = &UserGroupResource{}
	_ resource.ResourceWithImportState    = &UserGroupResource{}
	_ resource.ResourceWithValidateConfig = &UserGroupResource{}
)

func NewUserGroupResource() resource.Resource {
	return &UserGroupResource{}
}

type UserGroupResource struct {
	client slackExt.Client
}

type UserGroupResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	Handle                types.String `tfsdk:"handle"`
	Channels              types.List   `tfsdk:"channels"`
	Users                 types.List   `tfsdk:"users"`
	PreventDuplicateNames types.Bool   `tfsdk:"prevent_duplicate_names"`
}

func (r *UserGroupResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_usergroup"
}

func (r *UserGroupResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Slack user group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"handle": schema.StringAttribute{
				Optional: true,
			},
			"channels": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Channels shared by the user group.",
				Default: listdefault.StaticValue(
					types.ListValueMust(types.StringType, []attr.Value{}),
				),
			},
			"users": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: listdefault.StaticValue(
					types.ListValueMust(types.StringType, []attr.Value{}),
				),
				MarkdownDescription: "List of user IDs in the user group.",
			},
			"prevent_duplicate_names": schema.BoolAttribute{
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Optional:    true,
				Description: "If true, the plan fails if there's an enabled user group with the same name or handle (checked during plan).",
			},
		},
	}
}

func (r *UserGroupResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}
	pd, ok := req.ProviderData.(*SlackProviderData)
	if !ok || pd.Client == nil {
		resp.Diagnostics.AddError("Invalid Provider Data", "Could not create Slack client.")
		return
	}
	r.client = pd.Client
}

func (r *UserGroupResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	if r.client == nil {
		return
	}

	var plan UserGroupResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newResource := plan.ID.IsNull() || plan.ID.IsUnknown()
	if !plan.PreventDuplicateNames.ValueBool() || !newResource {
		return
	}

	name := plan.Name.ValueString()
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() && name != "" {
		existingByName, err := findUserGroupByField(ctx, name, "name", false, r.client)
		if err == nil {
			resp.Diagnostics.AddError(
				"Conflict: Existing Enabled Group With Same Name",
				fmt.Sprintf("An enabled user group named %q already exists (ID: %s).", existingByName.Name, existingByName.ID),
			)
		} else if !strings.Contains(err.Error(), "no usergroup with name") {
			resp.Diagnostics.AddError("Error Checking Name Conflict", err.Error())
		}
	}

	handle := plan.Handle.ValueString()
	if !plan.Handle.IsNull() && !plan.Handle.IsUnknown() && handle != "" {
		existingByHandle, err := findUserGroupByField(ctx, handle, "handle", false, r.client)
		if err == nil {
			resp.Diagnostics.AddError(
				"Conflict: Existing Enabled Group With Same Handle",
				fmt.Sprintf("An enabled user group with handle %q already exists (ID: %s).", existingByHandle.Handle, existingByHandle.ID),
			)
		} else if !strings.Contains(err.Error(), "no usergroup with handle") {
			resp.Diagnostics.AddError("Error Checking Handle Conflict", err.Error())
		}
	}
}

func (r *UserGroupResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan UserGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Handle.IsNull() || plan.Handle.ValueString() == "" {
		plan.Handle = plan.Name
	}

	desiredName := plan.Name.ValueString()
	desiredHandle := plan.Handle.ValueString()

	channels := listToStringSlice(plan.Channels)
	users := listToStringSlice(plan.Users)

	createReq := slack.UserGroup{
		Name:        desiredName,
		Description: plan.Description.ValueString(),
		Handle:      desiredHandle,
		Prefs: slack.UserGroupPrefs{
			Channels: channels,
		},
	}

	created, err := r.client.CreateUserGroup(ctx, createReq)
	if err != nil {
		var existingGroup slack.UserGroup
		var err2 error
		if err.Error() == "name_already_exists" {
			existingGroup, err2 = findUserGroupByField(ctx, desiredName, "name", true, r.client)
			if err2 != nil {
				resp.Diagnostics.AddError(
					"Create Error",
					fmt.Sprintf("Slack returned %q, then could not find group: %s", err.Error(), err2),
				)
				return
			}
		}
		if err.Error() == "handle_already_exists" {
			existingGroup, err2 = findUserGroupByField(ctx, desiredHandle, "handle", true, r.client)
			if err2 != nil {
				resp.Diagnostics.AddError(
					"Create Error",
					fmt.Sprintf("Slack returned %q, then could not find group: %s", err.Error(), err2),
				)
				return
			}
		}

		if existingGroup.DateDelete == 0 {
			resp.Diagnostics.AddError(
				"Create Error",
				fmt.Sprintf(
					"Conflict when creating group '%s': %q (ID: %s). Cannot reuse an enabled group.",
					desiredName, err.Error(), existingGroup.ID,
				),
			)
			return
		}

		_, err2 = r.client.EnableUserGroup(ctx, existingGroup.ID)
		if err2 != nil && err2.Error() != "already_enabled" {
			resp.Diagnostics.AddError("Enable Error", fmt.Sprintf("Could not enable usergroup %s: %s", existingGroup.ID, err2))
			return
		}
		opts := []slack.UpdateUserGroupsOption{
			slack.UpdateUserGroupsOptionName(plan.Name.ValueString()),
			slack.UpdateUserGroupsOptionHandle(plan.Handle.ValueString()),
			slack.UpdateUserGroupsOptionDescription(&[]string{plan.Description.ValueString()}[0]),
			slack.UpdateUserGroupsOptionChannels(channels),
		}
		_, err2 = r.client.UpdateUserGroup(ctx, existingGroup.ID, opts...)
		if err2 != nil {
			resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Could not update usergroup %s: %s", existingGroup.ID, err2))
			return
		}
		plan.ID = types.StringValue(existingGroup.ID)
	} else {
		plan.ID = types.StringValue(created.ID)
	}

	usersParam := strings.Join(users, ",")
	if len(users) == 0 {
		usersParam = "[]"
	}
	_, err = r.client.UpdateUserGroupMembers(ctx, plan.ID.ValueString(), usersParam)
	if err != nil {
		resp.Diagnostics.AddError("Members Update Error", fmt.Sprintf("Could not update usergroup members: %s", err))
		return
	}

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserGroupResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state UserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groups, err := r.client.GetUserGroups(ctx, slack.GetUserGroupsOptionIncludeUsers(true))
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Could not retrieve user groups: %s", err))
		return
	}

	found := findGroupByID(groups, state.ID.ValueString())
	if found == nil {
		tflog.Warn(ctx, "User group not found in Slack; removing from state", map[string]interface{}{
			"id": state.ID.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	setStateFromUserGroup(found, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserGroupResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state UserGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Handle.IsNull() || plan.Handle.ValueString() == "" {
		plan.Handle = plan.Name
	}

	channels := listToStringSlice(plan.Channels)
	users := listToStringSlice(plan.Users)

	opts := []slack.UpdateUserGroupsOption{
		slack.UpdateUserGroupsOptionName(plan.Name.ValueString()),
		slack.UpdateUserGroupsOptionHandle(plan.Handle.ValueString()),
		slack.UpdateUserGroupsOptionDescription(&[]string{plan.Description.ValueString()}[0]),
		slack.UpdateUserGroupsOptionChannels(channels),
	}

	_, err := r.client.UpdateUserGroup(ctx, state.ID.ValueString(), opts...)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Could not update usergroup: %s", err))
		return
	}

	if !plan.Users.Equal(state.Users) {
		usersParam := strings.Join(users, ",")
		if len(users) == 0 {
			usersParam = "[]"
		}
		_, err = r.client.UpdateUserGroupMembers(ctx, state.ID.ValueString(), usersParam)
		if err != nil {
			resp.Diagnostics.AddError("Members Update Error", fmt.Sprintf("Could not update usergroup members: %s", err))
			return
		}
	}

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserGroupResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state UserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DisableUserGroup(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Could not disable usergroup: %s", err))
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *UserGroupResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *UserGroupResource) readIntoModel(ctx context.Context, model *UserGroupResourceModel) error {
	groups, err := r.client.GetUserGroups(ctx, slack.GetUserGroupsOptionIncludeUsers(true))
	if err != nil {
		return fmt.Errorf("could not read user group: %w", err)
	}
	found := findGroupByID(groups, model.ID.ValueString())
	if found == nil {
		tflog.Warn(ctx, "User group not found after create/update", map[string]interface{}{
			"id": model.ID.ValueString(),
		})
		return nil
	}
	setStateFromUserGroup(found, model)
	return nil
}

func listToStringSlice(l types.List) []string {
	if l.IsNull() || l.IsUnknown() {
		return []string{}
	}
	elems := l.Elements()
	result := make([]string, 0, len(elems))
	for _, e := range elems {
		if s, ok := e.(types.String); ok && !s.IsNull() && !s.IsUnknown() {
			result = append(result, s.ValueString())
		}
	}
	return result
}

func stringSliceToList(list []string) types.List {
	if len(list) == 0 {
		emptyVal, _ := types.ListValue(types.StringType, []attr.Value{})
		return emptyVal
	}
	attrValues := make([]attr.Value, len(list))
	for i, s := range list {
		attrValues[i] = types.StringValue(s)
	}
	res, diags := types.ListValue(types.StringType, attrValues)
	if diags.HasError() {
		return types.ListNull(types.StringType)
	}
	return res
}

func findGroupByID(groups []slack.UserGroup, id string) *slack.UserGroup {
	for i := range groups {
		if groups[i].ID == id {
			return &groups[i]
		}
	}
	return nil
}

func findUserGroupByField(
	ctx context.Context,
	searchVal, searchField string,
	includeDisabled bool,
	client slackExt.Client,
) (slack.UserGroup, error) {
	groups, err := client.GetUserGroups(ctx,
		slack.GetUserGroupsOptionIncludeDisabled(includeDisabled),
		slack.GetUserGroupsOptionIncludeUsers(true),
	)
	if err != nil {
		return slack.UserGroup{}, err
	}

	for _, g := range groups {
		var matches bool
		switch searchField {
		case "name":
			matches = (g.Name == searchVal)
		case "handle":
			matches = (g.Handle == searchVal)
		default:
			continue
		}

		if matches {
			if !includeDisabled && g.DateDelete == 0 {
				return g, nil
			} else if includeDisabled {
				return g, nil
			}
		}
	}

	return slack.UserGroup{}, fmt.Errorf("no usergroup with %s %q found", searchField, searchVal)
}

func setStateFromUserGroup(ug *slack.UserGroup, model *UserGroupResourceModel) {
	model.ID = types.StringValue(ug.ID)
	model.Name = types.StringValue(ug.Name)
	model.Description = types.StringValue(ug.Description)
	model.Handle = types.StringValue(ug.Handle)
	model.Channels = stringSliceToList(ug.Prefs.Channels)
	model.Users = stringSliceToList(ug.Users)
}
