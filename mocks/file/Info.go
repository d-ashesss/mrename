// Code generated by mockery v3.0.0-alpha.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Info is an autogenerated mock type for the Info type
type Info struct {
	mock.Mock
}

// Name provides a mock function with given fields:
func (_m *Info) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Path provides a mock function with given fields:
func (_m *Info) Path() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewInfo interface {
	mock.TestingT
	Cleanup(func())
}

// NewInfo creates a new instance of Info. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewInfo(t mockConstructorTestingTNewInfo) *Info {
	mock := &Info{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
