# 08 - Security Checklist

> Security requirements and OWASP compliance for VacayTracker

## Table of Contents

1. [Authentication Security](#authentication-security)
2. [Authorization Security](#authorization-security)
3. [Input Validation](#input-validation)
4. [Data Protection](#data-protection)
5. [API Security](#api-security)
6. [Frontend Security](#frontend-security)
7. [Infrastructure Security](#infrastructure-security)
8. [OWASP Top 10 Compliance](#owasp-top-10-compliance)

---

## Authentication Security

### Password Requirements

| Requirement | Implementation | Status |
|-------------|----------------|--------|
| **Minimum length** | 6 characters | - [ ] |
| **Maximum length** | 72 characters (bcrypt limit) | - [ ] |
| **Hashing algorithm** | bcrypt | - [ ] |
| **Cost factor** | 10 rounds | - [ ] |
| **No plaintext storage** | Only hash stored in DB | - [ ] |

### Implementation

```go
// internal/service/auth.go
func (s *AuthService) HashPassword(password string) (string, error) {
    // Validate length before hashing
    if len(password) < 6 {
        return "", errors.New("password must be at least 6 characters")
    }
    if len(password) > 72 {
        return "", errors.New("password cannot exceed 72 characters")
    }

    // bcrypt with cost 10
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    return string(bytes), err
}
```

### JWT Token Security

| Requirement | Implementation | Status |
|-------------|----------------|--------|
| **Algorithm** | HS256 (HMAC-SHA256) | - [ ] |
| **Secret length** | Minimum 32 characters | - [ ] |
| **Token expiry** | 24 hours | - [ ] |
| **Claims validation** | Expiry, issuer, subject | - [ ] |
| **No sensitive data in payload** | Only ID, email, role | - [ ] |

### Implementation

```go
// Validate JWT secret at startup
func Load() *Config {
    cfg := &Config{
        JWTSecret: mustGetEnv("JWT_SECRET"),
    }

    if len(cfg.JWTSecret) < 32 {
        log.Fatal("JWT_SECRET must be at least 32 characters")
    }

    return cfg
}

// Generate token with minimal claims
func (s *AuthService) GenerateToken(user *domain.User) (string, error) {
    claims := jwt.MapClaims{
        "sub":   user.ID,           // Subject (user ID)
        "email": user.Email,         // Email for display
        "role":  user.Role,          // Role for authorization
        "exp":   time.Now().Add(24 * time.Hour).Unix(),
        "iat":   time.Now().Unix(),
    }
    // Note: NO password hash, balance, or other sensitive data

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.cfg.JWTSecret))
}
```

### Checklist

- [ ] JWT secret is at least 32 characters
- [ ] JWT secret is randomly generated (not hardcoded)
- [ ] JWT secret is stored in environment variable
- [ ] Token expiry is enforced (24h)
- [ ] Invalid tokens return 401 Unauthorized
- [ ] Expired tokens return 401 Unauthorized
- [ ] Token validation checks signature algorithm

---

## Authorization Security

### Role-Based Access Control

| Role | Permissions |
|------|-------------|
| **Employee** | View own requests, create requests, cancel own pending requests, view team calendar |
| **Admin** | All employee permissions + manage all users + approve/reject requests + manage settings |

### Route Protection

| Route Pattern | Required Role | Middleware |
|---------------|---------------|------------|
| `/api/auth/login` | None | - |
| `/api/auth/*` (authenticated) | Any authenticated | RequireAuth |
| `/api/vacation/*` | Employee or Admin | RequireAuth |
| `/api/admin/*` | Admin only | RequireAuth + RequireAdmin |

### Implementation

```go
// Middleware chain for admin routes
admin := api.Group("/admin")
admin.Use(m.Auth.RequireAuth())   // First: validate JWT
admin.Use(m.Auth.RequireAdmin())  // Second: check role

// RequireAdmin middleware
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("userRole")
        if !exists {
            c.AbortWithStatusJSON(401, dto.ErrorResponse{
                Code:    "AUTH_TOKEN_MISSING",
                Message: "Authentication required",
            })
            return
        }

        if role.(domain.Role) != domain.RoleAdmin {
            c.AbortWithStatusJSON(403, dto.ErrorResponse{
                Code:    "ADMIN_REQUIRED",
                Message: "Admin role required",
            })
            return
        }

        c.Next()
    }
}
```

### Business Rule Enforcement

| Rule | Implementation | Status |
|------|----------------|--------|
| **Cannot delete self** | Check `userID != currentUserID` | - [ ] |
| **Cannot delete last admin** | Check admin count > 1 | - [ ] |
| **Cannot demote last admin** | Check admin count > 1 | - [ ] |
| **Cannot modify own role** | Check `userID != currentUserID` | - [ ] |
| **Employee can only cancel own requests** | Check `request.UserID == userID` | - [ ] |
| **Employee can only view own requests** | Filter by `userID` | - [ ] |

### Checklist

- [ ] All protected routes use RequireAuth middleware
- [ ] Admin routes use RequireAdmin middleware
- [ ] Resource ownership is verified before modification
- [ ] Business rules prevent privilege escalation
- [ ] Error messages don't reveal sensitive information

---

## Input Validation

### Request Validation

| Field | Validations | Status |
|-------|-------------|--------|
| **Email** | Required, valid format, unique | - [ ] |
| **Password** | Required, 6-72 chars | - [ ] |
| **Name** | Required, 1-100 chars | - [ ] |
| **Role** | Required, enum (admin/employee) | - [ ] |
| **Start Date** | Required, DD/MM/YYYY format, not in past | - [ ] |
| **End Date** | Required, DD/MM/YYYY format, >= start date | - [ ] |
| **Reason** | Optional, max 500 chars | - [ ] |

### Implementation (Go Struct Tags)

```go
// internal/dto/request.go
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type CreateUserRequest struct {
    Email           string  `json:"email" binding:"required,email"`
    Password        string  `json:"password" binding:"required,min=6,max=72"`
    Name            string  `json:"name" binding:"required,min=1,max=100"`
    Role            string  `json:"role" binding:"required,oneof=admin employee"`
    VacationBalance *int    `json:"vacationBalance" binding:"omitempty,min=0"`
    StartDate       *string `json:"startDate"` // Custom validation
}
```

### SQL Injection Prevention

```go
// SAFE: Using parameterized queries
query := "SELECT * FROM users WHERE email = ?"
db.QueryRow(query, email)

// UNSAFE: String concatenation (NEVER DO THIS)
query := "SELECT * FROM users WHERE email = '" + email + "'"  // SQL INJECTION!
```

### Checklist

- [ ] All user inputs are validated
- [ ] SQL queries use parameterized statements
- [ ] Email format is validated
- [ ] Date formats are validated
- [ ] String lengths are enforced
- [ ] Enum values are validated
- [ ] Numeric ranges are validated

---

## Data Protection

### Sensitive Data Handling

| Data Type | Storage | Transmission | Display |
|-----------|---------|--------------|---------|
| **Passwords** | bcrypt hash only | HTTPS only | Never displayed |
| **JWT Secret** | Environment variable | N/A | Never exposed |
| **Email** | Plaintext (DB) | HTTPS | Visible to admin |
| **Vacation balance** | Plaintext (DB) | HTTPS | Visible to user/admin |

### Data Masking

```go
// User domain - password hash never serialized
type User struct {
    ID           string `json:"id"`
    Email        string `json:"email"`
    PasswordHash string `json:"-"` // NEVER exposed in JSON
    Name         string `json:"name"`
    // ...
}
```

### Environment Security

```bash
# Required secrets (generate securely)
JWT_SECRET=<random-32-plus-chars>
ADMIN_PASSWORD=<secure-password>

# Generate secure secret:
openssl rand -base64 32
```

### Checklist

- [ ] Password hashes never exposed in API responses
- [ ] JWT secret stored in environment variable
- [ ] .env file in .gitignore
- [ ] Database file has appropriate permissions
- [ ] Sensitive logs are filtered
- [ ] Production secrets differ from development

---

## API Security

### CORS Configuration

```go
// internal/middleware/cors.go
func CORS(allowOrigin string) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")

        // Only allow specific origins
        allowedOrigins := map[string]bool{
            allowOrigin:               true,
            "http://localhost:5173":   true, // Dev frontend
            "http://localhost:3000":   true, // Dev backend
        }

        if allowedOrigins[origin] {
            c.Header("Access-Control-Allow-Origin", origin)
        }

        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Max-Age", "86400")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

### Rate Limiting (Production)

```go
// Implement for production
// Consider: golang.org/x/time/rate or gin-contrib/limiter

// Per-IP limits
// - Login: 5 attempts per minute
// - API: 100 requests per minute
// - Admin actions: 30 per minute
```

### Security Headers

```go
// Add in production
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Next()
    }
}
```

### Checklist

- [ ] CORS configured with specific origins
- [ ] OPTIONS requests handled for preflight
- [ ] Rate limiting implemented (production)
- [ ] Security headers set (production)
- [ ] HTTPS enforced (production)

---

## Frontend Security

### XSS Prevention

```svelte
<!-- Svelte automatically escapes by default -->
<p>{userInput}</p>  <!-- SAFE: escaped -->

<!-- Only use @html with sanitized content -->
{@html sanitizedHtml}  <!-- DANGEROUS if not sanitized -->
```

### Token Storage

```typescript
// Store JWT in localStorage (acceptable for this app)
// For higher security apps, consider httpOnly cookies

export function setAuthToken(token: string): void {
    localStorage.setItem('auth_token', token);
}

export function clearAuthToken(): void {
    localStorage.removeItem('auth_token');
}

// Token is sent via Authorization header (not cookies)
// This prevents CSRF attacks
```

### Input Sanitization

```typescript
// Validate on frontend before sending
function validateEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

function validatePassword(password: string): boolean {
    return password.length >= 6 && password.length <= 72;
}

function validateDateRange(start: string, end: string): boolean {
    const startDate = new Date(start);
    const endDate = new Date(end);
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    return startDate >= today && endDate >= startDate;
}
```

### Checklist

- [ ] All user inputs validated before API calls
- [ ] Svelte's default escaping used (no raw HTML)
- [ ] JWT stored in localStorage, sent via header
- [ ] Form validation provides user feedback
- [ ] Error messages don't expose sensitive info

---

## Infrastructure Security

### SQLite Security

```go
// Enable foreign keys and WAL mode
dsn := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)", dbPath)

// File permissions (Unix)
// Database file: 600 (rw-------)
// Data directory: 700 (rwx------)
```

### Environment Variables

| Variable | Required | Security Level |
|----------|----------|----------------|
| `JWT_SECRET` | Yes | High - must be random, 32+ chars |
| `ADMIN_PASSWORD` | Yes | High - initial admin password |
| `RESEND_API_KEY` | No | Medium - email service key |
| `DB_PATH` | No | Low - file path |
| `PORT` | No | Low - server port |

### Production Checklist

- [ ] HTTPS enforced (TLS 1.2+)
- [ ] JWT_SECRET is cryptographically random
- [ ] Database file has restricted permissions
- [ ] Logs don't contain sensitive data
- [ ] Error pages don't leak stack traces
- [ ] .env files are not deployed
- [ ] Docker runs as non-root user

---

## OWASP Top 10 Compliance

### A01:2021 - Broken Access Control

| Control | Implementation | Status |
|---------|----------------|--------|
| Deny by default | All routes require auth except login | - [ ] |
| Enforce ownership | Users can only access own resources | - [ ] |
| Rate limiting | Prevent brute force (production) | - [ ] |
| CORS | Restrict to known origins | - [ ] |

### A02:2021 - Cryptographic Failures

| Control | Implementation | Status |
|---------|----------------|--------|
| Password hashing | bcrypt with cost 10 | - [ ] |
| TLS in transit | HTTPS enforced (production) | - [ ] |
| Secure secrets | Environment variables | - [ ] |

### A03:2021 - Injection

| Control | Implementation | Status |
|---------|----------------|--------|
| SQL injection | Parameterized queries | - [ ] |
| NoSQL injection | N/A (SQLite) | N/A |
| Command injection | No shell commands | N/A |

### A04:2021 - Insecure Design

| Control | Implementation | Status |
|---------|----------------|--------|
| Threat modeling | Business rules enforced | - [ ] |
| Secure defaults | Auth required by default | - [ ] |
| Input validation | All inputs validated | - [ ] |

### A05:2021 - Security Misconfiguration

| Control | Implementation | Status |
|---------|----------------|--------|
| Security headers | Set in production | - [ ] |
| Error handling | Generic errors to client | - [ ] |
| Default credentials | Changed on first run | - [ ] |

### A06:2021 - Vulnerable Components

| Control | Implementation | Status |
|---------|----------------|--------|
| Dependency audit | `go mod verify`, `npm audit` | - [ ] |
| Update policy | Regular dependency updates | - [ ] |

### A07:2021 - Authentication Failures

| Control | Implementation | Status |
|---------|----------------|--------|
| Password policy | 6+ chars required | - [ ] |
| Brute force protection | Rate limiting (production) | - [ ] |
| Session management | JWT with expiry | - [ ] |

### A08:2021 - Software and Data Integrity

| Control | Implementation | Status |
|---------|----------------|--------|
| CI/CD security | GitHub Actions with secrets | - [ ] |
| Dependency integrity | Lock files committed | - [ ] |

### A09:2021 - Security Logging

| Control | Implementation | Status |
|---------|----------------|--------|
| Login attempts | Logged with timestamp | - [ ] |
| Failed auth | Logged without password | - [ ] |
| Admin actions | Logged with user ID | - [ ] |

### A10:2021 - Server-Side Request Forgery

| Control | Implementation | Status |
|---------|----------------|--------|
| External requests | Only to Resend API | - [ ] |
| URL validation | N/A (no user-supplied URLs) | N/A |

---

## Security Review Checklist

### Pre-Deployment

- [ ] All secrets are environment variables
- [ ] JWT secret is 32+ characters
- [ ] Password hashing uses bcrypt
- [ ] All routes have appropriate auth
- [ ] SQL queries are parameterized
- [ ] CORS is properly configured
- [ ] Input validation is complete
- [ ] Error messages are generic

### Post-Deployment

- [ ] HTTPS is enforced
- [ ] Security headers are set
- [ ] Rate limiting is active
- [ ] Logs are monitored
- [ ] Dependencies are up to date
- [ ] Backup strategy is in place

---

## Related Documents

- [02-api-specification.md](./02-api-specification.md) - API authentication details
- [04-backend-tasks.md](./04-backend-tasks.md) - Security implementation
- [09-deployment-guide.md](./09-deployment-guide.md) - Production security
