<script lang="ts">
	import { Avatar } from 'bits-ui';
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

	const containerClasses = $derived(
		clsx(
			'relative inline-flex items-center justify-center rounded-full bg-ocean-500 overflow-hidden',
			sizeStyles[size],
			className
		)
	);
</script>

<Avatar.Root class={containerClasses}>
	{#if src}
		<Avatar.Image
			{src}
			alt={name}
			class="h-full w-full object-cover transition-opacity duration-200
				data-[status=loading]:opacity-0 data-[status=loaded]:opacity-100 data-[status=error]:opacity-0"
		/>
	{/if}
	<Avatar.Fallback
		class="absolute inset-0 flex items-center justify-center text-white font-medium
			data-[status=loaded]:hidden"
	>
		{initials || '?'}
	</Avatar.Fallback>
</Avatar.Root>
