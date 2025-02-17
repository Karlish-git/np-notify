package models

import "time"

type Player struct {
	UserID          string       `json:"user_id"`
	UserEmail       string       `json:"user_email"`
	Alias           string       `json:"alias"`
	DollarsPaid     int          `json:"dollars_paid"`
	Credits         int          `json:"credits"`
	SubscribedUntil time.Time    `json:"subscribed_until"`
	Created         time.Time    `json:"created"`
	Nagged          time.Time    `json:"nagged"`
	Badges          string       `json:"badges"`
	Score           int          `json:"score"`
	Karma           int          `json:"karma"`
	GamesWon        int          `json:"games_won"`
	GamesIn         int          `json:"games_in"`
	Emailing        bool         `json:"emailing"`
	HasPassword     bool         `json:"has_password"`
	Verified        bool         `json:"verified"`
	OpenGames       []PlayerGame `json:"open_games"`
	CompleteGames   []PlayerGame `json:"complete_games"`
}
