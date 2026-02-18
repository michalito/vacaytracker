import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import StatusBadge from '../StatusBadge.svelte';

describe('StatusBadge', () => {
	it('pending status shows Awaiting', () => {
		render(StatusBadge, { status: 'pending' });
		expect(screen.getByText('Awaiting')).toBeInTheDocument();
	});

	it('approved status shows Confirmed', () => {
		render(StatusBadge, { status: 'approved' });
		expect(screen.getByText('Confirmed')).toBeInTheDocument();
	});

	it('rejected status shows Declined', () => {
		render(StatusBadge, { status: 'rejected' });
		expect(screen.getByText('Declined')).toBeInTheDocument();
	});

	it('showLabel=false hides label text', () => {
		render(StatusBadge, { status: 'pending', showLabel: false });
		expect(screen.queryByText('Awaiting')).not.toBeInTheDocument();
	});

	it('completed status shows Enjoyed', () => {
		render(StatusBadge, { status: 'completed' });
		expect(screen.getByText('Enjoyed')).toBeInTheDocument();
	});
});
