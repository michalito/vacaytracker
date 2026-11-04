import { render, screen, waitFor } from '@testing-library/svelte';
import { userEvent } from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import type { User } from '$lib/types';

vi.mock('$lib/stores/admin.svelte', () => ({
	admin: { createUser: vi.fn(), updateUser: vi.fn() }
}));

import { admin } from '$lib/stores/admin.svelte';
import UserModal from '../UserModal.svelte';
import UserModalHarness from './UserModalHarness.svelte';

const alice: User = {
	id: 'u1',
	email: 'alice@example.com',
	name: 'Alice Smith',
	role: 'employee',
	vacationBalance: 20,
	startDate: '2024-03-01',
	createdAt: '2024-03-01T00:00:00Z',
	updatedAt: '2024-03-01T00:00:00Z'
} as User;

function nameInput() {
	return screen.findByPlaceholderText<HTMLInputElement>('John Doe');
}

function balanceInput() {
	return screen.findByLabelText<HTMLInputElement>('Vacation Balance');
}

describe('UserModal', () => {
	beforeEach(() => {
		vi.mocked(admin.createUser).mockResolvedValue(undefined as never);
		vi.mocked(admin.updateUser).mockResolvedValue(undefined as never);
	});

	it('reports onClosed when cancelled', async () => {
		const user = userEvent.setup();
		const onClosed = vi.fn();
		render(UserModal, { props: { open: true, user: null, onClosed } });

		await user.click(await screen.findByRole('button', { name: 'Cancel' }));

		await waitFor(() => expect(onClosed).toHaveBeenCalledTimes(1));
	});

	it('discards unsaved edits when the same user is opened again', async () => {
		const user = userEvent.setup();
		// The harness only notifies props whose value changed, like the users page
		const { rerender } = render(UserModalHarness, { props: { open: true, user: alice } });

		const input = await nameInput();
		await user.clear(input);
		await user.type(input, 'Unsaved Edit');
		await user.click(screen.getByRole('button', { name: 'Cancel' }));
		await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());

		// Same user object, only `open` toggles, exactly like the users page does
		await rerender({ open: false });
		await rerender({ open: true });
		expect((await nameInput()).value).toBe('Alice Smith');
	});

	it('opens an empty create form with the company default balance after an edit', async () => {
		const user = userEvent.setup();
		const { rerender } = render(UserModalHarness, {
			props: { open: true, user: alice, defaultVacationDays: 30 }
		});
		await user.click(await screen.findByRole('button', { name: 'Cancel' }));
		await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());

		await rerender({ open: false });
		await rerender({ user: null, open: true });

		expect(await screen.findByRole('dialog', { name: 'Create User' })).toBeInTheDocument();
		expect((await nameInput()).value).toBe('');
		expect((await balanceInput()).valueAsNumber).toBe(30);
	});

	it('adopts a company default that arrives after the create form opened', async () => {
		// Settings not loaded yet: the page passes undefined, so the prop falls back to 25
		const { rerender } = render(UserModalHarness, {
			props: { open: true, user: null, defaultVacationDays: undefined }
		});
		expect((await balanceInput()).valueAsNumber).toBe(25);

		await rerender({ defaultVacationDays: 30 });

		await waitFor(async () => expect((await balanceInput()).valueAsNumber).toBe(30));
	});

	it('keeps a balance the admin typed when the company default arrives later', async () => {
		const user = userEvent.setup();
		const { rerender } = render(UserModalHarness, {
			props: { open: true, user: null, defaultVacationDays: undefined }
		});
		const balance = await balanceInput();
		await user.clear(balance);
		await user.type(balance, '12');

		await rerender({ defaultVacationDays: 30 });

		expect((await balanceInput()).valueAsNumber).toBe(12);
	});

	it('rejects an empty or out-of-range balance without calling the API', async () => {
		const user = userEvent.setup();
		render(UserModal, { props: { open: true, user: alice } });

		const balance = await balanceInput();
		await user.clear(balance);
		await user.click(screen.getByRole('button', { name: 'Save Changes' }));

		expect(await screen.findByText('Enter a whole number between 0 and 365')).toBeInTheDocument();
		expect(admin.updateUser).not.toHaveBeenCalled();
	});

	it('saves trimmed values and closes', async () => {
		const user = userEvent.setup();
		const onClosed = vi.fn();
		render(UserModal, { props: { open: true, user: alice, onClosed } });

		const input = await nameInput();
		await user.clear(input);
		await user.type(input, '  Alice Brown  ');
		await user.click(screen.getByRole('button', { name: 'Save Changes' }));

		expect(admin.updateUser).toHaveBeenCalledWith(
			'u1',
			expect.objectContaining({
				name: 'Alice Brown',
				email: 'alice@example.com',
				vacationBalance: 20
			})
		);
		await waitFor(() => expect(onClosed).toHaveBeenCalledTimes(1));
	});
});
