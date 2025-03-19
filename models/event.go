package models

import (
	"time"

	"example.com/rest-api/db"
)

type Event struct {
	ID       int64
	Name     string    `binding:"required"`
	Descr    string    `binding:"required"` // description of an event
	Local    string    `binding:"required"` // place, where event should take place
	DateTime time.Time `binding:"required"`
	UserID   int64
}

// var events = []Event{}

func (e *Event) Save() error {
	query := `
	INSERT INTO events(name, descr, local, dateTime, user_id)
	VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Descr, e.Local, e.DateTime, e.UserID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	e.ID = id
	return err
}

func GetAllEvets() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Descr, &event.Local, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Descr, &event.Local, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (event Event) Update() error {
	query := `
	UPDATE events
	SET name = ?, descr = ?, local = ?, dateTime = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Descr, event.Local, event.DateTime, event.ID)
	return err
}

func (event Event) Delete() error {
	query := `
	DELETE FROM events
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	return err
}

func (event Event) Register(userId int64) error {
	query := `
	INSERT INTO registrations(event_id, user_id)
	VALUES (?, ?)
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID, userId)
	return err
}

func (event Event) CancelRegistration(userId int64) error {
	query := `
	DELETE FROM registrations
	WHERE event_id = ?
	AND user_id = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID, userId)
	return err
}
