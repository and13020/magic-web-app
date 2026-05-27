package repository

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewUserRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	if repo == nil {
		t.Fatal("expected non-nil repository")
	}
	if repo.db != db {
		t.Fatal("expected repository db to match provided db")
	}
}

func TestValidate_EmailAlreadyTaken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("SELECT \\* FROM users WHERE email = \\?").
		ExpectQuery().
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "password"}).
			AddRow("1", "test@example.com", "testuser", "hashedpass"))

	err = repo.Validate("test@example.com", "newuser")
	if err == nil {
		t.Fatal("expected error for email already taken")
	}
	if err.Error() != "Email already taken" {
		t.Fatalf("expected 'Email already taken', got '%s'", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestValidate_UsernameAlreadyUsed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("SELECT \\* FROM users WHERE email = \\?").
		ExpectQuery().
		WithArgs("new@example.com").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectPrepare("SELECT \\* FROM users WHERE username = \\?").
		ExpectQuery().
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "password"}).
			AddRow("1", "test@example.com", "testuser", "hashedpass"))

	err = repo.Validate("new@example.com", "testuser")
	if err == nil {
		t.Fatal("expected error for username already used")
	}
	if err.Error() != "Username already used" {
		t.Fatalf("expected 'Username already used', got '%s'", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestValidate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("SELECT \\* FROM users WHERE email = \\?").
		ExpectQuery().
		WithArgs("new@example.com").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectPrepare("SELECT \\* FROM users WHERE username = \\?").
		ExpectQuery().
		WithArgs("newuser").
		WillReturnError(sql.ErrNoRows)

	err = repo.Validate("new@example.com", "newuser")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestValidate_EmailQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	dbError := errors.New("database connection error")
	mock.ExpectPrepare("SELECT \\* FROM users WHERE email = \\?").
		ExpectQuery().
		WithArgs("test@example.com").
		WillReturnError(dbError)

	err = repo.Validate("test@example.com", "newuser")
	if err == nil {
		t.Fatal("expected error from database")
	}
	if err != dbError {
		t.Fatalf("expected database error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestValidate_UsernameQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("SELECT \\* FROM users WHERE email = \\?").
		ExpectQuery().
		WithArgs("new@example.com").
		WillReturnError(sql.ErrNoRows)

	dbError := errors.New("database connection error")
	mock.ExpectPrepare("SELECT \\* FROM users WHERE username = \\?").
		ExpectQuery().
		WithArgs("testuser").
		WillReturnError(dbError)

	err = repo.Validate("new@example.com", "testuser")
	if err == nil {
		t.Fatal("expected error from database")
	}
	if err != dbError {
		t.Fatalf("expected database error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAdd_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WithArgs("test@example.com", sqlmock.AnyArg(), "testuser").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Add("test@example.com", "password123", "testuser")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAdd_PrepareContextError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	prepareError := errors.New("prepare context error")
	mock.ExpectPrepare("INSERT INTO users").WillReturnError(prepareError)

	err = repo.Add("test@example.com", "password123", "testuser")
	if err == nil {
		t.Fatal("expected error from prepare context")
	}
	if err != prepareError {
		t.Fatalf("expected prepare error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAdd_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	execError := errors.New("exec error")
	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WithArgs("test@example.com", sqlmock.AnyArg(), "testuser").
		WillReturnError(execError)

	err = repo.Add("test@example.com", "password123", "testuser")
	if err == nil {
		t.Fatal("expected error from exec")
	}
	if err != execError {
		t.Fatalf("expected exec error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAdd_LastInsertIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	lastInsertIdError := errors.New("last insert id error")
	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WithArgs("test@example.com", sqlmock.AnyArg(), "testuser").
		WillReturnResult(sqlmock.NewErrorResult(lastInsertIdError))

	err = repo.Add("test@example.com", "password123", "testuser")
	if err == nil {
		t.Fatal("expected error from last insert id")
	}
	if err != lastInsertIdError {
		t.Fatalf("expected last insert id error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserByField_ByUsername_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	expectedUser := &User{
		ID:       "1",
		Email:    "test@example.com",
		Username: "testuser",
		Password: "hashedpass",
	}

	mock.ExpectPrepare("SELECT \\* FROM users WHERE username = \\?").
		ExpectQuery().
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "password"}).
			AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Username, expectedUser.Password))

	user, err := repo.GetUserByField("username", "testuser")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != expectedUser.ID {
		t.Fatalf("expected ID %s, got %s", expectedUser.ID, user.ID)
	}
	if user.Email != expectedUser.Email {
		t.Fatalf("expected email %s, got %s", expectedUser.Email, user.Email)
	}
	if user.Username != expectedUser.Username {
		t.Fatalf("expected username %s, got %s", expectedUser.Username, user.Username)
	}
	if user.Password != expectedUser.Password {
		t.Fatalf("expected password %s, got %s", expectedUser.Password, user.Password)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserByField_ByEmail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	expectedUser := &User{
		ID:       "2",
		Email:    "test2@example.com",
		Username: "testuser2",
		Password: "hashedpass2",
	}

	mock.ExpectPrepare("SELECT \\* FROM users WHERE email = \\?").
		ExpectQuery().
		WithArgs("test2@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "password"}).
			AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Username, expectedUser.Password))

	user, err := repo.GetUserByField("email", "test2@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != expectedUser.ID {
		t.Fatalf("expected ID %s, got %s", expectedUser.ID, user.ID)
	}
	if user.Email != expectedUser.Email {
		t.Fatalf("expected email %s, got %s", expectedUser.Email, user.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserByField_PrepareContextError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	prepareError := errors.New("prepare context error")
	mock.ExpectPrepare("SELECT \\* FROM users WHERE username = \\?").
		WillReturnError(prepareError)

	user, err := repo.GetUserByField("username", "testuser")
	if err == nil {
		t.Fatal("expected error from prepare context")
	}
	if user != nil {
		t.Fatal("expected nil user")
	}
	if err != prepareError {
		t.Fatalf("expected prepare error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserByField_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("SELECT \\* FROM users WHERE email = \\?").
		ExpectQuery().
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "password"}).
			AddRow("1", "test@example.com", "testuser", nil).
			RowError(0, errors.New("scan error")))

	user, err := repo.GetUserByField("email", "test@example.com")
	if err == nil {
		t.Fatal("expected error from scan")
	}
	if user != nil {
		t.Fatal("expected nil user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserByField_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("SELECT \\* FROM users WHERE username = \\?").
		ExpectQuery().
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByField("username", "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if user != nil {
		t.Fatal("expected nil user")
	}
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserByField_EmptyPrepStatement(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("").
		WillReturnError(errors.New("empty statement"))

	user, err := repo.GetUserByField("invalid_field", "value")
	if err == nil {
		t.Fatal("expected error for invalid field")
	}
	if user != nil {
		t.Fatal("expected nil user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAdd_HashPasswordError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	longPassword := string(make([]byte, 73))

	mock.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WithArgs("test@example.com", sqlmock.AnyArg(), "testuser").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Add("test@example.com", longPassword, "testuser")
	if err != nil {
		t.Fatalf("expected no error even with hash error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAdd_ContextTimeout(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("INSERT INTO users").
		WillDelayFor(5 * time.Second).
		WillReturnError(errors.New("context deadline exceeded"))

	err = repo.Add("test@example.com", "password123", "testuser")
	if err == nil {
		t.Fatal("expected error from context timeout")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserByField_ContextTimeout(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("SELECT \\* FROM users WHERE email = \\?").
		WillDelayFor(5 * time.Second).
		WillReturnError(errors.New("context deadline exceeded"))

	user, err := repo.GetUserByField("email", "test@example.com")
	if err == nil {
		t.Fatal("expected error from context timeout")
	}
	if user != nil {
		t.Fatal("expected nil user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
