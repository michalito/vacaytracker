<script lang="ts">
	import { createCheckbox, melt } from '@melt-ui/svelte';
	import { Check, Minus } from 'lucide-svelte';

	interface Props {
		checked?: boolean;
		disabled?: boolean;
		name?: string;
		value?: string;
		onchange?: (checked: boolean) => void;
	}

	let {
		checked = $bindable(false),
		disabled = false,
		name,
		value,
		onchange
	}: Props = $props();

	const {
		elements: { root, input },
		states: { checked: checkedState, disabled: disabledState },
		helpers: { isChecked, isIndeterminate }
	} = createCheckbox({
		defaultChecked: checked
	});

	// Sync external checked prop with internal state
	$effect.pre(() => {
		// Always update internal state when external checked prop changes
		checkedState.set(checked);
	});

	// Sync disabled prop reactively
	$effect(() => {
		disabledState.set(disabled);
	});

	// Handle click - call onchange with the NEW value (opposite of current)
	function handleClick() {
		if (disabled) return;
		const newValue = !checked;
		onchange?.(newValue);
	}
</script>

<button
	use:melt={$root}
	type="button"
	onclick={handleClick}
	class="h-5 w-5 shrink-0 rounded border-2 transition-all duration-200 flex items-center justify-center cursor-pointer
		data-[state=unchecked]:border-ocean-400 data-[state=unchecked]:bg-white data-[state=unchecked]:hover:border-ocean-500
		data-[state=checked]:border-ocean-500 data-[state=checked]:bg-ocean-500
		data-[state=indeterminate]:border-ocean-500 data-[state=indeterminate]:bg-ocean-500
		data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed
		focus-visible:ring-2 focus-visible:ring-ocean-500 focus-visible:ring-offset-2"
>
	{#if $isChecked}
		<Check class="h-3.5 w-3.5 text-white" />
	{:else if $isIndeterminate}
		<Minus class="h-3.5 w-3.5 text-white" />
	{/if}
</button>
<input use:melt={$input} class="sr-only" {name} {value} />
