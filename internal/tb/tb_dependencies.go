package tb

import (
	"github.com/essent/terraform-provider-slack/internal/provider/dependencies"
	"github.com/essent/terraform-provider-slack/internal/slackExt"
	"github.com/essent/terraform-provider-slack/internal/tb/mock_slackExt"
	"go.uber.org/mock/gomock"
)

type dependenciesImpl struct {
	c *gomock.Controller

	mockSlackClient *mock_slackExt.MockClient
}

func (d *dependenciesImpl) CreateSlackClient(token string) slackExt.Client {
	if d.mockSlackClient != nil {
		return d.mockSlackClient
	}

	d.mockSlackClient = mock_slackExt.NewMockClient(d.useMockController())
	return d.mockSlackClient
}

func (d *dependenciesImpl) useMockController() *gomock.Controller {
	if d.c == nil {
		panic("Mock controller not set")
	}

	return d.c
}

func NewDependencies() dependencies.Dependencies {
	return &global
}
