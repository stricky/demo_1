package main

import (
	"demo_1/broker"
	"demo_1/db"
	"demo_1/handlers"
	"github.com/nats-io/nats.go"
)

func main() {

	nc := broker.NewBrokerConnection(nats.DefaultURL, "test")
	sdb := db.NewDB("../../golang-sqlite/demo.db")
	sdb.Init()
	defer sdb.Close()
	defer serviceRecover()


	mh := handlers.NewProducerHandler(nc, sdb, 3)
	mh.Run()

}


func serviceRecover(){

}