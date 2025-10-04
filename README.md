# Expensio - Expense Management System

A full-stack expense management application built with Go (Fiber) backend and Next.js 15 frontend.

## Features

### 🔐 Authentication & Authorization

- JWT-based authentication with access & refresh tokens
- Role-based access control (Admin, Manager, Employee)
- NextAuth integration for session management

### 👥 User Management (Admin)

- Create, edit, and delete users
- Assign roles (Admin, Manager, Employee)
- Assign managers to employees
- View all team members

### 💰 Expense Management

- Create, view, edit, and delete expenses
- Multi-currency support with automatic conversion
- Expense categories: Meals, Travel, Accommodation, Entertainment, Office Supplies, Software, Other
- Attach receipt URLs
- Track expense status (Pending, Approved, Rejected)

### ✅ Approval Workflow (Manager/Admin)

- View pending expenses requiring approval
- Approve or reject expenses with comments
- Multi-level approval system
- Approval history tracking

### 📸 OCR Receipt Processing

- Upload receipt images
- Automatic text extraction using Tesseract OCR
- Auto-populate expense details from receipt
- Extract merchant, amount, date, and category

### 📊 Dashboard

- Real-time expense statistics
- Total expenses and amount
- Pending and approved counts
- Role-based navigation
- Quick action cards

## Tech Stack

### Backend

- **Framework**: Go 1.21+ with Fiber v2.52.0
- **Database**: MongoDB 7.0+
- **Cache**: Redis 7.0+
- **Authentication**: JWT tokens
- **External APIs**:
  - RestCountries API (currency detection)
  - ExchangeRate API (currency conversion)
  - Tesseract OCR (receipt processing)

### Frontend

- **Framework**: Next.js 15.5.4 with React 19
- **Authentication**: NextAuth
- **UI Components**: shadcn/ui + Tailwind CSS
- **Language**: TypeScript
- **Package Manager**: pnpm
- **Date Handling**: date-fns
- **Notifications**: Sonner toast

## Project Structure

```
Expensio/
├── backend/                    # Go backend
│   ├── cmd/
│   │   └── main.go            # Application entry point
│   ├── internal/
│   │   ├── config/            # Configuration management
│   │   ├── domain/            # Domain models
│   │   ├── handler/           # HTTP handlers (controllers)
│   │   ├── middleware/        # JWT, CORS, logging middleware
│   │   ├── repository/        # Data access layer
│   │   ├── routes/            # Route definitions
│   │   ├── service/           # Business logic
│   │   └── utils/             # Helper functions
│   ├── .env                   # Environment variables
│   ├── go.mod                 # Go dependencies
│   └── Dockerfile             # Docker configuration
│
└── frontend/                   # Next.js frontend
    ├── src/
    │   ├── app/
    │   │   ├── (authenticated)/ # Protected routes
    │   │   │   ├── dashboard/   # Dashboard page
    │   │   │   ├── users/       # User management
    │   │   │   ├── expenses/    # Expense management
    │   │   │   └── approvals/   # Approval workflow
    │   │   ├── login/          # Login page
    │   │   └── signup/         # Signup page
    │   ├── components/         # React components
    │   │   ├── ui/            # shadcn/ui components
    │   │   ├── app-sidebar.tsx
    │   │   ├── navbar.tsx
    │   │   └── expenses-client.tsx
    │   ├── lib/               # Utilities
    │   │   ├── api.ts         # API service functions
    │   │   └── api-client.ts  # HTTP client with auth
    │   ├── hooks/             # Custom React hooks
    │   │   └── use-toast.ts
    │   └── types/             # TypeScript types
    │       └── api.ts
    ├── types/
    │   └── auth.d.ts          # NextAuth type extensions
    ├── auth.ts                # NextAuth configuration
    ├── package.json
    └── next.config.ts
```

## Setup Instructions

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ and pnpm
- MongoDB 7.0+
- Redis 7.0+
- Tesseract OCR (for receipt processing)

### Backend Setup

1. **Navigate to backend directory:**

   ```bash
   cd backend
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Configure environment variables:**
   Create a `.env` file in the backend directory:

   ```env
   # Server
   PORT=8080

   # MongoDB
   MONGODB_URI=mongodb://localhost:27017
   MONGODB_DATABASE=expensio

   # Redis
   REDIS_URL=localhost:6379
   REDIS_PASSWORD=

   # JWT
   JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
   JWT_ACCESS_EXPIRY=15m
   JWT_REFRESH_EXPIRY=168h

   # External APIs
   EXCHANGE_RATE_API_KEY=your-exchange-rate-api-key
   TESSERACT_PATH=/usr/bin/tesseract
   ```

4. **Start MongoDB and Redis:**

   ```bash
   # Using Docker
   docker run -d -p 27017:27017 --name mongodb mongo:7.0
   docker run -d -p 6379:6379 --name redis redis:7.0
   ```

5. **Run the backend:**

   ```bash
   go run cmd/main.go
   ```

   The backend will start on `http://localhost:8080`

### Frontend Setup

1. **Navigate to frontend directory:**

   ```bash
   cd frontend
   ```

2. **Install dependencies:**

   ```bash
   pnpm install
   ```

3. **Configure environment variables:**
   Create a `.env.local` file:

   ```env
   NEXTAUTH_URL=http://localhost:3000
   NEXTAUTH_SECRET=your-nextauth-secret-change-this-in-production
   NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
   ```

4. **Run the frontend:**

   ```bash
   pnpm dev
   ```

   The frontend will start on `http://localhost:3000`

## API Documentation

### Authentication Endpoints

#### Signup

```
POST /api/v1/auth/signup
Content-Type: application/json

{
  "email": "admin@company.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "company_name": "Acme Corp",
  "country": "US"
}
```

#### Login

```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@company.com",
  "password": "password123"
}
```

#### Refresh Token

```
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

### User Endpoints

#### List Users

```
GET /api/v1/users
Authorization: Bearer {access_token}
```

#### Create User

```
POST /api/v1/users
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "email": "user@company.com",
  "password": "password123",
  "first_name": "Jane",
  "last_name": "Smith",
  "role": "employee"
}
```

#### Update User Role

```
PUT /api/v1/users/{id}/role
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "role": "manager"
}
```

#### Assign Manager

```
PUT /api/v1/users/{id}/manager
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "manager_id": "manager-user-id"
}
```

### Expense Endpoints

#### List Expenses

```
GET /api/v1/expenses?page=1&limit=10
Authorization: Bearer {access_token}
```

#### Create Expense

```
POST /api/v1/expenses
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "amount": 150.50,
  "currency": "USD",
  "category": "meals",
  "description": "Team lunch",
  "expense_date": "2025-10-04",
  "merchant": "Restaurant ABC"
}
```

#### Get Pending Expenses

```
GET /api/v1/expenses/pending
Authorization: Bearer {access_token}
```

### Approval Endpoints

#### Get Pending Approvals

```
GET /api/v1/approvals/pending
Authorization: Bearer {access_token}
```

#### Approve Expense

```
POST /api/v1/approvals/{id}/approve
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "comments": "Approved - within budget"
}
```

#### Reject Expense

```
POST /api/v1/approvals/{id}/reject
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "comments": "Rejected - missing receipt"
}
```

### OCR Endpoints

#### Upload Receipt

```
POST /api/v1/ocr/upload
Authorization: Bearer {access_token}
Content-Type: multipart/form-data

receipt: [image file]
create_expense: true
```

## Testing with Postman

A complete Postman collection is available in `backend/docs/postman/`. Import the collection and environment to test all endpoints.

## Architecture

### Backend Architecture

```
┌─────────────┐
│   Routes    │  (HTTP routing)
└──────┬──────┘
       │
┌──────▼──────┐
│  Middleware │  (JWT, CORS, Logging)
└──────┬──────┘
       │
┌──────▼──────┐
│  Handlers   │  (HTTP controllers)
└──────┬──────┘
       │
┌──────▼──────┐
│  Services   │  (Business logic)
└──────┬──────┘
       │
┌──────▼──────┐
│ Repository  │  (Data access)
└──────┬──────┘
       │
┌──────▼──────┐
│  MongoDB    │  (Database)
│   Redis     │  (Cache)
└─────────────┘
```

### Frontend Architecture

```
┌─────────────┐
│   Pages     │  (Next.js App Router)
└──────┬──────┘
       │
┌──────▼──────┐
│ Components  │  (React components)
└──────┬──────┘
       │
┌──────▼──────┐
│  API Layer  │  (API client with auth)
└──────┬──────┘
       │
┌──────▼──────┐
│   Backend   │  (Go API)
└─────────────┘
```

## User Roles & Permissions

### Admin

- Full access to all features
- User management (create, edit, delete users)
- View all expenses across organization
- Approve/reject expenses
- Configure approval rules

### Manager

- View team members' expenses
- Approve/reject expenses for direct reports
- Create own expenses
- Cannot manage users

### Employee

- Create and view own expenses
- Upload receipts
- Track expense status
- Cannot approve expenses

## Development

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
pnpm test
```

### Building for Production

#### Backend

```bash
cd backend
go build -o expensio cmd/main.go
./expensio
```

#### Frontend

```bash
cd frontend
pnpm build
pnpm start
```

### Docker Deployment

```bash
# Build and run with Docker Compose
docker-compose up -d
```

## Environment Variables

### Backend (.env)

| Variable              | Description               | Default                   |
| --------------------- | ------------------------- | ------------------------- |
| PORT                  | Server port               | 8080                      |
| MONGODB_URI           | MongoDB connection string | mongodb://localhost:27017 |
| MONGODB_DATABASE      | Database name             | expensio                  |
| REDIS_URL             | Redis connection URL      | localhost:6379            |
| JWT_SECRET            | JWT signing secret        | -                         |
| JWT_ACCESS_EXPIRY     | Access token expiry       | 15m                       |
| JWT_REFRESH_EXPIRY    | Refresh token expiry      | 168h                      |
| EXCHANGE_RATE_API_KEY | Exchange rate API key     | -                         |
| TESSERACT_PATH        | Tesseract binary path     | /usr/bin/tesseract        |

### Frontend (.env.local)

| Variable            | Description     | Default                      |
| ------------------- | --------------- | ---------------------------- |
| NEXTAUTH_URL        | Application URL | http://localhost:3000        |
| NEXTAUTH_SECRET     | NextAuth secret | -                            |
| NEXT_PUBLIC_API_URL | Backend API URL | http://localhost:8080/api/v1 |

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Support

For issues and questions:

- Open an issue on GitHub
- Contact: support@expensio.com

## Changelog

### Version 1.0.0 (2025-10-04)

- Initial release
- Complete expense management system
- User management with roles
- Multi-level approval workflow
- OCR receipt processing
- Multi-currency support
- Dashboard with statistics
