package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/kubotak-is/rabbit_mq_sample/protocol"
	"github.com/streadway/amqp"
)

var (
	queueName = flag.String("queueName", "test_queue", "Your RabbtMQ Queue")
)

func main() {
	flag.Parse()

	if *queueName == "" {
		log.Fatalln("[ERROR] require queueName")
	}

	log.Println("consumer start " + *queueName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	conn, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume(*queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return
	}

	go func() {
	CONSUMER_FOR:
		for {
			select {
			case <-ctx.Done():
				break CONSUMER_FOR
			case m, ok := <-msgs:
				if ok {
					var p protocol.Protocol
					if err := json.Unmarshal(m.Body, &p); err != nil {
						log.Printf("[ERROR] %s", err.Error())
						continue CONSUMER_FOR
					}

					log.Printf("[INFO] success consumed. tag: %d, msg: %v", m.DeliveryTag, &p)
				}
			}
		}
	}()

	<-signals

	log.Println("consumer stop")
}
