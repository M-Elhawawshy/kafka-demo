package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"kafka-demo/internal/database"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type application struct {
	logger *slog.Logger
	DBPool *pgxpool.Pool
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env")
	}
	pool, err := OpenDB()
	if err != nil {
		log.Fatal("could not open a connection to DB")
	}
	defer pool.Close()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	app := application{
		logger: logger,
		DBPool: pool,
	}

	channel := app.fetchEvents(context.Background())
	app.publishEvents(context.Background(), channel)

	log.Fatal(http.ListenAndServe("localhost:8080", app.routes()))
}

func (app *application) fetchEvents(ctx context.Context) <-chan database.Outbox {
	outChan := make(chan database.Outbox, 100)

	go func() {
		defer close(outChan)

		for {
			// Exit if context is canceled
			if ctx.Err() != nil {
				return
			}

			conn, err := app.DBPool.Acquire(ctx)
			if err != nil {
				app.logger.Error("Could not acquire a conn inside fetchEvents", "err", err)
				time.Sleep(time.Second)
				continue
			}

			q := database.New(conn)

			events, err := q.GetUnpublishedEvents(ctx)
			conn.Release() // make sure it's released even on success/failure

			if err != nil {
				app.logger.Error("Failed to fetch events", "err", err)
				time.Sleep(time.Second)
				continue
			}

			for _, event := range events {
				select {
				case outChan <- event:
					app.logger.Info("event fetched and passed into channel")
				case <-ctx.Done():
					return
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return outChan
}

func (app *application) publishEvents(ctx context.Context, eventsChan <-chan database.Outbox) {
	go func() {
		writer := &kafka.Writer{
			Addr:     kafka.TCP(os.Getenv("KAFKA_URL")),
			Balancer: &kafka.LeastBytes{},
		}
		defer writer.Close()

		for {
			select {
			case <-ctx.Done():
				return

			case event, ok := <-eventsChan:
				if !ok {
					return // channel closed
				}

				err := writer.WriteMessages(ctx, kafka.Message{
					Topic: event.EventType,
					Key:   []byte(event.AggregateID),
					Value: event.Payload,
				})
				if err != nil {
					app.logger.Error("Failed to write the event into Kafka", "event_id", event.ID, "err", err)
					continue
				}

				// Update DB to mark event as published
				DBConn, err := app.DBPool.Acquire(ctx)
				if err != nil {
					app.logger.Error("Could not acquire DB connection", "err", err)
					time.Sleep(time.Second)
					continue
				}

				q := database.New(DBConn)
				err = q.SetEventAsPublished(ctx, event.ID)
				DBConn.Release()

				if err != nil {
					app.logger.Error("Could not mark event as published", "event_id", event.ID, "err", err)
					time.Sleep(time.Second)
				}
			}
		}
	}()
}
