package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"kafka-demo/internal/database"
	"net/http"
)

func (app *application) testHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := app.DBPool.Acquire(r.Context())
	if err != nil {
		http.Error(w, fmt.Errorf("conn from pool error").Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	tx, err := conn.Begin(r.Context())
	if err != nil {
		http.Error(w, fmt.Errorf("transaction begin error").Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())
	qtx := database.New(tx)
	test, err := qtx.InsertTest(r.Context(), database.InsertTestParams{
		ID: uuid.New(),
		Content: pgtype.Text{
			String: "A very unique entry string",
			Valid:  true,
		},
	})
	if err != nil {
		http.Error(w, "Failed to create a test entry", http.StatusInternalServerError)
		return
	}
	testBytes, err := json.Marshal(test)
	if err != nil {
		http.Error(w, "failed to serialize test", http.StatusInternalServerError)
		return
	}

	err = qtx.InsertEvent(r.Context(), database.InsertEventParams{
		ID:            uuid.New(),
		EventType:     "test.created",
		AggregateType: "test",
		AggregateID:   test.ID.String(),
		Payload:       testBytes,
	})
	if err != nil {
		http.Error(w, "failed to insert an event to db", http.StatusInternalServerError)
		return
	}
	_ = tx.Commit(r.Context())
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Test event created!"))
}
