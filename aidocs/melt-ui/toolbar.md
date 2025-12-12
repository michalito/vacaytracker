# Toolbar

A container for grouping action buttons and controls.

## Use Cases

- Text editor toolbars
- Action button groups
- Control panels
- Form action bars

## Installation

```typescript
import { createToolbar } from '@melt-ui/svelte';
```

## API Reference

### createToolbar Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `loop` | `boolean` | `true` | Loop keyboard navigation |
| `orientation` | `'horizontal' \| 'vertical'` | `'horizontal'` | Layout direction |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Toolbar container |
| `button` | Toolbar button |
| `link` | Toolbar link |
| `separator` | Visual separator |

### Returned Builders

| Builder | Description |
|---------|-------------|
| `createToolbarGroup` | Create a button toggle group |

## Data Attributes

### Root
- `[data-orientation]` - `'horizontal'` or `'vertical'`
- `[data-melt-toolbar]` - Present on root

### Button
- `[data-melt-toolbar-button]` - Present on buttons

### Link
- `[data-melt-toolbar-link]` - Present on links

### Separator
- `[data-orientation]` - Separator orientation
- `[data-melt-toolbar-separator]` - Present on separator

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Tab` | Move focus in/out of toolbar |
| `Arrow Left/Up` | Previous item |
| `Arrow Right/Down` | Next item |
| `Home` | First item |
| `End` | Last item |

## Accessibility

Follows [WAI-ARIA Toolbar Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/toolbar/):
- `role="toolbar"` on container
- Proper keyboard navigation
- Focus management

## Example

```svelte
<script lang="ts">
  import { createToolbar, melt } from '@melt-ui/svelte';
  import { Bold, Italic, Underline, Link, AlignLeft, AlignCenter, AlignRight } from 'lucide-svelte';

  const {
    elements: { root, button, separator },
    builders: { createToolbarGroup }
  } = createToolbar();

  const alignGroup = createToolbarGroup({
    type: 'single',
    defaultValue: 'left'
  });
</script>

<div
  use:melt={$root}
  class="flex items-center gap-1 p-1 bg-white border border-ocean-200 rounded-lg"
>
  <!-- Text formatting buttons -->
  <button
    use:melt={$button}
    class="p-2 rounded text-ocean-600 hover:bg-ocean-100 transition-colors"
  >
    <Bold class="h-4 w-4" />
  </button>
  <button
    use:melt={$button}
    class="p-2 rounded text-ocean-600 hover:bg-ocean-100 transition-colors"
  >
    <Italic class="h-4 w-4" />
  </button>
  <button
    use:melt={$button}
    class="p-2 rounded text-ocean-600 hover:bg-ocean-100 transition-colors"
  >
    <Underline class="h-4 w-4" />
  </button>

  <div use:melt={$separator} class="w-px h-6 bg-ocean-200 mx-1"></div>

  <!-- Alignment toggle group -->
  <div use:melt={alignGroup.elements.root} class="flex">
    <button
      use:melt={alignGroup.elements.item('left')}
      class="p-2 rounded transition-colors
        data-[state=off]:text-ocean-600 data-[state=off]:hover:bg-ocean-100
        data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
    >
      <AlignLeft class="h-4 w-4" />
    </button>
    <button
      use:melt={alignGroup.elements.item('center')}
      class="p-2 rounded transition-colors
        data-[state=off]:text-ocean-600 data-[state=off]:hover:bg-ocean-100
        data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
    >
      <AlignCenter class="h-4 w-4" />
    </button>
    <button
      use:melt={alignGroup.elements.item('right')}
      class="p-2 rounded transition-colors
        data-[state=off]:text-ocean-600 data-[state=off]:hover:bg-ocean-100
        data-[state=on]:bg-ocean-500 data-[state=on]:text-white"
    >
      <AlignRight class="h-4 w-4" />
    </button>
  </div>

  <div use:melt={$separator} class="w-px h-6 bg-ocean-200 mx-1"></div>

  <!-- Link button -->
  <button
    use:melt={$button}
    class="p-2 rounded text-ocean-600 hover:bg-ocean-100 transition-colors"
  >
    <Link class="h-4 w-4" />
  </button>
</div>
```

## Styling with Tailwind

```css
/* Toolbar container */
[data-melt-toolbar] {
  @apply flex items-center gap-1 p-1 rounded-lg;
}

/* Toolbar button */
[data-melt-toolbar-button] {
  @apply p-2 rounded text-ocean-600 hover:bg-ocean-100 transition-colors;
}

[data-melt-toolbar-button]:focus-visible {
  @apply ring-2 ring-ocean-500 ring-offset-2;
}

/* Separator */
[data-melt-toolbar-separator][data-orientation="vertical"] {
  @apply w-px h-6 bg-ocean-200 mx-1;
}

[data-melt-toolbar-separator][data-orientation="horizontal"] {
  @apply h-px w-full bg-ocean-200 my-1;
}
```

## With Links

```svelte
<script lang="ts">
  import { createToolbar, melt } from '@melt-ui/svelte';
  import { ExternalLink } from 'lucide-svelte';

  const { elements: { root, button, link, separator } } = createToolbar();
</script>

<div use:melt={$root} class="flex items-center gap-1 p-2 bg-white rounded-lg">
  <button use:melt={$button} class="px-3 py-1.5 rounded hover:bg-ocean-100">
    Save
  </button>

  <div use:melt={$separator} class="w-px h-6 bg-ocean-200"></div>

  <a
    use:melt={$link}
    href="https://docs.example.com"
    target="_blank"
    class="flex items-center gap-1 px-3 py-1.5 rounded text-ocean-600 hover:bg-ocean-100"
  >
    <span>Docs</span>
    <ExternalLink class="h-3 w-3" />
  </a>
</div>
```

## Vertical Toolbar

```svelte
<script lang="ts">
  const { elements: { root, button, separator } } = createToolbar({
    orientation: 'vertical'
  });
</script>

<div
  use:melt={$root}
  class="flex flex-col items-center gap-1 p-1 bg-white border border-ocean-200 rounded-lg"
>
  <button use:melt={$button} class="p-2 rounded hover:bg-ocean-100">
    <Home class="h-4 w-4" />
  </button>
  <button use:melt={$button} class="p-2 rounded hover:bg-ocean-100">
    <Search class="h-4 w-4" />
  </button>

  <div use:melt={$separator} class="h-px w-6 bg-ocean-200"></div>

  <button use:melt={$button} class="p-2 rounded hover:bg-ocean-100">
    <Settings class="h-4 w-4" />
  </button>
</div>
```

## Action Bar

```svelte
<script lang="ts">
  import { createToolbar, melt } from '@melt-ui/svelte';
  import { Save, Undo, Redo, Trash } from 'lucide-svelte';

  const { elements: { root, button, separator } } = createToolbar();
</script>

<div class="fixed bottom-4 left-1/2 -translate-x-1/2 z-50">
  <div
    use:melt={$root}
    class="flex items-center gap-2 px-4 py-2 bg-white/95 backdrop-blur shadow-lg rounded-full"
  >
    <button use:melt={$button} class="p-2 rounded-full hover:bg-ocean-100" title="Save">
      <Save class="h-5 w-5 text-ocean-600" />
    </button>

    <div use:melt={$separator} class="w-px h-6 bg-ocean-200"></div>

    <button use:melt={$button} class="p-2 rounded-full hover:bg-ocean-100" title="Undo">
      <Undo class="h-5 w-5 text-ocean-600" />
    </button>
    <button use:melt={$button} class="p-2 rounded-full hover:bg-ocean-100" title="Redo">
      <Redo class="h-5 w-5 text-ocean-600" />
    </button>

    <div use:melt={$separator} class="w-px h-6 bg-ocean-200"></div>

    <button use:melt={$button} class="p-2 rounded-full hover:bg-red-100" title="Delete">
      <Trash class="h-5 w-5 text-red-600" />
    </button>
  </div>
</div>
```

## With Toggle Groups

```svelte
<script lang="ts">
  const { elements, builders: { createToolbarGroup } } = createToolbar();

  // Single selection group (alignment)
  const alignGroup = createToolbarGroup({ type: 'single', defaultValue: 'left' });

  // Multiple selection group (formatting)
  const formatGroup = createToolbarGroup({ type: 'multiple', defaultValue: [] });
</script>

<div use:melt={elements.root}>
  <!-- Format toggles (multiple) -->
  <div use:melt={formatGroup.elements.root} class="flex">
    <button use:melt={formatGroup.elements.item('bold')}>B</button>
    <button use:melt={formatGroup.elements.item('italic')}>I</button>
    <button use:melt={formatGroup.elements.item('underline')}>U</button>
  </div>

  <div use:melt={elements.separator}></div>

  <!-- Alignment toggles (single) -->
  <div use:melt={alignGroup.elements.root} class="flex">
    <button use:melt={alignGroup.elements.item('left')}>L</button>
    <button use:melt={alignGroup.elements.item('center')}>C</button>
    <button use:melt={alignGroup.elements.item('right')}>R</button>
  </div>
</div>
```
