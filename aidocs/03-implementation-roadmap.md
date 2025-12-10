# 03 - Implementation Roadmap

> Phased implementation plan with dependencies, critical path, and priorities

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Phase Overview](#phase-overview)
3. [Dependency Graph](#dependency-graph)
4. [Critical Path](#critical-path)
5. [Phase Details](#phase-details)
6. [MVP Definition](#mvp-definition)
7. [Risk Assessment](#risk-assessment)

---

## Executive Summary

### Project Scope

| Metric | Value |
|--------|-------|
| Total Phases | 10 |
| Total Tasks | ~150 |
| MVP Phases | 1-6 |
| Backend Files | ~20 |
| Frontend Files | ~40 |
| Database Tables | 3 |
| API Endpoints | 20 |

### Priority Legend

| Priority | Description |
|----------|-------------|
| **P0 - Critical** | Must complete for MVP, blocks other work |
| **P1 - High** | Core functionality, MVP requirement |
| **P2 - Medium** | Important features, post-MVP acceptable |
| **P3 - Low** | Nice-to-have, can defer |

### Complexity Legend

| Complexity | Description |
|------------|-------------|
| **Low** | Straightforward, well-defined, 1-2 files |
| **Medium** | Moderate complexity, 3-5 files |
| **High** | Complex logic, many dependencies, 5+ files |

---

## Phase Overview

```
Phase 1: Backend Foundation     ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Phase 2: Authentication         ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Phase 3: Vacation Core          ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Phase 4: Frontend Foundation    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Phase 5: Employee Dashboard     ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Phase 6: Admin Dashboard        ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Phase 7: Calendar System        ━━━━━━━━━━━━━━━━━━━━━━━━━━
Phase 8: Email Notifications    ━━━━━━━━━━━━━━━━━━━━
Phase 9: Newsletter System      ━━━━━━━━━━━━━━━━
Phase 10: Polish & Settings     ━━━━━━━━━━━━━━
                               MVP ────────────────────┘
```

| Phase | Name | Priority | Complexity | Tasks | Dependencies |
|-------|------|----------|------------|-------|--------------|
| 1 | Backend Foundation | P0 | Medium | 15 | None |
| 2 | Authentication | P0 | Medium | 12 | Phase 1 |
| 3 | Vacation Core | P1 | High | 18 | Phase 2 |
| 4 | Frontend Foundation | P1 | Medium | 20 | Phase 2 |
| 5 | Employee Dashboard | P1 | High | 15 | Phase 3, 4 |
| 6 | Admin Dashboard | P1 | High | 20 | Phase 3, 4 |
| 7 | Calendar System | P2 | Medium | 12 | Phase 3, 4 |
| 8 | Email Notifications | P2 | Low | 10 | Phase 2, 3 |
| 9 | Newsletter System | P3 | Medium | 10 | Phase 8 |
| 10 | Polish & Settings | P3 | Low | 8 | Phase 6 |

---

## Dependency Graph

```
                    ┌─────────────────────────────────────────────────────────┐
                    │                                                         │
                    v                                                         │
┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐       │
│   Phase 1   │──>│   Phase 2   │──>│   Phase 3   │──>│   Phase 8   │───────┤
│  Foundation │   │    Auth     │   │   Vacation  │   │    Email    │       │
└─────────────┘   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘       │
                         │                 │                 │               │
                         │                 │                 v               │
                         │                 │          ┌─────────────┐       │
                         │                 │          │   Phase 9   │       │
                         │                 │          │  Newsletter │       │
                         │                 │          └─────────────┘       │
                         v                 │                                 │
                  ┌─────────────┐          │                                 │
                  │   Phase 4   │<─────────┤                                 │
                  │  Frontend   │          │                                 │
                  └──────┬──────┘          │                                 │
                         │                 │                                 │
         ┌───────────────┼─────────────────┤                                 │
         │               │                 │                                 │
         v               v                 v                                 │
  ┌─────────────┐ ┌─────────────┐  ┌─────────────┐                          │
  │   Phase 5   │ │   Phase 6   │  │   Phase 7   │                          │
  │  Employee   │ │    Admin    │  │  Calendar   │                          │
  │  Dashboard  │ │  Dashboard  │  └─────────────┘                          │
  └─────────────┘ └──────┬──────┘                                            │
                         │                                                   │
                         v                                                   │
                  ┌─────────────┐                                            │
                  │  Phase 10   │────────────────────────────────────────────┘
                  │   Polish    │
                  └─────────────┘
```

### Dependency Matrix

| Phase | Depends On | Blocks |
|-------|------------|--------|
| 1 - Foundation | - | 2, 4 |
| 2 - Auth | 1 | 3, 4, 8 |
| 3 - Vacation | 2 | 5, 6, 7, 8 |
| 4 - Frontend | 2 | 5, 6, 7 |
| 5 - Employee | 3, 4 | - |
| 6 - Admin | 3, 4 | 10 |
| 7 - Calendar | 3, 4 | - |
| 8 - Email | 2, 3 | 9 |
| 9 - Newsletter | 8 | - |
| 10 - Polish | 6 | - |

---

## Critical Path

The critical path represents the longest sequence of dependent tasks that determines minimum project duration.

### Critical Path Sequence

```
Phase 1 → Phase 2 → Phase 3 → Phase 4 → Phase 6 → Phase 10
   │          │          │          │         │          │
   v          v          v          v         v          v
Foundation  Auth    Vacation   Frontend   Admin     Polish
```

### Critical Path Tasks (Must not slip)

1. **Phase 1:** Database connection, user repository, config system
2. **Phase 2:** JWT service, auth middleware, login handler
3. **Phase 3:** Vacation repository, business day calculation, request handlers
4. **Phase 4:** SvelteKit setup, API client, auth store, base UI components
5. **Phase 6:** Admin route guards, pending requests view, approval workflow
6. **Phase 10:** Settings UI, final integrations

### Parallel Work Opportunities

| After Phase | Can Run In Parallel |
|-------------|---------------------|
| Phase 2 | Phase 3 + Phase 4 (different engineers) |
| Phase 4 | Phase 5 + Phase 6 + Phase 7 (if backend ready) |
| Phase 3 | Phase 8 (email is independent after vacation exists) |

---

## Phase Details

### Phase 1: Backend Foundation

**Priority:** P0 - Critical | **Complexity:** Medium | **Blocks:** Everything

#### Objectives
- [ ] Establish Go project structure
- [ ] Configure development environment
- [ ] Set up SQLite database with migrations
- [ ] Implement core domain models
- [ ] Create repository layer foundation

#### Tasks

| # | Task | File(s) | Complexity | Acceptance Criteria |
|---|------|---------|------------|---------------------|
| 1.1 | Initialize Go module | `go.mod`, `go.sum` | Low | Module compiles, deps installed |
| 1.2 | Create project structure | All directories | Low | Standard layout created |
| 1.3 | Implement config loader | `internal/config/config.go` | Low | Loads from .env and env vars |
| 1.4 | Create main.go entry point | `cmd/server/main.go` | Low | Server starts on PORT |
| 1.5 | Set up Gin router | `cmd/server/main.go` | Low | Health endpoint responds |
| 1.6 | Implement health handler | `internal/handler/health.go` | Low | GET /health returns 200 |
| 1.7 | Create SQLite connection | `internal/repository/sqlite/sqlite.go` | Medium | WAL mode, foreign keys enabled |
| 1.8 | Write migration SQL | `migrations/001_init.sql` | Medium | All 3 tables created |
| 1.9 | Implement migration runner | `internal/repository/sqlite/migrate.go` | Medium | Migrations run on startup |
| 1.10 | Define User domain | `internal/domain/user.go` | Low | Struct with all fields |
| 1.11 | Define VacationRequest domain | `internal/domain/vacation.go` | Low | Struct with all fields |
| 1.12 | Define Settings domain | `internal/domain/settings.go` | Low | Struct with all fields |
| 1.13 | Create error types | `internal/dto/errors.go` | Low | AppError with codes |
| 1.14 | Create Makefile | `Makefile` | Low | build, run, test targets |
| 1.15 | Create .env.example | `.env.example` | Low | All env vars documented |

#### Deliverables
- Working Go server with health endpoint
- SQLite database with schema
- Domain models defined
- Configuration system functional

---

### Phase 2: Authentication

**Priority:** P0 - Critical | **Complexity:** Medium | **Depends On:** Phase 1

#### Objectives
- [ ] Implement secure password hashing
- [ ] Create JWT token generation and validation
- [ ] Build authentication middleware
- [ ] Implement login/logout handlers
- [ ] Create user CRUD operations

#### Tasks

| # | Task | File(s) | Complexity | Acceptance Criteria |
|---|------|---------|------------|---------------------|
| 2.1 | Implement password hashing | `internal/service/auth.go` | Medium | bcrypt cost 10, hash/verify |
| 2.2 | Implement JWT generation | `internal/service/auth.go` | Medium | HS256, 24h expiry, claims |
| 2.3 | Implement JWT validation | `internal/service/auth.go` | Medium | Validates signature, expiry |
| 2.4 | Create auth middleware | `internal/middleware/auth.go` | Medium | Extracts user from token |
| 2.5 | Create admin middleware | `internal/middleware/auth.go` | Low | Checks role == admin |
| 2.6 | Implement user repository | `internal/repository/sqlite/user.go` | Medium | CRUD operations |
| 2.7 | Implement login handler | `internal/handler/auth.go` | Medium | Returns JWT + user |
| 2.8 | Implement /auth/me handler | `internal/handler/auth.go` | Low | Returns current user |
| 2.9 | Implement password change | `internal/handler/auth.go` | Medium | Verifies current, updates |
| 2.10 | Implement email prefs update | `internal/handler/auth.go` | Low | Updates JSON field |
| 2.11 | Create request DTOs | `internal/dto/request.go` | Low | Login, password change |
| 2.12 | Create response DTOs | `internal/dto/response.go` | Low | User, token responses |

#### Deliverables
- POST /api/auth/login working
- GET /api/auth/me working
- PUT /api/auth/password working
- PUT /api/auth/email-preferences working
- Protected routes require valid JWT

---

### Phase 3: Vacation Core

**Priority:** P1 - High | **Complexity:** High | **Depends On:** Phase 2

#### Objectives
- [ ] Implement vacation request CRUD
- [ ] Build business day calculation
- [ ] Create approval workflow
- [ ] Implement balance management
- [ ] Build team calendar query

#### Tasks

| # | Task | File(s) | Complexity | Acceptance Criteria |
|---|------|---------|------------|---------------------|
| 3.1 | Implement vacation repository | `internal/repository/sqlite/vacation.go` | High | CRUD + queries |
| 3.2 | Implement business day calc | `internal/service/vacation.go` | High | Excludes weekends if enabled |
| 3.3 | Implement balance checking | `internal/service/vacation.go` | Medium | Returns error if insufficient |
| 3.4 | Create vacation request | `internal/handler/vacation.go` | Medium | POST /api/vacation/request |
| 3.5 | List user requests | `internal/handler/vacation.go` | Low | GET /api/vacation/requests |
| 3.6 | Get single request | `internal/handler/vacation.go` | Low | GET /api/vacation/requests/:id |
| 3.7 | Cancel request | `internal/handler/vacation.go` | Medium | DELETE, only pending |
| 3.8 | Team calendar endpoint | `internal/handler/vacation.go` | Medium | GET /api/vacation/team |
| 3.9 | Settings repository | `internal/repository/sqlite/settings.go` | Medium | CRUD for settings |
| 3.10 | Get settings endpoint | `internal/handler/admin.go` | Low | GET /api/admin/settings |
| 3.11 | Update settings endpoint | `internal/handler/admin.go` | Low | PUT /api/admin/settings |
| 3.12 | Pending requests endpoint | `internal/handler/admin.go` | Medium | GET /api/admin/vacation/pending |
| 3.13 | Approve request | `internal/handler/admin.go` | High | PUT, deduct balance |
| 3.14 | Reject request | `internal/handler/admin.go` | Medium | PUT, optional reason |
| 3.15 | Date parsing (DD/MM/YYYY) | `internal/service/vacation.go` | Low | EU format to ISO |
| 3.16 | Vacation request DTOs | `internal/dto/request.go` | Low | Create, approve, reject |
| 3.17 | Vacation response DTOs | `internal/dto/response.go` | Low | Request, list, team |
| 3.18 | Vacation service tests | `internal/service/vacation_test.go` | High | Business logic coverage |

#### Deliverables
- Full vacation request lifecycle
- Business day calculation respects settings
- Balance automatically managed
- Admin can approve/reject
- Team calendar shows approved vacations

---

### Phase 4: Frontend Foundation

**Priority:** P1 - High | **Complexity:** Medium | **Depends On:** Phase 2

#### Objectives
- [ ] Initialize SvelteKit project
- [ ] Configure Tailwind v4
- [ ] Build base UI components
- [ ] Create auth store and API client
- [ ] Implement login page

#### Tasks

| # | Task | File(s) | Complexity | Acceptance Criteria |
|---|------|---------|------------|---------------------|
| 4.1 | Initialize SvelteKit | Project root | Low | npm create svelte |
| 4.2 | Configure TypeScript | `tsconfig.json` | Low | Strict mode enabled |
| 4.3 | Install Tailwind v4 | `package.json`, `app.css` | Medium | @theme working |
| 4.4 | Define color theme | `src/app.css` | Low | Ocean/beach palette |
| 4.5 | Create Button component | `src/lib/components/ui/Button.svelte` | Medium | Variants, loading state |
| 4.6 | Create Input component | `src/lib/components/ui/Input.svelte` | Medium | Error state, label |
| 4.7 | Create Card component | `src/lib/components/ui/Card.svelte` | Low | Header, body, footer |
| 4.8 | Create Badge component | `src/lib/components/ui/Badge.svelte` | Low | Status colors |
| 4.9 | Create Avatar component | `src/lib/components/ui/Avatar.svelte` | Low | Initials fallback |
| 4.10 | Create ProgressRing | `src/lib/components/ui/ProgressRing.svelte` | Medium | SVG ring |
| 4.11 | Create auth store | `src/lib/stores/auth.svelte.ts` | High | Login, logout, persist |
| 4.12 | Create toast store | `src/lib/stores/toast.svelte.ts` | Medium | Show/dismiss toasts |
| 4.13 | Create API client | `src/lib/api/client.ts` | High | Fetch wrapper, error handling |
| 4.14 | Create auth API | `src/lib/api/auth.ts` | Medium | Login, me, password |
| 4.15 | Create vacation API | `src/lib/api/vacation.ts` | Medium | CRUD operations |
| 4.16 | Create TypeScript types | `src/lib/types/index.ts` | Medium | All DTOs as types |
| 4.17 | Create root layout | `src/routes/+layout.svelte` | Medium | Toast container |
| 4.18 | Create login page | `src/routes/+page.svelte` | High | Form, validation, redirect |
| 4.19 | Create employee layout guard | `src/routes/employee/+layout.server.ts` | Medium | Redirect if not auth |
| 4.20 | Create admin layout guard | `src/routes/admin/+layout.server.ts` | Medium | Redirect if not admin |

#### Deliverables
- SvelteKit project running
- Tailwind v4 with beach theme
- Base UI component library
- Login page functional
- Route guards in place

---

### Phase 5: Employee Dashboard

**Priority:** P1 - High | **Complexity:** High | **Depends On:** Phase 3, Phase 4

#### Objectives
- [ ] Build vacation balance display
- [ ] Create request submission modal
- [ ] Show request history
- [ ] Display team calendar overview
- [ ] Implement profile settings

#### Tasks

| # | Task | File(s) | Complexity | Acceptance Criteria |
|---|------|---------|------------|---------------------|
| 5.1 | Create employee layout | `src/routes/employee/+layout.svelte` | Medium | Header, nav |
| 5.2 | Create dashboard page | `src/routes/employee/+page.svelte` | High | Balance + requests |
| 5.3 | Create BalanceDisplay | `src/lib/components/features/vacation/BalanceDisplay.svelte` | Medium | Progress ring, days |
| 5.4 | Create RequestModal | `src/lib/components/features/vacation/RequestModal.svelte` | High | Melt Dialog, form |
| 5.5 | Create DateRangePicker | `src/lib/components/features/vacation/DateRangePicker.svelte` | High | DD/MM/YYYY, validation |
| 5.6 | Create RequestList | `src/lib/components/features/vacation/RequestList.svelte` | Medium | Status badges |
| 5.7 | Create RequestCard | `src/lib/components/features/vacation/RequestCard.svelte` | Medium | Details, cancel |
| 5.8 | Create team overview | `src/routes/employee/team/+page.svelte` | Medium | Calendar summary |
| 5.9 | Create profile page | `src/routes/employee/settings/+page.svelte` | Medium | Email prefs |
| 5.10 | Create password change | `src/lib/components/features/auth/PasswordChange.svelte` | Medium | Form, validation |
| 5.11 | Create EmailPreferences | `src/lib/components/features/auth/EmailPreferences.svelte` | Low | Toggle switches |
| 5.12 | Employee header | `src/lib/components/layout/EmployeeHeader.svelte` | Medium | Logo, nav, user menu |
| 5.13 | User dropdown menu | `src/lib/components/layout/UserMenu.svelte` | Medium | Melt Dropdown |
| 5.14 | Mobile navigation | `src/lib/components/layout/MobileNav.svelte` | Medium | Hamburger, drawer |
| 5.15 | Vacation store | `src/lib/stores/vacation.svelte.ts` | Medium | Requests cache |

#### Deliverables
- Employee dashboard with balance
- Can submit vacation requests
- Can view request history
- Can cancel pending requests
- Can update profile settings

---

### Phase 6: Admin Dashboard

**Priority:** P1 - High | **Complexity:** High | **Depends On:** Phase 3, Phase 4

#### Objectives
- [ ] Build admin statistics overview
- [ ] Create pending request management
- [ ] Implement user management
- [ ] Build vacation balance editing
- [ ] Create settings management

#### Tasks

| # | Task | File(s) | Complexity | Acceptance Criteria |
|---|------|---------|------------|---------------------|
| 6.1 | Create admin layout | `src/routes/admin/+layout.svelte` | Medium | Sidebar nav |
| 6.2 | Create admin dashboard | `src/routes/admin/+page.svelte` | High | Stats + pending |
| 6.3 | Create StatsCard | `src/lib/components/features/admin/StatsCard.svelte` | Low | Icon, number, label |
| 6.4 | Create PendingRequests | `src/lib/components/features/admin/PendingRequests.svelte` | High | List + actions |
| 6.5 | Create ApprovalModal | `src/lib/components/features/admin/ApprovalModal.svelte` | Medium | Confirm approve |
| 6.6 | Create RejectionModal | `src/lib/components/features/admin/RejectionModal.svelte` | Medium | Reason input |
| 6.7 | Create users page | `src/routes/admin/users/+page.svelte` | High | Table + CRUD |
| 6.8 | Create UserTable | `src/lib/components/features/admin/UserTable.svelte` | High | Sort, filter |
| 6.9 | Create UserModal | `src/lib/components/features/admin/UserModal.svelte` | High | Create/edit form |
| 6.10 | Create DeleteConfirm | `src/lib/components/features/admin/DeleteConfirm.svelte` | Medium | Warning dialog |
| 6.11 | Create balance page | `src/routes/admin/balances/+page.svelte` | Medium | Balance editor |
| 6.12 | Create BalanceEditor | `src/lib/components/features/admin/BalanceEditor.svelte` | Medium | Inline edit |
| 6.13 | Create settings page | `src/routes/admin/settings/+page.svelte` | Medium | All settings |
| 6.14 | Create WeekendPolicy | `src/lib/components/features/admin/WeekendPolicy.svelte` | Medium | Day toggles |
| 6.15 | Create NewsletterSettings | `src/lib/components/features/admin/NewsletterSettings.svelte` | Medium | Frequency, day |
| 6.16 | Admin API module | `src/lib/api/admin.ts` | Medium | All admin endpoints |
| 6.17 | Admin sidebar | `src/lib/components/layout/AdminSidebar.svelte` | Medium | Nav links |
| 6.18 | Admin header | `src/lib/components/layout/AdminHeader.svelte` | Low | Title, user |
| 6.19 | Pagination component | `src/lib/components/ui/Pagination.svelte` | Medium | Page controls |
| 6.20 | Admin store | `src/lib/stores/admin.svelte.ts` | Medium | Users, settings |

#### Deliverables
- Admin dashboard with statistics
- Can approve/reject requests
- Can manage users (CRUD)
- Can edit vacation balances
- Can configure settings

---

### Phase 7: Calendar System

**Priority:** P2 - Medium | **Complexity:** Medium | **Depends On:** Phase 3, Phase 4

#### Objectives
- [ ] Build week and month views
- [ ] Render vacation events
- [ ] Add filtering capabilities
- [ ] Implement caching

#### Tasks

| # | Task | File(s) | Complexity | Acceptance Criteria |
|---|------|---------|------------|---------------------|
| 7.1 | Create calendar page | `src/routes/employee/calendar/+page.svelte` | Medium | View selector |
| 7.2 | Create WeekView | `src/lib/components/features/calendar/WeekView.svelte` | High | 7 day grid |
| 7.3 | Create MonthView | `src/lib/components/features/calendar/MonthView.svelte` | High | Month grid |
| 7.4 | Create CalendarHeader | `src/lib/components/features/calendar/CalendarHeader.svelte` | Medium | Nav, view switch |
| 7.5 | Create VacationEvent | `src/lib/components/features/calendar/VacationEvent.svelte` | Medium | Event display |
| 7.6 | Create DayCell | `src/lib/components/features/calendar/DayCell.svelte` | Medium | Day + events |
| 7.7 | Create FilterPanel | `src/lib/components/features/calendar/FilterPanel.svelte` | Medium | User filter |
| 7.8 | Calendar store | `src/lib/stores/calendar.svelte.ts` | High | Data + navigation |
| 7.9 | Date utility functions | `src/lib/utils/date.ts` | Medium | Week/month helpers |
| 7.10 | Calendar API | `src/lib/api/calendar.ts` | Low | Team vacation fetch |
| 7.11 | Event coloring | `src/lib/utils/colors.ts` | Low | User color assignment |
| 7.12 | Calendar caching | `src/lib/stores/calendar.svelte.ts` | Medium | Cache by month |

#### Deliverables
- Week view calendar
- Month view calendar
- Vacation events displayed
- Filter by user
- Smooth navigation

---

### Phase 8: Email Notifications

**Priority:** P2 - Medium | **Complexity:** Low | **Depends On:** Phase 2, Phase 3

#### Objectives
- [ ] Integrate Resend API
- [ ] Send welcome emails
- [ ] Send status update emails
- [ ] Build email templates

#### Tasks

| # | Task | File(s) | Complexity | Acceptance Criteria |
|---|------|---------|------------|---------------------|
| 8.1 | Create email service | `internal/service/email.go` | Medium | Resend integration |
| 8.2 | Welcome email template | `internal/service/email.go` | Low | HTML + plain |
| 8.3 | Request submitted email | `internal/service/email.go` | Low | To user |
| 8.4 | Request approved email | `internal/service/email.go` | Low | To user |
| 8.5 | Request rejected email | `internal/service/email.go` | Low | To user + reason |
| 8.6 | Admin notification email | `internal/service/email.go` | Low | New request alert |
| 8.7 | Email config | `internal/config/config.go` | Low | Resend API key |
| 8.8 | Trigger emails from handlers | `internal/handler/*.go` | Medium | Async send |
| 8.9 | Email preference check | `internal/service/email.go` | Low | Respect user prefs |
| 8.10 | Email error handling | `internal/service/email.go` | Medium | Retry, logging |

#### Deliverables
- Welcome email on user creation
- Status emails on request changes
- Admin alerts for new requests
- Preferences respected

---

### Phase 9: Newsletter System

**Priority:** P3 - Low | **Complexity:** Medium | **Depends On:** Phase 8

#### Objectives
- [ ] Build monthly summary generator
- [ ] Implement scheduler
- [ ] Create admin controls
- [ ] Add preview functionality

#### Tasks

| # | Task | File(s) | Complexity | Acceptance Criteria |
|---|------|---------|------------|---------------------|
| 9.1 | Newsletter template | `internal/service/newsletter.go` | Medium | Monthly summary |
| 9.2 | Summary data aggregation | `internal/service/newsletter.go` | Medium | Stats, upcoming |
| 9.3 | Scheduler setup | `internal/service/scheduler.go` | Medium | Cron-like timing |
| 9.4 | Manual send endpoint | `internal/handler/admin.go` | Low | POST trigger |
| 9.5 | Preview mode | `internal/handler/admin.go` | Medium | Admin-only send |
| 9.6 | Newsletter settings | `internal/repository/sqlite/settings.go` | Low | Enable, frequency |
| 9.7 | Frontend send button | Admin settings page | Low | Trigger + preview |
| 9.8 | Last sent tracking | `internal/domain/settings.go` | Low | Prevent duplicates |
| 9.9 | Recipient list | `internal/service/newsletter.go` | Low | Active users |
| 9.10 | Newsletter tests | `internal/service/newsletter_test.go` | Medium | Template + schedule |

#### Deliverables
- Monthly newsletter generation
- Automated scheduling
- Manual send with preview
- Settings configuration

---

### Phase 10: Polish & Settings

**Priority:** P3 - Low | **Complexity:** Low | **Depends On:** Phase 6

#### Objectives
- [ ] Complete settings UI
- [ ] Add year-end reset functionality
- [ ] Implement UI polish
- [ ] Final integrations

#### Tasks

| # | Task | File(s) | Complexity | Acceptance Criteria |
|---|------|---------|------------|---------------------|
| 10.1 | Year-end balance reset | `internal/service/vacation.go` | Medium | Batch update |
| 10.2 | Reset settings UI | Admin settings | Low | Month selector |
| 10.3 | Default days setting | Admin settings | Low | Number input |
| 10.4 | Loading states | All components | Low | Skeleton/spinner |
| 10.5 | Error boundaries | `src/routes/+error.svelte` | Medium | Graceful errors |
| 10.6 | Empty states | Feature components | Low | Helpful messages |
| 10.7 | Animations | `src/app.css` | Low | Transitions |
| 10.8 | Final testing | All | Medium | E2E scenarios |

#### Deliverables
- All settings functional
- Year-end reset works
- Polish throughout UI
- Production ready

---

## MVP Definition

### MVP Includes (Phases 1-6)

| Feature | Phase | Endpoints | Components |
|---------|-------|-----------|------------|
| User authentication | 2 | 4 | Login page, auth store |
| Vacation requests | 3, 5 | 5 | Request modal, list, balance |
| Request approval | 3, 6 | 3 | Pending list, approve/reject |
| User management | 2, 6 | 5 | User table, CRUD modals |
| Basic settings | 3, 6 | 2 | Weekend policy |

### MVP Excludes (Phases 7-10)

| Feature | Phase | Reason |
|---------|-------|--------|
| Calendar views | 7 | Nice-to-have, team overview sufficient |
| Email notifications | 8 | Can operate without |
| Newsletter | 9 | Low priority |
| Advanced polish | 10 | Post-MVP refinement |

### MVP Acceptance Criteria

- [ ] Users can log in and log out
- [ ] Employees can submit vacation requests
- [ ] Employees can view their balance and history
- [ ] Employees can cancel pending requests
- [ ] Admins can approve or reject requests
- [ ] Admins can create and manage users
- [ ] Admins can edit vacation balances
- [ ] Weekend exclusion setting works
- [ ] All actions require appropriate authorization

---

## Risk Assessment

### Technical Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| SQLite performance at scale | Medium | Low | Index optimization, WAL mode |
| JWT token compromise | High | Low | Short expiry, secure secret |
| Business day calculation edge cases | Medium | Medium | Comprehensive tests |
| Svelte 5 runes learning curve | Low | Medium | Follow docs strictly |
| Melt UI breaking changes | Medium | Low | Pin version, check changelog |

### Schedule Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Auth complexity underestimated | High | Low | Start with simple JWT |
| Calendar implementation complex | Medium | Medium | Defer to post-MVP if needed |
| Integration issues | Medium | Medium | Test API early with frontend |

### Dependency Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Resend API changes | Low | Low | Abstract email service |
| Go module compatibility | Low | Low | Pin versions |
| Tailwind v4 instability | Medium | Low | Use stable features only |

---

## Related Documents

- [04-backend-tasks.md](./04-backend-tasks.md) - Detailed backend task list
- [05-frontend-tasks.md](./05-frontend-tasks.md) - Detailed frontend task list
- [06-component-inventory.md](./06-component-inventory.md) - Component specifications
- [07-testing-strategy.md](./07-testing-strategy.md) - Testing requirements per phase
