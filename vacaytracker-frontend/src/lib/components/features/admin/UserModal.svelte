<script lang="ts">
	import { untrack } from 'svelte';
	import type { User, Role } from '#lib/types/index.js';
	import { admin } from '#lib/stores/admin.svelte.js';
	import { toast } from '#lib/stores/toast.svelte.js';
	import Modal from '#lib/components/ui/Modal.svelte';
	import Button from '#lib/components/ui/Button.svelte';
	import Input from '#lib/components/ui/Input.svelte';
	import Select from '#lib/components/ui/Select.svelte';
	import { UserPlus, Edit } from 'lucide-svelte';

	const roleOptions = [
		{ value: 'employee', label: 'Employee' },
		{ value: 'admin', label: 'Admin' }
	];

	interface Props {
		open: boolean;
		/** User to edit; `null` creates a new user. */
		user?: User | null;
		/** Balance pre-filled for new users (the company default). */
		defaultVacationDays?: number;
		/** Runs once the modal has fully closed, however it was closed. */
		onClosed?: () => void;
	}

	let {
		open = $bindable(false),
		user = null,
		defaultVacationDays = 25,
		onClosed
	}: Props = $props();

	let name = $state('');
	let email = $state('');
	let password = $state('');
	let role = $state<Role>('employee');
	// Number inputs bind an empty field as null/undefined
	let vacationBalance = $state<number | null | undefined>(25);
	let startDate = $state('');
	let isSubmitting = $state(false);
	let errors = $state<Record<string, string>>({});

	const isEditing = $derived(!!user);

	// Whether the admin has typed a balance since the form was loaded (not reactive on purpose)
	let balanceEdited = false;

	// Load the form every time the modal opens, so a cancelled edit never leaks into the next one
	$effect.pre(() => {
		if (!open) return;
		const current = user;
		untrack(() => loadForm(current));
	});

	// Company settings load asynchronously: if the default arrives after the create form opened,
	// adopt it unless the admin already entered a balance
	$effect.pre(() => {
		const companyDefault = defaultVacationDays;
		if (open && !user && !balanceEdited) vacationBalance = companyDefault;
	});

	function loadForm(current: User | null) {
		name = current?.name ?? '';
		email = current?.email ?? '';
		password = '';
		role = current?.role ?? 'employee';
		vacationBalance = current?.vacationBalance ?? defaultVacationDays;
		balanceEdited = false;
		startDate = current?.startDate ? current.startDate.split('T')[0] : '';
		errors = {};
	}

	function validate(): boolean {
		const next: Record<string, string> = {};

		if (!name.trim()) {
			next.name = 'Name is required';
		}

		if (!email.trim()) {
			next.email = 'Email is required';
		} else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.trim())) {
			next.email = 'Invalid email format';
		}

		if (!isEditing && !password) {
			next.password = 'Password is required for new users';
		} else if (password && password.length < 8) {
			next.password = 'Password must be at least 8 characters';
		}

		if (
			typeof vacationBalance !== 'number' ||
			!Number.isInteger(vacationBalance) ||
			vacationBalance < 0 ||
			vacationBalance > 365
		) {
			next.vacationBalance = 'Enter a whole number between 0 and 365';
		}

		errors = next;
		return Object.keys(next).length === 0;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!validate()) return;

		const balance = vacationBalance ?? 0;
		isSubmitting = true;
		try {
			if (user) {
				await admin.updateUser(user.id, {
					name: name.trim(),
					email: email.trim(),
					role,
					vacationBalance: balance,
					startDate: startDate || undefined
				});
				toast.success('User updated');
			} else {
				await admin.createUser({
					name: name.trim(),
					email: email.trim(),
					password,
					role,
					vacationBalance: balance,
					startDate: startDate || undefined
				});
				toast.success('User created');
			}
			open = false;
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Failed to save user');
		} finally {
			isSubmitting = false;
		}
	}
</script>

<Modal
	bind:open
	title={isEditing ? 'Edit User' : 'Create User'}
	icon={isEditing ? Edit : UserPlus}
	closeLabel="Close user form"
	{onClosed}
>
	<form onsubmit={handleSubmit} class="space-y-4 p-4" novalidate>
		<Input label="Name" placeholder="John Doe" bind:value={name} error={errors.name} required />

		<Input
			type="email"
			label="Email"
			placeholder="john@company.com"
			bind:value={email}
			error={errors.email}
			required
		/>

		{#if !isEditing}
			<Input
				type="password"
				label="Password"
				placeholder="Enter password"
				bind:value={password}
				error={errors.password}
				required
			/>
		{/if}

		<Select label="Role" bind:value={role} options={roleOptions} />

		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
			<div>
				<label
					for="user-vacation-balance"
					class="mb-1.5 block text-sm font-semibold text-ocean-800"
				>
					Vacation Balance
				</label>
				<input
					id="user-vacation-balance"
					type="number"
					bind:value={vacationBalance}
					oninput={() => (balanceEdited = true)}
					min="0"
					max="365"
					step="1"
					aria-invalid={errors.vacationBalance ? 'true' : undefined}
					aria-describedby={errors.vacationBalance ? 'user-vacation-balance-error' : undefined}
					class="w-full rounded-lg border-2 bg-white px-4 py-2.5 text-ocean-900 transition-all duration-200
						hover:border-ocean-500 focus:border-ocean-500 focus:outline-none focus:ring-2 focus:ring-ocean-500/30
						{errors.vacationBalance ? 'border-error' : 'border-ocean-500/40'}"
				/>
				{#if errors.vacationBalance}
					<p id="user-vacation-balance-error" class="mt-1.5 text-sm text-error">
						{errors.vacationBalance}
					</p>
				{/if}
			</div>
			<Input type="date" label="Start Date" bind:value={startDate} />
		</div>

		<div class="flex gap-3 pt-2">
			<Button type="button" variant="outline" class="flex-1" onclick={() => (open = false)}>
				Cancel
			</Button>
			<Button type="submit" variant="primary" class="flex-1" loading={isSubmitting}>
				{isEditing ? 'Save Changes' : 'Create User'}
			</Button>
		</div>
	</form>
</Modal>
