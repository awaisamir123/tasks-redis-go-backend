# Task Management System

This project is a task management system that provides an API for managing tasks. It uses Go for the back-end, Redis for caching, and MongoDB as the database. The front-end communicates with the back-end API to create, update, delete, and fetch tasks. This system also supports bulk task creation.

## Project Structure

- `go-service/`: Contains the Go back-end service.
- `redis/`: Contains Redis setup instructions.
- `frontend/`: Contains the front-end (React) application.
- `mongodb/`: Instructions for setting up MongoDB.
- `.env`: Contains environment variables for configuring the Go back-end.

## Prerequisites

Before you start, ensure that you have the following installed:

- Go (version 1.18+)
- Redis (running locally or in a container)
- MongoDB (running locally or in a container)
- Node.js and npm (if working with the front-end)

## Environment Variables

Create a `.env` file in the `go-service/` directory with the following content:

```env
# MongoDB configuration
MONGO_URI=
MONGO_DB=

# Redis configuration
REDIS_ADDR=

# Server configuration
SERVER_PORT=

# Front-end URL (for CORS configuration)
FRONTEND_URL=


## Running the Go Service

1. Navigate to the go-service directory:
cd go-service

2. Install Go dependencies:
go mod tidy

3. Run the Go service:
go run main.go

The Go server should now be running on http://localhost:8080.


