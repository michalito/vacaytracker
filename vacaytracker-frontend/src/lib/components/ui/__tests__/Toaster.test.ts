import { render, screen } from '@testing-library/svelte';
import { userEvent } from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { toast } from '$lib/stores/toast.svelte';
import Toaster from '../Toaster.svelte';

describe('Toaster', () => {
	beforeEach(() => {
		vi.useFakeTimers();
	});

	afterEach(() => {
		toast.dismissAll();
		vi.useRealTimers();
	});

	it('renders empty when no toasts', () => {
		render(Toaster);
		expect(screen.queryAllByRole('alert')).toHaveLength(0);
	});

	it('renders toast with title', () => {
		toast.add('success', 'Done!');
		render(Toaster);
		expect(screen.getByText('Done!')).toBeInTheDocument();
	});

	it('renders toast with description', () => {
		toast.add('error', 'Error', 'Something failed');
		render(Toaster);
		expect(screen.getByText('Error')).toBeInTheDocument();
		expect(screen.getByText('Something failed')).toBeInTheDocument();
	});

	it('dismiss button removes toast', async () => {
		const user = userEvent.setup({ advanceTimers: vi.advanceTimersByTime });
		toast.add('info', 'Dismissable');
		render(Toaster);
		expect(screen.getByText('Dismissable')).toBeInTheDocument();
		await user.click(screen.getByLabelText('Dismiss notification'));
		expect(screen.queryByText('Dismissable')).not.toBeInTheDocument();
	});

	it('renders multiple toasts', () => {
		toast.add('success', 'First toast');
		toast.add('error', 'Second toast');
		render(Toaster);
		expect(screen.getByText('First toast')).toBeInTheDocument();
		expect(screen.getByText('Second toast')).toBeInTheDocument();
		expect(screen.getAllByRole('alert')).toHaveLength(2);
	});
});
