# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

VacayTracker is an employee vacation tracking application with role-based access control (Admin/Employee), email notifications via Resend, and team calendar visualization. Uses a beach/vacation theme throughout the UI.

## Tech Stack

**Frontend (`vacaytracker-frontend/`):**
- Svelte 5 with Runes (`$state`, `$derived`, `$effect`, `$props`)
- SvelteKit 2.10.x for routing and SSR
- Bits UI v2.x for headless accessible components (compound component API, no preprocessor)
- Tailwind CSS v4 with CSS-native `@theme` configuration
- TypeScript 5.x
- Lucide Svelte for icons

**Backend (`vacaytracker-api/`):**
- Go 1.23+
- Gin web framework
- SQLite with modernc.org/sqlite (CGo-free, `CGO_ENABLED=0`)
- golang-jwt/jwt v5 for authentication
- Resend API for email (resend-go v2)

## Development Commands

**Frontend:**
```bash
cd vacaytracker-frontend
npm run dev                # Vite dev server at localhost:5173
npm run check              # svelte-check type checking
npm run check:watch        # Watch mode
npm run build              # Production build
npm run lint               # ESLint (flat config: eslint.config.js)
npm run format             # Prettier (.prettierrc: tabs, single quotes, width 100)
npm test                   # Vitest unit + component tests
```

**Backend:**
```bash
cd vacaytracker-api
make run                   # Run server (or: go run ./cmd/server/main.go)
make build                 # CGO_ENABLED=0 go build
make test                  # go test -v ./...
make test-coverage         # Test with coverage report
make lint                  # golangci-lint run
make migrate               # Run database migrations

# Run tests for a specific package
go test -v ./internal/service/...

# Run specific test function
go test -v -run TestFunctionName ./internal/service/...
```

**Docker (from project root):**
```bash
docker compose up --build          # Start full stack (API: localhost:32804, Frontend: localhost:32805)
docker compose down                # Stop containers
docker compose logs -f frontend    # Follow frontend logs
docker compose logs -f api         # Follow API logs
```

## Architecture

### Frontend

**Routing** uses SvelteKit file-based routing with a `(app)` group for authenticated pages:
- `/` — Login page
- `/(app)/dashboard` — Employee dashboard
- `/(app)/calendar` — Team calendar
- `/(app)/settings` — User settings
- `/admin/*` — Admin routes (users, balances, settings)

The `(app)/+layout.svelte` enforces authentication (redirects to `/` if not authenticated) and renders the shared shell (UnifiedHeader, DecorativeBackground, Footer).

**Path aliases** configured in `svelte.config.js`:
- `$lib` → `./src/lib`
- `$components` → `./src/lib/components`

**Bits UI** provides headless accessible components via a compound component API (e.g., `<Dialog.Root>`, `<Dialog.Content>`). No preprocessor needed. Components use data attributes for styling (`data-[state=open]`, `data-[highlighted]`, etc.) and `bind:open`/`bind:value` for state management. A `<Tooltip.Provider>` wraps the app in the root `+layout.svelte` (component tests wrap tooltip users in `ui/__tests__/TooltipHarness.svelte`).

**Dialogs** use the shared `ui/Modal.svelte` (Dialog) or `ui/ConfirmDialog.svelte` (AlertDialog) rather than raw Bits parts, which gives every dialog the labelled close button, enter/exit transitions and the close contract below:
- Bits calls `onOpenChange` **only for its own closes** (Escape, overlay, close button), never when code sets `open = false`. Don't put cleanup there.
- Initialise form state when the dialog **opens** (`$effect.pre(() => { if (open) untrack(reset) })`, or in the function that opens it).
- Put post-close cleanup in `onClosed`, which wraps Bits `onOpenChangeComplete` and fires after **every** close.

**Other Bits UI conventions:**
- `ui/Tooltip.svelte` takes a `child` snippet: spread its `props` onto a single element instead of nesting a button inside Bits' trigger button.
- To combine two triggers on one element (tooltip + popover), render `<Popover.Trigger {...tooltipProps}>` inside `<Tooltip.Trigger>`'s `child` snippet.
- Floating content inside a dialog (e.g. the date range calendar) must be portalled (`<Portal>`). The dialog's transform and `backdrop-filter` create a containing block that would clip it.
- Enter/exit transitions use the `presence-pop` (content) and `presence-fade` (overlays) utilities from `app.css`. Bits sets `data-starting-style`/`data-ending-style` and waits for the transition before unmounting.

**Navigation** uses `page` from `$app/state` (not the deprecated `$app/stores`) and wraps internal hrefs and `goto()` targets in `resolve()` from `$app/paths` (enforced by ESLint).

**API proxy** in `vite.config.ts`: `/api` requests proxy to `VITE_API_URL || 'http://localhost:3000'`.

**Stores** (`src/lib/stores/*.svelte.ts`) are module-level singletons using runes, not context-based providers:
```typescript
function createSomeStore() {
  let value = $state(initialValue);
  const computed = $derived(/* ... */);
  return { get value() { return value; }, get computed() { return computed; } };
}
export const someStore = createSomeStore();
```

**API client** (`src/lib/api/client.ts`) provides a generic `request<T>()` function with automatic JWT Bearer token injection. Errors throw `ApiException` with `code`, `message`, `status`, and optional `details`. Each domain has its own API module (auth.ts, vacation.ts, admin.ts, calendar.ts, settings.ts).

**Vacation store** uses a 5-minute cache TTL and distinguishes `usedDays` (started/past approved) from `upcomingDays` (future approved) for balance visualization.

### Backend

**Dependency injection** wiring in `cmd/server/main.go`:
```
DB → Repositories (user, vacation, settings)
  → Services (auth, vacation, user, email, newsletter)
    → Handlers (health, auth, vacation, admin, settings)
```

**Middleware chain**: Logger → Recovery → ErrorMiddleware → SecurityHeaders → SecurityLogging → RateLimiter → CORS → (per-group: AuthMiddleware, AdminMiddleware)

**Route groups**:
- `/health` — Public health check
- `/api/auth/login` — Public with stricter rate limiting
- `/api/auth/*`, `/api/vacation/*`, `/api/settings/*` — Authenticated (AuthMiddleware)
- `/api/admin/*` — Authenticated + admin role (AuthMiddleware + AdminMiddleware)

**Error handling**: Centralized `AppError` type in `internal/dto/errors.go` with HTTP status, error code constants, and structured JSON response. Handlers check for `AppError` to return appropriate status codes.

**Email sends** are non-blocking — handlers dispatch emails via goroutines.

**Newsletter scheduler**: Background goroutine (not cron), started/stopped with the server lifecycle.

**Migrations**: Single SQL file at `migrations/001_init.sql`, auto-run at server startup.

## Svelte 5 Patterns

Use runes for reactivity:
```svelte
<script lang="ts">
  let count = $state(0);                    // Reactive state
  let doubled = $derived(count * 2);         // Computed value
  let { name, onclick }: Props = $props();   // Component props
  let { open = $bindable(false) } = $props(); // Two-way bindable prop
  $effect(() => { /* side effects */ });
</script>
```

For runes in `.ts` files, use `.svelte.ts` extension:
```typescript
// auth.svelte.ts
function createAuthStore() {
  let user = $state<User | null>(null);
  const isAuthenticated = $derived(user !== null);

  return {
    get user() { return user; },
    get isAuthenticated() { return isAuthenticated; },
  };
}
export const auth = createAuthStore();
```

**Toast notifications**: `import { toast } from '$lib/stores/toast.svelte'` then `toast.success('Title')`, `toast.error('Title', 'Description')`, or `toast.add('info', 'Title', 'Desc', duration)`.

## Tailwind v4

Configure theme in CSS (no `tailwind.config.js`):
```css
@import "tailwindcss";

@theme {
  --color-ocean-500: oklch(0.55 0.18 220);
  --color-sand-200: oklch(0.9 0.04 80);
  --color-success: oklch(0.65 0.20 145);
}
```

Theme colors: `ocean-*` (primary blue), `sand-*` (neutral warm), semantic (`success`, `warning`, `error`, `pending`, `approved`, `rejected`).

Custom animations are defined in `app.css` (fadeIn, slideUp, scaleIn, wave, float, shimmer, etc.).

## User Roles

- **Admin (Captain)**: Full access, approve/reject requests, manage users
- **Employee (Crew)**: Submit requests, view balance and team calendar

## Date Handling

- **API input format**: DD/MM/YYYY (EU) — the backend `parseDDMMYYYY` parses this
- **HTML date inputs**: YYYY-MM-DD (ISO) — converted via `toEUFormat()` before API calls
- **Storage/display format**: YYYY-MM-DD (ISO) in database and frontend state
- Business days calculation respects the weekend policy (`excludeWeekends` + `excludedDays`). The frontend preview (`calculateBusinessDays`) mirrors the API and gets the policy from `GET /api/settings/public`, which employees can read (`/api/admin/settings` is admin-only).
- Parse `YYYY-MM-DD` strings with `parseISODate()` (local) or compare them as strings. `new Date('YYYY-MM-DD')` is UTC midnight and shifts the day in non-UTC timezones.
- Bits UI date pickers use `@internationalized/date` `DateValue` — converted via `dateValueToApiFormat()`

## Environment Variables

Required:
- `JWT_SECRET` — Token signing key (32+ chars, enforced)
- `ADMIN_PASSWORD` — Initial admin password

Email (optional):
- `RESEND_API_KEY`, `EMAIL_FROM_ADDRESS`, `EMAIL_FROM_NAME`

Server:
- `PORT` (default: 3000), `ENV`, `APP_URL`

Database:
- `DB_PATH` (default: ./data/vacaytracker.db)

## API Reference

- Backend runs on `http://localhost:3000` (dev) with `/api/` prefix
- Frontend proxies API requests in dev mode via Vite
- Auth: JWT Bearer token in `Authorization` header, stored in localStorage
- Auth middleware stores claims in Gin context (`ContextKeyUserID`, `ContextKeyEmail`, `ContextKeyRole`)
- See `aidocs/02-api-specification.md` for full endpoint documentation

## aidocs/ Reference

The `aidocs/` directory contains detailed implementation documentation:
- `00-architecture-overview.md` — System architecture diagrams
- `01-database-schema.md` — SQLite tables and migrations
- `02-api-specification.md` — Complete REST API documentation
- `03-implementation-roadmap.md` — Feature phases and dependencies
- `04-backend-tasks.md` — Backend implementation checklist
- `05-frontend-tasks.md` — Frontend implementation roadmap
- `06-component-inventory.md` — Component specs (entities, repos, services, handlers)
- `07-testing-strategy.md` — Testing pyramid and stack (testify, httptest, Vitest)
- `08-security-checklist.md` — Security guidelines
- `09-deployment-guide.md` — Production deployment
- `10-development-workflow.md` — Git workflow and branch strategy

## Important Guidelines

- If uncertain about Svelte implementation, read the latest Svelte 5 docs first
- If uncertain about Go implementation, research the latest Go docs first
- Always ensure our approach aligns with the documentation in `aidocs/`
- Bits UI uses compound component patterns — no preprocessor or build transforms needed
- After changes, ensure to rebuild the containers
