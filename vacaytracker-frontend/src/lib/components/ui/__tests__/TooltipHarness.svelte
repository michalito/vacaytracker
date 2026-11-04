<script lang="ts">
	// Test-only wrapper: the app gets Tooltip.Provider from the root layout,
	// so components that use tooltips need one when rendered on their own.
	import { Tooltip } from 'bits-ui';
	import type { Component } from 'svelte';

	interface Props {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		component: Component<any>;
		props?: Record<string, unknown>;
	}

	let { component: Inner, props = {} }: Props = $props();
</script>

<Tooltip.Provider>
	<Inner {...props} />
</Tooltip.Provider>
