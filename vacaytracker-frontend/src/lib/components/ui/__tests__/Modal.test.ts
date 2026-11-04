import { render, screen, waitFor } from '@testing-library/svelte';
import { userEvent } from '@testing-library/user-event';
import { createRawSnippet } from 'svelte';
import { describe, it, expect, vi } from 'vitest';
import Modal from '../Modal.svelte';

const body = createRawSnippet(() => ({ render: () => '<p>Modal body</p>' }));

function setup(props: Record<string, unknown> = {}) {
	const onClosed = vi.fn();
	const utils = render(Modal, {
		props: { open: true, title: 'Edit User', children: body, onClosed, ...props }
	});
	return { ...utils, onClosed, user: userEvent.setup() };
}

describe('Modal', () => {
	it('renders a dialog named by its title with the body and description', async () => {
		setup({ description: 'Update the details below' });

		const dialog = await screen.findByRole('dialog', { name: 'Edit User' });
		expect(dialog).toHaveAccessibleDescription('Update the details below');
		expect(screen.getByText('Modal body')).toBeInTheDocument();
	});

	it('gives the icon-only close button an accessible name', async () => {
		setup({ closeLabel: 'Close user form' });
		expect(await screen.findByRole('button', { name: 'Close user form' })).toBeInTheDocument();
	});

	it('defaults the close button label', async () => {
		setup();
		expect(await screen.findByRole('button', { name: 'Close dialog' })).toBeInTheDocument();
	});

	it('reports onClosed after the close button closes it', async () => {
		const { user, onClosed } = setup();
		await user.click(await screen.findByRole('button', { name: 'Close dialog' }));

		await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
		await waitFor(() => expect(onClosed).toHaveBeenCalledTimes(1));
	});

	it('reports onClosed when closed from code, not only from its own controls', async () => {
		const { rerender, onClosed } = setup();
		await screen.findByRole('dialog');
		await rerender({ open: false });

		await waitFor(() => expect(onClosed).toHaveBeenCalledTimes(1));
	});

	it('reports onClosed after Escape', async () => {
		const { user, onClosed } = setup();
		await screen.findByRole('dialog');
		await user.keyboard('{Escape}');

		await waitFor(() => expect(onClosed).toHaveBeenCalledTimes(1));
	});
});
