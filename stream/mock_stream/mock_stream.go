// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rusq/slackdump/v3/stream (interfaces: Slacker)
//
// Generated by this command:
//
//	mockgen -destination mock_stream/mock_stream.go . Slacker
//

// Package mock_stream is a generated GoMock package.
package mock_stream

import (
	context "context"
	reflect "reflect"

	slack "github.com/rusq/slack"
	gomock "go.uber.org/mock/gomock"
)

// MockSlacker is a mock of Slacker interface.
type MockSlacker struct {
	ctrl     *gomock.Controller
	recorder *MockSlackerMockRecorder
	isgomock struct{}
}

// MockSlackerMockRecorder is the mock recorder for MockSlacker.
type MockSlackerMockRecorder struct {
	mock *MockSlacker
}

// NewMockSlacker creates a new mock instance.
func NewMockSlacker(ctrl *gomock.Controller) *MockSlacker {
	mock := &MockSlacker{ctrl: ctrl}
	mock.recorder = &MockSlackerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSlacker) EXPECT() *MockSlackerMockRecorder {
	return m.recorder
}

// AuthTestContext mocks base method.
func (m *MockSlacker) AuthTestContext(arg0 context.Context) (*slack.AuthTestResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthTestContext", arg0)
	ret0, _ := ret[0].(*slack.AuthTestResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthTestContext indicates an expected call of AuthTestContext.
func (mr *MockSlackerMockRecorder) AuthTestContext(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthTestContext", reflect.TypeOf((*MockSlacker)(nil).AuthTestContext), arg0)
}

// GetConversationHistoryContext mocks base method.
func (m *MockSlacker) GetConversationHistoryContext(ctx context.Context, params *slack.GetConversationHistoryParameters) (*slack.GetConversationHistoryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConversationHistoryContext", ctx, params)
	ret0, _ := ret[0].(*slack.GetConversationHistoryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConversationHistoryContext indicates an expected call of GetConversationHistoryContext.
func (mr *MockSlackerMockRecorder) GetConversationHistoryContext(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConversationHistoryContext", reflect.TypeOf((*MockSlacker)(nil).GetConversationHistoryContext), ctx, params)
}

// GetConversationInfoContext mocks base method.
func (m *MockSlacker) GetConversationInfoContext(ctx context.Context, input *slack.GetConversationInfoInput) (*slack.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConversationInfoContext", ctx, input)
	ret0, _ := ret[0].(*slack.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConversationInfoContext indicates an expected call of GetConversationInfoContext.
func (mr *MockSlackerMockRecorder) GetConversationInfoContext(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConversationInfoContext", reflect.TypeOf((*MockSlacker)(nil).GetConversationInfoContext), ctx, input)
}

// GetConversationRepliesContext mocks base method.
func (m *MockSlacker) GetConversationRepliesContext(ctx context.Context, params *slack.GetConversationRepliesParameters) ([]slack.Message, bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConversationRepliesContext", ctx, params)
	ret0, _ := ret[0].([]slack.Message)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetConversationRepliesContext indicates an expected call of GetConversationRepliesContext.
func (mr *MockSlackerMockRecorder) GetConversationRepliesContext(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConversationRepliesContext", reflect.TypeOf((*MockSlacker)(nil).GetConversationRepliesContext), ctx, params)
}

// GetConversationsContext mocks base method.
func (m *MockSlacker) GetConversationsContext(ctx context.Context, params *slack.GetConversationsParameters) ([]slack.Channel, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConversationsContext", ctx, params)
	ret0, _ := ret[0].([]slack.Channel)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetConversationsContext indicates an expected call of GetConversationsContext.
func (mr *MockSlackerMockRecorder) GetConversationsContext(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConversationsContext", reflect.TypeOf((*MockSlacker)(nil).GetConversationsContext), ctx, params)
}

// GetFileInfoContext mocks base method.
func (m *MockSlacker) GetFileInfoContext(ctx context.Context, fileID string, count, page int) (*slack.File, []slack.Comment, *slack.Paging, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileInfoContext", ctx, fileID, count, page)
	ret0, _ := ret[0].(*slack.File)
	ret1, _ := ret[1].([]slack.Comment)
	ret2, _ := ret[2].(*slack.Paging)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetFileInfoContext indicates an expected call of GetFileInfoContext.
func (mr *MockSlackerMockRecorder) GetFileInfoContext(ctx, fileID, count, page any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileInfoContext", reflect.TypeOf((*MockSlacker)(nil).GetFileInfoContext), ctx, fileID, count, page)
}

// GetStarredContext mocks base method.
func (m *MockSlacker) GetStarredContext(ctx context.Context, params slack.StarsParameters) ([]slack.StarredItem, *slack.Paging, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStarredContext", ctx, params)
	ret0, _ := ret[0].([]slack.StarredItem)
	ret1, _ := ret[1].(*slack.Paging)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetStarredContext indicates an expected call of GetStarredContext.
func (mr *MockSlackerMockRecorder) GetStarredContext(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStarredContext", reflect.TypeOf((*MockSlacker)(nil).GetStarredContext), ctx, params)
}

// GetUserInfoContext mocks base method.
func (m *MockSlacker) GetUserInfoContext(ctx context.Context, user string) (*slack.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfoContext", ctx, user)
	ret0, _ := ret[0].(*slack.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInfoContext indicates an expected call of GetUserInfoContext.
func (mr *MockSlackerMockRecorder) GetUserInfoContext(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfoContext", reflect.TypeOf((*MockSlacker)(nil).GetUserInfoContext), ctx, user)
}

// GetUsersInConversationContext mocks base method.
func (m *MockSlacker) GetUsersInConversationContext(ctx context.Context, params *slack.GetUsersInConversationParameters) ([]string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersInConversationContext", ctx, params)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUsersInConversationContext indicates an expected call of GetUsersInConversationContext.
func (mr *MockSlackerMockRecorder) GetUsersInConversationContext(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersInConversationContext", reflect.TypeOf((*MockSlacker)(nil).GetUsersInConversationContext), ctx, params)
}

// GetUsersPaginated mocks base method.
func (m *MockSlacker) GetUsersPaginated(options ...slack.GetUsersOption) slack.UserPagination {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUsersPaginated", varargs...)
	ret0, _ := ret[0].(slack.UserPagination)
	return ret0
}

// GetUsersPaginated indicates an expected call of GetUsersPaginated.
func (mr *MockSlackerMockRecorder) GetUsersPaginated(options ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersPaginated", reflect.TypeOf((*MockSlacker)(nil).GetUsersPaginated), options...)
}

// ListBookmarks mocks base method.
func (m *MockSlacker) ListBookmarks(channelID string) ([]slack.Bookmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBookmarks", channelID)
	ret0, _ := ret[0].([]slack.Bookmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBookmarks indicates an expected call of ListBookmarks.
func (mr *MockSlackerMockRecorder) ListBookmarks(channelID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBookmarks", reflect.TypeOf((*MockSlacker)(nil).ListBookmarks), channelID)
}

// SearchFilesContext mocks base method.
func (m *MockSlacker) SearchFilesContext(ctx context.Context, query string, params slack.SearchParameters) (*slack.SearchFiles, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchFilesContext", ctx, query, params)
	ret0, _ := ret[0].(*slack.SearchFiles)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchFilesContext indicates an expected call of SearchFilesContext.
func (mr *MockSlackerMockRecorder) SearchFilesContext(ctx, query, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchFilesContext", reflect.TypeOf((*MockSlacker)(nil).SearchFilesContext), ctx, query, params)
}

// SearchMessagesContext mocks base method.
func (m *MockSlacker) SearchMessagesContext(ctx context.Context, query string, params slack.SearchParameters) (*slack.SearchMessages, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchMessagesContext", ctx, query, params)
	ret0, _ := ret[0].(*slack.SearchMessages)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchMessagesContext indicates an expected call of SearchMessagesContext.
func (mr *MockSlackerMockRecorder) SearchMessagesContext(ctx, query, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchMessagesContext", reflect.TypeOf((*MockSlacker)(nil).SearchMessagesContext), ctx, query, params)
}
