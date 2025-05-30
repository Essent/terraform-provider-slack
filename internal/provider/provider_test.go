// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/essent/terraform-provider-slack/internal/tb"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/slack-go/slack"
	"go.uber.org/mock/gomock"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"slack": providerserver.NewProtocol6WithError(New("test", tb.NewDependencies())()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	tb.Init(t)
}

func testAccPreCheckWithSlackAuth(t *testing.T) {
	testAccPreCheck(t)

	m := tb.MockSlackClient()
	m.EXPECT().AuthTest(gomock.Any()).Return(&slack.AuthTestResponse{}, nil).AnyTimes()
}

func testConfig(t *testing.T, step resource.TestStep) {
	defer tb.Finish()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithSlackAuth(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{step},
	})
}
