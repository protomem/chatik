// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/protomem/chatik/internal/core/data (interfaces: UserAccessor,UserMutator)
//
// Generated by this command:
//
//	mockgen -destination mocks/user.go -package mocks . UserAccessor,UserMutator
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	data "github.com/protomem/chatik/internal/core/data"
	entity "github.com/protomem/chatik/internal/core/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockUserAccessor is a mock of UserAccessor interface.
type MockUserAccessor struct {
	ctrl     *gomock.Controller
	recorder *MockUserAccessorMockRecorder
}

// MockUserAccessorMockRecorder is the mock recorder for MockUserAccessor.
type MockUserAccessorMockRecorder struct {
	mock *MockUserAccessor
}

// NewMockUserAccessor creates a new mock instance.
func NewMockUserAccessor(ctrl *gomock.Controller) *MockUserAccessor {
	mock := &MockUserAccessor{ctrl: ctrl}
	mock.recorder = &MockUserAccessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserAccessor) EXPECT() *MockUserAccessorMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockUserAccessor) Get(arg0 context.Context, arg1 uuid.UUID) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserAccessorMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserAccessor)(nil).Get), arg0, arg1)
}

// GetByNickname mocks base method.
func (m *MockUserAccessor) GetByNickname(arg0 context.Context, arg1 string) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNickname", arg0, arg1)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNickname indicates an expected call of GetByNickname.
func (mr *MockUserAccessorMockRecorder) GetByNickname(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNickname", reflect.TypeOf((*MockUserAccessor)(nil).GetByNickname), arg0, arg1)
}

// Select mocks base method.
func (m *MockUserAccessor) Select(arg0 context.Context) ([]entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Select", arg0)
	ret0, _ := ret[0].([]entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Select indicates an expected call of Select.
func (mr *MockUserAccessorMockRecorder) Select(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Select", reflect.TypeOf((*MockUserAccessor)(nil).Select), arg0)
}

// MockUserMutator is a mock of UserMutator interface.
type MockUserMutator struct {
	ctrl     *gomock.Controller
	recorder *MockUserMutatorMockRecorder
}

// MockUserMutatorMockRecorder is the mock recorder for MockUserMutator.
type MockUserMutatorMockRecorder struct {
	mock *MockUserMutator
}

// NewMockUserMutator creates a new mock instance.
func NewMockUserMutator(ctrl *gomock.Controller) *MockUserMutator {
	mock := &MockUserMutator{ctrl: ctrl}
	mock.recorder = &MockUserMutatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserMutator) EXPECT() *MockUserMutatorMockRecorder {
	return m.recorder
}

// Insert mocks base method.
func (m *MockUserMutator) Insert(arg0 context.Context, arg1 data.InsertUserDTO) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockUserMutatorMockRecorder) Insert(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUserMutator)(nil).Insert), arg0, arg1)
}
