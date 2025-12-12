<script lang="ts">
	import { clsx } from 'clsx';
	import type { SvelteComponent } from 'svelte';
	import type { IconProps } from 'lucide-svelte';

	interface Props {
		title: string;
		value: number | string;
		icon?: typeof SvelteComponent<IconProps>;
		color?: 'ocean' | 'green' | 'yellow' | 'red' | 'coral';
	}

	let { title, value, icon: Icon, color = 'ocean' }: Props = $props();

	const colorStyles = {
		ocean: 'bg-ocean-500/10 text-ocean-600',
		green: 'bg-ocean-400/10 text-ocean-500',
		yellow: 'bg-coral-400/10 text-coral-500',
		red: 'bg-coral-500/15 text-coral-600',
		coral: 'bg-coral-400/10 text-coral-500'
	};
</script>

<div class="bg-white/80 backdrop-blur-sm rounded-xl shadow-md border border-white/30 p-5 transition-all duration-200 hover:shadow-lg hover:-translate-y-0.5">
	<div class="flex items-center gap-4">
		{#if Icon}
			<div class={clsx('p-3 rounded-xl', colorStyles[color])}>
				<Icon class="w-6 h-6" />
			</div>
		{/if}
		<div>
			<p class="text-2xl font-bold text-ocean-900">{value}</p>
			<p class="text-sm text-ocean-500">{title}</p>
		</div>
	</div>
</div>
