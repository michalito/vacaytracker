/**
 * Toast notification store using Svelte 5 runes.
 * Provides a simple API for displaying notifications with auto-dismiss.
 */

export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface ToastData {
	id: string;
	type: ToastType;
	title: string;
	description?: string;
}

function createToastStore() {
	let toasts = $state<ToastData[]>([]);

	/**
	 * Add a toast notification.
	 * Supports multiple signatures for backward compatibility:
	 * - add(type, title)
	 * - add(type, title, duration)
	 * - add(type, title, description)
	 * - add(type, title, description, duration)
	 */
	function add(
		type: ToastType,
		title: string,
		descOrDuration?: string | number,
		duration?: number
	): string {
		const id = crypto.randomUUID();

		let description: string | undefined;
		let closeDelay = 5000;

		if (typeof descOrDuration === 'number') {
			closeDelay = descOrDuration;
		} else if (typeof descOrDuration === 'string') {
			description = descOrDuration;
			closeDelay = duration ?? 5000;
		}

		const toast: ToastData = { id, type, title, description };
		toasts = [...toasts, toast];

		if (closeDelay > 0) {
			setTimeout(() => dismiss(id), closeDelay);
		}

		return id;
	}

	function dismiss(id: string): void {
		toasts = toasts.filter((t) => t.id !== id);
	}

	function dismissAll(): void {
		toasts = [];
	}

	// Convenience methods
	const success = (title: string, descOrDuration?: string | number, duration?: number) =>
		add('success', title, descOrDuration, duration);

	const error = (title: string, descOrDuration?: string | number, duration?: number) =>
		add('error', title, descOrDuration, duration);

	const warning = (title: string, descOrDuration?: string | number, duration?: number) =>
		add('warning', title, descOrDuration, duration);

	const info = (title: string, descOrDuration?: string | number, duration?: number) =>
		add('info', title, descOrDuration, duration);

	return {
		get toasts() {
			return toasts;
		},
		add,
		dismiss,
		dismissAll,
		success,
		error,
		warning,
		info
	};
}

export const toast = createToastStore();
