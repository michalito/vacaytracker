<script lang="ts">
	import { clsx } from 'clsx';
	import type { Snippet } from 'svelte';

	interface Props {
		variant?:
			| 'default'
			| 'success'
			| 'warning'
			| 'error'
			| 'info'
			| 'ocean'
			| 'coral';
		size?: 'sm' | 'md';
		class?: string;
		children: Snippet;
	}

	let { variant = 'default', size = 'md', class: className = '', children }: Props = $props();

	// Unified semantic color variants using theme tokens
	const variantStyles = {
		default: 'bg-slate-100 text-slate-700 border border-slate-200',
		success: 'bg-success-light text-green-800 border border-green-200',
		warning: 'bg-warning-light text-yellow-800 border border-yellow-200',
		error: 'bg-error-light text-red-800 border border-red-200',
		info: 'bg-info-light text-blue-800 border border-blue-200',
		ocean: 'bg-ocean-100 text-ocean-800 border border-ocean-200',
		coral: 'bg-coral-300/30 text-coral-600 border border-coral-300'
	};

	const sizeStyles = {
		sm: 'px-2 py-0.5 text-xs',
		md: 'px-2.5 py-1 text-sm'
	};

	const classes = $derived(
		clsx(
			'inline-flex items-center font-medium rounded-md',
			variantStyles[variant],
			sizeStyles[size],
			className
		)
	);
</script>

<span class={classes}>
	{@render children()}
</span>
