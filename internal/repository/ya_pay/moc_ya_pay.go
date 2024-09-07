// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/ya_pay/ya_pay.go

// Package ya_pay is a generated GoMock package.
package ya_pay

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockYaPayRepository is a mock of YaPayRepository interface.
type MockYaPayRepository struct {
	ctrl     *gomock.Controller
	recorder *MockYaPayRepositoryMockRecorder
}

// MockYaPayRepositoryMockRecorder is the mock recorder for MockYaPayRepository.
type MockYaPayRepositoryMockRecorder struct {
	mock *MockYaPayRepository
}

// NewMockYaPayRepository creates a new mock instance.
func NewMockYaPayRepository(ctrl *gomock.Controller) *MockYaPayRepository {
	mock := &MockYaPayRepository{ctrl: ctrl}
	mock.recorder = &MockYaPayRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockYaPayRepository) EXPECT() *MockYaPayRepositoryMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockYaPayRepository) Add(ctx context.Context, idAzs int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, idAzs)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockYaPayRepositoryMockRecorder) Add(ctx, idAzs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockYaPayRepository)(nil).Add), ctx, idAzs)
}

// CreateTable mocks base method.
func (m *MockYaPayRepository) CreateTable(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTable", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTable indicates an expected call of CreateTable.
func (mr *MockYaPayRepositoryMockRecorder) CreateTable(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTable", reflect.TypeOf((*MockYaPayRepository)(nil).CreateTable), ctx)
}

// Delete mocks base method.
func (m *MockYaPayRepository) Delete(ctx context.Context, idAzs int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, idAzs)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockYaPayRepositoryMockRecorder) Delete(ctx, idAzs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockYaPayRepository)(nil).Delete), ctx, idAzs)
}

// DeleteTable mocks base method.
func (m *MockYaPayRepository) DeleteTable(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTable", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTable indicates an expected call of DeleteTable.
func (mr *MockYaPayRepositoryMockRecorder) DeleteTable(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTable", reflect.TypeOf((*MockYaPayRepository)(nil).DeleteTable), ctx)
}

// Get mocks base method.
func (m *MockYaPayRepository) Get(ctx context.Context, idAzs int) (YaPay, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, idAzs)
	ret0, _ := ret[0].(YaPay)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockYaPayRepositoryMockRecorder) Get(ctx, idAzs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockYaPayRepository)(nil).Get), ctx, idAzs)
}

// GetAll mocks base method.
func (m *MockYaPayRepository) GetAll(ctx context.Context) ([]YaPay, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]YaPay)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockYaPayRepositoryMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockYaPayRepository)(nil).GetAll), ctx)
}

// Update mocks base method.
func (m *MockYaPayRepository) Update(ctx context.Context, idAzs, value int, data string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, idAzs, value, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockYaPayRepositoryMockRecorder) Update(ctx, idAzs, value, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockYaPayRepository)(nil).Update), ctx, idAzs, value, data)
}
