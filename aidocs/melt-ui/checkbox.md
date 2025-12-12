# Checkbox

A control enabling users to toggle between checked, unchecked, and optional indeterminate states.

## Use Cases

- Form agreements and terms
- Multi-select filters
- Todo list items
- Bulk selection controls

## Installation

```typescript
import { createCheckbox } from '@melt-ui/svelte';
```

## API Reference

### createCheckbox Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `disabled` | `boolean` | `false` | Disable checkbox |
| `required` | `boolean` | `false` | Mark as required |
| `name` | `string` | — | Form input name |
| `value` | `string` | — | Form input value |
| `defaultChecked` | `boolean \| 'indeterminate'` | `false` | Initial state |
| `checked` | `Writable<boolean \| 'indeterminate'>` | — | Controlled state |
| `onCheckedChange` | `ChangeFn` | — | State change callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Checkbox container element |
| `input` | Hidden native input for forms |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `checked` | `Writable<boolean \| 'indeterminate'>` | Current state |

### Returned Helpers

| Helper | Description |
|--------|-------------|
| `isChecked` | Derived boolean for checked state |
| `isIndeterminate` | Derived boolean for indeterminate state |

## Data Attributes

### Root
- `[data-state]` - `'checked'`, `'unchecked'`, or `'indeterminate'`
- `[data-disabled]` - Present when disabled
- `[data-melt-checkbox]` - Present on root element

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Space` | Toggle checkbox state |

## Accessibility

Follows the [WAI-ARIA Checkbox Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/checkbox/):
- Tri-state support (checked, unchecked, indeterminate)
- Proper `aria-checked` attribute
- Associated label support

## Example

```svelte
<script lang="ts">
  import { createCheckbox, melt } from '@melt-ui/svelte';
  import { Check, Minus } from 'lucide-svelte';

  const {
    elements: { root, input },
    states: { checked },
    helpers: { isChecked, isIndeterminate }
  } = createCheckbox();
</script>

<label class="flex items-center gap-3 cursor-pointer">
  <button
    use:melt={$root}
    class="h-5 w-5 rounded border-2 transition-colors flex items-center justify-center
      data-[state=unchecked]:border-ocean-400 data-[state=unchecked]:bg-white
      data-[state=checked]:border-ocean-500 data-[state=checked]:bg-ocean-500
      data-[state=indeterminate]:border-ocean-500 data-[state=indeterminate]:bg-ocean-500
      data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed"
  >
    {#if $isChecked}
      <Check class="h-3.5 w-3.5 text-white" />
    {:else if $isIndeterminate}
      <Minus class="h-3.5 w-3.5 text-white" />
    {/if}
  </button>
  <input use:melt={$input} />
  <span class="text-ocean-800">I agree to the terms and conditions</span>
</label>
```

## Styling with Tailwind

```css
/* Unchecked state */
[data-melt-checkbox][data-state="unchecked"] {
  @apply border-ocean-400 bg-white;
}

/* Checked state */
[data-melt-checkbox][data-state="checked"] {
  @apply border-ocean-500 bg-ocean-500;
}

/* Indeterminate state */
[data-melt-checkbox][data-state="indeterminate"] {
  @apply border-ocean-500 bg-ocean-500;
}

/* Disabled state */
[data-melt-checkbox][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}

/* Focus ring */
[data-melt-checkbox]:focus-visible {
  @apply ring-2 ring-ocean-500 ring-offset-2;
}
```

## Checkbox Group

```svelte
<script lang="ts">
  import { createCheckbox, melt } from '@melt-ui/svelte';

  const options = ['Email', 'SMS', 'Push notifications'];
  let selected = $state<string[]>([]);

  function toggleOption(option: string) {
    if (selected.includes(option)) {
      selected = selected.filter(s => s !== option);
    } else {
      selected = [...selected, option];
    }
  }
</script>

<fieldset class="space-y-3">
  <legend class="text-sm font-semibold text-ocean-800 mb-2">
    Notification Preferences
  </legend>
  {#each options as option}
    {@const checkbox = createCheckbox({
      checked: selected.includes(option)
    })}
    <label class="flex items-center gap-3 cursor-pointer">
      <button
        use:melt={checkbox.elements.root}
        onclick={() => toggleOption(option)}
        class="h-5 w-5 rounded border-2 flex items-center justify-center ..."
      >
        <!-- Check icon -->
      </button>
      <span class="text-ocean-700">{option}</span>
    </label>
  {/each}
</fieldset>
```

## Indeterminate State (Select All)

```svelte
<script lang="ts">
  import { createCheckbox, melt } from '@melt-ui/svelte';

  let items = $state([
    { id: 1, name: 'Item 1', checked: false },
    { id: 2, name: 'Item 2', checked: false },
    { id: 3, name: 'Item 3', checked: false }
  ]);

  const allChecked = $derived(items.every(i => i.checked));
  const someChecked = $derived(items.some(i => i.checked) && !allChecked);

  const selectAll = createCheckbox({
    checked: allChecked ? true : someChecked ? 'indeterminate' : false
  });

  function toggleAll() {
    const newValue = !allChecked;
    items = items.map(i => ({ ...i, checked: newValue }));
  }
</script>

<label class="flex items-center gap-3">
  <button use:melt={selectAll.elements.root} onclick={toggleAll}>
    <!-- Select all checkbox -->
  </button>
  <span>Select All</span>
</label>

{#each items as item}
  <!-- Individual checkboxes -->
{/each}
```

## With Form Submission

```svelte
<form onsubmit={handleSubmit}>
  <label class="flex items-center gap-3">
    <button use:melt={$root}>
      <!-- Checkbox visual -->
    </button>
    <input use:melt={$input} name="terms" value="accepted" />
    <span>Accept terms</span>
  </label>
  <button type="submit">Submit</button>
</form>
```
