# 📬 Postman Collection - Import Instructions

## 📦 What's Included

✅ **Expensio.postman_collection.json** - Complete API collection (25+ endpoints)  
✅ **Expensio.postman_environment.json** - Local environment configuration  
✅ **POSTMAN_GUIDE.md** - Comprehensive usage guide  
✅ **API_QUICK_REFERENCE.md** - Quick reference for curl commands

---

## 🚀 Import in 3 Steps

### Step 1: Open Postman

Launch Postman Desktop App or open [Postman Web](https://web.postman.co/)

### Step 2: Import Collection

1. Click **Import** button (top left corner)
2. Drag and drop these 2 files:
   - `Expensio.postman_collection.json`
   - `Expensio.postman_environment.json`
3. Click **Import**

### Step 3: Select Environment

1. Click environment dropdown (top right)
2. Select **"Expensio Local Environment"**

---

## ✅ Verify Import

You should see:

📁 **Collections Panel (left side):**

```
Expensio Backend API
├── Health Check
├── Authentication (4 endpoints)
├── User Management (6 endpoints)
├── Expense Management (6 endpoints)
├── Approval Workflow (4 endpoints)
└── OCR Receipt Processing (1 endpoint)
```

🌍 **Environment (top right):**

```
Expensio Local Environment ✓
```

---

## 🎯 First Test

### 1. Start Backend Server

```bash
cd d:\Expensio\backend
go run cmd/server/main.go
```

### 2. Run Health Check

- Open: **Health Check** → **Health Check**
- Click **Send**
- Should return: `"status": "ok"`

### 3. Create Admin Account

- Open: **Authentication** → **Signup (Create Company & Admin)**
- Update email/password if needed
- Click **Send**
- ✅ Tokens automatically saved to environment!

### 4. Test Protected Endpoint

- Open: **User Management** → **List Users**
- Click **Send**
- Should return your admin user (token auto-applied!)

---

## 🔧 Environment Variables

Auto-configured in environment:

| Variable       | Value                   | Description       |
| -------------- | ----------------------- | ----------------- |
| `baseUrl`      | `http://localhost:8080` | API base URL      |
| `accessToken`  | (auto-set on login)     | JWT access token  |
| `refreshToken` | (auto-set on login)     | JWT refresh token |
| `userId`       | (auto-set on login)     | Your user ID      |
| `companyId`    | (auto-set on login)     | Your company ID   |

---

## 📚 Next Steps

1. ✅ **Read**: `POSTMAN_GUIDE.md` for detailed usage
2. ✅ **Test**: Run through all folders systematically
3. ✅ **Experiment**: Modify requests and see responses
4. ✅ **Reference**: Use `API_QUICK_REFERENCE.md` for curl examples

---

## 💡 Pro Tips

### Automatic Token Management

The collection includes test scripts that automatically:

- Extract tokens from login/signup responses
- Save to environment variables
- Apply to all authenticated requests
- Clear on logout

**You just login once and forget about tokens!** 🎉

### Check Console

- Open: **View** → **Show Postman Console**
- See: Token save confirmations and request details

### Run Entire Collection

- Click: **Collections** → **Expensio Backend API** → **⋮** → **Run collection**
- Automated testing of all endpoints!

---

## 🐛 Troubleshooting

### Collection Not Visible?

- Check **Collections** panel on left
- Try searching "Expensio" in search bar

### Environment Not Available?

- Check top-right environment dropdown
- Click **Environments** tab to verify import

### Requests Failing?

1. Ensure backend is running: `http://localhost:8080/health`
2. Check environment is selected (top right)
3. Try login again to refresh tokens

---

## 🎨 Collection Features

✨ **Organized Folders** - Grouped by feature  
🔐 **Auto Authentication** - Bearer token auto-applied  
📝 **Pre-filled Examples** - Ready-to-use request bodies  
🧪 **Test Scripts** - Automatic token extraction  
📖 **Descriptions** - Each request documented  
🔄 **Variables** - Reusable across requests

---

## 📞 Need Help?

- **API Details**: See `API_DOCS.md`
- **Postman Usage**: See `POSTMAN_GUIDE.md`
- **Setup Issues**: See `SETUP_GUIDE.md`
- **Quick Commands**: See `API_QUICK_REFERENCE.md`

---

**You're all set! Start testing the Expensio API! 🚀**

**Happy Testing!** 🎉
