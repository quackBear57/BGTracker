package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Player struct {
	ID				int
	Name			string
	UserID			int
	Created			time.Time
}

type PlayerModel struct {
	DB	*sql.DB
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func (m *PlayerModel) Insert(name string, userID int) error {
	// Add new player to database
	
	stmt := `INSERT INTO players (name, user_id, created)
	VALUES(?, ?, UTC_TIMESTAMP())`

	// NEED TO CONFIRM USERID EXSITS BEFORE CALLING THIS
	// also check that non-zero user id's are unique
	
	_, err = m.DB.Exec(stmt, name, userID)
	if err != nil {
		return err
	}

	return nil
}

func (m *PlayerModel) Exists(id int) (bool, error) {
	// Check if player ID exists

	// Will need to confrim id is non-zero before calling

	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM players WHERE id = ?)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	
	return exists, err
}

func (m *PlayerModel) Get(id int) (User, error) {
	// Get all information about player from DB

	stmt := `SELECT id, name, user_id, created FROM players WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	var p Player

	err := row.Scan(&p.ID, &p.Name, &p.UserID, &p.Created)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return User{}, ErrNoRecord
		default:
			return User{}, err
		}
	}

	return p, nil
}

func (m *PlayerModel) UpdateName(id int, newName string) error {
	// Update a name for a given player ID

	stmt := `UPDATE players SET name = ? WHERE id = ?`

	_, err = m.DB.Exec(stmt, newName, id)
	if err != nil {
		return err
	}

	return nil
}