<script lang="ts">
	import { clsx } from 'clsx';
	import { Eye, EyeOff } from 'lucide-svelte';

	interface Props {
		type?: 'text' | 'email' | 'password' | 'number' | 'date';
		value?: string;
		placeholder?: string;
		label?: string;
		error?: string;
		disabled?: boolean;
		required?: boolean;
		id?: string;
		name?: string;
		class?: string;
		showPasswordToggle?: boolean;
		oninput?: (e: Event) => void;
		onblur?: (e: FocusEvent) => void;
	}

	let {
		type = 'text',
		value = $bindable(''),
		placeholder = '',
		label,
		error,
		disabled = false,
		required = false,
		id,
		name,
		class: className = '',
		showPasswordToggle = false,
		oninput,
		onblur
	}: Props = $props();

	let showPassword = $state(false);

	const inputId = $derived(id || `input-${crypto.randomUUID().slice(0, 8)}`);

	const effectiveType = $derived(
		type === 'password' && showPasswordToggle && showPassword ? 'text' : type
	);

	const inputClasses = $derived(
		clsx(
			'w-full px-4 py-2.5 rounded-lg border-2 transition-all duration-200 bg-white',
			'focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500',
			'hover:border-ocean-500',
			error
				? 'border-error text-error placeholder-error/50'
				: 'border-ocean-500/40 text-ocean-900 placeholder-ocean-500/50',
			disabled && 'bg-ocean-50 cursor-not-allowed opacity-60',
			showPasswordToggle && type === 'password' && 'pr-11',
			className
		)
	);

	function togglePasswordVisibility() {
		showPassword = !showPassword;
	}
</script>

<div class="w-full">
	{#if label}
		<label for={inputId} class="block text-sm font-semibold text-ocean-800 mb-1.5">
			{label}
			{#if required}
				<span class="text-coral-500">*</span>
			{/if}
		</label>
	{/if}

	<div class="relative">
		<input
			id={inputId}
			type={effectiveType}
			{name}
			{placeholder}
			{disabled}
			{required}
			class={inputClasses}
			bind:value
			{oninput}
			{onblur}
		/>

		{#if showPasswordToggle && type === 'password'}
			<button
				type="button"
				onclick={togglePasswordVisibility}
				class="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-ocean-400 hover:text-ocean-600 transition-colors duration-200 cursor-pointer focus:outline-none focus:text-ocean-500 rounded"
				aria-label={showPassword ? 'Hide password' : 'Show password'}
			>
				{#if showPassword}
					<EyeOff class="size-5" />
				{:else}
					<Eye class="size-5" />
				{/if}
			</button>
		{/if}
	</div>

	{#if error}
		<p class="mt-1 text-sm text-error">{error}</p>
	{/if}
</div>
