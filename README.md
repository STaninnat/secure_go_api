# Capstone Project

![code coverage badge](https://github.com/STaninnat/capstone_project/actions/workflows/ci.yml/badge.svg)
![code coverage badge](https://github.com/STaninnat/capstone_project/actions/workflows/cd.yml/badge.svg)

A Go-based web API with user and post management, built using Chi for routing, and connected to a Turso database.

## Features

- **Environment Configuration**: The application loads environment variables from a `.env` file for configuration, including database connection URL, JWT secrets, and refresh tokens.
- **User Authentication**:
  - JWT-based authentication for secure API access.
  - Login, logout, and token refresh endpoints.
  - Custom middleware for authenticated routes.
- **Database Integration**:
  - Connects to a database using SQLc-generated queries for type-safe SQL code.
  - Supports CRUD operations for user and post management when the database URL is provided.
  - Gracefully handles situations when the database is not configured, disabling CRUD functionality.
- **CORS Middleware**: Configures CORS to allow cross-origin requests from all origins, enabling the API to be accessed from various clients.
- **Static File Handling**: Serves static files like HTML, CSS, JavaScript, and JSON from a `/static/` endpoint.
- **Health Check**: Exposes a `/healthz` endpoint for readiness checks to ensure the API is functioning.
- **Error Handling**: A custom `/err` endpoint to simulate errors for testing purposes.

## Installation and Tools Used

- **[Go](https://golang.org/dl/)**: The primary language for building the API.
- **[SQLC](https://github.com/sqlc-dev/sqlc/)**: A Go package to generate type-safe Go code from SQL queries.
- **[Goose](https://github.com/pressly/goose/)**: A tool to manage database migrations.
- **[Chi](https://github.com/go-chi/chi/)**: A lightweight, idiomatic HTTP router for Go.
- **[CORS](https://github.com/go-chi/cors/)**: Middleware to handle cross-origin resource sharing.
- **[godotenv](https://github.com/joho/godotenv/)**: A Go package used to load environment variables from a `.env` file.
- **[golang-jwt/jwt](https://github.com/golang-jwt/jwt)**: A library for working with JSON Web Tokens (JWT) for authentication.
- **[google/uuid](https://github.com/google/uuid)**: A package to generate and handle UUIDs.
- **[golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)**: A collection of cryptographic algorithms and utilities for Go.

## Local Development

Create a `.env` file in the root of the project with the following contents:

```bash
PORT="8080"
```

```bash
JWT_SECRET="YOUR JWT SECRET"
```

```bash
REFRESH_SECRET="YOUR REFRESH SECRET"
```

Run the server:

```bash
go build -o out && ./out
```

## Installation

1. Forking and clone the repository.
2. Configure your environment variables.
3. Set up database with **[Turso](https://turso.tech/)** *or use other databases, but you may need to change some things in the code.*
4. Run the necessary database migrations with Goose.
5. Start the local service worker and API using Go.

```bash
go build -o out && ./out
```

## Note

Finally, I create a new Docker image from this project to google artifact registry and deploy to cloud run revision with the new image and serves the app to the public internet. This works with continuous integration(CI) and continuous deployment(CD). And I choose to use Turso database because turso is a cloud provider that specializes in hosting serverless SQLite-like databases. It's a great fit for this project because it has a very generous free tier, and it's easy to use.
