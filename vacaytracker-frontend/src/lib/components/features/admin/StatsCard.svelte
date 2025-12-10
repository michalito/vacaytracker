<script lang="ts">
	import { clsx } from 'clsx';
	import type { SvelteComponent } from 'svelte';
	import type { IconProps } from 'lucide-svelte';

	interface Props {
		title: string;
		value: number | string;
		icon: typeof SvelteComponent<IconProps>;
		color?: 'ocean' | 'green' | 'yellow' | 'red';
	}

	let { title, value, icon: Icon, color = 'ocean' }: Props = $props();

	const colorStyles = {
		ocean: 'bg-ocean-50 text-ocean-600',
		green: 'bg-green-50 text-green-600',
		yellow: 'bg-yellow-50 text-yellow-600',
		red: 'bg-red-50 text-red-600'
	};
</script>

<div class="bg-white rounded-lg shadow-md border border-sand-200 p-4">
	<div class="flex items-center gap-3">
		<div class={clsx('p-2 rounded-lg', colorStyles[color])}>
			<Icon class="w-6 h-6" />
		</div>
		<div>
			<p class="text-2xl font-bold text-ocean-900">{value}</p>
			<p class="text-sm text-ocean-500">{title}</p>
		</div>
	</div>
</div>
