// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ImageCompressorRepository is an autogenerated mock type for the ImageCompressorRepository type
type ImageCompressorRepository struct {
	mock.Mock
}

// Compress provides a mock function with given fields: filepath, filename, extension
func (_m *ImageCompressorRepository) Compress(filepath string, filename string, extension string) error {
	ret := _m.Called(filepath, filename, extension)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(filepath, filename, extension)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: dirpath, filename
func (_m *ImageCompressorRepository) Delete(dirpath string, filename string) error {
	ret := _m.Called(dirpath, filename)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(dirpath, filename)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
