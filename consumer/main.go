package main

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"kafka-demo/internal/database"
	"log"
	"log/slog"
	"os"
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
	inChan := app.consumeEvents(context.Background())
	app.processEvent(context.Background(), inChan)

}

func (app *application) consumeEvents(ctx context.Context) <-chan database.Test {
	outChan := make(chan database.Test, 100)
	go func() {
		defer close(outChan)
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{os.Getenv("KAFKA_URL")},
			Topic:    "test.created",
			MaxBytes: 10e6,
		})
		defer r.Close()
		for {
			if ctx.Err() != nil {
				return
			}
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				app.logger.Error("Could not read", "err", err)
				break
			}
			var test database.Test
			err = json.Unmarshal(m.Value, &test)
			if err != nil {
				app.logger.Error("Could not unmarshal", "err", err)
				break
			}
			select {
			case outChan <- test:
				app.logger.Debug("Event consumed and sent to be processed")
			case <-ctx.Done():
				return
			}
		}
	}()
	return outChan
}

func (app *application) processEvent(ctx context.Context, inChan <-chan database.Test) {
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-inChan:
			// do something with t
			conn, err := app.DBPool.Acquire(ctx)
			if err != nil {
				app.logger.Error("could not acquire conn", "err", err)
				continue
			}
			q := database.New(conn)
			ok, err := q.IsProcessed(ctx, t.ID)
			if err != nil {
				app.logger.Error("could not query for status of processed", "err", err)
				continue
			}

			if !ok {
				app.logger.Debug("Processing event now!", "event id", t.ID.String())
				err = q.InsertProcessed(ctx, database.InsertProcessedParams{
					ID:      t.ID,
					Content: t.Content,
				})
				if err != nil {
					app.logger.Error("could not inser processed into db", "error", err)
					continue
				}
			}
			conn.Release()
		}
	}
}
