# Toggle Group

A group of toggle buttons where one or multiple can be selected.

## Use Cases

- Text alignment options
- View mode selection
- Filter groups
- Segmented controls

## Installation

```typescript
import { createToggleGroup } from '@melt-ui/svelte';
```

## API Reference

### createToggleGroup Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `type` | `'single' \| 'multiple'` | `'single'` | Selection type |
| `disabled` | `boolean` | `false` | Disable all toggles |
| `loop` | `boolean` | `true` | Loop keyboard navigation |
| `orientation` | `'horizontal' \| 'vertical'` | `'horizontal'` | Layout direction |
| `rovingFocus` | `boolean` | `true` | Enable roving focus |
| `defaultValue` | `string \| string[]` | — | Initial value(s) |
| `value` | `Writable<string \| string[]>` | — | Controlled value(s) |
| `onValueChange` | `ChangeFn` | — | Value change callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Group container |
| `item` | Individual toggle item |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `value` | `Writable<string \| string[]>` | Currently selected value(s) |

### Returned Helpers

| Helper | Description |
|--------|-------------|
| `isPressed(value)` | Check if an item is pressed |

## Data Attributes

### Root
- `[data-orientation]` - `'horizontal'` or `'vertical'`
- `[data-melt-toggle-group]` - Present on root

### Item
- `[data-state]` - `'on'` or `'off'`
- `[data-disabled]` - Present when disabled
- `[data-value]` - Item value
- `[data-orientation]` - Inherited orientation
- `[data-melt-toggle-group-item]` - Present on items

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Tab` | Focus group |
| `Arrow Left/Up` | Previous item |
| `Arrow Right/Down` | Next item |
| `Home` | First item |
| `End` | Last item |
| `Space` / `Enter` | Toggle item |

## Accessibility

- Uses proper button group semantics
- `aria-pressed` indicates selection state
- Full keyboard navigation support

## Example

```svelte
<script lang="ts">
  import { createToggleGroup, melt } from '@melt-ui/svelte';
  import { AlignLeft, AlignCenter, AlignRight, AlignJustify } from 'lucide-svelte';

  const {
    elements: { root, item },
    states: { value },
    helpers: { isPressed }
  } = createToggleGroup({
    type: 'single',
    defaultValue: 'left'
  });

  const alignments = [
    { value: 'left', icon: AlignLeft, label: 'Align left' },
    { value: 'center', icon: AlignCenter, label: 'Align center' },
    { value: 'right', icon: AlignRight, label: 'Align right' },
    { value: 'justify', icon: AlignJustify, label: 'Justify' }
  ];
</script>

<div use:melt={$root} class="flex p-1 bg-ocean-100 rounded-lg">
  {#each alignments as { value: v, icon: Icon, label }}
    <button
      use:melt={$item(v)}
      class="p-2 rounded transition-colors
        data-[state=off]:text-ocean-600 data-[state=off]:hover:bg-ocean-200
        data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
      title={label}
    >
      <Icon class="h-4 w-4" />
    </button>
  {/each}
</div>
```

## Styling with Tailwind

```css
/* Group container */
[data-melt-toggle-group] {
  @apply flex p-1 rounded-lg;
}

/* Toggle item off */
[data-melt-toggle-group-item][data-state="off"] {
  @apply text-ocean-600 hover:bg-ocean-200;
}

/* Toggle item on */
[data-melt-toggle-group-item][data-state="on"] {
  @apply bg-ocean-500 text-white;
}

/* Disabled item */
[data-melt-toggle-group-item][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}

/* Focus ring */
[data-melt-toggle-group-item]:focus-visible {
  @apply ring-2 ring-ocean-500 ring-offset-2;
}
```

## Multiple Selection

```svelte
<script lang="ts">
  import { createToggleGroup, melt } from '@melt-ui/svelte';
  import { Bold, Italic, Underline } from 'lucide-svelte';

  const {
    elements: { root, item },
    states: { value },
    helpers: { isPressed }
  } = createToggleGroup({
    type: 'multiple',
    defaultValue: ['bold']
  });

  const formats = [
    { value: 'bold', icon: Bold },
    { value: 'italic', icon: Italic },
    { value: 'underline', icon: Underline }
  ];
</script>

<div use:melt={$root} class="flex gap-1">
  {#each formats as { value: v, icon: Icon }}
    <button
      use:melt={$item(v)}
      class="p-2 rounded border transition-colors
        data-[state=off]:border-ocean-200 data-[state=off]:text-ocean-600
        data-[state=on]:border-ocean-500 data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
    >
      <Icon class="h-4 w-4" />
    </button>
  {/each}
</div>

<!-- Current selections: {$value.join(', ')} -->
```

## View Mode Selector

```svelte
<script lang="ts">
  import { createToggleGroup, melt } from '@melt-ui/svelte';
  import { Grid, List, Calendar } from 'lucide-svelte';

  const {
    elements: { root, item },
    states: { value }
  } = createToggleGroup({
    type: 'single',
    defaultValue: 'list'
  });

  const views = [
    { value: 'list', icon: List, label: 'List view' },
    { value: 'grid', icon: Grid, label: 'Grid view' },
    { value: 'calendar', icon: Calendar, label: 'Calendar view' }
  ];
</script>

<div use:melt={$root} class="inline-flex p-1 bg-ocean-100 rounded-lg">
  {#each views as { value: v, icon: Icon, label }}
    <button
      use:melt={$item(v)}
      class="flex items-center gap-2 px-3 py-1.5 rounded text-sm font-medium transition-colors
        data-[state=off]:text-ocean-600
        data-[state=on]:bg-white data-[state=on]:text-ocean-800 data-[state=on]:shadow-sm"
      title={label}
    >
      <Icon class="h-4 w-4" />
      <span class="hidden sm:inline">{label}</span>
    </button>
  {/each}
</div>
```

## Period Selector

```svelte
<script lang="ts">
  const { elements: { root, item } } = createToggleGroup({
    type: 'single',
    defaultValue: 'monthly'
  });

  const periods = ['Daily', 'Weekly', 'Monthly', 'Yearly'];
</script>

<div use:melt={$root} class="inline-flex p-1 bg-ocean-100 rounded-lg">
  {#each periods as period}
    <button
      use:melt={$item(period.toLowerCase())}
      class="px-4 py-2 rounded text-sm font-medium transition-all
        data-[state=off]:text-ocean-600
        data-[state=on]:bg-white data-[state=on]:text-ocean-800 data-[state=on]:shadow-sm"
    >
      {period}
    </button>
  {/each}
</div>
```

## Vertical Group

```svelte
<script lang="ts">
  const { elements: { root, item } } = createToggleGroup({
    orientation: 'vertical',
    type: 'single',
    defaultValue: 'option1'
  });
</script>

<div use:melt={$root} class="flex flex-col p-1 bg-ocean-100 rounded-lg w-48">
  {#each options as opt}
    <button
      use:melt={$item(opt.value)}
      class="px-4 py-2 text-left text-sm rounded transition-colors
        data-[state=off]:text-ocean-600 data-[state=off]:hover:bg-ocean-200
        data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
    >
      {opt.label}
    </button>
  {/each}
</div>
```

## Filter Chips

```svelte
<script lang="ts">
  const { elements: { root, item }, states: { value } } = createToggleGroup({
    type: 'multiple',
    defaultValue: ['approved']
  });

  const statuses = ['pending', 'approved', 'rejected'];
</script>

<div use:melt={$root} class="flex flex-wrap gap-2">
  {#each statuses as status}
    <button
      use:melt={$item(status)}
      class="px-4 py-1.5 rounded-full text-sm font-medium border transition-colors capitalize
        data-[state=off]:border-ocean-300 data-[state=off]:text-ocean-600 data-[state=off]:hover:border-ocean-400
        data-[state=on]:border-ocean-500 data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
    >
      {status}
    </button>
  {/each}
</div>
```

## Controlled State

```svelte
<script lang="ts">
  import { writable } from 'svelte/store';

  const selectedValue = writable('option1');

  const { elements: { root, item } } = createToggleGroup({
    type: 'single',
    value: selectedValue,
    onValueChange: ({ next }) => {
      console.log('Selection changed:', next);
      return next;
    }
  });
</script>
```
