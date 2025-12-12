# Link Preview

A hover card that displays a preview of link content when hovering over a link.

## Use Cases

- Social media link previews
- Document link previews
- User profile hover cards
- Article snippet previews

## Installation

```typescript
import { createLinkPreview } from '@melt-ui/svelte';
```

## API Reference

### createLinkPreview Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `positioning` | `FloatingConfig` | — | Popup positioning |
| `arrowSize` | `number` | `8` | Arrow size in pixels |
| `defaultOpen` | `boolean` | `false` | Initial open state |
| `open` | `Writable<boolean>` | — | Controlled open state |
| `openDelay` | `number` | `700` | Delay before opening (ms) |
| `closeDelay` | `number` | `300` | Delay before closing (ms) |
| `portal` | `string \| HTMLElement \| null` | `'body'` | Portal target |
| `forceVisible` | `boolean` | `false` | Force visibility |
| `closeOnOutsideClick` | `boolean` | `true` | Close on click outside |
| `onOpenChange` | `ChangeFn` | — | Open state callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `trigger` | Link element that triggers preview |
| `content` | Preview popup content |
| `arrow` | Optional arrow element |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `open` | `Writable<boolean>` | Preview visibility |

## Data Attributes

### Trigger
- `[data-state]` - `'open'` or `'closed'`
- `[data-melt-link-preview-trigger]` - Present on trigger

### Content
- `[data-state]` - `'open'` or `'closed'`
- `[data-side]` - Placement side
- `[data-align]` - Placement alignment
- `[data-melt-link-preview-content]` - Present on content

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Tab` | Opens preview, moves focus |
| `Escape` | Closes preview |

## Accessibility

- Preview is accessible via keyboard focus
- Content is announced to screen readers
- Proper ARIA relationships

## Example

```svelte
<script lang="ts">
  import { createLinkPreview, melt } from '@melt-ui/svelte';

  const {
    elements: { trigger, content, arrow },
    states: { open }
  } = createLinkPreview({
    forceVisible: true,
    positioning: { placement: 'bottom' }
  });
</script>

<a
  use:melt={$trigger}
  href="https://example.com"
  class="text-ocean-600 underline hover:text-ocean-800"
>
  Visit Example.com
</a>

{#if $open}
  <div
    use:melt={$content}
    class="z-50 w-80 rounded-xl bg-white shadow-xl border border-ocean-200 overflow-hidden"
  >
    <div
      use:melt={$arrow}
      class="absolute bg-white border-l border-t border-ocean-200 h-3 w-3 rotate-45"
    ></div>

    <img
      src="/preview-image.jpg"
      alt="Preview"
      class="w-full h-40 object-cover"
    />
    <div class="p-4">
      <h3 class="font-semibold text-ocean-800">Example Website</h3>
      <p class="text-sm text-ocean-600 mt-1">
        A brief description of what this website is about and what visitors can expect to find.
      </p>
    </div>
  </div>
{/if}
```

## Styling with Tailwind

```css
/* Content animation */
[data-melt-link-preview-content] {
  @apply transition-all duration-200;
}

[data-melt-link-preview-content][data-state="open"] {
  @apply opacity-100 scale-100;
}

[data-melt-link-preview-content][data-state="closed"] {
  @apply opacity-0 scale-95;
}

/* Arrow positioning */
[data-melt-link-preview-content][data-side="top"] [data-arrow] {
  @apply bottom-[-6px] rotate-[225deg];
}

[data-melt-link-preview-content][data-side="bottom"] [data-arrow] {
  @apply top-[-6px] rotate-45;
}
```

## User Profile Preview

```svelte
<script lang="ts">
  import { createLinkPreview, melt } from '@melt-ui/svelte';

  interface UserPreviewProps {
    user: {
      id: string;
      name: string;
      email: string;
      avatar?: string;
      role: string;
    };
  }

  let { user }: UserPreviewProps = $props();

  const {
    elements: { trigger, content },
    states: { open }
  } = createLinkPreview({
    forceVisible: true,
    openDelay: 500
  });

  // Get initials for avatar fallback
  const initials = user.name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);
</script>

<a
  use:melt={$trigger}
  href="/users/{user.id}"
  class="text-ocean-600 hover:underline"
>
  {user.name}
</a>

{#if $open}
  <div
    use:melt={$content}
    class="z-50 w-64 rounded-xl bg-white shadow-xl border border-ocean-200 p-4"
  >
    <div class="flex items-center gap-3">
      <div class="h-12 w-12 rounded-full bg-ocean-200 flex items-center justify-center overflow-hidden">
        {#if user.avatar}
          <img src={user.avatar} alt={user.name} class="h-full w-full object-cover" />
        {:else}
          <span class="text-ocean-600 font-semibold">{initials}</span>
        {/if}
      </div>
      <div>
        <h3 class="font-semibold text-ocean-800">{user.name}</h3>
        <p class="text-sm text-ocean-500 capitalize">{user.role}</p>
      </div>
    </div>
    <p class="text-sm text-ocean-600 mt-3">{user.email}</p>
  </div>
{/if}
```

## With Loading State

```svelte
<script lang="ts">
  import { createLinkPreview, melt } from '@melt-ui/svelte';

  let previewData = $state<PreviewData | null>(null);
  let loading = $state(false);

  const {
    elements: { trigger, content },
    states: { open }
  } = createLinkPreview({
    onOpenChange: async ({ next }) => {
      if (next && !previewData) {
        loading = true;
        previewData = await fetchPreviewData(url);
        loading = false;
      }
      return next;
    }
  });
</script>

{#if $open}
  <div use:melt={$content} class="w-80 rounded-xl bg-white shadow-xl p-4">
    {#if loading}
      <div class="animate-pulse space-y-2">
        <div class="h-32 bg-ocean-200 rounded"></div>
        <div class="h-4 bg-ocean-200 rounded w-3/4"></div>
        <div class="h-4 bg-ocean-200 rounded w-1/2"></div>
      </div>
    {:else if previewData}
      <img src={previewData.image} alt="" class="w-full h-32 object-cover rounded" />
      <h3 class="font-semibold mt-2">{previewData.title}</h3>
      <p class="text-sm text-ocean-600">{previewData.description}</p>
    {/if}
  </div>
{/if}
```

## Custom Delays

```typescript
// Quick preview for navigation
const quickPreview = createLinkPreview({
  openDelay: 300,
  closeDelay: 100
});

// Slower preview for detailed content
const detailedPreview = createLinkPreview({
  openDelay: 1000,
  closeDelay: 500
});
```
