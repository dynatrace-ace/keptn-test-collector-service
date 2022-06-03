// Code generated by MockGen. DO NOT EDIT.
// Source: ../collector/collector.go

// Package eventHandler is a generated GoMock package.
package eventHandler

import (
	reflect "reflect"
	time "time"

	v2 "github.com/cloudevents/sdk-go/v2"
	gomock "github.com/golang/mock/gomock"
)

// MockCollectorIface is a mock of CollectorIface interface.
type MockCollectorIface struct {
	ctrl     *gomock.Controller
	recorder *MockCollectorIfaceMockRecorder
}

// MockCollectorIfaceMockRecorder is the mock recorder for MockCollectorIface.
type MockCollectorIfaceMockRecorder struct {
	mock *MockCollectorIface
}

// NewMockCollectorIface creates a new mock instance.
func NewMockCollectorIface(ctrl *gomock.Controller) *MockCollectorIface {
	mock := &MockCollectorIface{ctrl: ctrl}
	mock.recorder = &MockCollectorIfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCollectorIface) EXPECT() *MockCollectorIfaceMockRecorder {
	return m.recorder
}

// CollectBatchIds mocks base method.
func (m *MockCollectorIface) CollectBatchIds(events []v2.Event) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CollectBatchIds", events)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CollectBatchIds indicates an expected call of CollectBatchIds.
func (mr *MockCollectorIfaceMockRecorder) CollectBatchIds(events interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectBatchIds", reflect.TypeOf((*MockCollectorIface)(nil).CollectBatchIds), events)
}

// CollectEarliestTime mocks base method.
func (m *MockCollectorIface) CollectEarliestTime(events []v2.Event, isFloored bool) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CollectEarliestTime", events, isFloored)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CollectEarliestTime indicates an expected call of CollectEarliestTime.
func (mr *MockCollectorIfaceMockRecorder) CollectEarliestTime(events, isFloored interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectEarliestTime", reflect.TypeOf((*MockCollectorIface)(nil).CollectEarliestTime), events, isFloored)
}

// CollectExecutionIds mocks base method.
func (m *MockCollectorIface) CollectExecutionIds(events []v2.Event) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CollectExecutionIds", events)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CollectExecutionIds indicates an expected call of CollectExecutionIds.
func (mr *MockCollectorIfaceMockRecorder) CollectExecutionIds(events interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectExecutionIds", reflect.TypeOf((*MockCollectorIface)(nil).CollectExecutionIds), events)
}

// CollectLatestTime mocks base method.
func (m *MockCollectorIface) CollectLatestTime(events []v2.Event, isCeiled bool) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CollectLatestTime", events, isCeiled)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CollectLatestTime indicates an expected call of CollectLatestTime.
func (mr *MockCollectorIfaceMockRecorder) CollectLatestTime(events, isCeiled interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectLatestTime", reflect.TypeOf((*MockCollectorIface)(nil).CollectLatestTime), events, isCeiled)
}

// GetEvents mocks base method.
func (m *MockCollectorIface) GetEvents(keptnContext string) ([]v2.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvents", keptnContext)
	ret0, _ := ret[0].([]v2.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvents indicates an expected call of GetEvents.
func (mr *MockCollectorIfaceMockRecorder) GetEvents(keptnContext interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvents", reflect.TypeOf((*MockCollectorIface)(nil).GetEvents), keptnContext)
}

// GetEventsOfType mocks base method.
func (m *MockCollectorIface) GetEventsOfType(eventType, keptnContext string) ([]v2.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventsOfType", eventType, keptnContext)
	ret0, _ := ret[0].([]v2.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventsOfType indicates an expected call of GetEventsOfType.
func (mr *MockCollectorIfaceMockRecorder) GetEventsOfType(eventType, keptnContext interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventsOfType", reflect.TypeOf((*MockCollectorIface)(nil).GetEventsOfType), eventType, keptnContext)
}

// MustParseEventsOfType mocks base method.
func (m *MockCollectorIface) MustParseEventsOfType(events []v2.Event, filterType string) ([]v2.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MustParseEventsOfType", events, filterType)
	ret0, _ := ret[0].([]v2.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MustParseEventsOfType indicates an expected call of MustParseEventsOfType.
func (mr *MockCollectorIfaceMockRecorder) MustParseEventsOfType(events, filterType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MustParseEventsOfType", reflect.TypeOf((*MockCollectorIface)(nil).MustParseEventsOfType), events, filterType)
}

// ParseEventsOfType mocks base method.
func (m *MockCollectorIface) ParseEventsOfType(events []v2.Event, filterType string) []v2.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseEventsOfType", events, filterType)
	ret0, _ := ret[0].([]v2.Event)
	return ret0
}

// ParseEventsOfType indicates an expected call of ParseEventsOfType.
func (mr *MockCollectorIfaceMockRecorder) ParseEventsOfType(events, filterType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseEventsOfType", reflect.TypeOf((*MockCollectorIface)(nil).ParseEventsOfType), events, filterType)
}
