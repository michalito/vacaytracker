<script lang="ts">
	import { createDialog, melt } from '@melt-ui/svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { AlertTriangle, X } from 'lucide-svelte';

	interface Props {
		open: boolean;
		title?: string;
		message?: string;
		confirmText?: string;
		isLoading?: boolean;
		onConfirm: () => void;
		onCancel: () => void;
	}

	let {
		open = $bindable(false),
		title = 'Confirm Delete',
		message = 'Are you sure you want to delete this item? This action cannot be undone.',
		confirmText = 'Delete',
		isLoading = false,
		onConfirm,
		onCancel
	}: Props = $props();

	// Create Melt-UI dialog with alertdialog role for confirmations
	const {
		elements: { overlay, content, title: dialogTitle, description, close, portalled },
		states: { open: dialogOpen }
	} = createDialog({
		role: 'alertdialog',
		forceVisible: true,
		closeOnOutsideClick: false, // Require explicit action
		onOpenChange: ({ next }) => {
			open = next;
			if (!next) onCancel();
			return next;
		}
	});

	// Sync external open prop with dialog state
	$effect(() => {
		dialogOpen.set(open);
	});

	function handleCancel() {
		dialogOpen.set(false);
	}
</script>

{#if open}
	<div use:melt={$portalled}>
		<!-- Overlay -->
		<div
			use:melt={$overlay}
			class="fixed inset-0 z-50 bg-ocean-900/50 backdrop-blur-sm transition-opacity duration-200
				data-[state=open]:opacity-100 data-[state=closed]:opacity-0"
		></div>

		<!-- Content Container -->
		<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
			<div
				use:melt={$content}
				class="bg-white/95 backdrop-blur-md rounded-2xl shadow-xl border border-white/30 w-full max-w-sm
					transition-all duration-200 data-[state=open]:opacity-100 data-[state=open]:scale-100
					data-[state=closed]:opacity-0 data-[state=closed]:scale-95"
			>
				<!-- Header -->
				<div class="flex items-center justify-between p-4 border-b border-ocean-100/50">
					<div class="flex items-center gap-2 text-error">
						<AlertTriangle class="w-5 h-5" />
						<h2 use:melt={$dialogTitle} class="text-lg font-semibold">{title}</h2>
					</div>
					<button
						use:melt={$close}
						type="button"
						aria-label="Close dialog"
						class="p-1 rounded-lg text-ocean-400 hover:text-ocean-600 hover:bg-ocean-500/10 transition-all duration-200 cursor-pointer"
					>
						<X class="w-5 h-5" />
					</button>
				</div>

				<!-- Body -->
				<div class="p-4">
					<p use:melt={$description} class="text-ocean-600">{message}</p>
				</div>

				<!-- Actions -->
				<div class="flex gap-3 p-4 border-t border-ocean-100/50">
					<Button type="button" variant="outline" class="flex-1" onclick={handleCancel}>
						Cancel
					</Button>
					<Button
						type="button"
						variant="danger"
						class="flex-1"
						onclick={onConfirm}
						loading={isLoading}
					>
						{confirmText}
					</Button>
				</div>
			</div>
		</div>
	</div>
{/if}
