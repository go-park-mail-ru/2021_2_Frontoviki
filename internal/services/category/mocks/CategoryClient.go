// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"
	category "yula/proto/generated/category"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// CategoryClient is an autogenerated mock type for the CategoryClient type
type CategoryClient struct {
	mock.Mock
}

// GetCategories provides a mock function with given fields: ctx, in, opts
func (_m *CategoryClient) GetCategories(ctx context.Context, in *category.Nothing, opts ...grpc.CallOption) (*category.Categories, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *category.Categories
	if rf, ok := ret.Get(0).(func(context.Context, *category.Nothing, ...grpc.CallOption) *category.Categories); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*category.Categories)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *category.Nothing, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
