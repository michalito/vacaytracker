# Dropdown Menu

A trigger-activated menu displaying a list of actions or options.

## Use Cases

- User account menus
- Navigation dropdowns
- Action menus
- Settings menus
- Context-specific options

## Installation

```typescript
import { createDropdownMenu } from '@melt-ui/svelte';
```

## API Reference

### createDropdownMenu Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `arrowSize` | `number` | `8` | Arrow size in pixels |
| `dir` | `'ltr' \| 'rtl'` | `'ltr'` | Text direction |
| `preventScroll` | `boolean` | `true` | Block page scroll |
| `closeOnOutsideClick` | `boolean` | `true` | Close on click outside |
| `loop` | `boolean` | `false` | Loop keyboard navigation |
| `forceVisible` | `boolean` | `false` | Force visibility |
| `portal` | `string \| HTMLElement \| null` | `'body'` | Portal target |
| `positioning` | `FloatingConfig` | — | Floating positioning |
| `escapeBehavior` | `EscapeBehavior` | `'close'` | Escape key behavior |
| `defaultOpen` | `boolean` | `false` | Initial open state |
| `open` | `Writable<boolean>` | — | Controlled state |
| `onOpenChange` | `ChangeFn` | — | Open state callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `trigger` | Button that opens the menu |
| `menu` | Menu container |
| `item` | Standard menu item |
| `checkboxItem` | Toggleable menu item |
| `radioGroup` | Radio group container |
| `radioItem` | Radio selection item |
| `separator` | Visual divider |
| `submenu` | Nested menu container |
| `subTrigger` | Submenu trigger |
| `arrow` | Optional arrow element |
| `overlay` | Optional backdrop |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `open` | `Writable<boolean>` | Menu visibility |

### Returned Builders

| Builder | Description |
|---------|-------------|
| `createCheckboxItem` | Create checkbox item with state |
| `createRadioGroup` | Create radio group with state |
| `createSubmenu` | Create nested submenu |

## Data Attributes

### Trigger
- `[data-state]` - `'open'` or `'closed'`
- `[data-melt-dropdown-menu-trigger]` - Present on trigger

### Item
- `[data-highlighted]` - Currently highlighted
- `[data-disabled]` - Item is disabled
- `[data-melt-dropdown-menu-item]` - Present on items

### Checkbox Item
- `[data-state]` - `'checked'` or `'unchecked'`
- `[data-highlighted]` - Currently highlighted

### Radio Item
- `[data-state]` - `'checked'` or `'unchecked'`
- `[data-highlighted]` - Currently highlighted

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Space` / `Enter` | Open menu / Select item |
| `Arrow Down` | Next item |
| `Arrow Up` | Previous item |
| `Arrow Right` | Open submenu |
| `Arrow Left` | Close submenu |
| `Escape` | Close menu |
| `Home` | First item |
| `End` | Last item |

## Accessibility

Follows [WAI-ARIA Menu Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/menu/):
- Proper ARIA roles (menu, menuitem, menuitemcheckbox, menuitemradio)
- Full keyboard navigation
- Focus management

## Example

```svelte
<script lang="ts">
  import { createDropdownMenu, melt } from '@melt-ui/svelte';
  import { Shield, ChevronDown, Users, Wallet, Settings, LogOut } from 'lucide-svelte';

  const {
    elements: { trigger, menu, item, separator },
    states: { open }
  } = createDropdownMenu({
    forceVisible: true,
    positioning: { placement: 'bottom-end' }
  });

  const menuItems = [
    { icon: Users, label: 'Users', href: '/admin/users' },
    { icon: Wallet, label: 'Balances', href: '/admin/balances' },
    { icon: Settings, label: 'Settings', href: '/admin/settings' }
  ];
</script>

<div class="relative">
  <button
    use:melt={$trigger}
    class="flex items-center gap-2 px-4 py-2 rounded-xl font-medium text-ocean-600 hover:text-ocean-800 hover:bg-ocean-500/10 transition-all"
  >
    <Shield class="w-4 h-4" />
    <span>Admin</span>
    <ChevronDown class="w-4 h-4 transition-transform {$open ? 'rotate-180' : ''}" />
  </button>

  {#if $open}
    <div
      use:melt={$menu}
      class="absolute right-0 mt-2 w-48 bg-white/95 backdrop-blur-md rounded-xl shadow-xl border border-white/30 py-1 z-50"
    >
      {#each menuItems as { icon: Icon, label, href }}
        <a
          use:melt={$item}
          {href}
          class="flex items-center gap-3 px-4 py-2 text-ocean-600 transition-colors
            data-[highlighted]:bg-ocean-500/10 data-[highlighted]:text-ocean-800"
        >
          <Icon class="w-4 h-4" />
          <span>{label}</span>
        </a>
      {/each}

      <div use:melt={$separator} class="my-1 h-px bg-ocean-200"></div>

      <button
        use:melt={$item}
        class="flex items-center gap-3 w-full px-4 py-2 text-red-600 transition-colors
          data-[highlighted]:bg-red-50"
      >
        <LogOut class="w-4 h-4" />
        <span>Sign Out</span>
      </button>
    </div>
  {/if}
</div>
```

## Styling with Tailwind

```css
/* Trigger states */
[data-melt-dropdown-menu-trigger][data-state="open"] {
  @apply bg-ocean-500/15 text-ocean-700;
}

/* Item states */
[data-melt-dropdown-menu-item][data-highlighted] {
  @apply bg-ocean-500/10 text-ocean-800;
}

[data-melt-dropdown-menu-item][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}

/* Checkbox/Radio checked state */
[data-melt-dropdown-menu-checkbox-item][data-state="checked"],
[data-melt-dropdown-menu-radio-item][data-state="checked"] {
  @apply bg-ocean-50;
}
```

## VacayTracker Admin Dropdown

```svelte
<script lang="ts">
  import { createDropdownMenu, melt } from '@melt-ui/svelte';
  import { page } from '$app/stores';
  import { Shield, ChevronDown, Users, Wallet, Settings } from 'lucide-svelte';
  import { clsx } from 'clsx';

  const {
    elements: { trigger, menu, item },
    states: { open }
  } = createDropdownMenu({
    forceVisible: true,
    positioning: { placement: 'bottom-end' }
  });

  const adminNavItems = [
    { href: '/admin/users', icon: Users, label: 'Users' },
    { href: '/admin/balances', icon: Wallet, label: 'Balances' },
    { href: '/admin/settings', icon: Settings, label: 'Settings' }
  ];

  function isActive(href: string): boolean {
    return $page.url.pathname.startsWith(href);
  }
</script>

<button
  use:melt={$trigger}
  class={clsx(
    'flex items-center gap-2 px-4 py-2 rounded-xl font-medium transition-all cursor-pointer',
    $page.url.pathname.startsWith('/admin')
      ? 'bg-ocean-500/15 text-ocean-700 shadow-sm'
      : 'text-ocean-600 hover:text-ocean-800 hover:bg-ocean-500/10'
  )}
>
  <Shield class="w-4 h-4" />
  <span>Admin</span>
  <ChevronDown class="w-4 h-4 transition-transform {$open ? 'rotate-180' : ''}" />
</button>

{#if $open}
  <div
    use:melt={$menu}
    class="absolute right-0 mt-2 w-48 bg-white/95 backdrop-blur-md rounded-xl shadow-xl border border-white/30 py-1 z-50"
  >
    {#each adminNavItems as { href, icon: Icon, label }}
      {@const active = isActive(href)}
      <a
        use:melt={$item}
        {href}
        class={clsx(
          'flex items-center gap-3 px-4 py-2 transition-all',
          active
            ? 'bg-ocean-500/15 text-ocean-700'
            : 'text-ocean-600 data-[highlighted]:bg-ocean-500/10 data-[highlighted]:text-ocean-800'
        )}
      >
        <Icon class="w-4 h-4" />
        <span>{label}</span>
      </a>
    {/each}
  </div>
{/if}
```

## With Checkbox Items

```svelte
<script lang="ts">
  import { createDropdownMenu, melt } from '@melt-ui/svelte';
  import { Check } from 'lucide-svelte';

  const {
    elements: { trigger, menu },
    builders: { createCheckboxItem }
  } = createDropdownMenu();

  const showNotifications = createCheckboxItem({ defaultChecked: true });
  const darkMode = createCheckboxItem({ defaultChecked: false });
</script>

<div use:melt={$menu}>
  <button
    use:melt={showNotifications.elements.checkboxItem}
    class="flex items-center gap-2 w-full px-4 py-2"
  >
    <span class="w-4">
      {#if $showNotifications.states.checked}
        <Check class="h-4 w-4" />
      {/if}
    </span>
    <span>Show Notifications</span>
  </button>

  <button
    use:melt={darkMode.elements.checkboxItem}
    class="flex items-center gap-2 w-full px-4 py-2"
  >
    <span class="w-4">
      {#if $darkMode.states.checked}
        <Check class="h-4 w-4" />
      {/if}
    </span>
    <span>Dark Mode</span>
  </button>
</div>
```

## With Radio Group

```svelte
<script lang="ts">
  const {
    elements: { menu },
    builders: { createRadioGroup }
  } = createDropdownMenu();

  const themeGroup = createRadioGroup({ defaultValue: 'system' });
</script>

<div use:melt={$menu}>
  <div use:melt={themeGroup.elements.radioGroup}>
    <span class="px-4 py-1 text-xs font-semibold text-ocean-500 uppercase">Theme</span>
    {#each ['light', 'dark', 'system'] as theme}
      <button
        use:melt={themeGroup.elements.radioItem({ value: theme })}
        class="flex items-center gap-2 w-full px-4 py-2 capitalize"
      >
        <span class="w-4">
          {#if $themeGroup.states.value === theme}
            <div class="h-2 w-2 rounded-full bg-ocean-500"></div>
          {/if}
        </span>
        <span>{theme}</span>
      </button>
    {/each}
  </div>
</div>
```

## With Submenu

```svelte
<script lang="ts">
  const {
    elements: { menu, item },
    builders: { createSubmenu }
  } = createDropdownMenu();

  const exportSubmenu = createSubmenu();
</script>

<div use:melt={$menu}>
  <button use:melt={$item}>New</button>
  <button use:melt={$item}>Edit</button>

  <button use:melt={exportSubmenu.elements.subTrigger} class="flex items-center justify-between">
    <span>Export As</span>
    <ChevronRight class="h-4 w-4" />
  </button>

  <div use:melt={exportSubmenu.elements.submenu} class="min-w-[120px] bg-white shadow-lg rounded-lg p-1">
    <button use:melt={$item}>PDF</button>
    <button use:melt={$item}>CSV</button>
    <button use:melt={$item}>JSON</button>
  </div>
</div>
```
