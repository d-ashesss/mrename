// Code generated by mockery v3.0.0-alpha.0. DO NOT EDIT.

package mocks

import (
	observer "github.com/d-ashesss/mrename/observer"
	mock "github.com/stretchr/testify/mock"
)

// Subscriber is an autogenerated mock type for the Subscriber type
type Subscriber struct {
	mock.Mock
}

// Notify provides a mock function with given fields: _a0
func (_m *Subscriber) Notify(_a0 observer.Event) {
	_m.Called(_a0)
}

type mockConstructorTestingTNewSubscriber interface {
	mock.TestingT
	Cleanup(func())
}

// NewSubscriber creates a new instance of Subscriber. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSubscriber(t mockConstructorTestingTNewSubscriber) *Subscriber {
	mock := &Subscriber{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}