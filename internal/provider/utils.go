// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/slack-go/slack"
)

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

func setToStringSlice(s types.Set) []string {
	res := make([]string, 0, len(s.Elements()))
	for _, e := range s.Elements() {
		if str, ok := e.(types.String); ok {
			res = append(res, str.ValueString())
		}
	}
	return res
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

func stringSliceToSet(list []string) types.Set {
	attrValues := make([]attr.Value, len(list))
	for i, v := range list {
		attrValues[i] = types.StringValue(v)
	}
	return types.SetValueMust(types.StringType, attrValues)
}

// Slack does not preserve the order of prefs.channels, so a read-back keeps the
// order Terraform already holds. Returning a different order for the same
// channels breaks the guarantee that a known planned value survives apply.
func orderedLike(reference, actual []string) []string {
	remaining := make(map[string]bool, len(actual))
	for _, s := range actual {
		remaining[s] = true
	}

	ordered := make([]string, 0, len(actual))
	for _, group := range [][]string{reference, actual} {
		for _, s := range group {
			if remaining[s] {
				ordered = append(ordered, s)
				delete(remaining, s)
			}
		}
	}
	return ordered
}

func (m *UserGroupResourceModel) UpdateFromUserGroup(ug *slack.UserGroup) {
	m.ID = types.StringValue(ug.ID)
	m.Name = types.StringValue(ug.Name)
	m.Description = types.StringValue(ug.Description)
	m.Handle = types.StringValue(ug.Handle)
	m.Channels = stringSliceToList(orderedLike(listToStringSlice(m.Channels), ug.Prefs.Channels))
	m.Users = stringSliceToSet(ug.Users)
}
