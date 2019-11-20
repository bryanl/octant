// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/vmware-tanzu/octant/internal/octant (interfaces: StateClient)

// Package fake is a generated GoMock package.
package fake

import (
	gomock "github.com/golang/mock/gomock"
	octant "github.com/vmware-tanzu/octant/internal/octant"
	reflect "reflect"
)

// MockStateClient is a mock of StateClient interface
type MockStateClient struct {
	ctrl     *gomock.Controller
	recorder *MockStateClientMockRecorder
}

// MockStateClientMockRecorder is the mock recorder for MockStateClient
type MockStateClientMockRecorder struct {
	mock *MockStateClient
}

// NewMockStateClient creates a new mock instance
func NewMockStateClient(ctrl *gomock.Controller) *MockStateClient {
	mock := &MockStateClient{ctrl: ctrl}
	mock.recorder = &MockStateClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStateClient) EXPECT() *MockStateClientMockRecorder {
	return m.recorder
}

// ID mocks base method
func (m *MockStateClient) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockStateClientMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockStateClient)(nil).ID))
}

// Send mocks base method
func (m *MockStateClient) Send(arg0 octant.Event) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Send", arg0)
}

// Send indicates an expected call of Send
func (mr *MockStateClientMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockStateClient)(nil).Send), arg0)
}
