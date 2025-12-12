# Melt-UI Component Guide

## Overview

Melt-UI is a headless, accessible component builder library for Svelte. It provides unstyled, fully accessible UI primitives that integrate seamlessly with Tailwind CSS and Svelte 5 runes.

**Version**: 0.86.x
**Official Docs**: https://melt-ui.com

## Installation

Melt-UI is already installed in VacayTracker:

```bash
npm install @melt-ui/svelte
```

## Core Concepts

### Builder Pattern

Melt-UI uses a builder pattern where you call a `create*` function to initialize a component:

```svelte
<script lang="ts">
  import { createDialog } from '@melt-ui/svelte';

  const dialog = createDialog({
    // Configuration options
    defaultOpen: false,
    closeOnOutsideClick: true
  });

  // Destructure elements and states
  const { elements: { trigger, content, overlay }, states: { open } } = dialog;
</script>
```

### The `melt` Directive

Apply Melt-UI behavior to elements using the `melt` directive:

```svelte
<button use:melt={$trigger}>Open Dialog</button>
<div use:melt={$overlay} class="fixed inset-0 bg-black/50"></div>
<div use:melt={$content} class="bg-white rounded-lg p-6">
  <!-- Content -->
</div>
```

### Data Attributes for Styling

Melt-UI applies data attributes for state-based styling:

```css
/* Style based on state */
[data-state="open"] { opacity: 1; }
[data-state="closed"] { opacity: 0; }
[data-highlighted] { background-color: var(--color-ocean-100); }
[data-disabled] { opacity: 0.5; cursor: not-allowed; }
```

## Integration with Svelte 5 Runes

Melt-UI stores work naturally with Svelte 5:

```svelte
<script lang="ts">
  import { createSelect } from '@melt-ui/svelte';

  const select = createSelect({
    defaultSelected: { value: 'option1', label: 'Option 1' }
  });

  const { elements: { trigger, menu, option }, states: { selected } } = select;

  // React to selection changes
  $effect(() => {
    console.log('Selected:', $selected);
  });
</script>
```

## Common Patterns

### Controlled vs Uncontrolled

**Uncontrolled** (internal state):
```typescript
const dialog = createDialog({ defaultOpen: false });
```

**Controlled** (external state):
```typescript
import { writable } from 'svelte/store';
const openStore = writable(false);

const dialog = createDialog({ open: openStore });
```

### Portal Rendering

Most overlay components render to `<body>` by default:

```typescript
const dialog = createDialog({
  portal: 'body', // Default
  // portal: '#my-container', // Custom target
  // portal: null, // Disable portalling
});
```

### Focus Management

Melt-UI automatically handles:
- Focus trapping in modals/dialogs
- Focus restoration when closing
- Arrow key navigation in menus/selects

### Escape Behavior

Configure how components respond to Escape key:

```typescript
const dialog = createDialog({
  escapeBehavior: 'close', // Default: closes immediately
  // escapeBehavior: 'ignore', // Ignores escape key
  // escapeBehavior: 'defer-otherwise-close', // Delegates to parent
});
```

## VacayTracker Theme Integration

Apply VacayTracker's ocean/sand theme to Melt-UI components:

```svelte
<!-- Dialog with VacayTracker styling -->
<div
  use:melt={$overlay}
  class="fixed inset-0 bg-ocean-900/50 backdrop-blur-sm"
></div>

<div
  use:melt={$content}
  class="bg-white/95 backdrop-blur-md rounded-2xl shadow-xl border border-white/30"
>
  <!-- Content -->
</div>
```

## Keyboard Navigation Summary

| Component | Keys |
|-----------|------|
| Dialog | `Escape` to close, `Tab` to navigate |
| Dropdown | `Arrow` keys, `Enter/Space` to select |
| Select | `Arrow` keys, `Enter` to select, type to search |
| Tabs | `Arrow` keys to switch tabs |
| Accordion | `Arrow` keys, `Enter/Space` to toggle |

## Available Components

### Layout & Navigation
- [Accordion](./accordion.md) - Expandable content sections
- [Collapsible](./collapsible.md) - Simple expand/collapse
- [Tabs](./tabs.md) - Tabbed content panels
- [Menubar](./menubar.md) - Horizontal menu bar
- [Toolbar](./toolbar.md) - Button toolbar

### Overlays
- [Dialog](./dialog.md) - Modal dialogs and alerts
- [Popover](./popover.md) - Floating content panels
- [Tooltip](./tooltip.md) - Hover tooltips
- [Context Menu](./context-menu.md) - Right-click menus
- [Dropdown Menu](./dropdown-menu.md) - Trigger menus

### Form Controls
- [Checkbox](./checkbox.md) - Checkbox with tri-state
- [Radio Group](./radio-group.md) - Radio button groups
- [Switch](./switch.md) - Toggle switches
- [Select](./select.md) - Dropdown selection
- [Combobox](./combobox.md) - Searchable select
- [Slider](./slider.md) - Range slider
- [Tags Input](./tags-input.md) - Tag input field
- [PIN Input](./pin-input.md) - OTP/PIN input

### Date & Time
- [Calendar](./calendar.md) - Date grid
- [Date Field](./date-field.md) - Date input
- [Date Picker](./date-picker.md) - Calendar popup
- [Date Range Field](./date-range-field.md) - Range input
- [Date Range Picker](./date-range-picker.md) - Range calendar
- [Range Calendar](./range-calendar.md) - Date range selection

### Display
- [Avatar](./avatar.md) - User images
- [Progress](./progress.md) - Progress bar
- [Toast](./toast.md) - Notifications
- [Separator](./separator.md) - Visual dividers
- [Label](./label.md) - Form labels
- [Link Preview](./link-preview.md) - Hover previews

### Utility
- [Pagination](./pagination.md) - Page navigation
- [Scroll Area](./scroll-area.md) - Custom scrollbars
- [Toggle](./toggle.md) - Toggle button
- [Toggle Group](./toggle-group.md) - Button group
- [Tree](./tree.md) - Hierarchical tree
- [Table of Contents](./table-of-contents.md) - Document nav

## Accessibility

All Melt-UI components follow WAI-ARIA design patterns:
- Proper ARIA attributes automatically applied
- Full keyboard navigation support
- Screen reader compatible
- Focus management handled

## Best Practices

1. **Use semantic elements**: Apply builders to appropriate HTML elements
2. **Preserve styling**: Use data attributes for state-based styles
3. **Handle transitions**: Use `forceVisible` for CSS transitions
4. **Test accessibility**: Verify keyboard and screen reader support
