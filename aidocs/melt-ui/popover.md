# Popover

A floating panel triggered by a button for displaying rich content.

## Use Cases

- Filter panels
- Quick actions
- Form inputs
- Info panels
- User menus

## Installation

```typescript
import { createPopover } from '@melt-ui/svelte';
```

## API Reference

### createPopover Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `positioning` | `FloatingConfig` | `{ placement: 'bottom' }` | Floating positioning |
| `defaultOpen` | `boolean` | `false` | Initial open state |
| `open` | `Writable<boolean>` | — | Controlled open state |
| `preventScroll` | `boolean` | `true` | Block page scroll |
| `closeOnOutsideClick` | `boolean` | `true` | Close on click outside |
| `closeOnEscape` | `boolean` | `true` | Close on Escape key |
| `disableFocusTrap` | `boolean` | `false` | Disable focus trap |
| `arrowSize` | `number` | `8` | Arrow size in pixels |
| `escapeBehavior` | `EscapeBehavior` | `'close'` | Escape key behavior |
| `forceVisible` | `boolean` | `false` | Force visibility |
| `portal` | `string \| HTMLElement \| null` | `'body'` | Portal target |
| `onOpenChange` | `ChangeFn` | — | Open state callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `trigger` | Button that opens popover |
| `content` | Popover container |
| `arrow` | Optional arrow element |
| `close` | Close button |
| `overlay` | Optional backdrop overlay |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `open` | `Writable<boolean>` | Current open state |

## Data Attributes

### Trigger
- `[data-state]` - `'open'` or `'closed'`
- `[data-melt-popover-trigger]` - Present on trigger

### Content
- `[data-state]` - `'open'` or `'closed'`
- `[data-side]` - Placement side
- `[data-align]` - Placement alignment
- `[data-melt-popover-content]` - Present on content

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Space` / `Enter` | Toggle popover |
| `Tab` | Navigate focusable elements |
| `Shift + Tab` | Navigate backwards |
| `Escape` | Close popover |

## Accessibility

Follows [WAI-ARIA Dialog Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/dialog-modal/):
- Focus trapped within popover
- Focus returns to trigger on close
- Proper ARIA attributes

## Example

```svelte
<script lang="ts">
  import { createPopover, melt } from '@melt-ui/svelte';
  import { Settings, X } from 'lucide-svelte';

  const {
    elements: { trigger, content, arrow, close },
    states: { open }
  } = createPopover({
    forceVisible: true,
    positioning: { placement: 'bottom-end' }
  });
</script>

<button
  use:melt={$trigger}
  class="p-2 rounded-lg text-ocean-600 hover:bg-ocean-100 transition-colors"
>
  <Settings class="h-5 w-5" />
</button>

{#if $open}
  <div
    use:melt={$content}
    class="z-50 w-72 rounded-xl bg-white shadow-xl border border-ocean-200 p-4"
  >
    <div
      use:melt={$arrow}
      class="absolute h-3 w-3 bg-white border-l border-t border-ocean-200 rotate-45"
    ></div>

    <div class="flex items-center justify-between mb-4">
      <h3 class="font-semibold text-ocean-800">Settings</h3>
      <button
        use:melt={$close}
        class="p-1 rounded text-ocean-400 hover:text-ocean-600 hover:bg-ocean-100"
      >
        <X class="h-4 w-4" />
      </button>
    </div>

    <div class="space-y-4">
      <div class="flex items-center justify-between">
        <label class="text-sm text-ocean-700">Email notifications</label>
        <input type="checkbox" class="toggle" />
      </div>
      <div class="flex items-center justify-between">
        <label class="text-sm text-ocean-700">Dark mode</label>
        <input type="checkbox" class="toggle" />
      </div>
    </div>
  </div>
{/if}
```

## Styling with Tailwind

```css
/* Content animation */
[data-melt-popover-content] {
  @apply transition-all duration-200;
}

[data-melt-popover-content][data-state="open"] {
  @apply opacity-100 scale-100;
}

[data-melt-popover-content][data-state="closed"] {
  @apply opacity-0 scale-95;
}

/* Arrow based on side */
[data-melt-popover-content][data-side="top"] [data-arrow] {
  @apply bottom-[-6px] rotate-[225deg];
}

[data-melt-popover-content][data-side="bottom"] [data-arrow] {
  @apply top-[-6px] rotate-45;
}

[data-melt-popover-content][data-side="left"] [data-arrow] {
  @apply right-[-6px] rotate-[135deg];
}

[data-melt-popover-content][data-side="right"] [data-arrow] {
  @apply left-[-6px] rotate-[-45deg];
}
```

## VacayTracker Filter Panel

```svelte
<script lang="ts">
  import { createPopover, melt } from '@melt-ui/svelte';
  import { Filter, X } from 'lucide-svelte';

  const {
    elements: { trigger, content, close },
    states: { open }
  } = createPopover({
    forceVisible: true,
    positioning: { placement: 'bottom-start' }
  });

  let showApproved = $state(true);
  let showPending = $state(true);
  let showRejected = $state(false);
</script>

<button
  use:melt={$trigger}
  class="flex items-center gap-2 px-4 py-2 rounded-lg border-2 border-ocean-500/40 text-ocean-700 hover:bg-ocean-50"
>
  <Filter class="h-4 w-4" />
  <span>Filters</span>
</button>

{#if $open}
  <div
    use:melt={$content}
    class="z-50 w-64 rounded-xl bg-white/95 backdrop-blur-md shadow-xl border border-white/30 p-4"
  >
    <div class="flex items-center justify-between mb-4">
      <h3 class="font-semibold text-ocean-800">Filter Requests</h3>
      <button use:melt={$close} class="p-1 rounded text-ocean-400 hover:bg-ocean-100">
        <X class="h-4 w-4" />
      </button>
    </div>

    <div class="space-y-3">
      <label class="flex items-center gap-3 cursor-pointer">
        <input type="checkbox" bind:checked={showApproved} class="rounded border-ocean-400" />
        <span class="text-sm text-ocean-700">Approved</span>
      </label>
      <label class="flex items-center gap-3 cursor-pointer">
        <input type="checkbox" bind:checked={showPending} class="rounded border-ocean-400" />
        <span class="text-sm text-ocean-700">Pending</span>
      </label>
      <label class="flex items-center gap-3 cursor-pointer">
        <input type="checkbox" bind:checked={showRejected} class="rounded border-ocean-400" />
        <span class="text-sm text-ocean-700">Rejected</span>
      </label>
    </div>

    <div class="flex gap-2 mt-4">
      <button
        onclick={() => { showApproved = true; showPending = true; showRejected = false; }}
        class="flex-1 px-3 py-1.5 text-sm border border-ocean-200 rounded-lg hover:bg-ocean-50"
      >
        Reset
      </button>
      <button
        use:melt={$close}
        class="flex-1 px-3 py-1.5 text-sm bg-ocean-500 text-white rounded-lg hover:bg-ocean-600"
      >
        Apply
      </button>
    </div>
  </div>
{/if}
```

## With Overlay

```svelte
<script lang="ts">
  const {
    elements: { trigger, content, overlay }
  } = createPopover();
</script>

{#if $open}
  <div
    use:melt={$overlay}
    class="fixed inset-0 bg-black/20 z-40"
  ></div>

  <div use:melt={$content} class="z-50 ...">
    <!-- Content -->
  </div>
{/if}
```

## Positioning Options

```typescript
// Bottom left
const popover = createPopover({
  positioning: { placement: 'bottom-start' }
});

// Right side
const popover = createPopover({
  positioning: { placement: 'right' }
});

// With offset
const popover = createPopover({
  positioning: {
    placement: 'bottom',
    offset: { mainAxis: 8, crossAxis: 0 }
  }
});
```

## Controlled State

```svelte
<script lang="ts">
  import { writable } from 'svelte/store';

  const isOpen = writable(false);

  const popover = createPopover({
    open: isOpen,
    onOpenChange: ({ next }) => {
      console.log('Popover:', next ? 'opening' : 'closing');
      return next;
    }
  });

  // Programmatic control
  function openPopover() {
    isOpen.set(true);
  }
</script>
```

## Without Focus Trap

```typescript
// For popovers with simple content
const popover = createPopover({
  disableFocusTrap: true
});
```
