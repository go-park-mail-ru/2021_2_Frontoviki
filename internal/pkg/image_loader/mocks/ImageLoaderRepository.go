// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import multipart "mime/multipart"

// ImageLoaderRepository is an autogenerated mock type for the ImageLoaderRepository type
type ImageLoaderRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: filePath
func (_m *ImageLoaderRepository) Delete(filePath string) error {
	ret := _m.Called(filePath)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(filePath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Insert provides a mock function with given fields: fileHeader, dir, name
func (_m *ImageLoaderRepository) Insert(fileHeader *multipart.FileHeader, dir string, name string) error {
	ret := _m.Called(fileHeader, dir, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(*multipart.FileHeader, string, string) error); ok {
		r0 = rf(fileHeader, dir, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
