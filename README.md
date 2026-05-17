# Simple Whalebone Microservice

This application is a simple microservice written in Go, providing two main user-related endpoints.

## Prerequisites

Make sure you have the following installed on your system:
*   [Docker](https://docs.docker.com/get-docker/) and Docker Compose
*   [Go 1.25+](https://go.dev/doc/install) (only if running or building locally)
*   `make` utility

## How to Run 🚀

Use the `Makefile` to easily run and manage the application. It provides several targets:

### Running via Docker 🐋
Run following to build the image and start the application along with its dependencies (e.g., **PostgreSQL**) in the background:
```bash
make run
```

To stop the application and remove the containers:
```bash
make stop
```

### Local Development 🛠️
If you prefer to run the application directly on your host machine without Docker:
```bash
# Build the executable binary
make build-local

# Run the application directly
make run-local
```

### Quality Assurance 🛠️
To ensure code quality, run the test suite and linter:
```bash
# Run all unit and integration tests
make test

# Run the linter to automatically fix common issues
make lint
```

## API Endpoints 
The service exposes the following main routes under the `/api/v1` prefix:
- `GET /api/v1/users/{id}` - Retrieves a specific user by their (system) ID.
- `POST /api/v1/users/save` - Creates and persists a new user in the system.
