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
			| 'pending'
			| 'approved'
			| 'rejected';
		size?: 'sm' | 'md';
		class?: string;
		children: Snippet;
	}

	let { variant = 'default', size = 'md', class: className = '', children }: Props = $props();

	const variantStyles = {
		default: 'bg-sand-200 text-ocean-700',
		success: 'bg-green-100 text-green-800',
		warning: 'bg-yellow-100 text-yellow-800',
		error: 'bg-red-100 text-red-800',
		info: 'bg-blue-100 text-blue-800',
		pending: 'bg-yellow-100 text-yellow-800',
		approved: 'bg-green-100 text-green-800',
		rejected: 'bg-red-100 text-red-800'
	};

	const sizeStyles = {
		sm: 'px-2 py-0.5 text-xs',
		md: 'px-2.5 py-1 text-sm'
	};

	const classes = $derived(
		clsx(
			'inline-flex items-center font-medium rounded-full',
			variantStyles[variant],
			sizeStyles[size],
			className
		)
	);
</script>

<span class={classes}>
	{@render children()}
</span>
