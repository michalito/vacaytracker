# Slider

An input control for selecting values from a range.

## Use Cases

- Volume controls
- Price range filters
- Rating inputs
- Zoom controls
- Progress adjustments

## Installation

```typescript
import { createSlider } from '@melt-ui/svelte';
```

## API Reference

### createSlider Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `min` | `number` | `0` | Minimum value |
| `max` | `number` | `100` | Maximum value |
| `step` | `number` | `1` | Step increment |
| `orientation` | `'horizontal' \| 'vertical'` | `'horizontal'` | Slider direction |
| `dir` | `'ltr' \| 'rtl'` | `'ltr'` | Text direction |
| `autoSort` | `boolean` | `true` | Auto-sort values |
| `disabled` | `boolean` | `false` | Disable slider |
| `defaultValue` | `number[]` | `[]` | Initial value(s) |
| `value` | `Writable<number[]>` | — | Controlled value(s) |
| `onValueChange` | `ChangeFn` | — | Value change callback |
| `onValueCommitted` | `(value: number[]) => void` | — | Called when interaction ends |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Slider container |
| `range` | Filled track portion |
| `thumb` | Draggable handle(s) |
| `tick` | Optional tick marks |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `value` | `Writable<number[]>` | Current value(s) |
| `active` | `Readable<boolean>` | Currently being dragged |

## Data Attributes

### Root
- `[data-orientation]` - `'horizontal'` or `'vertical'`
- `[data-disabled]` - Present when disabled
- `[data-melt-slider]` - Present on root

### Thumb
- `[data-value]` - Current thumb value
- `[data-melt-slider-thumb]` - Present on thumb

### Range
- `[data-melt-slider-range]` - Present on range

### Tick
- `[data-bounded]` - Within active range
- `[data-melt-slider-tick]` - Present on tick

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Arrow Right/Up` | Increase by step |
| `Arrow Left/Down` | Decrease by step |
| `Home` | Set to minimum |
| `End` | Set to maximum |
| `Page Up` | Large increase |
| `Page Down` | Large decrease |

## Accessibility

Follows [WAI-ARIA Slider Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/slider/):
- Proper slider role
- `aria-valuemin`, `aria-valuemax`, `aria-valuenow`
- Keyboard accessible

## Example

```svelte
<script lang="ts">
  import { createSlider, melt } from '@melt-ui/svelte';

  const {
    elements: { root, range, thumb },
    states: { value }
  } = createSlider({
    defaultValue: [50],
    min: 0,
    max: 100,
    step: 1
  });
</script>

<div class="flex flex-col gap-2">
  <label class="text-sm font-semibold text-ocean-800">
    Volume: {$value[0]}%
  </label>

  <div
    use:melt={$root}
    class="relative flex h-5 w-full items-center"
  >
    <div class="h-2 w-full rounded-full bg-ocean-200">
      <div
        use:melt={$range}
        class="h-full rounded-full bg-ocean-500"
      ></div>
    </div>

    <div
      use:melt={$thumb(0)}
      class="h-5 w-5 rounded-full bg-white border-2 border-ocean-500 shadow focus:outline-none focus:ring-2 focus:ring-ocean-500/50"
    ></div>
  </div>
</div>
```

## Styling with Tailwind

```css
/* Slider track */
[data-melt-slider] {
  @apply relative flex items-center;
}

/* Range fill */
[data-melt-slider-range] {
  @apply h-full rounded-full bg-ocean-500;
}

/* Thumb */
[data-melt-slider-thumb] {
  @apply h-5 w-5 rounded-full bg-white border-2 border-ocean-500 shadow cursor-grab;
}

[data-melt-slider-thumb]:active {
  @apply cursor-grabbing;
}

[data-melt-slider-thumb]:focus {
  @apply outline-none ring-2 ring-ocean-500/50;
}

/* Disabled state */
[data-melt-slider][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}

[data-melt-slider][data-disabled] [data-melt-slider-thumb] {
  @apply cursor-not-allowed;
}
```

## Range Slider (Two Thumbs)

```svelte
<script lang="ts">
  import { createSlider, melt } from '@melt-ui/svelte';

  const {
    elements: { root, range, thumb },
    states: { value }
  } = createSlider({
    defaultValue: [20, 80],
    min: 0,
    max: 100
  });
</script>

<div class="flex flex-col gap-2">
  <div class="flex justify-between text-sm text-ocean-600">
    <span>Min: {$value[0]}</span>
    <span>Max: {$value[1]}</span>
  </div>

  <div use:melt={$root} class="relative flex h-5 w-full items-center">
    <div class="h-2 w-full rounded-full bg-ocean-200">
      <div use:melt={$range} class="h-full rounded-full bg-ocean-500"></div>
    </div>

    <div
      use:melt={$thumb(0)}
      class="h-5 w-5 rounded-full bg-white border-2 border-ocean-500 shadow"
    ></div>
    <div
      use:melt={$thumb(1)}
      class="h-5 w-5 rounded-full bg-white border-2 border-ocean-500 shadow"
    ></div>
  </div>
</div>
```

## With Ticks

```svelte
<script lang="ts">
  import { createSlider, melt } from '@melt-ui/svelte';

  const {
    elements: { root, range, thumb, tick },
    states: { value, ticks }
  } = createSlider({
    defaultValue: [50],
    min: 0,
    max: 100,
    step: 25
  });
</script>

<div use:melt={$root} class="relative flex h-5 w-full items-center">
  <div class="h-2 w-full rounded-full bg-ocean-200">
    <div use:melt={$range} class="h-full rounded-full bg-ocean-500"></div>
  </div>

  <!-- Tick marks -->
  {#each $ticks as t, i}
    <div
      use:melt={$tick(i)}
      class="absolute h-3 w-1 -translate-x-1/2 rounded bg-ocean-300
        data-[bounded]:bg-ocean-600"
      style="left: {(t / 100) * 100}%"
    ></div>
  {/each}

  <div use:melt={$thumb(0)} class="h-5 w-5 rounded-full bg-white border-2 border-ocean-500"></div>
</div>
```

## Vertical Slider

```svelte
<script lang="ts">
  const {
    elements: { root, range, thumb }
  } = createSlider({
    orientation: 'vertical',
    defaultValue: [50]
  });
</script>

<div use:melt={$root} class="relative flex h-48 w-5 flex-col items-center">
  <div class="h-full w-2 rounded-full bg-ocean-200">
    <div use:melt={$range} class="w-full rounded-full bg-ocean-500"></div>
  </div>
  <div use:melt={$thumb(0)} class="h-5 w-5 rounded-full bg-white border-2 border-ocean-500"></div>
</div>
```

## Price Range Filter

```svelte
<script lang="ts">
  import { createSlider, melt } from '@melt-ui/svelte';

  const {
    elements: { root, range, thumb },
    states: { value }
  } = createSlider({
    defaultValue: [100, 500],
    min: 0,
    max: 1000,
    step: 10
  });

  const formatPrice = (v: number) => `$${v}`;
</script>

<div class="space-y-4">
  <div class="flex justify-between">
    <span class="text-sm font-medium text-ocean-800">Price Range</span>
    <span class="text-sm text-ocean-600">
      {formatPrice($value[0])} - {formatPrice($value[1])}
    </span>
  </div>

  <div use:melt={$root} class="relative flex h-5 w-full items-center">
    <div class="h-2 w-full rounded-full bg-ocean-200">
      <div use:melt={$range} class="h-full bg-ocean-500"></div>
    </div>
    <div use:melt={$thumb(0)} class="thumb"></div>
    <div use:melt={$thumb(1)} class="thumb"></div>
  </div>

  <div class="flex justify-between text-xs text-ocean-400">
    <span>$0</span>
    <span>$1000</span>
  </div>
</div>
```

## Controlled Value

```svelte
<script lang="ts">
  import { writable } from 'svelte/store';

  const sliderValue = writable([50]);

  const slider = createSlider({
    value: sliderValue,
    onValueChange: ({ next }) => {
      console.log('Value changed:', next);
      return next;
    },
    onValueCommitted: (value) => {
      console.log('Final value:', value);
      // Save to server, etc.
    }
  });
</script>
```
