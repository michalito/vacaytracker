import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import type { TeamVacation } from '$lib/types';
import VacationEvent from '$lib/components/features/calendar/VacationEvent.svelte';
import VacationBalanceCard from '$lib/components/features/dashboard/VacationBalanceCard.svelte';
import TooltipHarness from './TooltipHarness.svelte';

// The Tooltip wrapper hands its trigger props to the consumer's own element
// instead of wrapping it in another <button>.
describe('Tooltip trigger composition', () => {
	it('does not nest the info button inside another button', () => {
		const { container } = render(TooltipHarness, {
			props: { component: VacationBalanceCard, props: { total: 25, used: 5, upcoming: 3 } }
		});

		expect(container.querySelector('button button')).toBeNull();
		const info = screen.getByRole('button', { name: 'About your vacation balance' });
		expect(info).toHaveAttribute('data-tooltip-trigger');
	});

	it('makes the calendar event itself the tooltip trigger', () => {
		const vacation: TeamVacation = {
			id: 'v1',
			userId: 'u1',
			userName: 'Alice Smith',
			startDate: '2026-06-01',
			endDate: '2026-06-05',
			totalDays: 5
		};
		const { container } = render(TooltipHarness, {
			props: { component: VacationEvent, props: { vacation } }
		});

		expect(container.querySelector('button')).toBeNull();
		const trigger = container.querySelector('[data-tooltip-trigger]');
		expect(trigger?.tagName).toBe('DIV');
		expect(trigger).toHaveTextContent('Alice Smith');
		expect(trigger).not.toHaveAttribute('type');
	});
});
