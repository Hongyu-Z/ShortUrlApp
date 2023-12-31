// Code generated by mockery v2.36.1. DO NOT EDIT.

package dao

import (
	models "ShortUrlApp/models"

	mock "github.com/stretchr/testify/mock"
)

// MockUrlRecordDao is an autogenerated mock type for the UrlRecordDao type
type MockUrlRecordDao struct {
	mock.Mock
}

// Delete provides a mock function with given fields: shortUrl
func (_m *MockUrlRecordDao) Delete(shortUrl string) error {
	ret := _m.Called(shortUrl)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(shortUrl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: shortUrl
func (_m *MockUrlRecordDao) Find(shortUrl string) (*models.UrlRecord, error) {
	ret := _m.Called(shortUrl)

	var r0 *models.UrlRecord
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.UrlRecord, error)); ok {
		return rf(shortUrl)
	}
	if rf, ok := ret.Get(0).(func(string) *models.UrlRecord); ok {
		r0 = rf(shortUrl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UrlRecord)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(shortUrl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: record
func (_m *MockUrlRecordDao) Save(record *models.UrlRecord) error {
	ret := _m.Called(record)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.UrlRecord) error); ok {
		r0 = rf(record)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockUrlRecordDao creates a new instance of MockUrlRecordDao. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUrlRecordDao(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUrlRecordDao {
	mock := &MockUrlRecordDao{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
