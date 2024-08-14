// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repositories/user_repo.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "gitlab.com/jkozhemiaka/web-layout/internal/models"
)

// MockUserRepoInterface is a mock of UserRepoInterface interface.
type MockUserRepoInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoInterfaceMockRecorder
}

// MockUserRepoInterfaceMockRecorder is the mock recorder for MockUserRepoInterface.
type MockUserRepoInterfaceMockRecorder struct {
	mock *MockUserRepoInterface
}

// NewMockUserRepoInterface creates a new mock instance.
func NewMockUserRepoInterface(ctrl *gomock.Controller) *MockUserRepoInterface {
	mock := &MockUserRepoInterface{ctrl: ctrl}
	mock.recorder = &MockUserRepoInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepoInterface) EXPECT() *MockUserRepoInterfaceMockRecorder {
	return m.recorder
}

// CountUsers mocks base method.
func (m *MockUserRepoInterface) CountUsers(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountUsers", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountUsers indicates an expected call of CountUsers.
func (mr *MockUserRepoInterfaceMockRecorder) CountUsers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountUsers", reflect.TypeOf((*MockUserRepoInterface)(nil).CountUsers), ctx)
}

// GetVote mocks base method.
func (m *MockUserRepoInterface) GetVote(ctx context.Context, userID uint, profileID uint) (*models.Vote, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVote", ctx, userID, profileID)
	ret0, _ := ret[0].(*models.Vote)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVote indicates an expected call of GetVote.
func (mr *MockUserRepoInterfaceMockRecorder) GetVote(ctx, userID, profileID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVote", reflect.TypeOf((*MockUserRepoInterface)(nil).GetVote), ctx, userID, profileID)
}

// CreateVote mocks base method.
func (m *MockUserRepoInterface) CreateVote(ctx context.Context, vote *models.Vote) (*models.Vote, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVote", ctx, vote)
	ret0, _ := ret[0].(*models.Vote)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateVote indicates an expected call of CreateVote.
func (mr *MockUserRepoInterfaceMockRecorder) CreateVote(ctx, vote interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVote", reflect.TypeOf((*MockUserRepoInterface)(nil).CreateVote), ctx, vote)
}

// UpdateVote mocks base method.
func (m *MockUserRepoInterface) UpdateVote(ctx context.Context, vote *models.Vote) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVote", ctx, vote)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateVote indicates an expected call of UpdateVote.
func (mr *MockUserRepoInterfaceMockRecorder) UpdateVote(ctx, vote interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVote", reflect.TypeOf((*MockUserRepoInterface)(nil).UpdateVote), ctx, vote)
}

// DeleteVote mocks base method
func (m *MockUserRepoInterface) DeleteVote(ctx context.Context, userID uint, profileID uint) error {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "DeleteVote", ctx, userID, profileID)
    ret0, _ := ret[0].(error)
    return ret0
}

// DeleteVote indicates an expected call of DeleteVote
func (mr *MockUserRepoInterfaceMockRecorder) DeleteVote(ctx, userID, profileID interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVote", reflect.TypeOf((*MockUserRepoInterface)(nil).DeleteVote), ctx, userID, profileID)
}


// GetUserByID mocks base method
func (m *MockUserRepoInterface) GetUserByID(ctx context.Context, userID uint, user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, userID, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetUserByID indicates an expected call of GetUserByID
func (mr *MockUserRepoInterfaceMockRecorder) GetUserByID(ctx, userID, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserRepoInterface)(nil).GetUserByID), ctx, userID, user)
}
// CreateUser mocks base method.
func (m *MockUserRepoInterface) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepoInterfaceMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepoInterface)(nil).CreateUser), ctx, user)
}

// DeleteUser mocks base method.
func (m *MockUserRepoInterface) DeleteUser(ctx context.Context, userID string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, userID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserRepoInterfaceMockRecorder) DeleteUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserRepoInterface)(nil).DeleteUser), ctx, userID)
}

// GetUser mocks base method.
func (m *MockUserRepoInterface) GetUser(ctx context.Context, userID string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, userID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserRepoInterfaceMockRecorder) GetUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserRepoInterface)(nil).GetUser), ctx, userID)
}

// GetUserByEmail mocks base method.
func (m *MockUserRepoInterface) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserRepoInterfaceMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserRepoInterface)(nil).GetUserByEmail), ctx, email)
}

// ListUsers mocks base method.
func (m *MockUserRepoInterface) ListUsers(ctx context.Context, page, pageSize int) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", ctx, page, pageSize)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockUserRepoInterfaceMockRecorder) ListUsers(ctx, page, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockUserRepoInterface)(nil).ListUsers), ctx, page, pageSize)
}

// UpdateUser mocks base method.
func (m *MockUserRepoInterface) UpdateUser(ctx context.Context, userID string, updatedData *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, userID, updatedData)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserRepoInterfaceMockRecorder) UpdateUser(ctx, userID, updatedData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepoInterface)(nil).UpdateUser), ctx, userID, updatedData)
}
