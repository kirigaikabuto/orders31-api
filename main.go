package main

import (
	"github.com/djumanoff/amqp"
	"github.com/kirigaikabuto/common-lib31"
	"github.com/kirigaikabuto/orders31"
	"log"
)

func main() {
	ordersMongoStore, err := orders31.NewOrdersStore(common.MongoConfig{
		Host:           "localhost",
		Port:           "27017",
		Database:       "ivi",
		CollectionName: "orders",
	})
	if err != nil {
		log.Fatal(err)
	}
	ordersAmqpEndpoints := orders31.NewOrdersAmqpEndpoints(ordersMongoStore)
	rabbitConfig := amqp.Config{
		Host:     "localhost",
		Port:     5672,
		LogLevel: 5,
	}
	serverConfig := amqp.ServerConfig{
		ResponseX: "response",
		RequestX:  "request",
	}

	sess := amqp.NewSession(rabbitConfig)
	err = sess.Connect()
	if err != nil {
		panic(err)
		return
	}
	srv, err := sess.Server(serverConfig)
	if err != nil {
		panic(err)
		return
	}
	srv.Endpoint("orders.create", ordersAmqpEndpoints.CreateOrderAmqpEndpoint())
	srv.Endpoint("orders.list", ordersAmqpEndpoints.ListOrderAmqpEndpoint())
	err = srv.Start()
	if err != nil {
		panic(err)
		return
	}
}
