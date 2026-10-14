package models

import (
	"database/sql"
	"errors"
	"time"
)

type  Game struct {
	ID				int
	Title			string
	Bgg_ID			int
	Price			float64
	Added			time.Time
	State			string
	Year_pub		int
	Description		string
	Player_min		int
	Player_max		int
	Playtime_min	int
	Playtime_max	int
	Age_min			int
	Scoring			bool
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
	age_min, scoring) VALUES(?, ?, ?, UTC_TIMESTAMP(), ?, ?,
	?, ?, ?, ?, ?,
	?, ?)`

	result, err := m.DB.Exec(stmt, title,  bgg_id, price, state, year_pub,
	description, player_min, player_max, playtime_min, playtime_max,
	age_min, scoring)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *GameModel) Get(id int) (Game, error) {
	stmt := `SELECT id, title, bgg_id, price, added, state, year_pub,
	description, player_min, player_max, playtime_min, playtime_max,
	age_min, scoring FROM games WHERE id = ?`

	// Creating temporary variable to handle NULL values for the description
	// field
	var desc_nullstr sql.NullString

	row := m.DB.QueryRow(stmt, id)

	var g Game

	err := row.Scan(&g.ID, &g.Title, &g.Bgg_ID, &g.Price, &g.Added, &g.State,
	&g.Year_pub, &desc_nullstr, &g.Player_min, &g.Player_max, &g.Playtime_min,
	&g.Playtime_max, &g.Age_min, &g.Scoring)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Game{}, ErrNoRecord
		default:
			return Game{}, err
		}
	}

	// Convert temporary var to the g description field if not NULL
	if desc_nullstr.Valid {
		g.Description = desc_nullstr.String
	}

	return g, nil
}

func (m *GameModel) Latest() ([]Game, error) {
	stmt := `SELECT id, title, bgg_id, price, added, state, year_pub,
	description, player_min, player_max, playtime_min, playtime_max,
	age_min, scoring FROM games WHERE added < DATE_ADD(UTC_TIMESTAMP(), 
	INTERVAL ? DAY) ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt, 30)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []Game

	for rows.Next() {
		var g Game
		var desc_nullstr sql.NullString
		
		err = rows.Scan(&g.ID, &g.Title, &g.Bgg_ID, &g.Price, &g.Added, &g.State,
		&g.Year_pub, &desc_nullstr, &g.Player_min, &g.Player_max, &g.Playtime_min,
		&g.Playtime_max, &g.Age_min, &g.Scoring)
		
		if err != nil {
			return nil, err
		}

		if desc_nullstr.Valid {
			g.Description = desc_nullstr.String
		}

		games = append(games, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return games, nil
}