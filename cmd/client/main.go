package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"


	"github.com/66james99/learn-pub-sub-starter/internal/gamelogic"
	"github.com/66james99/learn-pub-sub-starter/internal/pubsub"
	"github.com/66james99/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, fmt.Sprint(routing.PauseKey,".",username), routing.PauseKey, pubsub.Transient)
	if err != nil {
		log.Fatalf("could not declare and bind queue: %v", err)
		return
	}


		// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")

	
}
