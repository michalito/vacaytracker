# Dialog

A modal window that overlays the primary content, rendering underlying content inert.

## Use Cases

- Confirmation dialogs
- Form modals (create/edit)
- Alert dialogs
- Image lightboxes
- Settings panels

## Installation

```typescript
import { createDialog } from '@melt-ui/svelte';
```

## API Reference

### createDialog Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `role` | `'dialog' \| 'alertdialog'` | `'dialog'` | Semantic role |
| `preventScroll` | `boolean` | `true` | Block body scrolling |
| `escapeBehavior` | `EscapeBehavior` | `'close'` | Escape key behavior |
| `closeOnOutsideClick` | `boolean` | `true` | Close on overlay click |
| `portal` | `string \| HTMLElement \| null` | `'body'` | Portal target |
| `forceVisible` | `boolean` | `false` | Force visibility for transitions |
| `defaultOpen` | `boolean` | `false` | Initial open state |
| `open` | `Writable<boolean>` | — | Controlled open state |
| `onOpenChange` | `ChangeFn` | — | Open state callback |

### EscapeBehavior Options

- `'close'` - Close immediately on Escape
- `'ignore'` - Ignore Escape key
- `'defer-otherwise-close'` - Delegate to parent, then close
- `'defer-otherwise-ignore'` - Delegate to parent, then ignore

### Returned Elements

| Element | Description |
|---------|-------------|
| `trigger` | Button that opens the dialog |
| `portalled` | Container for portal rendering |
| `overlay` | Background overlay element |
| `content` | Main dialog container |
| `title` | Dialog heading (linked via aria-labelledby) |
| `description` | Dialog description (linked via aria-describedby) |
| `close` | Close button element |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `open` | `Writable<boolean>` | Current open state |

## Data Attributes

### Overlay & Content
- `[data-state]` - `'open'` or `'closed'`
- `[data-melt-dialog-overlay]` - Present on overlay
- `[data-melt-dialog-content]` - Present on content

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Escape` | Close dialog (configurable) |
| `Tab` | Navigate focusable elements (trapped) |
| `Shift + Tab` | Navigate backwards |

## Accessibility

Follows [WAI-ARIA Dialog Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/dialog-modal/):
- Focus is trapped within dialog
- Focus moves to first focusable element on open
- Focus returns to trigger on close
- `aria-labelledby` links to title
- `aria-describedby` links to description

## Example

```svelte
<script lang="ts">
  import { createDialog, melt } from '@melt-ui/svelte';
  import { X } from 'lucide-svelte';

  const {
    elements: { trigger, overlay, content, title, description, close, portalled },
    states: { open }
  } = createDialog({
    forceVisible: true
  });
</script>

<button
  use:melt={$trigger}
  class="px-4 py-2 rounded-lg bg-ocean-500 text-white font-medium hover:bg-ocean-600 transition-colors"
>
  Open Dialog
</button>

{#if $open}
  <div use:melt={$portalled}>
    <div
      use:melt={$overlay}
      class="fixed inset-0 z-50 bg-ocean-900/50 backdrop-blur-sm"
    ></div>

    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        use:melt={$content}
        class="bg-white/95 backdrop-blur-md rounded-2xl shadow-xl border border-white/30 w-full max-w-md"
      >
        <!-- Header -->
        <div class="flex items-center justify-between p-4 border-b border-ocean-100/50">
          <h2 use:melt={$title} class="text-lg font-semibold text-ocean-800">
            Dialog Title
          </h2>
          <button
            use:melt={$close}
            class="p-1 rounded-lg text-ocean-400 hover:text-ocean-600 hover:bg-ocean-500/10 transition-colors"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Body -->
        <div class="p-4">
          <p use:melt={$description} class="text-ocean-600">
            This is the dialog description.
          </p>
        </div>

        <!-- Footer -->
        <div class="flex gap-3 p-4 border-t border-ocean-100/50">
          <button
            use:melt={$close}
            class="flex-1 px-4 py-2 rounded-lg border-2 border-ocean-500/40 text-ocean-700 hover:bg-ocean-50"
          >
            Cancel
          </button>
          <button
            class="flex-1 px-4 py-2 rounded-lg bg-ocean-500 text-white hover:bg-ocean-600"
          >
            Confirm
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
```

## Styling with Tailwind

```css
/* Overlay fade-in */
[data-melt-dialog-overlay] {
  @apply transition-opacity duration-200;
}

[data-melt-dialog-overlay][data-state="open"] {
  @apply opacity-100;
}

[data-melt-dialog-overlay][data-state="closed"] {
  @apply opacity-0;
}

/* Content animation */
[data-melt-dialog-content] {
  @apply transition-all duration-200;
}

[data-melt-dialog-content][data-state="open"] {
  @apply opacity-100 scale-100;
}

[data-melt-dialog-content][data-state="closed"] {
  @apply opacity-0 scale-95;
}
```

## Alert Dialog

For actions requiring confirmation:

```svelte
<script lang="ts">
  const {
    elements: { trigger, overlay, content, title, description, close }
  } = createDialog({
    role: 'alertdialog', // Important: changes semantic role
    closeOnOutsideClick: false, // Don't dismiss on outside click
    escapeBehavior: 'ignore' // Require explicit action
  });
</script>

<div use:melt={$content}>
  <h2 use:melt={$title}>Delete User?</h2>
  <p use:melt={$description}>
    This action cannot be undone. This will permanently delete the user account.
  </p>
  <div class="flex gap-3">
    <button use:melt={$close}>Cancel</button>
    <button onclick={handleDelete} class="bg-red-500 text-white">
      Delete
    </button>
  </div>
</div>
```

## VacayTracker Request Modal

```svelte
<script lang="ts">
  import { createDialog, melt } from '@melt-ui/svelte';
  import { vacation } from '$lib/stores/vacation.svelte';
  import { toast } from '$lib/stores/toast.svelte';
  import { Palmtree, X } from 'lucide-svelte';

  interface Props {
    open: boolean;
  }

  let { open = $bindable(false) }: Props = $props();

  const dialog = createDialog({
    forceVisible: true,
    onOpenChange: ({ next }) => {
      open = next;
      return next;
    }
  });

  const { elements: { overlay, content, title, close, portalled } } = dialog;

  let startDate = $state('');
  let endDate = $state('');
  let isSubmitting = $state(false);

  // Sync external open state with dialog
  $effect(() => {
    dialog.states.open.set(open);
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    isSubmitting = true;
    try {
      await vacation.createRequest({ startDate, endDate });
      toast.success('Request submitted!');
      dialog.states.open.set(false);
    } catch (error) {
      toast.error('Failed to submit request');
    } finally {
      isSubmitting = false;
    }
  }
</script>

{#if open}
  <div use:melt={$portalled}>
    <div
      use:melt={$overlay}
      class="fixed inset-0 z-50 bg-ocean-900/50 backdrop-blur-sm"
    ></div>

    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        use:melt={$content}
        class="bg-white/95 backdrop-blur-md rounded-2xl shadow-xl w-full max-w-md"
      >
        <div class="flex items-center justify-between p-4 border-b border-ocean-100/50">
          <div class="flex items-center gap-2">
            <Palmtree class="w-5 h-5 text-ocean-500" />
            <h2 use:melt={$title} class="text-lg font-semibold text-ocean-800">
              Request Vacation
            </h2>
          </div>
          <button use:melt={$close} class="p-1 rounded-lg text-ocean-400 hover:text-ocean-600">
            <X class="w-5 h-5" />
          </button>
        </div>

        <form onsubmit={handleSubmit} class="p-4 space-y-4">
          <!-- Form fields -->
          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Submitting...' : 'Submit Request'}
          </button>
        </form>
      </div>
    </div>
  </div>
{/if}
```

## Nested Dialogs

```svelte
<script lang="ts">
  const parentDialog = createDialog();
  const childDialog = createDialog();
</script>

<!-- Parent Dialog -->
<div use:melt={parentDialog.elements.content}>
  <button use:melt={childDialog.elements.trigger}>
    Open Nested Dialog
  </button>
</div>

<!-- Child Dialog -->
<div use:melt={childDialog.elements.content}>
  <!-- Child content -->
</div>
```

## Controlled State

```svelte
<script lang="ts">
  import { writable } from 'svelte/store';

  const isOpen = writable(false);

  const dialog = createDialog({
    open: isOpen,
    onOpenChange: ({ next }) => {
      console.log('Dialog state:', next);
      return next;
    }
  });

  // Programmatic control
  function openDialog() {
    isOpen.set(true);
  }

  function closeDialog() {
    isOpen.set(false);
  }
</script>
```
