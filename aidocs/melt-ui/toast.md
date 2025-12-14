# Toast

A notification system for displaying temporary messages with auto-dismiss.

## Use Cases

- Success notifications (e.g., "Settings saved")
- Error messages (e.g., "Failed to save changes")
- Informational alerts (e.g., "Session expiring soon")
- Action confirmations with optional descriptions

## VacayTracker Implementation

VacayTracker uses a Svelte 5 rune-based toast store with a clean API supporting both simple messages and title+description format.

### Toast Store (`src/lib/stores/toast.svelte.ts`)

```typescript
export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface ToastData {
  id: string;
  type: ToastType;
  title: string;
  description?: string;
}

function createToastStore() {
  let toasts = $state<ToastData[]>([]);

  function add(type, title, descOrDuration?, duration?): string { ... }
  function dismiss(id: string): void { ... }
  function dismissAll(): void { ... }

  return {
    get toasts() { return toasts; },
    add, dismiss, dismissAll,
    success, error, warning, info  // Convenience methods
  };
}

export const toast = createToastStore();
```

### Usage

```typescript
import { toast } from '$lib/stores/toast.svelte';

// Simple message
toast.success('Settings saved');
toast.error('Failed to save');

// With description
toast.success('Request Approved', 'John\'s vacation from Dec 20-25 has been approved.');

// With custom duration (milliseconds)
toast.error('Connection lost', 10000);  // 10 seconds

// With description and custom duration
toast.info('New update available', 'Click to refresh the page.', 8000);
```

### Toaster Component (`src/lib/components/ui/Toaster.svelte`)

```svelte
<script lang="ts">
  import { toast } from '$lib/stores/toast.svelte';
  import { X, CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-svelte';

  const icons = {
    success: CheckCircle,
    error: AlertCircle,
    warning: AlertTriangle,
    info: Info
  };

  const styles = {
    success: { container: 'bg-success/10 border-success', icon: 'text-success' },
    error: { container: 'bg-error/10 border-error', icon: 'text-error' },
    warning: { container: 'bg-warning/10 border-warning', icon: 'text-warning' },
    info: { container: 'bg-info/10 border-info', icon: 'text-info' }
  };
</script>

<div class="fixed bottom-4 left-1/2 -translate-x-1/2 z-[9999] ..." role="region" aria-label="Notifications">
  {#each toast.toasts as t (t.id)}
    <div role="alert" aria-live="polite" class="... {styles[t.type].container}">
      <Icon class="{styles[t.type].icon}" />
      <p class="font-semibold">{t.title}</p>
      {#if t.description}<p>{t.description}</p>{/if}
      <button onclick={() => toast.dismiss(t.id)}>×</button>
    </div>
  {/each}
</div>
```

## API Reference

### Toast Store Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `success` | `(title, descOrDuration?, duration?)` | Show success toast |
| `error` | `(title, descOrDuration?, duration?)` | Show error toast |
| `warning` | `(title, descOrDuration?, duration?)` | Show warning toast |
| `info` | `(title, descOrDuration?, duration?)` | Show info toast |
| `dismiss` | `(id: string)` | Dismiss specific toast |
| `dismissAll` | `()` | Dismiss all toasts |

### Toast Properties

| Property | Type | Description |
|----------|------|-------------|
| `id` | `string` | Unique identifier (auto-generated) |
| `type` | `ToastType` | 'success' \| 'error' \| 'warning' \| 'info' |
| `title` | `string` | Main message text |
| `description` | `string?` | Optional secondary text |

### Default Behavior

- **Auto-dismiss**: 5 seconds (5000ms)
- **Position**: Bottom center
- **Stacking**: Newest at bottom (reverse column)
- **Z-index**: 9999

## Accessibility

- `role="region"` on container with `aria-label="Notifications"`
- `role="alert"` and `aria-live="polite"` on each toast
- `aria-hidden="true"` on decorative icons
- Accessible dismiss button with `aria-label`

## Styling

### Color Scheme

Uses semantic colors from `app.css`:
- **Success**: `--color-success` (green) - `oklch(0.65 0.2 145)`
- **Error**: `--color-error` (red) - `oklch(0.6 0.22 25)`
- **Warning**: `--color-warning` (amber) - `oklch(0.75 0.18 85)`
- **Info**: `--color-info` (blue) - `oklch(0.65 0.15 250)`

### Toast Animations (`app.css`)

```css
@keyframes toastIn {
  from { opacity: 0; transform: translateY(100%) scale(0.95); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

@keyframes toastOut {
  from { opacity: 1; transform: translateY(0) scale(1); }
  to { opacity: 0; transform: translateY(100%) scale(0.95); }
}

.animate-toast-in {
  animation: toastIn 0.3s ease-out forwards;
}
```

## Layout Integration

The Toaster component is mounted globally in `+layout.svelte`:

```svelte
<script lang="ts">
  import Toaster from '$lib/components/ui/Toaster.svelte';
</script>

<Toaster />
{@render children()}
```

## Positioning Options

Default is bottom-center. To change position, modify the container classes:

```svelte
<!-- Top right -->
<div class="fixed top-4 right-4 z-[9999] flex flex-col gap-2 ...">

<!-- Top center -->
<div class="fixed top-4 left-1/2 -translate-x-1/2 z-[9999] flex flex-col gap-2 ...">

<!-- Bottom right -->
<div class="fixed bottom-4 right-4 z-[9999] flex flex-col-reverse gap-2 ...">
```
