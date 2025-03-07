// Code generated by MockGen. DO NOT EDIT.
// Source: ./route.go

// Package mock_security is a generated GoMock package.
package mock_security

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	osc "github.com/outscale/osc-sdk-go/v2"
)

// MockOscRouteTableInterface is a mock of OscRouteTableInterface interface.
type MockOscRouteTableInterface struct {
	ctrl     *gomock.Controller
	recorder *MockOscRouteTableInterfaceMockRecorder
}

// MockOscRouteTableInterfaceMockRecorder is the mock recorder for MockOscRouteTableInterface.
type MockOscRouteTableInterfaceMockRecorder struct {
	mock *MockOscRouteTableInterface
}

// NewMockOscRouteTableInterface creates a new mock instance.
func NewMockOscRouteTableInterface(ctrl *gomock.Controller) *MockOscRouteTableInterface {
	mock := &MockOscRouteTableInterface{ctrl: ctrl}
	mock.recorder = &MockOscRouteTableInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOscRouteTableInterface) EXPECT() *MockOscRouteTableInterfaceMockRecorder {
	return m.recorder
}

// CreateRoute mocks base method.
func (m *MockOscRouteTableInterface) CreateRoute(ctx context.Context, destinationIpRange, routeTableId, resourceId, resourceType string) (*osc.RouteTable, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRoute", ctx, destinationIpRange, routeTableId, resourceId, resourceType)
	ret0, _ := ret[0].(*osc.RouteTable)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRoute indicates an expected call of CreateRoute.
func (mr *MockOscRouteTableInterfaceMockRecorder) CreateRoute(ctx, destinationIpRange, routeTableId, resourceId, resourceType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRoute", reflect.TypeOf((*MockOscRouteTableInterface)(nil).CreateRoute), ctx, destinationIpRange, routeTableId, resourceId, resourceType)
}

// CreateRouteTable mocks base method.
func (m *MockOscRouteTableInterface) CreateRouteTable(ctx context.Context, netId, clusterName, routeTableName string) (*osc.RouteTable, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRouteTable", ctx, netId, clusterName, routeTableName)
	ret0, _ := ret[0].(*osc.RouteTable)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRouteTable indicates an expected call of CreateRouteTable.
func (mr *MockOscRouteTableInterfaceMockRecorder) CreateRouteTable(ctx, netId, clusterName, routeTableName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRouteTable", reflect.TypeOf((*MockOscRouteTableInterface)(nil).CreateRouteTable), ctx, netId, clusterName, routeTableName)
}

// DeleteRoute mocks base method.
func (m *MockOscRouteTableInterface) DeleteRoute(ctx context.Context, destinationIpRange, routeTableId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRoute", ctx, destinationIpRange, routeTableId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRoute indicates an expected call of DeleteRoute.
func (mr *MockOscRouteTableInterfaceMockRecorder) DeleteRoute(ctx, destinationIpRange, routeTableId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRoute", reflect.TypeOf((*MockOscRouteTableInterface)(nil).DeleteRoute), ctx, destinationIpRange, routeTableId)
}

// DeleteRouteTable mocks base method.
func (m *MockOscRouteTableInterface) DeleteRouteTable(ctx context.Context, routeTableId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRouteTable", ctx, routeTableId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRouteTable indicates an expected call of DeleteRouteTable.
func (mr *MockOscRouteTableInterfaceMockRecorder) DeleteRouteTable(ctx, routeTableId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRouteTable", reflect.TypeOf((*MockOscRouteTableInterface)(nil).DeleteRouteTable), ctx, routeTableId)
}

// GetRouteTable mocks base method.
func (m *MockOscRouteTableInterface) GetRouteTable(ctx context.Context, routeTableId []string) (*osc.RouteTable, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRouteTable", ctx, routeTableId)
	ret0, _ := ret[0].(*osc.RouteTable)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRouteTable indicates an expected call of GetRouteTable.
func (mr *MockOscRouteTableInterfaceMockRecorder) GetRouteTable(ctx, routeTableId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRouteTable", reflect.TypeOf((*MockOscRouteTableInterface)(nil).GetRouteTable), ctx, routeTableId)
}

// GetRouteTableFromRoute mocks base method.
func (m *MockOscRouteTableInterface) GetRouteTableFromRoute(ctx context.Context, routeTableId, resourceId, resourceType string) (*osc.RouteTable, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRouteTableFromRoute", ctx, routeTableId, resourceId, resourceType)
	ret0, _ := ret[0].(*osc.RouteTable)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRouteTableFromRoute indicates an expected call of GetRouteTableFromRoute.
func (mr *MockOscRouteTableInterfaceMockRecorder) GetRouteTableFromRoute(ctx, routeTableId, resourceId, resourceType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRouteTableFromRoute", reflect.TypeOf((*MockOscRouteTableInterface)(nil).GetRouteTableFromRoute), ctx, routeTableId, resourceId, resourceType)
}

// GetRouteTableIdsFromNetIds mocks base method.
func (m *MockOscRouteTableInterface) GetRouteTableIdsFromNetIds(ctx context.Context, netId string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRouteTableIdsFromNetIds", ctx, netId)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRouteTableIdsFromNetIds indicates an expected call of GetRouteTableIdsFromNetIds.
func (mr *MockOscRouteTableInterfaceMockRecorder) GetRouteTableIdsFromNetIds(ctx, netId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRouteTableIdsFromNetIds", reflect.TypeOf((*MockOscRouteTableInterface)(nil).GetRouteTableIdsFromNetIds), ctx, netId)
}

// LinkRouteTable mocks base method.
func (m *MockOscRouteTableInterface) LinkRouteTable(ctx context.Context, routeTableId, subnetId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LinkRouteTable", ctx, routeTableId, subnetId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LinkRouteTable indicates an expected call of LinkRouteTable.
func (mr *MockOscRouteTableInterfaceMockRecorder) LinkRouteTable(ctx, routeTableId, subnetId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LinkRouteTable", reflect.TypeOf((*MockOscRouteTableInterface)(nil).LinkRouteTable), ctx, routeTableId, subnetId)
}

// UnlinkRouteTable mocks base method.
func (m *MockOscRouteTableInterface) UnlinkRouteTable(ctx context.Context, linkRouteTableId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlinkRouteTable", ctx, linkRouteTableId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnlinkRouteTable indicates an expected call of UnlinkRouteTable.
func (mr *MockOscRouteTableInterfaceMockRecorder) UnlinkRouteTable(ctx, linkRouteTableId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlinkRouteTable", reflect.TypeOf((*MockOscRouteTableInterface)(nil).UnlinkRouteTable), ctx, linkRouteTableId)
}
