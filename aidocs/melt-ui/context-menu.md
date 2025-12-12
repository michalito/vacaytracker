# Context Menu

A right-click triggered menu displaying actions relevant to the clicked element.

## Use Cases

- File/folder context actions
- Table row actions
- Image editing options
- Custom right-click menus

## Installation

```typescript
import { createContextMenu } from '@melt-ui/svelte';
```

## API Reference

### createContextMenu Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `arrowSize` | `number` | `8` | Arrow size in pixels |
| `dir` | `'ltr' \| 'rtl'` | `'ltr'` | Text direction |
| `preventScroll` | `boolean` | `true` | Block page scroll when open |
| `closeOnOutsideClick` | `boolean` | `true` | Close on click outside |
| `loop` | `boolean` | `false` | Loop focus navigation |
| `forceVisible` | `boolean` | `false` | Force visibility |
| `portal` | `string \| HTMLElement \| null` | `'body'` | Portal target |
| `positioning` | `FloatingConfig` | — | Floating positioning |
| `escapeBehavior` | `EscapeBehavior` | `'close'` | Escape key behavior |
| `onOpenChange` | `ChangeFn` | — | Open state callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `trigger` | Element that triggers context menu on right-click |
| `menu` | Menu container |
| `item` | Menu item |
| `checkboxItem` | Toggleable menu item |
| `radioGroup` | Radio group container |
| `radioItem` | Radio selection item |
| `separator` | Visual divider |
| `submenu` | Nested menu |
| `subTrigger` | Submenu trigger |
| `arrow` | Optional arrow element |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `open` | `Writable<boolean>` | Menu visibility state |

## Data Attributes

### Trigger
- `[data-state]` - `'open'` or `'closed'`
- `[data-melt-context-menu-trigger]` - Present on trigger

### Item
- `[data-highlighted]` - Currently highlighted
- `[data-disabled]` - Item is disabled
- `[data-melt-context-menu-item]` - Present on items

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Arrow Down` | Next item |
| `Arrow Up` | Previous item |
| `Arrow Right` | Open submenu |
| `Arrow Left` | Close submenu |
| `Enter` / `Space` | Select item |
| `Escape` | Close menu |

## Accessibility

Follows [WAI-ARIA Menu Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/menu/):
- Proper menu roles
- Keyboard navigation
- Focus management

## Example

```svelte
<script lang="ts">
  import { createContextMenu, melt } from '@melt-ui/svelte';
  import { Copy, Trash, Edit, Share } from 'lucide-svelte';

  const {
    elements: { trigger, menu, item, separator },
    states: { open }
  } = createContextMenu();

  const menuItems = [
    { icon: Edit, label: 'Edit', action: () => console.log('Edit') },
    { icon: Copy, label: 'Duplicate', action: () => console.log('Duplicate') },
    { icon: Share, label: 'Share', action: () => console.log('Share') },
    { icon: Trash, label: 'Delete', action: () => console.log('Delete'), destructive: true }
  ];
</script>

<div
  use:melt={$trigger}
  class="w-64 h-32 border-2 border-dashed border-ocean-300 rounded-xl flex items-center justify-center text-ocean-500 cursor-context-menu"
>
  Right-click me
</div>

{#if $open}
  <div
    use:melt={$menu}
    class="z-50 min-w-[180px] rounded-xl bg-white shadow-lg border border-ocean-200 p-1"
  >
    {#each menuItems as { icon: Icon, label, action, destructive }, i}
      {#if i === menuItems.length - 1}
        <div use:melt={$separator} class="my-1 h-px bg-ocean-200"></div>
      {/if}
      <button
        use:melt={$item}
        onclick={action}
        class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm outline-none
          data-[highlighted]:bg-ocean-100
          {destructive ? 'text-red-600 data-[highlighted]:bg-red-50' : 'text-ocean-700'}"
      >
        <Icon class="h-4 w-4" />
        <span>{label}</span>
      </button>
    {/each}
  </div>
{/if}
```

## Styling with Tailwind

```css
/* Menu item states */
[data-melt-context-menu-item][data-highlighted] {
  @apply bg-ocean-100;
}

[data-melt-context-menu-item][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}

/* Destructive item */
[data-melt-context-menu-item].destructive {
  @apply text-red-600;
}

[data-melt-context-menu-item].destructive[data-highlighted] {
  @apply bg-red-50;
}
```

## With Checkbox Items

```svelte
<script lang="ts">
  import { createContextMenu, melt } from '@melt-ui/svelte';
  import { Check } from 'lucide-svelte';

  const {
    elements: { trigger, menu, checkboxItem },
    builders: { createCheckboxItem }
  } = createContextMenu();

  let showHidden = $state(false);
  let showExtensions = $state(true);
</script>

<div use:melt={$menu}>
  <button
    use:melt={$checkboxItem}
    data-checked={showHidden}
    onclick={() => (showHidden = !showHidden)}
    class="flex items-center gap-2 px-3 py-2 rounded-lg"
  >
    <span class="w-4">
      {#if showHidden}<Check class="h-4 w-4" />{/if}
    </span>
    Show Hidden Files
  </button>
</div>
```

## With Submenu

```svelte
<script lang="ts">
  const {
    elements: { trigger, menu, item, submenu, subTrigger }
  } = createContextMenu();
</script>

<div use:melt={$menu}>
  <button use:melt={$item}>New File</button>

  <button use:melt={$subTrigger} class="flex items-center justify-between">
    <span>Export As</span>
    <ChevronRight class="h-4 w-4" />
  </button>

  <div use:melt={$submenu} class="min-w-[140px] rounded-xl bg-white shadow-lg p-1">
    <button use:melt={$item}>PDF</button>
    <button use:melt={$item}>PNG</button>
    <button use:melt={$item}>SVG</button>
  </div>

  <button use:melt={$item}>Delete</button>
</div>
```

## With Radio Group

```svelte
<script lang="ts">
  const {
    elements: { menu, radioGroup, radioItem },
    builders: { createRadioGroup }
  } = createContextMenu();

  const sortGroup = createRadioGroup({ defaultValue: 'name' });
</script>

<div use:melt={$menu}>
  <div use:melt={$radioGroup}>
    <span class="px-3 py-1 text-xs font-semibold text-ocean-500">Sort By</span>
    <button use:melt={$radioItem({ value: 'name' })}>
      Name
    </button>
    <button use:melt={$radioItem({ value: 'date' })}>
      Date Modified
    </button>
    <button use:melt={$radioItem({ value: 'size' })}>
      Size
    </button>
  </div>
</div>
```
