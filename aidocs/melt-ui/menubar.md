# Menubar

A horizontal menu bar with dropdown menus, commonly used for application navigation.

## Use Cases

- Application menu bars
- Desktop-style navigation
- Complex navigation structures
- Multi-level menu systems

## Installation

```typescript
import { createMenubar } from '@melt-ui/svelte';
```

## API Reference

### createMenubar Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `loop` | `boolean` | `true` | Loop through menu triggers |
| `closeOnEscape` | `boolean` | `true` | Close on Escape key |

### Returned Elements

| Element | Description |
|---------|-------------|
| `menubar` | Root menubar container |

### Returned Builders

| Builder | Description |
|---------|-------------|
| `createMenu` | Create a dropdown menu for the menubar |

### Menu Builder Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `positioning` | `FloatingConfig` | — | Menu positioning |
| `arrowSize` | `number` | `8` | Arrow size |
| `preventScroll` | `boolean` | `true` | Block scroll when open |
| `closeOnOutsideClick` | `boolean` | `true` | Close on outside click |
| `loop` | `boolean` | `false` | Loop through items |
| `forceVisible` | `boolean` | `false` | Force visibility |
| `portal` | `string \| HTMLElement \| null` | `'body'` | Portal target |

### Menu Returned Elements

| Element | Description |
|---------|-------------|
| `trigger` | Menu trigger button |
| `menu` | Menu container |
| `item` | Menu item |
| `checkboxItem` | Checkbox menu item |
| `radioGroup` | Radio group |
| `radioItem` | Radio item |
| `separator` | Visual separator |
| `submenu` | Nested submenu |
| `subTrigger` | Submenu trigger |

## Data Attributes

### Menubar
- `[data-melt-menubar]` - Present on menubar

### Trigger
- `[data-state]` - `'open'` or `'closed'`
- `[data-highlighted]` - Currently highlighted
- `[data-melt-menubar-trigger]` - Present on triggers

### Item
- `[data-highlighted]` - Currently highlighted
- `[data-disabled]` - Item is disabled
- `[data-melt-menubar-menu-item]` - Present on items

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Arrow Left/Right` | Navigate between menus |
| `Arrow Down` | Open menu / Next item |
| `Arrow Up` | Previous item |
| `Enter` / `Space` | Select item |
| `Escape` | Close menu |

## Accessibility

Follows [WAI-ARIA Menubar Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/menubar/):
- Proper roles (menubar, menu, menuitem)
- Full keyboard navigation
- Focus management between menus

## Example

```svelte
<script lang="ts">
  import { createMenubar, melt } from '@melt-ui/svelte';
  import { ChevronRight, Check } from 'lucide-svelte';

  const {
    elements: { menubar },
    builders: { createMenu }
  } = createMenubar();

  const fileMenu = createMenu({ forceVisible: true });
  const editMenu = createMenu({ forceVisible: true });
  const viewMenu = createMenu({ forceVisible: true });
</script>

<div
  use:melt={$menubar}
  class="flex items-center gap-1 bg-white border border-ocean-200 rounded-lg p-1"
>
  <!-- File Menu -->
  <div>
    <button
      use:melt={fileMenu.elements.trigger}
      class="px-3 py-1.5 rounded text-sm font-medium text-ocean-700
        data-[highlighted]:bg-ocean-100 data-[state=open]:bg-ocean-100"
    >
      File
    </button>

    {#if $fileMenu.states.open}
      <div
        use:melt={fileMenu.elements.menu}
        class="z-50 min-w-[180px] rounded-lg bg-white shadow-xl border border-ocean-200 p-1"
      >
        <button use:melt={fileMenu.elements.item} class="menu-item">
          New File
          <span class="text-ocean-400 text-xs">Ctrl+N</span>
        </button>
        <button use:melt={fileMenu.elements.item} class="menu-item">
          Open...
          <span class="text-ocean-400 text-xs">Ctrl+O</span>
        </button>
        <div use:melt={fileMenu.elements.separator} class="my-1 h-px bg-ocean-200"></div>
        <button use:melt={fileMenu.elements.item} class="menu-item">
          Save
          <span class="text-ocean-400 text-xs">Ctrl+S</span>
        </button>
        <button use:melt={fileMenu.elements.item} class="menu-item">
          Save As...
        </button>
      </div>
    {/if}
  </div>

  <!-- Edit Menu -->
  <div>
    <button
      use:melt={editMenu.elements.trigger}
      class="px-3 py-1.5 rounded text-sm font-medium text-ocean-700
        data-[highlighted]:bg-ocean-100 data-[state=open]:bg-ocean-100"
    >
      Edit
    </button>

    {#if $editMenu.states.open}
      <div
        use:melt={editMenu.elements.menu}
        class="z-50 min-w-[180px] rounded-lg bg-white shadow-xl border border-ocean-200 p-1"
      >
        <button use:melt={editMenu.elements.item} class="menu-item">
          Undo
          <span class="text-ocean-400 text-xs">Ctrl+Z</span>
        </button>
        <button use:melt={editMenu.elements.item} class="menu-item">
          Redo
          <span class="text-ocean-400 text-xs">Ctrl+Y</span>
        </button>
        <div use:melt={editMenu.elements.separator} class="my-1 h-px bg-ocean-200"></div>
        <button use:melt={editMenu.elements.item} class="menu-item">Cut</button>
        <button use:melt={editMenu.elements.item} class="menu-item">Copy</button>
        <button use:melt={editMenu.elements.item} class="menu-item">Paste</button>
      </div>
    {/if}
  </div>

  <!-- View Menu -->
  <div>
    <button
      use:melt={viewMenu.elements.trigger}
      class="px-3 py-1.5 rounded text-sm font-medium text-ocean-700
        data-[highlighted]:bg-ocean-100 data-[state=open]:bg-ocean-100"
    >
      View
    </button>

    {#if $viewMenu.states.open}
      <div
        use:melt={viewMenu.elements.menu}
        class="z-50 min-w-[180px] rounded-lg bg-white shadow-xl border border-ocean-200 p-1"
      >
        <button use:melt={viewMenu.elements.item} class="menu-item">
          Zoom In
        </button>
        <button use:melt={viewMenu.elements.item} class="menu-item">
          Zoom Out
        </button>
      </div>
    {/if}
  </div>
</div>

<style>
  .menu-item {
    @apply flex w-full items-center justify-between rounded px-3 py-1.5 text-sm text-ocean-700 outline-none;
    @apply data-[highlighted]:bg-ocean-100;
  }
</style>
```

## Styling with Tailwind

```css
/* Menubar container */
[data-melt-menubar] {
  @apply flex items-center gap-1;
}

/* Trigger states */
[data-melt-menubar-trigger][data-state="open"],
[data-melt-menubar-trigger][data-highlighted] {
  @apply bg-ocean-100;
}

/* Menu item states */
[data-melt-menubar-menu-item][data-highlighted] {
  @apply bg-ocean-100;
}

[data-melt-menubar-menu-item][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}
```

## With Submenus

```svelte
<script lang="ts">
  const shareSubmenu = fileMenu.builders.createSubmenu();
</script>

<div use:melt={fileMenu.elements.menu}>
  <button use:melt={fileMenu.elements.item}>New</button>

  <button
    use:melt={shareSubmenu.elements.subTrigger}
    class="flex items-center justify-between"
  >
    <span>Share</span>
    <ChevronRight class="h-4 w-4" />
  </button>

  <div
    use:melt={shareSubmenu.elements.submenu}
    class="min-w-[140px] rounded-lg bg-white shadow-xl p-1"
  >
    <button use:melt={fileMenu.elements.item}>Email</button>
    <button use:melt={fileMenu.elements.item}>Slack</button>
    <button use:melt={fileMenu.elements.item}>Copy Link</button>
  </div>
</div>
```

## With Checkbox Items

```svelte
<script lang="ts">
  const showLineNumbers = viewMenu.builders.createCheckboxItem({
    defaultChecked: true
  });
</script>

<button
  use:melt={showLineNumbers.elements.checkboxItem}
  class="flex items-center gap-2 px-3 py-1.5"
>
  <span class="w-4">
    {#if $showLineNumbers.states.checked}
      <Check class="h-4 w-4" />
    {/if}
  </span>
  <span>Show Line Numbers</span>
</button>
```

## With Radio Group

```svelte
<script lang="ts">
  const themeGroup = viewMenu.builders.createRadioGroup({
    defaultValue: 'light'
  });
</script>

<div use:melt={themeGroup.elements.radioGroup}>
  <span class="px-3 py-1 text-xs text-ocean-500">Theme</span>
  {#each ['light', 'dark', 'system'] as theme}
    <button
      use:melt={themeGroup.elements.radioItem({ value: theme })}
      class="flex items-center gap-2 px-3 py-1.5 capitalize"
    >
      <span class="w-4">
        {#if $themeGroup.states.value === theme}
          <div class="h-2 w-2 rounded-full bg-ocean-500"></div>
        {/if}
      </span>
      {theme}
    </button>
  {/each}
</div>
```
