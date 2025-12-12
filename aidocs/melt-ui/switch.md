# Switch

A toggle control for switching between on and off states.

## Use Cases

- Settings toggles
- Feature flags
- Dark mode toggle
- Enable/disable options

## Installation

```typescript
import { createSwitch } from '@melt-ui/svelte';
```

## API Reference

### createSwitch Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `disabled` | `boolean` | `false` | Disable the switch |
| `required` | `boolean` | `false` | Mark as required |
| `name` | `string` | — | Form input name |
| `value` | `string` | — | Form input value |
| `defaultChecked` | `boolean` | `false` | Initial checked state |
| `checked` | `Writable<boolean>` | — | Controlled checked state |
| `onCheckedChange` | `ChangeFn` | — | Checked change callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Switch container/track |
| `input` | Hidden form input |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `checked` | `Writable<boolean>` | Current checked state |

## Data Attributes

### Root
- `[data-state]` - `'checked'` or `'unchecked'`
- `[data-disabled]` - Present when disabled
- `[data-melt-switch]` - Present on root

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Space` | Toggle switch |
| `Enter` | Toggle switch |

## Accessibility

Follows [WAI-ARIA Switch Pattern](https://www.w3.org/WAI/ARIA/apg/patterns/switch/):
- `role="switch"`
- `aria-checked` state
- Keyboard accessible

## Example

```svelte
<script lang="ts">
  import { createSwitch, melt } from '@melt-ui/svelte';

  const {
    elements: { root, input },
    states: { checked }
  } = createSwitch({
    defaultChecked: false
  });
</script>

<label class="flex items-center gap-3 cursor-pointer">
  <button
    use:melt={$root}
    class="relative h-6 w-11 rounded-full transition-colors
      data-[state=unchecked]:bg-ocean-300
      data-[state=checked]:bg-ocean-500
      data-[disabled]:opacity-50"
  >
    <span
      class="absolute left-0.5 top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform
        {$checked ? 'translate-x-5' : 'translate-x-0'}"
    ></span>
  </button>
  <input use:melt={$input} />
  <span class="text-ocean-800">Enable notifications</span>
</label>
```

## Styling with Tailwind

```css
/* Switch track */
[data-melt-switch] {
  @apply relative h-6 w-11 rounded-full transition-colors;
}

[data-melt-switch][data-state="unchecked"] {
  @apply bg-ocean-300;
}

[data-melt-switch][data-state="checked"] {
  @apply bg-ocean-500;
}

[data-melt-switch][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}

/* Focus ring */
[data-melt-switch]:focus-visible {
  @apply ring-2 ring-ocean-500 ring-offset-2;
}
```

## VacayTracker Settings Toggle

```svelte
<script lang="ts">
  import { createSwitch, melt } from '@melt-ui/svelte';

  interface Props {
    label: string;
    description?: string;
    checked?: boolean;
    onchange?: (checked: boolean) => void;
  }

  let { label, description, checked = false, onchange }: Props = $props();

  const switchBuilder = createSwitch({
    defaultChecked: checked,
    onCheckedChange: ({ next }) => {
      onchange?.(next);
      return next;
    }
  });

  const { elements: { root, input }, states } = switchBuilder;
</script>

<label class="flex items-center justify-between p-4 rounded-lg hover:bg-ocean-50 cursor-pointer transition-colors">
  <div class="flex-1">
    <span class="font-medium text-ocean-800">{label}</span>
    {#if description}
      <p class="text-sm text-ocean-500 mt-0.5">{description}</p>
    {/if}
  </div>

  <button
    use:melt={$root}
    class="relative h-6 w-11 rounded-full transition-colors shrink-0 ml-4
      data-[state=unchecked]:bg-ocean-300
      data-[state=checked]:bg-ocean-500
      focus:outline-none focus:ring-2 focus:ring-ocean-500/50"
  >
    <span
      class="absolute left-0.5 top-0.5 h-5 w-5 rounded-full bg-white shadow-sm transition-transform
        {$states.checked ? 'translate-x-5' : 'translate-x-0'}"
    ></span>
  </button>
  <input use:melt={$input} />
</label>
```

## Settings Panel

```svelte
<script lang="ts">
  import { createSwitch, melt } from '@melt-ui/svelte';

  const settings = [
    {
      key: 'emailNotifications',
      label: 'Email Notifications',
      description: 'Receive email when your request is approved or rejected',
      defaultChecked: true
    },
    {
      key: 'weeklyDigest',
      label: 'Weekly Digest',
      description: 'Receive a weekly summary of team vacation schedule',
      defaultChecked: false
    },
    {
      key: 'darkMode',
      label: 'Dark Mode',
      description: 'Use dark theme for the application',
      defaultChecked: false
    }
  ];

  const switches = settings.map(s => ({
    ...s,
    switch: createSwitch({ defaultChecked: s.defaultChecked })
  }));
</script>

<div class="divide-y divide-ocean-100">
  {#each switches as { key, label, description, switch: sw }}
    <label class="flex items-center justify-between p-4 cursor-pointer">
      <div>
        <span class="font-medium text-ocean-800">{label}</span>
        <p class="text-sm text-ocean-500">{description}</p>
      </div>
      <button
        use:melt={sw.elements.root}
        class="relative h-6 w-11 rounded-full transition-colors
          data-[state=unchecked]:bg-ocean-300
          data-[state=checked]:bg-ocean-500"
      >
        <span
          class="absolute left-0.5 top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform
            {$sw.states.checked ? 'translate-x-5' : 'translate-x-0'}"
        ></span>
      </button>
    </label>
  {/each}
</div>
```

## With Icon

```svelte
<script lang="ts">
  import { createSwitch, melt } from '@melt-ui/svelte';
  import { Sun, Moon } from 'lucide-svelte';

  const {
    elements: { root },
    states: { checked }
  } = createSwitch();
</script>

<button
  use:melt={$root}
  class="relative h-8 w-14 rounded-full transition-colors
    data-[state=unchecked]:bg-ocean-200
    data-[state=checked]:bg-ocean-800"
>
  <span
    class="absolute left-1 top-1 h-6 w-6 rounded-full bg-white shadow flex items-center justify-center transition-transform
      {$checked ? 'translate-x-6' : 'translate-x-0'}"
  >
    {#if $checked}
      <Moon class="h-4 w-4 text-ocean-800" />
    {:else}
      <Sun class="h-4 w-4 text-orange-500" />
    {/if}
  </span>
</button>
```

## Form Integration

```svelte
<script lang="ts">
  const {
    elements: { root, input }
  } = createSwitch({
    name: 'notifications',
    value: 'enabled'
  });
</script>

<form onsubmit={handleSubmit}>
  <label class="flex items-center gap-3">
    <button use:melt={$root}>...</button>
    <input use:melt={$input} />
    <span>Enable notifications</span>
  </label>
  <button type="submit">Save</button>
</form>
```

## Controlled State

```svelte
<script lang="ts">
  import { writable } from 'svelte/store';

  const isEnabled = writable(false);

  const sw = createSwitch({
    checked: isEnabled,
    onCheckedChange: ({ next }) => {
      console.log('Switch changed to:', next);
      // Could make API call here
      return next;
    }
  });

  // Programmatic control
  function toggle() {
    isEnabled.update(v => !v);
  }
</script>
```
