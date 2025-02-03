package tb

import "github.com/slack-go/slack"

type UsergroupBuilder struct {
	result *slack.UserGroup
}

func (b *UsergroupBuilder) Build() *slack.UserGroup {
	return b.result
}

func (b *UsergroupBuilder) WithID(id string) *UsergroupBuilder {
	b.result.ID = id
	return b
}

func (b *UsergroupBuilder) WithName(name string) *UsergroupBuilder {
	b.result.Name = name
	return b
}

func (b *UsergroupBuilder) WithDescription(description string) *UsergroupBuilder {
	b.result.Description = description
	return b
}

func (b *UsergroupBuilder) WithHandle(handle string) *UsergroupBuilder {
	b.result.Handle = handle
	return b
}

func (b *UsergroupBuilder) WithChannels(channels []string) *UsergroupBuilder {
	b.result.Prefs.Channels = channels
	return b
}

func (b *UsergroupBuilder) WithUsers(users []string) *UsergroupBuilder {
	b.result.Users = users
	return b
}

func NewUsergroupBuilder() *UsergroupBuilder {
	return &UsergroupBuilder{
		result: &slack.UserGroup{},
	}
}
