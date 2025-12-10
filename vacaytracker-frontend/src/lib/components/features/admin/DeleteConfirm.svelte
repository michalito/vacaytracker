<script lang="ts">
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
</script>

{#if open}
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<div class="fixed inset-0 z-50">
		<div class="fixed inset-0 bg-black/50 backdrop-blur-sm" onclick={onCancel}></div>
		<div class="fixed inset-0 flex items-center justify-center p-4">
			<div
				class="bg-white rounded-xl shadow-xl w-full max-w-sm"
				onclick={(e) => e.stopPropagation()}
			>
				<!-- Header -->
				<div class="flex items-center justify-between p-4 border-b border-sand-200">
					<div class="flex items-center gap-2 text-error">
						<AlertTriangle class="w-5 h-5" />
						<h2 class="text-lg font-semibold">{title}</h2>
					</div>
					<button type="button" onclick={onCancel} class="text-ocean-400 hover:text-ocean-600">
						<X class="w-5 h-5" />
					</button>
				</div>

				<!-- Body -->
				<div class="p-4">
					<p class="text-ocean-600">{message}</p>
				</div>

				<!-- Actions -->
				<div class="flex gap-3 p-4 border-t border-sand-200">
					<Button type="button" variant="outline" class="flex-1" onclick={onCancel}>
						Cancel
					</Button>
					<Button
						type="button"
						variant="primary"
						class="flex-1 !bg-error hover:!bg-red-600"
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
