package models

import (
	"database/sql"
	// "errors"
	// "encoding/json"
	"fmt"
)

type ScoreConfig struct {
	Name		string
	Description	string
	Multiplier	float64
}

type ScoreSheet struct {
	ID				int
	GameID			int
	ScoreConfigs	[]ScoreConfig
}

type ScoresheetModel struct {
	DB		*sql.DB
}

func (m *ScoresheetModel) Insert(gameID int, configs ScoreConfig) (int, error) {
	// Save score settings into the scoresheet table

	stmt := `INSERT INTO scoresheets (game_id, scoredata, updated)
	VALUES(?, ?, UTC_TIMESTAMP())`

	fmt.Println(stmt)


	// result, err := m.DB.Exec(stmt, title,  bgg_id, price, state, year_pub,
	// description, player_min, player_max, playtime_min, playtime_max,
	// age_min, scoring)

	// if err != nil {
	// 	return 0, err
	// }

	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }

	return 0, nil
	// return int(id), nil
}