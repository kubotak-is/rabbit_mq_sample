package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/kubotak-is/rabbit_mq_sample/protocol"
	"github.com/streadway/amqp"
)

func main() {
	log.Println("producer start")

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

	for i := 0; i < 3; i++ {
		p := &protocol.Protocol{
			Message:   fmt.Sprintf("Hello. No%d", i),
			Timestamp: time.Now().UnixNano(),
		}
		bytes, err := json.Marshal(p)
		if err != nil {
			log.Printf("[ERROR] %s", err.Error())
			continue
		}
		if err := ch.Publish("test", "", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		}); err != nil {
			log.Printf("[ERROR] %s", err.Error())
			continue
		}

		log.Printf("[INFO] send message. msg: %v", p)
	}

	log.Println("producer stop")
}
