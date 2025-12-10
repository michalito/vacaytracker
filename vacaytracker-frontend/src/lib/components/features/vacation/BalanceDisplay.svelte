<script lang="ts">
	import ProgressRing from '$lib/components/ui/ProgressRing.svelte';

	interface Props {
		current: number;
		total: number;
		size?: number;
	}

	let { current, total, size = 140 }: Props = $props();
</script>

<ProgressRing value={current} max={total} {size} />
