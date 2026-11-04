<script lang="ts">
	import { Checkbox } from 'bits-ui';
	import { Check, Minus } from 'lucide-svelte';

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
</script>

<Checkbox.Root
	bind:checked
	onCheckedChange={(v) => {
		const nextBool = typeof v === 'boolean' ? v : false;
		onchange?.(nextBool);
	}}
	{disabled}
	{required}
	{name}
	{value}
	class="h-5 w-5 shrink-0 rounded border-2 transition-all duration-200 flex items-center justify-center cursor-pointer
		data-[state=unchecked]:border-ocean-400 data-[state=unchecked]:bg-white data-[state=unchecked]:hover:border-ocean-500
		data-[state=checked]:border-ocean-500 data-[state=checked]:bg-ocean-500
		data-[state=indeterminate]:border-ocean-500 data-[state=indeterminate]:bg-ocean-500
		data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed
		focus-visible:ring-2 focus-visible:ring-ocean-500 focus-visible:ring-offset-2"
>
	{#snippet children({ checked: isChecked, indeterminate })}
		{#if isChecked}
			<Check class="h-3.5 w-3.5 text-white" />
		{:else if indeterminate}
			<Minus class="h-3.5 w-3.5 text-white" />
		{/if}
	{/snippet}
</Checkbox.Root>
