# 06 - Component Inventory

> Complete inventory of backend and frontend components with specifications

## Table of Contents

1. [Backend Components](#backend-components)
2. [Frontend Components](#frontend-components)
3. [Component Dependencies](#component-dependencies)

---

## Backend Components

### Domain Entities

| Entity | File | Fields | Methods |
|--------|------|--------|---------|
| **User** | `internal/domain/user.go` | id, email, passwordHash, name, role, vacationBalance, startDate, emailPreferences, createdAt, updatedAt | IsAdmin(), IsEmployee() |
| **VacationRequest** | `internal/domain/vacation.go` | id, userId, startDate, endDate, totalDays, reason, status, reviewedBy, reviewedAt, rejectionReason, createdAt, updatedAt | IsPending(), IsApproved(), IsRejected(), CanBeCancelled() |
| **TeamVacation** | `internal/domain/vacation.go` | id, userId, userName, startDate, endDate, totalDays | - |
| **Settings** | `internal/domain/settings.go` | id, weekendPolicy, newsletter, defaultVacationDays, vacationResetMonth, updatedAt | - |
| **EmailPreferences** | `internal/domain/user.go` | vacationUpdates, weeklyDigest, teamNotifications | ToJSONString() |
| **WeekendPolicy** | `internal/domain/settings.go` | excludeWeekends, excludedDays | ToJSONString() |
| **NewsletterConfig** | `internal/domain/settings.go` | enabled, frequency, dayOfMonth | ToJSONString() |

### Repositories

| Repository | File | Methods |
|------------|------|---------|
| **UserRepository** | `internal/repository/sqlite/user.go` | Create, GetByID, GetByEmail, List, Update, UpdatePassword, UpdateBalance, UpdateEmailPreferences, Delete, CountByRole, EmailExists |
| **VacationRepository** | `internal/repository/sqlite/vacation.go` | Create, GetByID, ListByUser, ListPending, ListTeam, UpdateStatus, Delete |
| **SettingsRepository** | `internal/repository/sqlite/settings.go` | Get, Update |

### Services

| Service | File | Dependencies | Methods |
|---------|------|--------------|---------|
| **AuthService** | `internal/service/auth.go` | UserRepository, Config | HashPassword, VerifyPassword, GenerateToken, ValidateToken, Login, ChangePassword, CreateInitialAdmin |
| **UserService** | `internal/service/user.go` | UserRepository, AuthService | Create, Update, Delete |
| **VacationService** | `internal/service/vacation.go` | VacationRepository, UserRepository, SettingsRepository | Create, Cancel, Approve, Reject |
| **EmailService** | `internal/service/email.go` | Config | SendWelcome, SendRequestSubmitted, SendRequestApproved, SendRequestRejected, SendAdminNotification |

### Handlers

| Handler | File | Methods | Routes |
|---------|------|---------|--------|
| **HealthHandler** | `internal/handler/health.go` | Check | GET /health |
| **AuthHandler** | `internal/handler/auth.go` | Login, Me, ChangePassword, UpdateEmailPreferences | POST /api/auth/login, GET /api/auth/me, PUT /api/auth/password, PUT /api/auth/email-preferences |
| **VacationHandler** | `internal/handler/vacation.go` | Create, List, Get, Cancel, Team | POST /api/vacation/request, GET /api/vacation/requests, GET /api/vacation/requests/:id, DELETE /api/vacation/requests/:id, GET /api/vacation/team |
| **AdminHandler** | `internal/handler/admin.go` | ListUsers, CreateUser, GetUser, UpdateUser, DeleteUser, PendingRequests, ApproveRequest, RejectRequest, GetSettings, UpdateSettings | /api/admin/* |

### Middleware

| Middleware | File | Purpose |
|------------|------|---------|
| **RequireAuth** | `internal/middleware/auth.go` | Validates JWT token, sets user context |
| **RequireAdmin** | `internal/middleware/auth.go` | Checks admin role |
| **CORS** | `internal/middleware/cors.go` | Handles cross-origin requests |
| **ErrorHandler** | `internal/middleware/error.go` | Recovers from panics, handles errors |

### DTOs

| DTO | File | Fields |
|-----|------|--------|
| **LoginRequest** | `internal/dto/request.go` | email, password |
| **ChangePasswordRequest** | `internal/dto/request.go` | currentPassword, newPassword |
| **UpdateEmailPreferencesRequest** | `internal/dto/request.go` | vacationUpdates, weeklyDigest, teamNotifications |
| **CreateUserRequest** | `internal/dto/request.go` | email, password, name, role, vacationBalance, startDate |
| **UpdateUserRequest** | `internal/dto/request.go` | email, name, role, vacationBalance, startDate |
| **CreateVacationRequest** | `internal/dto/request.go` | startDate, endDate, reason |
| **RejectVacationRequest** | `internal/dto/request.go` | reason |
| **UpdateSettingsRequest** | `internal/dto/request.go` | weekendPolicy, newsletter, defaultVacationDays, vacationResetMonth |
| **LoginResponse** | `internal/dto/response.go` | token, user |
| **UserResponse** | `internal/dto/response.go` | id, email, name, role, vacationBalance, startDate, emailPreferences, createdAt, updatedAt |
| **ErrorResponse** | `internal/dto/errors.go` | code, message, details |
| **JWTClaims** | `internal/dto/response.go` | userID, email, name, role |

---

## Frontend Components

### Base UI Components

| Component | File | Props | State | Events |
|-----------|------|-------|-------|--------|
| **Button** | `src/lib/components/ui/Button.svelte` | variant, size, disabled, loading, type, class, children | - | onclick |
| **Input** | `src/lib/components/ui/Input.svelte` | type, value, placeholder, label, error, disabled, required, id, name, class | - | oninput, onblur |
| **Card** | `src/lib/components/ui/Card.svelte` | padding, class, header (snippet), footer (snippet), children | - | - |
| **Badge** | `src/lib/components/ui/Badge.svelte` | variant, size, class, children | - | - |
| **Avatar** | `src/lib/components/ui/Avatar.svelte` | name, src, size, class | initials (derived) | - |
| **ProgressRing** | `src/lib/components/ui/ProgressRing.svelte` | value, max, size, strokeWidth, class | percentage, radius, circumference, offset, color (all derived) | - |
| **Toast** | `src/lib/components/ui/Toast.svelte` | - | toasts (from store) | - |
| **Pagination** | `src/lib/components/ui/Pagination.svelte` | page, totalPages, class | - | onPageChange |

### Layout Components

| Component | File | Props | Description |
|-----------|------|-------|-------------|
| **EmployeeHeader** | `src/lib/components/layout/EmployeeHeader.svelte` | - | Navigation header for employee pages |
| **AdminSidebar** | `src/lib/components/layout/AdminSidebar.svelte` | - | Navigation sidebar for admin pages |
| **AdminHeader** | `src/lib/components/layout/AdminHeader.svelte` | - | Top header bar for admin pages |
| **UserMenu** | `src/lib/components/layout/UserMenu.svelte` | - | Dropdown menu with user actions |
| **MobileNav** | `src/lib/components/layout/MobileNav.svelte` | - | Mobile navigation drawer |

### Feature Components - Vacation

| Component | File | Props | State | Description |
|-----------|------|-------|-------|-------------|
| **BalanceDisplay** | `src/lib/components/features/vacation/BalanceDisplay.svelte` | current, total, size | - | Shows vacation balance as progress ring |
| **RequestModal** | `src/lib/components/features/vacation/RequestModal.svelte` | open (bindable) | startDate, endDate, reason, isSubmitting, errors | Modal form for creating vacation request |
| **RequestList** | `src/lib/components/features/vacation/RequestList.svelte` | requests | - | List of vacation request cards |
| **RequestCard** | `src/lib/components/features/vacation/RequestCard.svelte` | request | isCancelling | Individual vacation request display |
| **DateRangePicker** | `src/lib/components/features/vacation/DateRangePicker.svelte` | startDate, endDate, minDate | - | Date range selection with validation |

### Feature Components - Admin

| Component | File | Props | State | Description |
|-----------|------|-------|-------|-------------|
| **StatsCard** | `src/lib/components/features/admin/StatsCard.svelte` | title, value, icon, color | - | Dashboard statistics card |
| **PendingRequests** | `src/lib/components/features/admin/PendingRequests.svelte` | requests, onUpdate | processingId | List of pending requests with approve/reject |
| **ApprovalModal** | `src/lib/components/features/admin/ApprovalModal.svelte` | open, request, onConfirm | - | Confirmation modal for approval |
| **RejectionModal** | `src/lib/components/features/admin/RejectionModal.svelte` | open, request, onConfirm | reason | Modal with rejection reason input |
| **UserTable** | `src/lib/components/features/admin/UserTable.svelte` | users, pagination, onEdit, onDelete | sortField, sortDirection | Sortable user list table |
| **UserModal** | `src/lib/components/features/admin/UserModal.svelte` | open, user, onSave | form, errors, isSubmitting | Create/edit user form |
| **DeleteConfirm** | `src/lib/components/features/admin/DeleteConfirm.svelte` | open, user, onConfirm | isDeleting | Delete confirmation dialog |
| **BalanceEditor** | `src/lib/components/features/admin/BalanceEditor.svelte` | users, onUpdate | editingId | Inline balance editing |
| **WeekendPolicy** | `src/lib/components/features/admin/WeekendPolicy.svelte` | policy, onUpdate | - | Weekend day toggles |
| **NewsletterSettings** | `src/lib/components/features/admin/NewsletterSettings.svelte` | config, onUpdate | - | Newsletter configuration |

### Feature Components - Calendar

| Component | File | Props | State | Description |
|-----------|------|-------|-------|-------------|
| **MonthView** | `src/lib/components/features/calendar/MonthView.svelte` | year, month, vacations | days (derived) | Month grid calendar view |
| **WeekView** | `src/lib/components/features/calendar/WeekView.svelte` | year, week, vacations | - | Week grid calendar view |
| **DayCell** | `src/lib/components/features/calendar/DayCell.svelte` | day, vacations, isToday | - | Single day in calendar |
| **VacationEvent** | `src/lib/components/features/calendar/VacationEvent.svelte` | vacation | - | Vacation event display |
| **CalendarHeader** | `src/lib/components/features/calendar/CalendarHeader.svelte` | currentDate, view, onNavigate, onViewChange | - | Calendar navigation controls |
| **FilterPanel** | `src/lib/components/features/calendar/FilterPanel.svelte` | users, selectedUsers, onFilter | - | User filter for calendar |

### Feature Components - Auth

| Component | File | Props | State | Description |
|-----------|------|-------|-------|-------------|
| **PasswordChange** | `src/lib/components/features/auth/PasswordChange.svelte` | - | currentPassword, newPassword, confirmPassword, isSubmitting, errors | Password change form |
| **EmailPreferences** | `src/lib/components/features/auth/EmailPreferences.svelte` | preferences, onUpdate | - | Email notification toggles |

### Stores

| Store | File | State | Derived | Methods |
|-------|------|-------|---------|---------|
| **auth** | `src/lib/stores/auth.svelte.ts` | user, isLoading, error | isAuthenticated, isAdmin, isEmployee | initialize, login, logout, updateUser |
| **toast** | `src/lib/stores/toast.svelte.ts` | toasts | - | add, dismiss, dismissAll, success, error, warning, info |
| **vacation** | `src/lib/stores/vacation.svelte.ts` | requests, isLoading, error | pendingRequests, approvedRequests, rejectedRequests | fetchRequests, createRequest, cancelRequest, updateRequest |
| **admin** | `src/lib/stores/admin.svelte.ts` | users, settings, isLoading | - | fetchUsers, fetchSettings, updateSettings |
| **calendar** | `src/lib/stores/calendar.svelte.ts` | vacations, currentDate, view, cache | currentMonth, currentYear | fetchVacations, navigate, setView |

### API Modules

| Module | File | Methods |
|--------|------|---------|
| **client** | `src/lib/api/client.ts` | request, setAuthToken, clearAuthToken |
| **authApi** | `src/lib/api/auth.ts` | login, logout, me, changePassword, updateEmailPreferences |
| **vacationApi** | `src/lib/api/vacation.ts` | create, list, get, cancel, team |
| **adminApi** | `src/lib/api/admin.ts` | listUsers, createUser, getUser, updateUser, deleteUser, pendingRequests, approveRequest, rejectRequest, getSettings, updateSettings, sendNewsletter |

### Routes

| Route | File | Auth | Description |
|-------|------|------|-------------|
| `/` | `src/routes/+page.svelte` | Public | Login page |
| `/employee` | `src/routes/employee/+page.svelte` | Employee | Employee dashboard |
| `/employee/calendar` | `src/routes/employee/calendar/+page.svelte` | Employee | Team calendar |
| `/employee/team` | `src/routes/employee/team/+page.svelte` | Employee | Team overview |
| `/employee/settings` | `src/routes/employee/settings/+page.svelte` | Employee | Profile settings |
| `/admin` | `src/routes/admin/+page.svelte` | Admin | Admin dashboard |
| `/admin/users` | `src/routes/admin/users/+page.svelte` | Admin | User management |
| `/admin/balances` | `src/routes/admin/balances/+page.svelte` | Admin | Balance management |
| `/admin/calendar` | `src/routes/admin/calendar/+page.svelte` | Admin | Full calendar view |
| `/admin/settings` | `src/routes/admin/settings/+page.svelte` | Admin | App settings |

---

## Component Dependencies

### Backend Dependency Graph

```
main.go
├── config.Config
├── sqlite.DB
│   ├── sqlite.UserRepository
│   ├── sqlite.VacationRepository
│   └── sqlite.SettingsRepository
├── service.AuthService
│   └── sqlite.UserRepository
├── service.UserService
│   ├── sqlite.UserRepository
│   └── service.AuthService
├── service.VacationService
│   ├── sqlite.VacationRepository
│   ├── sqlite.UserRepository
│   └── sqlite.SettingsRepository
├── handler.HealthHandler
├── handler.AuthHandler
│   ├── service.AuthService
│   └── sqlite.UserRepository
├── handler.VacationHandler
│   ├── service.VacationService
│   └── sqlite.VacationRepository
├── handler.AdminHandler
│   ├── service.UserService
│   ├── service.VacationService
│   ├── sqlite.UserRepository
│   ├── sqlite.VacationRepository
│   └── sqlite.SettingsRepository
└── middleware.AuthMiddleware
    └── service.AuthService
```

### Frontend Dependency Graph

```
+layout.svelte
├── auth store
├── Toast component
│   └── toast store
└── routes/*
    ├── +page.svelte (login)
    │   ├── auth store
    │   ├── toast store
    │   ├── Button, Input, Card
    │   └── Lucide icons
    ├── employee/+layout.svelte
    │   ├── auth store
    │   └── EmployeeHeader
    │       └── Avatar, Button
    ├── employee/+page.svelte
    │   ├── auth store
    │   ├── vacation store
    │   ├── BalanceDisplay
    │   │   └── ProgressRing
    │   ├── RequestList
    │   │   └── RequestCard
    │   │       └── Badge, Button
    │   └── RequestModal
    │       └── Button, Input
    ├── employee/calendar/+page.svelte
    │   ├── vacationApi
    │   ├── MonthView
    │   │   └── DayCell
    │   └── Button, Card
    ├── admin/+layout.svelte
    │   ├── auth store
    │   ├── AdminSidebar
    │   └── AdminHeader
    ├── admin/+page.svelte
    │   ├── adminApi
    │   ├── StatsCard
    │   ├── PendingRequests
    │   │   ├── Avatar, Badge, Button
    │   │   └── toast store
    │   └── Card
    └── admin/users/+page.svelte
        ├── adminApi
        ├── UserTable
        ├── UserModal
        ├── DeleteConfirm
        └── Pagination
```

---

## Related Documents

- [04-backend-tasks.md](./04-backend-tasks.md) - Backend implementation details
- [05-frontend-tasks.md](./05-frontend-tasks.md) - Frontend implementation details
- [07-testing-strategy.md](./07-testing-strategy.md) - Component testing
