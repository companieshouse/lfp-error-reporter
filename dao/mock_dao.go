// Code generated by MockGen. DO NOT EDIT.
// Source: dao/dao.go

// Package mock_dao is a generated GoMock package.
package dao

import (
	reflect "reflect"

	models "github.com/companieshouse/lfp-error-reporter/models"
	gomock "github.com/golang/mock/gomock"
)

// MockDAO is a mock of DAO interface
type MockDAO struct {
	ctrl     *gomock.Controller
	recorder *MockDAOMockRecorder
}

// MockDAOMockRecorder is the mock recorder for MockDAO
type MockDAOMockRecorder struct {
	mock *MockDAO
}

// NewMockDAO creates a new mock instance
func NewMockDAO(ctrl *gomock.Controller) *MockDAO {
	mock := &MockDAO{ctrl: ctrl}
	mock.recorder = &MockDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDAO) EXPECT() *MockDAOMockRecorder {
	return m.recorder
}

// GetLFPData mocks base method
func (m *MockDAO) GetLFPData(reconciliationMetaData *models.ReconciliationMetaData) (models.PenaltyList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLFPData", reconciliationMetaData)
	ret0, _ := ret[0].(models.PenaltyList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLFPData indicates an expected call of GetLFPData
func (mr *MockDAOMockRecorder) GetLFPData(reconciliationMetaData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLFPData", reflect.TypeOf((*MockDAO)(nil).GetLFPData), reconciliationMetaData)
}
