<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { formatDateLong } from '$lib/utils/date';
	import Card from '$lib/components/ui/Card.svelte';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import PasswordChange from '$lib/components/features/auth/PasswordChange.svelte';
	import EmailPreferences from '$lib/components/features/auth/EmailPreferences.svelte';
	import { User, Calendar, Briefcase } from 'lucide-svelte';
</script>

<svelte:head>
	<title>Settings - VacayTracker</title>
</svelte:head>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold text-ocean-800">Settings</h1>
		<p class="text-ocean-600">Manage your account and preferences</p>
	</div>

	<!-- Profile Card -->
	<Card>
		{#snippet header()}
			<div class="flex items-center gap-2">
				<User class="w-5 h-5 text-ocean-500" />
				<h2 class="text-lg font-semibold text-ocean-700">Profile</h2>
			</div>
		{/snippet}

		<div class="flex items-start gap-6">
			<Avatar name={auth.user?.name} size="lg" />
			<div class="flex-1 space-y-4">
				<div>
					<p class="text-sm text-ocean-500">Name</p>
					<p class="font-medium text-ocean-800">{auth.user?.name}</p>
				</div>
				<div>
					<p class="text-sm text-ocean-500">Email</p>
					<p class="font-medium text-ocean-800">{auth.user?.email}</p>
				</div>
				<div class="flex gap-8">
					<div>
						<p class="text-sm text-ocean-500">Role</p>
						<p class="font-medium text-ocean-800 capitalize">{auth.user?.role}</p>
					</div>
					<div>
						<p class="text-sm text-ocean-500">Start Date</p>
						<p class="font-medium text-ocean-800">{formatDateLong(auth.user?.startDate)}</p>
					</div>
				</div>
			</div>
		</div>
	</Card>

	<!-- Vacation Balance Card -->
	<Card>
		{#snippet header()}
			<div class="flex items-center gap-2">
				<Briefcase class="w-5 h-5 text-ocean-500" />
				<h2 class="text-lg font-semibold text-ocean-700">Vacation Balance</h2>
			</div>
		{/snippet}

		<div class="flex items-center justify-between">
			<div>
				<p class="text-3xl font-bold text-ocean-800">{auth.user?.vacationBalance ?? 0}</p>
				<p class="text-sm text-ocean-500">days remaining this year</p>
			</div>
			<div class="p-4 bg-ocean-50 rounded-lg">
				<Calendar class="w-8 h-8 text-ocean-500" />
			</div>
		</div>
		<p class="mt-4 text-sm text-ocean-500">
			Contact your administrator if you believe there's an error with your balance.
		</p>
	</Card>

	<!-- Password Change -->
	<PasswordChange />

	<!-- Email Preferences -->
	<EmailPreferences />
</div>
