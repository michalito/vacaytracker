# Progress

A visual indicator showing task completion or loading state.

## Use Cases

- File upload progress
- Form completion indicators
- Loading states
- Step completion tracking

## Installation

```typescript
import { createProgress } from '@melt-ui/svelte';
```

## API Reference

### createProgress Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `max` | `number` | `100` | Maximum progress value |
| `defaultValue` | `number` | `0` | Initial progress value |
| `value` | `Writable<number>` | — | Controlled progress value |
| `onValueChange` | `ChangeFn` | — | Value change callback |

### Returned Elements

| Element | Description |
|---------|-------------|
| `root` | Progress bar container |

### Returned States

| State | Type | Description |
|-------|------|-------------|
| `value` | `Writable<number>` | Current progress value |

### Returned Options

| Option | Type | Description |
|--------|------|-------------|
| `max` | `Writable<number>` | Maximum value |

## Data Attributes

### Root
- `[data-value]` - Current progress value
- `[data-state]` - `'indeterminate'`, `'complete'`, or `'loading'`
- `[data-max]` - Maximum value
- `[data-melt-progress]` - Present on root

## Accessibility

- Provides assistive technology support for progress bar
- Announces progress changes to screen readers
- Uses `role="progressbar"` with proper ARIA attributes

## Example

```svelte
<script lang="ts">
  import { createProgress, melt } from '@melt-ui/svelte';

  const {
    elements: { root },
    states: { value },
    options: { max }
  } = createProgress({
    defaultValue: 45
  });

  const percentage = $derived(($value / $max) * 100);
</script>

<div class="w-full max-w-md">
  <div class="flex justify-between mb-2">
    <span class="text-sm font-medium text-ocean-700">Progress</span>
    <span class="text-sm text-ocean-500">{percentage.toFixed(0)}%</span>
  </div>

  <div
    use:melt={$root}
    class="h-3 w-full overflow-hidden rounded-full bg-ocean-200"
  >
    <div
      class="h-full bg-ocean-500 transition-all duration-300"
      style="width: {percentage}%"
    ></div>
  </div>
</div>
```

## Styling with Tailwind

```css
/* Progress container */
[data-melt-progress] {
  @apply h-3 w-full overflow-hidden rounded-full bg-ocean-200;
}

/* Progress states */
[data-melt-progress][data-state="loading"] {
  /* Active loading state */
}

[data-melt-progress][data-state="complete"] {
  /* Completed state */
}

[data-melt-progress][data-state="indeterminate"] {
  /* Unknown progress */
}
```

## Animated Progress

```svelte
<script lang="ts">
  import { createProgress, melt } from '@melt-ui/svelte';
  import { tweened } from 'svelte/motion';
  import { cubicOut } from 'svelte/easing';

  const animatedValue = tweened(0, {
    duration: 400,
    easing: cubicOut
  });

  const {
    elements: { root },
    states: { value }
  } = createProgress();

  // Sync value with animation
  $effect(() => {
    animatedValue.set($value);
  });
</script>

<div use:melt={$root} class="h-3 w-full overflow-hidden rounded-full bg-ocean-200">
  <div
    class="h-full bg-ocean-500"
    style="width: {$animatedValue}%"
  ></div>
</div>
```

## Indeterminate State

```svelte
<script lang="ts">
  import { createProgress, melt } from '@melt-ui/svelte';

  const {
    elements: { root }
  } = createProgress({
    defaultValue: null // Indeterminate
  });
</script>

<div use:melt={$root} class="h-3 w-full overflow-hidden rounded-full bg-ocean-200">
  <div class="h-full w-1/3 bg-ocean-500 animate-progress-indeterminate"></div>
</div>

<style>
  @keyframes progress-indeterminate {
    0% { transform: translateX(-100%); }
    100% { transform: translateX(400%); }
  }

  .animate-progress-indeterminate {
    animation: progress-indeterminate 1.5s ease-in-out infinite;
  }
</style>
```

## With Steps

```svelte
<script lang="ts">
  import { createProgress, melt } from '@melt-ui/svelte';

  const steps = ['Details', 'Payment', 'Confirmation'];
  let currentStep = $state(1);

  const { elements: { root }, states: { value } } = createProgress({
    max: steps.length,
    defaultValue: currentStep
  });
</script>

<div class="space-y-4">
  <div use:melt={$root} class="h-2 w-full overflow-hidden rounded-full bg-ocean-200">
    <div
      class="h-full bg-ocean-500 transition-all"
      style="width: {($value / steps.length) * 100}%"
    ></div>
  </div>

  <div class="flex justify-between">
    {#each steps as step, i}
      <div class="flex flex-col items-center gap-1">
        <div
          class="h-8 w-8 rounded-full flex items-center justify-center text-sm font-medium
            {i < currentStep ? 'bg-ocean-500 text-white' : i === currentStep ? 'bg-ocean-100 text-ocean-800' : 'bg-ocean-100 text-ocean-400'}"
        >
          {i + 1}
        </div>
        <span class="text-xs text-ocean-600">{step}</span>
      </div>
    {/each}
  </div>
</div>
```

## Circular Progress

```svelte
<script lang="ts">
  import { createProgress, melt } from '@melt-ui/svelte';

  const { states: { value }, options: { max } } = createProgress({
    defaultValue: 70
  });

  const percentage = $derived(($value / $max) * 100);
  const circumference = 2 * Math.PI * 45; // radius = 45
  const offset = $derived(circumference - (percentage / 100) * circumference);
</script>

<div class="relative h-32 w-32">
  <svg class="h-full w-full -rotate-90" viewBox="0 0 100 100">
    <!-- Background circle -->
    <circle
      cx="50"
      cy="50"
      r="45"
      fill="none"
      stroke-width="8"
      class="stroke-ocean-200"
    />
    <!-- Progress circle -->
    <circle
      cx="50"
      cy="50"
      r="45"
      fill="none"
      stroke-width="8"
      stroke-linecap="round"
      class="stroke-ocean-500 transition-all duration-300"
      stroke-dasharray={circumference}
      stroke-dashoffset={offset}
    />
  </svg>
  <div class="absolute inset-0 flex items-center justify-center">
    <span class="text-2xl font-bold text-ocean-800">{percentage.toFixed(0)}%</span>
  </div>
</div>
```

## File Upload Progress

```svelte
<script lang="ts">
  import { createProgress, melt } from '@melt-ui/svelte';

  interface UploadProgress {
    filename: string;
    progress: number;
  }

  let uploads = $state<UploadProgress[]>([
    { filename: 'document.pdf', progress: 100 },
    { filename: 'image.png', progress: 67 },
    { filename: 'video.mp4', progress: 23 }
  ]);
</script>

<div class="space-y-4">
  {#each uploads as upload}
    {@const progress = createProgress({ defaultValue: upload.progress })}
    <div class="space-y-1">
      <div class="flex justify-between text-sm">
        <span class="text-ocean-700 truncate">{upload.filename}</span>
        <span class="text-ocean-500">{upload.progress}%</span>
      </div>
      <div
        use:melt={progress.elements.root}
        class="h-2 w-full overflow-hidden rounded-full bg-ocean-200"
      >
        <div
          class="h-full transition-all duration-300
            {upload.progress === 100 ? 'bg-green-500' : 'bg-ocean-500'}"
          style="width: {upload.progress}%"
        ></div>
      </div>
    </div>
  {/each}
</div>
```
