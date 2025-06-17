# Kafka Outbox Demo

This project demonstrates the **Transactional Outbox Pattern** using **Go**, **PostgreSQL**, and **Kafka**. It simulates a basic microservice architecture where events are published to Kafka **only after the related database transaction succeeds**, ensuring consistency between your database and your message broker.

### 🔧 Key Features

- Atomic event publishing via transactional outbox table
- Kafka producer that polls unpublished events and sends them to a Kafka topic
- Kafka consumer that processes events and applies deduplication logic
- Idempotent message handling using a processed-events table
- Written in Go using `pgx`, `kafka-go`, and `slog` for structured logging

### 📦 Stack

- Go
- PostgreSQL
- Apache Kafka
- Segment's `kafka-go` client
- pgx / pgxpool for DB access
- Standard Go net/http

### 💡 Why This Matters

In distributed systems, it's crucial to ensure that data changes and event emissions are in sync. The outbox pattern prevents lost messages or data inconsistencies by coupling event creation with the database transaction — and this project shows how to build that correctly with Go.

### 📂 Project Structure

- `/consumer/` – Contains consumer logic
- `/producer/` – Contains producer logic
- `/internal/database/` – SQLC-generated queries (optional)
- `.env` – Environment configuration for Kafka and Postgres

### 📬 Running the Project

1. Start Kafka and PostgreSQL
2. Set up environment variables (`KAFKA_URL`, `DB_URL`)
3. Run the producer and consumer and then hit the `/test` endpoint
4. Events will be stored in the outbox, published to Kafka, and processed with deduplication

