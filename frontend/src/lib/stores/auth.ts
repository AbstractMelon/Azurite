import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

export interface User {
	id: number;
	username: string;
	email: string;
	display_name: string;
	avatar?: string;
	bio?: string;
	role: 'user' | 'admin' | 'community_moderator' | 'wiki_maintainer';
	is_active: boolean;
	email_verified: boolean;
	notify_email: boolean;
	notify_in_site: boolean;
	created_at: string;
	updated_at: string;
	last_login_at?: string;
}

export interface AuthState {
	user: User | null;
	token: string | null;
	isLoading: boolean;
	error: string | null;
}

// Initial state
const initialState: AuthState = {
	user: null,
	token: null,
	isLoading: false,
	error: null
};

// Create the auth store
function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>(initialState);

	// Load auth data from localStorage on initialization
	if (browser) {
		const stored = localStorage.getItem('azurite_auth');
		if (stored) {
			try {
				const authData = JSON.parse(stored);
				set({
					...initialState,
					user: authData.user,
					token: authData.token
				});
			} catch (e) {
				console.error('Failed to parse stored auth data:', e);
				localStorage.removeItem('azurite_auth');
			}
		}
	}

	return {
		subscribe,

		// Set loading state
		setLoading: (loading: boolean) => {
			update((state) => ({
				...state,
				isLoading: loading,
				error: loading ? null : state.error
			}));
		},

		// Set error
		setError: (error: string | null) => {
			update((state) => ({
				...state,
				error,
				isLoading: false
			}));
		},

		// Login user
		login: (user: User, token: string) => {
			const authData = { user, token };

			// Update store
			set({
				user,
				token,
				isLoading: false,
				error: null
			});

			// Persist to localStorage
			if (browser) {
				localStorage.setItem('azurite_auth', JSON.stringify(authData));
			}
		},

		// Update user profile
		updateUser: (updatedUser: Partial<User>) => {
			update((state) => {
				if (!state.user) return state;

				const newUser = { ...state.user, ...updatedUser };
				const authData = { user: newUser, token: state.token };

				// Persist to localStorage
				if (browser) {
					localStorage.setItem('azurite_auth', JSON.stringify(authData));
				}

				return {
					...state,
					user: newUser
				};
			});
		},

		// Logout user
		logout: () => {
			set(initialState);

			// Clear localStorage
			if (browser) {
				localStorage.removeItem('azurite_auth');
			}
		},

		// Clear error
		clearError: () => {
			update((state) => ({
				...state,
				error: null
			}));
		},

		// Refresh token (placeholder for future implementation)
		refreshToken: async () => {
			// TODO: Implement token refresh logic
			return true;
		}
	};
}

// Export the auth store
export const auth = createAuthStore();

// Derived stores for convenience
export const user = derived(auth, ($auth) => $auth.user);
export const isAuthenticated = derived(auth, ($auth) => !!$auth.user && !!$auth.token);
export const isLoading = derived(auth, ($auth) => $auth.isLoading);
export const authError = derived(auth, ($auth) => $auth.error);
export const token = derived(auth, ($auth) => $auth.token);

// Role-based derived stores
export const isAdmin = derived(user, ($user) => $user?.role === 'admin');
export const isModerator = derived(
	user,
	($user) => $user?.role === 'admin' || $user?.role === 'community_moderator'
);
export const isWikiMaintainer = derived(
	user,
	($user) =>
		$user?.role === 'admin' ||
		$user?.role === 'community_moderator' ||
		$user?.role === 'wiki_maintainer'
);

// Helper function to check if user has permission for a specific action
export function hasPermission(userRole: string | undefined, requiredRoles: string[]): boolean {
	if (!userRole) return false;

	// Admin has all permissions
	if (userRole === 'admin') return true;

	// Check if user role is in required roles
	return requiredRoles.includes(userRole);
}
