# Expensio Backend - Setup Guide

Complete guide to set up and run the Expensio backend application.

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **MongoDB 6.0+** - [Download](https://www.mongodb.com/try/download/community)
- **Redis 7.0+** - [Download](https://redis.io/download)
- **Git** - [Download](https://git-scm.com/downloads)
- **Tesseract OCR** (Optional, for receipt processing) - [Download](https://github.com/tesseract-ocr/tesseract)

## 🚀 Quick Start

### Option 1: Using Docker (Recommended)

The easiest way to run the application with all dependencies:

```bash
# 1. Clone the repository (if not already done)
cd backend

# 2. Start all services with Docker Compose
docker-compose up -d

# 3. Check logs
docker-compose logs -f backend
```

The API will be available at `http://localhost:8080`

### Option 2: Local Development

#### Step 1: Install Dependencies

```bash
# Install Go dependencies
go mod download
go mod tidy
```

#### Step 2: Configure Environment

```bash
# Copy environment template
copy .env.example .env

# Edit .env file with your configuration
# Update MongoDB URI, Redis settings, and API keys
```

**Required Environment Variables:**

```env
# Server
PORT=8080
APP_ENV=development

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=expensio

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=7d

# External APIs
RESTCOUNTRIES_API_URL=https://restcountries.com/v3.1
EXCHANGERATE_API_URL=https://api.exchangerate-api.com/v4/latest
```

#### Step 3: Start MongoDB and Redis

**Windows:**

```bash
# Start MongoDB
mongod

# Start Redis (in another terminal)
redis-server
```

**Linux/Mac:**

```bash
# Start MongoDB
sudo systemctl start mongod

# Start Redis
sudo systemctl start redis
```

#### Step 4: Run the Application

```bash
# Using Go run
go run cmd/server/main.go

# Or build and run
make build
./bin/expensio
```

## 📦 Available Make Commands

```bash
make install        # Install dependencies
make build          # Build the application
make run            # Run the application
make dev            # Run with hot reload (requires air)
make test           # Run tests
make clean          # Clean build artifacts
make fmt            # Format code
make docker-build   # Build Docker image
make docker-run     # Run with Docker Compose
make setup          # Complete development setup
```

## 🧪 Testing the API

### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok",
  "message": "Expensio API is running"
}
```

### Create First Admin User

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@company.com",
    "password": "SecurePass123",
    "first_name": "Admin",
    "last_name": "User",
    "company_name": "My Company",
    "country": "US"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@company.com",
    "password": "SecurePass123"
  }'
```

Save the `access_token` from the response for authenticated requests.

### Create Expense

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "amount": 150.50,
    "currency": "USD",
    "category": "meals",
    "description": "Team lunch meeting",
    "expense_date": "2025-10-01T12:00:00Z"
  }'
```

## 🏗️ Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── domain/                  # Domain models and interfaces
│   ├── handler/                 # HTTP handlers (controllers)
│   ├── middleware/              # Custom middleware
│   ├── repository/              # Data access layer
│   ├── routes/                  # Route definitions
│   └── service/                 # Business logic layer
├── pkg/
│   ├── cache/                   # Redis cache utilities
│   ├── currency/                # Currency conversion
│   ├── database/                # MongoDB connection
│   ├── jwt/                     # JWT utilities
│   ├── ocr/                     # OCR processing
│   ├── response/                # Response formatters
│   └── validator/               # Request validation
├── .env.example                 # Environment template
├── .gitignore                   # Git ignore rules
├── docker-compose.yml           # Docker Compose config
├── Dockerfile                   # Docker build config
├── go.mod                       # Go module definition
├── Makefile                     # Build commands
└── README.md                    # Project documentation
```

## 🔐 Security Best Practices

1. **Change JWT Secret**: Update `JWT_SECRET` in `.env` with a strong random string
2. **Use HTTPS**: In production, always use HTTPS
3. **Set Strong Passwords**: Enforce password policies
4. **Rate Limiting**: Implement rate limiting for production
5. **Database Security**: Use authentication for MongoDB and Redis
6. **Environment Variables**: Never commit `.env` file

## 🎯 Features Implemented

✅ **Authentication & Authorization**

- JWT-based authentication with refresh tokens
- Role-based access control (Admin, Manager, Employee)
- Redis session management
- Secure password hashing with bcrypt

✅ **User & Company Management**

- Auto-create company on first signup
- Auto-fetch currency from RestCountries API
- Create employees and managers
- Assign manager relationships

✅ **Expense Management**

- Submit expense claims with multi-currency support
- Auto-convert to company's base currency
- View expense history with pagination
- Update/delete pending expenses
- Redis caching for faster queries

✅ **Approval Workflow**

- Multi-level approval sequences
- Conditional approval rules:
  - Sequential approval
  - Percentage-based approval
  - Specific approver rule
  - Hybrid rules
- Cached pending approvals
- Approval history tracking

✅ **OCR Receipt Processing**

- Upload receipt images
- Extract expense details (amount, date, merchant, category)
- Cache OCR results
- Auto-create expenses from OCR data

✅ **Performance & Scalability**

- Redis caching layer
- Connection pooling for MongoDB and Redis
- Efficient database indexing
- Cache invalidation policies
- Configurable TTL for different data types

✅ **Clean Architecture**

- Layered architecture (Handler → Service → Repository)
- Dependency injection
- Interface-based design
- Separation of concerns

## 🐛 Troubleshooting

### MongoDB Connection Error

```
Error: failed to connect to MongoDB
```

**Solution:**

- Ensure MongoDB is running: `mongod` or `sudo systemctl start mongod`
- Check `MONGODB_URI` in `.env`
- Verify MongoDB port (default: 27017)

### Redis Connection Error

```
Error: failed to connect to Redis
```

**Solution:**

- Ensure Redis is running: `redis-server` or `sudo systemctl start redis`
- Check `REDIS_HOST` and `REDIS_PORT` in `.env`
- Verify Redis port (default: 6379)

### Port Already in Use

```
Error: listen tcp :8080: bind: address already in use
```

**Solution:**

- Change `PORT` in `.env` to a different port
- Or stop the process using port 8080

### Module Import Errors

```
Error: could not import [package]
```

**Solution:**

```bash
go mod tidy
go mod download
```

## 📚 API Documentation

Detailed API documentation is available in [API_DOCS.md](./API_DOCS.md)

## 🔄 Development Workflow

1. Create a new branch for your feature
2. Make your changes
3. Run tests: `make test`
4. Format code: `make fmt`
5. Build: `make build`
6. Test locally: `make run`
7. Commit and push
8. Create pull request

## 🚢 Deployment

### Docker Deployment

```bash
# Build and push to registry
docker build -t your-registry/expensio-backend:latest .
docker push your-registry/expensio-backend:latest

# Deploy to server
docker pull your-registry/expensio-backend:latest
docker-compose up -d
```

### Manual Deployment

```bash
# Build for production
CGO_ENABLED=0 GOOS=linux go build -o expensio cmd/server/main.go

# Copy binary to server
scp expensio user@server:/opt/expensio/

# Run with systemd or supervisor
```

## 📊 Monitoring

Consider adding:

- **Logging**: Structured logging with levels
- **Metrics**: Prometheus metrics
- **Tracing**: Distributed tracing
- **Health Checks**: Liveness and readiness probes

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

Proprietary - All rights reserved

## 📞 Support

For issues or questions:

- Create an issue in the repository
- Contact the development team

---

**Happy Coding! 🚀**
