// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package user is a generated GoMock package.
package user

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockUserRepository) Add(ctx context.Context, name, surname, login, hashedPassword string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, name, surname, login, hashedPassword)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockUserRepositoryMockRecorder) Add(ctx, name, surname, login, hashedPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockUserRepository)(nil).Add), ctx, name, surname, login, hashedPassword)
}

// Delete mocks base method.
func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUserRepository)(nil).Delete), ctx, id)
}

// Find mocks base method.
func (m *MockUserRepository) Find(ctx context.Context, login string) (User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, login)
	ret0, _ := ret[0].(User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockUserRepositoryMockRecorder) Find(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockUserRepository)(nil).Find), ctx, login)
}

// Get mocks base method.
func (m *MockUserRepository) Get(ctx context.Context, id int) (User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserRepositoryMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserRepository)(nil).Get), ctx, id)
}

// GetAll mocks base method.
func (m *MockUserRepository) GetAll(ctx context.Context) ([]User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockUserRepositoryMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockUserRepository)(nil).GetAll), ctx)
}

// Login mocks base method.
func (m *MockUserRepository) Login(ctx context.Context, login, hashedPassword string) (User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, login, hashedPassword)
	ret0, _ := ret[0].(User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserRepositoryMockRecorder) Login(ctx, login, hashedPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserRepository)(nil).Login), ctx, login, hashedPassword)
}

// Update mocks base method.
func (m *MockUserRepository) Update(ctx context.Context, user User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserRepositoryMockRecorder) Update(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserRepository)(nil).Update), ctx, user)
}

// UpdateUserPassword mocks base method.
func (m *MockUserRepository) UpdateUserPassword(ctx context.Context, id int, hashedPassword string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", ctx, id, hashedPassword)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPassword indicates an expected call of UpdateUserPassword.
func (mr *MockUserRepositoryMockRecorder) UpdateUserPassword(ctx, id, hashedPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPassword", reflect.TypeOf((*MockUserRepository)(nil).UpdateUserPassword), ctx, id, hashedPassword)
}