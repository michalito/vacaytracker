<script lang="ts">
	// Test-only: calls `run` inside an $effect, the way a page or component calls a store command
	let { run }: { run: () => void } = $props();

	$effect(() => {
		run();
	});
</script>
