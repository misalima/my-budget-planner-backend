# My Budget Planner Backend

A RESTful API backend service for My Budget Planner, an expense tracker and budget management application built with Go.

## Description

This backend service provides APIs for managing personal finances, including user authentication, expense tracking, budget management, categories, and credit card management. The application follows clean architecture principles with clear separation between domain logic, infrastructure, and API layers.

## Technologies Used

- **Go 1.24.2** - Programming language
- **Echo v4** - Web framework for building REST APIs
- **PostgreSQL** - Primary database
- **pgx/v5** - PostgreSQL driver and toolkit
- **JWT** - Authentication and authorization
- **Docker & Docker Compose** - Containerization and orchestration
- **godotenv** - Environment variable management

## Key Dependencies

- `github.com/labstack/echo/v4` - HTTP web framework
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` - JWT implementation
- `github.com/joho/godotenv` - Environment configuration
- `github.com/google/uuid` - UUID generation
- `golang.org/x/crypto` - Cryptographic functions

## Prerequisites

- **Go 1.22+** installed on your system
- **Docker** and **Docker Compose** installed
- **PostgreSQL** (if running without Docker)

## Environment Setup

Create a `.env` file in the root directory with the following variables:

```env
# Database Configuration
MBP_PG_USER=your_db_user
MBP_PG_NAME=your_db_name
MBP_PG_PASSWORD=your_db_password
MBP_PG_PORT=5432
MBP_PG_HOST=localhost

# PgAdmin Configuration (optional)
PGADMIN_DEFAULT_EMAIL=admin@admin.com
PGADMIN_DEFAULT_PASSWORD=password
```

## How to Run

### Option 1: Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/misalima/my-budget-planner-backend.git
cd my-budget-planner-backend
```

2. Create your `.env` file with the required environment variables

3. Start the services:
```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port `5432`
- PgAdmin on port `8081` (optional database management tool)

4. Build and run the application:
```bash
go run cmd/app/main.go
```

The API will be available at `http://localhost:8000`

### Option 2: Manual Setup

1. Start PostgreSQL database manually or use existing instance

2. Install dependencies:
```bash
go mod download
```

3. Set up your `.env` file with database connection details

4. Run the application:
```bash
go run cmd/app/main.go
```

### Option 3: Using Docker for the Application

1. Build the Docker image:
```bash
docker build -t my-budget-planner-backend .
```

2. Run the container:
```bash
docker run -p 8000:8080 --env-file .env my-budget-planner-backend
```

## API Endpoints

The API runs on port `8000` and includes endpoints for:
- User management and authentication
- Category management
- Credit card management
- Budget and expense tracking

## Database Management

If using Docker Compose, PgAdmin is available at `http://localhost:8081` for database management with the credentials specified in your `.env` file.

## Project Structure

```
├── cmd/app/                 # Application entry point
├── internal/
│   ├── api/http/           # HTTP handlers and routing
│   ├── core/               # Business logic and domain models
│   └── infra/postgres/     # Database infrastructure
├── config/postgres/        # Database migrations
├── docker-compose.yml      # Docker services configuration
├── Dockerfile             # Application containerization
└── .env                   # Environment variables
```