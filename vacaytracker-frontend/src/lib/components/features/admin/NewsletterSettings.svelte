<script lang="ts">
	import type { NewsletterConfig } from '$lib/types';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Mail, Send, Eye, X } from 'lucide-svelte';
	import { admin } from '$lib/stores/admin.svelte';
	import { toast } from '$lib/stores/toast.svelte';

	interface Props {
		config: NewsletterConfig;
		onChange: (config: NewsletterConfig) => void;
	}

	let { config, onChange }: Props = $props();
	let isSending = $state(false);
	let isLoadingPreview = $state(false);
	let showPreview = $state(false);

	function toggleEnabled() {
		onChange({
			...config,
			enabled: !config.enabled
		});
	}

	function handleFrequencyChange(e: Event) {
		const target = e.target as HTMLSelectElement;
		onChange({
			...config,
			frequency: target.value as 'weekly' | 'monthly'
		});
	}

	function handleDayChange(e: Event) {
		const target = e.target as HTMLInputElement;
		onChange({
			...config,
			dayOfMonth: parseInt(target.value, 10)
		});
	}

	async function handleSendNow() {
		isSending = true;
		try {
			const result = await admin.sendNewsletter();
			toast.success(result.message);
		} catch (error) {
			toast.error('Failed to send newsletter');
		} finally {
			isSending = false;
		}
	}

	async function handlePreview() {
		isLoadingPreview = true;
		try {
			await admin.fetchNewsletterPreview();
			showPreview = true;
		} catch (error) {
			toast.error('Failed to load preview');
		} finally {
			isLoadingPreview = false;
		}
	}

	function formatLastSent(lastSentAt: string | null | undefined): string {
		if (!lastSentAt) return 'Never';
		return new Date(lastSentAt).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<Card>
	{#snippet header()}
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-2">
				<Mail class="w-5 h-5 text-ocean-500" />
				<h2 class="text-lg font-semibold text-ocean-700">Newsletter Settings</h2>
			</div>
			{#if config.enabled}
				<div class="flex gap-2">
					<Button variant="outline" size="sm" onclick={handlePreview} loading={isLoadingPreview}>
						<Eye class="w-4 h-4 mr-1" />
						Preview
					</Button>
					<Button size="sm" onclick={handleSendNow} loading={isSending}>
						<Send class="w-4 h-4 mr-1" />
						Send Now
					</Button>
				</div>
			{/if}
		</div>
	{/snippet}

	<div class="space-y-4">
		<p class="text-sm text-ocean-600">
			Configure automatic email digests for your team.
		</p>

		<label class="flex items-center gap-3 cursor-pointer">
			<input
				type="checkbox"
				checked={config.enabled}
				onchange={toggleEnabled}
				class="w-5 h-5 rounded border-sand-300 text-ocean-500 focus:ring-ocean-500"
			/>
			<div>
				<p class="font-medium text-ocean-800">Enable newsletter</p>
				<p class="text-sm text-ocean-500">
					Send periodic vacation summaries to all users
				</p>
			</div>
		</label>

		{#if config.enabled}
			<div class="pt-4 border-t border-sand-200 space-y-4">
				<div>
					<label for="frequency" class="block text-sm font-medium text-ocean-700 mb-1">
						Frequency
					</label>
					<select
						id="frequency"
						value={config.frequency}
						onchange={handleFrequencyChange}
						class="w-full px-3 py-2 rounded-md border border-sand-300 focus:ring-2 focus:ring-ocean-500 focus:border-transparent"
					>
						<option value="weekly">Weekly</option>
						<option value="monthly">Monthly</option>
					</select>
				</div>

				{#if config.frequency === 'monthly'}
					<div>
						<label for="dayOfMonth" class="block text-sm font-medium text-ocean-700 mb-1">
							Day of Month
						</label>
						<input
							id="dayOfMonth"
							type="number"
							value={config.dayOfMonth}
							onchange={handleDayChange}
							min="1"
							max="28"
							class="w-full px-3 py-2 rounded-md border border-sand-300 focus:ring-2 focus:ring-ocean-500 focus:border-transparent"
						/>
						<p class="text-xs text-ocean-500 mt-1">
							Choose a day between 1-28 to ensure delivery in all months.
						</p>
					</div>
				{:else}
					<p class="text-sm text-ocean-600">
						Newsletter will be sent every Monday.
					</p>
				{/if}

				<div class="pt-4 border-t border-sand-200">
					<p class="text-sm text-ocean-600">
						<strong>Last sent:</strong> {formatLastSent(config.lastSentAt)}
					</p>
				</div>
			</div>
		{/if}
	</div>
</Card>

<!-- Preview Modal -->
{#if showPreview && admin.newsletterPreview}
	<div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
		<div class="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-hidden flex flex-col">
			<div class="flex items-center justify-between p-4 border-b border-sand-200">
				<div>
					<h3 class="text-lg font-semibold text-ocean-800">Newsletter Preview</h3>
					<p class="text-sm text-ocean-600">
						Will be sent to {admin.newsletterPreview.recipientCount} recipient{admin.newsletterPreview.recipientCount !== 1 ? 's' : ''}
					</p>
				</div>
				<button
					onclick={() => (showPreview = false)}
					class="p-2 hover:bg-sand-100 rounded-md transition-colors"
				>
					<X class="w-5 h-5 text-ocean-600" />
				</button>
			</div>
			<div class="p-4 border-b border-sand-200 bg-sand-50">
				<p class="text-sm">
					<strong>Subject:</strong> {admin.newsletterPreview.subject}
				</p>
			</div>
			<div class="flex-1 overflow-auto p-4">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html admin.newsletterPreview.htmlBody}
			</div>
			<div class="flex justify-end gap-2 p-4 border-t border-sand-200 bg-sand-50">
				<Button variant="outline" onclick={() => (showPreview = false)}>
					Close
				</Button>
				<Button onclick={handleSendNow} loading={isSending}>
					<Send class="w-4 h-4 mr-1" />
					Send Now
				</Button>
			</div>
		</div>
	</div>
{/if}
