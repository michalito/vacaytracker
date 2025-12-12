<script lang="ts">
	import { createSelect, melt } from '@melt-ui/svelte';
	import { ChevronDown, Check } from 'lucide-svelte';

	interface Option {
		value: string;
		label: string;
		disabled?: boolean;
	}

	interface Props {
		value?: string;
		options: Option[];
		placeholder?: string;
		disabled?: boolean;
		label?: string;
		onchange?: (value: string) => void;
	}

	let {
		value = $bindable(''),
		options,
		placeholder = 'Select an option',
		disabled = false,
		label: labelText,
		onchange
	}: Props = $props();

	// Compute initial selected value - using untrack to avoid Svelte 5 warning
	// The actual reactivity is handled in $effect below
	const initialOptions = $state.snapshot(options);
	const defaultSelected = initialOptions.find((o) => o.value === value);

	const {
		elements: { trigger, menu, option, label },
		states: { open, selectedLabel, selected, disabled: disabledState },
		helpers: { isSelected }
	} = createSelect<string>({
		defaultSelected: defaultSelected ? { value: defaultSelected.value, label: defaultSelected.label } : undefined,
		forceVisible: true,
		positioning: { placement: 'bottom', sameWidth: true },
		onSelectedChange: ({ next }) => {
			if (next) {
				value = next.value;
				onchange?.(next.value);
			}
			return next;
		}
	});

	// Sync disabled prop reactively
	$effect(() => {
		disabledState.set(disabled);
	});

	// Sync external value prop with internal state
	$effect(() => {
		const current = $selected?.value;
		if (value && value !== current) {
			const opt = options.find((o) => o.value === value);
			if (opt) {
				selected.set({ value: opt.value, label: opt.label });
			}
		}
	});
</script>

<div class="flex flex-col gap-1.5">
	{#if labelText}
		<label use:melt={$label} class="text-sm font-semibold text-ocean-800">
			{labelText}
		</label>
	{/if}

	<button
		use:melt={$trigger}
		class="flex items-center justify-between px-4 py-2.5 rounded-lg border-2 border-ocean-500/40 bg-white text-ocean-900 cursor-pointer
			transition-all duration-200
			hover:border-ocean-500
			focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500
			data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed
			data-[placeholder]:text-ocean-500/50"
	>
		<span>{$selectedLabel || placeholder}</span>
		<ChevronDown class="h-5 w-5 text-ocean-400 transition-transform duration-200 {$open ? 'rotate-180' : ''}" />
	</button>

	{#if $open}
		<div
			use:melt={$menu}
			class="z-50 max-h-60 overflow-auto rounded-xl bg-white/95 backdrop-blur-md p-1 shadow-lg border border-ocean-200"
		>
			{#each options as opt}
				<div
					use:melt={$option({ value: opt.value, label: opt.label, disabled: opt.disabled })}
					class="relative flex items-center justify-between rounded-lg px-4 py-2.5 cursor-pointer
						text-ocean-800 outline-none transition-colors
						data-[highlighted]:bg-ocean-100
						data-[selected]:bg-ocean-500 data-[selected]:text-white
						data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed"
				>
					<span>{opt.label}</span>
					{#if $isSelected(opt.value)}
						<Check class="h-4 w-4" />
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>
