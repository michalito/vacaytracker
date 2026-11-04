import { render, screen, waitFor } from '@testing-library/svelte';
import { userEvent } from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import Select from '../Select.svelte';

const options = [
	{ value: 'employee', label: 'Employee' },
	{ value: 'admin', label: 'Admin' }
];

describe('Select', () => {
	it('names the trigger by its visible label', () => {
		render(Select, { props: { label: 'Role', options, value: 'admin' } });

		const trigger = screen.getByRole('combobox', { name: 'Role' });
		expect(trigger).toHaveTextContent('Admin');
	});

	it('uses aria-label when there is no visible label', () => {
		render(Select, { props: { 'aria-label': 'Filter by role', options, value: 'employee' } });

		expect(screen.getByRole('combobox', { name: 'Filter by role' })).toBeInTheDocument();
	});

	it('reports every selection through onchange and updates the trigger', async () => {
		const user = userEvent.setup();
		const onchange = vi.fn();
		// An empty-string value is the case the old `if (v)` guard silently dropped
		const roleOptions = [
			{ value: '', label: 'Any Role' },
			{ value: 'admin', label: 'Admin' }
		];
		render(Select, {
			props: { 'aria-label': 'Filter by role', options: roleOptions, value: '', onchange }
		});
		const trigger = screen.getByRole('combobox', { name: 'Filter by role' });

		await user.click(trigger);
		await user.click(await screen.findByText('Admin'));
		await waitFor(() => expect(trigger).toHaveTextContent('Admin'));

		await user.click(trigger);
		await user.click(await screen.findByText('Any Role'));
		await waitFor(() => expect(trigger).toHaveTextContent('Any Role'));

		expect(onchange.mock.calls).toEqual([['admin'], ['']]);
	});

	it('gives each instance a unique trigger id', () => {
		render(Select, { props: { label: 'Role', options } });
		render(Select, { props: { label: 'Frequency', options } });

		const ids = screen.getAllByRole('combobox').map((el) => el.id);
		expect(new Set(ids).size).toBe(2);
	});
});
