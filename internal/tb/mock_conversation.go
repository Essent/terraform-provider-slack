// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tb

import "github.com/slack-go/slack"

type ChannelBuilder struct {
	result *slack.Channel
}

func (b *ChannelBuilder) Build() *slack.Channel {
	return b.result
}

func (b *ChannelBuilder) WithID(id string) *ChannelBuilder {
	b.result.ID = id
	return b
}

func (b *ChannelBuilder) WithTopic(topic string) *ChannelBuilder {
	b.result.Topic.Value = topic
	return b
}

func (b *ChannelBuilder) WithPurpose(purpose string) *ChannelBuilder {
	b.result.Purpose.Value = purpose
	return b
}

func (b *ChannelBuilder) WithCreated(created slack.JSONTime) *ChannelBuilder {
	b.result.Created = created
	return b
}

func (b *ChannelBuilder) WithCreator(creator string) *ChannelBuilder {
	b.result.Creator = creator
	return b
}

func (b *ChannelBuilder) WithIsArchived(isArchived bool) *ChannelBuilder {
	b.result.IsArchived = isArchived
	return b
}

func (b *ChannelBuilder) WithIsShared(isShared bool) *ChannelBuilder {
	b.result.IsShared = isShared
	return b
}

func (b *ChannelBuilder) WithIsExtShared(isExtShared bool) *ChannelBuilder {
	b.result.IsExtShared = isExtShared
	return b
}

func (b *ChannelBuilder) WithIsOrgShared(isOrgShared bool) *ChannelBuilder {
	b.result.IsOrgShared = isOrgShared
	return b
}

func (b *ChannelBuilder) WithIsGeneral(isGeneral bool) *ChannelBuilder {
	b.result.IsGeneral = isGeneral
	return b
}

func NewChannelBuilder() *ChannelBuilder {
	return &ChannelBuilder{
		result: &slack.Channel{},
	}
}
