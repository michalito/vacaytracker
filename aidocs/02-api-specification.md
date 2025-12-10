# 02 - API Specification

> Complete REST API documentation for VacayTracker backend

## Table of Contents

1. [Overview](#overview)
2. [Base Configuration](#base-configuration)
3. [Authentication](#authentication)
4. [Error Handling](#error-handling)
5. [Public Endpoints](#public-endpoints)
6. [Authenticated Endpoints](#authenticated-endpoints)
7. [Employee Endpoints](#employee-endpoints)
8. [Admin Endpoints](#admin-endpoints)
9. [Error Code Reference](#error-code-reference)

---

## Overview

### API Design Principles

| Principle | Implementation |
|-----------|----------------|
| **REST Convention** | Resource-based URLs, HTTP verbs for actions |
| **JSON** | All requests/responses use `application/json` |
| **Authentication** | JWT Bearer tokens in Authorization header |
| **Versioning** | No versioning (v1 implied), prefix `/api/` |
| **Error Format** | Consistent `{code, message, details?}` structure |

### Endpoint Summary

| Category | Count | Base Path |
|----------|-------|-----------|
| Public | 2 | `/`, `/api/auth/` |
| Authenticated | 3 | `/api/auth/` |
| Employee | 5 | `/api/vacation/` |
| Admin | 10 | `/api/admin/` |
| **Total** | **20** | - |

---

## Base Configuration

### Server Settings

```
Base URL:       http://localhost:3000 (development)
                https://vacaytracker.app (production)
Content-Type:   application/json
Accept:         application/json
```

### Request Headers

| Header | Required | Description |
|--------|----------|-------------|
| `Content-Type` | Yes (POST/PUT/PATCH) | Must be `application/json` |
| `Authorization` | Conditional | `Bearer <jwt_token>` for protected routes |
| `Accept` | No | Defaults to `application/json` |

### Response Headers

| Header | Description |
|--------|-------------|
| `Content-Type` | Always `application/json` |
| `X-Request-ID` | Unique request identifier for debugging |

---

## Authentication

### JWT Token Structure

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "user_uuid",
    "email": "user@example.com",
    "role": "employee|admin",
    "name": "User Name",
    "exp": 1234567890,
    "iat": 1234567890
  }
}
```

### Token Lifecycle

| Event | Duration |
|-------|----------|
| Token Expiry | 24 hours from issue |
| Refresh Strategy | Re-login required (no refresh tokens) |

### Authorization Header Format

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Role-Based Access

| Role | Access Level |
|------|--------------|
| `admin` | All endpoints |
| `employee` | Public + Authenticated + Employee endpoints |

---

## Error Handling

### Error Response Format

```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable error description",
  "details": {
    "field": "Additional context (optional)"
  }
}
```

### HTTP Status Codes

| Status | Usage |
|--------|-------|
| `200` | Success (GET, PUT, PATCH) |
| `201` | Created (POST) |
| `204` | No Content (DELETE) |
| `400` | Bad Request (validation errors) |
| `401` | Unauthorized (missing/invalid token) |
| `403` | Forbidden (insufficient permissions) |
| `404` | Not Found (resource doesn't exist) |
| `409` | Conflict (duplicate entry) |
| `422` | Unprocessable Entity (business rule violation) |
| `500` | Internal Server Error |

---

## Public Endpoints

### GET /health

Health check endpoint for monitoring and load balancers.

**Auth Required:** No

**Request:**
```http
GET /health HTTP/1.1
Host: localhost:3000
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

**Implementation File:** `internal/handler/health.go`

---

### POST /api/auth/login

Authenticate user and receive JWT token.

**Auth Required:** No

**Request:**
```http
POST /api/auth/login HTTP/1.1
Host: localhost:3000
Content-Type: application/json

{
  "email": "user@company.com",
  "password": "securePassword123"
}
```

**Request Schema:**

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `email` | string | Yes | Valid email format |
| `password` | string | Yes | Min 6 chars |

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@company.com",
    "name": "John Doe",
    "role": "employee",
    "vacationBalance": 25,
    "startDate": "2023-01-15",
    "emailPreferences": {
      "vacationUpdates": true,
      "weeklyDigest": false,
      "teamNotifications": true
    }
  }
}
```

**Response Schema (User Object):**

| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | User unique identifier |
| `email` | string | User email address |
| `name` | string | Display name |
| `role` | string | `admin` or `employee` |
| `vacationBalance` | number | Remaining vacation days |
| `startDate` | string (date) | Employment start date (YYYY-MM-DD) |
| `emailPreferences` | object | Email notification settings |

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `INVALID_CREDENTIALS` | 401 | Email/password mismatch |
| `USER_NOT_FOUND` | 404 | Email not registered |
| `VALIDATION_ERROR` | 400 | Missing/invalid fields |

**Implementation File:** `internal/handler/auth.go`

---

## Authenticated Endpoints

> All endpoints require `Authorization: Bearer <token>` header

### GET /api/auth/me

Get current authenticated user profile.

**Auth Required:** Yes (any role)

**Request:**
```http
GET /api/auth/me HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@company.com",
  "name": "John Doe",
  "role": "employee",
  "vacationBalance": 25,
  "startDate": "2023-01-15",
  "emailPreferences": {
    "vacationUpdates": true,
    "weeklyDigest": false,
    "teamNotifications": true
  }
}
```

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `AUTH_TOKEN_MISSING` | 401 | No Authorization header |
| `AUTH_TOKEN_INVALID` | 401 | Malformed/expired token |

**Implementation File:** `internal/handler/auth.go`

---

### PUT /api/auth/password

Change current user's password.

**Auth Required:** Yes (any role)

**Request:**
```http
PUT /api/auth/password HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
Content-Type: application/json

{
  "currentPassword": "oldPassword123",
  "newPassword": "newSecurePassword456"
}
```

**Request Schema:**

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `currentPassword` | string | Yes | Must match current |
| `newPassword` | string | Yes | Min 6, max 72 chars |

**Response (200 OK):**
```json
{
  "message": "Password updated successfully"
}
```

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `INVALID_CREDENTIALS` | 401 | Current password wrong |
| `VALIDATION_ERROR` | 400 | New password too short/long |

**Implementation File:** `internal/handler/auth.go`

---

### PUT /api/auth/email-preferences

Update email notification preferences.

**Auth Required:** Yes (any role)

**Request:**
```http
PUT /api/auth/email-preferences HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
Content-Type: application/json

{
  "vacationUpdates": true,
  "weeklyDigest": false,
  "teamNotifications": true
}
```

**Request Schema:**

| Field | Type | Required | Default |
|-------|------|----------|---------|
| `vacationUpdates` | boolean | No | true |
| `weeklyDigest` | boolean | No | false |
| `teamNotifications` | boolean | No | true |

**Response (200 OK):**
```json
{
  "emailPreferences": {
    "vacationUpdates": true,
    "weeklyDigest": false,
    "teamNotifications": true
  }
}
```

**Implementation File:** `internal/handler/auth.go`

---

## Employee Endpoints

> Require `Authorization: Bearer <token>` with `employee` or `admin` role

### POST /api/vacation/request

Submit a new vacation request.

**Auth Required:** Yes (Employee or Admin)

**Request:**
```http
POST /api/vacation/request HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
Content-Type: application/json

{
  "startDate": "15/01/2024",
  "endDate": "19/01/2024",
  "reason": "Family vacation"
}
```

**Request Schema:**

| Field | Type | Required | Format | Validation |
|-------|------|----------|--------|------------|
| `startDate` | string | Yes | DD/MM/YYYY | Must be future date |
| `endDate` | string | Yes | DD/MM/YYYY | Must be >= startDate |
| `reason` | string | No | Text | Max 500 chars |

**Response (201 Created):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "startDate": "2024-01-15",
  "endDate": "2024-01-19",
  "totalDays": 5,
  "reason": "Family vacation",
  "status": "pending",
  "createdAt": "2024-01-10T14:30:00Z",
  "updatedAt": "2024-01-10T14:30:00Z"
}
```

**Response Schema:**

| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Request unique identifier |
| `userId` | string (UUID) | Requesting user's ID |
| `startDate` | string (date) | Start date (YYYY-MM-DD) |
| `endDate` | string (date) | End date (YYYY-MM-DD) |
| `totalDays` | number | Business days count |
| `reason` | string | Request reason (nullable) |
| `status` | string | `pending`, `approved`, `rejected` |
| `createdAt` | string (datetime) | ISO 8601 timestamp |
| `updatedAt` | string (datetime) | ISO 8601 timestamp |

**Business Rules:**
- `totalDays` calculated as business days (excludes weekends if `excludeWeekends` setting is true)
- Cannot request dates in the past
- Cannot exceed vacation balance
- Overlapping requests are allowed (admin discretion)

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `INSUFFICIENT_BALANCE` | 422 | Not enough vacation days |
| `INVALID_DATE_RANGE` | 400 | End date before start date |
| `DATE_IN_PAST` | 400 | Start date is in the past |
| `VALIDATION_ERROR` | 400 | Invalid date format |

**Implementation File:** `internal/handler/vacation.go`

---

### GET /api/vacation/requests

Get current user's vacation requests with optional filtering.

**Auth Required:** Yes (Employee or Admin)

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `status` | string | No | Filter by status: `pending`, `approved`, `rejected` |
| `year` | number | No | Filter by year (default: current year) |

**Request:**
```http
GET /api/vacation/requests?status=approved&year=2024 HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "requests": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "userId": "550e8400-e29b-41d4-a716-446655440000",
      "startDate": "2024-01-15",
      "endDate": "2024-01-19",
      "totalDays": 5,
      "reason": "Family vacation",
      "status": "approved",
      "reviewedBy": "770e8400-e29b-41d4-a716-446655440002",
      "reviewedAt": "2024-01-11T09:00:00Z",
      "createdAt": "2024-01-10T14:30:00Z",
      "updatedAt": "2024-01-11T09:00:00Z"
    }
  ],
  "total": 1
}
```

**Implementation File:** `internal/handler/vacation.go`

---

### GET /api/vacation/requests/:id

Get a specific vacation request by ID.

**Auth Required:** Yes (Employee - own requests only, Admin - any request)

**Request:**
```http
GET /api/vacation/requests/660e8400-e29b-41d4-a716-446655440001 HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "userName": "John Doe",
  "startDate": "2024-01-15",
  "endDate": "2024-01-19",
  "totalDays": 5,
  "reason": "Family vacation",
  "status": "pending",
  "createdAt": "2024-01-10T14:30:00Z",
  "updatedAt": "2024-01-10T14:30:00Z"
}
```

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `REQUEST_NOT_FOUND` | 404 | Request doesn't exist |
| `ACCESS_DENIED` | 403 | Employee accessing another's request |

**Implementation File:** `internal/handler/vacation.go`

---

### DELETE /api/vacation/requests/:id

Cancel a pending vacation request.

**Auth Required:** Yes (Employee - own requests only)

**Request:**
```http
DELETE /api/vacation/requests/660e8400-e29b-41d4-a716-446655440001 HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
```

**Response (204 No Content):**
```
(empty body)
```

**Business Rules:**
- Only `pending` requests can be cancelled
- Only the request owner can cancel (employees)
- Admins use the admin endpoint instead

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `REQUEST_NOT_FOUND` | 404 | Request doesn't exist |
| `CANNOT_CANCEL_APPROVED` | 422 | Request already approved |
| `CANNOT_CANCEL_REJECTED` | 422 | Request already rejected |
| `ACCESS_DENIED` | 403 | Not the request owner |

**Implementation File:** `internal/handler/vacation.go`

---

### GET /api/vacation/team

Get team vacation calendar (approved requests visible to all employees).

**Auth Required:** Yes (Employee or Admin)

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `month` | number | No | Month (1-12), default: current |
| `year` | number | No | Year, default: current |

**Request:**
```http
GET /api/vacation/team?month=1&year=2024 HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "vacations": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "userId": "550e8400-e29b-41d4-a716-446655440000",
      "userName": "John Doe",
      "startDate": "2024-01-15",
      "endDate": "2024-01-19",
      "totalDays": 5
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440002",
      "userId": "550e8400-e29b-41d4-a716-446655440003",
      "userName": "Jane Smith",
      "startDate": "2024-01-22",
      "endDate": "2024-01-26",
      "totalDays": 5
    }
  ],
  "month": 1,
  "year": 2024
}
```

**Note:** Only `approved` vacations are returned. Sensitive fields (`reason`, reviewer info) are excluded.

**Implementation File:** `internal/handler/vacation.go`

---

## Admin Endpoints

> Require `Authorization: Bearer <token>` with `admin` role

### GET /api/admin/users

List all users with pagination and filtering.

**Auth Required:** Yes (Admin only)

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | number | No | Page number (default: 1) |
| `limit` | number | No | Items per page (default: 20, max: 100) |
| `role` | string | No | Filter by role: `admin`, `employee` |
| `search` | string | No | Search by name or email |

**Request:**
```http
GET /api/admin/users?page=1&limit=20&role=employee HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john@company.com",
      "name": "John Doe",
      "role": "employee",
      "vacationBalance": 25,
      "startDate": "2023-01-15",
      "createdAt": "2023-01-15T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "totalPages": 3
  }
}
```

**Implementation File:** `internal/handler/admin.go`

---

### POST /api/admin/users

Create a new user (employee or admin).

**Auth Required:** Yes (Admin only)

**Request:**
```http
POST /api/admin/users HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
Content-Type: application/json

{
  "email": "newuser@company.com",
  "password": "temporaryPassword123",
  "name": "New User",
  "role": "employee",
  "vacationBalance": 25,
  "startDate": "15/01/2024"
}
```

**Request Schema:**

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `email` | string | Yes | Valid email, unique |
| `password` | string | Yes | Min 6, max 72 chars |
| `name` | string | Yes | Min 1, max 100 chars |
| `role` | string | Yes | `admin` or `employee` |
| `vacationBalance` | number | No | Default: 25, min: 0 |
| `startDate` | string | No | DD/MM/YYYY format |

**Response (201 Created):**
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440004",
  "email": "newuser@company.com",
  "name": "New User",
  "role": "employee",
  "vacationBalance": 25,
  "startDate": "2024-01-15",
  "createdAt": "2024-01-10T15:00:00Z"
}
```

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `EMAIL_ALREADY_EXISTS` | 409 | Email already registered |
| `VALIDATION_ERROR` | 400 | Invalid field values |

**Implementation File:** `internal/handler/admin.go`

---

### GET /api/admin/users/:id

Get a specific user's details.

**Auth Required:** Yes (Admin only)

**Request:**
```http
GET /api/admin/users/550e8400-e29b-41d4-a716-446655440000 HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@company.com",
  "name": "John Doe",
  "role": "employee",
  "vacationBalance": 25,
  "startDate": "2023-01-15",
  "emailPreferences": {
    "vacationUpdates": true,
    "weeklyDigest": false,
    "teamNotifications": true
  },
  "createdAt": "2023-01-15T10:00:00Z",
  "updatedAt": "2024-01-10T14:00:00Z"
}
```

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `USER_NOT_FOUND` | 404 | User doesn't exist |

**Implementation File:** `internal/handler/admin.go`

---

### PUT /api/admin/users/:id

Update a user's information.

**Auth Required:** Yes (Admin only)

**Request:**
```http
PUT /api/admin/users/550e8400-e29b-41d4-a716-446655440000 HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "John Updated",
  "role": "admin",
  "vacationBalance": 30
}
```

**Request Schema (all fields optional):**

| Field | Type | Validation |
|-------|------|------------|
| `email` | string | Valid email, unique |
| `name` | string | Min 1, max 100 chars |
| `role` | string | `admin` or `employee` |
| `vacationBalance` | number | Min 0 |
| `startDate` | string | DD/MM/YYYY format |

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@company.com",
  "name": "John Updated",
  "role": "admin",
  "vacationBalance": 30,
  "startDate": "2023-01-15",
  "updatedAt": "2024-01-10T16:00:00Z"
}
```

**Business Rules:**
- Cannot demote the last admin
- Cannot modify own role (prevents lockout)

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `USER_NOT_FOUND` | 404 | User doesn't exist |
| `EMAIL_ALREADY_EXISTS` | 409 | New email already taken |
| `CANNOT_REMOVE_LAST_ADMIN` | 422 | Would leave no admins |
| `CANNOT_MODIFY_OWN_ROLE` | 422 | Self-demotion attempt |

**Implementation File:** `internal/handler/admin.go`

---

### DELETE /api/admin/users/:id

Delete a user account.

**Auth Required:** Yes (Admin only)

**Request:**
```http
DELETE /api/admin/users/550e8400-e29b-41d4-a716-446655440000 HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
```

**Response (204 No Content):**
```
(empty body)
```

**Business Rules:**
- Cannot delete yourself
- Cannot delete the last admin
- Deleting a user deletes their vacation requests (CASCADE)

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `USER_NOT_FOUND` | 404 | User doesn't exist |
| `CANNOT_DELETE_SELF` | 422 | Self-deletion attempt |
| `CANNOT_DELETE_LAST_ADMIN` | 422 | Would leave no admins |

**Implementation File:** `internal/handler/admin.go`

---

### GET /api/admin/vacation/pending

Get all pending vacation requests for review.

**Auth Required:** Yes (Admin only)

**Request:**
```http
GET /api/admin/vacation/pending HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "requests": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "userId": "550e8400-e29b-41d4-a716-446655440000",
      "userName": "John Doe",
      "userEmail": "john@company.com",
      "startDate": "2024-01-15",
      "endDate": "2024-01-19",
      "totalDays": 5,
      "reason": "Family vacation",
      "status": "pending",
      "createdAt": "2024-01-10T14:30:00Z"
    }
  ],
  "total": 1
}
```

**Implementation File:** `internal/handler/admin.go`

---

### PUT /api/admin/vacation/:id/approve

Approve a pending vacation request.

**Auth Required:** Yes (Admin only)

**Request:**
```http
PUT /api/admin/vacation/660e8400-e29b-41d4-a716-446655440001/approve HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "status": "approved",
  "reviewedBy": "770e8400-e29b-41d4-a716-446655440002",
  "reviewedAt": "2024-01-11T09:00:00Z"
}
```

**Side Effects:**
- User's `vacationBalance` is decremented by `totalDays`
- Email notification sent to user (if enabled)

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `REQUEST_NOT_FOUND` | 404 | Request doesn't exist |
| `REQUEST_ALREADY_PROCESSED` | 422 | Not in pending status |

**Implementation File:** `internal/handler/admin.go`

---

### PUT /api/admin/vacation/:id/reject

Reject a pending vacation request.

**Auth Required:** Yes (Admin only)

**Request:**
```http
PUT /api/admin/vacation/660e8400-e29b-41d4-a716-446655440001/reject HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
Content-Type: application/json

{
  "reason": "Insufficient coverage during that period"
}
```

**Request Schema:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `reason` | string | No | Rejection reason (shown to user) |

**Response (200 OK):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "status": "rejected",
  "rejectionReason": "Insufficient coverage during that period",
  "reviewedBy": "770e8400-e29b-41d4-a716-446655440002",
  "reviewedAt": "2024-01-11T09:00:00Z"
}
```

**Side Effects:**
- Email notification sent to user (if enabled)

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `REQUEST_NOT_FOUND` | 404 | Request doesn't exist |
| `REQUEST_ALREADY_PROCESSED` | 422 | Not in pending status |

**Implementation File:** `internal/handler/admin.go`

---

### GET /api/admin/settings

Get application settings.

**Auth Required:** Yes (Admin only)

**Request:**
```http
GET /api/admin/settings HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": "settings",
  "weekendPolicy": {
    "excludeWeekends": true,
    "excludedDays": [0, 6]
  },
  "newsletter": {
    "enabled": true,
    "frequency": "monthly",
    "dayOfMonth": 1
  },
  "defaultVacationDays": 25,
  "vacationResetMonth": 1,
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Implementation File:** `internal/handler/admin.go`

---

### PUT /api/admin/settings

Update application settings.

**Auth Required:** Yes (Admin only)

**Request:**
```http
PUT /api/admin/settings HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
Content-Type: application/json

{
  "weekendPolicy": {
    "excludeWeekends": true,
    "excludedDays": [0, 6]
  },
  "newsletter": {
    "enabled": true,
    "frequency": "monthly",
    "dayOfMonth": 1
  },
  "defaultVacationDays": 25,
  "vacationResetMonth": 1
}
```

**Request Schema:**

| Field | Type | Description |
|-------|------|-------------|
| `weekendPolicy.excludeWeekends` | boolean | Exclude weekends from day count |
| `weekendPolicy.excludedDays` | number[] | Days to exclude (0=Sun, 6=Sat) |
| `newsletter.enabled` | boolean | Enable newsletter feature |
| `newsletter.frequency` | string | `weekly` or `monthly` |
| `newsletter.dayOfMonth` | number | Day to send (1-28) |
| `defaultVacationDays` | number | Default balance for new users |
| `vacationResetMonth` | number | Month to reset balances (1-12) |

**Response (200 OK):**
```json
{
  "id": "settings",
  "weekendPolicy": {
    "excludeWeekends": true,
    "excludedDays": [0, 6]
  },
  "newsletter": {
    "enabled": true,
    "frequency": "monthly",
    "dayOfMonth": 1
  },
  "defaultVacationDays": 25,
  "vacationResetMonth": 1,
  "updatedAt": "2024-01-10T16:30:00Z"
}
```

**Implementation File:** `internal/handler/admin.go`

---

### POST /api/admin/newsletter/send

Manually trigger newsletter send.

**Auth Required:** Yes (Admin only)

**Request:**
```http
POST /api/admin/newsletter/send HTTP/1.1
Host: localhost:3000
Authorization: Bearer <token>
Content-Type: application/json

{
  "preview": false
}
```

**Request Schema:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `preview` | boolean | No | If true, sends only to admins |

**Response (200 OK):**
```json
{
  "message": "Newsletter sent successfully",
  "recipientCount": 45,
  "sentAt": "2024-01-10T16:45:00Z"
}
```

**Errors:**

| Code | Status | Condition |
|------|--------|-----------|
| `EMAIL_SERVICE_UNAVAILABLE` | 503 | Resend API error |
| `NEWSLETTER_DISABLED` | 422 | Newsletter feature disabled |

**Implementation File:** `internal/handler/admin.go`

---

## Error Code Reference

### Authentication Errors (4xx)

| Code | HTTP Status | Message |
|------|-------------|---------|
| `AUTH_TOKEN_MISSING` | 401 | Authorization header required |
| `AUTH_TOKEN_INVALID` | 401 | Invalid or expired token |
| `AUTH_TOKEN_EXPIRED` | 401 | Token has expired |
| `INVALID_CREDENTIALS` | 401 | Invalid email or password |
| `ACCESS_DENIED` | 403 | Insufficient permissions |
| `ADMIN_REQUIRED` | 403 | Admin role required |

### Validation Errors (400)

| Code | HTTP Status | Message |
|------|-------------|---------|
| `VALIDATION_ERROR` | 400 | Request validation failed |
| `INVALID_DATE_FORMAT` | 400 | Date must be DD/MM/YYYY |
| `INVALID_DATE_RANGE` | 400 | End date must be after start date |
| `DATE_IN_PAST` | 400 | Date cannot be in the past |
| `INVALID_EMAIL_FORMAT` | 400 | Invalid email address |
| `INVALID_UUID` | 400 | Invalid UUID format |

### Resource Errors (404)

| Code | HTTP Status | Message |
|------|-------------|---------|
| `USER_NOT_FOUND` | 404 | User not found |
| `REQUEST_NOT_FOUND` | 404 | Vacation request not found |
| `SETTINGS_NOT_FOUND` | 404 | Settings not found |

### Conflict Errors (409)

| Code | HTTP Status | Message |
|------|-------------|---------|
| `EMAIL_ALREADY_EXISTS` | 409 | Email address already registered |
| `REQUEST_ALREADY_EXISTS` | 409 | Duplicate vacation request |

### Business Rule Errors (422)

| Code | HTTP Status | Message |
|------|-------------|---------|
| `INSUFFICIENT_BALANCE` | 422 | Insufficient vacation balance |
| `CANNOT_CANCEL_APPROVED` | 422 | Cannot cancel approved request |
| `CANNOT_CANCEL_REJECTED` | 422 | Cannot cancel rejected request |
| `REQUEST_ALREADY_PROCESSED` | 422 | Request already approved/rejected |
| `CANNOT_DELETE_SELF` | 422 | Cannot delete own account |
| `CANNOT_DELETE_LAST_ADMIN` | 422 | Cannot delete last admin |
| `CANNOT_REMOVE_LAST_ADMIN` | 422 | Cannot demote last admin |
| `CANNOT_MODIFY_OWN_ROLE` | 422 | Cannot modify own role |
| `NEWSLETTER_DISABLED` | 422 | Newsletter feature is disabled |

### Server Errors (5xx)

| Code | HTTP Status | Message |
|------|-------------|---------|
| `INTERNAL_ERROR` | 500 | Internal server error |
| `DATABASE_ERROR` | 500 | Database operation failed |
| `EMAIL_SERVICE_UNAVAILABLE` | 503 | Email service unavailable |

---

## Go Implementation Reference

### Handler Signature Pattern

```go
// internal/handler/vacation.go
func (h *VacationHandler) CreateRequest(c *gin.Context) {
    // 1. Parse and validate request body
    var req dto.CreateVacationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, dto.ErrorResponse{
            Code:    "VALIDATION_ERROR",
            Message: err.Error(),
        })
        return
    }

    // 2. Get user from context (set by auth middleware)
    userID := c.GetString("userID")

    // 3. Call service layer
    vacation, err := h.vacationService.Create(c.Request.Context(), userID, req)
    if err != nil {
        // Map domain errors to HTTP responses
        handleError(c, err)
        return
    }

    // 4. Return success response
    c.JSON(201, vacation)
}
```

### Error Handling Pattern

```go
// internal/dto/errors.go
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Status  int    `json:"-"`
}

var (
    ErrInsufficientBalance = &AppError{
        Code:    "INSUFFICIENT_BALANCE",
        Message: "Insufficient vacation balance",
        Status:  422,
    }
    ErrRequestNotFound = &AppError{
        Code:    "REQUEST_NOT_FOUND",
        Message: "Vacation request not found",
        Status:  404,
    }
)
```

### Route Registration

```go
// cmd/server/main.go
func setupRoutes(r *gin.Engine, h *Handlers, m *Middleware) {
    // Public routes
    r.GET("/health", h.Health.Check)

    api := r.Group("/api")
    {
        // Auth routes (public)
        auth := api.Group("/auth")
        {
            auth.POST("/login", h.Auth.Login)
        }

        // Authenticated routes
        authenticated := api.Group("/")
        authenticated.Use(m.Auth.RequireAuth())
        {
            authenticated.GET("/auth/me", h.Auth.Me)
            authenticated.PUT("/auth/password", h.Auth.ChangePassword)
            authenticated.PUT("/auth/email-preferences", h.Auth.UpdateEmailPrefs)

            // Employee routes
            vacation := authenticated.Group("/vacation")
            {
                vacation.POST("/request", h.Vacation.Create)
                vacation.GET("/requests", h.Vacation.List)
                vacation.GET("/requests/:id", h.Vacation.Get)
                vacation.DELETE("/requests/:id", h.Vacation.Cancel)
                vacation.GET("/team", h.Vacation.Team)
            }
        }

        // Admin routes
        admin := api.Group("/admin")
        admin.Use(m.Auth.RequireAuth(), m.Auth.RequireAdmin())
        {
            admin.GET("/users", h.Admin.ListUsers)
            admin.POST("/users", h.Admin.CreateUser)
            admin.GET("/users/:id", h.Admin.GetUser)
            admin.PUT("/users/:id", h.Admin.UpdateUser)
            admin.DELETE("/users/:id", h.Admin.DeleteUser)

            admin.GET("/vacation/pending", h.Admin.PendingRequests)
            admin.PUT("/vacation/:id/approve", h.Admin.ApproveRequest)
            admin.PUT("/vacation/:id/reject", h.Admin.RejectRequest)

            admin.GET("/settings", h.Admin.GetSettings)
            admin.PUT("/settings", h.Admin.UpdateSettings)
            admin.POST("/newsletter/send", h.Admin.SendNewsletter)
        }
    }
}
```

---

## Frontend API Client Reference

### TypeScript API Module

```typescript
// src/lib/api/client.ts
const API_BASE = '/api';

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getAuthToken();

  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json();
    throw new ApiError(error.code, error.message, response.status);
  }

  return response.json();
}

// Vacation API
export const vacationApi = {
  create: (data: CreateVacationRequest) =>
    request<VacationRequest>('/vacation/request', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  list: (params?: { status?: string; year?: number }) =>
    request<VacationListResponse>(`/vacation/requests?${new URLSearchParams(params)}`),

  cancel: (id: string) =>
    request<void>(`/vacation/requests/${id}`, { method: 'DELETE' }),

  team: (month: number, year: number) =>
    request<TeamVacationResponse>(`/vacation/team?month=${month}&year=${year}`),
};
```

---

## Related Documents

- [00-architecture-overview.md](./00-architecture-overview.md) - System architecture
- [01-database-schema.md](./01-database-schema.md) - Database design
- [04-backend-tasks.md](./04-backend-tasks.md) - Handler implementation tasks
