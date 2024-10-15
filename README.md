# Capstone Project

A Go-based web API with user and post management, built using Chi for routing, and connected to a PostgreSQL database.

## Features

- RESTful API: Build and expose endpoints to interact with the aggregator service.
- Database Integration: Using PostgreSQL for database management with SQLc for SQL code generation and Goose for database migrations.

## Installation and Tools Used

- **[Go](https://golang.org/dl/)**: The primary language for building the API.
- **[PostgreSQL](https://www.postgresql.org/download/)**: The database used to store user and post data.
- **[SQLC](https://github.com/sqlc-dev/sqlc/)**: A Go package to generate type-safe Go code from SQL queries.
- **[Goose](https://github.com/pressly/goose/)**: A tool to manage database migrations.
- **[Chi](https://github.com/go-chi/chi/)**: A lightweight, idiomatic HTTP router for Go.
- **[CORS](https://github.com/go-chi/cors/)**: Middleware to handle cross-origin resource sharing.
- **[godotenv](https://github.com/joho/godotenv/)**: A Go package used to load environment variables from a `.env` file.
