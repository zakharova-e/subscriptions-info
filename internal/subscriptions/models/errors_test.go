package models_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/zakharova-e/subscriptions-info/internal/subscriptions/models"
)

func TestValidationError_Error(t *testing.T) {
	type fields struct {
		Errors []error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test without errors inside", fields{}, "invalid request data: "},
		{"test with one error inside", fields{[]error{errors.New("error 1 message")}}, "invalid request data: error 1 message "},
		{"test with three errors inside", fields{
			[]error{
				errors.New("error 1 message"),
				errors.New("error 2 message"),
				errors.New("error 3 message"),
			},
		}, "invalid request data: error 1 message error 2 message error 3 message "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &models.ValidationError{
				Errors: tt.fields.Errors,
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("ValidationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationError_Unwrap(t *testing.T) {
	type fields struct {
		Errors []error
	}
	firstError, secondError, thirdError := errors.New("first error"), errors.New("second error"), errors.New("third error")
	tests := []struct {
		name     string
		fields   fields
		expected error
	}{
		{"test with empty errors slice", fields{Errors: nil}, nil},
		{"test with single error inside", fields{Errors: []error{firstError}}, errors.Join(firstError)},
		{"test with two errors inside", fields{Errors: []error{firstError, secondError}}, errors.Join(firstError, secondError)},
		{"test with three errors inside", fields{Errors: []error{firstError, secondError, thirdError}}, errors.Join(firstError, secondError, thirdError)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &models.ValidationError{
				Errors: tt.fields.Errors,
			}
			if errUnwrapped := err.Unwrap(); errUnwrapped != tt.expected {
				if errUnwrapped == nil || tt.expected == nil {
					t.Errorf("ValidationError.Unwrap() error = %v(%T), expected %v(%T)", errUnwrapped, errUnwrapped, tt.expected, tt.expected)
				} else if errUnwrapped.Error() != tt.expected.Error() {
					t.Errorf("ValidationError.Unwrap() error = %v(%T), expected %v(%T)", errUnwrapped.Error(), errUnwrapped, tt.expected.Error(), tt.expected)
				}
			}
		})
	}
}

func TestResourceNotFoundError_Error(t *testing.T) {
	type fields struct {
		Err error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test with nil error", fields{Err: nil}, "resource not found: <nil>"},
		{"test with sql no rows error", fields{Err: sql.ErrNoRows}, "resource not found: sql: no rows in result set"},
		{"test with custom error", fields{Err: errors.New("custom err")}, "resource not found: custom err"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &models.ResourceNotFoundError{
				Err: tt.fields.Err,
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("ResourceNotFoundError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceNotFoundError_Unwrap(t *testing.T) {
	type fields struct {
		Err error
	}
	customErr := errors.New("custom err")
	tests := []struct {
		name     string
		fields   fields
		expected error
	}{
		{"test with nil error", fields{Err: nil}, nil},
		{"test with sql no rows error", fields{Err: sql.ErrNoRows}, sql.ErrNoRows},
		{"test with custom error", fields{Err: customErr}, customErr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &models.ResourceNotFoundError{
				Err: tt.fields.Err,
			}
			if err := err.Unwrap(); err != tt.expected {
				t.Errorf("ResourceNotFoundError.Unwrap() error = %v, wantErr %v", err, tt.expected)
			}
		})
	}
}

func TestDatabaseError_Error(t *testing.T) {
	type fields struct {
		Query string
		Err   error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test with nil error", fields{Err: nil}, "error during the query execution:  <nil>"},
		{"test with sql no rows error", fields{Query: "select * from table", Err: sql.ErrNoRows}, "error during the query execution: select * from table sql: no rows in result set"},
		{"test with custom error", fields{Query: "select * from table1 join table2", Err: errors.New("custom err")}, "error during the query execution: select * from table1 join table2 custom err"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &models.DatabaseError{
				Query: tt.fields.Query,
				Err:   tt.fields.Err,
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("DatabaseError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseError_Unwrap(t *testing.T) {
	type fields struct {
		Query string
		Err   error
	}
	customErr := errors.New("custom err")
	tests := []struct {
		name     string
		fields   fields
		expected error
	}{
		{"test with nil error", fields{Err: nil}, nil},
		{"test with sql no rows error", fields{Query: "select * from table", Err: sql.ErrNoRows}, sql.ErrNoRows},
		{"test with custom error", fields{Query: "select * from table1 join table2", Err: customErr}, customErr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &models.DatabaseError{
				Query: tt.fields.Query,
				Err:   tt.fields.Err,
			}
			if err := err.Unwrap(); err != tt.expected {
				t.Errorf("DatabaseError.Unwrap() error = %v, wantErr %v", err, tt.expected)
			}
		})
	}
}

func TestJsonError_Error(t *testing.T) {
	type fields struct {
		Json string
		Err  error
	}
	customVal := reflect.ValueOf("string")
	customErr1, customErr2 := json.MarshalerError{Err: errors.New("custom err 1"), Type: customVal.Type()}, json.UnmarshalTypeError{Value: "test", Struct: "struct test", Type: customVal.Type()}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test with empty error", fields{Err: nil}, "error during json coding/encoding:  <nil>"},
		{"test with marshal error inside", fields{Json: "{\"test\":123}", Err: &customErr1}, "error during json coding/encoding: {\"test\":123} json: error calling MarshalJSON for type string: custom err 1"},
		{"test with unmarshal error inside", fields{Json: "{\"test\":123}", Err: &customErr2}, "error during json coding/encoding: {\"test\":123} json: cannot unmarshal test into Go struct field struct test. of type string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &models.JsonError{
				Json: tt.fields.Json,
				Err:  tt.fields.Err,
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("JsonError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonError_Unwrap(t *testing.T) {
	type fields struct {
		Json string
		Err  error
	}
	customVal := reflect.ValueOf("string")
	customErr1, customErr2 := json.MarshalerError{Err: errors.New("custom err 1"), Type: customVal.Type()}, json.UnmarshalTypeError{Value: "test", Struct: "struct test", Type: customVal.Type()}
	tests := []struct {
		name     string
		fields   fields
		expected error
	}{
		{"test with empty error", fields{Err: nil}, nil},
		{"test with marshal error inside", fields{Json: "{\"test\":123}", Err: &customErr1}, &customErr1},
		{"test with unmarshal error inside", fields{Json: "{\"test\":123}", Err: &customErr2}, &customErr2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &models.JsonError{
				Json: tt.fields.Json,
				Err:  tt.fields.Err,
			}
			if err := err.Unwrap(); err != tt.expected {
				t.Errorf("JsonError.Unwrap() error = %v, wantErr %v", err, tt.expected)
			}
		})
	}
}

func TestMethodNotAllowedError_Error(t *testing.T) {
	type fields struct {
		RequiredMethod string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test with empty error", fields{RequiredMethod: ""}, "Only  method is allowed!"},
		{"test with only get allowed error", fields{RequiredMethod: "GET"}, "Only GET method is allowed!"},
		{"test with only post allowed error", fields{RequiredMethod: "POST"}, "Only POST method is allowed!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &models.MethodNotAllowedError{
				RequiredMethod: tt.fields.RequiredMethod,
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("MethodNotAllowedError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMethodNotAllowedError_Unwrap(t *testing.T) {
	type fields struct {
		RequiredMethod string
	}
	tests := []struct {
		name     string
		fields   fields
		expected error
	}{
		{"test with empty error", fields{RequiredMethod: ""}, nil},
		{"test with only get allowed error", fields{RequiredMethod: "GET"}, nil},
		{"test with only post allowed error", fields{RequiredMethod: "POST"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &models.MethodNotAllowedError{
				RequiredMethod: tt.fields.RequiredMethod,
			}
			if err := err.Unwrap(); err != tt.expected {
				t.Errorf("MethodNotAllowedError.Unwrap() error = %v, wantErr %v", err, tt.expected)
			}
		})
	}
}

func TestInvalidParameterError_Error(t *testing.T) {
	type fields struct {
		ParamName string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test with empty error", fields{ParamName: ""}, "invalid parameter: "},
		{"test with page param error", fields{ParamName: "page"}, "invalid parameter: page"},
		{"test with date param error", fields{ParamName: "date"}, "invalid parameter: date"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &models.InvalidParameterError{
				ParamName: tt.fields.ParamName,
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("InvalidParameterError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidParameterError_Unwrap(t *testing.T) {
	type fields struct {
		ParamName string
	}
	tests := []struct {
		name     string
		fields   fields
		expected error
	}{
		{"test with empty error", fields{ParamName: ""}, nil},
		{"test with page param error", fields{ParamName: "page"}, nil},
		{"test with date param error", fields{ParamName: "date"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &models.InvalidParameterError{
				ParamName: tt.fields.ParamName,
			}
			if err := err.Unwrap(); err != tt.expected {
				t.Errorf("InvalidParameterError.Unwrap() error = %v, wantErr %v", err, tt.expected)
			}
		})
	}
}
