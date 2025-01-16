// Code generated by MockGen. DO NOT EDIT.
// Source: ../slackExt/client.go

// Package mock_slackExt is a generated GoMock package.
package mock_slackExt

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	slack "github.com/slack-go/slack"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// AuthTest mocks base method.
func (m *MockClient) AuthTest(ctx context.Context) (*slack.AuthTestResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthTest", ctx)
	ret0, _ := ret[0].(*slack.AuthTestResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthTest indicates an expected call of AuthTest.
func (mr *MockClientMockRecorder) AuthTest(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthTest", reflect.TypeOf((*MockClient)(nil).AuthTest), ctx)
}

// GetConversationInfo mocks base method.
func (m *MockClient) GetConversationInfo(ctx context.Context, input *slack.GetConversationInfoInput) (*slack.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConversationInfo", ctx, input)
	ret0, _ := ret[0].(*slack.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConversationInfo indicates an expected call of GetConversationInfo.
func (mr *MockClientMockRecorder) GetConversationInfo(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConversationInfo", reflect.TypeOf((*MockClient)(nil).GetConversationInfo), ctx, input)
}

// GetUserByEmail mocks base method.
func (m *MockClient) GetUserByEmail(ctx context.Context, email string) (*slack.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(*slack.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockClientMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockClient)(nil).GetUserByEmail), ctx, email)
}

// GetUserGroups mocks base method.
func (m *MockClient) GetUserGroups(ctx context.Context, options ...slack.GetUserGroupsOption) ([]slack.UserGroup, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUserGroups", varargs...)
	ret0, _ := ret[0].([]slack.UserGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserGroups indicates an expected call of GetUserGroups.
func (mr *MockClientMockRecorder) GetUserGroups(ctx interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserGroups", reflect.TypeOf((*MockClient)(nil).GetUserGroups), varargs...)
}

// GetUserInfo mocks base method.
func (m *MockClient) GetUserInfo(ctx context.Context, user string) (*slack.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfo", ctx, user)
	ret0, _ := ret[0].(*slack.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInfo indicates an expected call of GetUserInfo.
func (mr *MockClientMockRecorder) GetUserInfo(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfo", reflect.TypeOf((*MockClient)(nil).GetUserInfo), ctx, user)
}

// GetUsersContext mocks base method.
func (m *MockClient) GetUsersContext(ctx context.Context) ([]slack.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersContext", ctx)
	ret0, _ := ret[0].([]slack.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersContext indicates an expected call of GetUsersContext.
func (mr *MockClientMockRecorder) GetUsersContext(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersContext", reflect.TypeOf((*MockClient)(nil).GetUsersContext), ctx)
}
