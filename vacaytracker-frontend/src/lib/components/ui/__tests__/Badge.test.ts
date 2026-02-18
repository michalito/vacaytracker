import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import BadgeWrapper from './BadgeWrapper.svelte';

describe('Badge', () => {
	it('renders badge text', () => {
		render(BadgeWrapper, { text: 'Active' });
		expect(screen.getByText('Active')).toBeInTheDocument();
	});

	it('applies default variant classes', () => {
		render(BadgeWrapper, { text: 'Default' });
		const badge = screen.getByText('Default');
		expect(badge.className).toContain('bg-slate-100');
	});

	it('applies success variant classes', () => {
		render(BadgeWrapper, { text: 'Success', variant: 'success' });
		const badge = screen.getByText('Success');
		expect(badge.className).toContain('bg-success-light');
	});

	it('applies error variant classes', () => {
		render(BadgeWrapper, { text: 'Error', variant: 'error' });
		const badge = screen.getByText('Error');
		expect(badge.className).toContain('bg-error-light');
	});

	it('applies sm size classes', () => {
		render(BadgeWrapper, { text: 'Small', size: 'sm' });
		const badge = screen.getByText('Small');
		expect(badge.className).toContain('text-xs');
	});
});
