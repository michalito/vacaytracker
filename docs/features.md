# VacayTracker - Complete Features Documentation

> A production-ready employee vacation tracking application with role-based access control, email notifications, and team calendar visualization.

---

## Table of Contents

1. [Authentication & Authorization](#1-authentication--authorization)
2. [Employee Features](#2-employee-features)
3. [Admin Features](#3-admin-features)
4. [Calendar System](#4-calendar-system)
5. [Email Notification System](#5-email-notification-system)
6. [Newsletter System](#6-newsletter-system)
7. [User Management](#7-user-management)
8. [System Settings](#8-system-settings)
9. [UI/UX Features](#9-uiux-features)
10. [Security Features](#10-security-features)
11. [Infrastructure Features](#11-infrastructure-features)

---

## 1. Authentication & Authorization

### 1.1 Login System
- **Tab-based login interface** with separate Employee and Admin tabs
- **Form validation** requiring both username and password fields
- **Role verification** ensuring selected tab matches user's actual role
- **Loading states** during authentication with visual feedback
- **Error messaging** via toast notifications for failed login attempts
- **Auto-redirect** for already authenticated users to appropriate dashboard
- **Session persistence** using sessionStorage with JWT tokens

### 1.2 JWT Token Management
- **24-hour token expiration** for security
- **Automatic token injection** in all API requests via storage.js module
- **Token refresh handling** with automatic redirect to login on expiration
- **Bearer token format** in Authorization header

### 1.3 Role-Based Access Control (RBAC)
- **Two user roles**: Admin (Captain) and Employee (Crew Member)
- **Protected routes**:
  - Admin-only: `/api/admin/*` endpoints
  - Employee-only: `/vacation` (UI), `/api/vacation` (API)
  - Authenticated: All dashboard and API endpoints
  - Public: Login page, health check endpoint
- **Middleware enforcement** at API level with `AuthMiddleware` and `AdminMiddleware`/`EmployeeMiddleware`

---

## 2. Employee Features

### 2.1 Vacation Balance Dashboard
- **Circular SVG progress ring** showing percentage of vacation days used
- **Numeric display** of days remaining in center of ring
- **Usage statistics**: Total days available, days already used
- **Dynamic mood messaging** based on remaining balance:
  - 80%+ remaining: "You're fully charged! Plenty of vacation days"
  - 60-80%: "Looking good! Healthy balance"
  - 40-60%: "Time to start planning!"
  - 20-40%: "Running low! Make those days count"
  - <20%: "Almost out! Dream about next year"
- **Mood emoji indicator** that changes with status

### 2.2 Next Vacation Countdown
- **Countdown display** for next approved vacation (e.g., "5 days", "Tomorrow", "Today")
- **Vacation dates** shown in card format
- **Duration** in business days
- **Conditional rendering** - only shows when approved vacation exists

### 2.3 Vacation Request Workflow
- **Request Time Off button** prominently displayed on dashboard
- **Date selection modal** with:
  - EU date format input fields (DD/MM/YYYY) with auto-formatting
  - Full calendar datepicker for visual date selection
  - Date range selection (click start date, then end date)
  - Past dates disabled
  - Weekends visually distinguished
  - Today highlighted
  - Selected range shading
- **Duration preview** showing calculated business days in real-time
- **Weekend settings awareness** - calculation respects system weekend policy
- **Vacation reason field**:
  - Optional textarea (max 200 characters)
  - Quick suggestion chips: "Family Time", "Travel", "Staycation"
  - Character counter with warning at 180+ characters
- **Availability validation** - checks sufficient days before submission
- **Days preview** showing impact on vacation balance before submitting

### 2.4 Vacation History
- **Complete history** of all vacation requests (approved, pending, rejected)
- **Chronological sorting** (newest first)
- **Per-entry details**:
  - Date range
  - Status with color coding (green/approved, orange/pending, red/rejected)
  - Duration in business days
  - Request submission date
  - Vacation reason
- **Empty state messaging** when no history exists
- **Animated entry** with staggered delays for visual appeal

### 2.5 Team Vacation Visibility
- **Upcoming team vacations timeline** showing next 10 approved vacations
- **Per-vacation details**: Employee name, dates, duration, relative time, reason
- **Own vacation highlighting** for easy identification
- **"Starting soon" indicators** for vacations within 7 days
- **Team member count** showing how many people have upcoming vacation

### 2.6 Employee Settings
- **Password change functionality**:
  - Current password verification
  - New password with confirmation
  - Minimum 6 character requirement
  - Password match validation
- **Success celebration animation** with fireworks emoji cascade
- **Time-based greeting** in header

---

## 3. Admin Features

### 3.1 Admin Dashboard Overview
- **Animated greeting** based on time of day (Good morning/afternoon/evening)
- **Vacation-themed animated background** with floating icons (anchor, compass, ship, island)
- **Statistics cards**:
  - Total Crew: Employee count
  - On Vacation: Employees currently on approved vacation
  - Pending Requests: Number awaiting approval
  - Avg Days Remaining: Team average vacation balance

### 3.2 Vacation Request Management
- **Pending requests section** (conditionally displayed when requests exist)
- **Request cards** showing:
  - Employee name
  - Submission date
  - Vacation dates with business days calculated
  - Reason (if provided)
- **One-click approve/reject buttons** with visual feedback
- **"All caught up!" message** when no pending requests
- **Automatic notifications** sent to employees on status change

### 3.3 Crew Vacation Timeline
- **Chronological list** of upcoming team vacations (next 10)
- **Per-vacation display**:
  - Employee name
  - Status badge (Approved/Pending/Rejected)
  - Date range in compact format
  - Duration in days
  - Relative time indicator (Today, Tomorrow, X days/weeks/months)
  - Vacation reason
- **"Starting Soon" highlight** for vacations within 7 days
- **Empty state** when no upcoming vacations

### 3.4 Crew Vacation Balance Management
- **Employee cards grid** showing all employees
- **Per-employee card**:
  - Avatar with first initial
  - Employee name
  - Visual progress bar of vacation usage
  - Remaining days highlighted
  - "X of Y used" text
  - Vacation mood emoji
  - "Adjust" button for manual adjustments
- **Search functionality** to filter employees by name
- **Adjust Vacation Days modal**:
  - Employee name (read-only)
  - Vacation days input field
  - Save Changes button

### 3.5 Year-End Vacation Reset
- **"New Vacation Year" button** for bulk operations
- **Reset options**:
  - Set new vacation days total for all employees
  - Clear vacation history toggle
- **Impact preview** showing:
  - Number of affected employees
  - Current average days
- **Two-step confirmation** with detailed warning modal
- **"This action cannot be undone" warning**
- **Automatic email notifications** to all employees after reset

### 3.6 Admin Settings
- **Password change** (same as employee with current password verification)
- **Weekend Policy configuration**:
  - Toggle between "Exclude Weekends" and "All Days Count"
  - Visual impact display showing calculation difference
  - Save button for applying changes
- **Visual toggle switch** with office/beach icons

---

## 4. Calendar System

### 4.1 Calendar Views
- **Week View**: 7-day horizontal display
- **Month View**: Full month grid
- **View toggle dropdown** for switching between views

### 4.2 Calendar Navigation
- **Previous/Next buttons** for period navigation
- **"Today" button** for quick return to current date
- **Period indicator** showing current date range or month name

### 4.3 Calendar Display Features
- **Day headers** with weekend highlighting
- **Today highlight** for current date
- **Vacation events** color-coded by status:
  - Green: Approved
  - Orange: Pending
  - Red: Rejected
- **Employee identification** via initials/names on vacation blocks
- **Tooltip information** on hover showing:
  - Employee name
  - Date range
  - Duration
  - Status

### 4.4 Admin Calendar Filters
- **Employee filter dropdown** populated with all employees
- **Status filter** (Approved/Pending/Rejected)
- **Combined filtering** capability

### 4.5 Calendar Performance
- **5-minute caching system** for vacation data
- **Efficient data fetching** from appropriate endpoints
- **Loading states** during data retrieval

---

## 5. Email Notification System

### 5.1 Email Service Configuration
- **Resend API integration** for reliable delivery
- **Configurable sender** (name and email address via environment variables)
- **Retry logic**: 3 attempts with exponential backoff
- **Graceful degradation** when email service unavailable

### 5.2 Welcome Email
- **Trigger**: New user creation by admin
- **Recipient**: New user's email address
- **Content**:
  - Login credentials (username and temporary password)
  - Security warning to change password
  - Login link to application
- **Conditional**: Only sent if email address provided and preferences enabled

### 5.3 Vacation Request Notification
- **Trigger**: Employee submits vacation request
- **Recipients**: All admins with vacation request notifications enabled
- **Content**:
  - Employee name
  - Vacation dates
  - Duration in business days
  - Optional reason/notes
  - Review button linking to admin dashboard

### 5.4 Vacation Status Update
- **Trigger**: Admin approves or rejects vacation request
- **Recipient**: Requesting employee
- **Content**:
  - Status (Approved/Rejected) with color coding
  - Vacation dates and duration
  - Admin name who reviewed
  - Employee's original notes
  - Link to employee dashboard
- **Color coding**: Green (#28a745) for approved, Red (#dc3545) for rejected

### 5.5 Vacation Days Reset Notification
- **Trigger**: Admin performs year-end vacation reset
- **Recipients**: All employees with email addresses
- **Content**:
  - New vacation days allocation
  - Whether history was cleared
  - Effective date
  - Link to employee dashboard

---

## 6. Newsletter System

### 6.1 Monthly Vacation Summary Email
- **"Beach Report" themed** email with tropical styling
- **Content**:
  - Vacations taken in current month with employee names and durations
  - Next upcoming vacation (single entry)
  - Beach-themed statistics and emoji
  - Link to full schedule
- **Automated scheduling** (Go internal scheduler)
- **Manual send option** for admins

### 6.2 Newsletter Settings (Admin)
- **Enable/disable toggle** for automatic newsletters
- **Schedule configuration**:
  - Day of month (1-31)
  - Time of day (12 AM - 11 PM)
- **Next scheduled send preview**
- **Last sent date/time display**
- **Recipient count display**

### 6.3 Newsletter Preview
- **Preview modal** showing:
  - Email subject line
  - Recipient count
  - Full HTML email content (sanitized)
- **Send Now button** with confirmation dialog

### 6.4 User Subscription Management
- **Bulk actions**:
  - Select All / Deselect All
  - Enable Selected / Disable Selected
- **User subscription table** with:
  - Checkbox selection
  - Name, email, role columns
  - Subscription status (Subscribed/Unsubscribed)
  - Individual toggle button
- **Search functionality** for filtering users
- **Statistics display**: Subscribed/Unsubscribed/Total counts

---

## 7. User Management

### 7.1 User Registry (Admin)
- **Crew stats header** showing:
  - Total Crew count
  - Captains (admins) count
  - Crew Members (employees) count

### 7.2 User Cards Display
- **Per-user card**:
  - Avatar with first initial
  - User name
  - Role badge (Captain/Crew)
  - Username
  - Email address
  - Vacation days (employees only)
  - Edit and Delete buttons
- **Current user identification** with "(Current User)" label
- **Search functionality** with debounced input

### 7.3 Add User
- **Form fields**:
  - Full Name (required)
  - Username (required, must be unique)
  - Email Address (optional, for notifications)
  - Password (required)
  - Role dropdown (Crew Member/Captain)
  - Vacation Days (employees only)
- **Validation**: Name and username required
- **Welcome email** sent automatically if email provided

### 7.4 Edit User
- **Same form as Add** with pre-populated values
- **Password field hint**: "Leave blank to keep current"
- **Role change** triggers email preference updates
- **Self-edit restrictions**: Cannot change own role

### 7.5 Delete User
- **Confirmation required** before deletion
- **Cascade delete**: Removes all associated vacation requests
- **Self-delete prevention**: Cannot delete own account

---

## 8. System Settings

### 8.1 Weekend Policy
- **Toggle option**: Count weekends as vacation days or not
- **Default**: Weekends excluded (only business days count)
- **Impact**: Affects all vacation day calculations throughout system
- **Visual explanation** of impact on vacation duration

### 8.2 Email Preferences (Per User)
- **Master toggle**: Enable/disable all email notifications
- **Admin-specific preferences**:
  - Vacation request notifications (when employees submit)
  - User created notifications
- **Employee-specific preferences**:
  - Status update notifications (when requests approved/rejected)
- **Shared preferences**:
  - Vacation reset notifications
  - Monthly vacation summary (newsletter)

---

## 9. UI/UX Features

### 9.1 Toast Notification System
- **Types**: Success, Error, Warning, Info
- **Features**:
  - Auto-dismiss with configurable duration
  - Progress bar animation
  - Pause on hover
  - Stacking animation
  - Action buttons support
  - Randomized personality emoji
- **Accessibility**: ARIA attributes for screen readers

### 9.2 Modal System
- **Focus trapping** within modal (Tab key navigation)
- **Click-outside-to-close** behavior
- **Escape key** closes modal
- **Focus restoration** after modal closes
- **Implementation** using Melt UI Dialog builder

### 9.3 Visual Design Theme
- **Vacation/beach theme** throughout:
  - Beach emoji and metaphors
  - Tropical color palette
  - Vacation-related icons
- **Animated backgrounds** with sand particle effects
- **Floating animated icons** (anchor, compass, ship, palm tree)

### 9.4 Animations Library
- **30+ custom animations**:
  - spin, shake, bounce, float, gentle-float, sway
  - shimmer, pulse, pulse-glow, fadeIn, fadeInUp
  - slideInRight, slideInLeft, popIn, fly-around, wave
- **Accessibility support**: Respects prefers-reduced-motion setting

### 9.5 Responsive Design
- **Mobile-first approach**
- **Breakpoints**: 1024px, 768px, 480px
- **Adaptive layouts** for cards and grids
- **Touch-friendly interactions**

### 9.6 Accessibility Features
- **Semantic HTML** structure
- **ARIA labels** and descriptions
- **Visually hidden labels** for screen readers
- **Keyboard navigation** support
- **Focus management** in interactive components
- **Color contrast** compliance
- **Form validation feedback**

---

## 10. Security Features

### 10.1 Password Security
- **bcrypt hashing** with 10 salt rounds
- **No plaintext storage**
- **Password migration utility** for legacy data

### 10.2 Authentication Security
- **JWT tokens** with 24-hour expiration
- **HS256 algorithm** for token signing
- **Environment-based secrets** (not hardcoded)
- **Token validation** on every protected request

### 10.3 API Security
- **CORS configuration** with restricted origins
- **Role-based endpoint protection**
- **Input validation** on all endpoints
- **SQL injection prevention** (N/A - JSON database)

### 10.4 XSS Prevention
- **HTML sanitization** with allowed tag/attribute whitelist
- **DOMParser-based sanitization** for email previews
- **URL validation** preventing javascript: and data: protocols
- **HTML escaping** for user-generated content

### 10.5 Container Security
- **Non-root user** in production containers
- **Read-only source mounts** in development
- **Alpine base image** (minimal attack surface)
- **Environment separation** for dev/prod

---

## 11. Infrastructure Features

### 11.1 Frontend Architecture
- **Svelte 5** with Runes for reactivity
- **SvelteKit** for routing and SSR
- **Melt UI** for headless accessible components
- **Tailwind CSS v4** for utility-first styling
- **Lucide Svelte** for icons

### 11.2 Database
- **SQLite Database** (`modernc.org/sqlite`)
- **Embedded Database** file at `data/vacaytracker.db`
- **WAL Mode** enabled for performance and concurrency
- **Schema Migrations** via SQL files
- **Data persistence** via Docker volumes

### 11.3 Development Environment
- **Vite** for frontend tooling
- **Hot reload** via Air (Go) and Vite (SvelteKit)
- **Docker Compose** for containerized backend
- **Volume mounts** for live code changes
- **Port 5173** for frontend, **Port 3000** for API

### 11.4 Production Environment
- **Docker Compose** deployment
- **Health checks** at `/health` endpoint
- **Multi-stage Docker build** for minimal image size (Alpine based)
- **Restart policy**: unless-stopped
- **Port 81** for production access

### 11.5 Code Quality
- **ESLint/Prettier** for Frontend
- **golangci-lint** for Backend static analysis
- **Go Test** for backend unit/integration testing
- **Husky + lint-staged** for pre-commit hooks
- **GitHub Actions** for CI/CD pipeline

### 11.6 Cron Jobs
- **Monthly newsletter** automated scheduling
- **Configurable** day and hour
- **Environment-based** schedule configuration

# VacayTracker - Additional Documentation Sections

---

## 12. Data Models

### 12.1 User Schema

```json
{
  "id": "usr_a1b2c3d4",
  "name": "John Doe",
  "username": "john",
  "email": "john@example.com",
  "password": "$2b$10$hashedPasswordString...",
  "role": "employee",
  "vacationDays": 25,
  "usedVacationDays": 5,
  "emailPreferences": {
    "enabled": true,
    "vacationStatusUpdates": true,
    "vacationResetNotifications": true,
    "monthlyVacationSummary": true
  },
  "createdAt": "2024-01-15T09:30:00.000Z",
  "updatedAt": "2024-06-20T14:22:00.000Z"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | string | Yes | Unique identifier (auto-generated) |
| name | string | Yes | Full display name |
| username | string | Yes | Login username (unique) |
| email | string | No | Email for notifications |
| password | string | Yes | bcrypt hashed password |
| role | enum | Yes | `"admin"` or `"employee"` |
| vacationDays | integer | Yes* | Total annual vacation days (*employees only) |
| usedVacationDays | integer | Yes* | Days already used (*employees only) |
| emailPreferences | object | Yes | Notification settings (see below) |
| createdAt | ISO 8601 | Yes | Account creation timestamp |
| updatedAt | ISO 8601 | Yes | Last modification timestamp |

**Email Preferences Object (Employee)**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| enabled | boolean | true | Master toggle for all notifications |
| vacationStatusUpdates | boolean | true | Notify on request approval/rejection |
| vacationResetNotifications | boolean | true | Notify on year-end reset |
| monthlyVacationSummary | boolean | true | Receive monthly newsletter |

**Email Preferences Object (Admin)**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| enabled | boolean | true | Master toggle for all notifications |
| vacationRequestNotifications | boolean | true | Notify when employees submit requests |
| userCreatedNotifications | boolean | true | Confirm new user creation |
| vacationResetNotifications | boolean | true | Confirm year-end reset |
| monthlyVacationSummary | boolean | true | Receive monthly newsletter |

---

### 12.2 Vacation Request Schema

```json
{
  "id": "vac_x9y8z7w6",
  "userId": "usr_a1b2c3d4",
  "startDate": "2024-07-15",
  "endDate": "2024-07-19",
  "businessDays": 5,
  "status": "approved",
  "reason": "Family vacation to the coast",
  "createdAt": "2024-06-01T10:15:00.000Z",
  "reviewedBy": "usr_admin001",
  "reviewedAt": "2024-06-02T08:30:00.000Z"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | string | Yes | Unique identifier (auto-generated) |
| userId | string | Yes | Reference to requesting user |
| startDate | date (YYYY-MM-DD) | Yes | First day of vacation |
| endDate | date (YYYY-MM-DD) | Yes | Last day of vacation |
| businessDays | integer | Yes | Calculated working days (respects weekend policy) |
| status | enum | Yes | `"pending"`, `"approved"`, or `"rejected"` |
| reason | string | No | Optional notes (max 200 chars) |
| createdAt | ISO 8601 | Yes | Request submission timestamp |
| reviewedBy | string | No | Admin user ID who reviewed (null if pending) |
| reviewedAt | ISO 8601 | No | Review timestamp (null if pending) |

---

### 12.3 Settings Schema

```json
{
  "weekendPolicy": {
    "excludeWeekends": true
  },
  "newsletter": {
    "enabled": true,
    "dayOfMonth": 1,
    "hourOfDay": 9,
    "lastSentAt": "2024-06-01T09:00:00.000Z"
  }
}
```

**Weekend Policy**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| excludeWeekends | boolean | true | If true, weekends don't count as vacation days |

**Newsletter Settings**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| enabled | boolean | false | Enable automatic monthly newsletter |
| dayOfMonth | integer | 1 | Day to send (1-31) |
| hourOfDay | integer | 9 | Hour to send in 24h format (0-23) |
| lastSentAt | ISO 8601 | null | Timestamp of last newsletter sent |

---

### 12.4 Database Schema (SQLite)

The SQLite database (`data/vacaytracker.db`) consists of three tables:

**`users` table**
- Primary storage for user accounts and preferences
- JSON storage for `email_preferences` column

**`vacation_requests` table**
- Stores all vacation requests
- Foreign key to `users` table (`user_id`)
- Indexed by `user_id` and `status` for performance

**`settings` table**
- Single-row table (ID=1)
- Stores system-wide configuration
- JSON storage for `weekend_policy` and `newsletter` columns

---

## 13. Environment Variables Reference

### 13.1 Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `JWT_SECRET` | Secret key for signing JWT tokens. Use a strong random string (32+ chars). | `a7f8e2c9b4d1...` |
| `ADMIN_PASSWORD` | Initial password for default admin account. Change after first login. | `SecureP@ss123!` |

### 13.2 Email Configuration (Resend)

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `RESEND_API_KEY` | API key from Resend dashboard | None (emails disabled) | `re_abc123...` |
| `EMAIL_FROM_ADDRESS` | Sender email address | `noreply@example.com` | `vacay@company.com` |
| `EMAIL_FROM_NAME` | Sender display name | `VacayTracker` | `HR Team` |

### 13.3 Server Configuration

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `PORT` | HTTP server port | `3000` | `8080` |
| `ENV` | Environment mode | `development` | `production` |
| `APP_URL` | Base URL for email links | `http://localhost:3000` | `https://vacay.company.com` |

### 13.4 Database Configuration

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `DB_PATH` | Path to SQLite database file | `./data/vacaytracker.db` | `/var/data/vacay.db` |

### 13.5 Newsletter Defaults

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `NEWSLETTER_DAY` | Default day of month for newsletter | `1` | `15` |
| `NEWSLETTER_HOUR` | Default hour (24h) for newsletter | `9` | `14` |

### 13.6 Security Configuration

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `TOKEN_EXPIRY` | JWT token lifetime | `24h` | `8h` |
| `BCRYPT_COST` | Password hashing cost factor | `10` | `12` |
| `CORS_ORIGINS` | Allowed CORS origins (comma-separated) | `*` | `https://app.company.com` |

### 13.7 Example .env File

```bash
# ===================
# Required
# ===================
JWT_SECRET=your-super-secret-jwt-key-min-32-chars
ADMIN_PASSWORD=ChangeThisSecurePassword!

# ===================
# Email (Resend)
# ===================
RESEND_API_KEY=re_xxxxxxxxxxxxxxxxxxxx
EMAIL_FROM_ADDRESS=vacaytracker@yourcompany.com
EMAIL_FROM_NAME=VacayTracker

# ===================
# Server
# ===================
ENV=development
PORT=3000
APP_URL=http://localhost:3000

# ===================
# Database
# ===================
DB_PATH=./data/vacaytracker.db

# ===================
# Security
# ===================
TOKEN_EXPIRY=24h
BCRYPT_COST=10
CORS_ORIGINS=http://localhost:5173

# ===================
# Newsletter
# ===================
NEWSLETTER_DAY=1
NEWSLETTER_HOUR=9
```

---

## 14. API Error Response Format

### 14.1 Standard Error Response

All API errors follow a consistent JSON structure:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Human-readable error description",
    "details": {}
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| error.code | string | Machine-readable error code (see section 14.3) |
| error.message | string | User-friendly error description |
| error.details | object | Optional additional context (field errors, constraints, etc.) |

### 14.2 HTTP Status Codes

| Status | Meaning | When Used |
|--------|---------|-----------|
| 200 | OK | Successful GET, PUT requests |
| 201 | Created | Successful POST creating a resource |
| 204 | No Content | Successful DELETE |
| 400 | Bad Request | Invalid input, validation failure |
| 401 | Unauthorized | Missing or invalid JWT token |
| 403 | Forbidden | Valid token but insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Duplicate resource (e.g., username taken) |
| 422 | Unprocessable Entity | Business logic violation |
| 500 | Internal Server Error | Unexpected server error |

### 14.3 Error Codes Reference

**Authentication Errors (401)**

| Code | Message | Cause |
|------|---------|-------|
| `AUTH_TOKEN_MISSING` | Authentication token required | No Authorization header |
| `AUTH_TOKEN_INVALID` | Invalid authentication token | Malformed or tampered token |
| `AUTH_TOKEN_EXPIRED` | Authentication token has expired | Token past 24h lifetime |
| `AUTH_CREDENTIALS_INVALID` | Invalid username or password | Failed login attempt |
| `AUTH_ROLE_MISMATCH` | User role does not match selected login type | Employee trying admin tab or vice versa |

**Authorization Errors (403)**

| Code | Message | Cause |
|------|---------|-------|
| `FORBIDDEN_ADMIN_ONLY` | Admin access required | Employee accessing admin endpoint |
| `FORBIDDEN_SELF_DELETE` | Cannot delete your own account | Admin trying to delete themselves |
| `FORBIDDEN_SELF_ROLE_CHANGE` | Cannot change your own role | Admin trying to demote themselves |

**Validation Errors (400)**

| Code | Message | Details |
|------|---------|---------|
| `VALIDATION_REQUIRED_FIELD` | Required field missing | `{ "field": "username" }` |
| `VALIDATION_INVALID_FORMAT` | Invalid field format | `{ "field": "email", "expected": "valid email" }` |
| `VALIDATION_MIN_LENGTH` | Field too short | `{ "field": "password", "min": 6 }` |
| `VALIDATION_MAX_LENGTH` | Field too long | `{ "field": "reason", "max": 200 }` |
| `VALIDATION_INVALID_DATE` | Invalid date format | `{ "field": "startDate", "expected": "YYYY-MM-DD" }` |
| `VALIDATION_DATE_RANGE` | End date before start date | `{ "startDate": "...", "endDate": "..." }` |
| `VALIDATION_PAST_DATE` | Date cannot be in the past | `{ "field": "startDate" }` |

**Resource Errors (404/409)**

| Code | Message | Cause |
|------|---------|-------|
| `USER_NOT_FOUND` | User not found | Invalid user ID |
| `VACATION_REQUEST_NOT_FOUND` | Vacation request not found | Invalid request ID |
| `USERNAME_TAKEN` | Username already exists | Duplicate username on create/update |

**Business Logic Errors (422)**

| Code | Message | Details |
|------|---------|---------|
| `INSUFFICIENT_VACATION_DAYS` | Not enough vacation days remaining | `{ "requested": 5, "available": 3 }` |
| `VACATION_ALREADY_REVIEWED` | Request has already been processed | `{ "currentStatus": "approved" }` |
| `PASSWORD_MISMATCH` | Current password is incorrect | Password change verification failed |
| `OVERLAPPING_VACATION` | Vacation dates overlap existing request | `{ "conflictingRequestId": "..." }` |

**Server Errors (500)**

| Code | Message | Cause |
|------|---------|-------|
| `INTERNAL_ERROR` | An unexpected error occurred | Unhandled exception |
| `DATABASE_ERROR` | Database operation failed | LowDB read/write failure |
| `EMAIL_SERVICE_ERROR` | Failed to send email notification | Resend API failure |

### 14.4 Error Response Examples

**Validation Error (400)**

```json
{
  "error": {
    "code": "VALIDATION_REQUIRED_FIELD",
    "message": "Required field missing",
    "details": {
      "field": "username"
    }
  }
}
```

**Authentication Error (401)**

```json
{
  "error": {
    "code": "AUTH_TOKEN_EXPIRED",
    "message": "Authentication token has expired"
  }
}
```

**Business Logic Error (422)**

```json
{
  "error": {
    "code": "INSUFFICIENT_VACATION_DAYS",
    "message": "Not enough vacation days remaining",
    "details": {
      "requested": 10,
      "available": 7,
      "shortfall": 3
    }
  }
}
```

**Conflict Error (409)**

```json
{
  "error": {
    "code": "USERNAME_TAKEN",
    "message": "Username already exists",
    "details": {
      "username": "john"
    }
  }
}
```

### 14.5 Success Response Patterns

**Single Resource (GET, PUT)**

```json
{
  "data": {
    "id": "usr_a1b2c3d4",
    "name": "John Doe",
    ...
  }
}
```

**Collection (GET list)**

```json
{
  "data": [
    { "id": "usr_a1b2c3d4", ... },
    { "id": "usr_e5f6g7h8", ... }
  ],
  "meta": {
    "total": 25,
    "count": 2
  }
}
```

**Created Resource (POST)**

```json
{
  "data": {
    "id": "vac_newid123",
    ...
  },
  "message": "Vacation request submitted successfully"
}
```

**Delete (204 No Content)**

No response body.

---

*These sections should be inserted after Section 11 (Infrastructure Features) in the main documentation.*

---

## Default Credentials

| Role | Username | Password | Notes |
|------|----------|----------|-------|
| Admin | admin | ADMIN_PASSWORD env var (default: admin123) | Full system access |
| Employee | john | john123 | Test employee with 25 vacation days |
| Employee | jane | jane123 | Test employee with 25 vacation days |

**Important**: Change all default passwords after initial deployment.

---

## API Endpoints Summary

| Method | Endpoint | Auth | Role | Purpose |
|--------|----------|------|------|---------|
| GET | `/health` | No | - | Health check |
| POST | `/api/auth/login` | No | - | User login |
| GET | `/api/auth/me` | Yes | Any | Get current user info |
| PUT | `/api/users/:id/password` | Yes | Self | Change password |
| PUT | `/api/users/:id/email-preferences` | Yes | Self | Update email preferences |
| GET | `/api/vacation` | Yes | Employee | Get my vacation requests |
| POST | `/api/vacation` | Yes | Employee | Submit vacation request |
| DELETE | `/api/vacation/:id` | Yes | Employee | Cancel pending request |
| GET | `/api/admin/users` | Yes | Admin | List all users |
| POST | `/api/admin/users` | Yes | Admin | Create user |
| PUT | `/api/admin/users/:id` | Yes | Admin | Update user |
| DELETE | `/api/admin/users/:id` | Yes | Admin | Delete user |
| GET | `/api/admin/vacation` | Yes | Admin | List all vacation requests |
| GET | `/api/admin/vacation/pending` | Yes | Admin | List pending requests |
| PUT | `/api/admin/vacation/:id/review` | Yes | Admin | Approve/reject request |
| POST | `/api/admin/reset-vacation-days` | Yes | Admin | Reset all vacation days |

---

*Documentation generated from comprehensive codebase analysis of VacayTracker v1.0*
