package tb

import (
	"testing"

	"github.com/essent/terraform-provider-slack/internal/tb/mock_slackExt"
	"go.uber.org/mock/gomock"
)

var global dependenciesImpl = dependenciesImpl{}

func Init(t *testing.T) {
	global.c = gomock.NewController(t)
	global.mock_slack_client = nil
	global.mock_slack_queries = nil
}

func Finish() {
	if global.c != nil {
		global.c.Finish()
	}
}

func MockSlackClient() *mock_slackExt.MockClient {
	global.CreateSlackClient("<TOKEN>")
	return global.mock_slack_client
}

func MockSlackQueries() *mock_slackExt.MockQueries {
	global.CreateSlackQueries(global.mock_slack_client)
	return global.mock_slack_queries
}
