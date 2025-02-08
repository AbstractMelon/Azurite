import { BaseEntity } from "./repositories";

export enum UserRole {
  USER = "user",
  MOD_CREATOR = "mod_creator",
  ADMIN = "admin",
}

export interface User extends BaseEntity {
  username: string;
  email: string;
  password: string; // Hashed password
  role: UserRole;
  displayName?: string;
  avatarUrl?: string;
  bio?: string;
  mods: string[]; // Array of mod IDs created by the user
  favorites: string[]; // Array of favorite mod IDs
}

export interface Rating {
  userId: string;
  score: number; // 1-5 stars
  comment?: string;
  createdAt: Date;
}

export interface Comment {
  id: string;
  userId: string;
  content: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface ChangelogEntry {
  version: string;
  description: string;
  changes: string[];
  date: Date;
}

export interface Screenshot {
  id: string;
  url: string;
  caption?: string;
  order: number;
}

export interface Mod extends BaseEntity {
  name: string;
  description: string;
  shortDescription: string;
  version: string;
  gameId: string;
  creatorId: string;
  filePath: string;
  fileSize: number;
  tags: string[];
  downloads: number;
  ratings: Rating[];
  averageRating: number;
  comments: Comment[];
  changelog: ChangelogEntry[];
  screenshots: Screenshot[];
  requirements?: {
    minGameVersion?: string;
    dependencies?: string[]; // Array of required mod IDs
  };
  isPublished: boolean;
  isApproved: boolean;
  approvedBy?: string; // Admin user ID
  approvedAt?: Date;
}

export interface Game extends BaseEntity {
  name: string;
  description: string;
  shortDescription: string;
  coverImageUrl?: string;
  websiteUrl?: string;
  supportedModTypes: string[];
  modCount: number;
  latestVersion: string;
  supportedVersions: string[];
  isActive: boolean;
  tags: string[];
  categories: string[];
}

// API Response types
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: unknown;
  };
}

export interface PaginatedResponse<T> extends ApiResponse<T[]> {
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

// Sort types
export type ModSortField = "downloads" | "rating" | "date";
export type GameSortField = "modCount" | "name" | "date";

// Search types
export interface SearchQuery<T = ModSortField> {
  q?: string;
  tags?: string[];
  gameId?: string;
  creatorId?: string;
  sort?: T;
  order?: "asc" | "desc";
  page?: number;
  limit?: number;
}

// File upload types
export interface UploadedFile {
  originalname: string;
  filename: string;
  path: string;
  size: number;
  mimetype: string;
}
