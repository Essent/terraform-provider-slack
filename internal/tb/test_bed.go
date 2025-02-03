// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tb

import (
	"fmt"
	"testing"

	"github.com/essent/terraform-provider-slack/internal/provider/dependencies"
	"github.com/essent/terraform-provider-slack/internal/slackExt"
	"github.com/essent/terraform-provider-slack/internal/tb/mock_slackExt"
	"go.uber.org/mock/gomock"
)

type dependenciesImpl struct {
	c *gomock.Controller

	mockSlackClient *mock_slackExt.MockClient
}

var global dependenciesImpl = dependenciesImpl{}

func Init(t *testing.T) {
	global.c = gomock.NewController(t)
	global.mockSlackClient = nil
}

func Finish() {
	if global.c != nil {
		global.c.Finish()
	}
}

func MockSlackClient() *mock_slackExt.MockClient {
	global.CreateSlackClient("<TOKEN>")
	return global.mockSlackClient
}

func (d *dependenciesImpl) useMockController() *gomock.Controller {
	if d.c == nil {
		panic("Init() must be called before useMockController()")
	}

	return d.c
}

func (d *dependenciesImpl) CreateSlackClient(token string) slackExt.Client {
	if d.mockSlackClient != nil {
		return d.mockSlackClient
	}

	d.mockSlackClient = mock_slackExt.NewMockClient(d.useMockController())
	return d.mockSlackClient
}

func NewDependencies() dependencies.Dependencies {
	return &global
}

func ExpectString(value string) func(string) error {
	return func(actual string) error {
		if actual != value {
			return fmt.Errorf("expected %q, got %q", value, actual)
		}
		return nil
	}
}
