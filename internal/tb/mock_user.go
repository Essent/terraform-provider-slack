// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tb

import "github.com/slack-go/slack"

type UserBuilder struct {
	result *slack.User
}

func (b *UserBuilder) Build() *slack.User {
	return b.result
}

func (b *UserBuilder) WithID(id string) *UserBuilder {
	b.result.ID = id
	return b
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.result.Name = name
	return b
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.result.Profile.Email = email
	return b
}

func (b *UserBuilder) WithDeleted(deleted bool) *UserBuilder {
	b.result.Deleted = deleted
	return b
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		result: &slack.User{},
	}
}
