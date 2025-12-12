# Pagination

A navigation component for paginated content.

## Use Cases

- Data table pagination
- Search results navigation
- Gallery pagination
- List navigation

## Installation

```typescript
import { createPagination } from '@melt-ui/svelte';
```

## API Reference

### createPagination Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `count` | `number` | — | Total number of items |
| `perPage` | `number` | `1` | Items per page |
| `siblingCount` | `number` | `1` | Pages shown beside current |
| `defaultPage` | `number` | `1` | Initial page |
| `page` | `Writable<number>` | — | Controlled current page |
| `onPageChange` | `ChangeFn` | — | Page change callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Pagination container |
| `pageTrigger` | Page number button |
| `prevButton` | Previous page button |
| `nextButton` | Next page button |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `pages` | `Readable<PageItem[]>` | Page items including ellipsis |
| `range` | `Readable<{ start: number, end: number }>` | Current item range |
| `totalPages` | `Readable<number>` | Total page count |
| `page` | `Writable<number>` | Current page |

### PageItem Type

```typescript
type PageItem = {
  type: 'page' | 'ellipsis';
  value: number;
};
```

## Data Attributes

### Root
- `[data-melt-pagination]` - Present on root

### Page Trigger
- `[data-selected]` - Currently selected page
- `[data-melt-pagination-page]` - Present on page triggers

### Prev/Next Button
- `[data-disabled]` - When at first/last page
- `[data-melt-pagination-prev]` - Present on prev button
- `[data-melt-pagination-next]` - Present on next button

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Tab` | Navigate between elements |
| `Space` / `Enter` | Select page |
| `Arrow Left` | Previous focusable trigger |
| `Arrow Right` | Next focusable trigger |
| `Home` | First page |
| `End` | Last page |

## Accessibility

Follows accessibility guidelines:
- Proper `aria-label` on navigation
- Current page announced to screen readers
- Keyboard navigable

## Example

```svelte
<script lang="ts">
  import { createPagination, melt } from '@melt-ui/svelte';
  import { ChevronLeft, ChevronRight } from 'lucide-svelte';

  const {
    elements: { root, pageTrigger, prevButton, nextButton },
    states: { pages, range, page }
  } = createPagination({
    count: 100,
    perPage: 10,
    defaultPage: 1,
    siblingCount: 1
  });
</script>

<nav use:melt={$root} class="flex items-center justify-center gap-1">
  <button
    use:melt={$prevButton}
    class="p-2 rounded-lg text-ocean-600 hover:bg-ocean-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
  >
    <ChevronLeft class="h-5 w-5" />
  </button>

  {#each $pages as p}
    {#if p.type === 'ellipsis'}
      <span class="px-2 text-ocean-400">...</span>
    {:else}
      <button
        use:melt={$pageTrigger(p)}
        class="h-10 w-10 rounded-lg text-sm font-medium transition-colors
          data-[selected]:bg-ocean-500 data-[selected]:text-white
          text-ocean-700 hover:bg-ocean-100"
      >
        {p.value}
      </button>
    {/if}
  {/each}

  <button
    use:melt={$nextButton}
    class="p-2 rounded-lg text-ocean-600 hover:bg-ocean-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
  >
    <ChevronRight class="h-5 w-5" />
  </button>
</nav>

<p class="text-center text-sm text-ocean-500 mt-2">
  Showing {$range.start} - {$range.end} of 100
</p>
```

## Styling with Tailwind

```css
/* Page trigger */
[data-melt-pagination-page] {
  @apply h-10 w-10 rounded-lg text-sm font-medium transition-colors;
  @apply text-ocean-700 hover:bg-ocean-100;
}

/* Selected page */
[data-melt-pagination-page][data-selected] {
  @apply bg-ocean-500 text-white;
}

/* Disabled buttons */
[data-melt-pagination-prev][data-disabled],
[data-melt-pagination-next][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}
```

## VacayTracker Integration

```svelte
<script lang="ts">
  import { createPagination, melt } from '@melt-ui/svelte';
  import { ChevronLeft, ChevronRight } from 'lucide-svelte';

  interface Props {
    page: number;
    totalPages: number;
    onPageChange: (page: number) => void;
  }

  let { page, totalPages, onPageChange }: Props = $props();

  const pagination = createPagination({
    count: totalPages,
    perPage: 1,
    page: page,
    siblingCount: 2
  });

  const { elements: { root, pageTrigger, prevButton, nextButton }, states: { pages } } = pagination;

  // Sync external page changes
  $effect(() => {
    pagination.states.page.set(page);
  });

  // Notify parent of page changes
  $effect(() => {
    const currentPage = $pagination.states.page;
    if (currentPage !== page) {
      onPageChange(currentPage);
    }
  });
</script>

{#if totalPages > 1}
  <nav use:melt={$root} class="flex items-center justify-center gap-1">
    <button
      use:melt={$prevButton}
      class="p-2 rounded-lg text-ocean-600 hover:bg-sand-100 disabled:opacity-50 transition-colors"
    >
      <ChevronLeft class="h-4 w-4" />
    </button>

    {#each $pages as p}
      {#if p.type === 'ellipsis'}
        <span class="px-2 text-ocean-400">...</span>
      {:else}
        <button
          use:melt={$pageTrigger(p)}
          class="h-8 w-8 rounded-md text-sm font-medium transition-colors
            data-[selected]:bg-ocean-500 data-[selected]:text-white
            text-ocean-700 hover:bg-sand-100"
        >
          {p.value}
        </button>
      {/if}
    {/each}

    <button
      use:melt={$nextButton}
      class="p-2 rounded-lg text-ocean-600 hover:bg-sand-100 disabled:opacity-50 transition-colors"
    >
      <ChevronRight class="h-4 w-4" />
    </button>
  </nav>
{/if}
```

## With Item Count

```svelte
<script lang="ts">
  const { states: { range, page } } = createPagination({
    count: 156,
    perPage: 10
  });
</script>

<div class="flex items-center justify-between">
  <span class="text-sm text-ocean-600">
    Showing {$range.start} - {$range.end} of 156 items
  </span>

  <nav>
    <!-- Pagination controls -->
  </nav>
</div>
```

## Controlled Pagination

```svelte
<script lang="ts">
  import { writable } from 'svelte/store';

  const currentPage = writable(1);

  const { states: { page } } = createPagination({
    count: 100,
    perPage: 10,
    page: currentPage,
    onPageChange: ({ next }) => {
      // Could fetch data here
      console.log('Page changed to:', next);
      return next;
    }
  });

  // Programmatic navigation
  function goToPage(n: number) {
    currentPage.set(n);
  }
</script>
```

## Custom Sibling Count

```svelte
<script lang="ts">
  // Show more pages on larger screens
  const siblingCount = window.innerWidth > 768 ? 2 : 1;

  const pagination = createPagination({
    count: 100,
    perPage: 10,
    siblingCount
  });
</script>
```
