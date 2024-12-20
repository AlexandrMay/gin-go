package models

import (
	"fmt"
	"gin-go/db"
	"time"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

func (e *Event) Save() error {
	query := `INSERT INTO events(name, description, location, dateTime, user_id) 
	VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("Error addding a statement: %v", err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return fmt.Errorf("Error addding event into db: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("Error getting last inserted id: %v", err)
	}
	e.ID = id
	return nil
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("No events was found: %v", err)
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, fmt.Errorf("Error scan rows: %v", err)
		}
		events = append(events, event)
	}
	return events, nil
}

func GetEventById(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, fmt.Errorf("Error scan rows: %v", err)
	}
	return &event, nil
}

func (e Event) Update() error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare a query statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	if err != nil {
		return fmt.Errorf("failed to execute update statement: %w", err)
	}
	return nil
}

func (e Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare a query statement: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.ID)
	if err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}
	return nil
}

func (e Event) Register(userId int64) error {
	query := "INSERT INTO registrations(event_id, user_id) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare a query statement: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.ID, userId)
	if err != nil {
		return fmt.Errorf("failed to execute register statement: %w", err)
	}
	return nil
}

func (e Event) CancelRegistration(userId int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare a query statement: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.ID, userId)
	if err != nil {
		return fmt.Errorf("failed to execute cancel registration statement: %w", err)
	}
	return nil
}
