# Date Range Field

A structured input field for entering date ranges with start and end date segments.

## Use Cases

- Vacation request date ranges
- Report date filters
- Booking period inputs
- Project duration fields

## Installation

```typescript
import { createDateRangeField } from '@melt-ui/svelte';
```

## API Reference

### createDateRangeField Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `defaultValue` | `DateRange` | — | Initial date range |
| `value` | `Writable<DateRange>` | — | Controlled date range |
| `defaultPlaceholder` | `DateValue` | Today | Placeholder date |
| `placeholder` | `Writable<DateValue>` | — | Controlled placeholder |
| `minValue` | `DateValue` | — | Minimum selectable date |
| `maxValue` | `DateValue` | — | Maximum selectable date |
| `disabled` | `boolean` | `false` | Disable field |
| `readonly` | `boolean` | `false` | Make read-only |
| `locale` | `string` | `'en'` | Locale for formatting |
| `granularity` | `'day' \| 'hour' \| 'minute' \| 'second'` | `'day'` | Date precision |
| `isDateUnavailable` | `(date) => boolean` | — | Mark dates unavailable |
| `onValueChange` | `ChangeFn` | — | Value change callback |

### DateRange Type

```typescript
type DateRange = {
  start: DateValue;
  end: DateValue;
};
```

### Returned Elements

| Element | Description |
|---------|-------------|
| `field` | Container element |
| `startSegment` | Start date segment |
| `endSegment` | End date segment |
| `label` | Field label |
| `startHiddenInput` | Hidden input for start date |
| `endHiddenInput` | Hidden input for end date |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `value` | `Writable<DateRange>` | Current date range |
| `startSegments` | `Readable<Segment[]>` | Start date segments |
| `endSegments` | `Readable<Segment[]>` | End date segments |
| `isInvalid` | `Readable<boolean>` | Validation state |

## Data Attributes

### Field
- `[data-invalid]` - Present when invalid
- `[data-disabled]` - Present when disabled
- `[data-melt-date-range-field]` - Present on field

### Segment
- `[data-segment]` - Segment type
- `[data-placeholder]` - Has placeholder value
- `[data-disabled]` - Present when disabled

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Arrow Up` | Increment segment |
| `Arrow Down` | Decrement segment |
| `Arrow Left` | Previous segment |
| `Arrow Right` | Next segment |
| `Tab` | Move between start/end |
| `0-9` | Direct number input |

## Accessibility

- Both start and end fields are accessible
- Screen reader announces which part of the range
- Proper ARIA labels and descriptions

## Example

```svelte
<script lang="ts">
  import { createDateRangeField, melt } from '@melt-ui/svelte';
  import { CalendarDays } from 'lucide-svelte';

  const {
    elements: { field, startSegment, endSegment, label },
    states: { startSegments, endSegments, value, isInvalid }
  } = createDateRangeField({
    locale: 'en-GB'
  });
</script>

<div class="flex flex-col gap-1">
  <span use:melt={$label} class="text-sm font-semibold text-ocean-800">
    Vacation Dates
  </span>

  <div
    use:melt={$field}
    class="flex items-center gap-2 px-3 py-2 rounded-lg border-2 border-ocean-500/40 bg-white
      focus-within:ring-2 focus-within:ring-ocean-500/30 focus-within:border-ocean-500
      data-[invalid]:border-red-500"
  >
    <CalendarDays class="h-5 w-5 text-ocean-400" />

    <!-- Start Date -->
    <div class="flex items-center gap-0.5">
      {#each $startSegments as seg}
        {#if seg.part === 'literal'}
          <span class="text-ocean-400">{seg.value}</span>
        {:else}
          <span
            use:melt={$startSegment(seg.part)}
            class="rounded px-0.5 tabular-nums outline-none
              focus:bg-ocean-100 focus:ring-2 focus:ring-ocean-500
              data-[placeholder]:text-ocean-400"
          >
            {seg.value}
          </span>
        {/if}
      {/each}
    </div>

    <span class="text-ocean-400 mx-2">to</span>

    <!-- End Date -->
    <div class="flex items-center gap-0.5">
      {#each $endSegments as seg}
        {#if seg.part === 'literal'}
          <span class="text-ocean-400">{seg.value}</span>
        {:else}
          <span
            use:melt={$endSegment(seg.part)}
            class="rounded px-0.5 tabular-nums outline-none
              focus:bg-ocean-100 focus:ring-2 focus:ring-ocean-500
              data-[placeholder]:text-ocean-400"
          >
            {seg.value}
          </span>
        {/if}
      {/each}
    </div>
  </div>

  {#if $isInvalid}
    <span class="text-sm text-red-500">
      Please enter a valid date range
    </span>
  {/if}
</div>
```

## Styling with Tailwind

```css
/* Field container */
[data-melt-date-range-field] {
  @apply flex items-center gap-2;
}

/* Segment focus */
[data-melt-date-range-field] [data-segment]:focus {
  @apply bg-ocean-100 outline-none ring-2 ring-ocean-500;
}

/* Placeholder */
[data-melt-date-range-field] [data-placeholder] {
  @apply text-ocean-400;
}

/* Invalid state */
[data-melt-date-range-field][data-invalid] {
  @apply border-red-500;
}
```

## VacayTracker Integration

```svelte
<script lang="ts">
  import { createDateRangeField, melt } from '@melt-ui/svelte';
  import { today, getLocalTimeZone } from '@internationalized/date';

  const {
    elements: { field, startSegment, endSegment },
    states: { value, isInvalid }
  } = createDateRangeField({
    locale: 'en-GB',
    minValue: today(getLocalTimeZone()),
    isDateUnavailable: (date) => {
      // Disable weekends based on company policy
      const day = date.toDate('UTC').getDay();
      return day === 0 || day === 6;
    }
  });

  // Calculate estimated business days
  const estimatedDays = $derived(() => {
    if (!$value?.start || !$value?.end) return 0;
    // Business days calculation
    let days = 0;
    let current = $value.start;
    while (current.compare($value.end) <= 0) {
      const day = current.toDate('UTC').getDay();
      if (day !== 0 && day !== 6) days++;
      current = current.add({ days: 1 });
    }
    return days;
  });
</script>

<div class="space-y-2">
  <div use:melt={$field}>
    <!-- Date range inputs -->
  </div>

  {#if estimatedDays() > 0}
    <p class="text-sm text-ocean-600">
      Estimated: {estimatedDays()} business days
    </p>
  {/if}
</div>
```

## Form Submission

```svelte
<script lang="ts">
  const {
    elements: { field, startHiddenInput, endHiddenInput }
  } = createDateRangeField();
</script>

<form onsubmit={handleSubmit}>
  <div use:melt={$field}>
    <!-- Visible segments -->
  </div>
  <input use:melt={$startHiddenInput} name="startDate" />
  <input use:melt={$endHiddenInput} name="endDate" />
  <button type="submit">Submit Request</button>
</form>
```

## Controlled Value

```svelte
<script lang="ts">
  import { writable } from 'svelte/store';
  import { CalendarDate } from '@internationalized/date';

  const dateRange = writable({
    start: new CalendarDate(2024, 6, 15),
    end: new CalendarDate(2024, 6, 20)
  });

  const { states: { value } } = createDateRangeField({
    value: dateRange,
    onValueChange: ({ next }) => {
      console.log('Range changed:', next);
      return next;
    }
  });
</script>
```
