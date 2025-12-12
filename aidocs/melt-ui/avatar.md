# Avatar

A visual representation of a user with image and fallback support.

## Use Cases

- User profile pictures
- Team member displays
- Comment author avatars
- Navigation user menus

## Installation

```typescript
import { createAvatar } from '@melt-ui/svelte';
```

## API Reference

### createAvatar Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `src` | `string` | — | Image source URL |
| `delayMs` | `number` | `0` | Delay before showing fallback |
| `loadingStatus` | `Writable<'loading' \| 'loaded' \| 'error'>` | — | Controlled loading state |
| `onLoadingStatusChange` | `ChangeFn` | — | Loading state callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `image` | The avatar image element |
| `fallback` | Fallback content when image fails |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `loadingStatus` | `Readable<'loading' \| 'loaded' \| 'error'>` | Current loading state |

## Data Attributes

### Image
- `[data-state]` - `'loading'`, `'loaded'`, or `'error'`
- `[data-melt-avatar-image]` - Present on image element

### Fallback
- `[data-state]` - `'loading'`, `'loaded'`, or `'error'`
- `[data-melt-avatar-fallback]` - Present on fallback element

## Accessibility

- Uses semantic `<img>` element with proper `alt` text
- Fallback provides visual alternative when image unavailable
- Screen readers announce user name via alt text or fallback content

## Example

```svelte
<script lang="ts">
  import { createAvatar, melt } from '@melt-ui/svelte';

  interface Props {
    src?: string;
    name: string;
  }

  let { src, name }: Props = $props();

  const {
    elements: { image, fallback },
    states: { loadingStatus }
  } = createAvatar({ src });

  // Get initials for fallback
  const initials = $derived(
    name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  );
</script>

<div class="relative h-10 w-10 rounded-full bg-ocean-200 overflow-hidden">
  <img
    use:melt={$image}
    alt={name}
    class="h-full w-full object-cover"
  />
  <span
    use:melt={$fallback}
    class="absolute inset-0 flex items-center justify-center text-sm font-medium text-ocean-700"
  >
    {initials}
  </span>
</div>
```

## Styling with Tailwind

```svelte
<!-- Avatar with status indicator -->
<div class="relative">
  <div class="h-12 w-12 rounded-full bg-ocean-200 overflow-hidden ring-2 ring-white">
    <img
      use:melt={$image}
      alt={name}
      class="h-full w-full object-cover"
    />
    <span
      use:melt={$fallback}
      class="absolute inset-0 flex items-center justify-center bg-gradient-to-br from-ocean-400 to-ocean-600 text-white font-semibold"
    >
      {initials}
    </span>
  </div>
  <!-- Online status -->
  <span class="absolute bottom-0 right-0 h-3 w-3 rounded-full bg-green-500 ring-2 ring-white"></span>
</div>
```

## Avatar Sizes

```svelte
<script lang="ts">
  const sizes = {
    sm: 'h-8 w-8 text-xs',
    md: 'h-10 w-10 text-sm',
    lg: 'h-12 w-12 text-base',
    xl: 'h-16 w-16 text-lg'
  };

  let { size = 'md' }: { size?: keyof typeof sizes } = $props();
</script>

<div class="relative rounded-full bg-ocean-200 overflow-hidden {sizes[size]}">
  <!-- Avatar content -->
</div>
```

## Avatar Group

```svelte
<div class="flex -space-x-2">
  {#each users as user}
    <div class="relative h-10 w-10 rounded-full ring-2 ring-white">
      <!-- Individual avatars -->
    </div>
  {/each}
  {#if remainingCount > 0}
    <div class="flex h-10 w-10 items-center justify-center rounded-full bg-ocean-100 ring-2 ring-white text-sm font-medium text-ocean-700">
      +{remainingCount}
    </div>
  {/if}
</div>
```

## Loading State Handling

```svelte
<script lang="ts">
  const { states: { loadingStatus } } = createAvatar({ src });
</script>

{#if $loadingStatus === 'loading'}
  <div class="animate-pulse bg-ocean-200 h-full w-full"></div>
{:else if $loadingStatus === 'error'}
  <span use:melt={$fallback}>{initials}</span>
{:else}
  <img use:melt={$image} alt={name} />
{/if}
```
