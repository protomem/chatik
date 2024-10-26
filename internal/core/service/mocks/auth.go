// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/protomem/chatik/internal/core/service (interfaces: Auth)
//
// Generated by this command:
//
//	mockgen -destination mocks/auth.go -package mocks . Auth
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

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuth) Login(arg0 context.Context, arg1 service.LoginDTO) (entity.User, service.TokensDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(service.TokensDTO)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Login indicates an expected call of Login.
func (mr *MockAuthMockRecorder) Login(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuth)(nil).Login), arg0, arg1)
}

// Logout mocks base method.
func (m *MockAuth) Logout(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockAuthMockRecorder) Logout(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuth)(nil).Logout), arg0, arg1)
}

// Refresh mocks base method.
func (m *MockAuth) Refresh(arg0 context.Context, arg1 service.RefreshDTO) (service.TokensDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", arg0, arg1)
	ret0, _ := ret[0].(service.TokensDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh.
func (mr *MockAuthMockRecorder) Refresh(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockAuth)(nil).Refresh), arg0, arg1)
}

// Register mocks base method.
func (m *MockAuth) Register(arg0 context.Context, arg1 service.RegisterDTO) (entity.User, service.TokensDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(service.TokensDTO)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Register indicates an expected call of Register.
func (mr *MockAuthMockRecorder) Register(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuth)(nil).Register), arg0, arg1)
}

// Verify mocks base method.
func (m *MockAuth) Verify(arg0 context.Context, arg1 string) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", arg0, arg1)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Verify indicates an expected call of Verify.
func (mr *MockAuthMockRecorder) Verify(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockAuth)(nil).Verify), arg0, arg1)
}
