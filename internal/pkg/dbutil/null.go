package dbutil

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type NullRawMessage struct {
	RawMessage json.RawMessage
	Valid      bool
}

func (n *NullRawMessage) Scan(value interface{}) error {
	if value == nil {
		n.RawMessage, n.Valid = nil, false
		return nil
	}
	n.Valid = true
	switch v := value.(type) {
	case []byte:
		n.RawMessage = json.RawMessage(v)
	case string:
		n.RawMessage = json.RawMessage(v)
	default:
		n.RawMessage = nil
	}
	return nil
}

func (n NullRawMessage) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return []byte(n.RawMessage), nil
}

func NewNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func NewNullInt64(i int64, valid bool) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: valid}
}

func NewNullFloat64(f float64, valid bool) sql.NullFloat64 {
	return sql.NullFloat64{Float64: f, Valid: valid}
}

func NewNullBool(b bool, valid bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: valid}
}

func NullStringToPtr(ns sql.NullString) *string {
	if ns.Valid && ns.String != "" {
		return &ns.String
	}
	return nil
}

func NullStringValue(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func NullInt64ToIntPtr(ni sql.NullInt64) *int {
	if ni.Valid {
		v := int(ni.Int64)
		return &v
	}
	return nil
}

func NullInt64ToInt64Ptr(ni sql.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}
