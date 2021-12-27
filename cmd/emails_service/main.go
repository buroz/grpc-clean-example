package main

import (
	"encoding/json"
	"log"

	"github.com/buroz/grpc-clean-example/pkg/amqp"

	"github.com/buroz/grpc-clean-example/pkg/config"
	"github.com/buroz/grpc-clean-example/pkg/smtp"
	"github.com/buroz/grpc-clean-example/pkg/utils"
)

var (
	exchangeName = "EMAILS"
	queueName    = "EMAIL_SEND"
)

type emailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func main() {
	forever := make(chan bool)

	conf, err := config.NewConfig()
	utils.Check(err)

	amqpClient := amqp.NewAmqpClient(&conf.Amqp)

	err = amqpClient.Connect()
	utils.Check(err)

	smtpConn := smtp.NewSMTPClient(&conf.Smtp)

	err = smtpConn.Connect()
	utils.Check(err)

	err = amqpClient.Channel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	utils.Check(err)

	q, err := amqpClient.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	utils.Check(err)

	err = amqpClient.Channel.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // args
	)
	utils.Check(err)

	msgs, err := amqpClient.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.Check(err)

	go func() {
		for d := range msgs {
			msg := emailMessage{}
			json.Unmarshal(d.Body, &msg)

			err := smtpConn.Send(msg.To, msg.Subject, msg.Body)
			if err != nil {
				log.Println("ERROR", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
