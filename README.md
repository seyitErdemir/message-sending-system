# Message Sending System

A modern messaging system API developed with Go. This system sends and manages messages through scheduled tasks (cron jobs). Features:

- RESTful API endpoints
- Automatic API documentation with Swagger
- Redis caching
- MySQL database support
- Scheduled message delivery
- Easy setup and deployment with Docker
- Database management with PHPMyAdmin and Redis Stack UI

## Requirements 🛠

- Go 1.16+
- Docker & Docker Compose

## Installation and Running 🚀

### Running with Docker (Recommended)
```bash
# Clone the project
git clone <repository-url>
cd message-sending-system

# Start the application
docker-compose up --build
```

### Local Development
```bash
# Install dependencies
go mod tidy

# Run with hot-reload
air
```

## API Documentation 📚

### Swagger
```bash
# Update Swagger documentation
swag init -g cmd/api/main.go
```

### API Endpoints

#### Message Operations
- `POST /api/messages` - Create new message
- `GET /api/messages` - List sent messages

#### Cron Operations
- `POST /cron/start` - Start message sending cron job
- `POST /cron/stop` - Stop cron job
- `GET /cron/status` - Check cron job status
- `GET /cron/logs` - View cron logs

## Management Interfaces 🖥

### API Documentation
Swagger UI: http://localhost:3000/swagger/index.html

### Database Management
PHPMyAdmin: http://localhost:8080

### Redis Management
Redis Stack UI: http://localhost:8001

## Project Structure 📁

```
.
├── cmd/api/      # Main application
├── pkg/
│   ├── handlers/ # HTTP handlers
│   ├── models/   # Data models
│   ├── database/ # Database operations
│   └── cron/     # Scheduled tasks
└── docker/       # Docker configurations
```

## Environment Variables 🔧

The project uses the following environment variable files:
- `.env` - Main environment variables
- `.env.test` - Variables for test environment

All environment variables are shared openly, and no additional configuration is required to run the project.
