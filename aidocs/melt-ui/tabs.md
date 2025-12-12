# Tabs

A tabbed interface for organizing content into panels.

## Use Cases

- Settings pages
- Dashboard sections
- Content organization
- Form wizards

## Installation

```typescript
import { createTabs } from '@melt-ui/svelte';
```

## API Reference

### createTabs Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `loop` | `boolean` | `false` | Loop keyboard navigation |
| `activateOnFocus` | `boolean` | `true` | Activate tab on focus |
| `orientation` | `'horizontal' \| 'vertical'` | `'horizontal'` | Tab direction |
| `autoSet` | `boolean` | `true` | Auto-set first tab |
| `defaultValue` | `string` | — | Initial active tab |
| `value` | `Writable<string>` | — | Controlled active tab |
| `onValueChange` | `ChangeFn` | — | Value change callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Tabs container |
| `list` | Tab trigger container |
| `trigger` | Individual tab button |
| `content` | Tab panel content |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `value` | `Writable<string>` | Currently active tab |

## Data Attributes

### List
- `[data-orientation]` - `'horizontal'` or `'vertical'`
- `[data-melt-tabs-list]` - Present on list

### Trigger
- `[data-state]` - `'active'` or `'inactive'`
- `[data-value]` - Tab value
- `[data-disabled]` - Present when disabled
- `[data-orientation]` - Inherited orientation
- `[data-melt-tabs-trigger]` - Present on triggers

### Content
- `[data-melt-tabs-content]` - Present on content

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Tab` | Focus active trigger / Move to content |
| `Arrow Left/Up` | Previous trigger |
| `Arrow Right/Down` | Next trigger |
| `Home` | First trigger |
| `End` | Last trigger |

## Accessibility

Follows [WAI-ARIA Tabs Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/tabs/):
- `role="tablist"`, `role="tab"`, `role="tabpanel"`
- `aria-selected` state
- Keyboard navigation

## Example

```svelte
<script lang="ts">
  import { createTabs, melt } from '@melt-ui/svelte';

  const {
    elements: { root, list, trigger, content },
    states: { value }
  } = createTabs({
    defaultValue: 'account'
  });

  const tabs = [
    { id: 'account', label: 'Account' },
    { id: 'password', label: 'Password' },
    { id: 'notifications', label: 'Notifications' }
  ];
</script>

<div use:melt={$root} class="w-full">
  <div
    use:melt={$list}
    class="flex border-b border-ocean-200"
  >
    {#each tabs as tab}
      <button
        use:melt={$trigger(tab.id)}
        class="px-4 py-2 text-sm font-medium transition-colors
          data-[state=active]:text-ocean-700 data-[state=active]:border-b-2 data-[state=active]:border-ocean-500
          data-[state=inactive]:text-ocean-500 data-[state=inactive]:hover:text-ocean-700"
      >
        {tab.label}
      </button>
    {/each}
  </div>

  <div use:melt={$content('account')} class="p-4">
    <h3 class="font-semibold text-ocean-800 mb-2">Account Settings</h3>
    <p class="text-ocean-600">Manage your account details.</p>
  </div>

  <div use:melt={$content('password')} class="p-4">
    <h3 class="font-semibold text-ocean-800 mb-2">Password</h3>
    <p class="text-ocean-600">Update your password.</p>
  </div>

  <div use:melt={$content('notifications')} class="p-4">
    <h3 class="font-semibold text-ocean-800 mb-2">Notifications</h3>
    <p class="text-ocean-600">Configure notification preferences.</p>
  </div>
</div>
```

## Styling with Tailwind

```css
/* Tab list */
[data-melt-tabs-list] {
  @apply flex border-b border-ocean-200;
}

/* Tab trigger */
[data-melt-tabs-trigger] {
  @apply px-4 py-2 text-sm font-medium transition-colors;
}

[data-melt-tabs-trigger][data-state="active"] {
  @apply text-ocean-700 border-b-2 border-ocean-500 -mb-px;
}

[data-melt-tabs-trigger][data-state="inactive"] {
  @apply text-ocean-500 hover:text-ocean-700;
}

[data-melt-tabs-trigger][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}

/* Focus ring */
[data-melt-tabs-trigger]:focus-visible {
  @apply outline-none ring-2 ring-ocean-500/50 ring-inset;
}
```

## Pill-Style Tabs

```svelte
<script lang="ts">
  const { elements: { root, list, trigger, content } } = createTabs();
</script>

<div use:melt={$root}>
  <div
    use:melt={$list}
    class="inline-flex p-1 bg-ocean-100 rounded-lg"
  >
    {#each tabs as tab}
      <button
        use:melt={$trigger(tab.id)}
        class="px-4 py-2 text-sm font-medium rounded-md transition-all
          data-[state=active]:bg-white data-[state=active]:text-ocean-800 data-[state=active]:shadow-sm
          data-[state=inactive]:text-ocean-600"
      >
        {tab.label}
      </button>
    {/each}
  </div>

  <!-- Content panels -->
</div>
```

## Vertical Tabs

```svelte
<script lang="ts">
  const { elements: { root, list, trigger, content } } = createTabs({
    orientation: 'vertical',
    defaultValue: 'general'
  });
</script>

<div use:melt={$root} class="flex gap-6">
  <div
    use:melt={$list}
    class="flex flex-col w-48 border-r border-ocean-200"
  >
    {#each tabs as tab}
      <button
        use:melt={$trigger(tab.id)}
        class="px-4 py-2 text-left text-sm font-medium transition-colors
          data-[state=active]:bg-ocean-100 data-[state=active]:text-ocean-800 data-[state=active]:border-r-2 data-[state=active]:border-ocean-500
          data-[state=inactive]:text-ocean-500 data-[state=inactive]:hover:text-ocean-700"
      >
        {tab.label}
      </button>
    {/each}
  </div>

  <div class="flex-1">
    {#each tabs as tab}
      <div use:melt={$content(tab.id)}>
        <!-- Content -->
      </div>
    {/each}
  </div>
</div>
```

## VacayTracker Admin Tabs

```svelte
<script lang="ts">
  import { createTabs, melt } from '@melt-ui/svelte';
  import { Users, Wallet, Settings } from 'lucide-svelte';

  const { elements: { root, list, trigger, content } } = createTabs({
    defaultValue: 'users'
  });

  const tabs = [
    { id: 'users', label: 'Users', icon: Users },
    { id: 'balances', label: 'Balances', icon: Wallet },
    { id: 'settings', label: 'Settings', icon: Settings }
  ];
</script>

<div use:melt={$root} class="bg-white rounded-xl shadow-lg overflow-hidden">
  <div use:melt={$list} class="flex bg-ocean-50 border-b border-ocean-200">
    {#each tabs as tab}
      <button
        use:melt={$trigger(tab.id)}
        class="flex items-center gap-2 px-6 py-4 text-sm font-medium transition-colors
          data-[state=active]:bg-white data-[state=active]:text-ocean-800 data-[state=active]:border-b-2 data-[state=active]:border-ocean-500
          data-[state=inactive]:text-ocean-500 data-[state=inactive]:hover:text-ocean-700"
      >
        <tab.icon class="h-4 w-4" />
        {tab.label}
      </button>
    {/each}
  </div>

  <div use:melt={$content('users')} class="p-6">
    <!-- Users content -->
  </div>

  <div use:melt={$content('balances')} class="p-6">
    <!-- Balances content -->
  </div>

  <div use:melt={$content('settings')} class="p-6">
    <!-- Settings content -->
  </div>
</div>
```

## With Badge Count

```svelte
{#each tabs as tab}
  <button use:melt={$trigger(tab.id)} class="...">
    {tab.label}
    {#if tab.count > 0}
      <span class="ml-2 px-2 py-0.5 text-xs bg-ocean-200 text-ocean-700 rounded-full">
        {tab.count}
      </span>
    {/if}
  </button>
{/each}
```

## Controlled Tabs

```svelte
<script lang="ts">
  import { writable } from 'svelte/store';
  import { page } from '$app/stores';

  // Sync with URL
  const activeTab = writable($page.url.searchParams.get('tab') || 'account');

  const tabs = createTabs({
    value: activeTab,
    onValueChange: ({ next }) => {
      // Update URL without navigation
      const url = new URL(window.location.href);
      url.searchParams.set('tab', next);
      history.replaceState(null, '', url);
      return next;
    }
  });
</script>
```
