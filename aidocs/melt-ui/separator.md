# Separator

A visual divider between content sections.

## Use Cases

- Menu dividers
- Section separators
- List item dividers
- Form section breaks

## Installation

```typescript
import { createSeparator } from '@melt-ui/svelte';
```

## API Reference

### createSeparator Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `orientation` | `'horizontal' \| 'vertical'` | `'horizontal'` | Separator direction |
| `decorative` | `boolean` | `false` | If true, hidden from assistive tech |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Separator element |

## Data Attributes

### Root
- `[data-orientation]` - `'horizontal'` or `'vertical'`
- `[data-melt-separator]` - Present on separator

## Accessibility

- Uses `role="separator"` for semantic meaning
- `aria-orientation` indicates direction
- Set `decorative: true` for purely visual separators

## Example

```svelte
<script lang="ts">
  import { createSeparator, melt } from '@melt-ui/svelte';

  const { elements: { root } } = createSeparator();
</script>

<div class="space-y-4">
  <p class="text-ocean-700">Section one content</p>

  <div use:melt={$root} class="h-px bg-ocean-200"></div>

  <p class="text-ocean-700">Section two content</p>
</div>
```

## Styling with Tailwind

```css
/* Horizontal separator */
[data-melt-separator][data-orientation="horizontal"] {
  @apply h-px w-full bg-ocean-200;
}

/* Vertical separator */
[data-melt-separator][data-orientation="vertical"] {
  @apply w-px h-full bg-ocean-200;
}
```

## Horizontal Separator

```svelte
<script lang="ts">
  import { createSeparator, melt } from '@melt-ui/svelte';

  const { elements: { root } } = createSeparator({
    orientation: 'horizontal'
  });
</script>

<div use:melt={$root} class="h-px w-full bg-ocean-200 my-4"></div>
```

## Vertical Separator

```svelte
<script lang="ts">
  import { createSeparator, melt } from '@melt-ui/svelte';

  const { elements: { root } } = createSeparator({
    orientation: 'vertical'
  });
</script>

<div class="flex items-center h-8">
  <span class="text-ocean-700">Option A</span>
  <div use:melt={$root} class="w-px h-full bg-ocean-200 mx-4"></div>
  <span class="text-ocean-700">Option B</span>
</div>
```

## Menu Separator

```svelte
<script lang="ts">
  import { createSeparator, melt } from '@melt-ui/svelte';

  const { elements: { root } } = createSeparator({
    decorative: true // Purely visual, hide from screen readers
  });
</script>

<div class="py-2 bg-white rounded-lg shadow-lg">
  <button class="w-full px-4 py-2 text-left hover:bg-ocean-50">Edit</button>
  <button class="w-full px-4 py-2 text-left hover:bg-ocean-50">Duplicate</button>

  <div use:melt={$root} class="h-px bg-ocean-100 my-1"></div>

  <button class="w-full px-4 py-2 text-left text-red-600 hover:bg-red-50">Delete</button>
</div>
```

## With Text Label

```svelte
<script lang="ts">
  import { createSeparator, melt } from '@melt-ui/svelte';

  const { elements: { root } } = createSeparator();
</script>

<div class="flex items-center gap-4">
  <div use:melt={$root} class="flex-1 h-px bg-ocean-200"></div>
  <span class="text-sm text-ocean-500">OR</span>
  <div class="flex-1 h-px bg-ocean-200"></div>
</div>
```

## Form Section Separator

```svelte
<form class="space-y-6">
  <!-- Personal Info Section -->
  <div class="space-y-4">
    <h3 class="font-semibold text-ocean-800">Personal Information</h3>
    <!-- Form fields -->
  </div>

  <div use:melt={$root} class="h-px bg-ocean-200"></div>

  <!-- Contact Section -->
  <div class="space-y-4">
    <h3 class="font-semibold text-ocean-800">Contact Details</h3>
    <!-- Form fields -->
  </div>

  <div use:melt={$root} class="h-px bg-ocean-200"></div>

  <!-- Preferences Section -->
  <div class="space-y-4">
    <h3 class="font-semibold text-ocean-800">Preferences</h3>
    <!-- Form fields -->
  </div>
</form>
```

## Toolbar Separator

```svelte
<div class="flex items-center gap-2 p-2 bg-white rounded-lg shadow">
  <button class="p-2 rounded hover:bg-ocean-100">
    <Bold class="h-4 w-4" />
  </button>
  <button class="p-2 rounded hover:bg-ocean-100">
    <Italic class="h-4 w-4" />
  </button>

  <div use:melt={verticalSeparator.elements.root} class="w-px h-6 bg-ocean-200 mx-1"></div>

  <button class="p-2 rounded hover:bg-ocean-100">
    <AlignLeft class="h-4 w-4" />
  </button>
  <button class="p-2 rounded hover:bg-ocean-100">
    <AlignCenter class="h-4 w-4" />
  </button>
</div>
```

## Decorative vs Semantic

```typescript
// Semantic separator - announced by screen readers
const semanticSeparator = createSeparator({
  decorative: false // Default
});

// Decorative separator - hidden from screen readers
const decorativeSeparator = createSeparator({
  decorative: true
});
```

Use `decorative: true` when the separator is purely visual and doesn't convey meaningful content grouping.
