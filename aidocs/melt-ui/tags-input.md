# Tags Input

An input field for managing multiple tags or values.

## Use Cases

- Tag/label management
- Skill selection
- Email recipients
- Search filters

## Installation

```typescript
import { createTagsInput } from '@melt-ui/svelte';
```

## API Reference

### createTagsInput Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `defaultTags` | `string[]` | `[]` | Initial tags |
| `tags` | `Writable<string[]>` | — | Controlled tags |
| `unique` | `boolean` | `false` | Prevent duplicates |
| `trim` | `boolean` | `true` | Trim whitespace |
| `blur` | `'nothing' \| 'add' \| 'clear'` | `'nothing'` | Action on blur |
| `addOnPaste` | `boolean` | `false` | Add tags on paste |
| `maxTags` | `number` | `Infinity` | Maximum tag count |
| `allowed` | `string[]` | — | Allowed values only |
| `denied` | `string[]` | `[]` | Denied values |
| `add` | `(tag: string) => string \| undefined` | — | Custom add handler |
| `remove` | `(tag: string) => boolean` | — | Custom remove handler |
| `disabled` | `boolean` | `false` | Disable input |
| `name` | `string` | — | Form input name |
| `onTagsChange` | `ChangeFn` | — | Tags change callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Container element |
| `input` | Text input element |
| `tag` | Individual tag element |
| `deleteTrigger` | Tag delete button |
| `edit` | Tag edit input |
| `hiddenInput` | Hidden form input |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `tags` | `Writable<string[]>` | Current tags |

### Returned Helpers

| Helper | Description |
|--------|-------------|
| `addTag(tag)` | Programmatically add tag |
| `removeTag(tag)` | Programmatically remove tag |
| `updateTag(old, new)` | Update tag value |
| `isInputValid` | Check if current input is valid |

## Data Attributes

### Root
- `[data-disabled]` - Present when disabled
- `[data-focus]` - Present when input focused
- `[data-melt-tags-input]` - Present on root

### Tag
- `[data-tag]` - Tag value
- `[data-disabled]` - Present when disabled
- `[data-editing]` - Present when being edited
- `[data-melt-tags-input-tag]` - Present on tags

### Input
- `[data-invalid]` - Present when input invalid
- `[data-melt-tags-input-input]` - Present on input

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Enter` | Add current input as tag |
| `Backspace` | Delete last tag (when input empty) |
| `Arrow Left` | Select previous tag |
| `Arrow Right` | Select next tag |
| `Delete` | Delete selected tag |

## Accessibility

- Tags are focusable and deletable via keyboard
- Screen reader announces tag additions/removals
- Proper ARIA attributes for list behavior

## Example

```svelte
<script lang="ts">
  import { createTagsInput, melt } from '@melt-ui/svelte';
  import { X } from 'lucide-svelte';

  const {
    elements: { root, input, tag, deleteTrigger },
    states: { tags }
  } = createTagsInput({
    defaultTags: ['svelte', 'melt-ui'],
    unique: true
  });
</script>

<div
  use:melt={$root}
  class="flex flex-wrap items-center gap-2 p-2 rounded-lg border-2 border-ocean-500/40 bg-white
    focus-within:ring-2 focus-within:ring-ocean-500/30 focus-within:border-ocean-500"
>
  {#each $tags as t}
    <div
      use:melt={$tag(t)}
      class="flex items-center gap-1 px-2 py-1 bg-ocean-100 text-ocean-700 rounded-md text-sm"
    >
      <span>{t}</span>
      <button
        use:melt={$deleteTrigger(t)}
        class="p-0.5 rounded hover:bg-ocean-200 transition-colors"
      >
        <X class="h-3 w-3" />
      </button>
    </div>
  {/each}

  <input
    use:melt={$input}
    placeholder="Add tag..."
    class="flex-1 min-w-[100px] outline-none text-ocean-900 placeholder-ocean-400"
  />
</div>
```

## Styling with Tailwind

```css
/* Root container */
[data-melt-tags-input] {
  @apply flex flex-wrap items-center gap-2 p-2 rounded-lg border-2 border-ocean-500/40 bg-white;
}

[data-melt-tags-input][data-focus] {
  @apply ring-2 ring-ocean-500/30 border-ocean-500;
}

[data-melt-tags-input][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}

/* Tag */
[data-melt-tags-input-tag] {
  @apply flex items-center gap-1 px-2 py-1 bg-ocean-100 text-ocean-700 rounded-md text-sm;
}

[data-melt-tags-input-tag][data-editing] {
  @apply ring-2 ring-ocean-500;
}

/* Input */
[data-melt-tags-input-input] {
  @apply flex-1 min-w-[100px] outline-none bg-transparent;
}

[data-melt-tags-input-input][data-invalid] {
  @apply text-red-500;
}
```

## VacayTracker Skills Input

```svelte
<script lang="ts">
  import { createTagsInput, melt } from '@melt-ui/svelte';
  import { X } from 'lucide-svelte';

  const {
    elements: { root, input, tag, deleteTrigger },
    states: { tags }
  } = createTagsInput({
    unique: true,
    maxTags: 10,
    trim: true
  });
</script>

<div class="space-y-2">
  <label class="text-sm font-semibold text-ocean-800">Skills</label>

  <div
    use:melt={$root}
    class="flex flex-wrap items-center gap-2 p-3 rounded-lg border-2 border-ocean-500/40 bg-white
      data-[focus]:ring-2 data-[focus]:ring-ocean-500/30 data-[focus]:border-ocean-500"
  >
    {#each $tags as t}
      <span
        use:melt={$tag(t)}
        class="inline-flex items-center gap-1 px-3 py-1 bg-ocean-500/10 text-ocean-700 rounded-full text-sm"
      >
        {t}
        <button
          use:melt={$deleteTrigger(t)}
          class="p-0.5 rounded-full hover:bg-ocean-500/20 transition-colors"
        >
          <X class="h-3 w-3" />
        </button>
      </span>
    {/each}

    <input
      use:melt={$input}
      placeholder={$tags.length >= 10 ? 'Max 10 skills' : 'Add skill...'}
      disabled={$tags.length >= 10}
      class="flex-1 min-w-[120px] outline-none text-ocean-900 placeholder-ocean-400 bg-transparent"
    />
  </div>

  <p class="text-xs text-ocean-500">Press Enter to add a skill ({$tags.length}/10)</p>
</div>
```

## With Suggestions

```svelte
<script lang="ts">
  import { createTagsInput, melt } from '@melt-ui/svelte';

  const suggestions = ['JavaScript', 'TypeScript', 'Svelte', 'React', 'Vue'];

  const {
    elements: { root, input, tag, deleteTrigger },
    states: { tags }
  } = createTagsInput();

  let inputValue = $state('');
  let showSuggestions = $state(false);

  const filteredSuggestions = $derived(
    suggestions.filter(
      s => s.toLowerCase().includes(inputValue.toLowerCase()) &&
      !$tags.includes(s)
    )
  );
</script>

<div class="relative">
  <div use:melt={$root} class="...">
    {#each $tags as t}
      <!-- Tags -->
    {/each}

    <input
      use:melt={$input}
      bind:value={inputValue}
      onfocus={() => showSuggestions = true}
      onblur={() => setTimeout(() => showSuggestions = false, 200)}
    />
  </div>

  {#if showSuggestions && filteredSuggestions.length > 0}
    <div class="absolute z-10 mt-1 w-full bg-white rounded-lg shadow-lg border border-ocean-200 py-1">
      {#each filteredSuggestions as suggestion}
        <button
          onclick={() => {
            tags.update(t => [...t, suggestion]);
            inputValue = '';
          }}
          class="w-full px-4 py-2 text-left text-sm text-ocean-700 hover:bg-ocean-50"
        >
          {suggestion}
        </button>
      {/each}
    </div>
  {/if}
</div>
```

## Email Recipients

```svelte
<script lang="ts">
  const { elements, states, helpers } = createTagsInput({
    unique: true,
    add: (value) => {
      // Validate email format
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (emailRegex.test(value)) {
        return value.toLowerCase();
      }
      return undefined; // Invalid, don't add
    }
  });
</script>

<div use:melt={elements.root} class="...">
  {#each $states.tags as email}
    <span use:melt={elements.tag(email)} class="...">
      {email}
      <button use:melt={elements.deleteTrigger(email)}>×</button>
    </span>
  {/each}

  <input
    use:melt={elements.input}
    type="email"
    placeholder="Add email..."
  />
</div>
```

## Form Integration

```svelte
<script lang="ts">
  const {
    elements: { root, input, tag, deleteTrigger, hiddenInput }
  } = createTagsInput({
    name: 'skills'
  });
</script>

<form onsubmit={handleSubmit}>
  <div use:melt={$root}>
    <!-- Tags and input -->
  </div>
  <input use:melt={$hiddenInput} />
  <button type="submit">Submit</button>
</form>
```
