import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export interface Notification {
	id: string;
	type: 'success' | 'error' | 'warning' | 'info';
	title: string;
	message?: string;
	duration?: number;
	persistent?: boolean;
	actions?: {
		label: string;
		action: () => void;
	}[];
}

export interface NotificationState {
	notifications: Notification[];
}

const initialState: NotificationState = {
	notifications: []
};

function createNotificationStore() {
	const { subscribe, update } = writable<NotificationState>(initialState);

	const store = {
		subscribe,
		// Add a new notification
		add: (notification: Omit<Notification, 'id'>) => {
			const id = crypto.randomUUID();
			const newNotification: Notification = {
				id,
				duration: 5000, // 5 seconds default
				persistent: false,
				...notification
			};

			update((state) => ({
				...state,
				notifications: [...state.notifications, newNotification]
			}));

			// Auto-remove non-persistent notifications
			if (!newNotification.persistent && browser) {
				setTimeout(() => {
					store.remove(id); // Use store.remove instead of just remove
				}, newNotification.duration);
			}

			return id;
		},
		// Remove a notification by id
		remove: (id: string) => {
			update((state) => ({
				...state,
				notifications: state.notifications.filter((n) => n.id !== id)
			}));
		},
		// Clear all notifications
		clear: () => {
			update(() => initialState);
		},
		// Convenience methods for different notification types
		success: (title: string, message?: string, options?: Partial<Notification>) => {
			return store.add({
				// Use store.add instead of createNotificationStore().add
				type: 'success',
				title,
				message,
				...options
			});
		},
		error: (title: string, message?: string, options?: Partial<Notification>) => {
			return store.add({
				// Use store.add instead of createNotificationStore().add
				type: 'error',
				title,
				message,
				duration: 8000, // Errors stay longer
				...options
			});
		},
		warning: (title: string, message?: string, options?: Partial<Notification>) => {
			return store.add({
				// Use store.add instead of createNotificationStore().add
				type: 'warning',
				title,
				message,
				duration: 6000,
				...options
			});
		},
		info: (title: string, message?: string, options?: Partial<Notification>) => {
			return store.add({
				// Use store.add instead of createNotificationStore().add
				type: 'info',
				title,
				message,
				...options
			});
		}
	};

	return store;
}

export const notifications = createNotificationStore();

// Helper functions for quick notifications
export const toast = {
	success: (title: string, message?: string) =>
		notifications.add({ type: 'success', title, message }),
	error: (title: string, message?: string) =>
		notifications.add({ type: 'error', title, message, duration: 8000 }),
	warning: (title: string, message?: string) =>
		notifications.add({ type: 'warning', title, message, duration: 6000 }),
	info: (title: string, message?: string) => notifications.add({ type: 'info', title, message })
};
