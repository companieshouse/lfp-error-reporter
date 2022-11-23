// Code generated by MockGen. DO NOT EDIT.
// Source: service/service.go

// Package mock_service is a generated GoMock package.
package service

import (
	reflect "reflect"

	models "github.com/companieshouse/lfp-error-reporter/models"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetLFPCSV mocks base method
func (m *MockService) GetLFPCSV(reconciliationMetaData *models.ReconciliationMetaData) (models.CSV, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLFPCSV", reconciliationMetaData)
	ret0, _ := ret[0].(models.CSV)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLFPCSV indicates an expected call of GetLFPCSV
func (mr *MockServiceMockRecorder) GetLFPCSV(reconciliationMetaData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLFPCSV", reflect.TypeOf((*MockService)(nil).GetLFPCSV), reconciliationMetaData)
}