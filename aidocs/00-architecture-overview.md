# VacayTracker Architecture Overview

> **Purpose:** Comprehensive system architecture documentation for AI agent implementation
> **Audience:** Claude/AI agents executing implementation tasks
> **Format:** Self-contained, explicit, actionable

---

## Table of Contents

1. [System Overview](#1-system-overview)
2. [Technology Stack](#2-technology-stack)
3. [Architecture Diagrams](#3-architecture-diagrams)
4. [Component Relationships](#4-component-relationships)
5. [Data Flow](#5-data-flow)
6. [Directory Structures](#6-directory-structures)
7. [Key Design Decisions](#7-key-design-decisions)
8. [Cross-Reference Guide](#8-cross-reference-guide)

---

## 1. System Overview

### 1.1 Application Description

VacayTracker is an employee vacation tracking application with:
- **Role-based access control** (Admin/Employee)
- **Email notifications** via Resend API
- **Team calendar visualization**
- **Beach/vacation theme** throughout UI

### 1.2 Core Features

| Feature | Employee | Admin |
|---------|----------|-------|
| View vacation balance | Yes | Yes (all employees) |
| Submit vacation requests | Yes | No |
| Approve/reject requests | No | Yes |
| View team calendar | Yes | Yes |
| Manage users | No | Yes |
| Configure settings | No | Yes |
| Receive notifications | Yes | Yes |

### 1.3 User Roles

```
┌─────────────────────────────────────────────────────────────────┐
│                         USER ROLES                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│   ADMIN (Captain)                    EMPLOYEE (Crew Member)       │
│   ───────────────                    ─────────────────────        │
│   • Full system access               • Submit vacation requests   │
│   • Approve/reject requests          • View own balance           │
│   • Manage all users                 • View team calendar         │
│   • Configure settings               • Change own password        │
│   • Reset vacation days              • Manage email preferences   │
│   • Send newsletters                                              │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. Technology Stack

### 2.1 Frontend Stack

| Technology | Version | Purpose |
|------------|---------|---------|
| **Svelte** | 5.x | UI framework with runes-based reactivity |
| **SvelteKit** | 2.49.x | Full-stack framework (routing, SSR) |
| **Melt UI** | 0.42.x (`melt`) | Headless accessible components |
| **Tailwind CSS** | 4.x | Utility-first styling with `@theme` |
| **TypeScript** | 5.x | Type safety |
| **Vite** | 6.x | Build tool and dev server |
| **Lucide Svelte** | latest | Icon library |

### 2.2 Backend Stack

| Technology | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.23+ | Core runtime |
| **Gin** | 1.10+ | HTTP routing and middleware |
| **SQLite** | - | Embedded database |
| **modernc.org/sqlite** | 1.34+ | CGo-free SQLite driver |
| **golang-jwt/jwt** | v5 | JWT authentication |
| **golang.org/x/crypto** | - | bcrypt password hashing |
| **resend-go** | v2+ | Transactional emails |
| **godotenv** | 1.5+ | Environment variables |
| **go-playground/validator** | v10 | Input validation |
| **google/uuid** | 1.6+ | ID generation |

### 2.3 Infrastructure

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Database** | SQLite (WAL mode) | Data persistence |
| **Containerization** | Docker | Deployment |
| **Development** | Docker Compose | Local environment |
| **Email Service** | Resend API | Transactional emails |

---

## 3. Architecture Diagrams

### 3.1 High-Level System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           VACAYTRACKER SYSTEM                                │
└─────────────────────────────────────────────────────────────────────────────┘

                                    ┌──────────────┐
                                    │   Browser    │
                                    │   (Client)   │
                                    └──────┬───────┘
                                           │
                                           │ HTTPS
                                           │
                            ┌──────────────▼───────────────┐
                            │      FRONTEND (SvelteKit)     │
                            │      Port 5173 (dev)          │
                            │      ─────────────────────    │
                            │  • Svelte 5 components        │
                            │  • Client-side routing        │
                            │  • Auth state management      │
                            │  • API client                 │
                            └──────────────┬───────────────┘
                                           │
                                           │ HTTP/JSON
                                           │ Authorization: Bearer <JWT>
                                           │
                            ┌──────────────▼───────────────┐
                            │       BACKEND (Go/Gin)        │
                            │       Port 3000               │
                            │       ─────────────────────   │
                            │  • REST API handlers          │
                            │  • JWT authentication         │
                            │  • Business logic services    │
                            │  • Data access repositories   │
                            └──────────────┬───────────────┘
                                           │
                      ┌────────────────────┼────────────────────┐
                      │                    │                    │
                      ▼                    ▼                    ▼
             ┌────────────────┐   ┌────────────────┐   ┌────────────────┐
             │    SQLite      │   │  Resend API    │   │   Scheduler    │
             │   Database     │   │  (Email)       │   │  (Newsletter)  │
             │ ────────────── │   │ ────────────── │   │ ────────────── │
             │ • Users        │   │ • Welcome      │   │ • Monthly      │
             │ • Vacations    │   │ • Status       │   │   summary      │
             │ • Settings     │   │ • Reset        │   │                │
             └────────────────┘   └────────────────┘   └────────────────┘
```

### 3.2 Backend Layer Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         BACKEND ARCHITECTURE (Go)                            │
└─────────────────────────────────────────────────────────────────────────────┘

    HTTP Request
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              MIDDLEWARE LAYER                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │    CORS     │  │    Auth     │  │   Admin     │  │   Error     │        │
│  │  Middleware │  │  Middleware │  │  Middleware │  │  Middleware │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              HANDLER LAYER                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │    Auth     │  │    User     │  │  Vacation   │  │   Admin     │        │
│  │   Handler   │  │   Handler   │  │   Handler   │  │   Handler   │        │
│  │ ─────────── │  │ ─────────── │  │ ─────────── │  │ ─────────── │        │
│  │ POST /login │  │ PUT /pass   │  │ GET /vac    │  │ CRUD /users │        │
│  │ GET /me     │  │ PUT /prefs  │  │ POST /vac   │  │ PUT /review │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              SERVICE LAYER                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │    Auth     │  │    User     │  │  Vacation   │  │   Email     │        │
│  │   Service   │  │   Service   │  │   Service   │  │   Service   │        │
│  │ ─────────── │  │ ─────────── │  │ ─────────── │  │ ─────────── │        │
│  │ Login       │  │ Create      │  │ Create      │  │ SendApprove │        │
│  │ ValidateJWT │  │ Update      │  │ Review      │  │ SendReject  │        │
│  │ HashPass    │  │ Delete      │  │ CalcDays    │  │ SendWelcome │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            REPOSITORY LAYER                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │ UserRepository  │  │VacationRepository│ │SettingsRepository│             │
│  │ ─────────────── │  │ ─────────────── │  │ ─────────────── │              │
│  │ Create          │  │ Create          │  │ Get             │              │
│  │ GetByID         │  │ GetByID         │  │ Update          │              │
│  │ GetByUsername   │  │ GetByUserID     │  │                 │              │
│  │ GetAll          │  │ GetAll          │  │                 │              │
│  │ Update          │  │ GetPending      │  │                 │              │
│  │ Delete          │  │ Update          │  │                 │              │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘              │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              DATABASE LAYER                                  │
│                     ┌─────────────────────────────┐                          │
│                     │         SQLite              │                          │
│                     │     (modernc.org/sqlite)    │                          │
│                     │ ─────────────────────────── │                          │
│                     │ • WAL mode enabled          │                          │
│                     │ • Foreign keys ON           │                          │
│                     │ • Single-writer (1 conn)    │                          │
│                     └─────────────────────────────┘                          │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Frontend Component Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       FRONTEND ARCHITECTURE (Svelte)                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                              ROUTES LAYER                                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   / (Login) │  │  /employee  │  │   /admin    │  │  /settings  │        │
│  │   +page     │  │   +layout   │  │   +layout   │  │   +page     │        │
│  │             │  │   +page     │  │   +page     │  │             │        │
│  │             │  │   /history  │  │   /users    │  │             │        │
│  │             │  │   /settings │  │   /calendar │  │             │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           FEATURE COMPONENTS                                 │
│  ┌───────────────────┐  ┌───────────────────┐  ┌───────────────────┐       │
│  │       Auth        │  │     Vacation      │  │       Admin       │       │
│  │ ───────────────── │  │ ───────────────── │  │ ───────────────── │       │
│  │ • LoginForm       │  │ • BalanceDisplay  │  │ • StatsCard       │       │
│  │                   │  │ • RequestModal    │  │ • UserCard        │       │
│  │                   │  │ • RequestCard     │  │ • ResetModal      │       │
│  │                   │  │ • Timeline        │  │                   │       │
│  │                   │  │ • Calendar        │  │                   │       │
│  └───────────────────┘  └───────────────────┘  └───────────────────┘       │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                             BASE UI COMPONENTS                               │
│  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐   │
│  │ Button │  │ Input  │  │  Card  │  │ Badge  │  │ Avatar │  │Progress│   │
│  │        │  │        │  │        │  │        │  │        │  │  Ring  │   │
│  └────────┘  └────────┘  └────────┘  └────────┘  └────────┘  └────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              STATE LAYER                                     │
│  ┌─────────────────────────────┐  ┌─────────────────────────────┐          │
│  │         Auth Store          │  │        Toast Store          │          │
│  │ ─────────────────────────── │  │ ─────────────────────────── │          │
│  │ • user: User | null         │  │ • toasts: Toast[]           │          │
│  │ • token: string | null      │  │ • add(type, message)        │          │
│  │ • isAuthenticated           │  │ • remove(id)                │          │
│  │ • isAdmin / isEmployee      │  │ • success/error/warning     │          │
│  │ • setSession() / logout()   │  │                             │          │
│  └─────────────────────────────┘  └─────────────────────────────┘          │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              API LAYER                                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   client    │  │   authApi   │  │ vacationApi │  │   usersApi  │        │
│  │ ─────────── │  │ ─────────── │  │ ─────────── │  │ ─────────── │        │
│  │ get/post    │  │ login       │  │ getRequests │  │ getAll      │        │
│  │ put/delete  │  │ getMe       │  │ create      │  │ create      │        │
│  │ auth header │  │             │  │ cancel      │  │ update      │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Component Relationships

### 4.1 Backend Dependency Graph

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    BACKEND COMPONENT DEPENDENCIES                            │
└─────────────────────────────────────────────────────────────────────────────┘

main.go
    │
    ├──► config.Config
    │
    ├──► sqlite.DB ──────────────────────────────────────────┐
    │                                                         │
    ├──► UserRepository ◄─────────────────────────────────────┤
    │         │                                               │
    │         ├──► AuthService ──► AuthMiddleware             │
    │         │         │                                     │
    │         │         ├──► AuthHandler                      │
    │         │         │                                     │
    │         ├──► UserService ──► UserHandler                │
    │         │         │                                     │
    │         │         └──► AdminHandler ◄──┐                │
    │         │                              │                │
    ├──► VacationRepository ◄────────────────┤                │
    │         │                              │                │
    │         └──► VacationService ──────────┤                │
    │                   │                    │                │
    │                   └──► VacationHandler │                │
    │                                        │                │
    ├──► SettingsRepository ◄────────────────┤                │
    │         │                              │                │
    │         └──► VacationService (settings)│                │
    │                                        │                │
    └──► EmailService ◄──────────────────────┘                │
              │                                               │
              └── (handlers call for notifications)           │
                                                              │
              All repositories depend on ──────────────────────┘
```

### 4.2 Frontend Component Hierarchy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                   FRONTEND COMPONENT HIERARCHY                               │
└─────────────────────────────────────────────────────────────────────────────┘

+layout.svelte (Root)
    │
    ├── ToastProvider
    │
    └── +page.svelte (Login) ──► LoginForm
                                    │
                                    ├── Tabs (Melt UI)
                                    ├── Input
                                    └── Button

/employee/+layout.svelte ──► Header, Sidebar
    │
    ├── +page.svelte (Dashboard)
    │       │
    │       ├── BalanceDisplay ──► ProgressRing
    │       ├── RequestModal ──► Dialog (Melt UI)
    │       │       │               ├── Input (date)
    │       │       │               ├── Button
    │       │       │               └── Textarea
    │       └── Timeline ──► RequestCard
    │                           └── Badge
    │
    ├── /history/+page.svelte
    │       └── RequestCard[]
    │
    └── /settings/+page.svelte
            ├── Input (password)
            └── Button

/admin/+layout.svelte ──► Header, Sidebar
    │
    ├── +page.svelte (Dashboard)
    │       │
    │       ├── StatsCard[]
    │       ├── RequestCard[] (pending)
    │       └── Timeline
    │
    ├── /users/+page.svelte
    │       │
    │       ├── UserCard[]
    │       │       └── Avatar, Badge, Button
    │       └── UserModal (add/edit)
    │               ├── Input[]
    │               ├── Select (Melt UI)
    │               └── Button
    │
    ├── /calendar/+page.svelte
    │       └── Calendar (Melt UI)
    │
    └── /settings/+page.svelte
            ├── Toggle (weekend policy)
            ├── ResetModal
            └── NewsletterSettings
```

---

## 5. Data Flow

### 5.1 Authentication Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         AUTHENTICATION FLOW                                  │
└─────────────────────────────────────────────────────────────────────────────┘

1. LOGIN REQUEST
   ┌──────────┐     POST /api/auth/login      ┌──────────┐
   │  Client  │ ──────────────────────────────►│  Server  │
   │          │    {username, password}        │          │
   └──────────┘                                └────┬─────┘
                                                    │
                                                    ▼
                                        ┌───────────────────┐
                                        │   AuthService     │
                                        │   ───────────────  │
                                        │ 1. GetByUsername   │
                                        │ 2. bcrypt.Compare  │
                                        │ 3. Generate JWT    │
                                        └─────────┬─────────┘
                                                  │
                                                  ▼
   ┌──────────┐     {token, user}          ┌──────────┐
   │  Client  │ ◄──────────────────────────│  Server  │
   │          │                            │          │
   └────┬─────┘                            └──────────┘
        │
        ▼
   ┌────────────────────────────┐
   │   Store token in           │
   │   sessionStorage           │
   │   Update auth store        │
   │   Redirect to dashboard    │
   └────────────────────────────┘

2. AUTHENTICATED REQUESTS
   ┌──────────┐     GET /api/vacation         ┌──────────┐
   │  Client  │ ──────────────────────────────►│  Server  │
   │          │  Authorization: Bearer <jwt>   │          │
   └──────────┘                                └────┬─────┘
                                                    │
                                                    ▼
                                        ┌───────────────────┐
                                        │  AuthMiddleware   │
                                        │  ────────────────  │
                                        │ 1. Extract token   │
                                        │ 2. Validate JWT    │
                                        │ 3. Set claims ctx  │
                                        └─────────┬─────────┘
                                                  │
                                             ┌────┴────┐
                                             ▼         ▼
                                         Valid      Invalid
                                           │           │
                                           ▼           ▼
                                      Continue     401 Error
                                      to Handler   Response
```

### 5.2 Vacation Request Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        VACATION REQUEST FLOW                                 │
└─────────────────────────────────────────────────────────────────────────────┘

1. EMPLOYEE SUBMITS REQUEST
   ┌────────────┐                              ┌────────────┐
   │  Employee  │  POST /api/vacation          │   Server   │
   │   Client   │ ────────────────────────────►│            │
   │            │  {startDate, endDate,        │            │
   │            │   reason}                    │            │
   └────────────┘                              └─────┬──────┘
                                                     │
                                     ┌───────────────▼───────────────┐
                                     │      VacationService          │
                                     │ ───────────────────────────── │
                                     │ 1. Parse & validate dates     │
                                     │ 2. Check overlap (HasOverlap) │
                                     │ 3. Get weekend settings       │
                                     │ 4. Calculate business days    │
                                     │ 5. Check sufficient balance   │
                                     │ 6. Create pending request     │
                                     └───────────────┬───────────────┘
                                                     │
                                     ┌───────────────▼───────────────┐
                                     │       EmailService            │
                                     │ (async, non-blocking)         │
                                     │ ───────────────────────────── │
                                     │ Notify all admins with        │
                                     │ vacationRequestNotifications  │
                                     │ preference enabled            │
                                     └───────────────────────────────┘

2. ADMIN REVIEWS REQUEST
   ┌────────────┐                              ┌────────────┐
   │   Admin    │  PUT /api/admin/vacation/    │   Server   │
   │   Client   │      :id/review              │            │
   │            │ ────────────────────────────►│            │
   │            │  {status: "approved"}        │            │
   └────────────┘                              └─────┬──────┘
                                                     │
                                     ┌───────────────▼───────────────┐
                                     │      VacationService          │
                                     │ ───────────────────────────── │
                                     │ 1. Get request by ID          │
                                     │ 2. Check status is pending    │
                                     │ 3. Update status, reviewer    │
                                     │ 4. If approved:               │
                                     │    - Update user.usedDays     │
                                     └───────────────┬───────────────┘
                                                     │
                                     ┌───────────────▼───────────────┐
                                     │       EmailService            │
                                     │ ───────────────────────────── │
                                     │ Notify employee with          │
                                     │ vacationStatusUpdates         │
                                     │ preference enabled            │
                                     └───────────────────────────────┘
```

---

## 6. Directory Structures

### 6.1 Backend Directory Structure

```
C:\code\vacayv2\vacaytracker-api\
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go                  # Configuration loading from env
│   ├── domain/
│   │   ├── user.go                    # User entity, Role enum
│   │   ├── vacation.go                # VacationRequest entity, Status enum
│   │   └── settings.go                # Settings entity
│   ├── repository/
│   │   ├── repository.go              # Repository interfaces
│   │   └── sqlite/
│   │       ├── sqlite.go              # DB connection, migration
│   │       ├── user.go                # UserRepository implementation
│   │       ├── vacation.go            # VacationRepository implementation
│   │       └── settings.go            # SettingsRepository implementation
│   ├── service/
│   │   ├── auth.go                    # AuthService (JWT, password)
│   │   ├── user.go                    # UserService (CRUD, preferences)
│   │   ├── vacation.go                # VacationService (requests, balance)
│   │   └── email.go                   # EmailService (Resend integration)
│   ├── handler/
│   │   ├── handler.go                 # Handler dependencies struct
│   │   ├── auth.go                    # POST /login, GET /me
│   │   ├── user.go                    # PUT /password, PUT /preferences
│   │   ├── vacation.go                # Employee vacation endpoints
│   │   ├── admin.go                   # Admin endpoints
│   │   └── health.go                  # GET /health
│   ├── middleware/
│   │   ├── auth.go                    # JWT validation, role checking
│   │   ├── cors.go                    # CORS configuration
│   │   └── error.go                   # Error recovery, logging
│   └── dto/
│       ├── request.go                 # Request DTOs with validation
│       ├── response.go                # Response DTOs
│       └── errors.go                  # Error types and codes
├── migrations/
│   └── 001_init.sql                   # Initial database schema
├── data/
│   └── .gitkeep                       # Database directory
├── .env.example                       # Environment template
├── go.mod                             # Go module definition
├── go.sum                             # Dependency checksums
├── Makefile                           # Build commands
└── README.md                          # Backend documentation
```

### 6.2 Frontend Directory Structure

```
C:\code\vacayv2\vacaytracker\
├── src/
│   ├── lib/
│   │   ├── components/
│   │   │   ├── ui/                    # Base UI primitives
│   │   │   │   ├── Button.svelte
│   │   │   │   ├── Input.svelte
│   │   │   │   ├── Card.svelte
│   │   │   │   ├── Badge.svelte
│   │   │   │   ├── Avatar.svelte
│   │   │   │   ├── ProgressRing.svelte
│   │   │   │   ├── ToastProvider.svelte
│   │   │   │   └── index.ts           # Barrel export
│   │   │   ├── layout/                # Layout components
│   │   │   │   ├── Header.svelte
│   │   │   │   ├── Sidebar.svelte
│   │   │   │   ├── PageHeader.svelte
│   │   │   │   └── index.ts
│   │   │   └── features/              # Domain components
│   │   │       ├── auth/
│   │   │       │   └── LoginForm.svelte
│   │   │       ├── vacation/
│   │   │       │   ├── BalanceDisplay.svelte
│   │   │       │   ├── RequestModal.svelte
│   │   │       │   ├── RequestCard.svelte
│   │   │       │   ├── Timeline.svelte
│   │   │       │   └── Calendar.svelte
│   │   │       └── admin/
│   │   │           ├── UserCard.svelte
│   │   │           ├── StatsCard.svelte
│   │   │           └── ResetModal.svelte
│   │   ├── stores/                    # Global state
│   │   │   ├── auth.svelte.ts         # Auth state with runes
│   │   │   └── toast.svelte.ts        # Toast notifications
│   │   ├── api/                       # API client
│   │   │   ├── client.ts              # Base API client
│   │   │   ├── auth.ts                # Auth endpoints
│   │   │   ├── vacation.ts            # Vacation endpoints
│   │   │   └── users.ts               # User endpoints
│   │   ├── types/                     # TypeScript types
│   │   │   └── index.ts
│   │   └── utils/                     # Utilities
│   │       ├── date.ts                # Date formatting
│   │       └── format.ts              # Number formatting
│   ├── routes/
│   │   ├── +layout.svelte             # Root layout
│   │   ├── +page.svelte               # Login page
│   │   ├── employee/
│   │   │   ├── +layout.svelte         # Employee layout
│   │   │   ├── +layout.server.ts      # Auth guard
│   │   │   ├── +page.svelte           # Dashboard
│   │   │   ├── +page.server.ts        # Load data
│   │   │   ├── history/
│   │   │   │   └── +page.svelte
│   │   │   └── settings/
│   │   │       └── +page.svelte
│   │   └── admin/
│   │       ├── +layout.svelte
│   │       ├── +layout.server.ts      # Admin auth guard
│   │       ├── +page.svelte           # Admin dashboard
│   │       ├── users/
│   │       │   └── +page.svelte
│   │       ├── calendar/
│   │       │   └── +page.svelte
│   │       └── settings/
│   │           └── +page.svelte
│   ├── app.css                        # Tailwind v4 with @theme
│   ├── app.html                       # HTML template
│   └── app.d.ts                       # Type declarations
├── static/                            # Static assets
├── tests/                             # Test files
├── svelte.config.js                   # Svelte configuration
├── vite.config.ts                     # Vite configuration
├── package.json                       # Dependencies
├── tsconfig.json                      # TypeScript config
└── README.md                          # Frontend documentation
```

---

## 7. Key Design Decisions

### 7.1 Technology Choices

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Database** | SQLite (modernc.org/sqlite) | CGo-free, single-file, sufficient for 5-person team, no infrastructure overhead |
| **Auth** | JWT (HS256, 24h expiry) | Stateless, industry standard, simple implementation |
| **Password** | bcrypt (cost 10) | Secure, battle-tested, configurable difficulty |
| **Frontend Framework** | Svelte 5 with SvelteKit | Minimal bundle, excellent DX, runes-based reactivity |
| **Component Library** | Melt UI (headless) | Accessible by default, full styling control |
| **Styling** | Tailwind CSS v4 | Utility-first, @theme configuration, Vite plugin |
| **Email** | Resend API | Simple API, reliable delivery, good DX |

### 7.2 Architectural Patterns

| Pattern | Application | Benefit |
|---------|-------------|---------|
| **Repository Pattern** | Data access layer | Decouples business logic from database |
| **Service Layer** | Business logic | Centralized validation and rules |
| **DTO Pattern** | Request/Response | Type-safe API contracts |
| **Middleware Chain** | Cross-cutting concerns | Clean separation of auth, CORS, errors |
| **Store Pattern** | Frontend state | Reactive, centralized state management |
| **Builder Pattern** | Melt UI components | Flexible, accessible component composition |

### 7.3 Security Decisions

| Area | Decision | Implementation |
|------|----------|----------------|
| **Passwords** | bcrypt with cost 10 | golang.org/x/crypto/bcrypt |
| **Tokens** | 24h expiry, HS256 | golang-jwt/jwt v5 |
| **CORS** | Explicit origin list | Gin middleware |
| **Role Enforcement** | Middleware-based | AdminMiddleware, EmployeeMiddleware |
| **Input Validation** | DTO binding tags | go-playground/validator |

---

## 8. Cross-Reference Guide

### 8.1 Related Documentation Files

| File | Content | When to Reference |
|------|---------|-------------------|
| `01-database-schema.md` | Tables, columns, migrations | Database changes |
| `02-api-specification.md` | Endpoints, request/response | API development |
| `03-implementation-roadmap.md` | Phases, dependencies | Planning, priorities |
| `04-backend-tasks.md` | Go implementation tasks | Backend development |
| `05-frontend-tasks.md` | Svelte implementation tasks | Frontend development |
| `06-component-inventory.md` | All components with specs | Component development |
| `07-testing-strategy.md` | Test plans, coverage | Testing |
| `08-security-checklist.md` | Security requirements | Security review |

### 8.2 Source Documentation

| File | Path | Purpose |
|------|------|---------|
| **Features Spec** | `docs/features.md` | Complete feature specifications |
| **Go Guide** | `docs/go.md` | Backend implementation patterns |
| **Svelte Guide** | `docs/svelte-melt.md` | Frontend implementation patterns |
| **Project Config** | `CLAUDE.md` | Development commands, conventions |

---

*Document Version: 1.0*
*Generated for AI Agent Implementation*
