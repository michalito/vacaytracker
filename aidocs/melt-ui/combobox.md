# Combobox

A searchable select input with an associated popup containing filterable options. Commonly used for autocomplete functionality.

## Use Cases

- User search and selection
- Location/city pickers
- Tag selection with search
- Command palettes

## Installation

```typescript
import { createCombobox } from '@melt-ui/svelte';
```

## API Reference

### createCombobox Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `defaultSelected` | `ComboboxOption` | — | Initial selected item |
| `selected` | `Writable<ComboboxOption>` | — | Controlled selection |
| `defaultOpen` | `boolean` | `false` | Initial open state |
| `open` | `Writable<boolean>` | — | Controlled open state |
| `multiple` | `boolean` | `false` | Enable multi-select |
| `scrollAlignment` | `'nearest' \| 'center'` | `'nearest'` | Scroll behavior |
| `loop` | `boolean` | `false` | Loop keyboard navigation |
| `closeOnOutsideClick` | `boolean` | `true` | Close on click outside |
| `portal` | `string \| HTMLElement \| null` | `'body'` | Portal target |
| `forceVisible` | `boolean` | `false` | Force visibility |
| `onSelectedChange` | `ChangeFn` | — | Selection callback |
| `onOpenChange` | `ChangeFn` | — | Open state callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `menu` | Options popup container |
| `input` | Search input element |
| `option` | Individual option element |
| `label` | Input label |
| `group` | Option group container |
| `groupLabel` | Group heading |
| `hiddenInput` | Hidden input for forms |
| `arrow` | Optional arrow element |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `open` | `Writable<boolean>` | Menu visibility |
| `inputValue` | `Writable<string>` | Current input text |
| `touchedInput` | `Readable<boolean>` | Input has been typed in |
| `selected` | `Writable<ComboboxOption>` | Selected item(s) |
| `highlighted` | `Readable<ComboboxOption>` | Highlighted option |

### Returned Helpers

| Helper | Description |
|--------|-------------|
| `isSelected` | Check if option is selected |
| `isHighlighted` | Check if option is highlighted |

## Data Attributes

### Input
- `[data-state]` - `'open'` or `'closed'`
- `[data-disabled]` - Present when disabled
- `[data-melt-combobox-input]` - Present on input

### Menu
- `[data-side]` - Placement side
- `[data-align]` - Placement alignment
- `[data-melt-combobox-menu]` - Present on menu

### Option
- `[data-disabled]` - Option is disabled
- `[data-selected]` - Option is selected
- `[data-highlighted]` - Option is highlighted
- `[data-melt-combobox-option]` - Present on options

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Arrow Down` | Highlight next option |
| `Arrow Up` | Highlight previous option |
| `Enter` | Select highlighted option |
| `Escape` | Close menu |
| `Home` | Highlight first option |
| `End` | Highlight last option |

## Accessibility

Follows the [WAI-ARIA Combobox Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/combobox/):
- Proper ARIA roles and attributes
- Screen reader announcements
- Keyboard navigation support

## Example

```svelte
<script lang="ts">
  import { createCombobox, melt } from '@melt-ui/svelte';
  import { Check, ChevronDown } from 'lucide-svelte';

  const users = [
    { value: '1', label: 'John Doe' },
    { value: '2', label: 'Jane Smith' },
    { value: '3', label: 'Bob Wilson' },
    { value: '4', label: 'Alice Brown' }
  ];

  const {
    elements: { menu, input, option, label },
    states: { open, inputValue, selected },
    helpers: { isSelected, isHighlighted }
  } = createCombobox({ forceVisible: true });

  const filteredUsers = $derived(
    $inputValue
      ? users.filter((user) =>
          user.label.toLowerCase().includes($inputValue.toLowerCase())
        )
      : users
  );
</script>

<div class="flex flex-col gap-1">
  <label use:melt={$label} class="text-sm font-semibold text-ocean-800">
    Select User
  </label>

  <div class="relative">
    <input
      use:melt={$input}
      class="w-full px-4 py-2.5 pr-10 rounded-lg border-2 border-ocean-500/40 bg-white text-ocean-900 placeholder-ocean-500/50 focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500"
      placeholder="Search users..."
    />
    <ChevronDown
      class="absolute right-3 top-1/2 -translate-y-1/2 h-5 w-5 text-ocean-400 {$open ? 'rotate-180' : ''} transition-transform"
    />
  </div>

  {#if $open}
    <ul
      use:melt={$menu}
      class="z-50 max-h-60 overflow-auto rounded-xl bg-white p-1 shadow-lg border border-ocean-200"
    >
      {#if filteredUsers.length === 0}
        <li class="px-4 py-2 text-sm text-ocean-500">No results found</li>
      {:else}
        {#each filteredUsers as user}
          <li
            use:melt={$option({ value: user.value, label: user.label })}
            class="relative flex items-center justify-between rounded-lg px-4 py-2 text-ocean-800 cursor-pointer
              data-[highlighted]:bg-ocean-100
              data-[selected]:bg-ocean-500 data-[selected]:text-white
              data-[disabled]:opacity-50"
          >
            <span>{user.label}</span>
            {#if $isSelected(user)}
              <Check class="h-4 w-4" />
            {/if}
          </li>
        {/each}
      {/if}
    </ul>
  {/if}
</div>
```

## Styling with Tailwind

```css
/* Input focus state */
[data-melt-combobox-input]:focus {
  @apply ring-2 ring-ocean-500/30 border-ocean-500;
}

/* Menu positioning */
[data-melt-combobox-menu] {
  @apply absolute z-50 mt-1;
}

/* Highlighted option */
[data-melt-combobox-option][data-highlighted] {
  @apply bg-ocean-100;
}

/* Selected option */
[data-melt-combobox-option][data-selected] {
  @apply bg-ocean-500 text-white;
}

/* Disabled option */
[data-melt-combobox-option][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}
```

## Multiple Selection

```svelte
<script lang="ts">
  import { createCombobox, melt } from '@melt-ui/svelte';

  const {
    elements: { menu, input, option },
    states: { selected }
  } = createCombobox<string>({
    multiple: true
  });
</script>

<!-- Display selected items as chips -->
<div class="flex flex-wrap gap-2 mb-2">
  {#each $selected as item}
    <span class="inline-flex items-center gap-1 px-2 py-1 bg-ocean-100 text-ocean-700 rounded-full text-sm">
      {item.label}
      <button onclick={() => /* remove item */}>×</button>
    </span>
  {/each}
</div>

<input use:melt={$input} placeholder="Add more..." />
```

## With Async Search

```svelte
<script lang="ts">
  import { createCombobox, melt } from '@melt-ui/svelte';

  let users = $state<User[]>([]);
  let loading = $state(false);

  const {
    elements: { menu, input, option },
    states: { inputValue }
  } = createCombobox();

  // Debounced search
  let timeout: ReturnType<typeof setTimeout>;
  $effect(() => {
    clearTimeout(timeout);
    timeout = setTimeout(async () => {
      if ($inputValue.length >= 2) {
        loading = true;
        users = await searchUsers($inputValue);
        loading = false;
      }
    }, 300);
  });
</script>

{#if loading}
  <li class="px-4 py-2 text-sm text-ocean-500">Searching...</li>
{:else}
  {#each users as user}
    <li use:melt={$option({ value: user.id, label: user.name })}>
      {user.name}
    </li>
  {/each}
{/if}
```

## Grouped Options

```svelte
<script lang="ts">
  const departments = [
    {
      name: 'Engineering',
      users: [
        { value: '1', label: 'Alice' },
        { value: '2', label: 'Bob' }
      ]
    },
    {
      name: 'Design',
      users: [
        { value: '3', label: 'Carol' },
        { value: '4', label: 'Dave' }
      ]
    }
  ];
</script>

<ul use:melt={$menu}>
  {#each departments as dept}
    <div use:melt={$group(dept.name)}>
      <div use:melt={$groupLabel(dept.name)} class="px-4 py-2 text-xs font-semibold text-ocean-500 uppercase">
        {dept.name}
      </div>
      {#each dept.users as user}
        <li use:melt={$option({ value: user.value, label: user.label })}>
          {user.label}
        </li>
      {/each}
    </div>
  {/each}
</ul>
```
