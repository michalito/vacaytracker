# VacayTracker - Frontend Implementation Guide

> Technical specification for building VacayTracker with Svelte 5, SvelteKit, Melt UI, and Tailwind CSS v4.

---

## Table of Contents

1. [Technology Stack](#1-technology-stack)
2. [Project Setup](#2-project-setup)
3. [Svelte 5 Fundamentals](#3-svelte-5-fundamentals)
4. [Tailwind CSS v4 Configuration](#4-tailwind-css-v4-configuration)
5. [Melt UI Component Library](#5-melt-ui-component-library)
6. [Project Structure](#6-project-structure)
7. [Component Architecture](#7-component-architecture)
8. [State Management](#8-state-management)
9. [Routing & Authentication](#9-routing--authentication)
10. [API Integration](#10-api-integration)
11. [VacayTracker Component Implementations](#11-vacaytracker-component-implementations)
12. [Development Workflow](#12-development-workflow)

---

## 1. Technology Stack

### 1.1 Core Technologies

| Technology | Version | Purpose |
|------------|---------|---------|
| **Svelte** | 5.x | UI framework with runes-based reactivity |
| **SvelteKit** | 2.49.x | Full-stack framework with routing, SSR, API routes |
| **Melt UI** | 0.42.x (`melt`) | Headless, accessible component builders |
| **Tailwind CSS** | 4.x | Utility-first CSS framework |
| **TypeScript** | 5.x | Type safety |
| **Vite** | 6.x | Build tool and dev server |

### 1.2 Why This Stack?

**Svelte 5** — Ground-up rewrite released October 2024. Introduces "runes" for explicit, fine-grained reactivity. Compiles to minimal JavaScript with no virtual DOM overhead. Ideal for a small team wanting maximum performance with minimal boilerplate.

**Melt UI** — The `melt` package is built specifically for Svelte 5 runes. Provides headless, accessible component primitives (dialogs, tabs, calendars, etc.) that you style yourself. WAI-ARIA compliant out of the box.

**Tailwind CSS v4** — Major release with CSS-native configuration via `@theme`, simplified Vite plugin setup, and improved performance. No more `tailwind.config.js` required.

### 1.3 Package Versions

```json
{
  "devDependencies": {
    "@sveltejs/adapter-node": "^5.2.0",
    "@sveltejs/kit": "^2.49.0",
    "@sveltejs/vite-plugin-svelte": "^5.0.0",
    "@tailwindcss/vite": "^4.0.0",
    "melt": "^0.42.0",
    "svelte": "^5.0.0",
    "tailwindcss": "^4.0.0",
    "typescript": "^5.7.0",
    "vite": "^6.0.0"
  }
}
```

---

## 2. Project Setup

### 2.1 Create Project

```bash
# Create new SvelteKit project using the sv CLI
npx sv create vacaytracker

# Prompts:
# ┌ Welcome to the Svelte CLI!
# │
# ◇ Which template would you like?
# │ SvelteKit minimal
# │
# ◇ Add type checking with TypeScript?
# │ Yes, using TypeScript syntax
# │
# ◇ What would you like to add to your project?
# │ prettier, eslint, tailwindcss
# │
# ◇ Which package manager do you want to install dependencies with?
# │ npm

cd vacaytracker
```

### 2.2 Install Additional Dependencies

```bash
# Melt UI (Svelte 5 version)
npm install melt

# Optional: icons
npm install -D lucide-svelte
```

### 2.3 Vite Configuration

```typescript
// vite.config.ts
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit()
  ]
});
```

### 2.4 TypeScript Configuration

```json
// tsconfig.json
{
  "extends": "./.svelte-kit/tsconfig.json",
  "compilerOptions": {
    "strict": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "paths": {
      "$lib/*": ["./src/lib/*"]
    }
  }
}
```

---

## 3. Svelte 5 Fundamentals

### 3.1 Runes Overview

Runes are compiler instructions that explicitly declare reactivity. They replace Svelte 4's implicit reactivity.

| Rune | Purpose | Example |
|------|---------|---------|
| `$state()` | Declare reactive state | `let count = $state(0)` |
| `$derived()` | Computed values | `let doubled = $derived(count * 2)` |
| `$effect()` | Side effects | `$effect(() => console.log(count))` |
| `$props()` | Component props | `let { name } = $props()` |
| `$bindable()` | Two-way bindable props | `let { value = $bindable() } = $props()` |

### 3.2 Reactive State

```svelte
<script lang="ts">
  // Primitive state
  let count = $state(0);
  
  // Object state (deeply reactive)
  let user = $state({
    name: 'John',
    vacationDays: 25
  });
  
  // Array state
  let requests = $state<VacationRequest[]>([]);
  
  function increment() {
    count += 1;  // Direct mutation works
  }
  
  function addRequest(request: VacationRequest) {
    requests.push(request);  // Array methods work
  }
</script>
```

### 3.3 Derived Values

```svelte
<script lang="ts">
  let total = $state(25);
  let used = $state(5);
  
  // Simple derivation
  let remaining = $derived(total - used);
  let percentage = $derived(Math.round((remaining / total) * 100));
  
  // Complex derivation with $derived.by()
  let status = $derived.by(() => {
    if (percentage >= 80) return { emoji: '🌴', text: 'Fully charged!' };
    if (percentage >= 50) return { emoji: '😎', text: 'Looking good!' };
    if (percentage >= 20) return { emoji: '😅', text: 'Running low!' };
    return { emoji: '😰', text: 'Almost out!' };
  });
</script>
```

### 3.4 Side Effects

```svelte
<script lang="ts">
  let searchQuery = $state('');
  let results = $state([]);
  
  // Runs when dependencies change (after DOM update)
  $effect(() => {
    if (searchQuery.length < 2) {
      results = [];
      return;
    }
    
    // Fetch results
    fetchResults(searchQuery).then(data => {
      results = data;
    });
  });
  
  // With cleanup
  $effect(() => {
    const interval = setInterval(() => {
      console.log('tick');
    }, 1000);
    
    return () => clearInterval(interval);  // Cleanup
  });
</script>
```

### 3.5 Component Props

```svelte
<!-- Button.svelte -->
<script lang="ts">
  import type { Snippet } from 'svelte';
  
  interface Props {
    variant?: 'primary' | 'secondary' | 'danger';
    disabled?: boolean;
    onclick?: () => void;
    children: Snippet;
  }
  
  let { 
    variant = 'primary', 
    disabled = false,
    onclick,
    children 
  }: Props = $props();
</script>

<button {disabled} {onclick} class="btn btn-{variant}">
  {@render children()}
</button>
```

### 3.6 Event Handling

```svelte
<script lang="ts">
  function handleClick(event: MouseEvent) {
    console.log('Clicked!', event);
  }
  
  function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    // Handle form
  }
</script>

<!-- Svelte 5 uses standard DOM event attributes -->
<button onclick={handleClick}>Click me</button>
<button onclick={() => console.log('inline')}>Inline</button>

<form onsubmit={handleSubmit}>
  <button type="submit">Submit</button>
</form>
```

### 3.7 Snippets (Replacing Slots)

```svelte
<!-- Card.svelte -->
<script lang="ts">
  import type { Snippet } from 'svelte';
  
  interface Props {
    title?: string;
    header?: Snippet;
    children: Snippet;
    footer?: Snippet;
  }
  
  let { title, header, children, footer }: Props = $props();
</script>

<div class="card">
  {#if header}
    <div class="card-header">{@render header()}</div>
  {:else if title}
    <div class="card-header"><h3>{title}</h3></div>
  {/if}
  
  <div class="card-body">
    {@render children()}
  </div>
  
  {#if footer}
    <div class="card-footer">{@render footer()}</div>
  {/if}
</div>

<!-- Usage -->
<Card title="Vacation Balance">
  <p>You have 20 days remaining</p>
  
  {#snippet footer()}
    <button>Request Time Off</button>
  {/snippet}
</Card>
```

### 3.8 Runes in TypeScript Files

When using runes outside `.svelte` files, use `.svelte.ts` extension:

```typescript
// src/lib/stores/counter.svelte.ts
export function createCounter(initial = 0) {
  let count = $state(initial);
  
  return {
    get count() { return count; },
    increment: () => count++,
    decrement: () => count--,
    reset: () => count = initial
  };
}
```

---

## 4. Tailwind CSS v4 Configuration

### 4.1 Setup (Already done by sv CLI)

Tailwind v4 uses a Vite plugin instead of PostCSS:

```typescript
// vite.config.ts
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  plugins: [tailwindcss(), sveltekit()]
});
```

### 4.2 Global Styles

```css
/* src/app.css */
@import "tailwindcss";

/* Custom theme using CSS @theme (Tailwind v4) */
@theme {
  /* Brand Colors - Beach/Vacation Theme */
  --color-ocean-50: oklch(0.97 0.02 220);
  --color-ocean-100: oklch(0.93 0.04 220);
  --color-ocean-200: oklch(0.86 0.08 220);
  --color-ocean-300: oklch(0.76 0.12 220);
  --color-ocean-400: oklch(0.66 0.16 220);
  --color-ocean-500: oklch(0.55 0.18 220);
  --color-ocean-600: oklch(0.48 0.16 220);
  --color-ocean-700: oklch(0.40 0.14 220);
  
  --color-sand-100: oklch(0.96 0.03 85);
  --color-sand-200: oklch(0.92 0.05 85);
  --color-sand-500: oklch(0.75 0.12 85);
  
  --color-coral-500: oklch(0.70 0.18 30);
  --color-palm-500: oklch(0.65 0.20 145);
  
  /* Semantic Colors */
  --color-success: oklch(0.65 0.20 145);
  --color-warning: oklch(0.75 0.15 65);
  --color-error: oklch(0.60 0.22 25);
  --color-pending: oklch(0.70 0.18 55);
  
  /* Typography */
  --font-sans: 'Inter', ui-sans-serif, system-ui, sans-serif;
  
  /* Spacing */
  --spacing-card: 1.5rem;
  
  /* Border Radius */
  --radius-sm: 0.375rem;
  --radius-md: 0.5rem;
  --radius-lg: 0.75rem;
  --radius-xl: 1rem;
}

/* Component layer for reusable styles */
@layer components {
  .btn {
    @apply px-4 py-2 rounded-lg font-medium transition-colors;
  }
  
  .btn-primary {
    @apply bg-ocean-500 text-white hover:bg-ocean-600;
  }
  
  .btn-secondary {
    @apply bg-gray-100 text-gray-700 hover:bg-gray-200;
  }
  
  .card {
    @apply bg-white rounded-xl shadow-sm p-6;
  }
  
  .input {
    @apply w-full px-3 py-2 border border-gray-300 rounded-lg 
           focus:outline-none focus:ring-2 focus:ring-ocean-500 focus:border-transparent;
  }
}
```

### 4.3 Layout Integration

```svelte
<!-- src/routes/+layout.svelte -->
<script lang="ts">
  import '../app.css';
  
  let { children } = $props();
</script>

{@render children()}
```

### 4.4 Tailwind v4 Key Changes

| Feature | Tailwind v3 | Tailwind v4 |
|---------|-------------|-------------|
| Config | `tailwind.config.js` | `@theme` in CSS |
| Plugin | PostCSS | Vite plugin |
| Import | `@tailwind base/components/utilities` | `@import "tailwindcss"` |
| Colors | JS config | CSS custom properties |

---

## 5. Melt UI Component Library

### 5.1 Installation

Use the `melt` package (not `@melt-ui/svelte`) for Svelte 5:

```bash
npm install melt
```

### 5.2 Available Components

| Category | Components |
|----------|------------|
| **Layout** | Accordion, Collapsible, Tabs, Toolbar, Tree |
| **Overlays** | Dialog, Popover, Tooltip, Dropdown Menu, Context Menu |
| **Form Controls** | Checkbox, Radio Group, Select, Combobox, Switch, Toggle, Slider, PIN Input, Tags Input |
| **Date/Time** | Calendar, Date Field, Date Picker, Date Range Picker, Range Calendar |
| **Display** | Avatar, Progress, Separator, Scroll Area, Toast, Pagination |

### 5.3 Usage Patterns

**Builder Pattern** (More control, direct class instantiation)

```svelte
<script lang="ts">
  import { Dialog } from 'melt/builders';
  
  let open = $state(false);
  
  const dialog = new Dialog({
    open: () => open,
    onOpenChange: (v) => (open = v)
  });
</script>

<button {...dialog.trigger}>Open Dialog</button>

{#if dialog.open}
  <div {...dialog.overlay} class="fixed inset-0 bg-black/50" />
  <div {...dialog.content} class="fixed inset-0 flex items-center justify-center">
    <div class="bg-white rounded-xl p-6 max-w-md">
      <h2 {...dialog.title}>Dialog Title</h2>
      <p {...dialog.description}>Dialog content goes here.</p>
      <button {...dialog.close}>Close</button>
    </div>
  </div>
{/if}
```

**Component Pattern** (Simpler, supports `bind:`)

```svelte
<script lang="ts">
  import { Toggle } from 'melt/components';
  
  let pressed = $state(false);
</script>

<Toggle bind:value={pressed}>
  {#snippet children(toggle)}
    <button {...toggle.trigger} class="p-2 rounded {toggle.value ? 'bg-ocean-500' : 'bg-gray-200'}">
      {toggle.value ? 'On' : 'Off'}
    </button>
  {/snippet}
</Toggle>
```

### 5.4 Core Component Examples

**Tabs**

```svelte
<script lang="ts">
  import { Tabs } from 'melt/builders';
  
  let value = $state('dashboard');
  
  const tabs = new Tabs({
    value: () => value,
    onValueChange: (v) => (value = v)
  });
</script>

<div {...tabs.root}>
  <div {...tabs.list} class="flex border-b">
    <button 
      {...tabs.trigger('dashboard')}
      class="px-4 py-2 border-b-2 {tabs.value === 'dashboard' ? 'border-ocean-500' : 'border-transparent'}"
    >
      Dashboard
    </button>
    <button 
      {...tabs.trigger('requests')}
      class="px-4 py-2 border-b-2 {tabs.value === 'requests' ? 'border-ocean-500' : 'border-transparent'}"
    >
      Requests
    </button>
  </div>
  
  <div {...tabs.content('dashboard')} class="p-4">
    Dashboard content
  </div>
  <div {...tabs.content('requests')} class="p-4">
    Requests content
  </div>
</div>
```

**Select**

```svelte
<script lang="ts">
  import { Select } from 'melt/builders';
  
  let value = $state('');
  
  const select = new Select({
    value: () => value,
    onValueChange: (v) => (value = v)
  });
  
  const options = [
    { value: 'family', label: 'Family Time' },
    { value: 'travel', label: 'Travel' },
    { value: 'staycation', label: 'Staycation' }
  ];
</script>

<div {...select.root}>
  <button {...select.trigger} class="input flex justify-between items-center">
    <span>{select.valueLabel || 'Select reason...'}</span>
    <span>▼</span>
  </button>
  
  {#if select.open}
    <div {...select.content} class="absolute mt-1 bg-white border rounded-lg shadow-lg z-10">
      {#each options as option}
        <div 
          {...select.option(option.value)}
          class="px-4 py-2 hover:bg-gray-100 cursor-pointer"
        >
          {option.label}
        </div>
      {/each}
    </div>
  {/if}
</div>
```

**Calendar**

```svelte
<script lang="ts">
  import { Calendar } from 'melt/builders';
  
  let value = $state<Date | null>(null);
  
  const calendar = new Calendar({
    value: () => value,
    onValueChange: (v) => (value = v)
  });
</script>

<div {...calendar.root} class="p-4 bg-white rounded-xl shadow">
  <div class="flex justify-between items-center mb-4">
    <button {...calendar.prevButton}>←</button>
    <span {...calendar.heading}>{calendar.headingValue}</span>
    <button {...calendar.nextButton}>→</button>
  </div>
  
  <div {...calendar.grid} class="grid grid-cols-7 gap-1">
    {#each calendar.weekdays as day}
      <div class="text-center text-sm text-gray-500 py-2">{day}</div>
    {/each}
    
    {#each calendar.days as day}
      <button
        {...calendar.day(day)}
        class="p-2 text-center rounded hover:bg-ocean-100
               {day.isToday ? 'font-bold' : ''}
               {day.isSelected ? 'bg-ocean-500 text-white' : ''}
               {day.isOutsideMonth ? 'text-gray-300' : ''}"
      >
        {day.day}
      </button>
    {/each}
  </div>
</div>
```

---

## 6. Project Structure

```
vacaytracker/
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
│   │   │   │   └── index.ts
│   │   │   ├── layout/                # Layout components
│   │   │   │   ├── Header.svelte
│   │   │   │   ├── Sidebar.svelte
│   │   │   │   ├── PageHeader.svelte
│   │   │   │   └── index.ts
│   │   │   └── features/              # Feature components
│   │   │       ├── auth/
│   │   │       │   └── LoginForm.svelte
│   │   │       ├── vacation/
│   │   │       │   ├── RequestModal.svelte
│   │   │       │   ├── RequestCard.svelte
│   │   │       │   ├── Timeline.svelte
│   │   │       │   ├── Calendar.svelte
│   │   │       │   └── BalanceDisplay.svelte
│   │   │       └── admin/
│   │   │           ├── UserCard.svelte
│   │   │           ├── StatsCard.svelte
│   │   │           └── ResetModal.svelte
│   │   ├── stores/                    # Global state
│   │   │   ├── auth.svelte.ts
│   │   │   └── toast.svelte.ts
│   │   ├── api/                       # API client
│   │   │   ├── client.ts
│   │   │   ├── auth.ts
│   │   │   ├── vacation.ts
│   │   │   └── users.ts
│   │   ├── types/                     # TypeScript types
│   │   │   └── index.ts
│   │   └── utils/                     # Utilities
│   │       ├── date.ts
│   │       └── format.ts
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
│   ├── app.css
│   ├── app.html
│   └── app.d.ts
├── static/
├── svelte.config.js
├── vite.config.ts
├── package.json
└── tsconfig.json
```

---

## 7. Component Architecture

### 7.1 Base UI Components

```svelte
<!-- src/lib/components/ui/Button.svelte -->
<script lang="ts">
  import type { Snippet } from 'svelte';
  import type { HTMLButtonAttributes } from 'svelte/elements';
  
  interface Props extends HTMLButtonAttributes {
    variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
    size?: 'sm' | 'md' | 'lg';
    loading?: boolean;
    children: Snippet;
  }
  
  let { 
    variant = 'primary', 
    size = 'md', 
    loading = false,
    children,
    class: className = '',
    disabled,
    ...rest 
  }: Props = $props();
  
  const variants: Record<string, string> = {
    primary: 'bg-ocean-500 text-white hover:bg-ocean-600',
    secondary: 'bg-gray-100 text-gray-700 hover:bg-gray-200',
    ghost: 'text-gray-600 hover:bg-gray-100',
    danger: 'bg-error text-white hover:bg-red-600'
  };
  
  const sizes: Record<string, string> = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2',
    lg: 'px-6 py-3 text-lg'
  };
</script>

<button 
  class="rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed
         {variants[variant]} {sizes[size]} {className}"
  disabled={disabled || loading}
  {...rest}
>
  {#if loading}
    <span class="inline-block animate-spin mr-2">⏳</span>
  {/if}
  {@render children()}
</button>
```

```svelte
<!-- src/lib/components/ui/Input.svelte -->
<script lang="ts">
  import type { HTMLInputAttributes } from 'svelte/elements';
  
  interface Props extends HTMLInputAttributes {
    label?: string;
    error?: string;
  }
  
  let { 
    label,
    error,
    id,
    class: className = '',
    ...rest 
  }: Props = $props();
  
  const inputId = id ?? crypto.randomUUID();
</script>

<div class="space-y-1">
  {#if label}
    <label for={inputId} class="block text-sm font-medium text-gray-700">
      {label}
    </label>
  {/if}
  
  <input
    id={inputId}
    class="w-full px-3 py-2 border rounded-lg transition-colors
           focus:outline-none focus:ring-2 focus:ring-ocean-500 focus:border-transparent
           {error ? 'border-error' : 'border-gray-300'}
           {className}"
    {...rest}
  />
  
  {#if error}
    <p class="text-sm text-error">{error}</p>
  {/if}
</div>
```

```svelte
<!-- src/lib/components/ui/Card.svelte -->
<script lang="ts">
  import type { Snippet } from 'svelte';
  
  interface Props {
    title?: string;
    padding?: 'none' | 'sm' | 'md' | 'lg';
    header?: Snippet;
    children: Snippet;
    footer?: Snippet;
  }
  
  let { title, padding = 'md', header, children, footer }: Props = $props();
  
  const paddings: Record<string, string> = {
    none: '',
    sm: 'p-4',
    md: 'p-6',
    lg: 'p-8'
  };
</script>

<div class="bg-white rounded-xl shadow-sm {paddings[padding]}">
  {#if header}
    <div class="border-b pb-4 mb-4">{@render header()}</div>
  {:else if title}
    <h3 class="text-lg font-semibold text-gray-800 mb-4">{title}</h3>
  {/if}
  
  {@render children()}
  
  {#if footer}
    <div class="border-t pt-4 mt-4">{@render footer()}</div>
  {/if}
</div>
```

### 7.2 Barrel Exports

```typescript
// src/lib/components/ui/index.ts
export { default as Button } from './Button.svelte';
export { default as Input } from './Input.svelte';
export { default as Card } from './Card.svelte';
export { default as Badge } from './Badge.svelte';
export { default as Avatar } from './Avatar.svelte';
export { default as ProgressRing } from './ProgressRing.svelte';
```

---

## 8. State Management

### 8.1 Auth Store

```typescript
// src/lib/stores/auth.svelte.ts
import type { User } from '$lib/types';

function createAuthStore() {
  let user = $state<User | null>(null);
  let token = $state<string | null>(null);
  
  return {
    get user() { return user; },
    get token() { return token; },
    get isAuthenticated() { return !!token && !!user; },
    get isAdmin() { return user?.role === 'admin'; },
    get isEmployee() { return user?.role === 'employee'; },
    
    setSession(newUser: User, newToken: string) {
      user = newUser;
      token = newToken;
      if (typeof sessionStorage !== 'undefined') {
        sessionStorage.setItem('token', newToken);
        sessionStorage.setItem('user', JSON.stringify(newUser));
      }
    },
    
    logout() {
      user = null;
      token = null;
      if (typeof sessionStorage !== 'undefined') {
        sessionStorage.removeItem('token');
        sessionStorage.removeItem('user');
      }
    },
    
    init() {
      if (typeof sessionStorage !== 'undefined') {
        token = sessionStorage.getItem('token');
        const storedUser = sessionStorage.getItem('user');
        if (storedUser) {
          user = JSON.parse(storedUser);
        }
      }
    }
  };
}

export const auth = createAuthStore();
```

### 8.2 Toast Store

```typescript
// src/lib/stores/toast.svelte.ts
export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface Toast {
  id: string;
  type: ToastType;
  title: string;
  message?: string;
}

function createToastStore() {
  let toasts = $state<Toast[]>([]);
  
  function add(type: ToastType, title: string, message?: string, duration = 5000) {
    const id = crypto.randomUUID();
    toasts.push({ id, type, title, message });
    
    setTimeout(() => remove(id), duration);
    return id;
  }
  
  function remove(id: string) {
    toasts = toasts.filter(t => t.id !== id);
  }
  
  return {
    get list() { return toasts; },
    
    success: (title: string, message?: string) => add('success', title, message),
    error: (title: string, message?: string) => add('error', title, message),
    warning: (title: string, message?: string) => add('warning', title, message),
    info: (title: string, message?: string) => add('info', title, message),
    
    remove
  };
}

export const toast = createToastStore();
```

---

## 9. Routing & Authentication

### 9.1 Route Guards

```typescript
// src/routes/employee/+layout.server.ts
import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies, fetch }) => {
  const token = cookies.get('token');
  
  if (!token) {
    throw redirect(303, '/?error=unauthorized');
  }
  
  try {
    const response = await fetch('/api/me', {
      headers: { Authorization: `Bearer ${token}` }
    });
    
    if (!response.ok) {
      cookies.delete('token', { path: '/' });
      throw redirect(303, '/?error=session_expired');
    }
    
    const user = await response.json();
    
    if (user.role !== 'employee') {
      throw redirect(303, '/admin');
    }
    
    return { user };
  } catch (error) {
    cookies.delete('token', { path: '/' });
    throw redirect(303, '/');
  }
};
```

```typescript
// src/routes/admin/+layout.server.ts
import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies, fetch }) => {
  const token = cookies.get('token');
  
  if (!token) {
    throw redirect(303, '/?error=unauthorized');
  }
  
  const response = await fetch('/api/me', {
    headers: { Authorization: `Bearer ${token}` }
  });
  
  if (!response.ok) {
    cookies.delete('token', { path: '/' });
    throw redirect(303, '/');
  }
  
  const user = await response.json();
  
  if (user.role !== 'admin') {
    throw redirect(303, '/employee');
  }
  
  return { user };
};
```

### 9.2 Login Page

```svelte
<!-- src/routes/+page.svelte -->
<script lang="ts">
  import { goto } from '$app/navigation';
  import { Tabs } from 'melt/builders';
  import { Button, Input, Card } from '$lib/components/ui';
  import { toast } from '$lib/stores/toast.svelte';
  
  let role = $state<'employee' | 'admin'>('employee');
  let username = $state('');
  let password = $state('');
  let loading = $state(false);
  let error = $state('');
  
  const tabs = new Tabs({
    value: () => role,
    onValueChange: (v) => (role = v as 'employee' | 'admin')
  });
  
  async function handleSubmit() {
    loading = true;
    error = '';
    
    try {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password, role })
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        error = data.error?.message || 'Login failed';
        return;
      }
      
      toast.success('Welcome back!', `Logged in as ${data.user.name}`);
      goto(role === 'admin' ? '/admin' : '/employee');
    } catch (e) {
      error = 'Connection error. Please try again.';
    } finally {
      loading = false;
    }
  }
</script>

<div class="min-h-screen bg-gradient-to-b from-ocean-100 to-sand-100 flex items-center justify-center p-4">
  <Card class="w-full max-w-md">
    <div class="text-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">🏖️ VacayTracker</h1>
      <p class="text-gray-600">Track your time off with ease</p>
    </div>
    
    <div {...tabs.root}>
      <div {...tabs.list} class="flex mb-6 bg-gray-100 rounded-lg p-1">
        <button 
          {...tabs.trigger('employee')}
          class="flex-1 py-2 rounded-md transition-colors
                 {tabs.value === 'employee' ? 'bg-white shadow' : ''}"
        >
          🧑‍💼 Employee
        </button>
        <button 
          {...tabs.trigger('admin')}
          class="flex-1 py-2 rounded-md transition-colors
                 {tabs.value === 'admin' ? 'bg-white shadow' : ''}"
        >
          👨‍✈️ Captain
        </button>
      </div>
    </div>
    
    <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
      {#if error}
        <div class="bg-red-50 text-error p-3 rounded-lg text-sm">{error}</div>
      {/if}
      
      <Input 
        label="Username"
        type="text" 
        bind:value={username}
        required
        autocomplete="username"
      />
      
      <Input 
        label="Password"
        type="password" 
        bind:value={password}
        required
        autocomplete="current-password"
      />
      
      <Button type="submit" {loading} class="w-full">
        {loading ? 'Signing in...' : 'Sign In'}
      </Button>
    </form>
  </Card>
</div>
```

---

## 10. API Integration

### 10.1 API Client

```typescript
// src/lib/api/client.ts
interface ApiError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}

class ApiClient {
  private baseUrl = '/api';
  
  private getToken(): string | null {
    if (typeof document === 'undefined') return null;
    return document.cookie
      .split('; ')
      .find(row => row.startsWith('token='))
      ?.split('=')[1] ?? null;
  }
  
  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const token = this.getToken();
    
    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers
      }
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw data.error as ApiError;
    }
    
    return data;
  }
  
  get<T>(endpoint: string) {
    return this.request<T>(endpoint);
  }
  
  post<T>(endpoint: string, body: unknown) {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(body)
    });
  }
  
  put<T>(endpoint: string, body: unknown) {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(body)
    });
  }
  
  delete<T>(endpoint: string) {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

export const api = new ApiClient();
```

### 10.2 API Modules

```typescript
// src/lib/api/vacation.ts
import { api } from './client';
import type { VacationRequest } from '$lib/types';

export const vacationApi = {
  getRequests: () => 
    api.get<{ data: VacationRequest[] }>('/vacation-requests'),
  
  create: (data: { startDate: string; endDate: string; reason?: string }) =>
    api.post<{ data: VacationRequest }>('/vacation-requests', data),
  
  updateStatus: (id: string, status: 'approved' | 'rejected') =>
    api.put<{ data: VacationRequest }>(`/vacation-requests/${id}`, { status }),
  
  getTeamVacations: () =>
    api.get<{ data: VacationRequest[] }>('/team-vacation'),
  
  getBalance: (userId: string) =>
    api.get<{ total: number; used: number }>(`/vacation-days/${userId}`)
};
```

```typescript
// src/lib/api/users.ts
import { api } from './client';
import type { User } from '$lib/types';

export const usersApi = {
  getAll: () => 
    api.get<{ data: User[] }>('/users'),
  
  create: (data: Omit<User, 'id'>) =>
    api.post<{ data: User }>('/users', data),
  
  update: (id: string, data: Partial<User>) =>
    api.put<{ data: User }>(`/users/${id}`, data),
  
  delete: (id: string) =>
    api.delete<void>(`/users/${id}`),
  
  resetVacationDays: (data: { vacationDays: number; clearHistory?: boolean }) =>
    api.post<void>('/admin/reset-vacation-days', data)
};
```

---

## 11. VacayTracker Component Implementations

### 11.1 Vacation Balance Display

```svelte
<!-- src/lib/components/features/vacation/BalanceDisplay.svelte -->
<script lang="ts">
  interface Props {
    total: number;
    used: number;
  }
  
  let { total, used }: Props = $props();
  
  let remaining = $derived(total - used);
  let percentage = $derived(Math.round((remaining / total) * 100));
  
  // SVG circle calculations
  const radius = 45;
  const circumference = $derived(2 * Math.PI * radius);
  let strokeOffset = $derived(circumference - (percentage / 100) * circumference);
  
  let mood = $derived.by(() => {
    if (percentage >= 80) return { emoji: '🌴', text: "You're fully charged!", color: 'text-palm-500' };
    if (percentage >= 60) return { emoji: '😎', text: 'Looking good!', color: 'text-ocean-500' };
    if (percentage >= 40) return { emoji: '🤔', text: 'Time to start planning!', color: 'text-warning' };
    if (percentage >= 20) return { emoji: '😅', text: 'Running low!', color: 'text-coral-500' };
    return { emoji: '😰', text: 'Almost out!', color: 'text-error' };
  });
</script>

<div class="card">
  <h3 class="text-lg font-semibold text-gray-700 mb-4">Vacation Balance</h3>
  
  <div class="flex items-center gap-6">
    <!-- Progress Ring -->
    <div class="relative">
      <svg width="120" height="120" class="-rotate-90">
        <!-- Background circle -->
        <circle
          cx="60" cy="60" r={radius}
          fill="none"
          stroke="#e5e7eb"
          stroke-width="10"
        />
        <!-- Progress circle -->
        <circle
          cx="60" cy="60" r={radius}
          fill="none"
          stroke="currentColor"
          stroke-width="10"
          stroke-linecap="round"
          stroke-dasharray={circumference}
          stroke-dashoffset={strokeOffset}
          class="text-ocean-500 transition-all duration-700 ease-out"
        />
      </svg>
      <div class="absolute inset-0 flex flex-col items-center justify-center rotate-0">
        <span class="text-3xl font-bold text-gray-800">{remaining}</span>
        <span class="text-sm text-gray-500">days left</span>
      </div>
    </div>
    
    <!-- Status -->
    <div class="space-y-2">
      <div class="flex items-center gap-2">
        <span class="text-3xl">{mood.emoji}</span>
        <span class="{mood.color} font-medium">{mood.text}</span>
      </div>
      <div class="text-sm text-gray-500 space-y-1">
        <p>{used} of {total} days used</p>
        <p>{percentage}% remaining</p>
      </div>
    </div>
  </div>
</div>
```

### 11.2 Vacation Request Modal

```svelte
<!-- src/lib/components/features/vacation/RequestModal.svelte -->
<script lang="ts">
  import { Dialog } from 'melt/builders';
  import { Calendar } from 'melt/builders';
  import { Button, Input } from '$lib/components/ui';
  import { vacationApi } from '$lib/api/vacation';
  import { toast } from '$lib/stores/toast.svelte';
  
  interface Props {
    onSuccess?: () => void;
  }
  
  let { onSuccess }: Props = $props();
  
  let open = $state(false);
  let startDate = $state<Date | null>(null);
  let endDate = $state<Date | null>(null);
  let reason = $state('');
  let loading = $state(false);
  
  const dialog = new Dialog({
    open: () => open,
    onOpenChange: (v) => { 
      open = v;
      if (!v) reset();
    }
  });
  
  const reasonSuggestions = ['Family Time', 'Travel', 'Staycation', 'Personal'];
  
  let businessDays = $derived.by(() => {
    if (!startDate || !endDate) return 0;
    // Simplified business day calculation
    let count = 0;
    const current = new Date(startDate);
    while (current <= endDate) {
      const day = current.getDay();
      if (day !== 0 && day !== 6) count++;
      current.setDate(current.getDate() + 1);
    }
    return count;
  });
  
  function reset() {
    startDate = null;
    endDate = null;
    reason = '';
  }
  
  async function handleSubmit() {
    if (!startDate || !endDate) return;
    
    loading = true;
    try {
      await vacationApi.create({
        startDate: startDate.toISOString().split('T')[0],
        endDate: endDate.toISOString().split('T')[0],
        reason: reason || undefined
      });
      
      toast.success('Request submitted!', 'Your vacation request is pending approval.');
      open = false;
      onSuccess?.();
    } catch (error: any) {
      toast.error('Failed to submit', error.message);
    } finally {
      loading = false;
    }
  }
</script>

<Button onclick={() => (open = true)}>🏖️ Request Time Off</Button>

{#if dialog.open}
  <div {...dialog.overlay} class="fixed inset-0 bg-black/50 z-40" />
  
  <div {...dialog.content} class="fixed inset-0 flex items-center justify-center z-50 p-4">
    <div class="bg-white rounded-xl shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
      <div class="p-6">
        <h2 {...dialog.title} class="text-xl font-semibold mb-1">Request Vacation</h2>
        <p {...dialog.description} class="text-gray-600 text-sm mb-6">
          Select your dates and submit for approval
        </p>
        
        <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-6">
          <!-- Date Selection -->
          <div class="grid grid-cols-2 gap-4">
            <Input 
              label="Start Date"
              type="date"
              value={startDate?.toISOString().split('T')[0] ?? ''}
              onchange={(e) => (startDate = new Date(e.currentTarget.value))}
              min={new Date().toISOString().split('T')[0]}
              required
            />
            <Input 
              label="End Date"
              type="date"
              value={endDate?.toISOString().split('T')[0] ?? ''}
              onchange={(e) => (endDate = new Date(e.currentTarget.value))}
              min={startDate?.toISOString().split('T')[0] ?? new Date().toISOString().split('T')[0]}
              required
            />
          </div>
          
          {#if businessDays > 0}
            <div class="bg-ocean-50 text-ocean-700 p-3 rounded-lg text-sm">
              📅 {businessDays} business day{businessDays !== 1 ? 's' : ''}
            </div>
          {/if}
          
          <!-- Reason -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Reason (optional)
            </label>
            <div class="flex flex-wrap gap-2 mb-2">
              {#each reasonSuggestions as suggestion}
                <button
                  type="button"
                  onclick={() => (reason = suggestion)}
                  class="px-3 py-1 text-sm rounded-full border transition-colors
                         {reason === suggestion 
                           ? 'bg-ocean-500 text-white border-ocean-500' 
                           : 'border-gray-300 hover:border-ocean-300'}"
                >
                  {suggestion}
                </button>
              {/each}
            </div>
            <textarea
              bind:value={reason}
              maxlength={200}
              rows={2}
              class="input resize-none"
              placeholder="Any additional notes..."
            />
            <div class="text-right text-xs text-gray-500 mt-1">
              {reason.length}/200
            </div>
          </div>
          
          <!-- Actions -->
          <div class="flex justify-end gap-3 pt-4 border-t">
            <Button type="button" variant="ghost" onclick={() => (open = false)}>
              Cancel
            </Button>
            <Button type="submit" {loading} disabled={!startDate || !endDate}>
              Submit Request
            </Button>
          </div>
        </form>
      </div>
    </div>
  </div>
{/if}
```

### 11.3 Toast Provider

```svelte
<!-- src/lib/components/ui/ToastProvider.svelte -->
<script lang="ts">
  import { toast } from '$lib/stores/toast.svelte';
  
  const icons: Record<string, string> = {
    success: '✅',
    error: '❌',
    warning: '⚠️',
    info: 'ℹ️'
  };
  
  const colors: Record<string, string> = {
    success: 'border-l-success',
    error: 'border-l-error',
    warning: 'border-l-warning',
    info: 'border-l-ocean-500'
  };
</script>

<div class="fixed bottom-4 right-4 z-50 space-y-2 max-w-sm">
  {#each toast.list as t (t.id)}
    <div 
      class="bg-white rounded-lg shadow-lg p-4 border-l-4 {colors[t.type]}
             animate-in slide-in-from-right duration-300"
    >
      <div class="flex gap-3">
        <span class="text-lg">{icons[t.type]}</span>
        <div class="flex-1">
          <p class="font-medium text-gray-800">{t.title}</p>
          {#if t.message}
            <p class="text-sm text-gray-600 mt-1">{t.message}</p>
          {/if}
        </div>
        <button 
          onclick={() => toast.remove(t.id)}
          class="text-gray-400 hover:text-gray-600"
        >
          ✕
        </button>
      </div>
    </div>
  {/each}
</div>
```

### 11.4 Admin Stats Card

```svelte
<!-- src/lib/components/features/admin/StatsCard.svelte -->
<script lang="ts">
  interface Props {
    icon: string;
    label: string;
    value: number | string;
    trend?: { direction: 'up' | 'down'; value: string };
    color?: 'ocean' | 'palm' | 'coral' | 'sand';
  }
  
  let { icon, label, value, trend, color = 'ocean' }: Props = $props();
  
  const bgColors: Record<string, string> = {
    ocean: 'bg-ocean-50',
    palm: 'bg-green-50',
    coral: 'bg-orange-50',
    sand: 'bg-amber-50'
  };
  
  const iconColors: Record<string, string> = {
    ocean: 'text-ocean-500',
    palm: 'text-palm-500',
    coral: 'text-coral-500',
    sand: 'text-sand-500'
  };
</script>

<div class="card">
  <div class="flex items-start justify-between">
    <div class="{bgColors[color]} p-3 rounded-lg">
      <span class="text-2xl {iconColors[color]}">{icon}</span>
    </div>
    
    {#if trend}
      <span class="text-sm {trend.direction === 'up' ? 'text-palm-500' : 'text-error'}">
        {trend.direction === 'up' ? '↑' : '↓'} {trend.value}
      </span>
    {/if}
  </div>
  
  <div class="mt-4">
    <p class="text-2xl font-bold text-gray-800">{value}</p>
    <p class="text-sm text-gray-500">{label}</p>
  </div>
</div>
```

---

## 12. Development Workflow

### 12.1 Commands

```bash
# Development
npm run dev              # Start dev server at localhost:5173

# Type checking
npm run check            # Run svelte-check
npm run check:watch      # Watch mode

# Building
npm run build            # Production build
npm run preview          # Preview production build

# Code quality
npm run lint             # Run ESLint
npm run format           # Run Prettier
```

### 12.2 VS Code Setup

**Extensions**
- Svelte for VS Code (`svelte.svelte-vscode`)
- Tailwind CSS IntelliSense (`bradlc.vscode-tailwindcss`) — use pre-release for v4
- ESLint (`dbaeumer.vscode-eslint`)
- Prettier (`esbenp.prettier-vscode`)

**Settings** (`.vscode/settings.json`)

```json
{
  "svelte.enable-ts-plugin": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "[svelte]": {
    "editor.defaultFormatter": "svelte.svelte-vscode"
  },
  "tailwindCSS.experimental.configFile": null,
  "tailwindCSS.includeLanguages": {
    "svelte": "html"
  },
  "editor.quickSuggestions": {
    "strings": true
  }
}
```

### 12.3 Type Definitions

```typescript
// src/lib/types/index.ts
export interface User {
  id: string;
  name: string;
  username: string;
  email?: string;
  role: 'admin' | 'employee';
  vacationDays?: number;
  usedVacationDays?: number;
}

export interface VacationRequest {
  id: string;
  userId: string;
  startDate: string;
  endDate: string;
  businessDays: number;
  status: 'pending' | 'approved' | 'rejected';
  reason?: string;
  createdAt: string;
  reviewedBy?: string;
  reviewedAt?: string;
}

export interface Settings {
  weekendPolicy: {
    excludeWeekends: boolean;
  };
  newsletter: {
    enabled: boolean;
    dayOfMonth: number;
    hourOfDay: number;
  };
}
```

---

## Resources

### Official Documentation
- [Svelte 5 Documentation](https://svelte.dev/docs/svelte)
- [SvelteKit Documentation](https://svelte.dev/docs/kit)
- [Melt UI (Svelte 5)](https://next.melt-ui.com)
- [Tailwind CSS v4](https://tailwindcss.com/docs)

### Package References
- [Svelte npm](https://www.npmjs.com/package/svelte) — v5.x
- [SvelteKit npm](https://www.npmjs.com/package/@sveltejs/kit) — v2.49.x
- [Melt npm](https://www.npmjs.com/package/melt) — v0.42.x
- [sv CLI](https://www.npmjs.com/package/sv) — v0.10.x

### Community
- [Svelte Discord](https://svelte.dev/chat)
- [Melt UI Discord](https://melt-ui.com/discord)

---

*VacayTracker Frontend Implementation Guide — December 2025*