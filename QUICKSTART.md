# Quick Start Guide - Expensio

Get Expensio up and running in 5 minutes!

## Prerequisites Check

Make sure you have installed:

- ✅ Go 1.21+ (`go version`)
- ✅ Node.js 18+ (`node --version`)
- ✅ pnpm (`pnpm --version` or install with `npm install -g pnpm`)
- ✅ MongoDB 7.0+ (running on port 27017)
- ✅ Redis 7.0+ (running on port 6379)

## Quick Setup

### 1. Start MongoDB & Redis (Docker)

```bash
# MongoDB
docker run -d -p 27017:27017 --name expensio-mongo mongo:7.0

# Redis
docker run -d -p 6379:6379 --name expensio-redis redis:7.0
```

### 2. Backend Setup

```bash
cd backend

# Copy environment file
copy .env.example .env   # Windows
# OR
cp .env.example .env     # Linux/Mac

# Install dependencies
go mod download

# Run the server
go run cmd/main.go
```

**Backend should start on:** `http://localhost:8080` ✅

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
pnpm install

# Create environment file
# Create .env.local with:
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-secret-change-this
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# Start dev server
pnpm dev
```

**Frontend should start on:** `http://localhost:3000` ✅

## First Time Use

### 1. Create Company & Admin Account

Visit: `http://localhost:3000/signup`

```
Email: admin@mycompany.com
Password: password123
First Name: John
Last Name: Doe
Company Name: My Company
Country: United States (or your country)
```

### 2. Login

Visit: `http://localhost:3000/login`

Use the credentials you just created.

### 3. Explore the Dashboard

After login, you'll see:

- 📊 Dashboard with stats
- 👥 Users (create team members)
- 💰 Expenses (create and track expenses)
- ✅ Approvals (approve/reject expenses)

## Testing the Features

### Create Your First Employee

1. Go to **Users** page
2. Click **Add User**
3. Fill in:
   ```
   Email: employee@mycompany.com
   Password: password123
   First Name: Jane
   Last Name: Smith
   Role: Employee
   ```
4. Click **Create User**

### Create Your First Expense

1. Go to **Expenses** page
2. Click **Add Expense**
3. Fill in:
   ```
   Amount: 50.00
   Currency: USD
   Category: Meals
   Date: Today
   Description: Team lunch
   ```
4. Click **Create Expense**

### Test Approval Workflow

1. Login as the employee account
2. Create an expense
3. Logout and login as admin
4. Go to **Approvals** page
5. See the pending expense
6. Click **Approve** and add a comment
7. Approve the expense

### Test OCR Upload

1. Go to **Expenses** page
2. Click **Upload Receipt**
3. Select a receipt image
4. Watch it automatically create an expense!

## Default Ports

| Service     | Port  | URL                       |
| ----------- | ----- | ------------------------- |
| Backend API | 8080  | http://localhost:8080     |
| Frontend    | 3000  | http://localhost:3000     |
| MongoDB     | 27017 | mongodb://localhost:27017 |
| Redis       | 6379  | localhost:6379            |

## Common Issues

### Backend won't start

**Error:** "Failed to connect to MongoDB"

```bash
# Check if MongoDB is running
docker ps | grep mongo

# If not, start it:
docker start expensio-mongo
```

**Error:** "Failed to connect to Redis"

```bash
# Check if Redis is running
docker ps | grep redis

# If not, start it:
docker start expensio-redis
```

### Frontend won't start

**Error:** "ERR_PNPM_NO_IMPORTER_MANIFEST_FOUND"

```bash
# Make sure you're in the frontend directory
cd frontend
pnpm install
pnpm dev
```

**Error:** "NEXTAUTH_SECRET is not set"

```bash
# Create .env.local file with required variables
echo "NEXTAUTH_SECRET=your-secret-here" > .env.local
echo "NEXTAUTH_URL=http://localhost:3000" >> .env.local
echo "NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1" >> .env.local
```

### Can't login

**Issue:** "Login failed" or "Invalid credentials"

1. Make sure backend is running on port 8080
2. Check backend logs for errors
3. Verify MongoDB is running
4. Try creating a new account via signup

### API requests fail

**Issue:** CORS errors or network errors

1. Verify `NEXT_PUBLIC_API_URL` in `.env.local` is correct
2. Make sure backend is running
3. Check browser console for specific errors
4. Verify backend logs show incoming requests

## Environment Files

### Backend `.env`

```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=expensio
REDIS_URL=localhost:6379
REDIS_PASSWORD=
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h
EXCHANGE_RATE_API_KEY=
TESSERACT_PATH=/usr/bin/tesseract
```

### Frontend `.env.local`

```env
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-nextauth-secret-change-in-production
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## API Testing

### Using cURL

```bash
# Signup
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@company.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User",
    "company_name": "Test Company",
    "country": "US"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@company.com",
    "password": "password123"
  }'
```

### Using Postman

Import the Postman collection from:

```
backend/docs/postman/Expensio_API_Collection.json
```

## Next Steps

- 📚 Read the full [README.md](./README.md) for detailed documentation
- 🔐 Change default secrets in production
- 🧪 Explore all features
- 📸 Test OCR with receipt images
- 👥 Add team members
- 💰 Create expense reports

## Getting Help

- Check the [README.md](./README.md) for detailed docs
- Review the code comments
- Check backend logs: Terminal running `go run cmd/main.go`
- Check frontend logs: Terminal running `pnpm dev`
- Browser console for frontend errors

## Success Checklist

- [ ] Backend running on port 8080
- [ ] Frontend running on port 3000
- [ ] MongoDB connected
- [ ] Redis connected
- [ ] Signup working
- [ ] Login working
- [ ] Dashboard loads
- [ ] Can create users
- [ ] Can create expenses
- [ ] Can approve expenses

**Ready to go! 🚀**
