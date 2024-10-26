// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/protomem/chatik/internal/core/service (interfaces: Message)
//
// Generated by this command:
//
//	mockgen -destination mocks/message.go -package mocks . Message
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/protomem/chatik/internal/core/entity"
	service "github.com/protomem/chatik/internal/core/service"
	gomock "go.uber.org/mock/gomock"
)

// MockMessage is a mock of Message interface.
type MockMessage struct {
	ctrl     *gomock.Controller
	recorder *MockMessageMockRecorder
}

// MockMessageMockRecorder is the mock recorder for MockMessage.
type MockMessageMockRecorder struct {
	mock *MockMessage
}

// NewMockMessage creates a new mock instance.
func NewMockMessage(ctrl *gomock.Controller) *MockMessage {
	mock := &MockMessage{ctrl: ctrl}
	mock.recorder = &MockMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessage) EXPECT() *MockMessageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockMessage) Create(arg0 context.Context, arg1 service.CreateMessageDTO) (entity.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(entity.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockMessageMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMessage)(nil).Create), arg0, arg1)
}

// FindByFromAndTo mocks base method.
func (m *MockMessage) FindByFromAndTo(arg0 context.Context, arg1 service.FindMessageByFromAndTo) ([]entity.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByFromAndTo", arg0, arg1)
	ret0, _ := ret[0].([]entity.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByFromAndTo indicates an expected call of FindByFromAndTo.
func (mr *MockMessageMockRecorder) FindByFromAndTo(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByFromAndTo", reflect.TypeOf((*MockMessage)(nil).FindByFromAndTo), arg0, arg1)
}
