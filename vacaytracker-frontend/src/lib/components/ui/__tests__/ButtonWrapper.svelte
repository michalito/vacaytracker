<script lang="ts">
	import Button from '../Button.svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	interface Props {
		text?: string;
		variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger' | 'success' | 'warning';
		size?: 'sm' | 'md' | 'lg';
		disabled?: boolean;
		loading?: boolean;
		type?: 'button' | 'submit' | 'reset';
		onclick?: HTMLButtonAttributes['onclick'];
		class?: string;
	}

	let { text = 'Click me', ...rest }: Props = $props();
</script>

<Button {...rest}>{text}</Button>
