# Expensio Backend

A scalable Expense Management Web Application backend built with Golang, Fiber, MongoDB, and Redis.

## Features

- 🔐 JWT-based authentication with refresh tokens
- 👥 Multi-tenant company management
- 💰 Multi-currency expense tracking with auto-conversion
- ✅ Multi-level approval workflow with conditional rules
- 📸 OCR receipt processing
- ⚡ High-performance caching with Redis
- 🏗️ Clean architecture with layered design

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **Database**: MongoDB
- **Cache**: Redis
- **Authentication**: JWT
- **OCR**: Tesseract

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── domain/                  # Domain models and interfaces
│   ├── repository/              # Data access layer
│   ├── service/                 # Business logic layer
│   ├── handler/                 # HTTP handlers (controllers)
│   ├── middleware/              # Custom middleware
│   └── routes/                  # Route definitions
├── pkg/
│   ├── cache/                   # Redis cache utilities
│   ├── database/                # MongoDB connection
│   ├── jwt/                     # JWT utilities
│   ├── validator/               # Request validation
│   ├── response/                # Response formatters
│   ├── currency/                # Currency conversion
│   └── ocr/                     # OCR processing
├── .env.example                 # Example environment variables
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MongoDB 6.0+
- Redis 7.0+
- Tesseract OCR (optional)

### Installation

1. Clone the repository

```bash
cd backend
```

2. Copy environment variables

```bash
copy .env.example .env
```

3. Update `.env` with your configuration

4. Install dependencies

```bash
go mod download
```

5. Run the application

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /api/v1/auth/signup` - Register new user and create company
- `POST /api/v1/auth/login` - Login with credentials
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout user

### User Management (Admin only)

- `POST /api/v1/users` - Create employee/manager
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/:id` - Get user details
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `PUT /api/v1/users/:id/role` - Assign/change role

### Expense Management

- `POST /api/v1/expenses` - Submit expense claim
- `GET /api/v1/expenses` - List expenses (filtered by user/company)
- `GET /api/v1/expenses/:id` - Get expense details
- `PUT /api/v1/expenses/:id` - Update expense (before approval)
- `DELETE /api/v1/expenses/:id` - Delete expense

### Approval Workflow

- `GET /api/v1/approvals/pending` - List pending approvals
- `POST /api/v1/approvals/:id/approve` - Approve expense
- `POST /api/v1/approvals/:id/reject` - Reject expense
- `GET /api/v1/approvals/history` - Approval history

### OCR

- `POST /api/v1/ocr/upload` - Upload and process receipt

## 📬 Testing with Postman

A complete Postman collection is included for easy API testing:

- **Collection**: `Expensio.postman_collection.json`
- **Environment**: `Expensio.postman_environment.json`
- **Guide**: `POSTMAN_GUIDE.md`

### Quick Start:

1. Import both JSON files into Postman
2. Select "Expensio Local Environment"
3. Run "Signup" to create admin account (tokens auto-saved!)
4. Start testing all endpoints

See [POSTMAN_GUIDE.md](POSTMAN_GUIDE.md) for detailed usage instructions.

## Architecture

### Clean Architecture Layers

1. **Handler Layer** (`internal/handler`): HTTP request handling, validation, response formatting
2. **Service Layer** (`internal/service`): Business logic, orchestration, external API calls
3. **Repository Layer** (`internal/repository`): Data access, database operations
4. **Domain Layer** (`internal/domain`): Core business entities and interfaces

### Caching Strategy

- **Auth tokens**: Redis session store with TTL
- **Expense lists**: 15-minute cache with invalidation on updates
- **Pending approvals**: 5-minute cache for manager views
- **Currency rates**: 1-hour cache to minimize API calls
- **OCR results**: 24-hour cache for duplicate receipt prevention

### Approval Workflow

Multi-level approval with conditional rules:

1. **Sequential Approval**: Manager → Finance → Director
2. **Percentage Rule**: Auto-approve if X% of approvers approve
3. **Specific Approver Rule**: Auto-approve if specific person (e.g., CFO) approves
4. **Hybrid Rule**: Combination of above rules

## Development

### Build

```bash
go build -o bin/server cmd/server/main.go
```

### Run tests

```bash
go test ./...
```

### Format code

```bash
go fmt ./...
```

## License

Proprietary - All rights reserved
