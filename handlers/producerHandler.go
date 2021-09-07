package handlers

import (
	"demo_1/broker"
	"demo_1/db"
	"fmt"
	"log"
	"time"
)

type MessageProducerHandler struct {
	db   *db.DBConnection
	freq time.Duration

	sender *broker.BrokerConnection
	sentDB *db.DBConnection
	errC   chan error
	stopC  chan bool
}

func NewProducerHandler(c *broker.BrokerConnection, db *db.DBConnection, sendFreq time.Duration) *MessageProducerHandler {
	return &MessageProducerHandler{
		sender: c,
		freq:   sendFreq * time.Second,
		sentDB: db,
		errC:   make(chan error),
		stopC:  make(chan bool),
	}
}

func (m *MessageProducerHandler) Run() {
	log.Println("Producer started")
	m.run()
}

func (m *MessageProducerHandler) Stop() {
	log.Println("Producer was stopped")
	m.stopC <- true
}

func (m *MessageProducerHandler) run() {
	ticker := time.NewTicker(m.freq)
	defer func() {
		ticker.Stop()
	}()

	defer func() {
		if recoveryMessage := recover(); recoveryMessage != nil {
			fmt.Println(recoveryMessage)
		}
	}()
	for {
		select {
		case <-ticker.C:
			go m.SendMessage()
		case err := <-m.errC:
			log.Println("Caught error:", err)
			return
		case <-m.stopC:
			return
		}
	}
}

func (m *MessageProducerHandler) SendMessage() {
	msg, err := m.sentDB.GetOutMessageFromQueue()
	if err != nil {
		m.errC <- fmt.Errorf("failed to get message from queue: %v", err)
		return
	}

	log.Printf ("try to send message %v", msg)
	uid, err := m.sentDB.RecordOutMessage(msg.BodyText, msg.Hash)
	if err != nil {
		m.errC <- fmt.Errorf("failed to record message in stat: %v", err)
		return
	}
	err = m.sender.Send(msg)
	if err != nil {
		err = m.sentDB.UpdateOutMessageStatus(msg.Uid, db.ErrorStatus)
		m.errC <- fmt.Errorf("failed send message: %v", err)
		return
	}
	err = m.sentDB.DeleteFromQueue(msg.Uid)
	if err != nil {
		m.errC <- fmt.Errorf("failed remove message %v from queue: %v", msg.Uid, err)
		return
	}
	err = m.sentDB.UpdateOutMessageStatus(uid, db.SentStatus)
	if err != nil {
		m.errC <- fmt.Errorf("failed update stat message %v from queue: %v", msg.Uid, err)
		return
	}
}
