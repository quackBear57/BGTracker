package models

import (
	"database/sql"
	"time"
)

type  Game struct {
	ID				int
	Title			string
	Bgg_ID			int
	price			float64
	added			time.Time
	state			string
	year_pub		int
	description		string
	player_min		int
	player_max		int
	playtime_min	int
	playtime_max	int
	scoring			bool
}

type GameModel struct {
	DB				*sql.DB
}

func (m *GameModel) Insert(title string,
	bgg_id int,
	price float64,
	state string,
	year_pub int,
	description string,
	player_min int,
	player_max int,
	playtime_min int,
	playtime_max int,
	age_min int,
	scoring bool) (int, error) {
	
	stmt := `INSERT INTO games (title, bgg_id, price, added, state, year_pub,
	description, player_min, player_max, playtime_min, playtime_max,
	age_min, scoring) VALUES (?, ?, ?, UTC_TIMESTAMP(), ?, ?,
	?, ?, ?, ?, ?,
	?, ?))`

	result, err := m.DB.Exec(stmt, title,  bgg_id, price, state, year_pub,
	description, player_min, player_max, playtime_min, playtime_max,
	age_min, scoring)

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *GameModel) Get(id int) (Game, error) {
	return Game{}, nil
}

func (m *GameModel) Latest() ([]Game, error) {
	return nil, nil
}