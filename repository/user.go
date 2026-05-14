package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"magic/utils"
	"strconv"
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

type UserRepositoryInterface interface {
	Validate(email, username string) (bool, bool, error)
	GetByField(field, input string) (*User, error)
	Add(email, password, username string)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Validate will return errors if u has empty/invalid fields
func (u UserRepository) Validate(email, username string) error {

	// Email validation
	_, err := u.GetUserByField("email", email)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return err
	}
	if err == nil {
		return errors.New("Email already taken")
	}

	// if user isn't returned, cant compare values

	_, err = u.GetUserByField("username", username)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return err
	}
	if err == nil {
		return errors.New("Username already used")
	}

	return nil
}

// TODO: add to DB
func (u UserRepository) Add(email, password, username string) error {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := u.db.PrepareContext(c, "INSERT INTO users (email, password, username) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Could not prepare context: ", err)
		return err
	}
	defer stmt.Close()

	// bcrypt password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		fmt.Println("Failed to encrypt password: ", err)
	}

	result, err := stmt.Exec(email, hashedPassword, username)
	if err != nil {
		fmt.Println("Failed to execute stmt: ", err)
		return err
	}

	n, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Failed to get last insert id: ", err)
		return err
	}
	fmt.Println("Successfully added user: ", strconv.Itoa(int(n)))

	return nil
}

// GetUserByField(f, i) queries users for given f WHERE f = i and returns a user if present
func (u UserRepository) GetUserByField(field, input string) (*User, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var prep string

	switch field {
	case "username":
		prep = `SELECT * FROM users WHERE username = ?`
	case "email":
		prep = `SELECT * FROM users WHERE email = ?`
	}

	stmt, err := u.db.PrepareContext(c, prep)
	if err != nil {
		fmt.Println("Could not prepare context: ", err)
		return nil, err
	}
	defer stmt.Close()

	var user User
	row := stmt.QueryRow(input)
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, err
}
