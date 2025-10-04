# Expensio API Documentation

Base URL: `http://localhost:8080/api/v1`

## Authentication

All authenticated endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <access_token>
```

---

## Auth Endpoints

### 1. Signup

**POST** `/auth/signup`

Creates a new user and company.

**Request Body:**

```json
{
  "email": "admin@company.com",
  "password": "SecurePass123",
  "first_name": "John",
  "last_name": "Doe",
  "company_name": "Acme Corp",
  "country": "US"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Signup successful",
  "data": {
    "user": { ... },
    "company": { ... },
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc..."
  }
}
```

---

### 2. Login

**POST** `/auth/login`

Authenticates a user.

**Request Body:**

```json
{
  "email": "admin@company.com",
  "password": "SecurePass123"
}
```

**Response:** Same as Signup

---

### 3. Refresh Token

**POST** `/auth/refresh`

Generates a new access token using refresh token.

**Request Body:**

```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Response:**

```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGc..."
  }
}
```

---

### 4. Logout

**POST** `/auth/logout`
🔐 **Requires Authentication**

Logs out the current user.

**Response:**

```json
{
  "success": true,
  "message": "Logout successful"
}
```

---

## User Management Endpoints

### 5. Create User

**POST** `/users`
🔐 **Requires Authentication** (Admin only)

Creates a new employee or manager.

**Request Body:**

```json
{
  "email": "employee@company.com",
  "password": "SecurePass123",
  "first_name": "Jane",
  "last_name": "Smith",
  "role": "employee",
  "manager_id": "60d5ec49f1b2c8b1f8e4e1a1"
}
```

**Response:**

```json
{
  "success": true,
  "message": "User created successfully",
  "data": { ... }
}
```

---

### 6. Get All Users

**GET** `/users`
🔐 **Requires Authentication** (Admin/Manager)

Retrieves all users in the company.

**Response:**

```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [ ... ]
}
```

---

### 7. Get User by ID

**GET** `/users/:id`
🔐 **Requires Authentication**

Retrieves a specific user.

**Response:**

```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": { ... }
}
```

---

### 8. Update User Role

**PUT** `/users/:id/role`
🔐 **Requires Authentication** (Admin only)

Updates a user's role.

**Request Body:**

```json
{
  "role": "manager"
}
```

---

### 9. Assign Manager

**PUT** `/users/:id/manager`
🔐 **Requires Authentication** (Admin/Manager)

Assigns a manager to a user.

**Request Body:**

```json
{
  "manager_id": "60d5ec49f1b2c8b1f8e4e1a1"
}
```

---

### 10. Delete User

**DELETE** `/users/:id`
🔐 **Requires Authentication** (Admin only)

Deletes a user.

---

## Expense Management Endpoints

### 11. Create Expense

**POST** `/expenses`
🔐 **Requires Authentication**

Submits a new expense claim.

**Request Body:**

```json
{
  "amount": 150.5,
  "currency": "USD",
  "category": "meals",
  "description": "Team lunch meeting",
  "expense_date": "2025-10-01T12:00:00Z",
  "receipt_url": "/uploads/receipt123.jpg",
  "merchant": "Restaurant Name"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Expense created successfully",
  "data": {
    "id": "...",
    "amount": 150.50,
    "currency": "USD",
    "converted_amount": 150.50,
    "exchange_rate": 1.0,
    "status": "pending",
    ...
  }
}
```

---

### 12. Get Expenses

**GET** `/expenses?page=1&limit=10`
🔐 **Requires Authentication**

Retrieves expenses with pagination.

- Employees: See their own expenses
- Managers/Admins: See all company expenses

**Response:**

```json
{
  "success": true,
  "message": "Expenses retrieved successfully",
  "data": [ ... ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "totalPages": 10
  }
}
```

---

### 13. Get Expense by ID

**GET** `/expenses/:id`
🔐 **Requires Authentication**

Retrieves a specific expense.

---

### 14. Update Expense

**PUT** `/expenses/:id`
🔐 **Requires Authentication**

Updates an expense (only if status is pending).

**Request Body:** Same as Create Expense

---

### 15. Delete Expense

**DELETE** `/expenses/:id`
🔐 **Requires Authentication**

Deletes an expense (only if status is pending).

---

### 16. Get Pending Expenses

**GET** `/expenses/pending`
🔐 **Requires Authentication** (Manager/Admin)

Retrieves all pending expenses for approval.

---

## Approval Workflow Endpoints

### 17. Get Pending Approvals

**GET** `/approvals/pending`
🔐 **Requires Authentication** (Manager/Admin)

Retrieves pending approvals for the current user.

**Response:**

```json
{
  "success": true,
  "message": "Pending approvals retrieved successfully",
  "data": [
    {
      "id": "...",
      "expense_id": "...",
      "approver_id": "...",
      "level": 1,
      "status": "pending",
      ...
    }
  ]
}
```

---

### 18. Approve Expense

**POST** `/approvals/:id/approve`
🔐 **Requires Authentication** (Manager/Admin)

Approves an expense.

**Request Body:**

```json
{
  "comments": "Approved - Valid business expense"
}
```

---

### 19. Reject Expense

**POST** `/approvals/:id/reject`
🔐 **Requires Authentication** (Manager/Admin)

Rejects an expense.

**Request Body:**

```json
{
  "comments": "Rejected - Missing receipt"
}
```

---

### 20. Get Approval History

**GET** `/approvals/history/:expenseId`
🔐 **Requires Authentication**

Retrieves approval history for an expense.

---

## OCR Endpoints

### 21. Upload Receipt

**POST** `/ocr/upload?create_expense=true`
🔐 **Requires Authentication**

Uploads a receipt image for OCR processing.

**Request:**

- Content-Type: `multipart/form-data`
- Field name: `receipt`
- Supported formats: JPG, PNG, PDF
- Max size: 10MB

**Query Parameters:**

- `create_expense` (optional): Set to `true` to automatically create expense from OCR data

**Response:**

```json
{
  "success": true,
  "message": "Receipt processed successfully",
  "data": {
    "ocr_result": {
      "amount": 45.50,
      "currency": "USD",
      "merchant": "Coffee Shop",
      "date": "2025-10-01T10:30:00Z",
      "category": "meals",
      "confidence": 0.85
    },
    "expense": { ... }
  }
}
```

---

## Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

**Common HTTP Status Codes:**

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation error
- `500 Internal Server Error` - Server error

---

## Role-Based Access Control

**Roles:**

- **Admin**: Full access to all endpoints
- **Manager**: Can manage team expenses and approvals
- **Employee**: Can submit and view own expenses

---

## Approval Workflow

The system supports multiple approval workflow types:

1. **Sequential Approval**: Expenses go through a predefined sequence of approvers
2. **Percentage-based**: Auto-approve when X% of approvers approve
3. **Specific Approver**: Auto-approve when a specific person (e.g., CFO) approves
4. **Hybrid**: Combination of the above rules

---

## Caching Strategy

The API uses Redis caching for:

- **User lists**: 15 minutes TTL
- **Expense lists**: 15 minutes TTL
- **Pending approvals**: 5 minutes TTL
- **Currency rates**: 1 hour TTL
- **OCR results**: 24 hours TTL
- **Auth sessions**: Based on JWT expiry

Cache is automatically invalidated when:

- Related data is updated
- Approvals are processed
- Users are modified

---

## Currency Conversion

The system automatically converts all expenses to the company's base currency using real-time exchange rates from ExchangeRate API.

**Supported currencies:** All ISO 4217 currency codes

---

## Testing the API

### Using cURL:

```bash
# Signup
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@company.com",
    "password": "SecurePass123",
    "first_name": "John",
    "last_name": "Doe",
    "company_name": "Acme Corp",
    "country": "US"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@company.com",
    "password": "SecurePass123"
  }'

# Create expense (with auth token)
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "amount": 150.50,
    "currency": "USD",
    "category": "meals",
    "description": "Team lunch",
    "expense_date": "2025-10-01T12:00:00Z"
  }'
```

### Using Postman:

1. Import the API collection (if provided)
2. Set up environment variables:
   - `base_url`: `http://localhost:8080/api/v1`
   - `access_token`: Your JWT token
3. Test endpoints sequentially

---

## Rate Limiting

Currently not implemented. Consider adding rate limiting for production use.

---

## Webhooks (Future Enhancement)

Plan to add webhooks for:

- Expense approved/rejected notifications
- New expense submitted
- Payment processed
