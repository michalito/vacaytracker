<script lang="ts">
	import { Select } from 'bits-ui';
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
		/** Visible label, associated with the trigger. */
		label?: string;
		/** Accessible name when there is no visible label. */
		'aria-label'?: string;
		id?: string;
		onchange?: (value: string) => void;
	}

	const uid = $props.id();

	let {
		value = $bindable(''),
		options,
		placeholder = 'Select an option',
		disabled = false,
		label: labelText,
		'aria-label': ariaLabel,
		id = `select-${uid}`,
		onchange
	}: Props = $props();

	let selectOpen = $state(false);

	const displayLabel = $derived(options.find((o) => o.value === value)?.label ?? '');
</script>

<div class="flex flex-col gap-1.5">
	{#if labelText}
		<label for={id} class="text-sm font-semibold text-ocean-800">
			{labelText}
		</label>
	{/if}

	<Select.Root
		type="single"
		bind:value
		bind:open={selectOpen}
		onValueChange={(v) => onchange?.(v)}
		{disabled}
		items={options}
	>
		<Select.Trigger
			{id}
			aria-label={labelText ? undefined : ariaLabel}
			class="flex items-center justify-between rounded-lg border-2 border-ocean-500/40 bg-white px-4 py-2.5 text-ocean-900 cursor-pointer
				transition-all duration-200
				hover:border-ocean-500
				focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500
				data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed
				data-[placeholder]:text-ocean-500/50"
		>
			<span class="truncate">{displayLabel || placeholder}</span>
			<ChevronDown
				class="h-5 w-5 shrink-0 text-ocean-400 transition-transform duration-200 {selectOpen
					? 'rotate-180'
					: ''}"
				aria-hidden="true"
			/>
		</Select.Trigger>

		<Select.Portal>
			<Select.Content
				side="bottom"
				sideOffset={4}
				class="presence-pop z-50 max-h-60 w-[var(--bits-select-anchor-width)] overflow-auto rounded-xl border border-white/40 bg-white/95 p-1 shadow-xl backdrop-blur-md"
			>
				<Select.Viewport>
					{#each options as opt (opt.value)}
						<Select.Item
							value={opt.value}
							label={opt.label}
							disabled={opt.disabled}
							class="relative flex items-center justify-between rounded-lg px-4 py-2.5 cursor-pointer
								text-ocean-800 outline-none transition-colors
								data-[highlighted]:bg-ocean-500/10
								data-[selected]:bg-ocean-500 data-[selected]:text-white
								data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed"
						>
							{#snippet children({ selected })}
								<span>{opt.label}</span>
								{#if selected}
									<Check class="h-4 w-4" aria-hidden="true" />
								{/if}
							{/snippet}
						</Select.Item>
					{/each}
				</Select.Viewport>
			</Select.Content>
		</Select.Portal>
	</Select.Root>
</div>
