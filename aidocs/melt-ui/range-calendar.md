# Range Calendar

A calendar component for selecting date ranges with visual range highlighting.

## Use Cases

- Vacation date selection
- Report period filters
- Booking date ranges
- Analytics date selection

## Installation

```typescript
import { createRangeCalendar } from '@melt-ui/svelte';
```

## API Reference

### createRangeCalendar Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `defaultValue` | `DateRange` | — | Initial date range |
| `value` | `Writable<DateRange>` | — | Controlled date range |
| `defaultPlaceholder` | `DateValue` | Today | Initial displayed month |
| `placeholder` | `Writable<DateValue>` | — | Controlled placeholder |
| `minValue` | `DateValue` | — | Minimum selectable date |
| `maxValue` | `DateValue` | — | Maximum selectable date |
| `disabled` | `boolean` | `false` | Disable calendar |
| `readonly` | `boolean` | `false` | Make read-only |
| `locale` | `string` | `'en'` | Locale for formatting |
| `fixedWeeks` | `boolean` | `false` | Always show 6 weeks |
| `weekStartsOn` | `0-6` | `0` | First day of week |
| `numberOfMonths` | `number` | `1` | Months to display |
| `isDateDisabled` | `(date) => boolean` | — | Disable specific dates |
| `isDateUnavailable` | `(date) => boolean` | — | Mark dates unavailable |
| `calendarLabel` | `string` | — | Accessible label |
| `onValueChange` | `ChangeFn` | — | Value change callback |
| `onPlaceholderChange` | `ChangeFn` | — | Placeholder callback |

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
| `calendar` | Root calendar container |
| `heading` | Month/year heading |
| `grid` | Calendar grid wrapper |
| `cell` | Individual day cell |
| `prevButton` | Previous month button |
| `nextButton` | Next month button |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `value` | `Writable<DateRange>` | Selected date range |
| `placeholder` | `Writable<DateValue>` | Displayed month |
| `months` | `Readable<Month[]>` | Month data |
| `weekdays` | `Readable<string[]>` | Weekday labels |
| `isInvalid` | `Readable<boolean>` | Validation state |

## Data Attributes

### Cell
- `[data-selected]` - Date is in selected range
- `[data-selection-start]` - Range start date
- `[data-selection-end]` - Range end date
- `[data-highlighted]` - Between start and end (hover preview)
- `[data-today]` - Today's date
- `[data-disabled]` - Date is disabled
- `[data-unavailable]` - Date is unavailable
- `[data-outside-month]` - Outside current month
- `[data-outside-visible-months]` - Outside all visible months
- `[data-focused]` - Has keyboard focus

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Arrow Keys` | Navigate calendar |
| `Enter` / `Space` | Select date (start then end) |
| `Page Up/Down` | Previous/next month |
| `Home` | Start of week |
| `End` | End of week |

## Accessibility

- Full keyboard navigation for range selection
- Screen reader announces selection state
- Proper ARIA attributes for date range

## Example

```svelte
<script lang="ts">
  import { createRangeCalendar, melt } from '@melt-ui/svelte';
  import { ChevronLeft, ChevronRight } from 'lucide-svelte';

  const {
    elements: { calendar, heading, grid, cell, prevButton, nextButton },
    states: { value, months, weekdays }
  } = createRangeCalendar({
    locale: 'en-GB',
    weekStartsOn: 1,
    numberOfMonths: 2
  });
</script>

<div use:melt={$calendar} class="bg-white rounded-xl shadow-lg p-4">
  <div class="flex gap-8">
    {#each $months as month, i}
      <div class="flex-1">
        <!-- Header -->
        <div class="flex items-center justify-between mb-4">
          {#if i === 0}
            <button use:melt={$prevButton} class="p-1 rounded hover:bg-ocean-100">
              <ChevronLeft class="h-5 w-5 text-ocean-600" />
            </button>
          {:else}
            <div class="w-7"></div>
          {/if}

          <div use:melt={$heading} class="font-semibold text-ocean-800">
            {month.value.toDate('UTC').toLocaleString('en', { month: 'long', year: 'numeric' })}
          </div>

          {#if i === $months.length - 1}
            <button use:melt={$nextButton} class="p-1 rounded hover:bg-ocean-100">
              <ChevronRight class="h-5 w-5 text-ocean-600" />
            </button>
          {:else}
            <div class="w-7"></div>
          {/if}
        </div>

        <!-- Calendar Grid -->
        <table use:melt={$grid} class="w-full">
          <thead>
            <tr>
              {#each $weekdays as day}
                <th class="text-sm font-medium text-ocean-500 pb-2 w-10">{day}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each month.weeks as week}
              <tr>
                {#each week as day}
                  <td class="p-0.5">
                    <div
                      use:melt={$cell(day, month.value)}
                      class="h-10 w-10 flex items-center justify-center text-sm cursor-pointer transition-all
                        data-[outside-month]:text-ocean-300
                        data-[selection-start]:bg-ocean-500 data-[selection-start]:text-white data-[selection-start]:rounded-l-lg
                        data-[selection-end]:bg-ocean-500 data-[selection-end]:text-white data-[selection-end]:rounded-r-lg
                        data-[highlighted]:bg-ocean-100
                        data-[today]:ring-2 data-[today]:ring-ocean-400
                        data-[disabled]:opacity-40 data-[disabled]:cursor-not-allowed
                        data-[unavailable]:line-through data-[unavailable]:text-ocean-300
                        hover:bg-ocean-200"
                    >
                      {day.day}
                    </div>
                  </td>
                {/each}
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/each}
  </div>

  <!-- Selection Summary -->
  {#if $value?.start && $value?.end}
    <div class="mt-4 pt-4 border-t border-ocean-200 text-sm text-ocean-600">
      Selected: {$value.start.toString()} to {$value.end.toString()}
    </div>
  {/if}
</div>
```

## Styling with Tailwind

```css
/* Selection start */
[data-melt-range-calendar-cell][data-selection-start] {
  @apply bg-ocean-500 text-white rounded-l-lg;
}

/* Selection end */
[data-melt-range-calendar-cell][data-selection-end] {
  @apply bg-ocean-500 text-white rounded-r-lg;
}

/* Highlighted (between start and end) */
[data-melt-range-calendar-cell][data-highlighted] {
  @apply bg-ocean-100;
}

/* Today indicator */
[data-melt-range-calendar-cell][data-today] {
  @apply ring-2 ring-ocean-400;
}

/* Disabled dates */
[data-melt-range-calendar-cell][data-disabled] {
  @apply opacity-40 cursor-not-allowed;
}

/* Unavailable dates */
[data-melt-range-calendar-cell][data-unavailable] {
  @apply line-through text-ocean-300;
}

/* Selection (both start and end on same day) */
[data-melt-range-calendar-cell][data-selection-start][data-selection-end] {
  @apply rounded-lg;
}
```

## VacayTracker Integration

```svelte
<script lang="ts">
  import { createRangeCalendar, melt } from '@melt-ui/svelte';
  import { today, getLocalTimeZone } from '@internationalized/date';
  import { auth } from '$lib/stores/auth.svelte';

  const {
    elements: { calendar, cell },
    states: { value }
  } = createRangeCalendar({
    locale: 'en-GB',
    weekStartsOn: 1,
    numberOfMonths: 2,
    minValue: today(getLocalTimeZone()),
    isDateDisabled: (date) => {
      // Disable weekends based on company policy
      const day = date.toDate('UTC').getDay();
      return day === 0 || day === 6;
    }
  });

  // Calculate business days in selected range
  const businessDays = $derived(() => {
    if (!$value?.start || !$value?.end) return 0;
    let days = 0;
    let current = $value.start;
    while (current.compare($value.end) <= 0) {
      const day = current.toDate('UTC').getDay();
      if (day !== 0 && day !== 6) days++;
      current = current.add({ days: 1 });
    }
    return days;
  });

  // Check against user balance
  const exceedsBalance = $derived(
    businessDays() > (auth.user?.vacationBalance ?? 0)
  );
</script>

<div use:melt={$calendar}>
  <!-- Calendar content -->
</div>

{#if $value?.start && $value?.end}
  <div class="mt-4 p-4 rounded-lg bg-ocean-50">
    <div class="flex justify-between text-sm">
      <span class="text-ocean-600">Business days:</span>
      <span class="font-medium text-ocean-800">{businessDays()}</span>
    </div>
    {#if exceedsBalance}
      <p class="text-sm text-red-500 mt-2">
        Exceeds your available balance of {auth.user?.vacationBalance} days
      </p>
    {/if}
  </div>
{/if}
```

## With Preset Ranges

```svelte
<script lang="ts">
  import { today, getLocalTimeZone } from '@internationalized/date';

  const presets = [
    {
      label: 'This Week',
      getValue: () => {
        const now = today(getLocalTimeZone());
        return {
          start: now.subtract({ days: now.dayOfWeek - 1 }),
          end: now.add({ days: 7 - now.dayOfWeek })
        };
      }
    },
    {
      label: 'Next Week',
      getValue: () => {
        const now = today(getLocalTimeZone());
        const nextMonday = now.add({ days: 8 - now.dayOfWeek });
        return {
          start: nextMonday,
          end: nextMonday.add({ days: 4 })
        };
      }
    }
  ];

  function selectPreset(preset: typeof presets[0]) {
    value.set(preset.getValue());
  }
</script>

<div class="flex gap-4">
  <div class="w-40 space-y-1">
    {#each presets as preset}
      <button
        onclick={() => selectPreset(preset)}
        class="w-full text-left px-3 py-2 text-sm rounded-lg hover:bg-ocean-100 transition-colors"
      >
        {preset.label}
      </button>
    {/each}
  </div>
  <div use:melt={$calendar}>
    <!-- Calendar grids -->
  </div>
</div>
```

## Single Month

```typescript
const calendar = createRangeCalendar({
  numberOfMonths: 1 // Just one month
});
```

## Minimum Range

```svelte
<script lang="ts">
  // Enforce minimum 3-day range
  const { states: { value } } = createRangeCalendar({
    onValueChange: ({ curr, next }) => {
      if (next?.start && next?.end) {
        const days = next.end.compare(next.start);
        if (days < 3) {
          // Keep current value if range too short
          return curr;
        }
      }
      return next;
    }
  });
</script>
```
