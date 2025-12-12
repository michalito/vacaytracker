# Table of Contents

A navigation component that tracks and highlights the currently active section.

## Use Cases

- Documentation navigation
- Article outlines
- Settings page navigation
- Long form content

## Installation

```typescript
import { createTableOfContents } from '@melt-ui/svelte';
```

## API Reference

### createTableOfContents Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `selector` | `string` | `'h2, h3, h4'` | CSS selector for headings |
| `exclude` | `string[]` | `[]` | IDs to exclude |
| `scrollOffset` | `number` | `0` | Offset when scrolling |
| `scrollBehaviour` | `'auto' \| 'smooth'` | `'smooth'` | Scroll behavior |
| `activeType` | `'lowest' \| 'highest' \| 'all'` | `'lowest'` | Which heading is active |

### Returned Elements

| Element | Description |
|---------|-------------|
| `item` | TOC item element |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `headingsTree` | `Readable<HeadingTree[]>` | Nested heading structure |
| `activeHeadingIdxs` | `Readable<number[]>` | Active heading indices |

### HeadingTree Type

```typescript
type HeadingTree = {
  id: string;
  text: string;
  level: number;
  children: HeadingTree[];
  node: HTMLElement;
};
```

### Returned Helpers

| Helper | Description |
|--------|-------------|
| `isActive(id)` | Check if heading is active |

## Data Attributes

### Item
- `[data-active]` - Present when heading is active
- `[data-melt-toc-item]` - Present on items

## Example

```svelte
<script lang="ts">
  import { createTableOfContents, melt } from '@melt-ui/svelte';

  const {
    elements: { item },
    states: { headingsTree, activeHeadingIdxs },
    helpers: { isActive }
  } = createTableOfContents({
    selector: 'h2, h3',
    scrollOffset: 100
  });
</script>

<nav class="sticky top-4 w-64 p-4 rounded-xl bg-white shadow-lg">
  <h2 class="text-sm font-semibold text-ocean-800 mb-4">On this page</h2>

  <ul class="space-y-2">
    {#each $headingsTree as heading}
      <li>
        <a
          use:melt={$item(heading.id)}
          href="#{heading.id}"
          class="block text-sm transition-colors
            {$isActive(heading.id)
              ? 'text-ocean-600 font-medium'
              : 'text-ocean-500 hover:text-ocean-700'}"
        >
          {heading.text}
        </a>

        {#if heading.children.length > 0}
          <ul class="ml-4 mt-2 space-y-2">
            {#each heading.children as child}
              <li>
                <a
                  use:melt={$item(child.id)}
                  href="#{child.id}"
                  class="block text-sm transition-colors
                    {$isActive(child.id)
                      ? 'text-ocean-600 font-medium'
                      : 'text-ocean-400 hover:text-ocean-600'}"
                >
                  {child.text}
                </a>
              </li>
            {/each}
          </ul>
        {/if}
      </li>
    {/each}
  </ul>
</nav>
```

## Styling with Tailwind

```css
/* Active item */
[data-melt-toc-item][data-active] {
  @apply text-ocean-600 font-medium;
}

/* Inactive item */
[data-melt-toc-item]:not([data-active]) {
  @apply text-ocean-500 hover:text-ocean-700;
}

/* Level indicators */
[data-melt-toc-item][data-level="2"] {
  @apply pl-0;
}

[data-melt-toc-item][data-level="3"] {
  @apply pl-4;
}

[data-melt-toc-item][data-level="4"] {
  @apply pl-8;
}
```

## Documentation Page

```svelte
<script lang="ts">
  import { createTableOfContents, melt } from '@melt-ui/svelte';

  const {
    elements: { item },
    states: { headingsTree },
    helpers: { isActive }
  } = createTableOfContents({
    selector: 'h2, h3, h4',
    scrollOffset: 80, // Account for sticky header
    scrollBehaviour: 'smooth'
  });
</script>

<div class="flex gap-8">
  <!-- Main content -->
  <main class="flex-1 prose">
    <h2 id="introduction">Introduction</h2>
    <p>...</p>

    <h2 id="getting-started">Getting Started</h2>
    <p>...</p>

    <h3 id="installation">Installation</h3>
    <p>...</p>

    <h3 id="configuration">Configuration</h3>
    <p>...</p>

    <h2 id="api-reference">API Reference</h2>
    <p>...</p>
  </main>

  <!-- Table of contents -->
  <aside class="w-64 shrink-0">
    <nav class="sticky top-20 p-4 rounded-xl border border-ocean-200">
      <h3 class="text-xs font-semibold text-ocean-500 uppercase tracking-wider mb-3">
        Contents
      </h3>

      {#each $headingsTree as heading}
        <a
          use:melt={$item(heading.id)}
          href="#{heading.id}"
          class="block py-1 text-sm border-l-2 pl-3 transition-colors
            {$isActive(heading.id)
              ? 'border-ocean-500 text-ocean-700 font-medium'
              : 'border-transparent text-ocean-500 hover:text-ocean-700 hover:border-ocean-300'}"
        >
          {heading.text}
        </a>

        {#each heading.children as child}
          <a
            use:melt={$item(child.id)}
            href="#{child.id}"
            class="block py-1 text-sm border-l-2 pl-6 transition-colors
              {$isActive(child.id)
                ? 'border-ocean-500 text-ocean-700 font-medium'
                : 'border-transparent text-ocean-400 hover:text-ocean-600'}"
          >
            {child.text}
          </a>
        {/each}
      {/each}
    </nav>
  </aside>
</div>
```

## Progress Indicator

```svelte
<script lang="ts">
  const {
    states: { headingsTree, activeHeadingIdxs }
  } = createTableOfContents();

  // Calculate progress based on active heading
  const progress = $derived(() => {
    const totalHeadings = $headingsTree.reduce(
      (acc, h) => acc + 1 + h.children.length,
      0
    );
    const activeIdx = $activeHeadingIdxs[0] ?? 0;
    return totalHeadings > 0 ? (activeIdx / totalHeadings) * 100 : 0;
  });
</script>

<div class="fixed top-0 left-0 h-1 bg-ocean-500 transition-all duration-300"
  style="width: {progress()}%"
></div>
```

## Custom Selector

```typescript
// Only track h2 headings
const toc = createTableOfContents({
  selector: 'h2'
});

// Track all headings with specific class
const toc = createTableOfContents({
  selector: '[data-toc-heading]'
});

// Exclude specific IDs
const toc = createTableOfContents({
  selector: 'h2, h3',
  exclude: ['hidden-section', 'private-section']
});
```

## Scroll Offset

```typescript
// Account for sticky header
const toc = createTableOfContents({
  scrollOffset: 80 // Header height + padding
});
```
