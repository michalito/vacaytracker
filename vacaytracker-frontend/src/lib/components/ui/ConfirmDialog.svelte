<script lang="ts">
	import { AlertDialog } from 'bits-ui';
	import { clsx } from 'clsx';
	import { AlertTriangle } from 'lucide-svelte';
	import type { Snippet } from 'svelte';
	import Button from './Button.svelte';

	interface Props {
		/** Whether the dialog is open. Bindable. */
		open: boolean;
		title: string;
		description: string | Snippet;
		confirmText?: string;
		cancelText?: string;
		variant?: 'danger' | 'warning';
		/** Shows a spinner on the confirm button and blocks dismissal while the action runs. */
		isLoading?: boolean;
		onConfirm: () => void;
		/**
		 * Runs once the dialog has fully closed, however it was closed
		 * (Cancel, Escape, or `open = false` from code).
		 */
		onClosed?: () => void;
	}

	let {
		open = $bindable(false),
		title,
		description,
		confirmText = 'Confirm',
		cancelText = 'Cancel',
		variant = 'danger',
		isLoading = false,
		onConfirm,
		onClosed
	}: Props = $props();

	const tone = $derived(
		variant === 'danger'
			? { icon: 'bg-error/10 text-error', button: 'danger' as const }
			: { icon: 'bg-warning/10 text-warning', button: 'primary' as const }
	);
</script>

<AlertDialog.Root
	bind:open
	onOpenChangeComplete={(isOpen) => {
		if (!isOpen) onClosed?.();
	}}
>
	<AlertDialog.Portal>
		<AlertDialog.Overlay
			class="presence-fade fixed inset-0 z-50 bg-ocean-900/50 backdrop-blur-sm"
		/>

		<AlertDialog.Content
			escapeKeydownBehavior={isLoading ? 'ignore' : 'close'}
			class="presence-pop fixed left-1/2 top-1/2 z-50 -translate-x-1/2 -translate-y-1/2
				w-[calc(100%-2rem)] max-w-md max-h-[calc(100dvh-2rem)] overflow-y-auto
				rounded-2xl border border-white/40 bg-white/95 p-6 shadow-xl backdrop-blur-md focus:outline-none"
		>
			<div class="mb-4 flex items-center gap-3">
				<div class={clsx('shrink-0 rounded-full p-2', tone.icon)}>
					<AlertTriangle class="h-6 w-6" aria-hidden="true" />
				</div>
				<AlertDialog.Title class="text-lg font-semibold text-ocean-800">{title}</AlertDialog.Title>
			</div>

			<AlertDialog.Description class="mb-6 text-ocean-600">
				{#if typeof description === 'string'}
					{description}
				{:else}
					{@render description()}
				{/if}
			</AlertDialog.Description>

			<div class="flex gap-3">
				<AlertDialog.Cancel disabled={isLoading}>
					{#snippet child({ props })}
						<Button {...props} variant="outline" class="flex-1" disabled={isLoading}>
							{cancelText}
						</Button>
					{/snippet}
				</AlertDialog.Cancel>
				<AlertDialog.Action>
					{#snippet child({ props })}
						<Button
							{...props}
							variant={tone.button}
							class="flex-1"
							loading={isLoading}
							onclick={onConfirm}
						>
							{confirmText}
						</Button>
					{/snippet}
				</AlertDialog.Action>
			</div>
		</AlertDialog.Content>
	</AlertDialog.Portal>
</AlertDialog.Root>
