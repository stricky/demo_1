package db

import (
	"database/sql"
	"demo_1/models"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type Status string

const (
	SentStatus  Status = "sent"
	ErrorStatus Status = "failed"
	NewStatus   Status = "new"
)

type DBConnection struct {
	connection *sql.DB
	dbPath     string
}

func NewDB(dbPath string) *DBConnection {
	return &DBConnection{
		dbPath: dbPath,
	}
}

func (db *DBConnection) Init() {
	dbC, err := sql.Open("sqlite3", db.dbPath)
	if err != nil {
		log.Println(err)
	}

	db.connection = dbC
}

func (db *DBConnection) Close() error {
	return db.connection.Close()
}

func (db *DBConnection) RecordOutMessage(bodyText string, hash string) (int64, error) {
	fmt.Println()
	stmt, err := db.connection.Prepare("INSERT INTO messages_stat (body, hash, status, timestamp ) VALUES (?,?,?,datetime('now'))")
	if err != nil {
		return 0, fmt.Errorf("failed build request: %v", err)
	}
	res, err := stmt.Exec(bodyText, hash, NewStatus)
	if err != nil {
		return 0, fmt.Errorf("failed exec request: %v", err)
	}
	uid, err := res.LastInsertId()
	return uid, err
}

func (db *DBConnection) UpdateOutMessageStatus(uid int64, status Status) error {
	stmt, err := db.connection.Prepare("UPDATE messages_stat SET status=? WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(status, uid)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBConnection) GetOutMessageFromQueue() (*models.OutgoingMessage, error) {

	var msg models.OutgoingMessage
	stmt, err := db.connection.Query(
		"SELECT * FROM messages_queue WHERE id= (SELECT min(id) from messages_queue)")
	if err != nil {
		return nil, err
	}

	var t string
	defer stmt.Close()
	stmt.Next()

	err = stmt.Scan(&msg.Uid, &t, &msg.BodyText, &msg.Hash)
	log.Printf("selected message from queue %s", msg)
	if msg.Uid == 0 {
		return nil, fmt.Errorf("no more messages in queue")
	}
	msg.Time, err = time.Parse("2006-01-02 15:04:05", t)

	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (db *DBConnection) DeleteFromQueue(uid int64) error {
	stmt, err := db.connection.Prepare("DELETE FROM messages_queue WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(uid)
	return err
}

func (db *DBConnection) RecordInMessage(body []byte, t time.Time) (int64, error) {
	stmt, err := db.connection.Prepare("INSERT INTO incoming_messages( timestamp, body ) VALUES (datetime('now'),?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(t.String(), string(body))
	if err != nil {
		return 0, err
	}
	uid, err := res.LastInsertId()
	return uid, err
}
