// Code generated by MockGen. DO NOT EDIT.
// Source: transaction.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	repository "github.com/andreevym/gophermart/internal/repository"
	gomock "github.com/golang/mock/gomock"
)

// MockTransactionRepository is a mock of TransactionRepository interface.
type MockTransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionRepositoryMockRecorder
}

// MockTransactionRepositoryMockRecorder is the mock recorder for MockTransactionRepository.
type MockTransactionRepositoryMockRecorder struct {
	mock *MockTransactionRepository
}

// NewMockTransactionRepository creates a new mock instance.
func NewMockTransactionRepository(ctrl *gomock.Controller) *MockTransactionRepository {
	mock := &MockTransactionRepository{ctrl: ctrl}
	mock.recorder = &MockTransactionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionRepository) EXPECT() *MockTransactionRepositoryMockRecorder {
	return m.recorder
}

// AccrualAmount mocks base method.
func (m *MockTransactionRepository) AccrualAmount(ctx context.Context, userID int64, orderNumber string, accrual float32, orderStatus string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccrualAmount", ctx, userID, orderNumber, accrual, orderStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// AccrualAmount indicates an expected call of AccrualAmount.
func (mr *MockTransactionRepositoryMockRecorder) AccrualAmount(ctx, userID, orderNumber, accrual, orderStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccrualAmount", reflect.TypeOf((*MockTransactionRepository)(nil).AccrualAmount), ctx, userID, orderNumber, accrual, orderStatus)
}

// CreateTransaction mocks base method.
func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, transaction repository.Transaction) (*repository.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", ctx, transaction)
	ret0, _ := ret[0].(*repository.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockTransactionRepositoryMockRecorder) CreateTransaction(ctx, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).CreateTransaction), ctx, transaction)
}

// DeleteTransaction mocks base method.
func (m *MockTransactionRepository) DeleteTransaction(ctx context.Context, transactionID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTransaction", ctx, transactionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTransaction indicates an expected call of DeleteTransaction.
func (mr *MockTransactionRepositoryMockRecorder) DeleteTransaction(ctx, transactionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).DeleteTransaction), ctx, transactionID)
}

// GetTransactionByID mocks base method.
func (m *MockTransactionRepository) GetTransactionByID(ctx context.Context, transactionID int64) (*repository.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionByID", ctx, transactionID)
	ret0, _ := ret[0].(*repository.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionByID indicates an expected call of GetTransactionByID.
func (mr *MockTransactionRepositoryMockRecorder) GetTransactionByID(ctx, transactionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionByID", reflect.TypeOf((*MockTransactionRepository)(nil).GetTransactionByID), ctx, transactionID)
}

// GetTransactionsByUserID mocks base method.
func (m *MockTransactionRepository) GetTransactionsByUserID(ctx context.Context, userID int64) ([]repository.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionsByUserID", ctx, userID)
	ret0, _ := ret[0].([]repository.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionsByUserID indicates an expected call of GetTransactionsByUserID.
func (mr *MockTransactionRepositoryMockRecorder) GetTransactionsByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionsByUserID", reflect.TypeOf((*MockTransactionRepository)(nil).GetTransactionsByUserID), ctx, userID)
}

// GetTransactionsByUserIDAndOperationType mocks base method.
func (m *MockTransactionRepository) GetTransactionsByUserIDAndOperationType(ctx context.Context, userID int64, operationType string) ([]repository.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionsByUserIDAndOperationType", ctx, userID, operationType)
	ret0, _ := ret[0].([]repository.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionsByUserIDAndOperationType indicates an expected call of GetTransactionsByUserIDAndOperationType.
func (mr *MockTransactionRepositoryMockRecorder) GetTransactionsByUserIDAndOperationType(ctx, userID, operationType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionsByUserIDAndOperationType", reflect.TypeOf((*MockTransactionRepository)(nil).GetTransactionsByUserIDAndOperationType), ctx, userID, operationType)
}

// UpdateTransaction mocks base method.
func (m *MockTransactionRepository) UpdateTransaction(ctx context.Context, transaction repository.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransaction", ctx, transaction)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTransaction indicates an expected call of UpdateTransaction.
func (mr *MockTransactionRepositoryMockRecorder) UpdateTransaction(ctx, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).UpdateTransaction), ctx, transaction)
}
