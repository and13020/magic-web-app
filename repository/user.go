package repository

import (
	"errors"
	"time"
)

type User struct {
	ID        string    `json:"id" sql:"id"`
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	Username  string    `json:"username" sql:"username"`
	TokenHash string    `json:"tokenhash" sql:"tokenhash"`
	CreatedAt time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt time.Time `json:"updatedat" sql:"updatedat"`
}

// TODO: other validations need to be done
// ie email is unique, min/max length etc
func (usr User) Validate() []error {
	var e []error
	if usr.Email == "" {
		e = append(e, errors.New("Email is missing"))
	}
	if usr.Password == "" {
		e = append(e, errors.New("Password is missing"))
	}

	return e
}
