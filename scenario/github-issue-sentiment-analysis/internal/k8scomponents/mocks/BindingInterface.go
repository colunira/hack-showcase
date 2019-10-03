// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
import v1beta1 "github.com/poy/service-catalog/pkg/apis/servicecatalog/v1beta1"

// BindingInterface is an autogenerated mock type for the BindingInterface type
type BindingInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *BindingInterface) Create(_a0 *v1beta1.ServiceBinding) (*v1beta1.ServiceBinding, error) {
	ret := _m.Called(_a0)

	var r0 *v1beta1.ServiceBinding
	if rf, ok := ret.Get(0).(func(*v1beta1.ServiceBinding) *v1beta1.ServiceBinding); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1beta1.ServiceBinding)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1beta1.ServiceBinding) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: name, options
func (_m *BindingInterface) Delete(name string, options *v1.DeleteOptions) error {
	ret := _m.Called(name, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *v1.DeleteOptions) error); ok {
		r0 = rf(name, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}