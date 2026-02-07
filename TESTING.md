# Testing Guide for Event-Driven Auth Service

This guide provides step-by-step instructions for testing the authentication service, including Docker setup, API testing, and verification procedures.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Docker Setup & Service Launch](#docker-setup--service-launch)
3. [Database Initialization](#database-initialization)
4. [Authentication API Testing](#authentication-api-testing)
5. [Token Management Testing](#token-management-testing)
6. [Error Scenario Testing](#error-scenario-testing)
7. [Health Check Testing](#health-check-testing)
8. [Database Verification](#database-verification)
9. [Testing Checklist](#testing-checklist)
10. [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Software
- **Docker** & **Docker Compose** (latest versions)
- **Go 1.25.5+** (for local development)
- **curl** or **Postman** (for API testing)
- **psql** (PostgreSQL client - optional)

### Environment Setup
1. Clone the repository:
```bash
git clone <repository-url>
cd event-driven
```

2. Copy environment template:
```bash
cp services/auth/.env.example .env
```

3. Set your JWT secret in `.env`:
```bash
echo "JWT_SECRET=my-super-secret-jwt-key-for-testing-$(date +%s)" >> .env
```

## Docker Setup & Service Launch

### Step 1: Start Services
```bash
# From project root directory
docker-compose up -d
```

### Step 2: Verify Services are Running
```bash
# Check all containers
docker-compose ps

# Expected output:
# NAME              COMMAND                  SERVICE             STATUS              PORTS
# auth-postgres     "docker-entrypoint.s…"   postgres            running             0.0.0.0:5432->5432/tcp
# auth-service      "/app/auth-service"      auth                running             0.0.0.0:8081->8080/tcp
```

### Step 3: Check Service Health
```bash
# Wait 10-15 seconds for services to fully start
sleep 15

# Check auth service health
curl http://localhost:8081/health

# Expected response: OK
```

### Step 4: View Service Logs (if needed)
```bash
# View auth service logs
docker-compose logs auth

# View database logs
docker-compose logs postgres
```

## Database Initialization

### Step 1: Wait for Database Ready
```bash
# Check if database is ready
docker exec auth-postgres pg_isready -U authuser -d authdb

# Expected output: accepting connections
```

### Step 2: Run Database Migrations
```bash
# Navigate to auth service directory
cd services/auth

# Create admin user (this also runs migrations automatically)
go run scripts/create_admin.go

# Expected output:
# ✅ Admin user created successfully
# Username: admin
# Password: admin
```

### Step 3: Verify Database Tables
```bash
# Connect to database
docker exec -it auth-postgres psql -U authuser -d authdb

# List tables
\dt

# Expected output:
# List of relations
#  Schema |    Name             | Type  | Owner
# --------+---------------------+-------+---------
#  public | refresh_tokens      | table | authuser
#  public | users               | table | authuser

# Verify admin user exists
SELECT id, username, email, created_at FROM users WHERE username = 'admin';

# Exit database
\q
```

## Authentication API Testing

### Service Endpoints
- **Base URL**: `http://localhost:8081`
- **API Version**: `/api/v1`

### 1. Login Testing

#### Successful Login
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin"
  }'
```

**Expected Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Save the tokens for next tests:**
```bash
# Extract and save tokens (Linux/Mac)
ACCESS_TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin"}' | \
  jq -r '.access_token')

REFRESH_TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin"}' | \
  jq -r '.refresh_token')

echo "Access Token: $ACCESS_TOKEN"
echo "Refresh Token: $REFRESH_TOKEN"
```

#### Failed Login - Wrong Password
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "wrongpassword"
  }'
```

**Expected Response (401 Unauthorized):**
```text
unauthorized
```

#### Failed Login - Non-existent User
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "nonexistent",
    "password": "password"
  }'
```

**Expected Response (401 Unauthorized):**
```text
unauthorized
```

### 2. Protected Endpoint Testing

#### Access Protected Resource (Valid Token)
```bash
# Use the access token from login
curl -X GET http://localhost:8081/api/v1/protected \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Expected Response (200 OK):**
```text
protected content
```

#### Access Protected Resource (No Token)
```bash
curl -X GET http://localhost:8081/api/v1/protected
```

**Expected Response (401 Unauthorized):**
```text
missing token
```

#### Access Protected Resource (Invalid Token)
```bash
curl -X GET http://localhost:8081/api/v1/protected \
  -H "Authorization: Bearer invalid.token.here"
```

**Expected Response (401 Unauthorized):**
```text
invalid token
```

## Token Management Testing

### 1. Token Refresh Testing

#### Successful Token Refresh
```bash
curl -X POST http://localhost:8081/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }"
```

**Expected Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Failed Token Refresh - Invalid Token
```bash
curl -X POST http://localhost:8081/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "invalid.refresh.token"
  }'
```

**Expected Response (401 Unauthorized):**
```text
unauthorized
```

#### Failed Token Refresh - Missing Token
```bash
curl -X POST http://localhost:8081/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Expected Response (401 Unauthorized):**
```text
unauthorized
```

### 2. Logout Testing

#### Successful Logout
```bash
curl -X POST http://localhost:8081/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }"
```

**Expected Response (204 No Content):**
- No response body

#### Failed Logout - Invalid Token
```bash
curl -X POST http://localhost:8081/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "invalid.token"
  }'
```

**Expected Response (400 Bad Request):**
```text
failed to logout
```

## Error Scenario Testing

### 1. Invalid JSON Input
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin"'  # Invalid JSON
```

**Expected Response (400 Bad Request):**

### 2. Missing Required Fields
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin"}'  # Missing password
```

**Expected Response (401 Unauthorized):**

### 3. Empty Request Body
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Expected Response (401 Unauthorized):**

## Health Check Testing

### Service Health Check
```bash
curl -i http://localhost:8081/health
```

**Expected Response (200 OK):**
```http
HTTP/1.1 200 OK
Date: [current date]
Content-Length: 2
Content-Type: text/plain

OK
```

### Database Connection Test
```bash
# Stop database to test failure scenario
docker-compose stop postgres

# Check health endpoint
curl -i http://localhost:8081/health
```

**Expected Response (503 Service Unavailable):**
```http
HTTP/1.1 503 Service Unavailable
```

```bash
# Restart database
docker-compose start postgres
```

## Database Verification

### Check Users Table
```bash
docker exec -it auth-postgres psql -U authuser -d authdb -c "
SELECT 
  id,
  username,
  email,
  created_at,
  updated_at
FROM users;
"
```

### Check Refresh Tokens Table
```bash
docker exec -it auth-postgres psql -U authuser -d authdb -c "
SELECT 
  id,
  user_id,
  token,
  expires_at,
  created_at,
  revoked_at
FROM refresh_tokens;
"
```

### Verify Token Storage After Login
```bash
# Login to generate token
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin"}' > /dev/null

# Check if token was stored
docker exec auth-postgres psql -U authuser -d authdb -c "
SELECT COUNT(*) as token_count 
FROM refresh_tokens 
WHERE user_id = 1 AND revoked_at IS NULL;
"
```

### Verify Token Revocation After Logout
```bash
# Logout to revoke token
curl -X POST http://localhost:8081/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }" > /dev/null

# Check if token was revoked
docker exec auth-postgres psql -U authuser -d authdb -c "
SELECT 
  COUNT(*) as revoked_count,
  revoked_at IS NOT NULL as is_revoked
FROM refresh_tokens 
WHERE user_id = 1 
GROUP BY revoked_at IS NOT NULL;
"
```

## Testing Checklist

### ✅ Pre-Setup Checklist
- [ ] Docker and Docker Compose installed
- [ ] Environment variables configured
- [ ] JWT secret set in `.env`
- [ ] Repository cloned to local machine

### ✅ Service Startup Checklist
- [ ] `docker-compose up -d` completes successfully
- [ ] Both containers are running (`docker-compose ps`)
- [ ] Health endpoint returns `OK`
- [ ] Database is accepting connections
- [ ] Admin user created successfully

### ✅ Authentication Testing Checklist
- [ ] Login with valid credentials succeeds
- [ ] Login with wrong password fails (401)
- [ ] Login with non-existent user fails (401)
- [ ] Access token is returned in login response
- [ ] Refresh token is returned in login response
- [ ] Protected endpoint accessible with valid token
- [ ] Protected endpoint inaccessible without token
- [ ] Protected endpoint inaccessible with invalid token

### ✅ Token Management Checklist
- [ ] Token refresh works with valid refresh token
- [ ] Token refresh fails with invalid token
- [ ] Logout succeeds with valid refresh token
- [ ] Logout fails with invalid token
- [ ] Refresh tokens stored in database after login
- [ ] Refresh tokens marked as revoked after logout
- [ ] New access tokens generated after refresh

### ✅ Error Handling Checklist
- [ ] Invalid JSON returns appropriate error
- [ ] Missing fields return appropriate error
- [ ] Empty request body returns appropriate error
- [ ] Database connection failure reflected in health check

### ✅ Database Verification Checklist
- [ ] Users table created correctly
- [ ] Refresh tokens table created correctly
- [ ] Admin user exists in database
- [ ] Passwords are hashed (not plain text)
- [ ] Tokens stored with correct expiration
- [ ] Tokens marked as revoked after logout

## Troubleshooting

### Common Issues and Solutions

#### 1. Service Won't Start
**Problem**: `docker-compose up` fails or containers restart repeatedly
**Solutions**:
```bash
# Check logs
docker-compose logs auth
docker-compose logs postgres

# Rebuild auth service
docker-compose down
docker-compose up --build -d

# Check port conflicts
netstat -tulpn | grep :8081
netstat -tulpn | grep :5432
```

#### 2. Database Connection Issues
**Problem**: Authentication service can't connect to database
**Solutions**:
```bash
# Check database container status
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Manually test database connection
docker exec -it auth-postgres psql -U authuser -d authdb -c "SELECT 1;"

# Restart database
docker-compose restart postgres
```

#### 3. Login Always Fails
**Problem**: Login returns 401 even with correct credentials
**Solutions**:
```bash
# Verify admin user exists
docker exec auth-postgres psql -U authuser -d authdb -c "SELECT username FROM users WHERE username = 'admin';"

# If user doesn't exist, recreate admin user
cd services/auth && go run scripts/create_admin.go

# Check JWT secret consistency
grep JWT_SECRET .env
docker-compose exec auth printenv | grep JWT_SECRET
```

#### 4. Token Issues
**Problem**: Token validation fails or refresh doesn't work
**Solutions**:
```bash
# Check if tokens are being stored
docker exec auth-postgres psql -U authuser -d authdb -c "SELECT COUNT(*) FROM refresh_tokens;"

# Verify JWT secret is consistent between requests
# (should be the same for entire session)

# Check token expiration
docker exec auth-postgres psql -U authuser -d authdb -c "SELECT expires_at FROM refresh_tokens WHERE user_id = 1;"
```

#### 5. Port Conflicts
**Problem**: Services can't bind to their ports
**Solutions**:
```bash
# Find what's using the ports
sudo lsof -i :8081
sudo lsof -i :5432

# Kill conflicting processes or change ports in docker-compose.yml
```

#### 6. Permission Issues
**Problem**: Permission denied errors
**Solutions**:
```bash
# Check Docker permissions
sudo usermod -aG docker $USER
# Logout and login again

# Or run with sudo (not recommended for development)
sudo docker-compose up -d
```

### Getting Help

If you encounter issues not covered here:

1. **Check logs first**: `docker-compose logs [service-name]`
2. **Verify environment**: Ensure `.env` file exists with correct values
3. **Check networking**: Ensure containers can communicate
4. **Restart services**: Sometimes a simple restart fixes issues
5. **Rebuild from scratch**: 
   ```bash
   docker-compose down -v
   docker-compose up --build -d
   ```

## Advanced Testing

### Load Testing (Optional)
For basic load testing, you can use:
```bash
# Install Apache Bench
sudo apt-get install apache2-utils

# Run load test on login endpoint (100 requests, 10 concurrent)
ab -n 100 -c 10 -p login_data.json -T application/json http://localhost:8081/api/v1/auth/login
```

Create `login_data.json`:
```json
{"username": "admin", "password": "admin"}
```

### Security Testing (Optional)
For security testing, consider:
- SQL injection attempts
- JWT token manipulation
- Brute force attack simulation
- Input validation testing

---

## Next Steps

After completing all manual testing, consider implementing:

1. **Unit Tests**: For individual functions and methods
2. **Integration Tests**: For database interactions
3. **End-to-End Tests**: For complete API flows
4. **CI/CD Pipeline**: For automated testing
5. **Performance Tests**: For load and stress testing

This testing guide ensures your authentication service is working correctly and is ready for production use.