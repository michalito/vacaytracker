<script lang="ts">
	import type { NewsletterConfig } from '#lib/types/index.js';
	import Card from '#lib/components/ui/Card.svelte';
	import Button from '#lib/components/ui/Button.svelte';
	import Checkbox from '#lib/components/ui/Checkbox.svelte';
	import Select from '#lib/components/ui/Select.svelte';
	import Modal from '#lib/components/ui/Modal.svelte';
	import { Mail, Send, Eye } from 'lucide-svelte';
	import { admin } from '#lib/stores/admin.svelte.js';
	import { toast } from '#lib/stores/toast.svelte.js';

	interface Props {
		config: NewsletterConfig;
		onChange: (config: NewsletterConfig) => void;
	}

	let { config, onChange }: Props = $props();
	let isSending = $state(false);
	let isLoadingPreview = $state(false);
	let previewOpen = $state(false);

	const frequencyOptions = [
		{ value: 'weekly', label: 'Weekly' },
		{ value: 'monthly', label: 'Monthly' }
	];

	function toggleEnabled() {
		onChange({
			...config,
			enabled: !config.enabled
		});
	}

	function handleFrequencyChange(value: string) {
		onChange({
			...config,
			frequency: value as 'weekly' | 'monthly'
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
			previewOpen = true;
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
		<p class="text-sm text-ocean-600">Configure automatic email digests for your team.</p>

		<label class="flex items-center gap-3 cursor-pointer">
			<Checkbox checked={config.enabled} onchange={toggleEnabled} />
			<div>
				<p class="font-medium text-ocean-800">Enable newsletter</p>
				<p class="text-sm text-ocean-500">Send periodic vacation summaries to all users</p>
			</div>
		</label>

		{#if config.enabled}
			<div class="pt-4 border-t border-ocean-100/50 space-y-4">
				<Select
					label="Frequency"
					value={config.frequency}
					options={frequencyOptions}
					onchange={handleFrequencyChange}
				/>

				{#if config.frequency === 'monthly'}
					<div>
						<label for="dayOfMonth" class="block text-sm font-semibold text-ocean-800 mb-1.5">
							Day of Month
						</label>
						<input
							id="dayOfMonth"
							type="number"
							value={config.dayOfMonth}
							onchange={handleDayChange}
							min="1"
							max="28"
							class="w-full px-4 py-2.5 rounded-lg border-2 border-ocean-500/40 bg-white text-ocean-900 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500 hover:border-ocean-500"
						/>
						<p class="text-xs text-ocean-500 mt-1.5">
							Choose a day between 1-28 to ensure delivery in all months.
						</p>
					</div>
				{:else}
					<p class="text-sm text-ocean-600">Newsletter will be sent every Monday.</p>
				{/if}

				<div class="pt-4 border-t border-ocean-100/50">
					<p class="text-sm text-ocean-600">
						<strong>Last sent:</strong>
						{formatLastSent(config.lastSentAt)}
					</p>
				</div>
			</div>
		{/if}
	</div>
</Card>

<!-- Preview Dialog -->
<Modal
	bind:open={previewOpen}
	title="Newsletter Preview"
	icon={Eye}
	size="xl"
	closeLabel="Close preview"
>
	{#snippet description()}
		{#if admin.newsletterPreview}
			Will be sent to {admin.newsletterPreview.recipientCount}
			recipient{admin.newsletterPreview.recipientCount !== 1 ? 's' : ''}
		{/if}
	{/snippet}

	{#if admin.newsletterPreview}
		<div class="border-b border-ocean-100/50 bg-ocean-50/30 px-4 py-3">
			<p class="text-sm text-ocean-700">
				<strong>Subject:</strong>
				{admin.newsletterPreview.subject}
			</p>
		</div>
		<div class="p-4">
			<!-- Sandboxed: the newsletter is a full HTML document whose styles must not leak into the app -->
			<iframe
				title="Newsletter preview"
				srcdoc={admin.newsletterPreview.htmlBody}
				sandbox=""
				class="h-[60vh] w-full rounded-xl border-0 bg-white shadow-sm"
			></iframe>
		</div>
	{/if}

	{#snippet footer()}
		<Button variant="outline" onclick={() => (previewOpen = false)}>Close</Button>
		<Button onclick={handleSendNow} loading={isSending}>
			<Send class="mr-1 h-4 w-4" aria-hidden="true" />
			Send Now
		</Button>
	{/snippet}
</Modal>
