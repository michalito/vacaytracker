import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { toast } from '$lib/stores/toast.svelte';

beforeEach(() => {
	toast.dismissAll();
	vi.useFakeTimers();
});

afterEach(() => {
	vi.useRealTimers();
});

describe('toast store', () => {
	describe('add', () => {
		it('adds a toast with correct type and title, and returns an id', () => {
			const id = toast.add('info', 'Hello');

			expect(id).toBeDefined();
			expect(typeof id).toBe('string');
			expect(toast.toasts).toHaveLength(1);
			expect(toast.toasts[0]).toMatchObject({
				id,
				type: 'info',
				title: 'Hello'
			});
		});

		it('sets description when third argument is a string', () => {
			toast.add('success', 'Title', 'Some description');

			expect(toast.toasts[0].description).toBe('Some description');
		});

		it('sets custom duration when third argument is a number (no description)', () => {
			toast.add('warning', 'Quick toast', 2000);

			expect(toast.toasts).toHaveLength(1);
			expect(toast.toasts[0].description).toBeUndefined();

			// Should still be present before 2000ms
			vi.advanceTimersByTime(1999);
			expect(toast.toasts).toHaveLength(1);

			// Should be dismissed at 2000ms
			vi.advanceTimersByTime(1);
			expect(toast.toasts).toHaveLength(0);
		});

		it('sets both description and custom duration when both provided', () => {
			toast.add('error', 'Oops', 'Something went wrong', 3000);

			expect(toast.toasts[0]).toMatchObject({
				type: 'error',
				title: 'Oops',
				description: 'Something went wrong'
			});

			vi.advanceTimersByTime(2999);
			expect(toast.toasts).toHaveLength(1);

			vi.advanceTimersByTime(1);
			expect(toast.toasts).toHaveLength(0);
		});
	});

	describe('dismiss', () => {
		it('removes a specific toast by id', () => {
			const id1 = toast.add('info', 'First');
			toast.add('info', 'Second');

			expect(toast.toasts).toHaveLength(2);

			toast.dismiss(id1);

			expect(toast.toasts).toHaveLength(1);
			expect(toast.toasts[0].title).toBe('Second');
		});

		it('does not throw when dismissing a non-existent id', () => {
			toast.add('info', 'Keep me');

			expect(() => toast.dismiss('non-existent-id')).not.toThrow();
			expect(toast.toasts).toHaveLength(1);
			expect(toast.toasts[0].title).toBe('Keep me');
		});
	});

	describe('dismissAll', () => {
		it('clears all toasts', () => {
			toast.add('info', 'One');
			toast.add('warning', 'Two');
			toast.add('error', 'Three');

			expect(toast.toasts).toHaveLength(3);

			toast.dismissAll();

			expect(toast.toasts).toHaveLength(0);
		});
	});

	describe('convenience methods', () => {
		it('success adds a toast with type "success"', () => {
			const id = toast.success('Done!');

			expect(toast.toasts[0]).toMatchObject({
				id,
				type: 'success',
				title: 'Done!'
			});
		});

		it('error adds a toast with type "error"', () => {
			const id = toast.error('Failed');

			expect(toast.toasts[0]).toMatchObject({
				id,
				type: 'error',
				title: 'Failed'
			});
		});

		it('warning adds a toast with type "warning"', () => {
			const id = toast.warning('Watch out');

			expect(toast.toasts[0]).toMatchObject({
				id,
				type: 'warning',
				title: 'Watch out'
			});
		});

		it('info adds a toast with type "info"', () => {
			const id = toast.info('FYI');

			expect(toast.toasts[0]).toMatchObject({
				id,
				type: 'info',
				title: 'FYI'
			});
		});
	});

	describe('auto-dismiss', () => {
		it('removes toast after default 5000ms', () => {
			toast.add('info', 'Auto dismiss');

			expect(toast.toasts).toHaveLength(1);

			vi.advanceTimersByTime(4999);
			expect(toast.toasts).toHaveLength(1);

			vi.advanceTimersByTime(1);
			expect(toast.toasts).toHaveLength(0);
		});

		it('removes toast after custom duration', () => {
			toast.add('info', 'Quick', 1500);

			vi.advanceTimersByTime(1499);
			expect(toast.toasts).toHaveLength(1);

			vi.advanceTimersByTime(1);
			expect(toast.toasts).toHaveLength(0);
		});

		it('uses default 5000ms when description is provided without custom duration', () => {
			toast.add('info', 'Title', 'A description');

			vi.advanceTimersByTime(4999);
			expect(toast.toasts).toHaveLength(1);

			vi.advanceTimersByTime(1);
			expect(toast.toasts).toHaveLength(0);
		});
	});

	describe('multiple toasts', () => {
		it('can add multiple toasts, each with a unique id', () => {
			const id1 = toast.add('info', 'First');
			const id2 = toast.add('success', 'Second');
			const id3 = toast.add('error', 'Third');

			expect(toast.toasts).toHaveLength(3);
			expect(id1).not.toBe(id2);
			expect(id2).not.toBe(id3);
			expect(id1).not.toBe(id3);

			expect(toast.toasts[0].title).toBe('First');
			expect(toast.toasts[1].title).toBe('Second');
			expect(toast.toasts[2].title).toBe('Third');
		});
	});
});
