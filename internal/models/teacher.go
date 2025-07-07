package models

type Teacher struct {
	Id        int    `json:"id" db:"id,omitempty"`
	FirstName string `json:"first_name,omitempty" db:"first_name"`
	LastName  string `json:"last_name,omitempty" db:"last_name"`
	Class     string `json:"class,omitempty" db:"class"`
	Subject   string `json:"subject,omitempty" db:"subject"`
	Email     string `json:"email,omitempty" db:"email"`
}
