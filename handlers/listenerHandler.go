package handlers

import (
	"demo_1/broker"
	"demo_1/db"
	"demo_1/models"
	"fmt"
	"log"
	"time"
)

type ListenerMessageHandler struct {
	receiver *broker.BrokerConnection
	dbc      *db.DBConnection
	stopC    chan bool
	errorC   chan error
}

func NewListenerMessageHandler(bc *broker.BrokerConnection, dbConnection *db.DBConnection) *ListenerMessageHandler {
	return &ListenerMessageHandler{
		receiver: bc,
		dbc:      dbConnection,
		stopC:    make(chan bool),
		errorC:   make(chan error),
	}
}

func (rm *ListenerMessageHandler) Run() {
	log.Println("Listener started")
	for {
		select {
		case msg := <-rm.receiver.CaughtMsgC:
			go rm.processMessage(msg)
		case err := <-rm.errorC:
			log.Printf("error while receiving %v", err)
		case <-rm.stopC:
			return
		}
	}
}

func (rm *ListenerMessageHandler) Stop() {
	rm.stopC <- true
}

func (rm *ListenerMessageHandler) processMessage(message *models.IncomingMessage) {
	fmt.Printf("Incoming message %s is recording", message)
	uid, err := rm.dbc.RecordInMessage(message.Body, time.Now())
	if err != nil {
		rm.errorC <- fmt.Errorf("failed record %v", err)
		return
	} else {
		fmt.Printf("Registered new message uid: %v", uid)
	}
}
