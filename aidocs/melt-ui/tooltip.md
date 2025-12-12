# Tooltip

A popup displaying information when hovering or focusing on an element.

## Use Cases

- Button/icon explanations
- Shortened text full display
- Help hints
- Feature descriptions

## Installation

```typescript
import { createTooltip } from '@melt-ui/svelte';
```

## API Reference

### createTooltip Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `positioning` | `FloatingConfig` | `{ placement: 'top' }` | Popup positioning |
| `arrowSize` | `number` | `8` | Arrow size in pixels |
| `escapeBehavior` | `EscapeBehavior` | `'close'` | Escape key behavior |
| `forceVisible` | `boolean` | `false` | Force visibility |
| `portal` | `string \| HTMLElement \| null` | `'body'` | Portal target |
| `closeOnPointerDown` | `boolean` | `true` | Close on click |
| `openDelay` | `number` | `1000` | Delay before opening (ms) |
| `closeDelay` | `number` | `500` | Delay before closing (ms) |
| `disableHoverableContent` | `boolean` | `false` | Prevent hover on content |
| `group` | `boolean \| string` | — | Group tooltips |
| `defaultOpen` | `boolean` | `false` | Initial open state |
| `open` | `Writable<boolean>` | — | Controlled open state |
| `onOpenChange` | `ChangeFn` | — | Open state callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `trigger` | Element that triggers tooltip |
| `content` | Tooltip content container |
| `arrow` | Optional arrow element |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `open` | `Writable<boolean>` | Current open state |

## Data Attributes

### Trigger
- `[data-state]` - `'open'` or `'closed'`
- `[data-melt-tooltip-trigger]` - Present on trigger

### Content
- `[data-state]` - `'open'` or `'closed'`
- `[data-side]` - Placement side
- `[data-align]` - Placement alignment
- `[data-melt-tooltip-content]` - Present on content

### Arrow
- `[data-arrow]` - Present on arrow
- `[data-side]` - Arrow side
- `[data-melt-tooltip-arrow]` - Present on arrow

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Tab` | Open/close without delay |
| `Space` | Close if open |
| `Enter` | Close if open |
| `Escape` | Close if open |

## Accessibility

Follows [WAI-ARIA Tooltip Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/tooltip/):
- `role="tooltip"` on content
- Linked via `aria-describedby`
- Opens on hover and focus

## Example

```svelte
<script lang="ts">
  import { createTooltip, melt } from '@melt-ui/svelte';
  import { HelpCircle } from 'lucide-svelte';

  const {
    elements: { trigger, content, arrow },
    states: { open }
  } = createTooltip({
    forceVisible: true,
    positioning: { placement: 'top' }
  });
</script>

<button
  use:melt={$trigger}
  class="p-2 rounded-lg text-ocean-500 hover:bg-ocean-100 transition-colors"
>
  <HelpCircle class="h-5 w-5" />
</button>

{#if $open}
  <div
    use:melt={$content}
    class="z-50 px-3 py-2 rounded-lg bg-ocean-800 text-white text-sm shadow-lg"
  >
    <div
      use:melt={$arrow}
      class="absolute h-2 w-2 bg-ocean-800 rotate-45"
    ></div>
    Click here for help
  </div>
{/if}
```

## Styling with Tailwind

```css
/* Tooltip content */
[data-melt-tooltip-content] {
  @apply px-3 py-2 rounded-lg bg-ocean-800 text-white text-sm shadow-lg;
  @apply transition-opacity duration-200;
}

[data-melt-tooltip-content][data-state="open"] {
  @apply opacity-100;
}

[data-melt-tooltip-content][data-state="closed"] {
  @apply opacity-0;
}

/* Arrow positioning based on side */
[data-melt-tooltip-content][data-side="top"] [data-melt-tooltip-arrow] {
  @apply bottom-[-4px];
}

[data-melt-tooltip-content][data-side="bottom"] [data-melt-tooltip-arrow] {
  @apply top-[-4px];
}

[data-melt-tooltip-content][data-side="left"] [data-melt-tooltip-arrow] {
  @apply right-[-4px];
}

[data-melt-tooltip-content][data-side="right"] [data-melt-tooltip-arrow] {
  @apply left-[-4px];
}
```

## Button Tooltip

```svelte
<script lang="ts">
  import { createTooltip, melt } from '@melt-ui/svelte';
  import { Save } from 'lucide-svelte';

  const { elements: { trigger, content }, states: { open } } = createTooltip({
    openDelay: 500,
    closeDelay: 200
  });
</script>

<button
  use:melt={$trigger}
  class="p-2 rounded-lg bg-ocean-500 text-white hover:bg-ocean-600"
>
  <Save class="h-5 w-5" />
</button>

{#if $open}
  <div
    use:melt={$content}
    class="z-50 px-2 py-1 rounded bg-gray-900 text-white text-xs"
  >
    Save changes (Ctrl+S)
  </div>
{/if}
```

## Truncated Text Tooltip

```svelte
<script lang="ts">
  import { createTooltip, melt } from '@melt-ui/svelte';

  let { text }: { text: string } = $props();

  const { elements: { trigger, content }, states: { open } } = createTooltip({
    openDelay: 300
  });
</script>

<span
  use:melt={$trigger}
  class="block truncate cursor-default"
>
  {text}
</span>

{#if $open}
  <div
    use:melt={$content}
    class="z-50 max-w-xs px-3 py-2 rounded-lg bg-ocean-800 text-white text-sm"
  >
    {text}
  </div>
{/if}
```

## VacayTracker Info Tooltip

```svelte
<script lang="ts">
  import { createTooltip, melt } from '@melt-ui/svelte';
  import { Info } from 'lucide-svelte';

  interface Props {
    content: string;
  }

  let { content: tooltipContent }: Props = $props();

  const {
    elements: { trigger, content, arrow },
    states: { open }
  } = createTooltip({
    forceVisible: true,
    openDelay: 300,
    positioning: { placement: 'top' }
  });
</script>

<button
  use:melt={$trigger}
  class="inline-flex items-center justify-center w-5 h-5 rounded-full text-ocean-400 hover:text-ocean-600 hover:bg-ocean-100 transition-colors"
>
  <Info class="w-4 h-4" />
</button>

{#if $open}
  <div
    use:melt={$content}
    class="z-50 max-w-xs px-3 py-2 rounded-lg bg-ocean-800 text-white text-sm shadow-lg"
  >
    <div use:melt={$arrow} class="absolute h-2 w-2 bg-ocean-800 rotate-45"></div>
    {tooltipContent}
  </div>
{/if}
```

## Grouped Tooltips

```svelte
<script lang="ts">
  // Tooltips in the same group close when another opens
  const tooltip1 = createTooltip({ group: 'toolbar' });
  const tooltip2 = createTooltip({ group: 'toolbar' });
  const tooltip3 = createTooltip({ group: 'toolbar' });
</script>

<div class="flex gap-2">
  <button use:melt={tooltip1.elements.trigger}>
    <Bold class="h-4 w-4" />
  </button>
  <button use:melt={tooltip2.elements.trigger}>
    <Italic class="h-4 w-4" />
  </button>
  <button use:melt={tooltip3.elements.trigger}>
    <Underline class="h-4 w-4" />
  </button>
</div>
```

## Custom Delays

```typescript
// Quick tooltip for obvious actions
const quickTooltip = createTooltip({
  openDelay: 300,
  closeDelay: 100
});

// Slow tooltip for detailed explanations
const slowTooltip = createTooltip({
  openDelay: 1500,
  closeDelay: 300
});
```

## Positioning

```typescript
// Above element
createTooltip({ positioning: { placement: 'top' } });

// Below element
createTooltip({ positioning: { placement: 'bottom' } });

// To the right
createTooltip({ positioning: { placement: 'right' } });

// With offset
createTooltip({
  positioning: {
    placement: 'top',
    offset: { mainAxis: 8 }
  }
});
```

## Rich Content

```svelte
{#if $open}
  <div use:melt={$content} class="z-50 p-4 rounded-xl bg-white shadow-xl border border-ocean-200 max-w-xs">
    <h4 class="font-semibold text-ocean-800 mb-1">Vacation Balance</h4>
    <p class="text-sm text-ocean-600">
      Your remaining vacation days for this year. Unused days expire on December 31st.
    </p>
  </div>
{/if}
```
