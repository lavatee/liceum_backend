package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/liceum_backend/internal/model"
)

type EventsPostgres struct {
	db *sqlx.DB
}

func NewEventsPostgres(db *sqlx.DB) *EventsPostgres {
	return &EventsPostgres{
		db: db,
	}
}

func (r *EventsPostgres) CreateEvent(event model.Event) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, description) VALUES ($1, $2) RETURNING id", eventsTable)
	var id int
	row := r.db.QueryRow(query, event.Name, event.Description)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *EventsPostgres) CreateEventBlocks(blocks []model.EventBlock) error {
	query := fmt.Sprintf("INSERT INTO %s (event_id, name, description, link, start_date, end_date) VALUES ", eventBlocksTable)
	queryPieces := make([]string, len(blocks))
	argsCounter := 0
	argsArr := make([]interface{}, 0)
	for i, block := range blocks {
		queryPieces[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", argsCounter+1, argsCounter+2, argsCounter+3, argsCounter+4, argsCounter+5, argsCounter+6)
		argsCounter += 6
		argsArr = append(argsArr, block.EventID, block.Name, block.Description, block.Link, block.StartDate, block.EndDate)
	}
	query += strings.Join(queryPieces, ", ")
	_, err := r.db.Exec(query, argsArr...)
	return err
}
