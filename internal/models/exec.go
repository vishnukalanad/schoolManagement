package models

import "database/sql"

type Exec struct {
	Id                 int            `json:"id,omitempty" db:"id,omitempty"`
	FirstName          string         `json:"first_name,omitempty" db:"first_name,omitempty"`
	LastName           string         `json:"last_name,omitempty" db:"last_name,omitempty"`
	Email              string         `json:"email,omitempty" db:"email,omitempty"`
	Username           string         `json:"username,omitempty" db:"username,omitempty"`
	Password           string         `json:"password,omitempty" db:"password,omitempty"`
	PasswordUpdatedAt  sql.NullString `json:"password_updated_at,omitempty" db:"password_updated_at,omitempty"`
	CreatedAt          sql.NullString `json:"created_at,omitempty" db:"created_at,omitempty"`
	PasswordResetCode  sql.NullString `json:"password_reset_code,omitempty" db:"password_reset_code,omitempty"`
	PasswordCodeExpiry sql.NullString `json:"password_code_expiry,omitempty" db:"password_code_expiry,omitempty"`
	Inactive           bool           `json:"inactive_status,omitempty" db:"inactive_status,omitempty"`
	Role               string         `json:"role,omitempty" db:"role,omitempty"`
}
