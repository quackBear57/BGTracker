package main

import(
	"bgtracker.ndr54.com/internal/models"
)

type templateData struct {
	Game			models.Game
	Games			[]models.Game
}