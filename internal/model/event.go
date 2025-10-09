package model

import "time"

type Event struct {
	ID          int          `db:"id" json:"id"`
	Name        string       `db:"name" json:"name"`
	Description string       `db:"description" json:"description"`
	EventBlocks []EventBlock `json:"event_blocks"`
}

type EventBlock struct {
	ID          int       `db:"id" json:"id"`
	EventID     int       `db:"event_id" json:"event_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	StartDate   time.Time `db:"start_date" json:"start_date"`
	EndDate     time.Time `db:"end_date" json:"end_date"`
	Link        string    `db:"link" json:"link"`
}
