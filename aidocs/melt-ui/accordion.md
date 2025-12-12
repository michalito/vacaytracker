# Accordion

An interactive component enabling content organization through expandable/collapsible sections.

## Use Cases

- FAQ sections
- Settings panels with grouped options
- Content organization with multiple sections
- Expandable navigation menus

## Installation

```typescript
import { createAccordion } from '@melt-ui/svelte';
```

## API Reference

### createAccordion Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `multiple` | `boolean` | `false` | Allow multiple items open simultaneously |
| `disabled` | `boolean` | `false` | Disable the entire accordion |
| `forceVisible` | `boolean` | `false` | Force visibility for custom transitions |
| `defaultValue` | `string \| string[]` | — | Initial open item(s) |
| `value` | `Writable<string \| string[]>` | — | Controlled open state |
| `onValueChange` | `ChangeFn` | — | Callback when value changes |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Container wrapping all accordion items |
| `item` | Individual accordion section wrapper |
| `trigger` | Button toggling item expansion |
| `content` | Collapsible content area |
| `heading` | Optional semantic heading wrapper |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `value` | `Writable<string \| string[]>` | Currently open item(s) |

### Returned Helpers

| Helper | Description |
|--------|-------------|
| `isSelected(value)` | Check if an item is currently open |

## Data Attributes

### Root
- `[data-orientation]` - `'vertical'` or `'horizontal'`
- `[data-melt-accordion]` - Present on root element

### Item
- `[data-state]` - `'open'` or `'closed'`
- `[data-disabled]` - Present when disabled
- `[data-melt-accordion-item]` - Present on item elements

### Trigger
- `[data-state]` - `'open'` or `'closed'`
- `[data-disabled]` - Present when disabled
- `[data-melt-accordion-trigger]` - Present on triggers

### Content
- `[data-state]` - `'open'` or `'closed'`
- `[data-disabled]` - Present when disabled
- `[data-melt-accordion-content]` - Present on content areas

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Space` / `Enter` | Toggle focused section |
| `Arrow Down` | Move focus to next trigger |
| `Arrow Up` | Move focus to previous trigger |
| `Home` | Move focus to first trigger |
| `End` | Move focus to last trigger |

## Accessibility

Follows the [WAI-ARIA Accordion Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/accordion/):
- Proper `aria-expanded` state on triggers
- `aria-controls` linking trigger to content
- Semantic heading structure support

## Example

```svelte
<script lang="ts">
  import { createAccordion, melt } from '@melt-ui/svelte';
  import { ChevronDown } from 'lucide-svelte';

  const {
    elements: { root, item, trigger, content },
    helpers: { isSelected }
  } = createAccordion({
    defaultValue: 'item-1'
  });

  const items = [
    { id: 'item-1', title: 'Personal Information', content: 'Update your personal details.' },
    { id: 'item-2', title: 'Vacation Balance', content: 'View your remaining vacation days.' },
    { id: 'item-3', title: 'Notification Settings', content: 'Configure email preferences.' }
  ];
</script>

<div use:melt={$root} class="w-full max-w-md mx-auto">
  {#each items as { id, title, content: text }}
    <div use:melt={$item(id)} class="border-b border-ocean-200">
      <h3>
        <button
          use:melt={$trigger(id)}
          class="flex w-full items-center justify-between py-4 px-2 text-ocean-800 font-medium hover:bg-ocean-50 transition-colors"
        >
          {title}
          <ChevronDown
            class="h-4 w-4 transition-transform duration-200 {$isSelected(id) ? 'rotate-180' : ''}"
          />
        </button>
      </h3>
      {#if $isSelected(id)}
        <div use:melt={$content(id)} class="px-2 pb-4 text-ocean-600">
          {text}
        </div>
      {/if}
    </div>
  {/each}
</div>
```

## Styling with Tailwind

```css
/* Trigger hover state */
[data-melt-accordion-trigger]:hover {
  @apply bg-ocean-50;
}

/* Open state indicator */
[data-melt-accordion-trigger][data-state="open"] {
  @apply bg-ocean-100 text-ocean-900;
}

/* Content animation */
[data-melt-accordion-content] {
  @apply overflow-hidden;
}

[data-melt-accordion-content][data-state="open"] {
  animation: slideDown 200ms ease-out;
}

[data-melt-accordion-content][data-state="closed"] {
  animation: slideUp 200ms ease-out;
}

@keyframes slideDown {
  from { height: 0; }
  to { height: var(--melt-content-height); }
}

@keyframes slideUp {
  from { height: var(--melt-content-height); }
  to { height: 0; }
}
```

## Multiple Items Open

```typescript
const accordion = createAccordion({
  multiple: true,
  defaultValue: ['item-1', 'item-2']
});
```

## Disabled Items

```svelte
<div use:melt={$item('disabled-item')} data-disabled>
  <button use:melt={$trigger('disabled-item')} disabled>
    Unavailable Section
  </button>
</div>
```
