# Date Range Picker

A calendar popup for selecting date ranges with visual range highlighting.

## Use Cases

- Vacation request forms
- Report date filtering
- Booking date selection
- Analytics date ranges

## Installation

```typescript
import { createDateRangePicker } from '@melt-ui/svelte';
```

## API Reference

### createDateRangePicker Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `defaultValue` | `DateRange` | — | Initial date range |
| `value` | `Writable<DateRange>` | — | Controlled date range |
| `defaultPlaceholder` | `DateValue` | Today | Initial displayed month |
| `minValue` | `DateValue` | — | Minimum selectable date |
| `maxValue` | `DateValue` | — | Maximum selectable date |
| `disabled` | `boolean` | `false` | Disable picker |
| `readonly` | `boolean` | `false` | Make read-only |
| `locale` | `string` | `'en'` | Locale for formatting |
| `isDateDisabled` | `(date) => boolean` | — | Disable specific dates |
| `isDateUnavailable` | `(date) => boolean` | — | Mark dates unavailable |
| `weekStartsOn` | `0-6` | `0` | First day of week |
| `fixedWeeks` | `boolean` | `false` | Always show 6 weeks |
| `numberOfMonths` | `number` | `1` | Months to display |
| `portal` | `string \| HTMLElement \| null` | `'body'` | Portal target |
| `forceVisible` | `boolean` | `false` | Force visibility |
| `closeOnOutsideClick` | `boolean` | `true` | Close on click outside |
| `onValueChange` | `ChangeFn` | — | Value change callback |
| `onOpenChange` | `ChangeFn` | — | Open state callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `trigger` | Button to open picker |
| `content` | Calendar popup |
| `calendar` | Calendar container |
| `cell` | Individual day cell |
| `heading` | Month/year heading |
| `grid` | Calendar grid |
| `prevButton` | Previous month button |
| `nextButton` | Next month button |
| `field` | Date range field container |
| `startSegment` | Start date segment builder (function) |
| `endSegment` | End date segment builder (function) |
| `label` | Field label |
| `startHiddenInput` | Hidden input for start date (form submission) |
| `endHiddenInput` | Hidden input for end date (form submission) |
| `validation` | Validation message container |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `open` | `Writable<boolean>` | Popup visibility |
| `value` | `Writable<DateRange>` | Selected range `{ start, end }` |
| `months` | `Readable<Month[]>` | Month data for rendering |
| `weekdays` | `Readable<string[]>` | Weekday labels |
| `isInvalid` | `Readable<boolean>` | Validation state |
| `segmentContents` | `Readable<{ start: Segment[], end: Segment[] }>` | Segment arrays for start/end |
| `startValue` | `Writable<DateValue \| undefined>` | Start date value |
| `endValue` | `Writable<DateValue \| undefined>` | End date value |
| `headingValue` | `Readable<string>` | Formatted heading text |
| `placeholder` | `Writable<DateValue>` | Current placeholder date |

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
| `Space` / `Enter` | Select start/end date |
| `Arrow Keys` | Navigate calendar |
| `Page Up/Down` | Previous/next month |
| `Escape` | Close picker |

## Accessibility

- Full keyboard navigation for range selection
- Screen reader announces range selection state
- Proper ARIA labels for start/end selection

## Example

```svelte
<script lang="ts">
  import { createDateRangePicker, melt } from '@melt-ui/svelte';
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
      startSegment,
      endSegment,
      label
    },
    states: { open, value, months, weekdays, segmentContents }
  } = createDateRangePicker({
    locale: 'en-GB',
    weekStartsOn: 1,
    numberOfMonths: 2,
    forceVisible: true
  });
</script>

<div class="flex flex-col gap-1">
  <span use:melt={$label} class="text-sm font-semibold text-ocean-800">
    Vacation Period
  </span>

  <div class="flex items-center gap-2">
    <!-- Date Range Field -->
    <div
      use:melt={$field}
      class="flex items-center gap-2 px-3 py-2 rounded-lg border-2 border-ocean-500/40 bg-white"
    >
      <!-- Start segments -->
      <div class="flex items-center gap-0.5">
        {#each $segmentContents.start as seg}
          {#if seg.part === 'literal'}
            <span class="text-ocean-400">{seg.value}</span>
          {:else}
            <span use:melt={$startSegment(seg.part)} class="rounded px-0.5 tabular-nums outline-none focus:bg-ocean-100">
              {seg.value}
            </span>
          {/if}
        {/each}
      </div>

      <span class="text-ocean-400">to</span>

      <!-- End segments -->
      <div class="flex items-center gap-0.5">
        {#each $segmentContents.end as seg}
          {#if seg.part === 'literal'}
            <span class="text-ocean-400">{seg.value}</span>
          {:else}
            <span use:melt={$endSegment(seg.part)} class="rounded px-0.5 tabular-nums outline-none focus:bg-ocean-100">
              {seg.value}
            </span>
          {/if}
        {/each}
      </div>
    </div>

    <!-- Trigger Button -->
    <button
      use:melt={$trigger}
      class="p-2 rounded-lg border-2 border-ocean-500/40 hover:bg-ocean-50"
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
      <div use:melt={$calendar} class="flex gap-4">
        {#each $months as month, i}
          <div>
            <!-- Month Header -->
            <div class="flex items-center justify-between mb-4">
              {#if i === 0}
                <button use:melt={$prevButton} class="p-1 rounded hover:bg-ocean-100">
                  <ChevronLeft class="h-5 w-5" />
                </button>
              {:else}
                <div class="w-7"></div>
              {/if}
              <div use:melt={$heading} class="font-semibold text-ocean-800">
                {month.value.toDate('UTC').toLocaleString('en', { month: 'long', year: 'numeric' })}
              </div>
              {#if i === $months.length - 1}
                <button use:melt={$nextButton} class="p-1 rounded hover:bg-ocean-100">
                  <ChevronRight class="h-5 w-5" />
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
                    <th class="text-sm font-medium text-ocean-500 pb-2 w-9">{day}</th>
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
                          class="h-9 w-9 flex items-center justify-center text-sm cursor-pointer
                            data-[outside-month]:text-ocean-300
                            data-[selection-start]:bg-ocean-500 data-[selection-start]:text-white data-[selection-start]:rounded-l-lg
                            data-[selection-end]:bg-ocean-500 data-[selection-end]:text-white data-[selection-end]:rounded-r-lg
                            data-[highlighted]:bg-ocean-100
                            data-[today]:ring-2 data-[today]:ring-ocean-400
                            hover:bg-ocean-200 transition-colors"
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
    </div>
  {/if}
</div>
```

## Styling with Tailwind

```css
/* Selection start */
[data-melt-calendar-cell][data-selection-start] {
  @apply bg-ocean-500 text-white rounded-l-lg;
}

/* Selection end */
[data-melt-calendar-cell][data-selection-end] {
  @apply bg-ocean-500 text-white rounded-r-lg;
}

/* Highlighted (in range preview) */
[data-melt-calendar-cell][data-highlighted] {
  @apply bg-ocean-100;
}

/* Today indicator */
[data-melt-calendar-cell][data-today] {
  @apply ring-2 ring-ocean-400;
}

/* Disabled dates */
[data-melt-calendar-cell][data-disabled] {
  @apply opacity-40 cursor-not-allowed;
}

/* Unavailable dates */
[data-melt-calendar-cell][data-unavailable] {
  @apply text-ocean-300 line-through;
}
```

## VacayTracker Integration

```svelte
<script lang="ts">
  import { createDateRangePicker, melt } from '@melt-ui/svelte';
  import { today, getLocalTimeZone } from '@internationalized/date';
  import { auth } from '$lib/stores/auth.svelte';

  const {
    elements: { trigger, content, cell },
    states: { value }
  } = createDateRangePicker({
    locale: 'en-GB',
    weekStartsOn: 1,
    numberOfMonths: 2,
    minValue: today(getLocalTimeZone()),
    isDateDisabled: (date) => {
      // Disable weekends based on company setting
      const day = date.toDate('UTC').getDay();
      return day === 0 || day === 6;
    }
  });

  // Calculate business days for the selected range
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

  // Check against user's balance
  const exceedsBalance = $derived(
    businessDays() > (auth.user?.vacationBalance ?? 0)
  );
</script>
```

## With Preset Ranges

```svelte
<script lang="ts">
  import { today, getLocalTimeZone } from '@internationalized/date';

  const presets = [
    {
      label: 'This Week',
      range: {
        start: today(getLocalTimeZone()).subtract({ days: today(getLocalTimeZone()).dayOfWeek - 1 }),
        end: today(getLocalTimeZone()).add({ days: 7 - today(getLocalTimeZone()).dayOfWeek })
      }
    },
    {
      label: 'Next Week',
      range: {
        start: today(getLocalTimeZone()).add({ days: 8 - today(getLocalTimeZone()).dayOfWeek }),
        end: today(getLocalTimeZone()).add({ days: 14 - today(getLocalTimeZone()).dayOfWeek })
      }
    }
  ];

  function selectPreset(range: DateRange) {
    value.set(range);
  }
</script>

<div class="flex">
  <div class="border-r p-2 space-y-1">
    {#each presets as preset}
      <button
        onclick={() => selectPreset(preset.range)}
        class="w-full text-left px-3 py-2 text-sm rounded hover:bg-ocean-100"
      >
        {preset.label}
      </button>
    {/each}
  </div>
  <div class="p-4">
    <!-- Calendar grids -->
  </div>
</div>
```
