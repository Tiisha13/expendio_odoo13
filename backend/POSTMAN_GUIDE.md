# 📬 Postman Collection Guide - Expensio Backend API

## 📁 Files Included

1. **Expensio.postman_collection.json** - Complete API collection with all endpoints
2. **Expensio.postman_environment.json** - Local environment configuration

---

## 🚀 Quick Setup

### Step 1: Import into Postman

1. Open Postman
2. Click **Import** button (top left)
3. Drag and drop both files:
   - `Expensio.postman_collection.json`
   - `Expensio.postman_environment.json`
4. Click **Import**

### Step 2: Select Environment

1. Click the **Environment dropdown** (top right)
2. Select **Expensio Local Environment**

### Step 3: Start Backend Server

```bash
cd d:\Expensio\backend
go run cmd/server/main.go
```

Or use Docker:

```bash
docker-compose up -d
```

---

## 🎯 Testing Workflow

### 1️⃣ Health Check (No Auth Required)

**Request:** `GET /health`

**Expected Response:**

```json
{
  "success": true,
  "message": "Server is healthy",
  "data": {
    "status": "ok",
    "timestamp": "2025-10-04T12:00:00Z"
  }
}
```

---

### 2️⃣ Create First Account (Signup)

**Request:** `POST /api/v1/auth/signup`

**Body:**

```json
{
  "email": "admin@acmecorp.com",
  "password": "Admin@123",
  "first_name": "John",
  "last_name": "Doe",
  "company_name": "Acme Corporation",
  "country": "US"
}
```

**What Happens:**

- ✅ Creates user account
- ✅ Automatically creates company
- ✅ Auto-detects base currency (USD for US)
- ✅ Assigns admin role
- ✅ Returns access & refresh tokens
- ✅ **Tokens automatically saved to environment!**

**Expected Response:**

```json
{
  "success": true,
  "message": "Signup successful",
  "data": {
    "user": {
      "id": "67001234567890abcdef1234",
      "email": "admin@acmecorp.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "admin",
      "company_id": "67001234567890abcdef5678"
    },
    "company": {
      "id": "67001234567890abcdef5678",
      "name": "Acme Corporation",
      "base_currency": "USD",
      "country": "US"
    },
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc..."
  }
}
```

---

### 3️⃣ Login (Subsequent Sessions)

**Request:** `POST /api/v1/auth/login`

**Body:**

```json
{
  "email": "admin@acmecorp.com",
  "password": "Admin@123"
}
```

**What Happens:**

- ✅ Validates credentials
- ✅ Creates new session in Redis
- ✅ Returns fresh tokens
- ✅ **Tokens automatically saved to environment!**

---

### 4️⃣ Create Additional Users (Admin Only)

**Request:** `POST /api/v1/users`  
**Auth:** Bearer Token (Auto-applied from environment)

**Body - Create Manager:**

```json
{
  "email": "manager@acmecorp.com",
  "password": "Manager@123",
  "first_name": "Jane",
  "last_name": "Smith",
  "role": "manager"
}
```

**Body - Create Employee:**

```json
{
  "email": "employee@acmecorp.com",
  "password": "Employee@123",
  "first_name": "Bob",
  "last_name": "Johnson",
  "role": "employee"
}
```

**Available Roles:**

- `admin` - Full system access
- `manager` - Can approve expenses, manage users
- `employee` - Can submit expenses

---

### 5️⃣ Submit Expense

**Request:** `POST /api/v1/expenses`

**Body:**

```json
{
  "amount": 150.5,
  "currency": "USD",
  "category": "meals",
  "description": "Team lunch at downtown restaurant",
  "expense_date": "2025-10-04T12:30:00Z",
  "merchant": "The Fancy Restaurant",
  "receipt_url": "https://example.com/receipts/receipt123.jpg"
}
```

**Available Categories:**

- `meals`
- `travel`
- `accommodation`
- `entertainment`
- `office_supplies`
- `software`
- `other`

**What Happens:**

- ✅ Validates expense data
- ✅ Converts to company base currency (if different)
- ✅ Creates expense record
- ✅ Initializes approval workflow
- ✅ Invalidates expense cache

---

### 6️⃣ View Pending Approvals (Manager/Admin)

**Request:** `GET /api/v1/approvals/pending`

**Expected Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "67001234567890abcdef9999",
      "expense_id": "67001234567890abcdef8888",
      "expense": {
        "amount": 150.5,
        "currency": "USD",
        "category": "meals",
        "description": "Team lunch",
        "user_name": "Bob Johnson"
      },
      "level": 1,
      "status": "pending",
      "created_at": "2025-10-04T12:30:00Z"
    }
  ]
}
```

---

### 7️⃣ Approve or Reject Expense

**Approve:** `POST /api/v1/approvals/:id/approve`

**Body:**

```json
{
  "comments": "Approved - valid business expense"
}
```

**Reject:** `POST /api/v1/approvals/:id/reject`

**Body:**

```json
{
  "comments": "Rejected - insufficient documentation"
}
```

**What Happens on Approve:**

- ✅ Marks approval as approved
- ✅ Checks approval rules (Sequential/Percentage/Specific/Hybrid)
- ✅ Auto-approves expense if rules satisfied
- ✅ Moves to next approval level if needed
- ✅ Invalidates approval & expense caches

---

### 8️⃣ Upload Receipt with OCR

**Request:** `POST /api/v1/ocr/upload`  
**Content-Type:** `multipart/form-data`

**Form Data:**

- `receipt` (file) - Image file (JPG, PNG, PDF)
- `create_expense` (text) - "true" to auto-create expense

**What Happens:**

- ✅ Uploads file to server
- ✅ Processes with Tesseract OCR
- ✅ Extracts: amount, date, merchant, category, currency
- ✅ Caches OCR result (24 hours)
- ✅ Optionally creates expense automatically

**Expected Response:**

```json
{
  "success": true,
  "data": {
    "ocr_result": {
      "id": "67001234567890abcdef7777",
      "amount": 45.99,
      "currency": "USD",
      "merchant": "ACME STORE",
      "date": "2025-10-04",
      "category": "office_supplies",
      "confidence": 0.85
    },
    "expense_created": true,
    "expense_id": "67001234567890abcdef6666"
  }
}
```

---

## 🔐 Authentication

### How Tokens Work

1. **Access Token** (15 minutes)

   - Used for all authenticated requests
   - Automatically included in `Authorization` header
   - Format: `Bearer eyJhbGc...`

2. **Refresh Token** (7 days)
   - Used to get new access token
   - Call `/api/v1/auth/refresh` when access token expires

### Auto-Token Management

The collection includes test scripts that automatically:

- ✅ Extract tokens from login/signup responses
- ✅ Save to environment variables
- ✅ Apply to all subsequent requests
- ✅ Clear on logout

### Manual Token Refresh

If access token expires:

1. Go to **Authentication** → **Refresh Token**
2. Click **Send**
3. New access token automatically saved

---

## 📊 Response Format

### Success Response

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

### Paginated Response

```json
{
  "success": true,
  "data": [ ... ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

### Error Response

```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

---

## 🔍 Common Status Codes

| Code | Meaning      | Example                   |
| ---- | ------------ | ------------------------- |
| 200  | OK           | Successful GET/PUT/DELETE |
| 201  | Created      | Successful POST           |
| 400  | Bad Request  | Invalid input data        |
| 401  | Unauthorized | Missing/invalid token     |
| 403  | Forbidden    | Insufficient permissions  |
| 404  | Not Found    | Resource doesn't exist    |
| 500  | Server Error | Internal server error     |

---

## 🎨 Folder Structure

```
📁 Expensio Backend API
  ├── 📁 Health Check
  │   └── Health Check
  │
  ├── 📁 Authentication
  │   ├── Signup (Create Company & Admin)
  │   ├── Login
  │   ├── Refresh Token
  │   └── Logout
  │
  ├── 📁 User Management
  │   ├── Create User (Admin Only)
  │   ├── List Users (Admin/Manager)
  │   ├── Get User by ID
  │   ├── Update User Role (Admin Only)
  │   ├── Assign Manager (Admin/Manager)
  │   └── Delete User (Admin Only)
  │
  ├── 📁 Expense Management
  │   ├── Create Expense
  │   ├── Get My Expenses
  │   ├── Get Expense by ID
  │   ├── Update Expense
  │   ├── Delete Expense
  │   └── Get Pending Expenses (Manager/Admin)
  │
  ├── 📁 Approval Workflow
  │   ├── Get Pending Approvals
  │   ├── Approve Expense
  │   ├── Reject Expense
  │   └── Get Approval History
  │
  └── 📁 OCR Receipt Processing
      └── Upload Receipt
```

---

## 🧪 Testing Scenarios

### Scenario 1: Complete Expense Workflow

1. **Login as Employee**
   - Use employee credentials
2. **Create Expense**
   - Submit expense claim
3. **Login as Manager**
   - Switch to manager account
4. **View Pending Approvals**
   - See employee's expense
5. **Approve/Reject**
   - Make decision with comments
6. **View Approval History**
   - Check complete audit trail

### Scenario 2: Multi-Currency Flow

1. **Create Expense in EUR**

   ```json
   {
     "amount": 100,
     "currency": "EUR",
     "category": "travel"
   }
   ```

2. **System Auto-Converts**
   - Fetches EUR → USD rate
   - Stores converted amount
   - Caches rate for 1 hour

### Scenario 3: OCR to Expense

1. **Upload Receipt**
   - Set `create_expense=true`
2. **OCR Extracts Data**
   - Amount, merchant, date, category
3. **Expense Auto-Created**
   - Pre-filled with OCR data
   - Ready for approval

---

## 🛠️ Environment Variables

| Variable       | Description        | Auto-Set       |
| -------------- | ------------------ | -------------- |
| `baseUrl`      | API base URL       | No             |
| `accessToken`  | JWT access token   | Yes (on login) |
| `refreshToken` | JWT refresh token  | Yes (on login) |
| `userId`       | Current user ID    | Yes (on login) |
| `companyId`    | Current company ID | Yes (on login) |

---

## 🐛 Troubleshooting

### Issue: "Unauthorized" Error

**Solution:**

1. Check if access token is set: `{{accessToken}}`
2. Token might be expired - use **Refresh Token** endpoint
3. Or login again to get fresh tokens

### Issue: "Forbidden" Error

**Solution:**

- You don't have permission for this action
- Check your role (admin/manager/employee)
- Some endpoints require specific roles

### Issue: "Connection Refused"

**Solution:**

1. Ensure backend server is running: `http://localhost:8080/health`
2. Check MongoDB is running on `localhost:27017`
3. Check Redis is running on `localhost:6379`

### Issue: OCR Upload Fails

**Solution:**

1. Check file size < 10MB
2. Ensure Tesseract is installed
3. Supported formats: JPG, PNG, PDF

---

## 📝 Tips & Best Practices

### 1. Use Test Scripts

The collection includes automatic token management. Just login once!

### 2. Check Console

View saved variables in Postman Console:

- View → Show Postman Console
- See token save confirmations

### 3. Organize Requests

Use **Collections** to group related tests:

- Create folders for different test scenarios
- Use **Run Collection** for automated testing

### 4. Environment Switching

Create multiple environments:

- **Local** - `http://localhost:8080`
- **Development** - `https://dev.expensio.com`
- **Production** - `https://api.expensio.com`

### 5. Save Example Responses

Right-click request → **Save Response** → **Save as Example**

- Useful for documentation
- Helps understand expected format

---

## 🚀 Advanced Usage

### Run Collection Automatically

1. Click **Runner** button
2. Select **Expensio Backend API**
3. Select environment
4. Click **Run Expensio Backend API**
5. View results dashboard

### Export Variables

```javascript
// In Tests tab of any request
pm.environment.set("expenseId", jsonData.data.id);
```

### Chain Requests

Use response data in next request:

```
{{expenseId}}  // Use saved variable
```

### Pre-request Scripts

Add setup logic before request:

```javascript
// Generate random email
pm.environment.set("randomEmail", "user" + Math.random() + "@test.com");
```

---

## 📚 Additional Resources

- **API Documentation**: `API_DOCS.md`
- **Setup Guide**: `SETUP_GUIDE.md`
- **Project Summary**: `PROJECT_SUMMARY.md`
- **README**: `README.md`

---

## ✅ Quick Checklist

- [ ] Import collection and environment
- [ ] Select environment in Postman
- [ ] Start backend server
- [ ] Test health check
- [ ] Run signup to create admin account
- [ ] Verify tokens are saved
- [ ] Create test users
- [ ] Submit test expense
- [ ] Test approval workflow
- [ ] Try OCR upload

---

## 🎉 You're Ready!

Your Postman collection is fully configured with:

- ✅ 25+ API endpoints
- ✅ Automatic token management
- ✅ Environment variables
- ✅ Request examples
- ✅ Response documentation

**Happy Testing! 🚀**

---

## 📞 Support

For issues or questions:

1. Check troubleshooting section
2. Review API documentation
3. Check server logs
4. Verify environment configuration

**Version:** 1.0.0  
**Last Updated:** October 4, 2025
