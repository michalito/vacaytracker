<script lang="ts">
	import { createDialog, melt } from '@melt-ui/svelte';
	import { admin } from '$lib/stores/admin.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { adminApi } from '$lib/api/admin';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import WeekendPolicy from '$lib/components/features/admin/WeekendPolicy.svelte';
	import NewsletterSettings from '$lib/components/features/admin/NewsletterSettings.svelte';
	import EmailTestPanel from '$lib/components/features/admin/EmailTestPanel.svelte';
	import { Settings, Save, RotateCcw, AlertTriangle, X } from 'lucide-svelte';
	import type { WeekendPolicy as WeekendPolicyType, NewsletterConfig } from '$lib/types';

	let defaultVacationDays = $state(25);
	let vacationResetMonth = $state('1');
	let weekendPolicy = $state<WeekendPolicyType>({
		excludeWeekends: true,
		excludedDays: [0, 6]
	});
	let newsletter = $state<NewsletterConfig>({
		enabled: false,
		frequency: 'monthly',
		dayOfMonth: 1
	});
	let isSaving = $state(false);
	let hasChanges = $state(false);
	let isResetting = $state(false);

	// Reset Confirmation Dialog
	const {
		elements: { overlay, content, title, description, close },
		states: { open: resetDialogOpen }
	} = createDialog({
		forceVisible: true,
		role: 'alertdialog'
	});

	$effect(() => {
		admin.fetchSettings();
	});

	// Sync with store
	$effect(() => {
		if (admin.settings) {
			defaultVacationDays = admin.settings.defaultVacationDays;
			vacationResetMonth = String(admin.settings.vacationResetMonth);
			weekendPolicy = { ...admin.settings.weekendPolicy };
			newsletter = { ...admin.settings.newsletter };
			hasChanges = false;
		}
	});

	function markChanged() {
		hasChanges = true;
	}

	async function saveSettings() {
		isSaving = true;
		try {
			await admin.updateSettings({
				defaultVacationDays,
				vacationResetMonth: parseInt(vacationResetMonth, 10),
				weekendPolicy,
				newsletter
			});
			toast.success('Settings saved');
			hasChanges = false;
		} catch (error) {
			toast.error('Failed to save settings');
		} finally {
			isSaving = false;
		}
	}

	async function resetAllBalances() {
		isResetting = true;
		try {
			const result = await adminApi.resetAllBalances();
			toast.success(`Reset ${result.usersUpdated} employee balances to ${result.newBalance} days`);
			resetDialogOpen.set(false);
		} catch (error) {
			toast.error('Failed to reset balances');
		} finally {
			isResetting = false;
		}
	}

	const months = [
		{ value: '1', label: 'January' },
		{ value: '2', label: 'February' },
		{ value: '3', label: 'March' },
		{ value: '4', label: 'April' },
		{ value: '5', label: 'May' },
		{ value: '6', label: 'June' },
		{ value: '7', label: 'July' },
		{ value: '8', label: 'August' },
		{ value: '9', label: 'September' },
		{ value: '10', label: 'October' },
		{ value: '11', label: 'November' },
		{ value: '12', label: 'December' }
	];
</script>

<svelte:head>
	<title>Settings - VacayTracker</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-ocean-800">Settings</h1>
			<p class="text-ocean-600">Configure your VacayTracker instance</p>
		</div>
		<Button onclick={saveSettings} loading={isSaving} disabled={!hasChanges}>
			<Save class="w-4 h-4 mr-2" />
			Save Changes
		</Button>
	</div>

	<!-- General Settings -->
	<Card>
		{#snippet header()}
			<div class="flex items-center gap-2">
				<Settings class="w-5 h-5 text-ocean-500" />
				<h2 class="text-lg font-semibold text-ocean-700">General Settings</h2>
			</div>
		{/snippet}

		<div class="space-y-4">
			<div>
				<label for="defaultVacationDays" class="block text-sm font-semibold text-ocean-800 mb-1.5">
					Default Vacation Days
				</label>
				<input
					id="defaultVacationDays"
					type="number"
					bind:value={defaultVacationDays}
					oninput={markChanged}
					min="0"
					max="365"
					class="w-full max-w-xs px-4 py-2.5 rounded-lg border-2 border-ocean-500/40 bg-white text-ocean-900 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500 hover:border-ocean-500"
				/>
				<p class="text-xs text-ocean-500 mt-1.5">
					Number of vacation days new employees receive by default.
				</p>
			</div>

			<div class="max-w-xs">
				<Select
					label="Vacation Reset Month"
					bind:value={vacationResetMonth}
					options={months}
					onchange={() => markChanged()}
				/>
				<p class="text-xs text-ocean-500 mt-1.5">
					Month when vacation balances are reset to the default.
				</p>
			</div>
		</div>
	</Card>

	<!-- Year-End Reset -->
	<Card>
		{#snippet header()}
			<div class="flex items-center gap-2">
				<RotateCcw class="w-5 h-5 text-ocean-500" />
				<h2 class="text-lg font-semibold text-ocean-700">Year-End Balance Reset</h2>
			</div>
		{/snippet}

		<div class="space-y-4">
			<p class="text-sm text-ocean-600">
				Reset all employee vacation balances to the default value ({defaultVacationDays} days).
				This is typically done at the start of a new year.
			</p>
			<Button variant="outline" onclick={() => resetDialogOpen.set(true)}>
				<RotateCcw class="w-4 h-4 mr-2" />
				Reset All Balances
			</Button>
		</div>
	</Card>

	<!-- Weekend Policy -->
	<WeekendPolicy
		policy={weekendPolicy}
		onChange={(p) => {
			weekendPolicy = p;
			markChanged();
		}}
	/>

	<!-- Newsletter Settings -->
	<NewsletterSettings
		config={newsletter}
		onChange={(c) => {
			newsletter = c;
			markChanged();
		}}
	/>

	<!-- Email Test Panel -->
	<EmailTestPanel />
</div>

<!-- Reset Confirmation Dialog -->
{#if $resetDialogOpen}
	<div use:melt={$overlay} class="fixed inset-0 z-50 bg-ocean-900/50 backdrop-blur-sm"></div>
	<div
		use:melt={$content}
		class="fixed left-1/2 top-1/2 z-50 -translate-x-1/2 -translate-y-1/2 bg-white/95 backdrop-blur-md rounded-2xl shadow-xl border border-white/30 w-full max-w-md p-6"
	>
		<div class="flex items-center justify-between mb-4">
			<div class="flex items-center gap-3">
				<div class="p-2 bg-warning/10 rounded-full">
					<AlertTriangle class="w-6 h-6 text-warning" />
				</div>
				<h3 use:melt={$title} class="text-lg font-semibold text-ocean-800">Reset All Balances</h3>
			</div>
			<button
				use:melt={$close}
				type="button"
				aria-label="Close dialog"
				class="p-1 rounded-lg text-ocean-400 hover:text-ocean-600 hover:bg-ocean-500/10 transition-all cursor-pointer"
			>
				<X class="w-5 h-5" />
			</button>
		</div>
		<p use:melt={$description} class="text-ocean-600 mb-6">
			This will reset vacation balances for <strong>all employees</strong> to{' '}
			<strong>{defaultVacationDays} days</strong>. This action cannot be undone.
		</p>
		<div class="flex gap-3">
			<Button variant="outline" class="flex-1" onclick={() => resetDialogOpen.set(false)}>
				Cancel
			</Button>
			<Button variant="primary" class="flex-1" onclick={resetAllBalances} loading={isResetting}>
				Reset Balances
			</Button>
		</div>
	</div>
{/if}
