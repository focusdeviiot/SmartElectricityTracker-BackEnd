package models

import (
	"database/sql/driver"
	"fmt"
)

type Role string

const (
	USER  Role = "USER"
	ADMIN Role = "ADMIN"
)

// func (e *Role) Scan(value interface{}) error {
// 	*e = Role(value.([]byte))
// 	return nil
// }

func (e *Role) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*e = Role(v)
	case string:
		*e = Role(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (e Role) Value() (driver.Value, error) {
	return string(e), nil
}

func StringToRole(s string) (Role, error) {
	switch s {
	case "USER":
		return USER, nil
	case "ADMIN":
		return ADMIN, nil
	default:
		return "", fmt.Errorf("invalid role: %s", s)
	}
}
