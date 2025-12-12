# VacayTracker

Employee vacation tracking application with role-based access control, email notifications, and team calendar visualization.

## Features

- **Role-based access**: Admin (Captain) and Employee (Crew) roles
- **Vacation requests**: Submit, approve, reject, and track vacation requests
- **Team calendar**: View team vacation schedules
- **Email notifications**: Automated emails via Resend for request updates
- **Newsletter**: Monthly summary emails with team stats
- **Modern UI**: Beach/vacation themed interface built with Svelte 5

## Tech Stack

| Component | Technology |
|-----------|------------|
| Frontend | Svelte 5, SvelteKit, Tailwind CSS v4 |
| Backend | Go 1.23+, Gin framework |
| Database | SQLite (embedded) |
| Email | Resend API |

## Quick Start with Docker

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) 24+
- [Docker Compose](https://docs.docker.com/compose/install/) 2.20+

### 1. Clone the Repository

```bash
git clone <repository-url>
cd vacaytracker
```

### 2. Configure Environment

```bash
# Copy the example environment file
cp .env.example .env

# Edit with your settings
nano .env  # or use your preferred editor
```

**Required settings:**

```env
# Generate a secure JWT secret (at least 32 characters)
JWT_SECRET=your-secure-secret-key-minimum-32-characters-long

# Set your admin password
ADMIN_PASSWORD=your-secure-password
```

**Optional email settings (for notifications):**

```env
# Get your API key from https://resend.com
RESEND_API_KEY=re_your_api_key_here

# Use a verified domain email OR onboarding@resend.dev for testing
EMAIL_FROM_ADDRESS=noreply@yourdomain.com
```

### 3. Start the Application

```bash
docker compose up --build
```

### 4. Access the Application

| Service | URL |
|---------|-----|
| Frontend | http://localhost:32805 |
| API | http://localhost:32804 |
| Health Check | http://localhost:32804/health |

### 5. Login

Use your configured admin credentials:
- **Email**: Value of `ADMIN_EMAIL` (default: `admin@company.com`)
- **Password**: Value of `ADMIN_PASSWORD`

## Docker Commands

```bash
# Start in foreground (see logs)
docker compose up --build

# Start in background
docker compose up --build -d

# View logs
docker compose logs -f           # All services
docker compose logs -f api       # API only
docker compose logs -f frontend  # Frontend only

# Stop services
docker compose down

# Stop and remove volumes (WARNING: deletes database!)
docker compose down -v

# Rebuild after code changes
docker compose up --build

# Check service status
docker compose ps
```

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `JWT_SECRET` | Yes | - | JWT signing secret (32+ characters) |
| `ADMIN_PASSWORD` | Yes | - | Initial admin password |
| `ADMIN_EMAIL` | No | `admin@company.com` | Admin email address |
| `ADMIN_NAME` | No | `Captain Admin` | Admin display name |
| `RESEND_API_KEY` | No | - | Resend API key for emails |
| `EMAIL_FROM_ADDRESS` | No | - | Sender email (verified in Resend) |
| `EMAIL_FROM_NAME` | No | `VacayTracker` | Sender display name |

### Generating Secure Secrets

```bash
# Generate JWT secret
openssl rand -base64 32

# Or using Node.js
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
```

## Email Setup (Resend)

VacayTracker sends the following email notifications:

| Template | When Sent |
|----------|-----------|
| Welcome | New user created with login credentials |
| Request Submitted | Employee submits vacation request |
| Request Approved | Admin approves a request |
| Request Rejected | Admin rejects a request |
| Admin Notification | New request pending review |
| Newsletter | Monthly team summary |

### Setting Up Resend

1. Create a free account at [resend.com](https://resend.com)
2. Get your API key from the dashboard
3. Either:
   - **For testing**: Use `onboarding@resend.dev` as `EMAIL_FROM_ADDRESS`
   - **For production**: Verify your domain and use your domain email

```env
RESEND_API_KEY=re_abc123...
EMAIL_FROM_ADDRESS=onboarding@resend.dev  # For testing
# EMAIL_FROM_ADDRESS=noreply@yourdomain.com  # For production
EMAIL_FROM_NAME=VacayTracker
```

### Testing Emails

1. Log in as admin
2. Go to Admin Settings
3. Use the Email Templates panel to preview and send test emails

## Project Structure

```
vacaytracker/
├── docker-compose.yml      # Docker orchestration
├── .env.example            # Environment template
├── vacaytracker-api/       # Go backend
│   ├── cmd/server/         # Entry point
│   ├── internal/           # Application code
│   └── migrations/         # Database migrations
└── vacaytracker-frontend/  # Svelte frontend
    ├── src/
    │   ├── routes/         # SvelteKit pages
    │   └── lib/            # Components & utilities
    └── static/             # Static assets
```

## Development

### Local Development (Without Docker)

**Backend:**

```bash
cd vacaytracker-api
cp .env.example .env
# Edit .env with your settings
make run
```

**Frontend:**

```bash
cd vacaytracker-frontend
npm install
npm run dev
```

### Useful Commands

```bash
# Backend
cd vacaytracker-api
make run              # Run server
make build            # Build binary
make test             # Run tests
make lint             # Run linter

# Frontend
cd vacaytracker-frontend
npm run dev           # Dev server
npm run build         # Production build
npm run check         # Type checking
npm run lint          # ESLint
```

## Database

VacayTracker uses SQLite for simplicity. The database file is stored in a Docker volume (`vacaytracker-data`) and persists across container restarts.

### Backup

```bash
# Create backup
docker compose exec api cp /app/data/vacaytracker.db /app/data/backup.db
docker compose cp api:/app/data/backup.db ./backup.db

# Restore backup
docker compose cp ./backup.db api:/app/data/vacaytracker.db
docker compose restart api
```

## Ports

| Service | Host Port | Container Port |
|---------|-----------|----------------|
| Frontend | 32805 | 5173 |
| API | 32804 | 3000 |

To change ports, edit `docker-compose.yml`:

```yaml
services:
  api:
    ports:
      - "YOUR_PORT:3000"
  frontend:
    ports:
      - "YOUR_PORT:5173"
```

## Troubleshooting

### Container won't start

```bash
# Check logs for errors
docker compose logs api
docker compose logs frontend

# Ensure ports aren't in use
netstat -an | grep 32804
netstat -an | grep 32805
```

### Email not sending

1. Check that `RESEND_API_KEY` and `EMAIL_FROM_ADDRESS` are set:
   ```bash
   docker compose exec api printenv | grep -E "(RESEND|EMAIL)"
   ```

2. Verify the from address is verified in Resend (or use `onboarding@resend.dev`)

3. Check API logs for email errors:
   ```bash
   docker compose logs api | grep -i email
   ```

### Database reset

```bash
# WARNING: This deletes all data!
docker compose down -v
docker compose up --build
```

## License

MIT
