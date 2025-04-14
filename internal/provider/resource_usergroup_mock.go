// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UserGroupResourceModelBuilder struct {
	result *UserGroupResourceModel
}

func (b *UserGroupResourceModelBuilder) Build() *UserGroupResourceModel {
	return b.result
}

func NewUserGroupResourceModelBuilder() *UserGroupResourceModelBuilder {
	return &UserGroupResourceModelBuilder{&UserGroupResourceModel{}}
}

func (b *UserGroupResourceModelBuilder) WithID(value string) *UserGroupResourceModelBuilder {
	b.result.ID = types.StringValue(value)
	return b
}

func (b *UserGroupResourceModelBuilder) WithName(value string) *UserGroupResourceModelBuilder {
	b.result.Name = types.StringValue(value)
	return b
}

func (b *UserGroupResourceModelBuilder) WithDescription(value string) *UserGroupResourceModelBuilder {
	b.result.Description = types.StringValue(value)
	return b
}

func (b *UserGroupResourceModelBuilder) WithHandle(value string) *UserGroupResourceModelBuilder {
	b.result.Handle = types.StringValue(value)
	return b
}

func (b *UserGroupResourceModelBuilder) WithChannels(value []string) *UserGroupResourceModelBuilder {
	attrValues := make([]attr.Value, len(value))
	for i, v := range value {
		attrValues[i] = types.StringValue(v)
	}
	b.result.Channels = types.ListValueMust(types.StringType, attrValues)

	return b
}

func (b *UserGroupResourceModelBuilder) WithUsers(value []string) *UserGroupResourceModelBuilder {
	attrValues := make([]attr.Value, len(value))
	for i, v := range value {
		attrValues[i] = types.StringValue(v)
	}
	b.result.Users = types.SetValueMust(types.StringType, attrValues)

	return b
}

func (b *UserGroupResourceModelBuilder) WithPreventConflicts(value bool) *UserGroupResourceModelBuilder {
	b.result.PreventConflicts = types.BoolValue(value)
	return b
}
