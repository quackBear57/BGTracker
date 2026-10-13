package models

import (
	"database/sql"
	"time"
)

type  Game struct {
	ID				int
	Title			string
	Bgg_ID			int
	price			float
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
	price float,
	state string,
	year_pub int,
	description string,
	player_min int,
	player_max int,
	playtime_min int,
	playtime_max int,
	scoring bool) (int, error) {
	
	return 0, nil
}

func (m *GameModel) Get(id int) (Game, error) {
	return Game{}, nil
}

func (m *GameModel) Latest() ([]Game, error) {
	return nil, nil
}