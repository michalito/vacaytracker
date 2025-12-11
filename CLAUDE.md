# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

VacayTracker is an employee vacation tracking application with role-based access control (Admin/Employee), email notifications via Resend, and team calendar visualization. Uses a beach/vacation theme throughout the UI.

## Tech Stack

**Frontend (`vacaytracker-frontend/`):**
- Svelte 5 with Runes (`$state`, `$derived`, `$effect`, `$props`)
- SvelteKit 2.10.x for routing and SSR
- @melt-ui/svelte v0.86.x available for headless accessible components
- Tailwind CSS v4 with CSS-native `@theme` configuration
- TypeScript 5.x
- Lucide Svelte for icons

**Backend (`vacaytracker-api/`):**
- Go 1.23+
- Gin web framework
- SQLite with modernc.org/sqlite (CGo-free)
- golang-jwt/jwt v5 for authentication
- Resend API for email

## Development Commands

**Frontend:**
```bash
cd vacaytracker-frontend
npm run dev                # Vite dev server at localhost:5173
npm run check              # svelte-check type checking
npm run check:watch        # Watch mode
npm run build              # Production build
npm run lint               # ESLint
npm run format             # Prettier
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
docker compose up --build          # Start full stack
docker compose down                # Stop containers
docker compose logs -f frontend    # Follow frontend logs
docker compose logs -f api         # Follow API logs
```

## Architecture

### Frontend Structure

```
vacaytracker-frontend/src/
├── routes/                 # SvelteKit file-based routing
│   ├── +page.svelte       # Login page (/)
│   ├── employee/          # Employee routes
│   │   ├── +layout.svelte
│   │   ├── +page.svelte   # Dashboard
│   │   ├── team/          # Team calendar
│   │   └── settings/      # User settings
│   └── admin/             # Admin routes
│       ├── +layout.svelte
│       ├── +page.svelte   # Admin dashboard
│       ├── users/         # User management
│       ├── balances/      # Balance management
│       └── settings/      # Admin settings
├── lib/
│   ├── components/
│   │   ├── ui/            # Base primitives (Button, Input, Card, Badge, etc.)
│   │   ├── layout/        # Headers, sidebars
│   │   └── features/      # Domain components (vacation/, admin/, auth/)
│   ├── stores/            # Global state (*.svelte.ts for runes)
│   ├── api/               # API client modules
│   └── types/             # TypeScript types
└── app.css                # Tailwind v4 with @theme config
```

### Backend Structure

```
vacaytracker-api/
├── cmd/server/main.go     # Entry point
├── internal/
│   ├── config/            # Configuration loading
│   ├── domain/            # Entities (user.go, vacation.go, settings.go)
│   ├── repository/sqlite/ # Database layer
│   ├── service/           # Business logic
│   ├── handler/           # HTTP handlers
│   ├── middleware/        # Auth, CORS, error handling
│   └── dto/               # Request/response types
└── Makefile               # Build commands
```

## Svelte 5 Patterns

Use runes for reactivity:
```svelte
<script lang="ts">
  let count = $state(0);                    // Reactive state
  let doubled = $derived(count * 2);         // Computed value
  let { name, onclick }: Props = $props();   // Component props
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

## User Roles

- **Admin (Captain)**: Full access, approve/reject requests, manage users
- **Employee (Crew)**: Submit requests, view balance and team calendar

## Date Handling

- Input format: DD/MM/YYYY (EU) for API
- HTML date inputs: YYYY-MM-DD (ISO)
- Storage format: YYYY-MM-DD (ISO)
- Business days calculation respects `excludeWeekends` setting

## Environment Variables

Required:
- `JWT_SECRET` - Token signing key (32+ chars)
- `ADMIN_PASSWORD` - Initial admin password

Email (optional):
- `RESEND_API_KEY`, `EMAIL_FROM_ADDRESS`, `EMAIL_FROM_NAME`

Server:
- `PORT` (default: 3000), `ENV`, `APP_URL`

Database:
- `DB_PATH` (default: ./data/vacaytracker.db)

## Important Guidelines

- If uncertain about Svelte implementation, read the latest Svelte 5 docs first
- If uncertain about Go implementation, research the latest Go docs first
- Always ensure our approach aligns with the documentation in `aidocs/`
- After changes, ensure to rebuild the containers