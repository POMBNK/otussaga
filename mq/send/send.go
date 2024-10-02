package main

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math/rand"
	"time"
)

type AddTask struct {
	Number1 int
	Number2 int
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

func main() {
	conn, err := amqp.Dial("amqp://rmuser:rmpassword@localhost:5672/")
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	//queue, err := amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	err = amqpChannel.ExchangeDeclare(
		"add",    // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	handleError(err, "Could not declare `add` queue")

	err = amqpChannel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	handleError(err, "Could not declare `add` queue")

	rand.Seed(time.Now().UnixNano())
	addTask := AddTask{Number1: rand.Intn(999), Number2: rand.Intn(999)}
	body, err := json.Marshal(addTask)
	if err != nil {
		handleError(err, "Error encoding JSON")
	}

	err = amqpChannel.Publish("add", "", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	log.Printf("AddTask: %d+%d", addTask.Number1, addTask.Number2)

}
