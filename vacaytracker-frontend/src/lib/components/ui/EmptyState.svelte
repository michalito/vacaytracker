<script lang="ts">
	import type { SvelteComponent } from 'svelte';
	import type { IconProps } from 'lucide-svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		icon: typeof SvelteComponent<IconProps>;
		message: string;
		iconSize?: 'sm' | 'md' | 'lg';
		children?: Snippet;
	}

	let { icon: Icon, message, iconSize = 'md', children }: Props = $props();

	const iconSizeClasses = {
		sm: 'w-8 h-8',
		md: 'w-10 h-10',
		lg: 'w-12 h-12'
	};

	const paddingClasses = {
		sm: 'py-4',
		md: 'py-6',
		lg: 'py-8'
	};
</script>

<div class="{paddingClasses[iconSize]} text-center">
	<Icon class="{iconSizeClasses[iconSize]} mx-auto text-ocean-300 mb-2" />
	<p class="text-ocean-500">{message}</p>
	{#if children}
		<div class="mt-4">
			{@render children()}
		</div>
	{/if}
</div>
