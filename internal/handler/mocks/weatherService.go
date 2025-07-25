// Code generated by MockGen. DO NOT EDIT.
// Source: weatherService.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	handler "weather-service/internal/handler"

	gomock "github.com/golang/mock/gomock"
)

// MockForecastClient is a mock of ForecastClient interface.
type MockForecastClient struct {
	ctrl     *gomock.Controller
	recorder *MockForecastClientMockRecorder
}

// MockForecastClientMockRecorder is the mock recorder for MockForecastClient.
type MockForecastClientMockRecorder struct {
	mock *MockForecastClient
}

// NewMockForecastClient creates a new mock instance.
func NewMockForecastClient(ctrl *gomock.Controller) *MockForecastClient {
	mock := &MockForecastClient{ctrl: ctrl}
	mock.recorder = &MockForecastClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockForecastClient) EXPECT() *MockForecastClientMockRecorder {
	return m.recorder
}

// GetForecast mocks base method.
func (m *MockForecastClient) GetForecast(lat, long string) (handler.ForecastMap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForecast", lat, long)
	ret0, _ := ret[0].(handler.ForecastMap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetForecast indicates an expected call of GetForecast.
func (mr *MockForecastClientMockRecorder) GetForecast(lat, long interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForecast", reflect.TypeOf((*MockForecastClient)(nil).GetForecast), lat, long)
}

// MockCache is a mock of Cache interface.
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance.
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockCache) Get(key string) (*handler.CachedWeather, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(*handler.CachedWeather)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCacheMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCache)(nil).Get), key)
}

// Put mocks base method.
func (m *MockCache) Put(key string, weather *handler.CachedWeather) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", key, weather)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockCacheMockRecorder) Put(key, weather interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockCache)(nil).Put), key, weather)
}
