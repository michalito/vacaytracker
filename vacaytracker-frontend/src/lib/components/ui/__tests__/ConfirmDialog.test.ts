import { render, screen, waitFor } from '@testing-library/svelte';
import { userEvent } from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import ConfirmDialog from '../ConfirmDialog.svelte';

function setup(props: Record<string, unknown> = {}) {
	const onConfirm = vi.fn();
	const onClosed = vi.fn();
	const utils = render(ConfirmDialog, {
		props: {
			open: true,
			title: 'Delete User',
			description: 'Are you sure?',
			confirmText: 'Delete',
			onConfirm,
			onClosed,
			...props
		}
	});
	return { ...utils, onConfirm, onClosed, user: userEvent.setup() };
}

describe('ConfirmDialog', () => {
	it('renders an alertdialog named by its title', async () => {
		setup();
		expect(await screen.findByRole('alertdialog', { name: 'Delete User' })).toBeInTheDocument();
		expect(screen.getByText('Are you sure?')).toBeInTheDocument();
	});

	it('Cancel closes the dialog and reports onClosed', async () => {
		const { user, onClosed, onConfirm } = setup();
		await user.click(await screen.findByRole('button', { name: 'Cancel' }));

		await waitFor(() => expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument());
		await waitFor(() => expect(onClosed).toHaveBeenCalledTimes(1));
		expect(onConfirm).not.toHaveBeenCalled();
	});

	it('Escape closes the dialog and reports onClosed', async () => {
		const { user, onClosed } = setup();
		await screen.findByRole('alertdialog');
		await user.keyboard('{Escape}');

		await waitFor(() => expect(onClosed).toHaveBeenCalledTimes(1));
	});

	it('reports onClosed when closed from code', async () => {
		const { rerender, onClosed } = setup();
		await screen.findByRole('alertdialog');
		await rerender({ open: false });

		await waitFor(() => expect(onClosed).toHaveBeenCalledTimes(1));
	});

	it('runs onConfirm and leaves closing to the caller', async () => {
		const { user, onConfirm, onClosed } = setup();
		await user.click(await screen.findByRole('button', { name: 'Delete' }));

		expect(onConfirm).toHaveBeenCalledTimes(1);
		expect(screen.getByRole('alertdialog')).toBeInTheDocument();
		expect(onClosed).not.toHaveBeenCalled();
	});

	it('cannot be dismissed while the action is running', async () => {
		const { user, onClosed } = setup({ isLoading: true });
		await screen.findByRole('alertdialog');

		expect(screen.getByRole('button', { name: 'Cancel' })).toBeDisabled();
		await user.keyboard('{Escape}');

		expect(screen.getByRole('alertdialog')).toBeInTheDocument();
		expect(onClosed).not.toHaveBeenCalled();
	});
});
