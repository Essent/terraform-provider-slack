// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/essent/terraform-provider-slack/internal/tb"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// func Test_Resource_UserGroup(t *testing.T) {
// 	testConfig(t, tr.TestStep{
// 		// arrange
// 		PreConfig: func() {

// 		},
// 	})
// }

func Test_Resource_UserGroup_Error_WhenSlackClientNil(t *testing.T) {
	// arrange
	res := &resource.ConfigureResponse{}
	req := resource.ConfigureRequest{
		ProviderData: &SlackProviderData{
			Client:  nil,
			Queries: tb.MockSlackQueries(),
		},
	}

	test_instance := UserGroupResource{}

	// act
	test_instance.Configure(context.Background(), req, res)

	// assert
	if res.Diagnostics.Errors()[0].Summary() != "Invalid Provider Data" {
		t.Errorf("Expected error summary to be 'Invalid Provider Data', got: %s", res.Diagnostics.Errors()[0].Summary())
	}
}

func Test_Resource_UserGroup_Error_WhenSlackQueriesNil(t *testing.T) {
	// arrange
	res := &resource.ConfigureResponse{}
	req := resource.ConfigureRequest{
		ProviderData: &SlackProviderData{
			Client:  tb.MockSlackClient(),
			Queries: nil,
		},
	}

	test_instance := UserGroupResource{}

	// act
	test_instance.Configure(context.Background(), req, res)

	// assert
	if res.Diagnostics.Errors()[0].Summary() != "Invalid Provider Data" {
		t.Errorf("Expected error summary to be 'Invalid Provider Data', got: %s", res.Diagnostics.Errors()[0].Summary())
	}
}
