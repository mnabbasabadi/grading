// Code generated by MockGen. DO NOT EDIT.
// Source: stores.go

// Package postgres is a generated GoMock package.
package postgres

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/mnabbasabadi/grading/service/shared/domain"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetGrades mocks base method.
func (m *MockRepository) GetGrades(arg0 context.Context, arg1, arg2 int) ([]domain.Grade, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGrades", arg0, arg1, arg2)
	ret0, _ := ret[0].([]domain.Grade)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetGrades indicates an expected call of GetGrades.
func (mr *MockRepositoryMockRecorder) GetGrades(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGrades", reflect.TypeOf((*MockRepository)(nil).GetGrades), arg0, arg1, arg2)
}

// GetScales mocks base method.
func (m *MockRepository) GetScales(arg0 context.Context, arg1 domain.ScaleType) (domain.Scales, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScales", arg0, arg1)
	ret0, _ := ret[0].(domain.Scales)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScales indicates an expected call of GetScales.
func (mr *MockRepositoryMockRecorder) GetScales(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScales", reflect.TypeOf((*MockRepository)(nil).GetScales), arg0, arg1)
}
