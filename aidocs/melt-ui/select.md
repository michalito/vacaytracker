# Select

A dropdown control for selecting from a list of options.

## Use Cases

- Form dropdowns
- Filter selections
- Setting options
- Category selection

## Installation

```typescript
import { createSelect } from '@melt-ui/svelte';
```

## API Reference

### createSelect Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `defaultSelected` | `SelectOption` | — | Initial selected option |
| `selected` | `Writable<SelectOption>` | — | Controlled selection |
| `defaultOpen` | `boolean` | `false` | Initial open state |
| `open` | `Writable<boolean>` | — | Controlled open state |
| `multiple` | `boolean` | `false` | Enable multi-select |
| `disabled` | `boolean` | `false` | Disable select |
| `loop` | `boolean` | `false` | Loop keyboard navigation |
| `preventScroll` | `boolean` | `true` | Block page scroll |
| `closeOnOutsideClick` | `boolean` | `true` | Close on outside click |
| `portal` | `string \| HTMLElement \| null` | `'body'` | Portal target |
| `forceVisible` | `boolean` | `false` | Force visibility |
| `positioning` | `FloatingConfig` | — | Floating positioning |
| `name` | `string` | — | Form input name |
| `required` | `boolean` | `false` | Mark as required |
| `onSelectedChange` | `ChangeFn` | — | Selection callback |
| `onOpenChange` | `ChangeFn` | — | Open state callback |

### SelectOption Type

```typescript
type SelectOption<T = unknown> = {
  value: T;
  label?: string;
  disabled?: boolean;
};
```

### Returned Elements

| Element | Description |
|---------|-------------|
| `trigger` | Button that opens the menu |
| `menu` | Options container |
| `option` | Individual option element |
| `label` | Select label |
| `group` | Option group container |
| `groupLabel` | Group heading |
| `hiddenInput` | Hidden form input |
| `arrow` | Optional arrow element |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `open` | `Writable<boolean>` | Menu visibility |
| `selected` | `Writable<SelectOption>` | Selected option(s) |
| `selectedLabel` | `Readable<string>` | Selected label text |

### Returned Helpers

| Helper | Description |
|--------|-------------|
| `isSelected(value)` | Check if value is selected |

## Data Attributes

### Trigger
- `[data-state]` - `'open'` or `'closed'`
- `[data-disabled]` - Present when disabled
- `[data-placeholder]` - Present when showing placeholder
- `[data-melt-select-trigger]` - Present on trigger

### Option
- `[data-selected]` - Option is selected
- `[data-highlighted]` - Option is highlighted
- `[data-disabled]` - Option is disabled
- `[data-melt-select-option]` - Present on options

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Space` / `Enter` | Open menu / Select option |
| `Arrow Down` | Next option |
| `Arrow Up` | Previous option |
| `Home` | First option |
| `End` | Last option |
| `Escape` | Close menu |
| `A-Z` | Type to search |

## Accessibility

Follows [WAI-ARIA Listbox Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/listbox/):
- Proper listbox roles
- Full keyboard navigation
- Screen reader support

## Example

```svelte
<script lang="ts">
  import { createSelect, melt } from '@melt-ui/svelte';
  import { ChevronDown, Check } from 'lucide-svelte';

  const {
    elements: { trigger, menu, option, label },
    states: { open, selected, selectedLabel },
    helpers: { isSelected }
  } = createSelect<string>({
    forceVisible: true,
    positioning: { placement: 'bottom', sameWidth: true }
  });

  const roles = [
    { value: 'employee', label: 'Employee' },
    { value: 'admin', label: 'Admin' }
  ];
</script>

<div class="flex flex-col gap-1">
  <label use:melt={$label} class="text-sm font-semibold text-ocean-800">
    Role
  </label>

  <button
    use:melt={$trigger}
    class="flex items-center justify-between px-4 py-2.5 rounded-lg border-2 border-ocean-500/40 bg-white text-ocean-900
      focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500
      data-[placeholder]:text-ocean-500/50"
  >
    {$selectedLabel || 'Select a role'}
    <ChevronDown class="h-5 w-5 text-ocean-400 {$open ? 'rotate-180' : ''} transition-transform" />
  </button>

  {#if $open}
    <div
      use:melt={$menu}
      class="z-50 max-h-60 overflow-auto rounded-xl bg-white p-1 shadow-lg border border-ocean-200"
    >
      {#each roles as role}
        <div
          use:melt={$option({ value: role.value, label: role.label })}
          class="relative flex items-center justify-between rounded-lg px-4 py-2 cursor-pointer
            text-ocean-800 outline-none
            data-[highlighted]:bg-ocean-100
            data-[selected]:bg-ocean-500 data-[selected]:text-white
            data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed"
        >
          <span>{role.label}</span>
          {#if $isSelected(role.value)}
            <Check class="h-4 w-4" />
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>
```

## Styling with Tailwind

```css
/* Trigger states */
[data-melt-select-trigger] {
  @apply flex items-center justify-between px-4 py-2.5 rounded-lg border-2;
}

[data-melt-select-trigger][data-state="open"] {
  @apply ring-2 ring-ocean-500/30 border-ocean-500;
}

[data-melt-select-trigger][data-placeholder] {
  @apply text-ocean-500/50;
}

/* Option states */
[data-melt-select-option] {
  @apply flex items-center justify-between px-4 py-2 rounded-lg cursor-pointer outline-none;
}

[data-melt-select-option][data-highlighted] {
  @apply bg-ocean-100;
}

[data-melt-select-option][data-selected] {
  @apply bg-ocean-500 text-white;
}

[data-melt-select-option][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}
```

## VacayTracker Role Select

```svelte
<script lang="ts">
  import { createSelect, melt } from '@melt-ui/svelte';
  import { ChevronDown, Check, Shield, User } from 'lucide-svelte';

  let { value = $bindable('employee') }: { value: string } = $props();

  const {
    elements: { trigger, menu, option, label },
    states: { open, selectedLabel },
    helpers: { isSelected }
  } = createSelect<string>({
    defaultSelected: { value, label: value === 'admin' ? 'Admin' : 'Employee' },
    forceVisible: true,
    positioning: { placement: 'bottom', sameWidth: true },
    onSelectedChange: ({ next }) => {
      if (next) {
        value = next.value;
      }
      return next;
    }
  });

  const roles = [
    { value: 'employee', label: 'Employee', icon: User, description: 'Standard user access' },
    { value: 'admin', label: 'Admin', icon: Shield, description: 'Full administrative access' }
  ];
</script>

<div class="flex flex-col gap-1">
  <label use:melt={$label} class="text-sm font-semibold text-ocean-800">
    Role
  </label>

  <button
    use:melt={$trigger}
    class="flex items-center justify-between px-4 py-2.5 rounded-lg border-2 border-ocean-500/40 bg-white text-ocean-900 cursor-pointer
      focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500"
  >
    <span>{$selectedLabel || 'Select role'}</span>
    <ChevronDown class="h-5 w-5 text-ocean-400 {$open ? 'rotate-180' : ''} transition-transform" />
  </button>

  {#if $open}
    <div
      use:melt={$menu}
      class="z-50 rounded-xl bg-white p-1 shadow-lg border border-ocean-200"
    >
      {#each roles as role}
        <div
          use:melt={$option({ value: role.value, label: role.label })}
          class="flex items-center gap-3 rounded-lg px-3 py-2.5 cursor-pointer outline-none
            data-[highlighted]:bg-ocean-100
            data-[selected]:bg-ocean-500 data-[selected]:text-white"
        >
          <role.icon class="h-5 w-5" />
          <div class="flex-1">
            <div class="font-medium">{role.label}</div>
            <div class="text-sm opacity-70">{role.description}</div>
          </div>
          {#if $isSelected(role.value)}
            <Check class="h-4 w-4" />
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>
```

## Multiple Selection

```svelte
<script lang="ts">
  const {
    elements: { trigger, menu, option },
    states: { selected }
  } = createSelect<string>({
    multiple: true
  });
</script>

<button use:melt={$trigger}>
  {$selected.length > 0
    ? `${$selected.length} selected`
    : 'Select options'}
</button>
```

## With Groups

```svelte
<script lang="ts">
  const {
    elements: { trigger, menu, option, group, groupLabel }
  } = createSelect();

  const departments = [
    {
      name: 'Engineering',
      users: [
        { value: 'alice', label: 'Alice' },
        { value: 'bob', label: 'Bob' }
      ]
    },
    {
      name: 'Design',
      users: [
        { value: 'carol', label: 'Carol' },
        { value: 'dave', label: 'Dave' }
      ]
    }
  ];
</script>

<div use:melt={$menu}>
  {#each departments as dept}
    <div use:melt={$group(dept.name)}>
      <div use:melt={$groupLabel(dept.name)} class="px-3 py-2 text-xs font-semibold text-ocean-500">
        {dept.name}
      </div>
      {#each dept.users as user}
        <div use:melt={$option({ value: user.value, label: user.label })}>
          {user.label}
        </div>
      {/each}
    </div>
  {/each}
</div>
```

## Form Integration

```svelte
<script lang="ts">
  const {
    elements: { trigger, menu, option, hiddenInput }
  } = createSelect({
    name: 'role',
    required: true
  });
</script>

<form onsubmit={handleSubmit}>
  <button use:melt={$trigger}>Select role</button>
  <!-- Menu content -->
  <input use:melt={$hiddenInput} />
  <button type="submit">Submit</button>
</form>
```
