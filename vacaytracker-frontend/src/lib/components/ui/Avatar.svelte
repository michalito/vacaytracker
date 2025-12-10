<script lang="ts">
	import { clsx } from 'clsx';

	interface Props {
		name?: string;
		src?: string;
		size?: 'sm' | 'md' | 'lg';
		class?: string;
	}

	let { name = '', src, size = 'md', class: className = '' }: Props = $props();

	const initials = $derived(
		name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.toUpperCase()
			.slice(0, 2)
	);

	const sizeStyles = {
		sm: 'w-8 h-8 text-xs',
		md: 'w-10 h-10 text-sm',
		lg: 'w-12 h-12 text-base'
	};

	const classes = $derived(
		clsx(
			'inline-flex items-center justify-center rounded-full bg-ocean-500 text-white font-medium',
			sizeStyles[size],
			className
		)
	);
</script>

{#if src}
	<img {src} alt={name} class={clsx('rounded-full object-cover', sizeStyles[size], className)} />
{:else}
	<div class={classes}>
		{initials || '?'}
	</div>
{/if}
