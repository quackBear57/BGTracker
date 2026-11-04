package models

import (
	"database/sql"
	"errors"
	"encoding/json"
	"fmt"
	"time"
)

type ScoreConfig struct {
	Name		string
	Description	string
	Multiplier	float64
}

type Scoresheet struct {
	ID				int
	GameID			int
	ScoreConfigs	[]ScoreConfig
	Updated			time.Time
}

type ScoresheetModel struct {
	DB		*sql.DB
}

func (m *ScoresheetModel) Insert(gameID int, configs ScoreConfig) (int, error) {
	// Save score settings into the scoresheet table

	stmt := `INSERT INTO scoresheets (game_id, scoredata, updated)
	VALUES(?, ?, UTC_TIMESTAMP())`

	fmt.Println(stmt)

	scoreConfigsJSON, _ := json.Marshal(configs)

	result, err := m.DB.Exec(stmt, gameID, scoreConfigsJSON)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *ScoresheetModel) GetByGame(gameID int) (Scoresheet, error) {
	// Get all information about player from DB

	stmt := `SELECT id, game_id, scoredata, updated FROM scoresheets WHERE game_id = ?`

	row := m.DB.QueryRow(stmt, gameID)

	var s Scoresheet
	var scoreConfigsJSON []byte

	err := row.Scan(&s.ID, &s.GameID, &scoreConfigsJSON, &s.Updated)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Scoresheet{}, ErrNoRecord
		default:
			return Scoresheet{}, err
		}
	}

	if len(scoreConfigsJSON) > 0 {
		err = json.Unmarshal(scoreConfigsJSON, s.ScoreConfigs)
		if err != nil {
			fmt.Println("Error in unmarshal JSON")
			return Scoresheet{}, err
		}
	}

	return s, nil
}

func (m *ScoresheetModel) ExistsByGame(gameID int) (bool, error) {
	// Check if player user ID exists in players

	var exists bool

	stmt := "SELECT true WHERE EXISTS(SELECT true FROM scoresheets WHERE game_id = ?)"

	err := m.DB.QueryRow(stmt, gameID).Scan(&exists)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, nil
		default:
			return false, err
		}
	}
		
	return exists, nil
}