# Calendar

A date grid component for navigating and selecting dates.

## Use Cases

- Date selection in forms
- Event calendars
- Vacation calendar displays
- Scheduling interfaces

## Installation

```typescript
import { createCalendar } from '@melt-ui/svelte';
```

## API Reference

### createCalendar Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `defaultValue` | `DateValue` | — | Initial selected date |
| `value` | `Writable<DateValue>` | — | Controlled selected date |
| `defaultPlaceholder` | `DateValue` | Today | Initial displayed month |
| `placeholder` | `Writable<DateValue>` | — | Controlled displayed month |
| `minValue` | `DateValue` | — | Minimum selectable date |
| `maxValue` | `DateValue` | — | Maximum selectable date |
| `disabled` | `boolean` | `false` | Disable calendar interaction |
| `readonly` | `boolean` | `false` | Prevent date selection |
| `locale` | `string` | `'en'` | Locale for formatting |
| `fixedWeeks` | `boolean` | `false` | Always show 6 weeks |
| `weekStartsOn` | `0-6` | `0` | First day of week (0=Sunday) |
| `isDateDisabled` | `(date) => boolean` | — | Disable specific dates |
| `isDateUnavailable` | `(date) => boolean` | — | Mark dates unavailable |
| `calendarLabel` | `string` | `'Event Date'` | Accessible label |
| `onValueChange` | `ChangeFn` | — | Value change callback |
| `onPlaceholderChange` | `ChangeFn` | — | Displayed month callback |

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
| `value` | `Writable<DateValue>` | Selected date |
| `placeholder` | `Writable<DateValue>` | Displayed month |
| `months` | `Readable<Month[]>` | Month data for rendering |
| `weekdays` | `Readable<string[]>` | Weekday labels |

### Returned Helpers

| Helper | Description |
|--------|-------------|
| `isDateDisabled` | Check if date is disabled |
| `isDateUnavailable` | Check if date is unavailable |
| `isDateSelected` | Check if date is selected |

## Data Attributes

### Cell
- `[data-disabled]` - Date is disabled
- `[data-unavailable]` - Date is unavailable
- `[data-selected]` - Date is selected
- `[data-today]` - Date is today
- `[data-outside-month]` - Date is outside current month
- `[data-focused]` - Date has keyboard focus

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Arrow Up` | Previous week |
| `Arrow Down` | Next week |
| `Arrow Left` | Previous day |
| `Arrow Right` | Next day |
| `Home` | Start of week |
| `End` | End of week |
| `Page Up` | Previous month |
| `Page Down` | Next month |
| `Enter` / `Space` | Select focused date |

## Accessibility

Follows the [WAI-ARIA Calendar Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/dialog-modal/examples/datepicker-dialog/):
- Grid navigation with arrow keys
- Proper ARIA labels and roles
- Screen reader announcements for month changes

## Example

```svelte
<script lang="ts">
  import { createCalendar, melt } from '@melt-ui/svelte';
  import { ChevronLeft, ChevronRight } from 'lucide-svelte';

  const {
    elements: { calendar, heading, grid, cell, prevButton, nextButton },
    states: { months, weekdays }
  } = createCalendar({
    weekStartsOn: 1, // Monday
    locale: 'en-GB'
  });
</script>

<div use:melt={$calendar} class="w-full max-w-sm bg-white rounded-xl shadow-lg p-4">
  <!-- Header -->
  <div class="flex items-center justify-between mb-4">
    <button
      use:melt={$prevButton}
      class="p-2 rounded-lg hover:bg-ocean-100 text-ocean-600"
    >
      <ChevronLeft class="h-5 w-5" />
    </button>
    <div use:melt={$heading} class="font-semibold text-ocean-800">
      {$months[0].value.toDate('UTC').toLocaleString('en', { month: 'long', year: 'numeric' })}
    </div>
    <button
      use:melt={$nextButton}
      class="p-2 rounded-lg hover:bg-ocean-100 text-ocean-600"
    >
      <ChevronRight class="h-5 w-5" />
    </button>
  </div>

  <!-- Grid -->
  {#each $months as month}
    <table use:melt={$grid} class="w-full">
      <thead>
        <tr>
          {#each $weekdays as day}
            <th class="text-sm font-medium text-ocean-500 pb-2">{day}</th>
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
                  class="h-9 w-9 flex items-center justify-center rounded-lg text-sm
                    data-[outside-month]:text-ocean-300
                    data-[selected]:bg-ocean-500 data-[selected]:text-white
                    data-[today]:ring-2 data-[today]:ring-ocean-400
                    data-[disabled]:opacity-40 data-[disabled]:cursor-not-allowed
                    hover:bg-ocean-100 cursor-pointer transition-colors"
                >
                  {day.day}
                </div>
              </td>
            {/each}
          </tr>
        {/each}
      </tbody>
    </table>
  {/each}
</div>
```

## Styling with Tailwind

```css
/* Selected date */
[data-melt-calendar-cell][data-selected] {
  @apply bg-ocean-500 text-white font-medium;
}

/* Today indicator */
[data-melt-calendar-cell][data-today] {
  @apply ring-2 ring-ocean-400;
}

/* Outside month (faded) */
[data-melt-calendar-cell][data-outside-month] {
  @apply text-ocean-300;
}

/* Disabled dates */
[data-melt-calendar-cell][data-disabled] {
  @apply opacity-40 cursor-not-allowed;
}

/* Hover state */
[data-melt-calendar-cell]:not([data-disabled]):hover {
  @apply bg-ocean-100;
}

/* Focused state */
[data-melt-calendar-cell][data-focused] {
  @apply ring-2 ring-ocean-500;
}
```

## Disabling Dates

```typescript
// Disable weekends
const calendar = createCalendar({
  isDateDisabled: (date) => {
    const day = date.toDate('UTC').getDay();
    return day === 0 || day === 6; // Sunday or Saturday
  }
});

// Disable past dates
const calendar = createCalendar({
  isDateDisabled: (date) => {
    return date.compare(today(getLocalTimeZone())) < 0;
  }
});
```

## Multiple Months

```svelte
<script lang="ts">
  const { states: { months } } = createCalendar({
    numberOfMonths: 2
  });
</script>

<div class="flex gap-4">
  {#each $months as month}
    <!-- Render each month grid -->
  {/each}
</div>
```
