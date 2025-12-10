# 10 - Development Workflow

> Git workflow, code review, and CI/CD pipeline

## Table of Contents

1. [Git Workflow](#git-workflow)
2. [Branch Naming](#branch-naming)
3. [Commit Messages](#commit-messages)
4. [Pull Request Process](#pull-request-process)
5. [Code Review Guidelines](#code-review-guidelines)
6. [CI/CD Pipeline](#cicd-pipeline)
7. [Release Process](#release-process)

---

## Git Workflow

### Branch Strategy

```
main (production)
  │
  ├── develop (integration)
  │     │
  │     ├── feature/user-management
  │     ├── feature/vacation-calendar
  │     ├── bugfix/login-validation
  │     └── hotfix/security-patch
  │
  └── release/v1.0.0
```

### Branch Types

| Branch Type | Base | Merge To | Purpose |
|-------------|------|----------|---------|
| `main` | - | - | Production-ready code |
| `develop` | main | main | Integration branch |
| `feature/*` | develop | develop | New features |
| `bugfix/*` | develop | develop | Bug fixes |
| `hotfix/*` | main | main, develop | Critical production fixes |
| `release/*` | develop | main, develop | Release preparation |

### Workflow Steps

#### 1. Starting New Work

```bash
# Update develop branch
git checkout develop
git pull origin develop

# Create feature branch
git checkout -b feature/vacation-request-modal

# Work on feature...
git add .
git commit -m "feat(vacation): add request modal component"
```

#### 2. Keeping Branch Updated

```bash
# Rebase on develop regularly
git fetch origin
git rebase origin/develop

# Or merge if preferred
git merge origin/develop
```

#### 3. Creating Pull Request

```bash
# Push branch
git push origin feature/vacation-request-modal

# Create PR via GitHub/GitLab UI
# Base: develop
# Compare: feature/vacation-request-modal
```

#### 4. After PR Approval

```bash
# Squash and merge via GitHub UI
# Or:
git checkout develop
git merge --squash feature/vacation-request-modal
git commit -m "feat(vacation): add request modal with date validation"
git push origin develop

# Delete feature branch
git branch -d feature/vacation-request-modal
git push origin --delete feature/vacation-request-modal
```

---

## Branch Naming

### Format

```
<type>/<short-description>
```

### Examples

| Type | Example | Description |
|------|---------|-------------|
| `feature/` | `feature/admin-dashboard` | New functionality |
| `bugfix/` | `bugfix/date-parsing` | Non-critical bug fix |
| `hotfix/` | `hotfix/auth-bypass` | Critical production fix |
| `refactor/` | `refactor/vacation-service` | Code improvement |
| `docs/` | `docs/api-specification` | Documentation |
| `test/` | `test/vacation-coverage` | Test improvements |
| `chore/` | `chore/update-dependencies` | Maintenance tasks |

### Rules

- Use lowercase
- Use hyphens (not underscores)
- Keep short but descriptive
- Include ticket number if applicable: `feature/VT-123-user-management`

---

## Commit Messages

### Format

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation |
| `style` | Formatting (no code change) |
| `refactor` | Code refactoring |
| `test` | Adding tests |
| `chore` | Maintenance tasks |
| `perf` | Performance improvement |
| `ci` | CI/CD changes |

### Scopes

| Scope | Description |
|-------|-------------|
| `auth` | Authentication |
| `vacation` | Vacation requests |
| `admin` | Admin features |
| `calendar` | Calendar system |
| `email` | Email notifications |
| `ui` | UI components |
| `api` | API endpoints |
| `db` | Database |
| `config` | Configuration |

### Examples

```bash
# Feature
feat(vacation): add date range validation to request form

# Bug fix
fix(auth): handle expired token gracefully

# Documentation
docs(api): add endpoint documentation for vacation routes

# Refactoring
refactor(vacation): extract business day calculation to utility

# With body
feat(admin): implement user management table

- Add sortable columns
- Add pagination
- Add search filter
- Add role filter

Closes #42

# Breaking change
feat(api)!: change date format from ISO to DD/MM/YYYY

BREAKING CHANGE: API now expects dates in DD/MM/YYYY format
```

---

## Pull Request Process

### PR Template

- [ ] **Create PR template** `.github/pull_request_template.md`

```markdown
## Summary

<!-- Brief description of changes -->

## Type of Change

- [ ] Feature (new functionality)
- [ ] Bug fix (fixes an issue)
- [ ] Refactor (code improvement, no functional change)
- [ ] Documentation
- [ ] Test
- [ ] Chore (maintenance)

## Related Issues

<!-- Link to related issues: Closes #123 -->

## Changes Made

<!-- Bullet points of specific changes -->

-
-
-

## Testing

<!-- How were these changes tested? -->

- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing performed

## Screenshots (if applicable)

<!-- Add screenshots for UI changes -->

## Checklist

- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated (if needed)
- [ ] Tests added/updated
- [ ] No console.log or debug code
- [ ] No hardcoded values that should be configurable
```

### PR Requirements

| Requirement | Description |
|-------------|-------------|
| **Title** | Clear, follows commit message format |
| **Description** | Explains what and why |
| **Tests** | All tests pass |
| **Review** | At least 1 approval |
| **Conflicts** | Resolved before merge |
| **Size** | Preferably < 400 lines changed |

### PR Size Guidelines

| Size | Lines Changed | Review Time |
|------|---------------|-------------|
| Small | < 100 | Quick |
| Medium | 100-400 | Normal |
| Large | 400-1000 | Extended |
| Too Large | > 1000 | Split into smaller PRs |

---

## Code Review Guidelines

### For Authors

1. **Self-review first** - Read your own code before requesting review
2. **Small PRs** - Easier to review, faster feedback
3. **Clear description** - Explain the context and decisions
4. **Respond promptly** - Address feedback quickly
5. **Be open** - Accept constructive criticism

### For Reviewers

1. **Be respectful** - Critique code, not people
2. **Be specific** - Point to exact lines, suggest alternatives
3. **Ask questions** - Clarify intent before assuming
4. **Prioritize** - Focus on logic, security, then style
5. **Approve quickly** - Don't block on minor issues

### Review Checklist

#### Functionality
- [ ] Does the code do what it claims?
- [ ] Are edge cases handled?
- [ ] Are error cases handled?

#### Security
- [ ] No hardcoded secrets?
- [ ] Input validation present?
- [ ] Authorization checks in place?

#### Code Quality
- [ ] Code is readable?
- [ ] No unnecessary complexity?
- [ ] DRY principle followed?
- [ ] Consistent naming?

#### Testing
- [ ] Tests exist for new code?
- [ ] Tests are meaningful?
- [ ] Edge cases tested?

#### Documentation
- [ ] Code is self-documenting?
- [ ] Complex logic has comments?
- [ ] API changes documented?

### Comment Prefixes

| Prefix | Meaning |
|--------|---------|
| `[blocking]` | Must fix before approval |
| `[suggestion]` | Nice to have, not required |
| `[question]` | Need clarification |
| `[nitpick]` | Style preference |
| `[praise]` | Good job! |

---

## CI/CD Pipeline

### GitHub Actions

- [ ] **Create CI workflow** `.github/workflows/ci.yml`

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  # Backend tests
  backend-test:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ./vacaytracker-api

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: go test -v -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./vacaytracker-api/coverage.out

  # Backend lint
  backend-lint:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ./vacaytracker-api

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          working-directory: ./vacaytracker-api

  # Frontend tests
  frontend-test:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ./vacaytracker-frontend

    steps:
      - uses: actions/checkout@v4

      - name: Set up Node
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: ./vacaytracker-frontend/package-lock.json

      - name: Install dependencies
        run: npm ci

      - name: Run linter
        run: npm run lint

      - name: Run type check
        run: npm run check

      - name: Run tests
        run: npm run test

  # Frontend build
  frontend-build:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ./vacaytracker-frontend

    steps:
      - uses: actions/checkout@v4

      - name: Set up Node
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: ./vacaytracker-frontend/package-lock.json

      - name: Install dependencies
        run: npm ci

      - name: Build
        run: npm run build

  # E2E tests (on main only)
  e2e-test:
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    needs: [backend-test, frontend-build]

    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Compose
        run: docker-compose -f docker-compose.dev.yml up -d

      - name: Wait for services
        run: sleep 30

      - name: Run E2E tests
        run: npx playwright test

      - name: Upload test results
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: playwright-report
          path: playwright-report/
```

### Deployment Workflow

- [ ] **Create deploy workflow** `.github/workflows/deploy.yml`

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    steps:
      - uses: actions/checkout@v4

      - name: Deploy to production
        uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.DEPLOY_HOST }}
          username: ${{ secrets.DEPLOY_USER }}
          key: ${{ secrets.DEPLOY_KEY }}
          script: |
            cd /opt/vacaytracker
            git pull origin main
            docker-compose build
            docker-compose up -d
            docker system prune -f
```

### Pipeline Stages

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Commit    │───>│    Build    │───>│    Test     │───>│   Deploy    │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
      │                  │                  │                  │
      │                  │                  │                  │
      v                  v                  v                  v
   Trigger          Compile           Unit Tests         Production
   Webhook          Backend           Integration        Deployment
                    Build             E2E Tests          (main only)
                    Frontend          Lint
                    Docker            Type Check
```

---

## Release Process

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):

```
MAJOR.MINOR.PATCH

1.0.0  - Initial release
1.1.0  - New feature (backward compatible)
1.1.1  - Bug fix
2.0.0  - Breaking change
```

### Release Steps

#### 1. Prepare Release

```bash
# Create release branch from develop
git checkout develop
git pull origin develop
git checkout -b release/v1.1.0

# Update version in files
# - package.json
# - version constant in Go code

# Update CHANGELOG.md
# Commit changes
git commit -am "chore: prepare release v1.1.0"
```

#### 2. Final Testing

```bash
# Run all tests
make test
npm run test

# Build and verify
docker-compose build
docker-compose up -d

# Manual smoke testing
```

#### 3. Merge and Tag

```bash
# Merge to main
git checkout main
git merge release/v1.1.0
git push origin main

# Create tag
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0

# Merge back to develop
git checkout develop
git merge main
git push origin develop

# Delete release branch
git branch -d release/v1.1.0
git push origin --delete release/v1.1.0
```

#### 4. Create GitHub Release

1. Go to GitHub Releases
2. Click "Create new release"
3. Select tag `v1.1.0`
4. Add release notes from CHANGELOG
5. Publish release

### Changelog Format

```markdown
# Changelog

## [1.1.0] - 2024-01-15

### Added
- Calendar month view with vacation events
- Email notification preferences

### Changed
- Improved date picker UX
- Enhanced mobile responsiveness

### Fixed
- Date parsing for EU format
- Session timeout handling

### Security
- Updated JWT library to latest version

## [1.0.0] - 2024-01-01

### Added
- Initial release
- User authentication
- Vacation request management
- Admin dashboard
```

---

## Development Commands Reference

### Backend (Go)

```bash
# Run locally
make run

# Run tests
make test

# Run with coverage
make test-coverage

# Lint
make lint

# Build binary
make build

# Clean build artifacts
make clean
```

### Frontend (Svelte)

```bash
# Development server
npm run dev

# Type checking
npm run check

# Lint
npm run lint

# Format
npm run format

# Build
npm run build

# Test
npm run test

# E2E tests
npm run test:e2e
```

### Docker

```bash
# Start development
docker-compose -f docker-compose.dev.yml up

# Start production
docker-compose up -d

# View logs
docker-compose logs -f

# Rebuild
docker-compose build --no-cache

# Stop
docker-compose down

# Clean volumes
docker-compose down -v
```

---

## Related Documents

- [03-implementation-roadmap.md](./03-implementation-roadmap.md) - Development phases
- [07-testing-strategy.md](./07-testing-strategy.md) - Testing requirements
- [09-deployment-guide.md](./09-deployment-guide.md) - Deployment process
