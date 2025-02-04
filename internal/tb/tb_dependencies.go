package tb

import (
	"github.com/essent/terraform-provider-slack/internal/provider/dependencies"
	"github.com/essent/terraform-provider-slack/internal/slackExt"
	"github.com/essent/terraform-provider-slack/internal/tb/mock_slackExt"
	"go.uber.org/mock/gomock"
)

type dependenciesImpl struct {
	c *gomock.Controller

	mock_slack_client  *mock_slackExt.MockClient
	mock_slack_queries *mock_slackExt.MockQueries
}

func (d *dependenciesImpl) CreateSlackClient(token string) slackExt.Client {
	if d.mock_slack_client != nil {
		return d.mock_slack_client
	}

	d.mock_slack_client = mock_slackExt.NewMockClient(d.useMockController())
	return d.mock_slack_client
}

func (d *dependenciesImpl) CreateSlackQueries(client slackExt.Client) slackExt.Queries {
	if d.mock_slack_queries != nil {
		return d.mock_slack_queries
	}

	d.mock_slack_queries = mock_slackExt.NewMockQueries(d.useMockController())
	return d.mock_slack_queries
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
