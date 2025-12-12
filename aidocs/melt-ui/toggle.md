# Toggle

A button that can be toggled between on and off states.

## Use Cases

- Text formatting buttons (bold, italic)
- View mode toggles
- Feature toggles
- Filter buttons

## Installation

```typescript
import { createToggle } from '@melt-ui/svelte';
```

## API Reference

### createToggle Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `disabled` | `boolean` | `false` | Disable the toggle |
| `defaultPressed` | `boolean` | `false` | Initial pressed state |
| `pressed` | `Writable<boolean>` | — | Controlled pressed state |
| `onPressedChange` | `ChangeFn` | — | Pressed change callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Toggle button element |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `pressed` | `Writable<boolean>` | Current pressed state |

## Data Attributes

### Root
- `[data-state]` - `'on'` or `'off'`
- `[data-disabled]` - Present when disabled
- `[data-melt-toggle]` - Present on root

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Space` | Toggle state |
| `Enter` | Toggle state |

## Accessibility

- Uses `aria-pressed` attribute
- Screen reader announces state
- Keyboard accessible

## Example

```svelte
<script lang="ts">
  import { createToggle, melt } from '@melt-ui/svelte';
  import { Bold } from 'lucide-svelte';

  const {
    elements: { root },
    states: { pressed }
  } = createToggle();
</script>

<button
  use:melt={$root}
  class="p-2 rounded-lg transition-colors
    data-[state=off]:bg-ocean-100 data-[state=off]:text-ocean-600 data-[state=off]:hover:bg-ocean-200
    data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
>
  <Bold class="h-5 w-5" />
</button>
```

## Styling with Tailwind

```css
/* Toggle off state */
[data-melt-toggle][data-state="off"] {
  @apply bg-ocean-100 text-ocean-600 hover:bg-ocean-200;
}

/* Toggle on state */
[data-melt-toggle][data-state="on"] {
  @apply bg-ocean-500 text-white;
}

/* Disabled state */
[data-melt-toggle][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}

/* Focus ring */
[data-melt-toggle]:focus-visible {
  @apply ring-2 ring-ocean-500 ring-offset-2;
}
```

## Text Formatting Toolbar

```svelte
<script lang="ts">
  import { createToggle, melt } from '@melt-ui/svelte';
  import { Bold, Italic, Underline } from 'lucide-svelte';

  const boldToggle = createToggle();
  const italicToggle = createToggle();
  const underlineToggle = createToggle();
</script>

<div class="flex gap-1 p-1 bg-ocean-100 rounded-lg">
  <button
    use:melt={boldToggle.elements.root}
    class="p-2 rounded transition-colors
      data-[state=off]:text-ocean-600 data-[state=off]:hover:bg-ocean-200
      data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
    title="Bold"
  >
    <Bold class="h-4 w-4" />
  </button>

  <button
    use:melt={italicToggle.elements.root}
    class="p-2 rounded transition-colors
      data-[state=off]:text-ocean-600 data-[state=off]:hover:bg-ocean-200
      data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
    title="Italic"
  >
    <Italic class="h-4 w-4" />
  </button>

  <button
    use:melt={underlineToggle.elements.root}
    class="p-2 rounded transition-colors
      data-[state=off]:text-ocean-600 data-[state=off]:hover:bg-ocean-200
      data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
    title="Underline"
  >
    <Underline class="h-4 w-4" />
  </button>
</div>
```

## Filter Toggle

```svelte
<script lang="ts">
  import { createToggle, melt } from '@melt-ui/svelte';
  import { Check } from 'lucide-svelte';

  const { elements: { root }, states: { pressed } } = createToggle();
</script>

<button
  use:melt={$root}
  class="inline-flex items-center gap-2 px-4 py-2 rounded-full border transition-colors
    data-[state=off]:border-ocean-300 data-[state=off]:text-ocean-600 data-[state=off]:hover:border-ocean-400
    data-[state=on]:border-ocean-500 data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
>
  {#if $pressed}
    <Check class="h-4 w-4" />
  {/if}
  Approved Only
</button>
```

## View Mode Toggle

```svelte
<script lang="ts">
  import { createToggle, melt } from '@melt-ui/svelte';
  import { Grid, List } from 'lucide-svelte';

  const { elements: { root }, states: { pressed } } = createToggle({
    defaultPressed: false // false = list, true = grid
  });

  const viewMode = $derived($pressed ? 'grid' : 'list');
</script>

<button
  use:melt={$root}
  class="p-2 rounded-lg bg-ocean-100 text-ocean-600 hover:bg-ocean-200 transition-colors"
  title="{$pressed ? 'Switch to list view' : 'Switch to grid view'}"
>
  {#if $pressed}
    <List class="h-5 w-5" />
  {:else}
    <Grid class="h-5 w-5" />
  {/if}
</button>
```

## Controlled State

```svelte
<script lang="ts">
  import { writable } from 'svelte/store';

  const isPressed = writable(false);

  const toggle = createToggle({
    pressed: isPressed,
    onPressedChange: ({ next }) => {
      console.log('Toggle changed:', next);
      return next;
    }
  });

  // Programmatic control
  function turnOn() {
    isPressed.set(true);
  }

  function turnOff() {
    isPressed.set(false);
  }
</script>
```

## With Label

```svelte
<label class="flex items-center gap-3 cursor-pointer">
  <button
    use:melt={$root}
    class="p-2 rounded transition-colors
      data-[state=off]:bg-ocean-100 data-[state=off]:text-ocean-600
      data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
  >
    <Star class="h-4 w-4" />
  </button>
  <span class="text-ocean-700">Mark as favorite</span>
</label>
```
