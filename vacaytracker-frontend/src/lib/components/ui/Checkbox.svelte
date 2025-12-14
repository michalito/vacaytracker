<script lang="ts">
	import { createCheckbox, melt } from '@melt-ui/svelte';
	import { Check, Minus } from 'lucide-svelte';
	import { writable } from 'svelte/store';

	interface Props {
		checked?: boolean;
		disabled?: boolean;
		required?: boolean;
		name?: string;
		value?: string;
		onchange?: (checked: boolean) => void;
	}

	let {
		checked = $bindable(false),
		disabled = false,
		required = false,
		name,
		value = 'on',
		onchange
	}: Props = $props();

	// Controlled store (lets Melt-UI manage interactions while we keep `checked` bindable in sync)
	const checkedStore = writable<boolean | 'indeterminate'>(checked);

	const {
		elements: { root, input },
		helpers: { isChecked, isIndeterminate },
		options: { disabled: disabledOption, required: requiredOption, name: nameOption, value: valueOption }
	} = createCheckbox({
		checked: checkedStore,
		onCheckedChange: ({ next }) => {
			// Normalize Melt's tri-state to boolean for app usage
			const nextBool = next === 'indeterminate' ? false : next;
			checked = nextBool;
			onchange?.(nextBool);
			return next;
		}
	});

	// Sync external props into Melt option stores
	$effect(() => {
		disabledOption.set(disabled);
		requiredOption.set(required);
		nameOption.set(name);
		valueOption.set(value);
	});

	// Sync bindable checked into the controlled store (does not trigger onCheckedChange)
	$effect(() => {
		checkedStore.set(checked);
	});
</script>

<button
	use:melt={$root}
	type="button"
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
<input use:melt={$input} class="sr-only" />
