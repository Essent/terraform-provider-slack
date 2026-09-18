// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/slack-go/slack"
)

func Test_OrderedLike(t *testing.T) {
	cases := []struct {
		name      string
		reference []string
		actual    []string
		expected  []string
	}{
		{
			name:      "reorders a permutation to match the reference",
			reference: []string{"C1", "C2", "C3"},
			actual:    []string{"C3", "C1", "C2"},
			expected:  []string{"C1", "C2", "C3"},
		},
		{
			name:      "appends channels missing from the reference",
			reference: []string{"C2", "C1"},
			actual:    []string{"C4", "C1", "C3", "C2"},
			expected:  []string{"C2", "C1", "C4", "C3"},
		},
		{
			name:      "drops channels missing from the actual",
			reference: []string{"C1", "C2", "C3"},
			actual:    []string{"C3", "C1"},
			expected:  []string{"C1", "C3"},
		},
		{
			name:      "falls back to the actual order without a reference",
			reference: []string{},
			actual:    []string{"C2", "C1"},
			expected:  []string{"C2", "C1"},
		},
		{
			name:      "returns empty when the actual is empty",
			reference: []string{"C1"},
			actual:    []string{},
			expected:  []string{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// act
			result := orderedLike(c.reference, c.actual)

			// assert
			if !slices.Equal(result, c.expected) {
				t.Errorf("Expected %v, got: %v", c.expected, result)
			}
		})
	}
}

func Test_UpdateFromUserGroup_KeepsChannelOrder(t *testing.T) {
	// arrange
	model := NewUserGroupResourceModelBuilder().WithChannels([]string{"C1", "C2", "C3"}).Build()
	group := slack.UserGroup{
		Prefs: slack.UserGroupPrefs{Channels: []string{"C2", "C3", "C1"}},
	}

	// act
	model.UpdateFromUserGroup(&group)

	// assert
	expected := stringSliceToList([]string{"C1", "C2", "C3"})
	if !model.Channels.Equal(expected) {
		t.Errorf("Expected %v, got: %v", expected, model.Channels)
	}
}

func Test_UpdateFromUserGroup_UsesSlackOrder_WhenChannelsUnknown(t *testing.T) {
	// arrange
	model := &UserGroupResourceModel{Channels: types.ListUnknown(types.StringType)}
	group := slack.UserGroup{
		Prefs: slack.UserGroupPrefs{Channels: []string{"C2", "C1"}},
	}

	// act
	model.UpdateFromUserGroup(&group)

	// assert
	expected := stringSliceToList([]string{"C2", "C1"})
	if !model.Channels.Equal(expected) {
		t.Errorf("Expected %v, got: %v", expected, model.Channels)
	}
}
