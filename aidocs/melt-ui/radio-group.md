# Radio Group

A set of checkable buttons where only one button can be selected at a time.

## Use Cases

- Settings options
- Survey questions
- Payment method selection
- Plan/tier selection

## Installation

```typescript
import { createRadioGroup } from '@melt-ui/svelte';
```

## API Reference

### createRadioGroup Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `disabled` | `boolean` | `false` | Disable the entire group |
| `required` | `boolean` | `false` | Mark as required |
| `loop` | `boolean` | `false` | Loop keyboard navigation |
| `orientation` | `'horizontal' \| 'vertical'` | `'vertical'` | Layout direction |
| `name` | `string` | — | Hidden input name |
| `defaultValue` | `string` | — | Initial selected value |
| `value` | `Writable<string>` | — | Controlled value |
| `onValueChange` | `ChangeFn` | — | Value change callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Radio group container |
| `item` | Individual radio button |
| `hiddenInput` | Hidden form input |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `value` | `Writable<string>` | Currently selected value |

### Returned Helpers

| Helper | Description |
|--------|-------------|
| `isChecked(value)` | Check if a value is selected |

## Data Attributes

### Root
- `[data-orientation]` - `'vertical'` or `'horizontal'`
- `[data-melt-radio-group]` - Present on root

### Item
- `[data-state]` - `'checked'` or `'unchecked'`
- `[data-disabled]` - Present when disabled
- `[data-orientation]` - Inherited orientation
- `[data-melt-radio-group-item]` - Present on items

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Tab` | Focus checked item or first item |
| `Space` | Select focused item |
| `Arrow Down/Right` | Select next item |
| `Arrow Up/Left` | Select previous item |

## Accessibility

Follows [WAI-ARIA Radio Group Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/radio/):
- Proper `role="radiogroup"` and `role="radio"`
- `aria-checked` state management
- Keyboard navigation support

## Example

```svelte
<script lang="ts">
  import { createRadioGroup, melt } from '@melt-ui/svelte';

  const {
    elements: { root, item },
    states: { value },
    helpers: { isChecked }
  } = createRadioGroup({
    defaultValue: 'comfortable'
  });

  const options = [
    { value: 'default', label: 'Default', description: 'System default spacing' },
    { value: 'comfortable', label: 'Comfortable', description: 'More space between items' },
    { value: 'compact', label: 'Compact', description: 'Tighter spacing for more content' }
  ];
</script>

<fieldset>
  <legend class="text-sm font-semibold text-ocean-800 mb-3">Display Density</legend>

  <div use:melt={$root} class="space-y-2">
    {#each options as option}
      <label
        class="flex items-start gap-3 p-3 rounded-lg border-2 cursor-pointer transition-colors
          {$isChecked(option.value)
            ? 'border-ocean-500 bg-ocean-50'
            : 'border-ocean-200 hover:border-ocean-300'}"
      >
        <button
          use:melt={$item(option.value)}
          class="mt-0.5 h-5 w-5 rounded-full border-2 flex items-center justify-center transition-colors
            data-[state=checked]:border-ocean-500 data-[state=checked]:bg-ocean-500
            data-[state=unchecked]:border-ocean-400"
        >
          {#if $isChecked(option.value)}
            <div class="h-2 w-2 rounded-full bg-white"></div>
          {/if}
        </button>
        <div>
          <span class="font-medium text-ocean-800">{option.label}</span>
          <p class="text-sm text-ocean-500">{option.description}</p>
        </div>
      </label>
    {/each}
  </div>
</fieldset>
```

## Styling with Tailwind

```css
/* Radio button states */
[data-melt-radio-group-item] {
  @apply h-5 w-5 rounded-full border-2 border-ocean-400 flex items-center justify-center;
}

[data-melt-radio-group-item][data-state="checked"] {
  @apply border-ocean-500 bg-ocean-500;
}

[data-melt-radio-group-item][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}

/* Focus ring */
[data-melt-radio-group-item]:focus-visible {
  @apply ring-2 ring-ocean-500 ring-offset-2;
}
```

## VacayTracker Role Selection

```svelte
<script lang="ts">
  import { createRadioGroup, melt } from '@melt-ui/svelte';
  import { Shield, User } from 'lucide-svelte';

  const {
    elements: { root, item },
    states: { value },
    helpers: { isChecked }
  } = createRadioGroup({
    defaultValue: 'employee'
  });

  const roles = [
    {
      value: 'employee',
      label: 'Employee',
      description: 'Submit vacation requests and view team calendar',
      icon: User
    },
    {
      value: 'admin',
      label: 'Admin',
      description: 'Manage users, approve requests, and configure settings',
      icon: Shield
    }
  ];
</script>

<div use:melt={$root} class="space-y-3">
  {#each roles as role}
    <label
      class="flex items-center gap-4 p-4 rounded-xl border-2 cursor-pointer transition-all
        {$isChecked(role.value)
          ? 'border-ocean-500 bg-ocean-500/10 shadow-sm'
          : 'border-ocean-200 hover:border-ocean-300'}"
    >
      <button
        use:melt={$item(role.value)}
        class="h-5 w-5 rounded-full border-2 flex items-center justify-center shrink-0
          data-[state=checked]:border-ocean-500 data-[state=checked]:bg-ocean-500
          data-[state=unchecked]:border-ocean-400"
      >
        {#if $isChecked(role.value)}
          <div class="h-2 w-2 rounded-full bg-white"></div>
        {/if}
      </button>

      <role.icon class="h-5 w-5 text-ocean-500 shrink-0" />

      <div class="flex-1">
        <span class="font-medium text-ocean-800">{role.label}</span>
        <p class="text-sm text-ocean-500">{role.description}</p>
      </div>
    </label>
  {/each}
</div>
```

## Horizontal Layout

```svelte
<script lang="ts">
  const {
    elements: { root, item },
    helpers: { isChecked }
  } = createRadioGroup({
    orientation: 'horizontal',
    defaultValue: 'monthly'
  });

  const options = ['Weekly', 'Monthly', 'Yearly'];
</script>

<div use:melt={$root} class="flex gap-2 p-1 bg-ocean-100 rounded-lg">
  {#each options as option}
    <button
      use:melt={$item(option.toLowerCase())}
      class="flex-1 px-4 py-2 rounded-md text-sm font-medium transition-colors
        data-[state=checked]:bg-white data-[state=checked]:text-ocean-800 data-[state=checked]:shadow-sm
        data-[state=unchecked]:text-ocean-600"
    >
      {option}
    </button>
  {/each}
</div>
```

## With Disabled Options

```svelte
<script lang="ts">
  const options = [
    { value: 'free', label: 'Free', disabled: false },
    { value: 'pro', label: 'Pro', disabled: false },
    { value: 'enterprise', label: 'Enterprise', disabled: true }
  ];
</script>

<div use:melt={$root}>
  {#each options as option}
    <button
      use:melt={$item(option.value)}
      disabled={option.disabled}
      class="..."
    >
      {option.label}
      {#if option.disabled}
        <span class="text-xs text-ocean-400">(Coming soon)</span>
      {/if}
    </button>
  {/each}
</div>
```

## Form Integration

```svelte
<script lang="ts">
  const {
    elements: { root, item, hiddenInput }
  } = createRadioGroup({
    name: 'role',
    defaultValue: 'employee'
  });
</script>

<form onsubmit={handleSubmit}>
  <div use:melt={$root}>
    {#each roles as role}
      <button use:melt={$item(role.value)}>{role.label}</button>
    {/each}
  </div>
  <input use:melt={$hiddenInput} />
  <button type="submit">Submit</button>
</form>
```
