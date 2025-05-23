// Code generated by MockGen. DO NOT EDIT.
// Source: hydrate.go
//
// Generated by this command:
//
//	mockgen -destination=hydrate_mock_test.go -package=diag -source hydrate.go sourcer
//

// Package diag is a generated GoMock package.
package diag

import (
	context "context"
	iter "iter"
	reflect "reflect"

	slack "github.com/rusq/slack"
	gomock "go.uber.org/mock/gomock"
)

// Mocksourcer is a mock of sourcer interface.
type Mocksourcer struct {
	ctrl     *gomock.Controller
	recorder *MocksourcerMockRecorder
	isgomock struct{}
}

// MocksourcerMockRecorder is the mock recorder for Mocksourcer.
type MocksourcerMockRecorder struct {
	mock *Mocksourcer
}

// NewMocksourcer creates a new mock instance.
func NewMocksourcer(ctrl *gomock.Controller) *Mocksourcer {
	mock := &Mocksourcer{ctrl: ctrl}
	mock.recorder = &MocksourcerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocksourcer) EXPECT() *MocksourcerMockRecorder {
	return m.recorder
}

// AllMessages mocks base method.
func (m *Mocksourcer) AllMessages(ctx context.Context, channelID string) (iter.Seq2[slack.Message, error], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllMessages", ctx, channelID)
	ret0, _ := ret[0].(iter.Seq2[slack.Message, error])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllMessages indicates an expected call of AllMessages.
func (mr *MocksourcerMockRecorder) AllMessages(ctx, channelID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllMessages", reflect.TypeOf((*Mocksourcer)(nil).AllMessages), ctx, channelID)
}

// AllThreadMessages mocks base method.
func (m *Mocksourcer) AllThreadMessages(ctx context.Context, channelID, threadTimestamp string) (iter.Seq2[slack.Message, error], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllThreadMessages", ctx, channelID, threadTimestamp)
	ret0, _ := ret[0].(iter.Seq2[slack.Message, error])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllThreadMessages indicates an expected call of AllThreadMessages.
func (mr *MocksourcerMockRecorder) AllThreadMessages(ctx, channelID, threadTimestamp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllThreadMessages", reflect.TypeOf((*Mocksourcer)(nil).AllThreadMessages), ctx, channelID, threadTimestamp)
}

// Channels mocks base method.
func (m *Mocksourcer) Channels(ctx context.Context) ([]slack.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Channels", ctx)
	ret0, _ := ret[0].([]slack.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Channels indicates an expected call of Channels.
func (mr *MocksourcerMockRecorder) Channels(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Channels", reflect.TypeOf((*Mocksourcer)(nil).Channels), ctx)
}
