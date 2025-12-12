# Label

An accessible label element that properly associates with form controls.

## Use Cases

- Form field labels
- Checkbox/radio labels
- Custom input labels
- Screen reader accessible text

## Installation

```typescript
import { createLabel } from '@melt-ui/svelte';
```

## API Reference

### createLabel Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `for` | `string` | — | ID of associated element |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Label element |

## Data Attributes

### Root
- `[data-melt-label]` - Present on label element

## Accessibility

- Automatically associates with form controls via `for` attribute
- Click on label focuses associated input
- Screen readers announce label with form control

## Example

```svelte
<script lang="ts">
  import { createLabel, melt } from '@melt-ui/svelte';

  const { elements: { root } } = createLabel();
</script>

<div class="flex flex-col gap-1">
  <label use:melt={$root} class="text-sm font-semibold text-ocean-800">
    Email Address
  </label>
  <input
    type="email"
    class="px-4 py-2 rounded-lg border-2 border-ocean-500/40"
    placeholder="john@example.com"
  />
</div>
```

## With For Attribute

```svelte
<script lang="ts">
  import { createLabel, melt } from '@melt-ui/svelte';

  const emailLabel = createLabel();
</script>

<label
  use:melt={emailLabel.elements.root}
  for="email-input"
  class="text-sm font-semibold text-ocean-800"
>
  Email
</label>
<input
  id="email-input"
  type="email"
  class="mt-1 px-4 py-2 rounded-lg border-2 border-ocean-500/40"
/>
```

## Styling with Tailwind

```css
/* Base label styling */
[data-melt-label] {
  @apply text-sm font-semibold text-ocean-800;
}

/* Required indicator */
[data-melt-label][data-required]::after {
  content: ' *';
  @apply text-red-500;
}
```

## With Required Indicator

```svelte
<script lang="ts">
  interface Props {
    label: string;
    required?: boolean;
  }

  let { label, required = false }: Props = $props();

  const { elements: { root } } = createLabel();
</script>

<label use:melt={$root} class="text-sm font-semibold text-ocean-800">
  {label}
  {#if required}
    <span class="text-red-500">*</span>
  {/if}
</label>
```

## With Form Field

```svelte
<script lang="ts">
  import { createLabel, melt } from '@melt-ui/svelte';

  interface Props {
    id: string;
    label: string;
    error?: string;
  }

  let { id, label, error }: Props = $props();

  const labelBuilder = createLabel();
</script>

<div class="flex flex-col gap-1">
  <label
    use:melt={labelBuilder.elements.root}
    for={id}
    class="text-sm font-semibold text-ocean-800"
  >
    {label}
  </label>

  <input
    {id}
    class="px-4 py-2 rounded-lg border-2 {error ? 'border-red-500' : 'border-ocean-500/40'}"
  />

  {#if error}
    <span class="text-sm text-red-500">{error}</span>
  {/if}
</div>
```

## Clickable Label Area

```svelte
<script lang="ts">
  import { createLabel, createCheckbox, melt } from '@melt-ui/svelte';
  import { Check } from 'lucide-svelte';

  const label = createLabel();
  const checkbox = createCheckbox();
</script>

<label
  use:melt={label.elements.root}
  class="flex items-center gap-3 cursor-pointer p-3 rounded-lg hover:bg-ocean-50 transition-colors"
>
  <button
    use:melt={checkbox.elements.root}
    class="h-5 w-5 rounded border-2 border-ocean-400 flex items-center justify-center
      data-[state=checked]:bg-ocean-500 data-[state=checked]:border-ocean-500"
  >
    {#if $checkbox.states.checked}
      <Check class="h-3.5 w-3.5 text-white" />
    {/if}
  </button>
  <span class="text-ocean-700">I agree to the terms</span>
</label>
```
