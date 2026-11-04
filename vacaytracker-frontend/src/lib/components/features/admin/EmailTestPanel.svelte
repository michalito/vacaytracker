<script lang="ts">
	import type { EmailTemplateType, EmailPreviewResponse } from '#lib/types/index.js';
	import Card from '#lib/components/ui/Card.svelte';
	import Button from '#lib/components/ui/Button.svelte';
	import Modal from '#lib/components/ui/Modal.svelte';
	import {
		Mail,
		Send,
		UserPlus,
		Clock,
		CheckCircle,
		XCircle,
		Bell,
		Eye,
		Newspaper
	} from 'lucide-svelte';
	import { adminApi } from '#lib/api/admin.js';
	import { toast } from '#lib/stores/toast.svelte.js';

	interface EmailTemplate {
		id: EmailTemplateType;
		name: string;
		description: string;
		icon: typeof Mail;
		color: string;
	}

	const templates: EmailTemplate[] = [
		{
			id: 'welcome',
			name: 'Welcome Email',
			description: 'Sent to new users with login credentials',
			icon: UserPlus,
			color: 'text-ocean-500'
		},
		{
			id: 'request_submitted',
			name: 'Request Submitted',
			description: 'Confirmation when vacation is requested',
			icon: Clock,
			color: 'text-warning'
		},
		{
			id: 'request_approved',
			name: 'Request Approved',
			description: 'Notification when request is approved',
			icon: CheckCircle,
			color: 'text-success'
		},
		{
			id: 'request_rejected',
			name: 'Request Rejected',
			description: 'Notification when request is rejected',
			icon: XCircle,
			color: 'text-error'
		},
		{
			id: 'admin_notification',
			name: 'Admin Notification',
			description: 'Alert to admins for new pending requests',
			icon: Bell,
			color: 'text-purple-500'
		},
		{
			id: 'newsletter',
			name: 'Newsletter',
			description: 'Monthly summary sent to all users',
			icon: Newspaper,
			color: 'text-ocean-600'
		}
	];

	let sendingTemplate = $state<EmailTemplateType | null>(null);
	let previewingTemplate = $state<EmailTemplateType | null>(null);
	let previewData = $state<EmailPreviewResponse | null>(null);
	let previewTemplate = $state<EmailTemplateType | null>(null);
	let previewOpen = $state(false);

	async function sendTestEmail(template: EmailTemplateType) {
		sendingTemplate = template;
		try {
			const result = await adminApi.sendTestEmail(template);
			toast.success(result.message);
		} catch (error) {
			toast.error(
				`Failed to send test email: ${error instanceof Error ? error.message : 'Unknown error'}`
			);
		} finally {
			sendingTemplate = null;
		}
	}

	async function previewEmail(template: EmailTemplateType) {
		previewingTemplate = template;
		try {
			previewData = await adminApi.previewEmail(template);
			previewTemplate = template;
			previewOpen = true;
		} catch (error) {
			toast.error(
				`Failed to load preview: ${error instanceof Error ? error.message : 'Unknown error'}`
			);
		} finally {
			previewingTemplate = null;
		}
	}

	function sendFromPreview() {
		if (!previewTemplate) return;
		const template = previewTemplate;
		previewOpen = false;
		sendTestEmail(template);
	}
</script>

<Card>
	{#snippet header()}
		<div class="flex items-center gap-2">
			<Mail class="w-5 h-5 text-ocean-500" />
			<h2 class="text-lg font-semibold text-ocean-700">Email Templates</h2>
		</div>
	{/snippet}

	<div class="space-y-4">
		<p class="text-sm text-ocean-600">
			Preview or send test emails to yourself to see how each template looks.
		</p>

		<div class="grid gap-3">
			{#each templates as template (template.id)}
				{@const Icon = template.icon}
				<div
					class="flex items-center justify-between p-3 bg-sand-50 rounded-lg border border-sand-200"
				>
					<div class="flex items-center gap-3">
						<div class="p-2 bg-white rounded-md shadow-sm">
							<Icon class="w-5 h-5 {template.color}" />
						</div>
						<div>
							<p class="font-medium text-ocean-800">{template.name}</p>
							<p class="text-xs text-ocean-500">{template.description}</p>
						</div>
					</div>
					<div class="flex gap-2">
						<Button
							variant="outline"
							size="sm"
							onclick={() => previewEmail(template.id)}
							loading={previewingTemplate === template.id}
							disabled={sendingTemplate !== null || previewingTemplate !== null}
						>
							<Eye class="w-4 h-4 mr-1" />
							Preview
						</Button>
						<Button
							variant="primary"
							size="sm"
							onclick={() => sendTestEmail(template.id)}
							loading={sendingTemplate === template.id}
							disabled={sendingTemplate !== null || previewingTemplate !== null}
						>
							<Send class="w-4 h-4 mr-1" />
							Send
						</Button>
					</div>
				</div>
			{/each}
		</div>

		<p class="text-xs text-ocean-500 pt-2 border-t border-sand-200">
			<strong>Note:</strong> Test emails use sample data and are sent to your admin email address.
		</p>
	</div>
</Card>

<!-- Preview Dialog -->
<Modal
	bind:open={previewOpen}
	title={previewData ? `${previewData.template} Preview` : 'Email Preview'}
	icon={Eye}
	size="xl"
	closeLabel="Close preview"
	onClosed={() => {
		previewData = null;
		previewTemplate = null;
	}}
>
	{#snippet description()}
		{#if previewData}
			<strong>Subject:</strong> {previewData.subject}
		{/if}
	{/snippet}

	{#if previewData}
		<div class="bg-ocean-50/30 p-4">
			<!-- Sandboxed: the template is a full HTML document whose styles must not leak into the app -->
			<iframe
				title="{previewData.template} email preview"
				srcdoc={previewData.htmlBody}
				sandbox=""
				class="h-[60vh] w-full rounded-xl border-0 bg-white shadow-sm"
			></iframe>
		</div>
	{/if}

	{#snippet footer()}
		<Button variant="outline" onclick={() => (previewOpen = false)}>Close</Button>
		<Button
			onclick={sendFromPreview}
			loading={sendingTemplate !== null}
			disabled={!previewTemplate}
		>
			<Send class="mr-1 h-4 w-4" aria-hidden="true" />
			Send Test Email
		</Button>
	{/snippet}
</Modal>
