// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type RequestHistoryMethod string

const (
	RequestHistoryMethodGET    RequestHistoryMethod = "GET"
	RequestHistoryMethodPOST   RequestHistoryMethod = "POST"
	RequestHistoryMethodTRACE  RequestHistoryMethod = "TRACE"
	RequestHistoryMethodHEAD   RequestHistoryMethod = "HEAD"
	RequestHistoryMethodDELETE RequestHistoryMethod = "DELETE"
	RequestHistoryMethodPATCH  RequestHistoryMethod = "PATCH"
	RequestHistoryMethodPUT    RequestHistoryMethod = "PUT"
)

func (e *RequestHistoryMethod) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RequestHistoryMethod(s)
	case string:
		*e = RequestHistoryMethod(s)
	default:
		return fmt.Errorf("unsupported scan type for RequestHistoryMethod: %T", src)
	}
	return nil
}

type NullRequestHistoryMethod struct {
	RequestHistoryMethod RequestHistoryMethod
	Valid                bool // Valid is true if RequestHistoryMethod is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRequestHistoryMethod) Scan(value interface{}) error {
	if value == nil {
		ns.RequestHistoryMethod, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RequestHistoryMethod.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRequestHistoryMethod) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RequestHistoryMethod), nil
}

type RequestHistory struct {
	RequestID   sql.NullInt16
	Endpoint    string
	Headers     sql.NullString
	Method      NullRequestHistoryMethod
	Body        sql.NullString
	RequestTime time.Time
}
