export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface Toast {
	id: string;
	type: ToastType;
	message: string;
	duration?: number;
}

function createToastStore() {
	let toasts = $state<Toast[]>([]);

	function add(type: ToastType, message: string, duration: number = 5000): string {
		const id = crypto.randomUUID();
		const toast: Toast = { id, type, message, duration };

		toasts = [...toasts, toast];

		if (duration > 0) {
			setTimeout(() => dismiss(id), duration);
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
	const success = (message: string, duration?: number) => add('success', message, duration);
	const error = (message: string, duration?: number) => add('error', message, duration);
	const warning = (message: string, duration?: number) => add('warning', message, duration);
	const info = (message: string, duration?: number) => add('info', message, duration);

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
