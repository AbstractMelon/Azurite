// User types
export interface User {
	id: number;
	username: string;
	email: string;
	display_name: string;
	avatar: string;
	bio: string;
	role: string;
	is_active: boolean;
	email_verified: boolean;
	created_at: string;
	updated_at: string;
	last_login_at?: string;
	notify_email: boolean;
	notify_in_site: boolean;
	github_id?: string;
	discord_id?: string;
	google_id?: string;
}

// Game types
export interface Game {
	id: number;
	name: string;
	slug: string;
	description: string;
	icon: string;
	is_active: boolean;
	created_at: string;
	updated_at: string;
	mod_count: number;
}

// Tag types
export interface Tag {
	id: number;
	name: string;
	slug: string;
	game_id: number;
	created_at: string;
}

// ModFile types
export interface ModFile {
	id: number;
	mod_id: number;
	filename: string;
	file_path: string;
	file_size: number;
	mime_type: string;
	hash: string;
	is_main: boolean;
	created_at: string;
}

// Mod types
export interface Mod {
	id: number;
	name: string;
	slug: string;
	description: string;
	short_description: string;
	icon: string;
	version: string;
	game_version: string;
	game_id: number;
	owner_id: number;
	downloads: number;
	likes: number;
	source_website: string;
	contact_info: string;
	is_rejected: boolean;
	rejection_reason: string;
	is_scanned: boolean;
	scan_result: string;
	created_at: string;
	updated_at: string;
	game?: Game;
	owner?: User;
	tags?: Tag[];
	dependencies?: Mod[];
	files?: ModFile[];
	is_liked?: boolean;
	comments?: number;
}

// Comment types
export interface Comment {
	id: number;
	mod_id: number;
	user_id: number;
	content: string;
	parent_id?: number;
	is_active: boolean;
	created_at: string;
	updated_at: string;
	user?: User;
	replies?: Comment[];
}

// Game Request types
export interface GameRequest {
	id: number;
	name: string;
	description: string;
	icon: string;
	requested_by: number;
	status: 'pending' | 'approved' | 'denied';
	admin_notes: string;
	created_at: string;
	updated_at: string;
	user?: User;
}

// Documentation types
export interface Documentation {
	id: number;
	game_id: number;
	title: string;
	slug: string;
	content: string;
	author_id: number;
	created_at: string;
	updated_at: string;
	author?: User;
	game?: Game;
}

// Notification types
export interface Notification {
	id: number;
	user_id: number;
	type: string;
	title: string;
	message: string;
	data: string;
	is_read: boolean;
	created_at: string;
}

// Ban types
export interface Ban {
	id: number;
	user_id?: number;
	ip_address: string;
	game_id?: number;
	reason: string;
	banned_by: number;
	expires_at?: string;
	is_active: boolean;
	created_at: string;
	user?: User;
	banned_by_user?: User;
	game?: Game;
}

// API Response types
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

// Request types
export interface LoginRequest {
	email: string;
	password: string;
}

export interface RegisterRequest {
	username: string;
	email: string;
	password: string;
	display_name: string;
}

export interface ModCreateRequest {
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
}

export interface ModUpdateRequest {
	name: string;
	description: string;
	short_description: string;
	version: string;
	game_version: string;
	source_website?: string;
	contact_info?: string;
	tags?: string[];
	dependencies?: number[];
}

export interface CommentCreateRequest {
	content: string;
	parent_id?: number;
}

export interface UserUpdateRequest {
	display_name: string;
	bio?: string;
	notify_email?: boolean;
	notify_in_site?: boolean;
}

export interface GameRequestCreateRequest {
	name: string;
	description: string;
}

export interface BanCreateRequest {
	user_id?: number;
	ip_address?: string;
	game_id?: number;
	reason: string;
	duration?: number;
}

// Auth types
export interface AuthResponse {
	token: string;
	user: User;
}

export interface AuthState {
	user: User | null;
	token: string | null;
	isAuthenticated: boolean;
}

// Dashboard types
export interface DashboardStats {
	totalMods: number;
	totalDownloads: number;
	totalLikes: number;
	totalComments: number;
	pendingMods: number;
	approvedMods: number;
	rejectedMods: number;
}

export interface AdminStats {
	totalUsers: number;
	totalMods: number;
	totalGames: number;
	totalDownloads: number;
	pendingMods: number;
	totalComments: number;
	totalBans: number;
	recentActivity: ActivityItem[];
}

export interface ActivityItem {
	id: number;
	type: string;
	name: string;
	title: string;
	message: string;
	created_at: string;
	user?: User;
	mod?: Mod;
	game?: Game;
}

// Form types
export interface FormField {
	name: string;
	label: string;
	type: 'text' | 'email' | 'password' | 'textarea' | 'select' | 'checkbox' | 'file';
	value: string | boolean | File | null;
	error?: string;
	required?: boolean;
	placeholder?: string;
	options?: Array<{ value: string | number; label: string }>;
}

// UI types
export interface Toast {
	id: string;
	type: 'success' | 'error' | 'warning' | 'info';
	title: string;
	message?: string;
	duration?: number;
	dismissible?: boolean;
}

export interface ModalConfig {
	title: string;
	message?: string;
	confirmText?: string;
	cancelText?: string;
	type?: 'confirm' | 'alert' | 'prompt';
	onConfirm?: () => void | Promise<void>;
	onCancel?: () => void;
}

// Component props types
export interface LoadingProps {
	size?: 'sm' | 'md' | 'lg';
	text?: string;
	inline?: boolean;
}

export interface ButtonProps {
	variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
	size?: 'sm' | 'md' | 'lg';
	disabled?: boolean;
	loading?: boolean;
	type?: 'button' | 'submit' | 'reset';
	href?: string;
	target?: '_blank' | '_self' | '_parent' | '_top';
}

export interface CardProps {
	title?: string;
	subtitle?: string;
	hover?: boolean;
	padding?: 'none' | 'sm' | 'md' | 'lg';
}

// Search and filter types
export interface SearchFilters {
	query?: string;
	game_id?: number;
	tags?: string[];
	sort?: 'newest' | 'oldest' | 'popular' | 'downloads' | 'likes' | 'name';
	page?: number;
	per_page?: number;
}

export interface ModFilters extends SearchFilters {
	status?: 'all' | 'pending' | 'approved' | 'rejected';
	owner_id?: number;
}

// Pagination types
export interface PaginationInfo {
	page: number;
	per_page: number;
	total: number;
	total_pages: number;
	has_next: boolean;
	has_prev: boolean;
}

// File upload types
export interface FileUpload {
	file: File;
	progress: number;
	status: 'pending' | 'uploading' | 'success' | 'error';
	error?: string;
	url?: string;
}

// Constants
export const USER_ROLES = {
	USER: 'user',
	ADMIN: 'admin',
	COMMUNITY_MODERATOR: 'community_moderator',
	WIKI_MAINTAINER: 'wiki_maintainer'
} as const;

export const GAME_REQUEST_STATUS = {
	PENDING: 'pending',
	APPROVED: 'approved',
	DENIED: 'denied'
} as const;

export const NOTIFICATION_TYPES = {
	MOD_REJECTED: 'mod_rejected',
	NEW_COMMENT: 'new_comment',
	MOD_MILESTONE: 'mod_milestone',
	GAME_REQUEST: 'game_request',
	MOD_APPROVED: 'mod_approved'
} as const;

export const SCAN_RESULTS = {
	PENDING: 'pending',
	CLEAN: 'clean',
	THREAT: 'threat'
} as const;

// Type guards
export function isUser(obj: unknown): obj is User {
	return typeof obj === 'object' && obj !== null && 'id' in obj && 'username' in obj;
}

export function isMod(obj: unknown): obj is Mod {
	return typeof obj === 'object' && obj !== null && 'id' in obj && 'name' in obj && 'slug' in obj;
}

export function isGame(obj: unknown): obj is Game {
	return typeof obj === 'object' && obj !== null && 'id' in obj && 'name' in obj && 'slug' in obj;
}

export function isComment(obj: unknown): obj is Comment {
	return (
		typeof obj === 'object' && obj !== null && 'id' in obj && 'content' in obj && 'mod_id' in obj
	);
}

// Utility types
export type UserRole = (typeof USER_ROLES)[keyof typeof USER_ROLES];
export type GameRequestStatus = (typeof GAME_REQUEST_STATUS)[keyof typeof GAME_REQUEST_STATUS];
export type NotificationType = (typeof NOTIFICATION_TYPES)[keyof typeof NOTIFICATION_TYPES];
export type ScanResult = (typeof SCAN_RESULTS)[keyof typeof SCAN_RESULTS];

// Error types
export interface ApiError {
	message: string;
	code?: string;
	status?: number;
	details?: Record<string, unknown>;
}

export interface ValidationError {
	field: string;
	message: string;
	code: string;
}

// Theme types
export interface Theme {
	name: string;
	primary: string;
	secondary: string;
	background: string;
	surface: string;
	text: string;
	textSecondary: string;
	border: string;
}

// Route types for type-safe navigation
export interface RouteParams {
	slug?: string;
	id?: string;
	gameSlug?: string;
	modSlug?: string;
}

// Component event types
export interface CustomEvents {
	click: MouseEvent;
	submit: SubmitEvent;
	change: Event;
	input: Event;
	keydown: KeyboardEvent;
	keyup: KeyboardEvent;
}
