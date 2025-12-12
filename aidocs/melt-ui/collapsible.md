# Collapsible

A simple interactive component for expanding and collapsing content panels.

## Use Cases

- FAQ questions
- Settings sections
- Mobile navigation menus
- Additional information toggles

## Installation

```typescript
import { createCollapsible } from '@melt-ui/svelte';
```

## API Reference

### createCollapsible Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `disabled` | `boolean` | `false` | Disable toggle interaction |
| `forceVisible` | `boolean` | `false` | Force visibility for animations |
| `defaultOpen` | `boolean` | `false` | Initial open state |
| `open` | `Writable<boolean>` | — | Controlled open state |
| `onOpenChange` | `ChangeFn` | — | Open state callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Container element |
| `trigger` | Button to toggle open state |
| `content` | Collapsible content panel |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `open` | `Writable<boolean>` | Current open state |

## Data Attributes

### Root
- `[data-state]` - `'open'` or `'closed'`
- `[data-disabled]` - Present when disabled
- `[data-melt-collapsible]` - Present on root element

### Trigger
- `[data-state]` - `'open'` or `'closed'`
- `[data-disabled]` - Present when disabled
- `[data-melt-collapsible-trigger]` - Present on trigger

### Content
- `[data-state]` - `'open'` or `'closed'`
- `[data-disabled]` - Present when disabled
- `[data-melt-collapsible-content]` - Present on content

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Space` | Toggle open state |
| `Enter` | Toggle open state |

## Accessibility

Follows the [WAI-ARIA Disclosure Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/disclosure/):
- `aria-expanded` on trigger
- `aria-controls` linking trigger to content
- Proper focus management

## Example

```svelte
<script lang="ts">
  import { createCollapsible, melt } from '@melt-ui/svelte';
  import { ChevronDown } from 'lucide-svelte';

  const {
    elements: { root, trigger, content },
    states: { open }
  } = createCollapsible();
</script>

<div use:melt={$root} class="w-full max-w-md">
  <div class="flex items-center justify-between">
    <span class="text-sm font-semibold text-ocean-800">Advanced Options</span>
    <button
      use:melt={$trigger}
      class="p-2 rounded-lg hover:bg-ocean-100 text-ocean-600 transition-colors"
    >
      <ChevronDown
        class="h-4 w-4 transition-transform duration-200 {$open ? 'rotate-180' : ''}"
      />
    </button>
  </div>

  {#if $open}
    <div
      use:melt={$content}
      class="mt-2 rounded-lg bg-ocean-50 p-4 space-y-2"
    >
      <p class="text-sm text-ocean-700">
        Additional configuration options are available here.
      </p>
      <!-- More content -->
    </div>
  {/if}
</div>
```

## Styling with Tailwind

```css
/* Trigger rotation */
[data-melt-collapsible-trigger][data-state="open"] svg {
  transform: rotate(180deg);
}

/* Content animation */
[data-melt-collapsible-content] {
  overflow: hidden;
}

[data-melt-collapsible-content][data-state="open"] {
  animation: slideDown 200ms ease-out;
}

[data-melt-collapsible-content][data-state="closed"] {
  animation: slideUp 200ms ease-out;
}

@keyframes slideDown {
  from {
    height: 0;
    opacity: 0;
  }
  to {
    height: var(--melt-content-height);
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    height: var(--melt-content-height);
    opacity: 1;
  }
  to {
    height: 0;
    opacity: 0;
  }
}
```

## With Animation

```svelte
<script lang="ts">
  import { createCollapsible, melt } from '@melt-ui/svelte';
  import { slide } from 'svelte/transition';

  const {
    elements: { root, trigger, content },
    states: { open }
  } = createCollapsible({ forceVisible: true });
</script>

<div use:melt={$root}>
  <button use:melt={$trigger}>Toggle</button>

  {#if $open}
    <div use:melt={$content} transition:slide={{ duration: 200 }}>
      Content here
    </div>
  {/if}
</div>
```

## Card Collapsible

```svelte
<script lang="ts">
  import { createCollapsible, melt } from '@melt-ui/svelte';
  import { ChevronDown, Settings } from 'lucide-svelte';

  const {
    elements: { root, trigger, content },
    states: { open }
  } = createCollapsible();
</script>

<div use:melt={$root} class="bg-white rounded-xl shadow-lg overflow-hidden">
  <button
    use:melt={$trigger}
    class="w-full flex items-center justify-between p-4 hover:bg-ocean-50 transition-colors"
  >
    <div class="flex items-center gap-3">
      <Settings class="h-5 w-5 text-ocean-500" />
      <span class="font-medium text-ocean-800">Settings</span>
    </div>
    <ChevronDown
      class="h-5 w-5 text-ocean-400 transition-transform {$open ? 'rotate-180' : ''}"
    />
  </button>

  {#if $open}
    <div use:melt={$content} class="border-t border-ocean-100 p-4 space-y-4">
      <!-- Settings content -->
      <div class="flex items-center justify-between">
        <span class="text-sm text-ocean-700">Email notifications</span>
        <input type="checkbox" class="toggle" />
      </div>
      <div class="flex items-center justify-between">
        <span class="text-sm text-ocean-700">Dark mode</span>
        <input type="checkbox" class="toggle" />
      </div>
    </div>
  {/if}
</div>
```

## Controlled State

```svelte
<script lang="ts">
  import { createCollapsible, melt } from '@melt-ui/svelte';
  import { writable } from 'svelte/store';

  const isOpen = writable(false);

  const { elements: { root, trigger, content } } = createCollapsible({
    open: isOpen,
    onOpenChange: ({ curr, next }) => {
      console.log('Changing from', curr, 'to', next);
      return next;
    }
  });

  // Programmatic control
  function openPanel() {
    isOpen.set(true);
  }
</script>
```
