# Toast

A notification system for displaying temporary messages.

## Use Cases

- Success notifications
- Error messages
- Informational alerts
- Action confirmations

## Installation

```typescript
import { createToaster } from '@melt-ui/svelte';
```

## API Reference

### createToaster Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `closeDelay` | `number` | `5000` | Auto-close delay (ms) |
| `type` | `'foreground' \| 'background'` | `'foreground'` | Toast priority |
| `hover` | `'pause' \| 'pause-all' \| null` | `'pause'` | Hover behavior |

### Returned Elements

| Element | Description |
|---------|-------------|
| `content` | Toast container |
| `title` | Toast title |
| `description` | Toast description |
| `close` | Close button |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `toasts` | `Readable<Toast[]>` | Active toasts |

### Returned Helpers

| Helper | Description |
|--------|-------------|
| `addToast(data)` | Add a new toast |
| `removeToast(id)` | Remove a toast |
| `updateToast(id, data)` | Update toast data |

### Returned Actions

| Action | Description |
|--------|-------------|
| `portal` | Portal action for rendering |

### Toast Type

```typescript
type Toast<T> = {
  id: string;
  data: T;
  closeDelay?: number;
  type?: 'foreground' | 'background';
};
```

## Data Attributes

### Content
- `[data-melt-toast-content]` - Present on content

### Title
- `[data-melt-toast-title]` - Present on title

### Description
- `[data-melt-toast-description]` - Present on description

### Close
- `[data-melt-toast-close]` - Present on close button

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Tab` | Navigate focusable elements |
| `Escape` | Dismiss toast |

## Accessibility

Follows [WAI-ARIA Alert Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/alert/):
- `role="alert"` for important messages
- Focus management
- Screen reader announcements

## Setup

### 1. Create Toaster Component

```svelte
<!-- Toaster.svelte -->
<script lang="ts" module>
  import { createToaster, melt } from '@melt-ui/svelte';
  import { X, CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-svelte';

  export type ToastData = {
    title: string;
    description?: string;
    type: 'success' | 'error' | 'warning' | 'info';
  };

  const {
    elements: { content, title, description, close },
    helpers: { addToast },
    states: { toasts },
    actions: { portal }
  } = createToaster<ToastData>();

  export { addToast };
</script>

<script lang="ts">
  const icons = {
    success: CheckCircle,
    error: AlertCircle,
    warning: AlertTriangle,
    info: Info
  };

  const styles = {
    success: 'bg-green-50 border-green-500 text-green-800',
    error: 'bg-red-50 border-red-500 text-red-800',
    warning: 'bg-yellow-50 border-yellow-500 text-yellow-800',
    info: 'bg-blue-50 border-blue-500 text-blue-800'
  };
</script>

<div use:portal class="fixed top-4 right-4 z-50 flex flex-col gap-2 w-80">
  {#each $toasts as t (t.id)}
    {@const Icon = icons[t.data.type]}
    <div
      use:melt={$content(t.id)}
      class="flex items-start gap-3 p-4 rounded-lg border-l-4 shadow-lg {styles[t.data.type]}"
    >
      <Icon class="w-5 h-5 flex-shrink-0" />
      <div class="flex-1 min-w-0">
        <p use:melt={$title(t.id)} class="font-medium">{t.data.title}</p>
        {#if t.data.description}
          <p use:melt={$description(t.id)} class="text-sm mt-1 opacity-80">
            {t.data.description}
          </p>
        {/if}
      </div>
      <button
        use:melt={$close(t.id)}
        class="flex-shrink-0 p-1 rounded hover:bg-black/5 transition-colors"
      >
        <X class="w-4 h-4" />
      </button>
    </div>
  {/each}
</div>
```

### 2. Add to Layout

```svelte
<!-- +layout.svelte -->
<script lang="ts">
  import Toaster from '$lib/components/ui/Toaster.svelte';
</script>

<slot />
<Toaster />
```

### 3. Use the Toast

```svelte
<script lang="ts">
  import { addToast } from '$lib/components/ui/Toaster.svelte';

  function handleSuccess() {
    addToast({
      data: {
        title: 'Request submitted!',
        description: 'Your vacation request has been sent for approval.',
        type: 'success'
      }
    });
  }

  function handleError() {
    addToast({
      data: {
        title: 'Error',
        description: 'Something went wrong. Please try again.',
        type: 'error'
      },
      closeDelay: 10000 // Longer for errors
    });
  }
</script>
```

## Styling with Tailwind

```css
/* Toast animations */
[data-melt-toast-content] {
  animation: slideIn 200ms ease-out;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

/* Toast types */
[data-melt-toast-content].success {
  @apply bg-green-50 border-l-4 border-green-500;
}

[data-melt-toast-content].error {
  @apply bg-red-50 border-l-4 border-red-500;
}

[data-melt-toast-content].warning {
  @apply bg-yellow-50 border-l-4 border-yellow-500;
}

[data-melt-toast-content].info {
  @apply bg-blue-50 border-l-4 border-blue-500;
}
```

## VacayTracker Toast Store

```typescript
// toast.svelte.ts
import { createToaster, melt } from '@melt-ui/svelte';

export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface ToastData {
  title: string;
  description?: string;
  type: ToastType;
}

const toaster = createToaster<ToastData>({
  closeDelay: 5000,
  hover: 'pause'
});

export const {
  elements,
  helpers: { addToast, removeToast },
  states: { toasts },
  actions: { portal }
} = toaster;

// Convenience methods
export function toast(type: ToastType, title: string, description?: string) {
  addToast({
    data: { type, title, description }
  });
}

export const toastSuccess = (title: string, description?: string) =>
  toast('success', title, description);

export const toastError = (title: string, description?: string) =>
  toast('error', title, description);

export const toastWarning = (title: string, description?: string) =>
  toast('warning', title, description);

export const toastInfo = (title: string, description?: string) =>
  toast('info', title, description);
```

## With Progress Bar

```svelte
<script lang="ts">
  // Each toast provides getPercentage() helper
</script>

{#each $toasts as t (t.id)}
  {@const percentage = t.getPercentage()}
  <div use:melt={$content(t.id)} class="relative overflow-hidden ...">
    <!-- Toast content -->

    <!-- Progress bar -->
    <div
      class="absolute bottom-0 left-0 h-1 bg-current opacity-30 transition-all"
      style="width: {100 - percentage}%"
    ></div>
  </div>
{/each}
```

## With Actions

```svelte
<script lang="ts" module>
  export type ToastData = {
    title: string;
    description?: string;
    type: 'success' | 'error';
    action?: {
      label: string;
      onClick: () => void;
    };
  };
</script>

{#each $toasts as t (t.id)}
  <div use:melt={$content(t.id)} class="...">
    <div class="flex-1">
      <p use:melt={$title(t.id)}>{t.data.title}</p>
      {#if t.data.description}
        <p use:melt={$description(t.id)}>{t.data.description}</p>
      {/if}
    </div>

    {#if t.data.action}
      <button
        onclick={() => {
          t.data.action?.onClick();
          removeToast(t.id);
        }}
        class="px-3 py-1 text-sm font-medium bg-white/50 rounded hover:bg-white/80"
      >
        {t.data.action.label}
      </button>
    {/if}

    <button use:melt={$close(t.id)}>×</button>
  </div>
{/each}

<!-- Usage -->
<script>
  addToast({
    data: {
      title: 'File deleted',
      type: 'success',
      action: {
        label: 'Undo',
        onClick: () => restoreFile()
      }
    }
  });
</script>
```

## Positioning

```svelte
<!-- Top right (default) -->
<div use:portal class="fixed top-4 right-4 z-50">

<!-- Top left -->
<div use:portal class="fixed top-4 left-4 z-50">

<!-- Bottom right -->
<div use:portal class="fixed bottom-4 right-4 z-50">

<!-- Bottom center -->
<div use:portal class="fixed bottom-4 left-1/2 -translate-x-1/2 z-50">
```
