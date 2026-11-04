<script lang="ts">
	// Test-only parent for UserModal. testing-library's rerender() swaps the whole props object,
	// which makes every prop look changed. A real parent only notifies the props whose value
	// changed; routing each prop through $derived (which only propagates real changes) restores that,
	// so tests cannot pass just because an unrelated prop "changed".
	import type { User } from '#lib/types/index.js';
	import UserModal from '../UserModal.svelte';

	interface Props {
		open: boolean;
		user: User | null;
		defaultVacationDays?: number;
		onClosed?: () => void;
	}

	let props: Props = $props();

	const open = $derived(props.open);
	const user = $derived(props.user);
	const defaultVacationDays = $derived(props.defaultVacationDays);
</script>

<UserModal {open} {user} {defaultVacationDays} onClosed={props.onClosed} />
