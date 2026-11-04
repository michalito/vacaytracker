import { render, screen } from '@testing-library/svelte';
import { userEvent } from '@testing-library/user-event';
import { describe, it, expect, beforeEach } from 'vitest';
import { calendar } from '$lib/stores/calendar.svelte';
import CalendarHeader from '../CalendarHeader.svelte';

function viewToggle(label: 'Week' | 'Month') {
	const button = screen.getByText(label).closest('button');
	if (!button) throw new Error(`No toggle for ${label}`);
	return button;
}

describe('CalendarHeader view toggle', () => {
	beforeEach(() => {
		calendar.setViewType('month');
	});

	it('keeps the active view selected when it is clicked again', async () => {
		const user = userEvent.setup();
		render(CalendarHeader);

		expect(viewToggle('Month')).toHaveAttribute('data-state', 'on');
		await user.click(viewToggle('Month'));

		expect(viewToggle('Month')).toHaveAttribute('data-state', 'on');
		expect(viewToggle('Week')).toHaveAttribute('data-state', 'off');
		expect(calendar.viewType).toBe('month');
	});

	it('switches the calendar store view', async () => {
		const user = userEvent.setup();
		render(CalendarHeader);

		await user.click(viewToggle('Week'));

		expect(calendar.viewType).toBe('week');
		expect(viewToggle('Week')).toHaveAttribute('data-state', 'on');
		expect(viewToggle('Month')).toHaveAttribute('data-state', 'off');
	});

	it('reflects view changes made elsewhere', async () => {
		render(CalendarHeader);

		calendar.setViewType('week');

		await expect.poll(() => viewToggle('Week').getAttribute('data-state')).toBe('on');
	});
});
