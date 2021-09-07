package broker

import (
	"demo_1/models"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type BrokerConnection struct {
	url     string
	subject string

	CaughtMsgC chan *models.IncomingMessage
	Connection *nats.Conn
}

func NewBrokerConnection(url string, subject string) *BrokerConnection {
	return &BrokerConnection{
		url:        url,
		subject:    subject,
		CaughtMsgC: make(chan *models.IncomingMessage),
	}
}

func (bc *BrokerConnection) Listen() error {
	nc, err := nats.Connect(bc.url)
	if err != nil {
		log.Printf("failed to setup connection with broker: %v", err)
		return err
	}
	defer nc.Close()
	_, err = nc.Subscribe(bc.subject, func(msg *nats.Msg) {
		log.Printf("[#%d] Received [%s]: '%s'", 1, msg.Subject, string(msg.Data))
		bc.CaughtMsgC <- models.NewIncomingMessage(time.Now(), msg.Data)
	})
	if err != nil {
		log.Printf("failed to get message: %v", err)
		return err
	}
	return nil
}

func (bc *BrokerConnection) Send(msg *models.OutgoingMessage) error {
	nc, err := nats.Connect(bc.url)
	if err != nil {
		log.Printf("failed to setup connection with broker: %v", err)
		return err
	}
	defer nc.Close()

	bMsg, err := parse(msg)
	log.Printf("Sending message %s", bMsg)

	err = nc.Publish(bc.subject, bMsg)
	if err != nil {
		log.Printf("message %s sent", bMsg)
	}
	return err
}

func parse(msg *models.OutgoingMessage) ([]byte, error) {
	res, err := json.Marshal(msg)
	return res, err
}
