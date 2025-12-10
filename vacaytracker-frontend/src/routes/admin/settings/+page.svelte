<script lang="ts">
	import { admin } from '$lib/stores/admin.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { adminApi } from '$lib/api/admin';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import WeekendPolicy from '$lib/components/features/admin/WeekendPolicy.svelte';
	import NewsletterSettings from '$lib/components/features/admin/NewsletterSettings.svelte';
	import { Settings, Save, RotateCcw, AlertTriangle } from 'lucide-svelte';
	import type { WeekendPolicy as WeekendPolicyType, NewsletterConfig } from '$lib/types';

	let defaultVacationDays = $state(25);
	let vacationResetMonth = $state(1);
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
	let showResetModal = $state(false);
	let isResetting = $state(false);

	$effect(() => {
		admin.fetchSettings();
	});

	// Sync with store
	$effect(() => {
		if (admin.settings) {
			defaultVacationDays = admin.settings.defaultVacationDays;
			vacationResetMonth = admin.settings.vacationResetMonth;
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
				vacationResetMonth,
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
			showResetModal = false;
		} catch (error) {
			toast.error('Failed to reset balances');
		} finally {
			isResetting = false;
		}
	}

	const months = [
		{ value: 1, label: 'January' },
		{ value: 2, label: 'February' },
		{ value: 3, label: 'March' },
		{ value: 4, label: 'April' },
		{ value: 5, label: 'May' },
		{ value: 6, label: 'June' },
		{ value: 7, label: 'July' },
		{ value: 8, label: 'August' },
		{ value: 9, label: 'September' },
		{ value: 10, label: 'October' },
		{ value: 11, label: 'November' },
		{ value: 12, label: 'December' }
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
				<label for="defaultVacationDays" class="block text-sm font-medium text-ocean-700 mb-1">
					Default Vacation Days
				</label>
				<input
					id="defaultVacationDays"
					type="number"
					bind:value={defaultVacationDays}
					oninput={markChanged}
					min="0"
					max="365"
					class="w-full max-w-xs px-3 py-2 rounded-md border border-sand-300 focus:ring-2 focus:ring-ocean-500 focus:border-transparent"
				/>
				<p class="text-xs text-ocean-500 mt-1">
					Number of vacation days new employees receive by default.
				</p>
			</div>

			<div>
				<label for="vacationResetMonth" class="block text-sm font-medium text-ocean-700 mb-1">
					Vacation Reset Month
				</label>
				<select
					id="vacationResetMonth"
					bind:value={vacationResetMonth}
					onchange={markChanged}
					class="w-full max-w-xs px-3 py-2 rounded-md border border-sand-300 focus:ring-2 focus:ring-ocean-500 focus:border-transparent"
				>
					{#each months as month}
						<option value={month.value}>{month.label}</option>
					{/each}
				</select>
				<p class="text-xs text-ocean-500 mt-1">
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
			<Button variant="outline" onclick={() => (showResetModal = true)}>
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
</div>

<!-- Reset Confirmation Modal -->
{#if showResetModal}
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<div class="fixed inset-0 z-50">
		<div class="fixed inset-0 bg-black/50" onclick={() => (showResetModal = false)}></div>
		<div class="fixed inset-0 flex items-center justify-center p-4">
			<div
				class="bg-white rounded-xl shadow-xl w-full max-w-md p-6"
				onclick={(e) => e.stopPropagation()}
			>
				<div class="flex items-center gap-3 mb-4">
					<div class="p-2 bg-warning/10 rounded-full">
						<AlertTriangle class="w-6 h-6 text-warning" />
					</div>
					<h3 class="text-lg font-semibold text-ocean-800">Reset All Balances</h3>
				</div>
				<p class="text-ocean-600 mb-6">
					This will reset vacation balances for <strong>all employees</strong> to{' '}
					<strong>{defaultVacationDays} days</strong>. This action cannot be undone.
				</p>
				<div class="flex gap-3">
					<Button variant="outline" class="flex-1" onclick={() => (showResetModal = false)}>
						Cancel
					</Button>
					<Button variant="primary" class="flex-1" onclick={resetAllBalances} loading={isResetting}>
						Reset Balances
					</Button>
				</div>
			</div>
		</div>
	</div>
{/if}
