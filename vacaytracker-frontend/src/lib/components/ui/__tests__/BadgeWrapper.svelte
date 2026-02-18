<script lang="ts">
	import Badge from '../Badge.svelte';

	interface Props {
		text?: string;
		variant?: 'default' | 'success' | 'warning' | 'error' | 'info' | 'ocean' | 'coral';
		size?: 'sm' | 'md';
		class?: string;
	}

	let { text = 'Badge', ...rest }: Props = $props();
</script>

<Badge {...rest}>{text}</Badge>
