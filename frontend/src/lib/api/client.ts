import { browser } from '$app/environment';
import { get } from 'svelte/store';
import { auth } from '$lib/stores/auth';
import { toast } from '$lib/stores/notifications';

const API_BASE_URL = 'http://localhost:8080/api';

export interface ApiResponse<T = unknown> {
	success: boolean;
	data?: T;
	error?: string;
	message?: string;
}

export interface PaginatedResponse<T> {
	data: T[];
	page: number;
	per_page: number;
	total: number;
	total_pages: number;
}

class ApiClient {
	private baseURL: string;

	constructor(baseURL: string = API_BASE_URL) {
		this.baseURL = baseURL;
	}

	// Get authorization headers
	private getAuthHeaders(): Record<string, string> {
		const authState = get(auth);
		const headers: Record<string, string> = {
			'Content-Type': 'application/json'
		};

		if (authState.token) {
			headers.Authorization = `Bearer ${authState.token}`;
		}

		return headers;
	}

	private async handleResponse<T>(response: Response): Promise<ApiResponse<T>> {
		try {
			// Always read body once
			const raw = await response.text();
			let data: any;

			try {
				data = raw ? JSON.parse(raw) : {};
			} catch {
				data = { error: raw || 'Invalid JSON response' };
			}

			if (!response.ok) {
				console.error(`API request failed [${response.status} ${response.statusText}]`, data);

				if (response.status === 401) {
					if (browser) {
						auth.logout();
						toast.error('Session expired', 'Please log in again');
					}
				} else if (response.status >= 500) {
					toast.error('Server Error', 'Please try again later');
				}

				return {
					success: false,
					error: data.error || `HTTP ${response.status}: ${response.statusText}`
				};
			}

			return data as ApiResponse<T>;
		} catch (error) {
			console.error('API Response Error:', error);
			return {
				success: false,
				error: 'Failed to parse response'
			};
		}
	}

	// Generic GET request
	async get<T = unknown>(
		endpoint: string,
		params?: Record<string, unknown>
	): Promise<ApiResponse<T>> {
		try {
			const url = new URL(`${this.baseURL}${endpoint}`);

			// Add query parameters
			if (params) {
				Object.entries(params).forEach(([key, value]) => {
					if (value !== undefined && value !== null) {
						url.searchParams.append(key, String(value));
					}
				});
			}

			const response = await fetch(url.toString(), {
				method: 'GET',
				headers: this.getAuthHeaders()
			});

			return this.handleResponse<T>(response);
		} catch (error) {
			console.error('GET Error:', error);
			return {
				success: false,
				error: 'Network error'
			};
		}
	}

	// Generic POST request
	async post<T = unknown>(endpoint: string, data?: unknown): Promise<ApiResponse<T>> {
		try {
			const response = await fetch(`${this.baseURL}${endpoint}`, {
				method: 'POST',
				headers: this.getAuthHeaders(),
				body: data ? JSON.stringify(data) : undefined
			});

			return this.handleResponse<T>(response);
		} catch (error) {
			console.error('POST Error:', error);
			return {
				success: false,
				error: 'Network error'
			};
		}
	}

	// Generic PUT request
	async put<T = unknown>(endpoint: string, data?: unknown): Promise<ApiResponse<T>> {
		try {
			const response = await fetch(`${this.baseURL}${endpoint}`, {
				method: 'PUT',
				headers: this.getAuthHeaders(),
				body: data ? JSON.stringify(data) : undefined
			});

			return this.handleResponse<T>(response);
		} catch (error) {
			console.error('PUT Error:', error);
			return {
				success: false,
				error: 'Network error'
			};
		}
	}

	// Generic DELETE request
	async delete<T = unknown>(endpoint: string): Promise<ApiResponse<T>> {
		try {
			const response = await fetch(`${this.baseURL}${endpoint}`, {
				method: 'DELETE',
				headers: this.getAuthHeaders()
			});

			return this.handleResponse<T>(response);
		} catch (error) {
			console.error('DELETE Error:', error);
			return {
				success: false,
				error: 'Network error'
			};
		}
	}

	// File upload with FormData
	async uploadFile<T = unknown>(endpoint: string, formData: FormData): Promise<ApiResponse<T>> {
		try {
			const authState = get(auth);
			const headers: Record<string, string> = {};

			if (authState.token) {
				headers.Authorization = `Bearer ${authState.token}`;
			}

			const response = await fetch(`${this.baseURL}${endpoint}`, {
				method: 'POST',
				headers,
				body: formData
			});

			return this.handleResponse<T>(response);
		} catch (error) {
			console.error('Upload Error:', error);
			return {
				success: false,
				error: 'Upload failed'
			};
		}
	}
}

// Create singleton instance
export const api = new ApiClient();

// Authentication API
export const authApi = {
	async register(data: {
		username: string;
		email: string;
		password: string;
		display_name: string;
	}) {
		return api.post('/auth/register', data);
	},

	async login(data: { email: string; password: string }) {
		return api.post('/auth/login', data);
	},

	async logout() {
		return api.post('/auth/logout');
	},

	async getProfile() {
		return api.get('/auth/profile');
	},

	async updateProfile(data: {
		display_name?: string;
		bio?: string;
		notify_email?: boolean;
		notify_in_site?: boolean;
	}) {
		return api.put('/auth/profile', data);
	},

	async updatePassword(data: { current_password: string; new_password: string }) {
		return api.put('/auth/password', data);
	},

	async forgotPassword(email: string) {
		return api.post('/auth/forgot-password', { email });
	},

	async resetPassword(data: { token: string; password: string }) {
		return api.post('/auth/reset-password', data);
	},

	async getUserMods(params?: { page?: number; per_page?: number }) {
		return api.get('/auth/mods', params);
	}
};

// Games API
export const gamesApi = {
	async getGames(params?: { page?: number; per_page?: number; search?: string }) {
		return api.get('/games', params);
	},

	async getGame(slug: string) {
		return api.get(`/games/${slug}`);
	},

	async getGameById(id: number) {
		return api.get(`/games/id/${id}`);
	},

	async getGameMods(
		slug: string,
		params?: {
			page?: number;
			per_page?: number;
			sort?: string;
			tags?: string;
		}
	) {
		return api.get(`/games/${slug}/mods`, params);
	},

	async getGameTags(slug: string) {
		return api.get(`/games/${slug}/tags`);
	},

	async getGameModerators(slug: string) {
		return api.get(`/games/${slug}/moderators`);
	},

	async assignModerator(slug: string, userId: number) {
		return api.post(`/games/${slug}/moderators`, { user_id: userId });
	},

	async removeModerator(slug: string, userId: number) {
		return api.delete(`/games/${slug}/moderators/${userId}`);
	},

	async createGameRequest(data: {
		name: string;
		reason: string;
		description: string;
		icon?: string;
		existing_community?: string;
		mod_loader?: string;
		contact?: string;
	}) {
		return api.post('/games/requests', data);
	}
};

// Mods API
export const modsApi = {
	async getMod(gameSlug: string, modSlug: string) {
		return api.get(`/mods/${gameSlug}/${modSlug}`);
	},

	async getModById(id: number) {
		return api.get(`/mods/id/${id}`);
	},

	async createMod(data: {
		name: string;
		description: string;
		short_description: string;
		version: string;
		game_version: string;
		game_id: number;
		source_website?: string;
		contact_info?: string;
		tags?: string[];
		dependencies?: number[];
	}) {
		return api.post('/mods', data);
	},

	getDownloadUrl(gameSlug: string, modSlug: string): string {
		const authState = get(auth);
		const url = `${API_BASE_URL}/download/${gameSlug}/${modSlug}`;

		// If authenticated, append token as query param for download
		if (authState.token) {
			return `${url}?token=${authState.token}`;
		}

		return url;
	},

	async downloadMod(gameSlug: string, modSlug: string) {
		return api.get(`/download/${gameSlug}/${modSlug}`);
	},

	async updateMod(id: number, data: Record<string, unknown>) {
		return api.put(`/mods/${id}`, data);
	},

	async deleteMod(id: number) {
		return api.delete(`/mods/${id}`);
	},

	async likeMod(id: number) {
		return api.post(`/mods/${id}/like`);
	},

	async unlikeMod(id: number) {
		return api.delete(`/mods/${id}/like`);
	},

	async uploadModFile(modId: number, file: File, isMain: boolean = false) {
		const formData = new FormData();
		formData.append('file', file);
		formData.append('is_main', String(isMain));
		return api.uploadFile(`/mods/${modId}/files`, formData);
	},

	async approveMod(id: number) {
		return api.post(`/mods/${id}/approve`);
	},

	async rejectMod(id: number, reason: string) {
		return api.post(`/mods/${id}/reject`, { reason });
	}
};

// Comments API
export const commentsApi = {
	async getModComments(modId: number, params?: { page?: number; per_page?: number }) {
		return api.get(`/comments/mod/${modId}`, params);
	},

	async createComment(modId: number, data: { content: string; parent_id?: number }) {
		return api.post(`/comments?mod_id=${modId}`, data);
	},

	async updateComment(id: number, data: { content: string }) {
		return api.put(`/comments/${id}`, data);
	},

	async deleteComment(id: number) {
		return api.delete(`/comments/${id}`);
	}
};

// Notifications API
export const notificationsApi = {
	async getNotifications(params?: { page?: number; per_page?: number }) {
		return api.get('/notifications', params);
	},

	async getUnreadCount() {
		return api.get('/notifications/unread-count');
	},

	async markAsRead(id: number) {
		return api.put(`/notifications/${id}/read`);
	},

	async markAllAsRead() {
		return api.put('/notifications/read-all');
	},

	async deleteNotification(id: number) {
		return api.delete(`/notifications/${id}`);
	}
};

// Documentation API
export const docsApi = {
	async getGameDocs(gameSlug: string, params?: { page?: number; per_page?: number }) {
		return api.get(`/docs/${gameSlug}`, params);
	},

	async getDoc(gameSlug: string, docSlug: string) {
		return api.get(`/docs/${gameSlug}/${docSlug}`);
	},

	async createDoc(gameSlug: string, data: { title: string; content: string }) {
		return api.post(`/docs/${gameSlug}`, data);
	},

	async updateDoc(gameSlug: string, docSlug: string, data: { title?: string; content?: string }) {
		return api.put(`/docs/${gameSlug}/${docSlug}`, data);
	},

	async deleteDoc(gameSlug: string, docSlug: string) {
		return api.delete(`/docs/${gameSlug}/${docSlug}`);
	}
};

// Users API
export const usersApi = {
	async getUser(id: number) {
		return api.get(`/users/${id}`);
	},

	async getUserByUsername(username: string) {
		return api.get(`/users/username/${username}`);
	},

	async getUserMods(userId: number, params?: { page?: number; per_page?: number }) {
		return api.get(`/users/${userId}/mods`, params);
	}
};

// Admin API
export const adminApi = {
	async getStats() {
		return api.get('/admin/stats');
	},

	async getActivity(params?: { limit?: number }) {
		return api.get('/admin/activity', params);
	},

	async getPendingMods(params?: { page?: number; per_page?: number }) {
		return api.get('/admin/mods/pending', params);
	},

	async createBan(data: {
		user_id?: number;
		ip_address?: string;
		game_id?: number;
		reason: string;
		duration?: number;
	}) {
		return api.post('/admin/bans', data);
	},

	async getBans(params?: { page?: number; per_page?: number; active?: boolean }) {
		return api.get('/admin/bans', params);
	},

	async unbanUser(banId: number) {
		return api.post(`/admin/bans/${banId}/unban`);
	},

	async getGameRequests(params?: { page?: number; per_page?: number }) {
		return api.get('/games/requests', params);
	},

	async createGame(data: Record<string, unknown>) {
		return api.post('/games', data);
	},

	async updateGame(id: number, data: Record<string, unknown>) {
		return api.put(`/games/manage/${id}`, data);
	},

	async deleteGame(id: number) {
		return api.delete(`/games/manage/${id}`);
	},

	async approveGameRequest(requestId: number, adminNotes?: string) {
		return api.post(`/games/requests/${requestId}/approve`, { admin_notes: adminNotes });
	},

	async rejectGameRequest(requestId: number, adminNotes?: string) {
		return api.post(`/games/requests/${requestId}/deny`, { admin_notes: adminNotes });
	},

	async getAllGames(params?: { page?: number; per_page?: number }) {
		return api.get('/admin/games', params);
	},

	async getGameRequest(requestId: number) {
		return api.get(`/games/requests/${requestId}`);
	},

	async updateGameRequest(requestId: number, data: Record<string, unknown>) {
		return api.put(`/games/requests/${requestId}`, data);
	},

	async approveMod(modId: number) {
		return api.post(`/mods/${modId}/approve`);
	},

	async rejectMod(modId: number, reason: string) {
		return api.post(`/mods/${modId}/reject`, { reason });
	}
};

// Health check
export const healthApi = {
	async check() {
		return api.get('/health');
	}
};
