# Date Field

A structured input field for entering dates with individual segments for day, month, and year.

## Use Cases

- Date of birth inputs
- Document date fields
- Structured date entry
- Forms requiring specific date format

## Installation

```typescript
import { createDateField } from '@melt-ui/svelte';
```

## API Reference

### createDateField Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `defaultValue` | `DateValue` | — | Initial date value |
| `value` | `Writable<DateValue>` | — | Controlled date value |
| `defaultPlaceholder` | `DateValue` | Today | Placeholder date |
| `placeholder` | `Writable<DateValue>` | — | Controlled placeholder |
| `minValue` | `DateValue` | — | Minimum selectable date |
| `maxValue` | `DateValue` | — | Maximum selectable date |
| `disabled` | `boolean` | `false` | Disable field |
| `readonly` | `boolean` | `false` | Make read-only |
| `locale` | `string` | `'en'` | Locale for formatting |
| `granularity` | `'day' \| 'hour' \| 'minute' \| 'second'` | `'day'` | Date precision |
| `hourCycle` | `12 \| 24` | — | Hour format |
| `hideTimeZone` | `boolean` | `false` | Hide timezone |
| `isDateUnavailable` | `(date) => boolean` | — | Mark dates unavailable |
| `onValueChange` | `ChangeFn` | — | Value change callback |
| `onPlaceholderChange` | `ChangeFn` | — | Placeholder callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `field` | Container element |
| `segment` | Individual date segment (day, month, year) |
| `label` | Field label |
| `hiddenInput` | Hidden input for forms |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `value` | `Writable<DateValue>` | Current date value |
| `placeholder` | `Writable<DateValue>` | Current placeholder |
| `segments` | `Readable<Segment[]>` | Date segments for rendering |
| `isInvalid` | `Readable<boolean>` | Validation state |

## Data Attributes

### Field
- `[data-invalid]` - Present when invalid
- `[data-disabled]` - Present when disabled
- `[data-melt-date-field]` - Present on field

### Segment
- `[data-segment]` - Segment type (day, month, year)
- `[data-placeholder]` - Has placeholder value
- `[data-disabled]` - Present when disabled
- `[data-melt-date-field-segment]` - Present on segments

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Arrow Up` | Increment segment value |
| `Arrow Down` | Decrement segment value |
| `Arrow Left` | Previous segment |
| `Arrow Right` | Next segment |
| `Tab` | Next segment/field |
| `0-9` | Direct number input |

## Accessibility

- Individual segments are focusable
- Screen reader announces segment type and value
- Proper ARIA labels and descriptions

## Example

```svelte
<script lang="ts">
  import { createDateField, melt } from '@melt-ui/svelte';
  import { CalendarDays } from 'lucide-svelte';

  const {
    elements: { field, segment, label },
    states: { segments, value, isInvalid }
  } = createDateField({
    locale: 'en-GB' // DD/MM/YYYY format
  });
</script>

<div class="flex flex-col gap-1">
  <span use:melt={$label} class="text-sm font-semibold text-ocean-800">
    Start Date
  </span>

  <div
    use:melt={$field}
    class="flex items-center gap-1 px-3 py-2 rounded-lg border-2 border-ocean-500/40 bg-white
      focus-within:ring-2 focus-within:ring-ocean-500/30 focus-within:border-ocean-500
      data-[invalid]:border-red-500"
  >
    <CalendarDays class="h-5 w-5 text-ocean-400 mr-2" />

    {#each $segments as seg}
      {#if seg.part === 'literal'}
        <span class="text-ocean-400">{seg.value}</span>
      {:else}
        <span
          use:melt={$segment(seg.part)}
          class="rounded px-0.5 tabular-nums outline-none
            focus:bg-ocean-100 focus:ring-2 focus:ring-ocean-500
            data-[placeholder]:text-ocean-400"
        >
          {seg.value}
        </span>
      {/if}
    {/each}
  </div>

  {#if $isInvalid}
    <span class="text-sm text-red-500">Please enter a valid date</span>
  {/if}
</div>
```

## Styling with Tailwind

```css
/* Field container */
[data-melt-date-field] {
  @apply flex items-center gap-0.5;
}

/* Segment focus */
[data-melt-date-field-segment]:focus {
  @apply bg-ocean-100 outline-none ring-2 ring-ocean-500;
}

/* Placeholder styling */
[data-melt-date-field-segment][data-placeholder] {
  @apply text-ocean-400;
}

/* Invalid state */
[data-melt-date-field][data-invalid] {
  @apply border-red-500;
}

/* Disabled state */
[data-melt-date-field][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}
```

## With Min/Max Dates

```typescript
import { CalendarDate } from '@internationalized/date';

const dateField = createDateField({
  minValue: new CalendarDate(2024, 1, 1),
  maxValue: new CalendarDate(2024, 12, 31)
});
```

## With Unavailable Dates

```typescript
// Disable weekends
const dateField = createDateField({
  isDateUnavailable: (date) => {
    const day = date.toDate('UTC').getDay();
    return day === 0 || day === 6;
  }
});
```

## With Time

```svelte
<script lang="ts">
  const {
    elements: { field, segment },
    states: { segments }
  } = createDateField({
    granularity: 'minute',
    hourCycle: 24
  });
</script>

<!-- Segments will include hour and minute -->
<div use:melt={$field}>
  {#each $segments as seg}
    <!-- Renders: DD / MM / YYYY HH : MM -->
  {/each}
</div>
```

## Form Integration

```svelte
<script lang="ts">
  const {
    elements: { field, segment, hiddenInput },
    states: { value }
  } = createDateField();
</script>

<form onsubmit={handleSubmit}>
  <div use:melt={$field}>
    {#each $segments as seg}
      <!-- segments -->
    {/each}
  </div>
  <input use:melt={$hiddenInput} name="startDate" />
  <button type="submit">Submit</button>
</form>
```

## Controlled Value

```svelte
<script lang="ts">
  import { writable } from 'svelte/store';
  import { CalendarDate } from '@internationalized/date';

  const dateValue = writable(new CalendarDate(2024, 6, 15));

  const { states: { value } } = createDateField({
    value: dateValue,
    onValueChange: ({ next }) => {
      console.log('Date changed:', next);
      return next;
    }
  });
</script>
```
