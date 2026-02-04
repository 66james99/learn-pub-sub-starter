package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/66james99/learn-pub-sub-starter/internal/pubsub"
	"github.com/66james99/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not open RabbitMQ channel: %v", err)
	}
	defer channel.Close()

	var message routing.PlayingState
	message.IsPaused = true

	err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, message)
	if err != nil {
		log.Fatalf("could not publish message: %v", err)
	}
	fmt.Println("Message published successfully!")

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}
