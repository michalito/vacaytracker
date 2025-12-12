# Date Picker

A popup calendar interface for selecting dates, combining a date field trigger with a calendar popover.

## Use Cases

- Form date inputs
- Event scheduling
- Vacation request dates
- Booking systems

## Installation

```typescript
import { createDatePicker } from '@melt-ui/svelte';
```

## API Reference

### createDatePicker Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `defaultValue` | `DateValue` | — | Initial selected date |
| `value` | `Writable<DateValue>` | — | Controlled date value |
| `defaultPlaceholder` | `DateValue` | Today | Initial displayed month |
| `placeholder` | `Writable<DateValue>` | — | Controlled placeholder |
| `minValue` | `DateValue` | — | Minimum selectable date |
| `maxValue` | `DateValue` | — | Maximum selectable date |
| `disabled` | `boolean` | `false` | Disable picker |
| `readonly` | `boolean` | `false` | Make read-only |
| `locale` | `string` | `'en'` | Locale for formatting |
| `isDateDisabled` | `(date) => boolean` | — | Disable specific dates |
| `isDateUnavailable` | `(date) => boolean` | — | Mark dates unavailable |
| `weekStartsOn` | `0-6` | `0` | First day of week |
| `fixedWeeks` | `boolean` | `false` | Always show 6 weeks |
| `calendarLabel` | `string` | — | Calendar accessible label |
| `closeOnOutsideClick` | `boolean` | `true` | Close on click outside |
| `portal` | `string \| HTMLElement \| null` | `'body'` | Portal target |
| `forceVisible` | `boolean` | `false` | Force visibility |
| `positioning` | `FloatingConfig` | — | Popup positioning |
| `onValueChange` | `ChangeFn` | — | Value change callback |
| `onOpenChange` | `ChangeFn` | — | Open state callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `trigger` | Button to open calendar |
| `content` | Calendar popup container |
| `calendar` | Calendar grid |
| `cell` | Individual day cell |
| `heading` | Month/year heading |
| `grid` | Calendar grid wrapper |
| `prevButton` | Previous month button |
| `nextButton` | Next month button |
| `field` | Date field container |
| `segment` | Date segment |
| `label` | Field label |
| `hiddenInput` | Hidden form input |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `open` | `Writable<boolean>` | Popup visibility |
| `value` | `Writable<DateValue>` | Selected date |
| `placeholder` | `Writable<DateValue>` | Displayed month |
| `months` | `Readable<Month[]>` | Month data |
| `weekdays` | `Readable<string[]>` | Weekday labels |
| `isInvalid` | `Readable<boolean>` | Validation state |

## Data Attributes

### Cell
- `[data-selected]` - Date is selected
- `[data-today]` - Date is today
- `[data-disabled]` - Date is disabled
- `[data-outside-month]` - Outside current month
- `[data-focused]` - Has keyboard focus

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Space` / `Enter` | Open picker / Select date |
| `Arrow Keys` | Navigate calendar |
| `Page Up/Down` | Previous/next month |
| `Escape` | Close picker |

## Accessibility

Follows [WAI-ARIA Date Picker Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/dialog-modal/examples/datepicker-dialog/):
- Full keyboard navigation
- Screen reader support
- Focus management

## Example

```svelte
<script lang="ts">
  import { createDatePicker, melt } from '@melt-ui/svelte';
  import { CalendarDays, ChevronLeft, ChevronRight } from 'lucide-svelte';

  const {
    elements: {
      trigger,
      content,
      calendar,
      cell,
      heading,
      grid,
      prevButton,
      nextButton,
      field,
      segment,
      label
    },
    states: { open, value, months, weekdays, segments }
  } = createDatePicker({
    locale: 'en-GB',
    weekStartsOn: 1,
    forceVisible: true
  });
</script>

<div class="flex flex-col gap-1">
  <span use:melt={$label} class="text-sm font-semibold text-ocean-800">
    Select Date
  </span>

  <div class="flex items-center gap-2">
    <!-- Date Field -->
    <div
      use:melt={$field}
      class="flex items-center gap-1 px-3 py-2 rounded-lg border-2 border-ocean-500/40 bg-white"
    >
      {#each $segments as seg}
        {#if seg.part === 'literal'}
          <span class="text-ocean-400">{seg.value}</span>
        {:else}
          <span
            use:melt={$segment(seg.part)}
            class="rounded px-0.5 tabular-nums outline-none focus:bg-ocean-100"
          >
            {seg.value}
          </span>
        {/if}
      {/each}
    </div>

    <!-- Calendar Trigger -->
    <button
      use:melt={$trigger}
      class="p-2 rounded-lg border-2 border-ocean-500/40 bg-white hover:bg-ocean-50 transition-colors"
    >
      <CalendarDays class="h-5 w-5 text-ocean-500" />
    </button>
  </div>

  <!-- Calendar Popup -->
  {#if $open}
    <div
      use:melt={$content}
      class="z-50 bg-white rounded-xl shadow-xl border border-ocean-200 p-4"
    >
      <div use:melt={$calendar}>
        <!-- Header -->
        <div class="flex items-center justify-between mb-4">
          <button use:melt={$prevButton} class="p-1 rounded hover:bg-ocean-100">
            <ChevronLeft class="h-5 w-5 text-ocean-600" />
          </button>
          <div use:melt={$heading} class="font-semibold text-ocean-800">
            <!-- Month Year -->
          </div>
          <button use:melt={$nextButton} class="p-1 rounded hover:bg-ocean-100">
            <ChevronRight class="h-5 w-5 text-ocean-600" />
          </button>
        </div>

        <!-- Calendar Grid -->
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
                        class="h-9 w-9 flex items-center justify-center rounded-lg text-sm cursor-pointer
                          data-[outside-month]:text-ocean-300
                          data-[selected]:bg-ocean-500 data-[selected]:text-white
                          data-[today]:ring-2 data-[today]:ring-ocean-400
                          hover:bg-ocean-100 transition-colors"
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
    </div>
  {/if}
</div>
```

## Styling with Tailwind

```css
/* Calendar cell states */
[data-melt-datepicker-cell][data-selected] {
  @apply bg-ocean-500 text-white;
}

[data-melt-datepicker-cell][data-today] {
  @apply ring-2 ring-ocean-400;
}

[data-melt-datepicker-cell][data-outside-month] {
  @apply text-ocean-300;
}

[data-melt-datepicker-cell][data-disabled] {
  @apply opacity-40 cursor-not-allowed;
}

/* Hover state */
[data-melt-datepicker-cell]:not([data-disabled]):hover {
  @apply bg-ocean-100;
}
```

## VacayTracker Integration

```svelte
<script lang="ts">
  import { createDatePicker, melt } from '@melt-ui/svelte';
  import { today, getLocalTimeZone } from '@internationalized/date';

  const {
    elements: { trigger, content, cell },
    states: { value }
  } = createDatePicker({
    // Start from today
    minValue: today(getLocalTimeZone()),
    // Disable weekends (based on settings)
    isDateDisabled: (date) => {
      if (!excludeWeekends) return false;
      const day = date.toDate('UTC').getDay();
      return day === 0 || day === 6;
    },
    // Close after selection
    closeOnSelect: true
  });

  // Convert to API format when submitting
  function getApiDate() {
    if (!$value) return '';
    return $value.toString(); // YYYY-MM-DD format
  }
</script>
```

## With Presets

```svelte
<script lang="ts">
  import { today, getLocalTimeZone } from '@internationalized/date';

  const presets = [
    { label: 'Today', value: today(getLocalTimeZone()) },
    { label: 'Tomorrow', value: today(getLocalTimeZone()).add({ days: 1 }) },
    { label: 'Next Week', value: today(getLocalTimeZone()).add({ weeks: 1 }) },
    { label: 'Next Month', value: today(getLocalTimeZone()).add({ months: 1 }) }
  ];
</script>

<div use:melt={$content} class="flex">
  <div class="border-r border-ocean-200 p-2">
    {#each presets as preset}
      <button
        onclick={() => value.set(preset.value)}
        class="w-full text-left px-3 py-2 text-sm rounded hover:bg-ocean-100"
      >
        {preset.label}
      </button>
    {/each}
  </div>
  <div class="p-4">
    <!-- Calendar -->
  </div>
</div>
```
