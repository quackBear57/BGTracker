package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID				int
	Name			string
	Email			string
	HashedPassword	[]byte
	Created			time.Time
}

type UserModel struct {
	DB	*sql.DB
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	
	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLErr *mysql.MySQLError
		switch {
		case errors.As(err, &mySQLErr) && mySQLErr.Number == 1062 && 
			strings.Contains(mySQLErr.Message, "users_uc_email"):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrInvalidCredentials
		default:
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return 0, ErrInvalidCredentials
		default:
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	
	return exists, err
}

func (m *UserModel) Get(id int) (User, error) {
	stmt := `SELECT id, name, email, created FROM users WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	var u User

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Created)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return User{}, ErrNoRecord
		default:
			return User{}, err
		}
	}

	return u, nil
}

func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) (bool, error) {
	var hashedCurrentPassword []byte

	stmt := `SELECT hashed_password FROM users WHERE id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&hashedCurrentPassword)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, ErrNoUser
		default:
			return false, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedCurrentPassword, []byte(currentPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, ErrInvalidCredentials
		default:
			return false, err
		}
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return false, err
	}

	updateStmt := `UPDATE users SET hashed_password = ? WHERE id = ?`

	_, err = m.DB.Exec(updateStmt, string(hashedNewPassword), id)
	if err != nil {
		return false, err
	}

	return true, nil
}