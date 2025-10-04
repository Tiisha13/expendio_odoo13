# 🚀 Expensio API Quick Reference

## 📍 Base URL

```
http://localhost:8080
```

## 🔐 Authentication Header

```
Authorization: Bearer <access_token>
```

---

## 🏃 Quick Start (3 Steps)

### 1. Health Check

```bash
curl http://localhost:8080/health
```

### 2. Signup (Creates Admin + Company)

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@company.com",
    "password": "Admin@123",
    "first_name": "John",
    "last_name": "Doe",
    "company_name": "Acme Corp",
    "country": "US"
  }'
```

### 3. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@company.com",
    "password": "Admin@123"
  }'
```

---

## 🔑 Authentication

| Endpoint               | Method | Auth | Description                  |
| ---------------------- | ------ | ---- | ---------------------------- |
| `/api/v1/auth/signup`  | POST   | No   | Create account + company     |
| `/api/v1/auth/login`   | POST   | No   | Login and get tokens         |
| `/api/v1/auth/refresh` | POST   | No   | Refresh access token         |
| `/api/v1/auth/logout`  | POST   | Yes  | Logout and invalidate tokens |

### Token Expiry

- **Access Token**: 15 minutes
- **Refresh Token**: 7 days

---

## 👥 User Management

| Endpoint                    | Method | Auth | Role          | Description    |
| --------------------------- | ------ | ---- | ------------- | -------------- |
| `/api/v1/users`             | POST   | Yes  | Admin         | Create user    |
| `/api/v1/users`             | GET    | Yes  | Admin/Manager | List users     |
| `/api/v1/users/:id`         | GET    | Yes  | All           | Get user       |
| `/api/v1/users/:id/role`    | PUT    | Yes  | Admin         | Update role    |
| `/api/v1/users/:id/manager` | PUT    | Yes  | Admin/Manager | Assign manager |
| `/api/v1/users/:id`         | DELETE | Yes  | Admin         | Delete user    |

### Available Roles

- `admin` - Full access
- `manager` - Approvals + user management
- `employee` - Submit expenses

---

## 💰 Expense Management

| Endpoint                           | Method | Auth | Description             |
| ---------------------------------- | ------ | ---- | ----------------------- |
| `/api/v1/expenses`                 | POST   | Yes  | Create expense          |
| `/api/v1/expenses?page=1&limit=10` | GET    | Yes  | List my expenses        |
| `/api/v1/expenses/:id`             | GET    | Yes  | Get expense details     |
| `/api/v1/expenses/:id`             | PUT    | Yes  | Update expense          |
| `/api/v1/expenses/:id`             | DELETE | Yes  | Delete expense          |
| `/api/v1/expenses/pending`         | GET    | Yes  | Pending (Manager/Admin) |

### Expense Categories

- `meals`
- `travel`
- `accommodation`
- `entertainment`
- `office_supplies`
- `software`
- `other`

### Supported Currencies

Any ISO 4217 code: `USD`, `EUR`, `GBP`, `JPY`, `CAD`, etc.

---

## ✅ Approval Workflow

| Endpoint                               | Method | Auth | Role          | Description           |
| -------------------------------------- | ------ | ---- | ------------- | --------------------- |
| `/api/v1/approvals/pending`            | GET    | Yes  | Manager/Admin | My pending approvals  |
| `/api/v1/approvals/:id/approve`        | POST   | Yes  | Manager/Admin | Approve expense       |
| `/api/v1/approvals/:id/reject`         | POST   | Yes  | Manager/Admin | Reject expense        |
| `/api/v1/approvals/history/:expenseId` | GET    | Yes  | All           | View approval history |

### Approval Rules

1. **Sequential**: Manager → Finance → Director
2. **Percentage**: Auto-approve at X% threshold
3. **Specific Approver**: Auto-approve by CFO
4. **Hybrid**: Combination of rules

---

## 📸 OCR Receipt Processing

| Endpoint             | Method | Auth | Description              |
| -------------------- | ------ | ---- | ------------------------ |
| `/api/v1/ocr/upload` | POST   | Yes  | Upload & process receipt |

### Supported Formats

- JPG
- PNG
- PDF

### Max File Size

10 MB

### Extracted Data

- Amount
- Date
- Merchant
- Category
- Currency

---

## 📋 Example Requests

### Create Expense

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "amount": 150.50,
    "currency": "USD",
    "category": "meals",
    "description": "Team lunch",
    "expense_date": "2025-10-04T12:30:00Z",
    "merchant": "Restaurant Name"
  }'
```

### Create User (Admin)

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "email": "manager@company.com",
    "password": "Manager@123",
    "first_name": "Jane",
    "last_name": "Smith",
    "role": "manager"
  }'
```

### Approve Expense

```bash
curl -X POST http://localhost:8080/api/v1/approvals/APPROVAL_ID/approve \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "comments": "Approved - valid expense"
  }'
```

### Upload Receipt (OCR)

```bash
curl -X POST http://localhost:8080/api/v1/ocr/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "receipt=@/path/to/receipt.jpg" \
  -F "create_expense=true"
```

### Get Pending Approvals

```bash
curl http://localhost:8080/api/v1/approvals/pending \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

---

## 📊 Response Format

### Success (200/201)

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

### Success with Pagination

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

### Error (4xx/5xx)

```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error"
}
```

---

## 🔍 HTTP Status Codes

| Code | Status         | Meaning                  |
| ---- | -------------- | ------------------------ |
| 200  | OK             | Success                  |
| 201  | Created        | Resource created         |
| 400  | Bad Request    | Invalid input            |
| 401  | Unauthorized   | Invalid/missing token    |
| 403  | Forbidden      | Insufficient permissions |
| 404  | Not Found      | Resource not found       |
| 500  | Internal Error | Server error             |

---

## ⚡ Cache TTL

| Data Type         | TTL          | Invalidation          |
| ----------------- | ------------ | --------------------- |
| Auth Session      | Token expiry | Logout                |
| User List         | 15 min       | User create/update    |
| Expense List      | 15 min       | Expense create/update |
| Pending Approvals | 5 min        | Approval action       |
| Currency Rates    | 1 hour       | -                     |
| OCR Results       | 24 hours     | -                     |

---

## 🔐 Password Requirements

- Minimum 6 characters
- At least 1 uppercase letter
- At least 1 number
- Special characters allowed

---

## 🌍 Multi-Currency

### Auto-Conversion

- Expenses auto-convert to company base currency
- Uses ExchangeRate API
- Rates cached for 1 hour

### Currency Detection

- Company currency auto-detected from country code
- Uses RestCountries API

---

## 📝 Validation Rules

### Email

- Valid email format
- Unique in system

### Amount

- Greater than 0
- Up to 2 decimal places

### Currency

- Valid ISO 4217 code
- 3 letters uppercase

### Category

- Must be from allowed list

---

## 🛠️ Troubleshooting

### "Unauthorized" Error

- Check token is valid
- Token might be expired - use refresh
- Ensure `Authorization: Bearer TOKEN` header

### "Forbidden" Error

- Check your role permissions
- Admin: Full access
- Manager: Approvals + users
- Employee: Own expenses only

### "Bad Request" Error

- Check required fields
- Validate data formats
- Review error message

---

## 📚 Full Documentation

- **API Docs**: `API_DOCS.md`
- **Postman Guide**: `POSTMAN_GUIDE.md`
- **Setup Guide**: `SETUP_GUIDE.md`
- **Project Summary**: `PROJECT_SUMMARY.md`

---

## 🎯 Common Workflows

### 1. Employee Submits Expense

```
Login → Create Expense → View Status
```

### 2. Manager Approves

```
Login → View Pending → Approve/Reject
```

### 3. Admin Creates Team

```
Login → Create Manager → Create Employees → Assign Managers
```

### 4. OCR Quick Submit

```
Login → Upload Receipt → Auto-Create Expense → Done
```

---

## 💡 Pro Tips

1. **Use Postman Collection** - Automatic token management
2. **Cache Tokens** - Save refresh token securely
3. **Pagination** - Use `page` and `limit` parameters
4. **Filters** - Filter expenses by status, date, etc.
5. **Batch Operations** - Create multiple users efficiently

---

**Version:** 1.0.0  
**Last Updated:** October 4, 2025  
**Support:** Check documentation files for details
