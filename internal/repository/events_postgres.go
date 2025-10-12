package repository

import (
	"fmt"
	"strings"
	"time"

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

func (r *EventsPostgres) DeleteEventBlock(blockId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", eventBlocksTable)
	_, err := r.db.Exec(query, blockId)
	return err
}

func (r *EventsPostgres) DeleteEvent(eventId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE event_id = $1", eventBlocksTable)
	_, err = tx.Exec(query, eventId)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE id = $1", eventsTable)
	_, err = tx.Exec(query, eventId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *EventsPostgres) EditEventInfo(event model.Event) error {
	query := fmt.Sprintf("UPDATE %s SET name = $1, description = $2 WHERE id = $3", eventsTable)
	_, err := r.db.Exec(query, event.Name, event.Description, event.ID)
	return err
}

func (r *EventsPostgres) EditBlockInfo(block model.EventBlock) error {
	query := fmt.Sprintf("UPDATE %s SET name = $1, description = $2, start_date = $3, end_date = $4, link = $5 WHERE id = $6", eventBlocksTable)
	_, err := r.db.Exec(query, block.Name, block.Description, block.StartDate, block.EndDate, block.Link, block.ID)
	return err
}

func (r *EventsPostgres) GetCurrentEvents() ([]model.Event, error) {
	currentTime := time.Now()
	query := fmt.Sprintf(`
		SELECT DISTINCT e.id, e.name, e.description
		FROM %s e
		JOIN %s b ON e.id = b.event_id
		WHERE b.start_date <= $1 AND b.end_date >= $1
	`, eventsTable, eventBlocksTable)

	var events []model.Event
	if err := r.db.Select(&events, query, currentTime); err != nil {
		return nil, err
	}
	return events, nil
}

type eventWithBlocks struct {
	ID          int
	Name        string
	Description string
	Blocks      []model.EventBlock
}

func (r *EventsPostgres) GetAllEvents() ([]model.Event, error) {
	rows, err := r.db.Queryx(fmt.Sprintf(`
		SELECT 
			e.id as event_id, 
			e.name as event_name, 
			e.description as event_description, 
			b.id as block_id, 
			b.name as block_name, 
			b.description as block_description, 
			b.start_date as block_start_date, 
			b.end_date as block_end_date, 
			b.link as block_link
		FROM %s e
		LEFT JOIN %s b ON e.id = b.event_id
		ORDER BY e.id, b.id
	`, eventsTable, eventBlocksTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventsMap := make(map[int]*eventWithBlocks)

	for rows.Next() {
		var (
			eventID          int
			eventName        string
			eventDescription string
			blockID          *int
			blockName        *string
			blockDescription *string
			blockStartDate   *time.Time
			blockEndDate     *time.Time
			blockLink        *string
		)

		err := rows.Scan(
			&eventID,
			&eventName,
			&eventDescription,
			&blockID,
			&blockName,
			&blockDescription,
			&blockStartDate,
			&blockEndDate,
			&blockLink,
		)
		if err != nil {
			return nil, err
		}

		evt, ok := eventsMap[eventID]
		if !ok {
			evt = &eventWithBlocks{
				ID:          eventID,
				Name:        eventName,
				Description: eventDescription,
				Blocks:      []model.EventBlock{},
			}
			eventsMap[eventID] = evt
		}

		if blockID != nil {
			block := model.EventBlock{
				ID:          *blockID,
				EventID:     eventID,
				Name:        "",
				Description: "",
				StartDate:   time.Time{},
				EndDate:     time.Time{},
				Link:        "",
			}
			if blockName != nil {
				block.Name = *blockName
			}
			if blockDescription != nil {
				block.Description = *blockDescription
			}
			if blockStartDate != nil {
				block.StartDate = *blockStartDate
			}
			if blockEndDate != nil {
				block.EndDate = *blockEndDate
			}
			if blockLink != nil {
				block.Link = *blockLink
			}
			evt.Blocks = append(evt.Blocks, block)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	events := make([]model.Event, 0, len(eventsMap))
	for _, evt := range eventsMap {
		e := model.Event{
			ID:          evt.ID,
			Name:        evt.Name,
			Description: evt.Description,
			EventBlocks: evt.Blocks,
		}
		events = append(events, e)
	}

	return events, nil

}

func (r *EventsPostgres) GetOneEvent(eventId int) (model.Event, error) {
	var event model.Event
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", eventsTable)
	if err := r.db.Get(&event, query, eventId); err != nil {
		return model.Event{}, err
	}
	var eventBlocks []model.EventBlock
	query = fmt.Sprintf("SELECT * FROM %s WHERE event_id = $1", eventBlocksTable)
	if err := r.db.Select(&eventBlocks, query, eventId); err != nil {
		return model.Event{}, err
	}
	event.EventBlocks = eventBlocks
	return event, nil
}

func (r *EventsPostgres) GetOneBlock(blockId int) (model.EventBlock, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", eventBlocksTable)
	var block model.EventBlock
	if err := r.db.Get(&block, query, blockId); err != nil {
		return model.EventBlock{}, err
	}
	return block, nil
}
