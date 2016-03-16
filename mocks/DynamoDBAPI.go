package mocks

import "github.com/stretchr/testify/mock"

import "github.com/aws/aws-sdk-go/aws/request"
import "github.com/aws/aws-sdk-go/service/dynamodb"

type DynamoDBAPI struct {
	mock.Mock
}

// BatchGetItemRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) BatchGetItemRequest(_a0 *dynamodb.BatchGetItemInput) (*request.Request, *dynamodb.BatchGetItemOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.BatchGetItemInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.BatchGetItemOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.BatchGetItemInput) *dynamodb.BatchGetItemOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.BatchGetItemOutput)
		}
	}

	return r0, r1
}

// BatchGetItem provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) BatchGetItem(_a0 *dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.BatchGetItemOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.BatchGetItemInput) *dynamodb.BatchGetItemOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.BatchGetItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.BatchGetItemInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BatchGetItemPages provides a mock function with given fields: _a0, _a1
func (_m *DynamoDBAPI) BatchGetItemPages(_a0 *dynamodb.BatchGetItemInput, _a1 func(*dynamodb.BatchGetItemOutput, bool) bool) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*dynamodb.BatchGetItemInput, func(*dynamodb.BatchGetItemOutput, bool) bool) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BatchWriteItemRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) BatchWriteItemRequest(_a0 *dynamodb.BatchWriteItemInput) (*request.Request, *dynamodb.BatchWriteItemOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.BatchWriteItemInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.BatchWriteItemOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.BatchWriteItemInput) *dynamodb.BatchWriteItemOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.BatchWriteItemOutput)
		}
	}

	return r0, r1
}

// BatchWriteItem provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) BatchWriteItem(_a0 *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.BatchWriteItemOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.BatchWriteItemInput) *dynamodb.BatchWriteItemOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.BatchWriteItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.BatchWriteItemInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateTableRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) CreateTableRequest(_a0 *dynamodb.CreateTableInput) (*request.Request, *dynamodb.CreateTableOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.CreateTableInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.CreateTableOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.CreateTableInput) *dynamodb.CreateTableOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.CreateTableOutput)
		}
	}

	return r0, r1
}

// CreateTable provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) CreateTable(_a0 *dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.CreateTableOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.CreateTableInput) *dynamodb.CreateTableOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.CreateTableOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.CreateTableInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteItemRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) DeleteItemRequest(_a0 *dynamodb.DeleteItemInput) (*request.Request, *dynamodb.DeleteItemOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.DeleteItemInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.DeleteItemOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.DeleteItemInput) *dynamodb.DeleteItemOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.DeleteItemOutput)
		}
	}

	return r0, r1
}

// DeleteItem provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) DeleteItem(_a0 *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.DeleteItemOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.DeleteItemInput) *dynamodb.DeleteItemOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.DeleteItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.DeleteItemInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteTableRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) DeleteTableRequest(_a0 *dynamodb.DeleteTableInput) (*request.Request, *dynamodb.DeleteTableOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.DeleteTableInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.DeleteTableOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.DeleteTableInput) *dynamodb.DeleteTableOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.DeleteTableOutput)
		}
	}

	return r0, r1
}

// DeleteTable provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) DeleteTable(_a0 *dynamodb.DeleteTableInput) (*dynamodb.DeleteTableOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.DeleteTableOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.DeleteTableInput) *dynamodb.DeleteTableOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.DeleteTableOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.DeleteTableInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DescribeLimitsRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) DescribeLimitsRequest(_a0 *dynamodb.DescribeLimitsInput) (*request.Request, *dynamodb.DescribeLimitsOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.DescribeLimitsInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.DescribeLimitsOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.DescribeLimitsInput) *dynamodb.DescribeLimitsOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.DescribeLimitsOutput)
		}
	}

	return r0, r1
}

// DescribeLimits provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) DescribeLimits(_a0 *dynamodb.DescribeLimitsInput) (*dynamodb.DescribeLimitsOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.DescribeLimitsOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.DescribeLimitsInput) *dynamodb.DescribeLimitsOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.DescribeLimitsOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.DescribeLimitsInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DescribeTableRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) DescribeTableRequest(_a0 *dynamodb.DescribeTableInput) (*request.Request, *dynamodb.DescribeTableOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.DescribeTableInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.DescribeTableOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.DescribeTableInput) *dynamodb.DescribeTableOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.DescribeTableOutput)
		}
	}

	return r0, r1
}

// DescribeTable provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) DescribeTable(_a0 *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.DescribeTableOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.DescribeTableInput) *dynamodb.DescribeTableOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.DescribeTableOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.DescribeTableInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetItemRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) GetItemRequest(_a0 *dynamodb.GetItemInput) (*request.Request, *dynamodb.GetItemOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.GetItemInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.GetItemOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.GetItemInput) *dynamodb.GetItemOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.GetItemOutput)
		}
	}

	return r0, r1
}

// GetItem provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) GetItem(_a0 *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.GetItemOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.GetItemInput) *dynamodb.GetItemOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.GetItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.GetItemInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListTablesRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) ListTablesRequest(_a0 *dynamodb.ListTablesInput) (*request.Request, *dynamodb.ListTablesOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.ListTablesInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.ListTablesOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.ListTablesInput) *dynamodb.ListTablesOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.ListTablesOutput)
		}
	}

	return r0, r1
}

// ListTables provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) ListTables(_a0 *dynamodb.ListTablesInput) (*dynamodb.ListTablesOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.ListTablesOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.ListTablesInput) *dynamodb.ListTablesOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.ListTablesOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.ListTablesInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListTablesPages provides a mock function with given fields: _a0, _a1
func (_m *DynamoDBAPI) ListTablesPages(_a0 *dynamodb.ListTablesInput, _a1 func(*dynamodb.ListTablesOutput, bool) bool) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*dynamodb.ListTablesInput, func(*dynamodb.ListTablesOutput, bool) bool) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PutItemRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) PutItemRequest(_a0 *dynamodb.PutItemInput) (*request.Request, *dynamodb.PutItemOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.PutItemInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.PutItemOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.PutItemInput) *dynamodb.PutItemOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.PutItemOutput)
		}
	}

	return r0, r1
}

// PutItem provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) PutItem(_a0 *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.PutItemOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.PutItemInput) *dynamodb.PutItemOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.PutItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.PutItemInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) QueryRequest(_a0 *dynamodb.QueryInput) (*request.Request, *dynamodb.QueryOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.QueryInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.QueryOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.QueryInput) *dynamodb.QueryOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.QueryOutput)
		}
	}

	return r0, r1
}

// Query provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) Query(_a0 *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.QueryOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.QueryInput) *dynamodb.QueryOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.QueryOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.QueryInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryPages provides a mock function with given fields: _a0, _a1
func (_m *DynamoDBAPI) QueryPages(_a0 *dynamodb.QueryInput, _a1 func(*dynamodb.QueryOutput, bool) bool) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*dynamodb.QueryInput, func(*dynamodb.QueryOutput, bool) bool) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ScanRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) ScanRequest(_a0 *dynamodb.ScanInput) (*request.Request, *dynamodb.ScanOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.ScanInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.ScanOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.ScanInput) *dynamodb.ScanOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.ScanOutput)
		}
	}

	return r0, r1
}

// Scan provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) Scan(_a0 *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.ScanOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.ScanInput) *dynamodb.ScanOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.ScanOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.ScanInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScanPages provides a mock function with given fields: _a0, _a1
func (_m *DynamoDBAPI) ScanPages(_a0 *dynamodb.ScanInput, _a1 func(*dynamodb.ScanOutput, bool) bool) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*dynamodb.ScanInput, func(*dynamodb.ScanOutput, bool) bool) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateItemRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) UpdateItemRequest(_a0 *dynamodb.UpdateItemInput) (*request.Request, *dynamodb.UpdateItemOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.UpdateItemInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.UpdateItemOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.UpdateItemInput) *dynamodb.UpdateItemOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.UpdateItemOutput)
		}
	}

	return r0, r1
}

// UpdateItem provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) UpdateItem(_a0 *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.UpdateItemOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.UpdateItemInput) *dynamodb.UpdateItemOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.UpdateItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.UpdateItemInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTableRequest provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) UpdateTableRequest(_a0 *dynamodb.UpdateTableInput) (*request.Request, *dynamodb.UpdateTableOutput) {
	ret := _m.Called(_a0)

	var r0 *request.Request
	if rf, ok := ret.Get(0).(func(*dynamodb.UpdateTableInput) *request.Request); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*request.Request)
		}
	}

	var r1 *dynamodb.UpdateTableOutput
	if rf, ok := ret.Get(1).(func(*dynamodb.UpdateTableInput) *dynamodb.UpdateTableOutput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dynamodb.UpdateTableOutput)
		}
	}

	return r0, r1
}

// UpdateTable provides a mock function with given fields: _a0
func (_m *DynamoDBAPI) UpdateTable(_a0 *dynamodb.UpdateTableInput) (*dynamodb.UpdateTableOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.UpdateTableOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.UpdateTableInput) *dynamodb.UpdateTableOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.UpdateTableOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.UpdateTableInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
