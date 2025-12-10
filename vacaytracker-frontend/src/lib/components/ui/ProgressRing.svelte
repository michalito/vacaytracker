<script lang="ts">
	interface Props {
		value: number;
		max: number;
		size?: number;
		strokeWidth?: number;
		class?: string;
	}

	let { value, max, size = 120, strokeWidth = 8, class: className = '' }: Props = $props();

	const percentage = $derived(Math.min(100, Math.max(0, (value / max) * 100)));
	const radius = $derived((size - strokeWidth) / 2);
	const circumference = $derived(2 * Math.PI * radius);
	const offset = $derived(circumference - (percentage / 100) * circumference);

	const color = $derived(
		percentage > 70 ? 'text-success' : percentage > 30 ? 'text-warning' : 'text-error'
	);
</script>

<div class="relative inline-flex items-center justify-center {className}">
	<svg width={size} height={size} class="-rotate-90">
		<!-- Background circle -->
		<circle
			cx={size / 2}
			cy={size / 2}
			r={radius}
			fill="none"
			stroke="currentColor"
			stroke-width={strokeWidth}
			class="text-sand-200"
		/>
		<!-- Progress circle -->
		<circle
			cx={size / 2}
			cy={size / 2}
			r={radius}
			fill="none"
			stroke="currentColor"
			stroke-width={strokeWidth}
			stroke-linecap="round"
			stroke-dasharray={circumference}
			stroke-dashoffset={offset}
			class={color}
			style="transition: stroke-dashoffset 0.5s ease"
		/>
	</svg>
	<div class="absolute inset-0 flex flex-col items-center justify-center">
		<span class="text-2xl font-bold text-ocean-900">{value}</span>
		<span class="text-sm text-ocean-500">of {max}</span>
	</div>
</div>
