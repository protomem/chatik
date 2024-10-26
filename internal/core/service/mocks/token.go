// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/protomem/chatik/internal/core/service (interfaces: Token)
//
// Generated by this command:
//
//	mockgen -destination mocks/token.go -package mocks . Token
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	service "github.com/protomem/chatik/internal/core/service"
	gomock "go.uber.org/mock/gomock"
)

// MockToken is a mock of Token interface.
type MockToken struct {
	ctrl     *gomock.Controller
	recorder *MockTokenMockRecorder
}

// MockTokenMockRecorder is the mock recorder for MockToken.
type MockTokenMockRecorder struct {
	mock *MockToken
}

// NewMockToken creates a new mock instance.
func NewMockToken(ctrl *gomock.Controller) *MockToken {
	mock := &MockToken{ctrl: ctrl}
	mock.recorder = &MockTokenMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockToken) EXPECT() *MockTokenMockRecorder {
	return m.recorder
}

// GenerateJWT mocks base method.
func (m *MockToken) GenerateJWT(arg0 service.JWTTokenPayload, arg1 service.JWTTokenOptions) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateJWT", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateJWT indicates an expected call of GenerateJWT.
func (mr *MockTokenMockRecorder) GenerateJWT(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateJWT", reflect.TypeOf((*MockToken)(nil).GenerateJWT), arg0, arg1)
}

// GenerateRandom mocks base method.
func (m *MockToken) GenerateRandom(arg0 service.RandomTokenOptions) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRandom", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateRandom indicates an expected call of GenerateRandom.
func (mr *MockTokenMockRecorder) GenerateRandom(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRandom", reflect.TypeOf((*MockToken)(nil).GenerateRandom), arg0)
}

// VerifyJWT mocks base method.
func (m *MockToken) VerifyJWT(arg0 string, arg1 service.JWTTokenOptions) (service.JWTTokenPayload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyJWT", arg0, arg1)
	ret0, _ := ret[0].(service.JWTTokenPayload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyJWT indicates an expected call of VerifyJWT.
func (mr *MockTokenMockRecorder) VerifyJWT(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyJWT", reflect.TypeOf((*MockToken)(nil).VerifyJWT), arg0, arg1)
}
