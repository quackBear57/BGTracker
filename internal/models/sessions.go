package models

import (
	"database/sql"
	// "errors"
	"time"
	"encoding/json"
	"fmt"
)

type Session struct {
	ID				int
	Date			time.Time
	Game_ID			int
	Player_Scores	[]Player_Score
	Winner_ID		int
}

type Player_Score struct {
	Player_ID	int
	Score		int
}

type SessionModel struct {
	DB				*sql.DB
}

func (m *SessionModel) Insert(gameID int,
	date time.Time,
	playerScores []Player_Score,
	winnerID int) (int, error) {
	
	stmt := `INSERT INTO play_sessions (date, game_id, winner_id, scores) 
	VALUES(?, ?, ?, ?)`

	// fmt.Printf("Game ID: %d\n", gameID)
	// fmt.Printf("Date: %s\n", date.Format("02-Jan-2006 15:04"))
	// fmt.Printf("Winner ID: %d\n", winnerID)

	// fmt.Println("Json Data:")
	jsonBytes, _ := json.Marshal(playerScores)
	fmt.Println(string(jsonBytes))


	for _, playerScore := range playerScores{
		fmt.Printf("Player: %d // Score: %d", playerScore.Player_ID, playerScore.Score)
	}

	// need to store scores as json??

	result, err := m.DB.Exec(stmt, date, gameID, playerScores, winnerID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return 0, nil
}

// func (m *GameModel) Get(id int) (Game, error) {
// 	stmt := `SELECT id, title, bgg_id, price, added, state, year_pub,
// 	description, player_min, player_max, playtime_min, playtime_max,
// 	age_min, scoring FROM games WHERE id = ?`

// 	// Creating temporary variable to handle NULL values for the description
// 	// field
// 	var desc_nullstr sql.NullString

// 	row := m.DB.QueryRow(stmt, id)

// 	var g Game

// 	err := row.Scan(&g.ID, &g.Title, &g.Bgg_ID, &g.Price, &g.Added, &g.State,
// 	&g.Year_pub, &desc_nullstr, &g.Player_min, &g.Player_max, &g.Playtime_min,
// 	&g.Playtime_max, &g.Age_min, &g.Scoring)

// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return Game{}, ErrNoRecord
// 		default:
// 			return Game{}, err
// 		}
// 	}

// 	// Convert temporary var to the g description field if not NULL
// 	if desc_nullstr.Valid {
// 		g.Description = desc_nullstr.String
// 	}

// 	return g, nil
// }

// func (m *GameModel) Latest() ([]Game, error) {
// 	stmt := `SELECT id, title, bgg_id, price, added, state, year_pub,
// 	description, player_min, player_max, playtime_min, playtime_max,
// 	age_min, scoring FROM games WHERE added < DATE_ADD(UTC_TIMESTAMP(), 
// 	INTERVAL ? DAY) ORDER BY id DESC LIMIT 10`

// 	rows, err := m.DB.Query(stmt, 30)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	var games []Game

// 	for rows.Next() {
// 		var g Game
// 		var desc_nullstr sql.NullString
		
// 		err = rows.Scan(&g.ID, &g.Title, &g.Bgg_ID, &g.Price, &g.Added, &g.State,
// 		&g.Year_pub, &desc_nullstr, &g.Player_min, &g.Player_max, &g.Playtime_min,
// 		&g.Playtime_max, &g.Age_min, &g.Scoring)
		
// 		if err != nil {
// 			return nil, err
// 		}

// 		if desc_nullstr.Valid {
// 			g.Description = desc_nullstr.String
// 		}

// 		games = append(games, g)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}
	
// 	return games, nil
// }

func (m *SessionModel) GetList(limit int) ([]Session, error) {
	stmt := `SELECT id, date, game_id, winner_id, scores
	FROM play_sessions ORDER BY id DESC LIMIT ?`

	rows, err := m.DB.Query(stmt, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sessions []Session

	for rows.Next() {
		var s Session


		// Need to check how to properly unmarshal player scores from json to struct
		err = rows.Scan(&s.ID, &s.Date, &s.Game_ID, &s.Winner_ID, &s.Player_Scores)
		
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return sessions, nil
}