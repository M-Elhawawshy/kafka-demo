package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func main() {
	topic := "test-topic"
	partition := 0
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		"localhost:9092",
		topic, partition,
	)
	if err != nil {
		log.Fatal("could not initialize kafka producer")
	}

	_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	_, err = conn.WriteMessages(
		kafka.Message{Key: []byte("some key"), Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
