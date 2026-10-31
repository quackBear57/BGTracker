package models

import (
	"database/sql"
	"errors"
	// "strings"
	"time"

	// "github.com/go-sql-driver/mysql"
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

func (m *PlayerModel) Insert(name string, userID int) (int, error) {
	// Add new player to database
	
	stmt := `INSERT INTO players (name, user_id, created)
	VALUES(?, ?, UTC_TIMESTAMP())`

	// NEED TO CONFIRM USERID EXSITS BEFORE CALLING THIS
	// also check that non-zero user id's are unique
	
	result, err := m.DB.Exec(stmt, name, userID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PlayerModel) Exists(id int) (bool, error) {
	// Check if player ID exists

	// Will need to confrim id is non-zero before calling

	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM players WHERE id = ?)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	
	return exists, err
}

func (m *PlayerModel) UserIDTaken(id int) (bool, error) {
	// Check if player user ID exists in players

	// not currently working

	var result string
	// var taken bool

	stmt := "SELECT name FROM players WHERE user_id = ?"

	err := m.DB.QueryRow(stmt, id).Scan(&result)
	
	return result != "", err
}

func (m *PlayerModel) Get(id int) (Player, error) {
	// Get all information about player from DB

	stmt := `SELECT id, name, user_id, created FROM players WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	var p Player

	err := row.Scan(&p.ID, &p.Name, &p.UserID, &p.Created)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Player{}, ErrNoRecord
		default:
			return Player{}, err
		}
	}

	return p, nil
}

func (m *PlayerModel) UpdateName(id int, newName string) error {
	// Update a name for a given player ID

	stmt := `UPDATE players SET name = ? WHERE id = ?`

	_, err := m.DB.Exec(stmt, newName, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PlayerModel) Latest() ([]Player, error) {
	stmt := `SELECT id, name, user_id, created FROM players WHERE created < DATE_ADD(UTC_TIMESTAMP(), 
	INTERVAL ? DAY) ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt, 30)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var players []Player

	for rows.Next() {
		var p Player
		
		err = rows.Scan(&p.ID, &p.Name, &p.UserID, &p.Created)
		
		if err != nil {
			return nil, err
		}

		players = append(players, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return players, nil
}