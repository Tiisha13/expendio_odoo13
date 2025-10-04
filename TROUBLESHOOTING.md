# Troubleshooting Guide - Expensio

Common issues and solutions for Expensio expense management system.

## Table of Contents

1. [Backend Issues](#backend-issues)
2. [Frontend Issues](#frontend-issues)
3. [Database Issues](#database-issues)
4. [Authentication Issues](#authentication-issues)
5. [API Issues](#api-issues)
6. [OCR Issues](#ocr-issues)
7. [Performance Issues](#performance-issues)

---

## Backend Issues

### Backend won't start

**Error:** "Failed to connect to MongoDB"

**Solution:**

```bash
# Check if MongoDB is running
docker ps | grep mongo

# Start MongoDB if not running
docker start expensio-mongo

# Or start new instance
docker run -d -p 27017:27017 --name expensio-mongo mongo:7.0

# Verify connection
mongosh --host localhost --port 27017
```

**Error:** "Failed to connect to Redis"

**Solution:**

```bash
# Check if Redis is running
docker ps | grep redis

# Start Redis
docker start expensio-redis

# Or start new instance
docker run -d -p 6379:6379 --name expensio-redis redis:7.0

# Test connection
redis-cli ping
```

**Error:** "Port 8080 already in use"

**Solution:**

```bash
# Find process using port 8080
# Windows
netstat -ano | findstr :8080

# Linux/Mac
lsof -i :8080

# Kill the process or change PORT in .env
PORT=8081
```

### JWT Token Issues

**Error:** "Invalid token" or "Token expired"

**Solution:**

1. Check JWT_SECRET in `.env` matches
2. Verify token expiry times:
   ```env
   JWT_ACCESS_EXPIRY=15m
   JWT_REFRESH_EXPIRY=168h
   ```
3. Try logging in again
4. Clear browser cookies/storage

### CORS Errors

**Error:** "CORS policy blocked"

**Solution:**
Edit `backend/internal/config/config.go`:

```go
AllowOrigins: []string{
    "http://localhost:3000",
    "https://your-domain.com",
}
```

---

## Frontend Issues

### Frontend won't start

**Error:** "ERR_PNPM_NO_IMPORTER_MANIFEST_FOUND"

**Solution:**

```bash
# Make sure you're in the frontend directory
cd frontend

# Check package.json exists
ls package.json

# Reinstall dependencies
rm -rf node_modules pnpm-lock.yaml
pnpm install
```

**Error:** "Module not found" errors

**Solution:**

```bash
# Clear Next.js cache
rm -rf .next

# Reinstall dependencies
pnpm install

# Restart dev server
pnpm dev
```

**Error:** "Port 3000 already in use"

**Solution:**

```bash
# Windows: Find and kill process
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# Linux/Mac
lsof -ti:3000 | xargs kill -9

# Or use different port
pnpm dev -- -p 3001
```

### Build Errors

**Error:** TypeScript compilation errors

**Solution:**

```bash
# Check TypeScript version
pnpm list typescript

# Clear type cache
rm -rf .next
rm -rf node_modules/.cache

# Rebuild
pnpm build
```

**Error:** "Hydration failed"

**Solution:**

1. Check for mismatched HTML between server and client
2. Ensure no `localStorage` access during SSR
3. Use `useEffect` for client-only code
4. Check for missing `"use client"` directives

---

## Database Issues

### MongoDB Connection Issues

**Error:** "MongoServerError: Authentication failed"

**Solution:**

```bash
# Check credentials in .env
MONGODB_URI=mongodb://username:password@localhost:27017/expensio

# Or connect without auth for development
MONGODB_URI=mongodb://localhost:27017/expensio

# Test connection
mongosh "mongodb://localhost:27017/expensio"
```

**Error:** "Database not found"

**Solution:**
MongoDB creates databases automatically. Just verify:

```bash
mongosh
show dbs
use expensio
db.users.countDocuments()
```

### Redis Issues

**Error:** "NOAUTH Authentication required"

**Solution:**

```bash
# Set password in .env
REDIS_PASSWORD=your-password

# Or remove password requirement
redis-cli CONFIG SET requirepass ""
```

**Error:** "Connection timeout"

**Solution:**

```bash
# Check Redis is running
redis-cli ping

# Should return: PONG

# If not, restart Redis
docker restart expensio-redis
```

---

## Authentication Issues

### Can't Signup

**Error:** "Email already exists"

**Solution:**

```bash
# Delete existing user
mongosh
use expensio
db.users.deleteOne({email: "user@example.com"})

# Or use different email
```

**Error:** "Invalid country code"

**Solution:**
Use valid ISO country codes:

- US (United States)
- GB (United Kingdom)
- CA (Canada)
- AU (Australia)
- etc.

### Can't Login

**Error:** "Invalid credentials"

**Solution:**

1. Verify email and password are correct
2. Check user exists:
   ```bash
   mongosh
   use expensio
   db.users.findOne({email: "your@email.com"})
   ```
3. Reset password or create new account

**Error:** "NextAuth callback error"

**Solution:**

1. Check `NEXTAUTH_URL` in `.env.local`:
   ```env
   NEXTAUTH_URL=http://localhost:3000
   ```
2. Verify `NEXTAUTH_SECRET` is set
3. Clear browser cookies
4. Restart frontend

### Session Issues

**Error:** "Session expired" or "No session"

**Solution:**

1. Check token expiry times in backend
2. Verify NextAuth configuration
3. Clear browser storage:
   ```javascript
   localStorage.clear();
   sessionStorage.clear();
   ```
4. Login again

---

## API Issues

### 404 Not Found

**Error:** "Cannot GET /api/v1/..."

**Solution:**

1. Verify backend is running on correct port
2. Check API URL in frontend:
   ```env
   NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
   ```
3. Check route definition in `backend/internal/routes/routes.go`

### 401 Unauthorized

**Error:** "Unauthorized" or "Missing token"

**Solution:**

1. Verify user is logged in
2. Check Authorization header:
   ```javascript
   headers: {
     'Authorization': `Bearer ${token}`
   }
   ```
3. Token might be expired - login again
4. Check JWT middleware is configured

### 500 Internal Server Error

**Solution:**

1. Check backend logs:

   ```bash
   # If running with go run
   # Logs appear in terminal

   # If using systemd
   sudo journalctl -u expensio -n 50
   ```

2. Look for error messages
3. Check database connection
4. Verify all environment variables

### Network Errors

**Error:** "Failed to fetch" or "Network request failed"

**Solution:**

1. Verify backend is running
2. Check firewall/antivirus
3. Verify URLs are correct
4. Check CORS configuration
5. Try curl to test:
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

---

## OCR Issues

### Tesseract Not Found

**Error:** "Tesseract not installed"

**Solution:**

```bash
# Windows
choco install tesseract

# Linux
sudo apt install tesseract-ocr

# Mac
brew install tesseract

# Update .env
TESSERACT_PATH=/usr/bin/tesseract  # Linux
TESSERACT_PATH=C:\Program Files\Tesseract-OCR\tesseract.exe  # Windows
```

### OCR Not Extracting Data

**Issue:** Receipt uploaded but no data extracted

**Solution:**

1. Check image quality (clear, well-lit)
2. Supported formats: JPG, PNG, TIFF
3. Check Tesseract logs in backend
4. Try different image
5. Verify Tesseract language pack:
   ```bash
   tesseract --list-langs
   ```

### Upload Fails

**Error:** "File too large" or "Upload failed"

**Solution:**

1. Check file size (max 10MB typically)
2. Verify file format (image only)
3. Check server upload limits
4. Verify multipart/form-data headers

---

## Performance Issues

### Slow API Responses

**Solution:**

1. Check database indexes:
   ```javascript
   db.expenses.createIndex({ user_id: 1 });
   db.expenses.createIndex({ status: 1 });
   db.users.createIndex({ email: 1 });
   ```
2. Enable Redis caching
3. Check MongoDB query performance:
   ```javascript
   db.expenses.find({ user_id: "..." }).explain("executionStats");
   ```

### High Memory Usage

**Solution:**

1. Check for memory leaks
2. Limit query results (pagination)
3. Close database connections properly
4. Monitor with:

   ```bash
   # Backend
   go tool pprof http://localhost:8080/debug/pprof/heap

   # MongoDB
   db.serverStatus().mem
   ```

### Slow Page Load

**Solution:**

1. Enable Next.js caching
2. Optimize images
3. Lazy load components
4. Check network tab in browser
5. Reduce bundle size:
   ```bash
   pnpm analyze
   ```

---

## Common Error Messages

### "Cannot read property 'accessToken' of null"

**Cause:** Session not loaded yet

**Solution:**

```typescript
// Add loading check
if (!session?.accessToken) return null;

// Or use session status
const { data: session, status } = useSession();
if (status === "loading") return <Loading />;
```

### "Expense date is required"

**Cause:** Field name mismatch (date vs expense_date)

**Solution:**
Backend expects `expense_date`:

```json
{
  "expense_date": "2025-10-04"
}
```

### "Role must be admin, manager, or employee"

**Cause:** Invalid role value

**Solution:**
Use exact values:

```typescript
role: "admin" | "manager" | "employee";
```

---

## Debugging Tips

### Backend Debugging

1. **Enable verbose logging:**

   ```go
   log.SetLevel(log.DebugLevel)
   ```

2. **Check database queries:**

   ```go
   fmt.Printf("Query: %+v\n", filter)
   ```

3. **Test endpoints with curl:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@test.com","password":"pass"}'
   ```

### Frontend Debugging

1. **Check browser console:**

   - Press F12
   - Look for errors in Console tab
   - Check Network tab for failed requests

2. **React DevTools:**

   - Install React DevTools extension
   - Inspect component state
   - Check props

3. **Log API responses:**
   ```typescript
   console.log("API Response:", response);
   ```

### Database Debugging

1. **Check collections:**

   ```bash
   mongosh
   use expensio
   show collections
   db.users.countDocuments()
   db.expenses.find().limit(5)
   ```

2. **Check indexes:**

   ```bash
   db.expenses.getIndexes()
   ```

3. **Monitor queries:**
   ```bash
   db.setProfilingLevel(2)
   db.system.profile.find().sort({ts:-1}).limit(5)
   ```

---

## Getting Help

If issues persist:

1. **Check logs:**

   - Backend: Terminal output
   - Frontend: Browser console
   - MongoDB: `/var/log/mongodb/mongod.log`

2. **Verify environment:**

   - All .env variables set
   - Correct URLs and ports
   - Services running

3. **Test individually:**

   - Backend health: `curl http://localhost:8080/api/v1/health`
   - Frontend: Visit `http://localhost:3000`
   - Database: `mongosh`

4. **Common fixes:**

   - Restart all services
   - Clear caches
   - Reinstall dependencies
   - Check firewall/antivirus

5. **Documentation:**
   - Read README.md
   - Check QUICKSTART.md
   - Review API documentation

---

## Useful Commands

### Backend

```bash
# Restart backend
pkill -f "go run"
go run cmd/main.go

# Check Go version
go version

# Update dependencies
go mod tidy
```

### Frontend

```bash
# Clear cache and restart
rm -rf .next node_modules
pnpm install
pnpm dev

# Type check
pnpm tsc --noEmit

# Build test
pnpm build
```

### Database

```bash
# MongoDB
mongosh
show dbs
use expensio
db.stats()

# Redis
redis-cli
INFO
KEYS *
```

### Docker

```bash
# View logs
docker logs expensio-backend
docker logs expensio-frontend

# Restart containers
docker restart expensio-backend
docker restart expensio-mongo

# Remove and recreate
docker-compose down
docker-compose up -d
```

---

**Still stuck?** Create an issue on GitHub with:

- Error message (full text)
- Steps to reproduce
- Environment (OS, versions)
- Relevant logs
