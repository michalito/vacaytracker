import { render, screen } from '@testing-library/svelte';
import { userEvent } from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import ButtonWrapper from './ButtonWrapper.svelte';

describe('Button', () => {
	it('renders with text', () => {
		render(ButtonWrapper, { text: 'Submit' });
		expect(screen.getByRole('button', { name: 'Submit' })).toBeInTheDocument();
	});

	it('applies primary variant by default', () => {
		render(ButtonWrapper, { text: 'Primary' });
		const button = screen.getByRole('button');
		expect(button.className).toContain('bg-ocean-500');
	});

	it('applies secondary variant', () => {
		render(ButtonWrapper, { text: 'Secondary', variant: 'secondary' });
		const button = screen.getByRole('button');
		expect(button.className).toContain('bg-sand-200');
	});

	it('applies size classes', () => {
		render(ButtonWrapper, { text: 'Large', size: 'lg' });
		const button = screen.getByRole('button');
		expect(button.className).toContain('px-6');
	});

	it('disabled state', () => {
		render(ButtonWrapper, { text: 'Disabled', disabled: true });
		const button = screen.getByRole('button');
		expect(button).toBeDisabled();
	});

	it('loading state disables button and shows spinner', () => {
		render(ButtonWrapper, { text: 'Loading', loading: true });
		const button = screen.getByRole('button');
		expect(button).toBeDisabled();
		expect(button).toHaveAttribute('aria-busy', 'true');
		expect(button.querySelector('svg')).toBeInTheDocument();
	});

	it('calls onclick handler when clicked', async () => {
		const user = userEvent.setup();
		const handler = vi.fn();
		render(ButtonWrapper, { text: 'Click', onclick: handler });
		await user.click(screen.getByRole('button'));
		expect(handler).toHaveBeenCalledOnce();
	});

	it('disabled button does not call onclick handler', async () => {
		const user = userEvent.setup();
		const handler = vi.fn();
		render(ButtonWrapper, { text: 'No click', disabled: true, onclick: handler });
		await user.click(screen.getByRole('button'));
		expect(handler).not.toHaveBeenCalled();
	});
});
