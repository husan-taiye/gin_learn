// Code generated by MockGen. DO NOT EDIT.
// Source: webook/internal/repository/user.go
//
// Generated by this command:
//
//	mockgen -source=webook/internal/repository/user.go -package=repomocks -destination=webook/internal/repository/mocks/user_mock.go
//

// Package repomocks is a generated GoMock package.
package repomocks

import (
	context "context"
	domain "gin_learn/webook/internal/domain"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
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

// Create mocks base method.
func (m *MockUserRepository) Create(ctx context.Context, u domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserRepositoryMockRecorder) Create(ctx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserRepository)(nil).Create), ctx, u)
}

// FindByEmail mocks base method.
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserRepositoryMockRecorder) FindByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindByEmail), ctx, email)
}

// FindByPhone mocks base method.
func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPhone", ctx, phone)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPhone indicates an expected call of FindByPhone.
func (mr *MockUserRepositoryMockRecorder) FindByPhone(ctx, phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPhone", reflect.TypeOf((*MockUserRepository)(nil).FindByPhone), ctx, phone)
}

// FindByWechat mocks base method.
func (m *MockUserRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByWechat", ctx, openID)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByWechat indicates an expected call of FindByWechat.
func (mr *MockUserRepositoryMockRecorder) FindByWechat(ctx, openID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByWechat", reflect.TypeOf((*MockUserRepository)(nil).FindByWechat), ctx, openID)
}

// FindProfileByUserId mocks base method.
func (m *MockUserRepository) FindProfileByUserId(ctx context.Context, userId int64) (domain.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindProfileByUserId", ctx, userId)
	ret0, _ := ret[0].(domain.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindProfileByUserId indicates an expected call of FindProfileByUserId.
func (mr *MockUserRepositoryMockRecorder) FindProfileByUserId(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindProfileByUserId", reflect.TypeOf((*MockUserRepository)(nil).FindProfileByUserId), ctx, userId)
}

// Update mocks base method.
func (m *MockUserRepository) Update(ctx context.Context, up domain.UserProfile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, up)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserRepositoryMockRecorder) Update(ctx, up any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserRepository)(nil).Update), ctx, up)
}
