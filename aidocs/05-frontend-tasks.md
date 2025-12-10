# 05 - Frontend Tasks

> Complete Svelte 5 / SvelteKit implementation checklist with code examples

## Table of Contents

1. [Project Setup](#project-setup)
2. [Tailwind Configuration](#tailwind-configuration)
3. [TypeScript Types](#typescript-types)
4. [API Client](#api-client)
5. [Stores](#stores)
6. [Base UI Components](#base-ui-components)
7. [Layout Components](#layout-components)
8. [Auth Pages](#auth-pages)
9. [Employee Features](#employee-features)
10. [Admin Features](#admin-features)
11. [Calendar Components](#calendar-components)

---

## Project Setup

### Task 1.1: Initialize SvelteKit Project

- [ ] **Create SvelteKit project**

```bash
npm create svelte@latest vacaytracker-frontend
# Select: Skeleton project, TypeScript, ESLint, Prettier
cd vacaytracker-frontend
npm install
```

### Task 1.2: Install Dependencies

- [ ] **Install required packages** `package.json`

```bash
npm install melt lucide-svelte clsx
npm install -D tailwindcss @tailwindcss/vite
```

**package.json dependencies:**
```json
{
  "dependencies": {
    "melt": "^0.42.0",
    "lucide-svelte": "^0.300.0",
    "clsx": "^2.1.0"
  },
  "devDependencies": {
    "@sveltejs/adapter-auto": "^3.0.0",
    "@sveltejs/kit": "^2.49.0",
    "@sveltejs/vite-plugin-svelte": "^4.0.0",
    "svelte": "^5.0.0",
    "svelte-check": "^4.0.0",
    "tailwindcss": "^4.0.0",
    "@tailwindcss/vite": "^4.0.0",
    "typescript": "^5.0.0",
    "vite": "^6.0.0"
  }
}
```

### Task 1.3: Configure Vite

- [ ] **Update vite.config.ts** `vite.config.ts`

```typescript
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [tailwindcss(), sveltekit()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
    },
  },
});
```

### Task 1.4: Configure TypeScript

- [ ] **Update tsconfig.json** `tsconfig.json`

```json
{
  "extends": "./.svelte-kit/tsconfig.json",
  "compilerOptions": {
    "strict": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "paths": {
      "$lib/*": ["./src/lib/*"],
      "$components/*": ["./src/lib/components/*"]
    }
  }
}
```

---

## Tailwind Configuration

### Task 2.1: Create Global Styles

- [ ] **Create app.css** `src/app.css`

```css
@import "tailwindcss";

@theme {
  /* Beach/Ocean Theme Colors */
  --color-ocean-50: oklch(0.97 0.01 220);
  --color-ocean-100: oklch(0.93 0.03 220);
  --color-ocean-200: oklch(0.87 0.06 220);
  --color-ocean-300: oklch(0.78 0.10 220);
  --color-ocean-400: oklch(0.68 0.14 220);
  --color-ocean-500: oklch(0.55 0.18 220);
  --color-ocean-600: oklch(0.48 0.16 220);
  --color-ocean-700: oklch(0.40 0.14 220);
  --color-ocean-800: oklch(0.33 0.11 220);
  --color-ocean-900: oklch(0.27 0.08 220);

  /* Sand Colors */
  --color-sand-50: oklch(0.98 0.01 80);
  --color-sand-100: oklch(0.95 0.02 80);
  --color-sand-200: oklch(0.90 0.04 80);
  --color-sand-300: oklch(0.85 0.06 80);
  --color-sand-400: oklch(0.75 0.08 80);
  --color-sand-500: oklch(0.65 0.10 80);

  /* Semantic Colors */
  --color-success: oklch(0.65 0.20 145);
  --color-warning: oklch(0.75 0.18 85);
  --color-error: oklch(0.60 0.22 25);
  --color-info: oklch(0.65 0.15 250);

  /* Status Colors */
  --color-pending: oklch(0.75 0.15 85);
  --color-approved: oklch(0.65 0.20 145);
  --color-rejected: oklch(0.60 0.22 25);

  /* Typography */
  --font-family-sans: 'Inter', system-ui, sans-serif;
  --font-family-display: 'Cal Sans', var(--font-family-sans);

  /* Spacing Scale */
  --spacing-xs: 0.25rem;
  --spacing-sm: 0.5rem;
  --spacing-md: 1rem;
  --spacing-lg: 1.5rem;
  --spacing-xl: 2rem;

  /* Border Radius */
  --radius-sm: 0.375rem;
  --radius-md: 0.5rem;
  --radius-lg: 0.75rem;
  --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
  --shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1);
  --shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1);
}

/* Base Styles */
html {
  font-family: var(--font-family-sans);
}

body {
  @apply bg-sand-50 text-ocean-900;
}

/* Focus Styles */
*:focus-visible {
  @apply outline-2 outline-offset-2 outline-ocean-500;
}

/* Scrollbar Styles */
::-webkit-scrollbar {
  @apply w-2;
}

::-webkit-scrollbar-track {
  @apply bg-sand-100;
}

::-webkit-scrollbar-thumb {
  @apply bg-ocean-300 rounded-full;
}
```

---

## TypeScript Types

### Task 3.1: Create Type Definitions

- [ ] **Create types/index.ts** `src/lib/types/index.ts`

```typescript
// User Types
export type Role = 'admin' | 'employee';

export interface EmailPreferences {
  vacationUpdates: boolean;
  weeklyDigest: boolean;
  teamNotifications: boolean;
}

export interface User {
  id: string;
  email: string;
  name: string;
  role: Role;
  vacationBalance: number;
  startDate?: string;
  emailPreferences: EmailPreferences;
  createdAt?: string;
  updatedAt?: string;
}

// Vacation Types
export type VacationStatus = 'pending' | 'approved' | 'rejected';

export interface VacationRequest {
  id: string;
  userId: string;
  userName?: string;
  userEmail?: string;
  startDate: string;
  endDate: string;
  totalDays: number;
  reason?: string;
  status: VacationStatus;
  reviewedBy?: string;
  reviewedAt?: string;
  rejectionReason?: string;
  createdAt: string;
  updatedAt: string;
}

export interface TeamVacation {
  id: string;
  userId: string;
  userName: string;
  startDate: string;
  endDate: string;
  totalDays: number;
}

// Settings Types
export interface WeekendPolicy {
  excludeWeekends: boolean;
  excludedDays: number[];
}

export interface NewsletterConfig {
  enabled: boolean;
  frequency: 'weekly' | 'monthly';
  dayOfMonth: number;
}

export interface Settings {
  id: string;
  weekendPolicy: WeekendPolicy;
  newsletter: NewsletterConfig;
  defaultVacationDays: number;
  vacationResetMonth: number;
  updatedAt: string;
}

// API Types
export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}

export interface PaginationInfo {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface VacationListResponse {
  requests: VacationRequest[];
  total: number;
}

export interface TeamVacationResponse {
  vacations: TeamVacation[];
  month: number;
  year: number;
}

export interface UserListResponse {
  users: User[];
  pagination: PaginationInfo;
}

// Form Types
export interface CreateVacationForm {
  startDate: string;
  endDate: string;
  reason?: string;
}

export interface CreateUserForm {
  email: string;
  password: string;
  name: string;
  role: Role;
  vacationBalance?: number;
  startDate?: string;
}

export interface UpdateUserForm {
  email?: string;
  name?: string;
  role?: Role;
  vacationBalance?: number;
  startDate?: string;
}
```

---

## API Client

### Task 4.1: Create Base API Client

- [ ] **Create api/client.ts** `src/lib/api/client.ts`

```typescript
import type { ApiError } from '$lib/types';

const API_BASE = '/api';

export class ApiException extends Error {
  constructor(
    public code: string,
    message: string,
    public status: number,
    public details?: Record<string, unknown>
  ) {
    super(message);
    this.name = 'ApiException';
  }
}

function getAuthToken(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem('auth_token');
}

export function setAuthToken(token: string): void {
  localStorage.setItem('auth_token', token);
}

export function clearAuthToken(): void {
  localStorage.removeItem('auth_token');
}

export async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getAuthToken();

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    let error: ApiError;
    try {
      error = await response.json();
    } catch {
      error = {
        code: 'UNKNOWN_ERROR',
        message: 'An unknown error occurred',
      };
    }
    throw new ApiException(
      error.code,
      error.message,
      response.status,
      error.details
    );
  }

  // Handle 204 No Content
  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}
```

### Task 4.2: Create Auth API

- [ ] **Create api/auth.ts** `src/lib/api/auth.ts`

```typescript
import { request, setAuthToken, clearAuthToken } from './client';
import type { User, LoginResponse, EmailPreferences } from '$lib/types';

export const authApi = {
  login: async (email: string, password: string): Promise<LoginResponse> => {
    const response = await request<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    setAuthToken(response.token);
    return response;
  },

  logout: (): void => {
    clearAuthToken();
  },

  me: (): Promise<User> => {
    return request<User>('/auth/me');
  },

  changePassword: (currentPassword: string, newPassword: string): Promise<void> => {
    return request('/auth/password', {
      method: 'PUT',
      body: JSON.stringify({ currentPassword, newPassword }),
    });
  },

  updateEmailPreferences: (prefs: Partial<EmailPreferences>): Promise<{ emailPreferences: EmailPreferences }> => {
    return request('/auth/email-preferences', {
      method: 'PUT',
      body: JSON.stringify(prefs),
    });
  },
};
```

### Task 4.3: Create Vacation API

- [ ] **Create api/vacation.ts** `src/lib/api/vacation.ts`

```typescript
import { request } from './client';
import type {
  VacationRequest,
  VacationListResponse,
  TeamVacationResponse,
  CreateVacationForm,
} from '$lib/types';

export const vacationApi = {
  create: (data: CreateVacationForm): Promise<VacationRequest> => {
    return request('/vacation/request', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  list: (params?: { status?: string; year?: number }): Promise<VacationListResponse> => {
    const searchParams = new URLSearchParams();
    if (params?.status) searchParams.set('status', params.status);
    if (params?.year) searchParams.set('year', params.year.toString());

    const query = searchParams.toString();
    return request(`/vacation/requests${query ? `?${query}` : ''}`);
  },

  get: (id: string): Promise<VacationRequest> => {
    return request(`/vacation/requests/${id}`);
  },

  cancel: (id: string): Promise<void> => {
    return request(`/vacation/requests/${id}`, { method: 'DELETE' });
  },

  team: (month: number, year: number): Promise<TeamVacationResponse> => {
    return request(`/vacation/team?month=${month}&year=${year}`);
  },
};
```

### Task 4.4: Create Admin API

- [ ] **Create api/admin.ts** `src/lib/api/admin.ts`

```typescript
import { request } from './client';
import type {
  User,
  UserListResponse,
  VacationRequest,
  Settings,
  CreateUserForm,
  UpdateUserForm,
} from '$lib/types';

export const adminApi = {
  // Users
  listUsers: (params?: {
    page?: number;
    limit?: number;
    role?: string;
    search?: string;
  }): Promise<UserListResponse> => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set('page', params.page.toString());
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.role) searchParams.set('role', params.role);
    if (params?.search) searchParams.set('search', params.search);

    const query = searchParams.toString();
    return request(`/admin/users${query ? `?${query}` : ''}`);
  },

  createUser: (data: CreateUserForm): Promise<User> => {
    return request('/admin/users', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  getUser: (id: string): Promise<User> => {
    return request(`/admin/users/${id}`);
  },

  updateUser: (id: string, data: UpdateUserForm): Promise<User> => {
    return request(`/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  deleteUser: (id: string): Promise<void> => {
    return request(`/admin/users/${id}`, { method: 'DELETE' });
  },

  // Vacation Management
  pendingRequests: (): Promise<{ requests: VacationRequest[]; total: number }> => {
    return request('/admin/vacation/pending');
  },

  approveRequest: (id: string): Promise<VacationRequest> => {
    return request(`/admin/vacation/${id}/approve`, { method: 'PUT' });
  },

  rejectRequest: (id: string, reason?: string): Promise<VacationRequest> => {
    return request(`/admin/vacation/${id}/reject`, {
      method: 'PUT',
      body: JSON.stringify({ reason }),
    });
  },

  // Settings
  getSettings: (): Promise<Settings> => {
    return request('/admin/settings');
  },

  updateSettings: (data: Partial<Settings>): Promise<Settings> => {
    return request('/admin/settings', {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  // Newsletter
  sendNewsletter: (preview: boolean = false): Promise<{ message: string; recipientCount: number }> => {
    return request('/admin/newsletter/send', {
      method: 'POST',
      body: JSON.stringify({ preview }),
    });
  },
};
```

---

## Stores

### Task 5.1: Create Auth Store

- [ ] **Create stores/auth.svelte.ts** `src/lib/stores/auth.svelte.ts`

```typescript
import { authApi } from '$lib/api/auth';
import type { User } from '$lib/types';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

function createAuthStore() {
  let user = $state<User | null>(null);
  let isLoading = $state(true);
  let error = $state<string | null>(null);

  const isAuthenticated = $derived(user !== null);
  const isAdmin = $derived(user?.role === 'admin');
  const isEmployee = $derived(user?.role === 'employee');

  async function initialize(): Promise<void> {
    if (typeof window === 'undefined') return;

    const token = localStorage.getItem('auth_token');
    if (!token) {
      isLoading = false;
      return;
    }

    try {
      user = await authApi.me();
    } catch (e) {
      localStorage.removeItem('auth_token');
      user = null;
    } finally {
      isLoading = false;
    }
  }

  async function login(email: string, password: string): Promise<void> {
    error = null;
    isLoading = true;

    try {
      const response = await authApi.login(email, password);
      user = response.user;
    } catch (e) {
      error = e instanceof Error ? e.message : 'Login failed';
      throw e;
    } finally {
      isLoading = false;
    }
  }

  function logout(): void {
    authApi.logout();
    user = null;
  }

  function updateUser(updates: Partial<User>): void {
    if (user) {
      user = { ...user, ...updates };
    }
  }

  return {
    get user() { return user; },
    get isAuthenticated() { return isAuthenticated; },
    get isAdmin() { return isAdmin; },
    get isEmployee() { return isEmployee; },
    get isLoading() { return isLoading; },
    get error() { return error; },
    initialize,
    login,
    logout,
    updateUser,
  };
}

export const auth = createAuthStore();
```

### Task 5.2: Create Toast Store

- [ ] **Create stores/toast.svelte.ts** `src/lib/stores/toast.svelte.ts`

```typescript
export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface Toast {
  id: string;
  type: ToastType;
  message: string;
  duration?: number;
}

function createToastStore() {
  let toasts = $state<Toast[]>([]);

  function add(type: ToastType, message: string, duration: number = 5000): string {
    const id = crypto.randomUUID();
    const toast: Toast = { id, type, message, duration };

    toasts = [...toasts, toast];

    if (duration > 0) {
      setTimeout(() => dismiss(id), duration);
    }

    return id;
  }

  function dismiss(id: string): void {
    toasts = toasts.filter(t => t.id !== id);
  }

  function dismissAll(): void {
    toasts = [];
  }

  // Convenience methods
  const success = (message: string, duration?: number) => add('success', message, duration);
  const error = (message: string, duration?: number) => add('error', message, duration);
  const warning = (message: string, duration?: number) => add('warning', message, duration);
  const info = (message: string, duration?: number) => add('info', message, duration);

  return {
    get toasts() { return toasts; },
    add,
    dismiss,
    dismissAll,
    success,
    error,
    warning,
    info,
  };
}

export const toast = createToastStore();
```

### Task 5.3: Create Vacation Store

- [ ] **Create stores/vacation.svelte.ts** `src/lib/stores/vacation.svelte.ts`

```typescript
import { vacationApi } from '$lib/api/vacation';
import type { VacationRequest, VacationStatus } from '$lib/types';

function createVacationStore() {
  let requests = $state<VacationRequest[]>([]);
  let isLoading = $state(false);
  let error = $state<string | null>(null);

  const pendingRequests = $derived(requests.filter(r => r.status === 'pending'));
  const approvedRequests = $derived(requests.filter(r => r.status === 'approved'));
  const rejectedRequests = $derived(requests.filter(r => r.status === 'rejected'));

  async function fetchRequests(status?: VacationStatus, year?: number): Promise<void> {
    isLoading = true;
    error = null;

    try {
      const response = await vacationApi.list({ status, year });
      requests = response.requests;
    } catch (e) {
      error = e instanceof Error ? e.message : 'Failed to fetch requests';
    } finally {
      isLoading = false;
    }
  }

  async function createRequest(data: { startDate: string; endDate: string; reason?: string }): Promise<VacationRequest> {
    const newRequest = await vacationApi.create(data);
    requests = [newRequest, ...requests];
    return newRequest;
  }

  async function cancelRequest(id: string): Promise<void> {
    await vacationApi.cancel(id);
    requests = requests.filter(r => r.id !== id);
  }

  function updateRequest(id: string, updates: Partial<VacationRequest>): void {
    requests = requests.map(r => r.id === id ? { ...r, ...updates } : r);
  }

  return {
    get requests() { return requests; },
    get pendingRequests() { return pendingRequests; },
    get approvedRequests() { return approvedRequests; },
    get rejectedRequests() { return rejectedRequests; },
    get isLoading() { return isLoading; },
    get error() { return error; },
    fetchRequests,
    createRequest,
    cancelRequest,
    updateRequest,
  };
}

export const vacation = createVacationStore();
```

---

## Base UI Components

### Task 6.1: Create Button Component

- [ ] **Create Button.svelte** `src/lib/components/ui/Button.svelte`

```svelte
<script lang="ts">
  import { clsx } from 'clsx';
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
    size?: 'sm' | 'md' | 'lg';
    disabled?: boolean;
    loading?: boolean;
    type?: 'button' | 'submit' | 'reset';
    onclick?: (e: MouseEvent) => void;
    class?: string;
    children: Snippet;
  }

  let {
    variant = 'primary',
    size = 'md',
    disabled = false,
    loading = false,
    type = 'button',
    onclick,
    class: className = '',
    children,
  }: Props = $props();

  const baseStyles = 'inline-flex items-center justify-center font-medium rounded-md transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed';

  const variantStyles = {
    primary: 'bg-ocean-500 text-white hover:bg-ocean-600 focus:ring-ocean-500',
    secondary: 'bg-sand-200 text-ocean-900 hover:bg-sand-300 focus:ring-sand-400',
    outline: 'border-2 border-ocean-500 text-ocean-500 hover:bg-ocean-50 focus:ring-ocean-500',
    ghost: 'text-ocean-600 hover:bg-ocean-50 focus:ring-ocean-500',
    danger: 'bg-error text-white hover:bg-red-600 focus:ring-error',
  };

  const sizeStyles = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-base',
    lg: 'px-6 py-3 text-lg',
  };

  const classes = $derived(clsx(
    baseStyles,
    variantStyles[variant],
    sizeStyles[size],
    className
  ));
</script>

<button
  {type}
  class={classes}
  disabled={disabled || loading}
  {onclick}
>
  {#if loading}
    <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
  {/if}
  {@render children()}
</button>
```

### Task 6.2: Create Input Component

- [ ] **Create Input.svelte** `src/lib/components/ui/Input.svelte`

```svelte
<script lang="ts">
  import { clsx } from 'clsx';

  interface Props {
    type?: 'text' | 'email' | 'password' | 'number' | 'date';
    value?: string;
    placeholder?: string;
    label?: string;
    error?: string;
    disabled?: boolean;
    required?: boolean;
    id?: string;
    name?: string;
    class?: string;
    oninput?: (e: Event) => void;
    onblur?: (e: FocusEvent) => void;
  }

  let {
    type = 'text',
    value = $bindable(''),
    placeholder = '',
    label,
    error,
    disabled = false,
    required = false,
    id,
    name,
    class: className = '',
    oninput,
    onblur,
  }: Props = $props();

  const inputId = id || `input-${crypto.randomUUID().slice(0, 8)}`;

  const inputClasses = $derived(clsx(
    'w-full px-3 py-2 rounded-md border transition-colors',
    'focus:outline-none focus:ring-2 focus:ring-ocean-500 focus:border-transparent',
    error
      ? 'border-error text-error placeholder-error/50'
      : 'border-sand-300 text-ocean-900 placeholder-sand-400',
    disabled && 'bg-sand-100 cursor-not-allowed',
    className
  ));
</script>

<div class="w-full">
  {#if label}
    <label for={inputId} class="block text-sm font-medium text-ocean-700 mb-1">
      {label}
      {#if required}
        <span class="text-error">*</span>
      {/if}
    </label>
  {/if}

  <input
    id={inputId}
    {type}
    {name}
    {placeholder}
    {disabled}
    {required}
    class={inputClasses}
    bind:value
    {oninput}
    {onblur}
  />

  {#if error}
    <p class="mt-1 text-sm text-error">{error}</p>
  {/if}
</div>
```

### Task 6.3: Create Card Component

- [ ] **Create Card.svelte** `src/lib/components/ui/Card.svelte`

```svelte
<script lang="ts">
  import { clsx } from 'clsx';
  import type { Snippet } from 'svelte';

  interface Props {
    padding?: 'none' | 'sm' | 'md' | 'lg';
    class?: string;
    header?: Snippet;
    footer?: Snippet;
    children: Snippet;
  }

  let {
    padding = 'md',
    class: className = '',
    header,
    footer,
    children,
  }: Props = $props();

  const paddingStyles = {
    none: '',
    sm: 'p-3',
    md: 'p-4',
    lg: 'p-6',
  };

  const cardClasses = $derived(clsx(
    'bg-white rounded-lg shadow-md border border-sand-200',
    className
  ));

  const contentClasses = $derived(paddingStyles[padding]);
</script>

<div class={cardClasses}>
  {#if header}
    <div class="px-4 py-3 border-b border-sand-200">
      {@render header()}
    </div>
  {/if}

  <div class={contentClasses}>
    {@render children()}
  </div>

  {#if footer}
    <div class="px-4 py-3 border-t border-sand-200 bg-sand-50 rounded-b-lg">
      {@render footer()}
    </div>
  {/if}
</div>
```

### Task 6.4: Create Badge Component

- [ ] **Create Badge.svelte** `src/lib/components/ui/Badge.svelte`

```svelte
<script lang="ts">
  import { clsx } from 'clsx';
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'default' | 'success' | 'warning' | 'error' | 'info' | 'pending' | 'approved' | 'rejected';
    size?: 'sm' | 'md';
    class?: string;
    children: Snippet;
  }

  let {
    variant = 'default',
    size = 'md',
    class: className = '',
    children,
  }: Props = $props();

  const variantStyles = {
    default: 'bg-sand-200 text-ocean-700',
    success: 'bg-green-100 text-green-800',
    warning: 'bg-yellow-100 text-yellow-800',
    error: 'bg-red-100 text-red-800',
    info: 'bg-blue-100 text-blue-800',
    pending: 'bg-yellow-100 text-yellow-800',
    approved: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800',
  };

  const sizeStyles = {
    sm: 'px-2 py-0.5 text-xs',
    md: 'px-2.5 py-1 text-sm',
  };

  const classes = $derived(clsx(
    'inline-flex items-center font-medium rounded-full',
    variantStyles[variant],
    sizeStyles[size],
    className
  ));
</script>

<span class={classes}>
  {@render children()}
</span>
```

### Task 6.5: Create Avatar Component

- [ ] **Create Avatar.svelte** `src/lib/components/ui/Avatar.svelte`

```svelte
<script lang="ts">
  import { clsx } from 'clsx';

  interface Props {
    name?: string;
    src?: string;
    size?: 'sm' | 'md' | 'lg';
    class?: string;
  }

  let {
    name = '',
    src,
    size = 'md',
    class: className = '',
  }: Props = $props();

  const initials = $derived(
    name
      .split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  );

  const sizeStyles = {
    sm: 'w-8 h-8 text-xs',
    md: 'w-10 h-10 text-sm',
    lg: 'w-12 h-12 text-base',
  };

  const classes = $derived(clsx(
    'inline-flex items-center justify-center rounded-full bg-ocean-500 text-white font-medium',
    sizeStyles[size],
    className
  ));
</script>

{#if src}
  <img {src} alt={name} class={clsx('rounded-full object-cover', sizeStyles[size], className)} />
{:else}
  <div class={classes}>
    {initials || '?'}
  </div>
{/if}
```

### Task 6.6: Create ProgressRing Component

- [ ] **Create ProgressRing.svelte** `src/lib/components/ui/ProgressRing.svelte`

```svelte
<script lang="ts">
  interface Props {
    value: number;
    max: number;
    size?: number;
    strokeWidth?: number;
    class?: string;
  }

  let {
    value,
    max,
    size = 120,
    strokeWidth = 8,
    class: className = '',
  }: Props = $props();

  const percentage = $derived(Math.min(100, Math.max(0, (value / max) * 100)));
  const radius = $derived((size - strokeWidth) / 2);
  const circumference = $derived(2 * Math.PI * radius);
  const offset = $derived(circumference - (percentage / 100) * circumference);

  const color = $derived(
    percentage > 70 ? 'text-success' :
    percentage > 30 ? 'text-warning' :
    'text-error'
  );
</script>

<div class="relative inline-flex items-center justify-center {className}">
  <svg width={size} height={size} class="-rotate-90">
    <!-- Background circle -->
    <circle
      cx={size / 2}
      cy={size / 2}
      r={radius}
      fill="none"
      stroke="currentColor"
      stroke-width={strokeWidth}
      class="text-sand-200"
    />
    <!-- Progress circle -->
    <circle
      cx={size / 2}
      cy={size / 2}
      r={radius}
      fill="none"
      stroke="currentColor"
      stroke-width={strokeWidth}
      stroke-linecap="round"
      stroke-dasharray={circumference}
      stroke-dashoffset={offset}
      class={color}
      style="transition: stroke-dashoffset 0.5s ease"
    />
  </svg>
  <div class="absolute inset-0 flex flex-col items-center justify-center">
    <span class="text-2xl font-bold text-ocean-900">{value}</span>
    <span class="text-sm text-ocean-500">of {max}</span>
  </div>
</div>
```

### Task 6.7: Create Toast Component

- [ ] **Create Toast.svelte** `src/lib/components/ui/Toast.svelte`

```svelte
<script lang="ts">
  import { toast, type Toast } from '$lib/stores/toast.svelte';
  import { X, CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-svelte';
  import { clsx } from 'clsx';

  const icons = {
    success: CheckCircle,
    error: AlertCircle,
    warning: AlertTriangle,
    info: Info,
  };

  const styles = {
    success: 'bg-green-50 border-green-500 text-green-800',
    error: 'bg-red-50 border-red-500 text-red-800',
    warning: 'bg-yellow-50 border-yellow-500 text-yellow-800',
    info: 'bg-blue-50 border-blue-500 text-blue-800',
  };
</script>

<div class="fixed top-4 right-4 z-50 flex flex-col gap-2 w-80">
  {#each toast.toasts as t (t.id)}
    {@const Icon = icons[t.type]}
    <div
      class={clsx(
        'flex items-start gap-3 p-4 rounded-lg border-l-4 shadow-lg',
        styles[t.type]
      )}
      role="alert"
    >
      <Icon class="w-5 h-5 flex-shrink-0" />
      <p class="flex-1 text-sm">{t.message}</p>
      <button
        onclick={() => toast.dismiss(t.id)}
        class="flex-shrink-0 hover:opacity-70"
      >
        <X class="w-4 h-4" />
      </button>
    </div>
  {/each}
</div>
```

---

## Layout Components

### Task 7.1: Create Root Layout

- [ ] **Create +layout.svelte** `src/routes/+layout.svelte`

```svelte
<script lang="ts">
  import '../app.css';
  import { auth } from '$lib/stores/auth.svelte';
  import Toast from '$lib/components/ui/Toast.svelte';
  import type { Snippet } from 'svelte';

  interface Props {
    children: Snippet;
  }

  let { children }: Props = $props();

  $effect(() => {
    auth.initialize();
  });
</script>

<Toast />

{#if auth.isLoading}
  <div class="min-h-screen flex items-center justify-center">
    <div class="animate-spin rounded-full h-12 w-12 border-4 border-ocean-500 border-t-transparent"></div>
  </div>
{:else}
  {@render children()}
{/if}
```

### Task 7.2: Create Employee Layout

- [ ] **Create employee/+layout.svelte** `src/routes/employee/+layout.svelte`

```svelte
<script lang="ts">
  import { auth } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';
  import EmployeeHeader from '$lib/components/layout/EmployeeHeader.svelte';
  import type { Snippet } from 'svelte';

  interface Props {
    children: Snippet;
  }

  let { children }: Props = $props();

  $effect(() => {
    if (!auth.isLoading && !auth.isAuthenticated) {
      goto('/');
    }
  });
</script>

{#if auth.isAuthenticated}
  <div class="min-h-screen bg-sand-50">
    <EmployeeHeader />
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {@render children()}
    </main>
  </div>
{/if}
```

### Task 7.3: Create Admin Layout

- [ ] **Create admin/+layout.svelte** `src/routes/admin/+layout.svelte`

```svelte
<script lang="ts">
  import { auth } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';
  import AdminSidebar from '$lib/components/layout/AdminSidebar.svelte';
  import AdminHeader from '$lib/components/layout/AdminHeader.svelte';
  import type { Snippet } from 'svelte';

  interface Props {
    children: Snippet;
  }

  let { children }: Props = $props();

  $effect(() => {
    if (!auth.isLoading && (!auth.isAuthenticated || !auth.isAdmin)) {
      goto('/');
    }
  });
</script>

{#if auth.isAuthenticated && auth.isAdmin}
  <div class="min-h-screen bg-sand-50 flex">
    <AdminSidebar />
    <div class="flex-1 flex flex-col">
      <AdminHeader />
      <main class="flex-1 p-6">
        {@render children()}
      </main>
    </div>
  </div>
{/if}
```

### Task 7.4: Create Employee Header

- [ ] **Create EmployeeHeader.svelte** `src/lib/components/layout/EmployeeHeader.svelte`

```svelte
<script lang="ts">
  import { auth } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';
  import Avatar from '$lib/components/ui/Avatar.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { Palmtree, Calendar, User, LogOut, Settings } from 'lucide-svelte';

  function handleLogout() {
    auth.logout();
    goto('/');
  }
</script>

<header class="bg-white shadow-sm border-b border-sand-200">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="flex items-center justify-between h-16">
      <!-- Logo -->
      <a href="/employee" class="flex items-center gap-2">
        <Palmtree class="w-8 h-8 text-ocean-500" />
        <span class="text-xl font-bold text-ocean-700">VacayTracker</span>
      </a>

      <!-- Navigation -->
      <nav class="hidden md:flex items-center gap-6">
        <a href="/employee" class="text-ocean-600 hover:text-ocean-800 font-medium">
          Dashboard
        </a>
        <a href="/employee/calendar" class="text-ocean-600 hover:text-ocean-800 font-medium">
          <Calendar class="w-4 h-4 inline mr-1" />
          Calendar
        </a>
        <a href="/employee/team" class="text-ocean-600 hover:text-ocean-800 font-medium">
          Team
        </a>
      </nav>

      <!-- User Menu -->
      <div class="flex items-center gap-4">
        <a href="/employee/settings" class="text-ocean-600 hover:text-ocean-800">
          <Settings class="w-5 h-5" />
        </a>
        <div class="flex items-center gap-2">
          <Avatar name={auth.user?.name} size="sm" />
          <span class="hidden sm:block text-sm font-medium text-ocean-700">
            {auth.user?.name}
          </span>
        </div>
        <Button variant="ghost" size="sm" onclick={handleLogout}>
          <LogOut class="w-4 h-4" />
        </Button>
      </div>
    </div>
  </div>
</header>
```

### Task 7.5: Create Admin Sidebar

- [ ] **Create AdminSidebar.svelte** `src/lib/components/layout/AdminSidebar.svelte`

```svelte
<script lang="ts">
  import { page } from '$app/stores';
  import { clsx } from 'clsx';
  import {
    Palmtree,
    LayoutDashboard,
    Users,
    Calendar,
    Settings,
    Mail,
    Wallet,
  } from 'lucide-svelte';

  const navItems = [
    { href: '/admin', icon: LayoutDashboard, label: 'Dashboard' },
    { href: '/admin/users', icon: Users, label: 'Users' },
    { href: '/admin/balances', icon: Wallet, label: 'Balances' },
    { href: '/admin/calendar', icon: Calendar, label: 'Calendar' },
    { href: '/admin/settings', icon: Settings, label: 'Settings' },
  ];
</script>

<aside class="w-64 bg-ocean-900 text-white min-h-screen flex flex-col">
  <!-- Logo -->
  <div class="p-4 border-b border-ocean-700">
    <a href="/admin" class="flex items-center gap-2">
      <Palmtree class="w-8 h-8 text-ocean-300" />
      <div>
        <span class="text-lg font-bold">VacayTracker</span>
        <span class="block text-xs text-ocean-400">Captain's Deck</span>
      </div>
    </a>
  </div>

  <!-- Navigation -->
  <nav class="flex-1 p-4">
    <ul class="space-y-1">
      {#each navItems as item}
        {@const isActive = $page.url.pathname === item.href}
        <li>
          <a
            href={item.href}
            class={clsx(
              'flex items-center gap-3 px-3 py-2 rounded-md transition-colors',
              isActive
                ? 'bg-ocean-700 text-white'
                : 'text-ocean-300 hover:bg-ocean-800 hover:text-white'
            )}
          >
            <item.icon class="w-5 h-5" />
            <span>{item.label}</span>
          </a>
        </li>
      {/each}
    </ul>
  </nav>

  <!-- Footer -->
  <div class="p-4 border-t border-ocean-700 text-ocean-400 text-sm">
    <p>VacayTracker v1.0</p>
  </div>
</aside>
```

---

## Auth Pages

### Task 8.1: Create Login Page

- [ ] **Create +page.svelte** `src/routes/+page.svelte`

```svelte
<script lang="ts">
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth.svelte';
  import { toast } from '$lib/stores/toast.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Card from '$lib/components/ui/Card.svelte';
  import { Palmtree, Waves } from 'lucide-svelte';

  let email = $state('');
  let password = $state('');
  let isSubmitting = $state(false);
  let errors = $state<{ email?: string; password?: string }>({});

  // Redirect if already authenticated
  $effect(() => {
    if (auth.isAuthenticated) {
      goto(auth.isAdmin ? '/admin' : '/employee');
    }
  });

  function validate(): boolean {
    errors = {};

    if (!email) {
      errors.email = 'Email is required';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      errors.email = 'Invalid email format';
    }

    if (!password) {
      errors.password = 'Password is required';
    } else if (password.length < 6) {
      errors.password = 'Password must be at least 6 characters';
    }

    return Object.keys(errors).length === 0;
  }

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();

    if (!validate()) return;

    isSubmitting = true;

    try {
      await auth.login(email, password);
      toast.success('Welcome back!');
      goto(auth.isAdmin ? '/admin' : '/employee');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Login failed');
    } finally {
      isSubmitting = false;
    }
  }
</script>

<svelte:head>
  <title>Login - VacayTracker</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-ocean-400 to-ocean-600 p-4">
  <!-- Decorative waves -->
  <div class="absolute bottom-0 left-0 right-0 text-ocean-700/20">
    <Waves class="w-full h-32" />
  </div>

  <Card class="w-full max-w-md relative z-10">
    {#snippet header()}
      <div class="text-center">
        <div class="flex justify-center mb-2">
          <Palmtree class="w-12 h-12 text-ocean-500" />
        </div>
        <h1 class="text-2xl font-bold text-ocean-800">VacayTracker</h1>
        <p class="text-ocean-600">Sign in to manage your vacation</p>
      </div>
    {/snippet}

    <form onsubmit={handleSubmit} class="space-y-4">
      <Input
        type="email"
        label="Email"
        placeholder="you@company.com"
        bind:value={email}
        error={errors.email}
        required
      />

      <Input
        type="password"
        label="Password"
        placeholder="Enter your password"
        bind:value={password}
        error={errors.password}
        required
      />

      <Button
        type="submit"
        variant="primary"
        class="w-full"
        loading={isSubmitting}
      >
        Sign In
      </Button>
    </form>

    {#snippet footer()}
      <p class="text-center text-sm text-ocean-500">
        Contact your administrator if you need access
      </p>
    {/snippet}
  </Card>
</div>
```

---

## Employee Features

### Task 9.1: Create Employee Dashboard

- [ ] **Create employee/+page.svelte** `src/routes/employee/+page.svelte`

```svelte
<script lang="ts">
  import { auth } from '$lib/stores/auth.svelte';
  import { vacation } from '$lib/stores/vacation.svelte';
  import BalanceDisplay from '$lib/components/features/vacation/BalanceDisplay.svelte';
  import RequestList from '$lib/components/features/vacation/RequestList.svelte';
  import RequestModal from '$lib/components/features/vacation/RequestModal.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Card from '$lib/components/ui/Card.svelte';
  import { Plus, Umbrella } from 'lucide-svelte';

  let isRequestModalOpen = $state(false);

  $effect(() => {
    vacation.fetchRequests();
  });
</script>

<svelte:head>
  <title>Dashboard - VacayTracker</title>
</svelte:head>

<div class="space-y-6">
  <!-- Welcome Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold text-ocean-800">
        Welcome back, {auth.user?.name?.split(' ')[0]}!
      </h1>
      <p class="text-ocean-600">Here's your vacation overview</p>
    </div>
    <Button onclick={() => isRequestModalOpen = true}>
      <Plus class="w-4 h-4 mr-2" />
      Request Vacation
    </Button>
  </div>

  <!-- Balance Section -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
    <Card class="md:col-span-1">
      <div class="flex flex-col items-center py-4">
        <h2 class="text-lg font-semibold text-ocean-700 mb-4">Your Balance</h2>
        <BalanceDisplay
          current={auth.user?.vacationBalance ?? 0}
          total={25}
        />
        <p class="mt-4 text-sm text-ocean-500">
          Days remaining this year
        </p>
      </div>
    </Card>

    <!-- Quick Stats -->
    <Card class="md:col-span-2">
      <h2 class="text-lg font-semibold text-ocean-700 mb-4">Quick Stats</h2>
      <div class="grid grid-cols-3 gap-4">
        <div class="text-center p-4 bg-yellow-50 rounded-lg">
          <p class="text-2xl font-bold text-yellow-600">
            {vacation.pendingRequests.length}
          </p>
          <p class="text-sm text-yellow-700">Pending</p>
        </div>
        <div class="text-center p-4 bg-green-50 rounded-lg">
          <p class="text-2xl font-bold text-green-600">
            {vacation.approvedRequests.length}
          </p>
          <p class="text-sm text-green-700">Approved</p>
        </div>
        <div class="text-center p-4 bg-ocean-50 rounded-lg">
          <p class="text-2xl font-bold text-ocean-600">
            {vacation.approvedRequests.reduce((sum, r) => sum + r.totalDays, 0)}
          </p>
          <p class="text-sm text-ocean-700">Days Used</p>
        </div>
      </div>
    </Card>
  </div>

  <!-- Recent Requests -->
  <Card>
    {#snippet header()}
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold text-ocean-700">Your Requests</h2>
        <Umbrella class="w-5 h-5 text-ocean-400" />
      </div>
    {/snippet}

    {#if vacation.isLoading}
      <div class="py-8 text-center text-ocean-500">
        Loading requests...
      </div>
    {:else if vacation.requests.length === 0}
      <div class="py-8 text-center">
        <Umbrella class="w-12 h-12 mx-auto text-ocean-300 mb-2" />
        <p class="text-ocean-500">No vacation requests yet</p>
        <Button
          variant="outline"
          class="mt-4"
          onclick={() => isRequestModalOpen = true}
        >
          Request Your First Vacation
        </Button>
      </div>
    {:else}
      <RequestList requests={vacation.requests} />
    {/if}
  </Card>
</div>

<RequestModal bind:open={isRequestModalOpen} />
```

### Task 9.2: Create Balance Display Component

- [ ] **Create BalanceDisplay.svelte** `src/lib/components/features/vacation/BalanceDisplay.svelte`

```svelte
<script lang="ts">
  import ProgressRing from '$lib/components/ui/ProgressRing.svelte';

  interface Props {
    current: number;
    total: number;
    size?: number;
  }

  let { current, total, size = 140 }: Props = $props();
</script>

<ProgressRing value={current} max={total} {size} />
```

### Task 9.3: Create Request Modal

- [ ] **Create RequestModal.svelte** `src/lib/components/features/vacation/RequestModal.svelte`

```svelte
<script lang="ts">
  import { Dialog } from 'melt/builders';
  import { vacation } from '$lib/stores/vacation.svelte';
  import { auth } from '$lib/stores/auth.svelte';
  import { toast } from '$lib/stores/toast.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import { X, Calendar, Palmtree } from 'lucide-svelte';

  interface Props {
    open: boolean;
  }

  let { open = $bindable(false) }: Props = $props();

  let startDate = $state('');
  let endDate = $state('');
  let reason = $state('');
  let isSubmitting = $state(false);
  let errors = $state<{ startDate?: string; endDate?: string }>({});

  const dialog = new Dialog({
    open: () => open,
    onOpenChange: (v) => {
      open = v;
      if (!v) resetForm();
    },
  });

  function resetForm() {
    startDate = '';
    endDate = '';
    reason = '';
    errors = {};
  }

  function validate(): boolean {
    errors = {};

    if (!startDate) {
      errors.startDate = 'Start date is required';
    }

    if (!endDate) {
      errors.endDate = 'End date is required';
    }

    // Additional date validation can be added here

    return Object.keys(errors).length === 0;
  }

  async function handleSubmit() {
    if (!validate()) return;

    isSubmitting = true;

    try {
      // Convert YYYY-MM-DD to DD/MM/YYYY for API
      const formatDate = (d: string) => {
        const [year, month, day] = d.split('-');
        return `${day}/${month}/${year}`;
      };

      await vacation.createRequest({
        startDate: formatDate(startDate),
        endDate: formatDate(endDate),
        reason: reason || undefined,
      });

      // Update user balance
      auth.updateUser({
        vacationBalance: (auth.user?.vacationBalance ?? 0) - 1, // Approximate
      });

      toast.success('Vacation request submitted!');
      open = false;
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to submit request');
    } finally {
      isSubmitting = false;
    }
  }
</script>

{#if dialog.open}
  <div class="fixed inset-0 z-50">
    <!-- Overlay -->
    <div
      {...dialog.overlay}
      class="fixed inset-0 bg-black/50 backdrop-blur-sm"
    ></div>

    <!-- Content -->
    <div class="fixed inset-0 flex items-center justify-center p-4">
      <div
        {...dialog.content}
        class="bg-white rounded-xl shadow-xl w-full max-w-md"
      >
        <!-- Header -->
        <div class="flex items-center justify-between p-4 border-b border-sand-200">
          <div class="flex items-center gap-2">
            <Palmtree class="w-5 h-5 text-ocean-500" />
            <h2 class="text-lg font-semibold text-ocean-800">
              Request Vacation
            </h2>
          </div>
          <button
            {...dialog.close}
            class="text-ocean-400 hover:text-ocean-600"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Body -->
        <form
          onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}
          class="p-4 space-y-4"
        >
          <div class="grid grid-cols-2 gap-4">
            <Input
              type="date"
              label="Start Date"
              bind:value={startDate}
              error={errors.startDate}
              required
            />
            <Input
              type="date"
              label="End Date"
              bind:value={endDate}
              error={errors.endDate}
              required
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-ocean-700 mb-1">
              Reason (optional)
            </label>
            <textarea
              bind:value={reason}
              class="w-full px-3 py-2 rounded-md border border-sand-300 focus:ring-2 focus:ring-ocean-500 focus:border-transparent"
              rows="3"
              placeholder="Family vacation, personal time, etc."
            ></textarea>
          </div>

          <!-- Balance Info -->
          <div class="bg-ocean-50 rounded-lg p-3 text-sm">
            <p class="text-ocean-700">
              <strong>Current Balance:</strong> {auth.user?.vacationBalance} days
            </p>
          </div>

          <!-- Actions -->
          <div class="flex gap-3 pt-2">
            <Button
              type="button"
              variant="outline"
              class="flex-1"
              onclick={() => open = false}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              variant="primary"
              class="flex-1"
              loading={isSubmitting}
            >
              Submit Request
            </Button>
          </div>
        </form>
      </div>
    </div>
  </div>
{/if}
```

### Task 9.4: Create Request List Component

- [ ] **Create RequestList.svelte** `src/lib/components/features/vacation/RequestList.svelte`

```svelte
<script lang="ts">
  import type { VacationRequest } from '$lib/types';
  import RequestCard from './RequestCard.svelte';

  interface Props {
    requests: VacationRequest[];
  }

  let { requests }: Props = $props();
</script>

<div class="divide-y divide-sand-200">
  {#each requests as request (request.id)}
    <RequestCard {request} />
  {/each}
</div>
```

### Task 9.5: Create Request Card Component

- [ ] **Create RequestCard.svelte** `src/lib/components/features/vacation/RequestCard.svelte`

```svelte
<script lang="ts">
  import type { VacationRequest } from '$lib/types';
  import { vacation } from '$lib/stores/vacation.svelte';
  import { toast } from '$lib/stores/toast.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { Calendar, Trash2 } from 'lucide-svelte';

  interface Props {
    request: VacationRequest;
  }

  let { request }: Props = $props();

  let isCancelling = $state(false);

  function formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-GB', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    });
  }

  async function handleCancel() {
    if (!confirm('Are you sure you want to cancel this request?')) return;

    isCancelling = true;
    try {
      await vacation.cancelRequest(request.id);
      toast.success('Request cancelled');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to cancel');
    } finally {
      isCancelling = false;
    }
  }
</script>

<div class="flex items-center justify-between py-4">
  <div class="flex items-center gap-4">
    <div class="p-2 bg-ocean-100 rounded-lg">
      <Calendar class="w-5 h-5 text-ocean-600" />
    </div>
    <div>
      <p class="font-medium text-ocean-800">
        {formatDate(request.startDate)} - {formatDate(request.endDate)}
      </p>
      <p class="text-sm text-ocean-500">
        {request.totalDays} day{request.totalDays !== 1 ? 's' : ''}
        {#if request.reason}
          &middot; {request.reason}
        {/if}
      </p>
    </div>
  </div>

  <div class="flex items-center gap-3">
    <Badge variant={request.status}>
      {request.status}
    </Badge>

    {#if request.status === 'pending'}
      <Button
        variant="ghost"
        size="sm"
        onclick={handleCancel}
        loading={isCancelling}
      >
        <Trash2 class="w-4 h-4 text-error" />
      </Button>
    {/if}
  </div>
</div>
```

---

## Admin Features

### Task 10.1: Create Admin Dashboard

- [ ] **Create admin/+page.svelte** `src/routes/admin/+page.svelte`

```svelte
<script lang="ts">
  import { adminApi } from '$lib/api/admin';
  import Card from '$lib/components/ui/Card.svelte';
  import StatsCard from '$lib/components/features/admin/StatsCard.svelte';
  import PendingRequests from '$lib/components/features/admin/PendingRequests.svelte';
  import { Users, Clock, CheckCircle, XCircle } from 'lucide-svelte';
  import type { VacationRequest, User } from '$lib/types';

  let pendingRequests = $state<VacationRequest[]>([]);
  let totalUsers = $state(0);
  let isLoading = $state(true);

  $effect(() => {
    loadData();
  });

  async function loadData() {
    isLoading = true;
    try {
      const [requestsRes, usersRes] = await Promise.all([
        adminApi.pendingRequests(),
        adminApi.listUsers({ limit: 1 }),
      ]);
      pendingRequests = requestsRes.requests;
      totalUsers = usersRes.pagination.total;
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Admin Dashboard - VacayTracker</title>
</svelte:head>

<div class="space-y-6">
  <h1 class="text-2xl font-bold text-ocean-800">Captain's Dashboard</h1>

  <!-- Stats Grid -->
  <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
    <StatsCard
      title="Total Crew"
      value={totalUsers}
      icon={Users}
      color="ocean"
    />
    <StatsCard
      title="Pending Requests"
      value={pendingRequests.length}
      icon={Clock}
      color="yellow"
    />
    <StatsCard
      title="Approved Today"
      value={0}
      icon={CheckCircle}
      color="green"
    />
    <StatsCard
      title="Rejected Today"
      value={0}
      icon={XCircle}
      color="red"
    />
  </div>

  <!-- Pending Requests -->
  <Card>
    {#snippet header()}
      <h2 class="text-lg font-semibold text-ocean-700">
        Pending Requests ({pendingRequests.length})
      </h2>
    {/snippet}

    {#if isLoading}
      <div class="py-8 text-center text-ocean-500">Loading...</div>
    {:else}
      <PendingRequests
        requests={pendingRequests}
        onUpdate={loadData}
      />
    {/if}
  </Card>
</div>
```

### Task 10.2: Create Stats Card Component

- [ ] **Create StatsCard.svelte** `src/lib/components/features/admin/StatsCard.svelte`

```svelte
<script lang="ts">
  import type { ComponentType } from 'svelte';
  import { clsx } from 'clsx';

  interface Props {
    title: string;
    value: number | string;
    icon: ComponentType;
    color?: 'ocean' | 'green' | 'yellow' | 'red';
  }

  let { title, value, icon: Icon, color = 'ocean' }: Props = $props();

  const colorStyles = {
    ocean: 'bg-ocean-50 text-ocean-600',
    green: 'bg-green-50 text-green-600',
    yellow: 'bg-yellow-50 text-yellow-600',
    red: 'bg-red-50 text-red-600',
  };
</script>

<div class="bg-white rounded-lg shadow-md border border-sand-200 p-4">
  <div class="flex items-center gap-3">
    <div class={clsx('p-2 rounded-lg', colorStyles[color])}>
      <Icon class="w-6 h-6" />
    </div>
    <div>
      <p class="text-2xl font-bold text-ocean-900">{value}</p>
      <p class="text-sm text-ocean-500">{title}</p>
    </div>
  </div>
</div>
```

### Task 10.3: Create Pending Requests Component

- [ ] **Create PendingRequests.svelte** `src/lib/components/features/admin/PendingRequests.svelte`

```svelte
<script lang="ts">
  import type { VacationRequest } from '$lib/types';
  import { adminApi } from '$lib/api/admin';
  import { toast } from '$lib/stores/toast.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  import Avatar from '$lib/components/ui/Avatar.svelte';
  import { Check, X, Clock } from 'lucide-svelte';

  interface Props {
    requests: VacationRequest[];
    onUpdate: () => void;
  }

  let { requests, onUpdate }: Props = $props();

  let processingId = $state<string | null>(null);

  function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString('en-GB', {
      day: 'numeric',
      month: 'short',
    });
  }

  async function approve(id: string) {
    processingId = id;
    try {
      await adminApi.approveRequest(id);
      toast.success('Request approved');
      onUpdate();
    } catch (error) {
      toast.error('Failed to approve request');
    } finally {
      processingId = null;
    }
  }

  async function reject(id: string) {
    const reason = prompt('Rejection reason (optional):');
    processingId = id;
    try {
      await adminApi.rejectRequest(id, reason || undefined);
      toast.success('Request rejected');
      onUpdate();
    } catch (error) {
      toast.error('Failed to reject request');
    } finally {
      processingId = null;
    }
  }
</script>

{#if requests.length === 0}
  <div class="py-8 text-center">
    <Clock class="w-12 h-12 mx-auto text-ocean-300 mb-2" />
    <p class="text-ocean-500">No pending requests</p>
  </div>
{:else}
  <div class="divide-y divide-sand-200">
    {#each requests as request (request.id)}
      <div class="flex items-center justify-between py-4">
        <div class="flex items-center gap-4">
          <Avatar name={request.userName} size="md" />
          <div>
            <p class="font-medium text-ocean-800">{request.userName}</p>
            <p class="text-sm text-ocean-500">
              {formatDate(request.startDate)} - {formatDate(request.endDate)}
              ({request.totalDays} days)
            </p>
            {#if request.reason}
              <p class="text-xs text-ocean-400 mt-1">{request.reason}</p>
            {/if}
          </div>
        </div>

        <div class="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onclick={() => reject(request.id)}
            disabled={processingId !== null}
          >
            <X class="w-4 h-4 mr-1" />
            Reject
          </Button>
          <Button
            variant="primary"
            size="sm"
            onclick={() => approve(request.id)}
            loading={processingId === request.id}
            disabled={processingId !== null && processingId !== request.id}
          >
            <Check class="w-4 h-4 mr-1" />
            Approve
          </Button>
        </div>
      </div>
    {/each}
  </div>
{/if}
```

---

## Calendar Components

### Task 11.1: Create Calendar Page

- [ ] **Create calendar/+page.svelte** `src/routes/employee/calendar/+page.svelte`

```svelte
<script lang="ts">
  import { vacationApi } from '$lib/api/vacation';
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import MonthView from '$lib/components/features/calendar/MonthView.svelte';
  import { ChevronLeft, ChevronRight, Calendar as CalendarIcon } from 'lucide-svelte';
  import type { TeamVacation } from '$lib/types';

  let currentDate = $state(new Date());
  let vacations = $state<TeamVacation[]>([]);
  let isLoading = $state(true);

  const currentMonth = $derived(currentDate.getMonth() + 1);
  const currentYear = $derived(currentDate.getFullYear());
  const monthName = $derived(currentDate.toLocaleDateString('en-US', { month: 'long', year: 'numeric' }));

  $effect(() => {
    loadVacations();
  });

  async function loadVacations() {
    isLoading = true;
    try {
      const response = await vacationApi.team(currentMonth, currentYear);
      vacations = response.vacations;
    } finally {
      isLoading = false;
    }
  }

  function previousMonth() {
    currentDate = new Date(currentYear, currentMonth - 2, 1);
  }

  function nextMonth() {
    currentDate = new Date(currentYear, currentMonth, 1);
  }

  function goToToday() {
    currentDate = new Date();
  }
</script>

<svelte:head>
  <title>Team Calendar - VacayTracker</title>
</svelte:head>

<div class="space-y-6">
  <Card>
    {#snippet header()}
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <CalendarIcon class="w-5 h-5 text-ocean-500" />
          <h2 class="text-lg font-semibold text-ocean-700">Team Calendar</h2>
        </div>
        <div class="flex items-center gap-2">
          <Button variant="ghost" size="sm" onclick={previousMonth}>
            <ChevronLeft class="w-4 h-4" />
          </Button>
          <span class="text-ocean-800 font-medium min-w-[150px] text-center">
            {monthName}
          </span>
          <Button variant="ghost" size="sm" onclick={nextMonth}>
            <ChevronRight class="w-4 h-4" />
          </Button>
          <Button variant="outline" size="sm" onclick={goToToday}>
            Today
          </Button>
        </div>
      </div>
    {/snippet}

    {#if isLoading}
      <div class="py-8 text-center text-ocean-500">Loading calendar...</div>
    {:else}
      <MonthView
        year={currentYear}
        month={currentMonth}
        {vacations}
      />
    {/if}
  </Card>
</div>
```

### Task 11.2: Create Month View Component

- [ ] **Create MonthView.svelte** `src/lib/components/features/calendar/MonthView.svelte`

```svelte
<script lang="ts">
  import type { TeamVacation } from '$lib/types';
  import DayCell from './DayCell.svelte';

  interface Props {
    year: number;
    month: number;
    vacations: TeamVacation[];
  }

  let { year, month, vacations }: Props = $props();

  const weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

  const days = $derived(() => {
    const firstDay = new Date(year, month - 1, 1);
    const lastDay = new Date(year, month, 0);
    const daysInMonth = lastDay.getDate();
    const startWeekday = firstDay.getDay();

    const result: (number | null)[] = [];

    // Add empty cells for days before the first of the month
    for (let i = 0; i < startWeekday; i++) {
      result.push(null);
    }

    // Add all days of the month
    for (let day = 1; day <= daysInMonth; day++) {
      result.push(day);
    }

    return result;
  });

  function getVacationsForDay(day: number): TeamVacation[] {
    const dateStr = `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
    return vacations.filter(v => {
      return dateStr >= v.startDate && dateStr <= v.endDate;
    });
  }
</script>

<div class="grid grid-cols-7 gap-px bg-sand-200 rounded-lg overflow-hidden">
  <!-- Header -->
  {#each weekDays as day}
    <div class="bg-ocean-50 py-2 text-center text-sm font-medium text-ocean-700">
      {day}
    </div>
  {/each}

  <!-- Days -->
  {#each days() as day}
    {#if day === null}
      <div class="bg-white min-h-[100px]"></div>
    {:else}
      <DayCell
        {day}
        vacations={getVacationsForDay(day)}
        isToday={new Date().getDate() === day &&
                 new Date().getMonth() === month - 1 &&
                 new Date().getFullYear() === year}
      />
    {/if}
  {/each}
</div>
```

### Task 11.3: Create Day Cell Component

- [ ] **Create DayCell.svelte** `src/lib/components/features/calendar/DayCell.svelte`

```svelte
<script lang="ts">
  import type { TeamVacation } from '$lib/types';
  import { clsx } from 'clsx';

  interface Props {
    day: number;
    vacations: TeamVacation[];
    isToday?: boolean;
  }

  let { day, vacations, isToday = false }: Props = $props();

  // Generate consistent color based on user ID
  function getUserColor(userId: string): string {
    const colors = [
      'bg-ocean-200 text-ocean-800',
      'bg-green-200 text-green-800',
      'bg-purple-200 text-purple-800',
      'bg-pink-200 text-pink-800',
      'bg-yellow-200 text-yellow-800',
      'bg-red-200 text-red-800',
    ];
    let hash = 0;
    for (let i = 0; i < userId.length; i++) {
      hash = userId.charCodeAt(i) + ((hash << 5) - hash);
    }
    return colors[Math.abs(hash) % colors.length];
  }
</script>

<div class={clsx(
  'bg-white min-h-[100px] p-2',
  isToday && 'ring-2 ring-inset ring-ocean-500'
)}>
  <span class={clsx(
    'inline-flex items-center justify-center w-7 h-7 text-sm rounded-full',
    isToday ? 'bg-ocean-500 text-white font-bold' : 'text-ocean-700'
  )}>
    {day}
  </span>

  {#if vacations.length > 0}
    <div class="mt-1 space-y-1">
      {#each vacations.slice(0, 3) as vacation}
        <div class={clsx(
          'text-xs px-1.5 py-0.5 rounded truncate',
          getUserColor(vacation.userId)
        )}>
          {vacation.userName}
        </div>
      {/each}
      {#if vacations.length > 3}
        <div class="text-xs text-ocean-500">
          +{vacations.length - 3} more
        </div>
      {/if}
    </div>
  {/if}
</div>
```

---

## Related Documents

- [03-implementation-roadmap.md](./03-implementation-roadmap.md) - Phase dependencies
- [04-backend-tasks.md](./04-backend-tasks.md) - API implementation
- [06-component-inventory.md](./06-component-inventory.md) - Component specifications
- [07-testing-strategy.md](./07-testing-strategy.md) - Frontend testing
