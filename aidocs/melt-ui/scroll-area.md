# Scroll Area

A custom scrollbar container with consistent cross-browser styling.

## Use Cases

- Custom styled scrollbars
- Chat message containers
- Code blocks
- Long lists
- Modal content areas

## Installation

```typescript
import { createScrollArea } from '@melt-ui/svelte';
```

## API Reference

### createScrollArea Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `type` | `'auto' \| 'always' \| 'scroll' \| 'hover'` | `'hover'` | Scrollbar visibility behavior |
| `dir` | `'ltr' \| 'rtl'` | `'ltr'` | Text direction |
| `hideDelay` | `number` | `600` | Delay before hiding scrollbar (ms) |
| `scrollHideDelay` | `number` | `600` | Delay after scrolling stops |

### Type Options

- `'auto'` - Shows scrollbar when content overflows
- `'always'` - Always shows scrollbar
- `'scroll'` - Shows scrollbar while scrolling
- `'hover'` - Shows scrollbar on hover

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Scroll area container |
| `viewport` | Scrollable viewport |
| `content` | Content container |
| `scrollbarY` | Vertical scrollbar |
| `scrollbarX` | Horizontal scrollbar |
| `thumbY` | Vertical thumb |
| `thumbX` | Horizontal thumb |
| `corner` | Corner element (when both scrollbars) |

## Data Attributes

### Root
- `[data-melt-scroll-area]` - Present on root

### Scrollbar
- `[data-state]` - `'visible'` or `'hidden'`
- `[data-orientation]` - `'vertical'` or `'horizontal'`
- `[data-melt-scroll-area-scrollbar]` - Present on scrollbar

### Thumb
- `[data-melt-scroll-area-thumb]` - Present on thumb

## Accessibility

- Maintains native scroll behavior
- Keyboard scrolling works normally
- Screen readers access content normally

## Example

```svelte
<script lang="ts">
  import { createScrollArea, melt } from '@melt-ui/svelte';

  const {
    elements: { root, viewport, content, scrollbarY, thumbY }
  } = createScrollArea({
    type: 'hover'
  });
</script>

<div use:melt={$root} class="h-64 w-full overflow-hidden rounded-lg border border-ocean-200">
  <div use:melt={$viewport} class="h-full w-full">
    <div use:melt={$content} class="p-4">
      <h3 class="font-semibold text-ocean-800 mb-2">Scrollable Content</h3>
      <p class="text-ocean-600">
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <!-- Long content here -->
      </p>
    </div>
  </div>

  <div
    use:melt={$scrollbarY}
    class="flex h-full w-2.5 touch-none select-none border-l border-transparent p-0.5 transition-colors
      data-[state=visible]:bg-ocean-100"
  >
    <div
      use:melt={$thumbY}
      class="relative flex-1 rounded-full bg-ocean-400 hover:bg-ocean-500 transition-colors"
    ></div>
  </div>
</div>
```

## Styling with Tailwind

```css
/* Scrollbar track */
[data-melt-scroll-area-scrollbar] {
  @apply flex touch-none select-none transition-colors;
}

[data-melt-scroll-area-scrollbar][data-orientation="vertical"] {
  @apply h-full w-2.5 border-l border-transparent p-0.5;
}

[data-melt-scroll-area-scrollbar][data-orientation="horizontal"] {
  @apply w-full h-2.5 border-t border-transparent p-0.5 flex-col;
}

/* Scrollbar visibility */
[data-melt-scroll-area-scrollbar][data-state="visible"] {
  @apply bg-ocean-100;
}

[data-melt-scroll-area-scrollbar][data-state="hidden"] {
  @apply opacity-0;
}

/* Scrollbar thumb */
[data-melt-scroll-area-thumb] {
  @apply relative flex-1 rounded-full bg-ocean-400;
}

[data-melt-scroll-area-thumb]:hover {
  @apply bg-ocean-500;
}
```

## Chat Container

```svelte
<script lang="ts">
  import { createScrollArea, melt } from '@melt-ui/svelte';

  interface Message {
    id: string;
    content: string;
    sender: string;
    timestamp: Date;
  }

  let { messages }: { messages: Message[] } = $props();

  const {
    elements: { root, viewport, content, scrollbarY, thumbY }
  } = createScrollArea({
    type: 'auto'
  });

  // Auto-scroll to bottom on new messages
  let viewportEl: HTMLElement;

  $effect(() => {
    if (viewportEl && messages.length > 0) {
      viewportEl.scrollTop = viewportEl.scrollHeight;
    }
  });
</script>

<div use:melt={$root} class="h-96 rounded-xl bg-ocean-50 relative">
  <div
    use:melt={$viewport}
    bind:this={viewportEl}
    class="h-full w-full p-4"
  >
    <div use:melt={$content} class="space-y-4">
      {#each messages as message}
        <div class="flex gap-3">
          <div class="h-8 w-8 rounded-full bg-ocean-200 flex items-center justify-center">
            {message.sender[0]}
          </div>
          <div class="flex-1">
            <div class="flex items-baseline gap-2">
              <span class="font-medium text-ocean-800">{message.sender}</span>
              <span class="text-xs text-ocean-400">
                {message.timestamp.toLocaleTimeString()}
              </span>
            </div>
            <p class="text-ocean-700">{message.content}</p>
          </div>
        </div>
      {/each}
    </div>
  </div>

  <div
    use:melt={$scrollbarY}
    class="absolute right-0 top-0 flex h-full w-2.5 p-0.5"
  >
    <div use:melt={$thumbY} class="flex-1 rounded-full bg-ocean-300"></div>
  </div>
</div>
```

## Code Block

```svelte
<script lang="ts">
  const {
    elements: { root, viewport, content, scrollbarY, scrollbarX, thumbY, thumbX, corner }
  } = createScrollArea({
    type: 'always'
  });
</script>

<div
  use:melt={$root}
  class="max-h-64 w-full rounded-lg bg-gray-900 relative overflow-hidden"
>
  <div use:melt={$viewport} class="h-full w-full">
    <div use:melt={$content} class="p-4">
      <pre class="text-sm text-gray-100 whitespace-pre"><code>{code}</code></pre>
    </div>
  </div>

  <!-- Vertical scrollbar -->
  <div
    use:melt={$scrollbarY}
    class="absolute right-0 top-0 flex h-full w-2 p-0.5"
  >
    <div use:melt={$thumbY} class="flex-1 rounded-full bg-gray-600"></div>
  </div>

  <!-- Horizontal scrollbar -->
  <div
    use:melt={$scrollbarX}
    class="absolute bottom-0 left-0 flex h-2 w-full p-0.5 flex-col"
  >
    <div use:melt={$thumbX} class="flex-1 rounded-full bg-gray-600"></div>
  </div>

  <!-- Corner -->
  <div use:melt={$corner} class="absolute bottom-0 right-0 h-2 w-2 bg-gray-800"></div>
</div>
```

## Always Visible Scrollbar

```typescript
const scrollArea = createScrollArea({
  type: 'always' // Always show scrollbar
});
```

## With Custom Hide Delay

```typescript
const scrollArea = createScrollArea({
  type: 'hover',
  hideDelay: 1000, // Wait 1 second before hiding
  scrollHideDelay: 800 // Wait 0.8 seconds after scroll stops
});
```

## Modal Content

```svelte
<script lang="ts">
  import { createScrollArea, melt } from '@melt-ui/svelte';

  const {
    elements: { root, viewport, content, scrollbarY, thumbY }
  } = createScrollArea();
</script>

<!-- Inside a modal -->
<div class="modal-content">
  <div use:melt={$root} class="max-h-[60vh] relative">
    <div use:melt={$viewport} class="h-full w-full pr-4">
      <div use:melt={$content}>
        <!-- Long form content -->
      </div>
    </div>
    <div use:melt={$scrollbarY} class="absolute right-0 top-0 h-full w-2">
      <div use:melt={$thumbY} class="rounded-full bg-ocean-400"></div>
    </div>
  </div>
</div>
```
