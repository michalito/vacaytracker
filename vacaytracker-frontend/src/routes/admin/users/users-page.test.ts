import { render, screen, waitFor } from '@testing-library/svelte';
import { userEvent } from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import type { User } from '$lib/types';

// Only the HTTP layer is mocked: the real admin store (and its cache) and the real page run
vi.mock('$lib/api/admin', () => ({
	adminApi: { listUsers: vi.fn(), getSettings: vi.fn() }
}));

import { adminApi } from '$lib/api/admin';
import UsersPage from './+page.svelte';

function makeUser(id: string, name: string, role: User['role']): User {
	return {
		id,
		name,
		role,
		email: `${id}@example.com`,
		vacationBalance: 20,
		emailPreferences: { vacationUpdates: true, weeklyDigest: false, teamNotifications: true }
	} as User;
}

const everyone = [makeUser('u1', 'Alice Smith', 'employee'), makeUser('u2', 'Bob Jones', 'admin')];

async function chooseRole(label: string) {
	const user = userEvent.setup();
	await user.click(screen.getByRole('combobox', { name: 'Filter by role' }));
	// jsdom never finishes floating-ui positioning, so options stay `visibility: hidden` and get
	// an empty accessible name; match on text content instead
	const option = await waitFor(() => {
		const match = screen
			.getAllByRole('option', { hidden: true })
			.find((el) => el.textContent?.trim() === label);
		if (!match) throw new Error(`No "${label}" option`);
		return match;
	});
	await user.click(option);
}

describe('admin users page', () => {
	beforeEach(() => {
		vi.mocked(adminApi.listUsers).mockImplementation(async (params) => {
			const users = params?.role ? everyone.filter((u) => u.role === params.role) : everyone;
			return {
				users,
				pagination: { page: 1, limit: 10, total: users.length, totalPages: 1 }
			};
		});
		vi.mocked(adminApi.getSettings).mockRejectedValue(new Error('not needed here'));
	});

	it('shows everyone again after filtering by role and switching back to All Roles', async () => {
		render(UsersPage);
		expect(await screen.findByText('Alice Smith')).toBeInTheDocument();

		await chooseRole('Admin');
		await waitFor(() => expect(screen.queryByText('Alice Smith')).not.toBeInTheDocument());
		expect(adminApi.listUsers).toHaveBeenLastCalledWith(expect.objectContaining({ role: 'admin' }));

		await chooseRole('All Roles');
		expect(await screen.findByText('Alice Smith')).toBeInTheDocument();
		expect(screen.getByText('Bob Jones')).toBeInTheDocument();
		expect(adminApi.listUsers).toHaveBeenLastCalledWith(
			expect.objectContaining({ role: undefined })
		);
	});

	it('does not refetch on every keystroke in the search box', async () => {
		const user = userEvent.setup();
		render(UsersPage);
		await screen.findByText('Alice Smith');
		const callsAfterLoad = vi.mocked(adminApi.listUsers).mock.calls.length;

		await user.type(screen.getByRole('searchbox', { name: 'Search users' }), 'bob');
		expect(adminApi.listUsers).toHaveBeenCalledTimes(callsAfterLoad);

		await user.keyboard('{Enter}');
		await waitFor(() =>
			expect(adminApi.listUsers).toHaveBeenLastCalledWith(
				expect.objectContaining({ search: 'bob' })
			)
		);
	});
});
