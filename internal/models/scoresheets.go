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

func (m *ScoresheetModel) Insert(gameID int, configs []ScoreConfig) (error) {
	// Save score settings into the scoresheet table

	scoreConfigsJSON, _ := json.Marshal(configs)

	stmt := `INSERT INTO scoresheets (game_id, scoredata, updated)
	VALUES(?, ?, UTC_TIMESTAMP())`

	_, err := m.DB.Exec(stmt, gameID, scoreConfigsJSON)
	if err != nil {
		return err
	}

	return nil
}

func (m *ScoresheetModel) Update(gameID int, configs []ScoreConfig) (error) {
	// Update existing scoresheet

	scoreConfigsJSON, _ := json.Marshal(configs)

	stmt := `UPDATE scoresheets 
	SET scoredata = ?, updated = UTC_TIMESTAMP() 
	WHERE game_id = ?`

	_, err := m.DB.Exec(stmt, scoreConfigsJSON, gameID)
	if err != nil {
		return err
	}

	return nil
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
		err = json.Unmarshal(scoreConfigsJSON, &s.ScoreConfigs)
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