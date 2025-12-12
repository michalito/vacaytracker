# PIN Input

A segmented input for entering PIN codes, OTPs, and verification codes.

## Use Cases

- Two-factor authentication codes
- Phone verification
- Payment PIN entry
- Security codes

## Installation

```typescript
import { createPinInput } from '@melt-ui/svelte';
```

## API Reference

### createPinInput Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `placeholder` | `string` | `'○'` | Placeholder character |
| `defaultValue` | `string[]` | — | Initial values |
| `value` | `Writable<string[]>` | — | Controlled values |
| `disabled` | `boolean` | `false` | Disable input |
| `name` | `string` | — | Form input name |
| `type` | `'text' \| 'password'` | `'text'` | Input type |
| `onValueChange` | `ChangeFn` | — | Value change callback |
| `onComplete` | `(value: string) => void` | — | Called when all filled |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Container element |
| `input` | Individual input segment |
| `hiddenInput` | Hidden input for forms |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `value` | `Writable<string[]>` | Current values array |
| `valueStr` | `Readable<string>` | Joined value string |

## Data Attributes

### Root
- `[data-complete]` - All inputs filled
- `[data-melt-pin-input]` - Present on root

### Input
- `[data-complete]` - Input has value
- `[data-disabled]` - Input is disabled
- `[data-melt-pin-input-input]` - Present on inputs

## Keyboard Navigation

| Key | Action |
|-----|--------|
| `Arrow Left` | Previous input |
| `Arrow Right` | Next input |
| `Backspace` | Clear current and move back |
| `Delete` | Clear current input |
| `0-9` / `a-z` | Enter character |

## Accessibility

- Each input is properly labeled
- Screen reader announces position
- Paste support for full code

## Example

```svelte
<script lang="ts">
  import { createPinInput, melt } from '@melt-ui/svelte';

  const {
    elements: { root, input },
    states: { value, valueStr }
  } = createPinInput({
    placeholder: '○',
    onComplete: (code) => {
      console.log('Code entered:', code);
      verifyCode(code);
    }
  });

  const inputs = [0, 1, 2, 3, 4, 5]; // 6-digit code
</script>

<div class="flex flex-col items-center gap-4">
  <label class="text-sm font-semibold text-ocean-800">
    Enter verification code
  </label>

  <div use:melt={$root} class="flex gap-2">
    {#each inputs as i}
      <input
        use:melt={$input(i)}
        class="h-12 w-12 rounded-lg border-2 border-ocean-500/40 bg-white text-center text-xl font-semibold text-ocean-900
          focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500
          data-[complete]:border-ocean-500
          disabled:opacity-50"
        maxlength="1"
      />
    {/each}
  </div>

  <p class="text-sm text-ocean-500">
    Didn't receive a code? <button class="text-ocean-600 underline">Resend</button>
  </p>
</div>
```

## Styling with Tailwind

```css
/* Input focus */
[data-melt-pin-input-input]:focus {
  @apply ring-2 ring-ocean-500/30 border-ocean-500;
}

/* Completed input */
[data-melt-pin-input-input][data-complete] {
  @apply border-ocean-500;
}

/* Disabled state */
[data-melt-pin-input-input][data-disabled] {
  @apply opacity-50 cursor-not-allowed;
}

/* All complete indicator */
[data-melt-pin-input][data-complete] [data-melt-pin-input-input] {
  @apply border-green-500;
}
```

## Password/Hidden Input

```svelte
<script lang="ts">
  const { elements: { root, input } } = createPinInput({
    type: 'password',
    placeholder: '•'
  });
</script>

<div use:melt={$root} class="flex gap-2">
  {#each [0, 1, 2, 3] as i}
    <input
      use:melt={$input(i)}
      type="password"
      class="h-12 w-12 rounded-lg border-2 text-center text-2xl"
    />
  {/each}
</div>
```

## With Validation

```svelte
<script lang="ts">
  let error = $state('');
  let loading = $state(false);

  const {
    elements: { root, input },
    states: { valueStr }
  } = createPinInput({
    onComplete: async (code) => {
      loading = true;
      error = '';

      try {
        await verifyCode(code);
        // Success - redirect or show success
      } catch (e) {
        error = 'Invalid code. Please try again.';
        // Clear inputs
        value.set(['', '', '', '', '', '']);
      } finally {
        loading = false;
      }
    }
  });
</script>

<div class="flex flex-col items-center gap-4">
  <div use:melt={$root} class="flex gap-2">
    {#each [0, 1, 2, 3, 4, 5] as i}
      <input
        use:melt={$input(i)}
        class="h-12 w-12 rounded-lg border-2 text-center text-xl
          {error ? 'border-red-500' : 'border-ocean-500/40'}"
        disabled={loading}
      />
    {/each}
  </div>

  {#if error}
    <p class="text-sm text-red-500">{error}</p>
  {/if}

  {#if loading}
    <p class="text-sm text-ocean-500">Verifying...</p>
  {/if}
</div>
```

## Form Integration

```svelte
<script lang="ts">
  const {
    elements: { root, input, hiddenInput }
  } = createPinInput({
    name: 'verificationCode'
  });
</script>

<form onsubmit={handleSubmit}>
  <div use:melt={$root} class="flex gap-2">
    {#each [0, 1, 2, 3, 4, 5] as i}
      <input use:melt={$input(i)} class="..." />
    {/each}
  </div>
  <input use:melt={$hiddenInput} />
  <button type="submit">Verify</button>
</form>
```

## Custom Length

```svelte
<script lang="ts">
  // 4-digit PIN
  const pinInputs = [0, 1, 2, 3];

  // 8-character code
  const codeInputs = [0, 1, 2, 3, 4, 5, 6, 7];
</script>

<!-- 4-digit PIN -->
<div use:melt={$root}>
  {#each pinInputs as i}
    <input use:melt={$input(i)} />
  {/each}
</div>
```

## With Separator

```svelte
<div use:melt={$root} class="flex items-center gap-2">
  {#each [0, 1, 2] as i}
    <input use:melt={$input(i)} class="h-12 w-12 ..." />
  {/each}

  <span class="text-ocean-400 text-xl">-</span>

  {#each [3, 4, 5] as i}
    <input use:melt={$input(i)} class="h-12 w-12 ..." />
  {/each}
</div>
```
