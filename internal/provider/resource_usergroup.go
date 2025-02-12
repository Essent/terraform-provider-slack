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
	_ resource.Resource                = &UserGroupResource{}
	_ resource.ResourceWithImportState = &UserGroupResource{}
	_ resource.ResourceWithModifyPlan  = &UserGroupResource{}
)

func NewUserGroupResource() resource.Resource {
	return &UserGroupResource{}
}

type UserGroupResource struct {
	client slackExt.Client
}

type UserGroupResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Handle           types.String `tfsdk:"handle"`
	Channels         types.List   `tfsdk:"channels"`
	Users            types.List   `tfsdk:"users"`
	PreventConflicts types.Bool   `tfsdk:"prevent_conflicts"`
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
		MarkdownDescription: `Manages a Slack user group.

This resource requires the following scopes:

- usergroups:write
- usergroups:read

If you get missing_scope errors while using this resource check the scopes against the documentation for the methods above.`,

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
				Required: true,
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
				Description: "List of user IDs in the user group.",
			},
			"prevent_conflicts": schema.BoolAttribute{
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Optional:    true,
				Description: "If true, the plan fails if there's an enabled user group with the same name or handle.",
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

func (r *UserGroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !req.State.Raw.IsNull() {
		return
	}
	if r.client == nil {
		return
	}

	var plan UserGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newResource := plan.ID.IsNull() || plan.ID.IsUnknown()
	if !plan.PreventConflicts.ValueBool() || !newResource {
		return
	}

	name := plan.Name.ValueString()
	existingByName, errNameLookup := findUserGroupByField(ctx, name, "name", false, r.client)
	if errNameLookup == nil {
		resp.Diagnostics.AddError(
			"Conflict: Existing Enabled Group With Same Name",
			fmt.Sprintf("An enabled user group named %q already exists (ID: %s).", existingByName.Name, existingByName.ID),
		)
	} else if !strings.Contains(errNameLookup.Error(), "no usergroup with name") {
		resp.Diagnostics.AddError("Error Checking Name Conflict", errNameLookup.Error())
	}

	handle := plan.Handle.ValueString()
	existingByHandle, errHandleLookup := findUserGroupByField(ctx, handle, "handle", false, r.client)
	if errHandleLookup == nil {
		resp.Diagnostics.AddError(
			"Conflict: Existing Enabled Group With Same Handle",
			fmt.Sprintf("An enabled user group with handle %q already exists (ID: %s).", existingByHandle.Handle, existingByHandle.ID),
		)
	} else if !strings.Contains(errHandleLookup.Error(), "no usergroup with handle") {
		resp.Diagnostics.AddError("Error Checking Handle Conflict", errHandleLookup.Error())
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

	channels := listToStringSlice(plan.Channels)

	createReq := slack.UserGroup{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Handle:      plan.Handle.ValueString(),
		Prefs: slack.UserGroupPrefs{
			Channels: channels,
		},
	}

	created, errCreate := r.client.CreateUserGroup(ctx, createReq)
	if errCreate != nil {
		var lookupField, lookupValue string
		switch errCreate.Error() {
		case "name_already_exists":
			lookupField = "name"
			lookupValue = createReq.Name
		case "handle_already_exists":
			lookupField = "handle"
			lookupValue = createReq.Handle
		}

		var existingGroup slack.UserGroup
		if lookupField != "" {
			var errLookup error
			existingGroup, errLookup = findUserGroupByField(ctx, lookupValue, lookupField, true, r.client)
			if errLookup != nil {
				resp.Diagnostics.AddError(
					"Create Error",
					fmt.Sprintf("Slack returned %q, and %q when trying to find group with %s : %s", errCreate.Error(), errLookup.Error(), lookupField, lookupValue),
				)
				return
			}

			if existingGroup.DateDelete == 0 {
				resp.Diagnostics.AddError(
					"Create Error",
					fmt.Sprintf(
						"Conflict when creating group '%s' (conflicts with group ID: %s). Cannot reuse an enabled group.",
						createReq.Name, existingGroup.ID,
					),
				)
				return
			}

			if errEnable := r.enableAndUpdateUserGroup(ctx, existingGroup.ID, plan, channels); errEnable != nil {
				resp.Diagnostics.AddError("Enable/Update Error", errEnable.Error())
				return
			}
			plan.ID = types.StringValue(existingGroup.ID)
		} else {
			resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Error creating user group: %q", errCreate.Error()))
			return
		}
	} else {
		plan.ID = types.StringValue(created.ID)
		if err := r.updateUserGroupMembership(ctx, plan.ID.ValueString(), plan.Users); err != nil {
			resp.Diagnostics.AddError("Members Update Error", err.Error())
			return
		}
	}

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserGroupResource) enableAndUpdateUserGroup(
	ctx context.Context,
	groupID string,
	plan UserGroupResourceModel,
	channels []string,
) error {
	_, err := r.client.EnableUserGroup(ctx, groupID)
	if err != nil && err.Error() != "already_enabled" {
		return fmt.Errorf("could not enable usergroup %s: %w", groupID, err)
	}

	opts := []slack.UpdateUserGroupsOption{
		slack.UpdateUserGroupsOptionName(plan.Name.ValueString()),
		slack.UpdateUserGroupsOptionHandle(plan.Handle.ValueString()),
		slack.UpdateUserGroupsOptionDescription(&[]string{plan.Description.ValueString()}[0]),
		slack.UpdateUserGroupsOptionChannels(channels),
	}

	if _, err := r.client.UpdateUserGroup(ctx, groupID, opts...); err != nil {
		return fmt.Errorf("could not update usergroup %s: %w", groupID, err)
	}

	if err := r.updateUserGroupMembership(ctx, groupID, plan.Users); err != nil {
		return fmt.Errorf("could not update usergroup members: %w", err)
	}

	return nil
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

	found, err := findUserGroupByField(ctx, state.ID.ValueString(), "id", false, r.client)
	if err != nil {
		tflog.Warn(ctx, "Usergroup not found in Slack; removing from state", map[string]interface{}{
			"id": state.ID.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	state.UpdateFromUserGroup(&found)
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

	channels := listToStringSlice(plan.Channels)

	if err := r.enableAndUpdateUserGroup(ctx, state.ID.ValueString(), plan, channels); err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
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
	found, err := findUserGroupByField(ctx, model.ID.ValueString(), "id", false, r.client)
	if err != nil {
		tflog.Warn(ctx, "User group not found after create/update", map[string]interface{}{
			"id": model.ID.ValueString(),
		})
		return fmt.Errorf("user group with ID %s not found: %w", model.ID.ValueString(), err)
	}
	model.UpdateFromUserGroup(&found)
	return nil
}

func (r *UserGroupResource) updateUserGroupMembership(
	ctx context.Context,
	groupID string,
	userList types.List,
) error {
	users := listToStringSlice(userList)
	usersParam := strings.Join(users, ",")
	if len(users) == 0 {
		usersParam = "[]"
	}

	_, err := r.client.UpdateUserGroupMembers(ctx, groupID, usersParam)
	if err != nil {
		return fmt.Errorf("could not update usergroup members: %s", err)
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
		case "id":
			matches = (g.ID == searchVal)
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
