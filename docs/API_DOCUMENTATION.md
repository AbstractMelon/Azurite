# Azurite Backend API Documentation

This documentation provides a guide to the Azurite backend API for frontend developers.

## Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Response Format](#response-format)
- [Error Handling](#error-handling)
- [Endpoints](#endpoints)
  - [Authentication](#authentication-endpoints)
  - [Users](#user-endpoints)
  - [Games](#game-endpoints)
  - [Mods](#mod-endpoints)
  - [Comments](#comment-endpoints)
  - [Notifications](#notification-endpoints)
  - [Documentation](#documentation-endpoints)
  - [Admin](#admin-endpoints)
  - [File Serving](#file-serving)
- [Data Models](#data-models)
- [Constants](#constants)

## Overview

The Azurite API is a RESTful API built with Go and Gin framework. All API endpoints are prefixed with `/api` and return JSON responses.

**Base URL**: `http://localhost:8080/api` (development)

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. After login, include the token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

### OAuth Providers
- GitHub
- Google
- Discord

## Response Format

All API responses follow a consistent format:

```json
{
  "success": true,
  "data": {},
  "message": "Optional message"
}
```

For paginated responses:

```json
{
  "success": true,
  "data": {
    "data": [],
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

## Error Handling

Error responses include an error message:

```json
{
  "success": false,
  "error": "Error description"
}
```

Common HTTP status codes:
- `200` - OK
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

## Endpoints

### Authentication Endpoints

#### Register User
```
POST /api/auth/register
```

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "display_name": "John Doe"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "jwt_token_here",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "display_name": "John Doe",
      "role": "user",
      "created_at": "2023-12-01T10:00:00Z"
    }
  }
}
```

#### Login User
```
POST /api/auth/login
```

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "jwt_token_here",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "display_name": "John Doe",
      "role": "user"
    }
  }
}
```

#### Request Password Reset
```
POST /api/auth/forgot-password
```

**Request Body:**
```json
{
  "email": "john@example.com"
}
```

#### Reset Password
```
POST /api/auth/reset-password
```

**Request Body:**
```json
{
  "token": "reset_token",
  "password": "new_password123"
}
```

#### OAuth Endpoints
```
GET /api/auth/github
GET /api/auth/google
GET /api/auth/discord
GET /api/auth/callback/github
GET /api/auth/callback/google
GET /api/auth/callback/discord
```

#### Get User Profile
```
GET /api/auth/profile
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "display_name": "John Doe",
    "avatar": "avatar_url",
    "bio": "User bio",
    "role": "user",
    "notify_email": true,
    "notify_in_site": true
  }
}
```

#### Update Profile
```
PUT /api/auth/profile
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "display_name": "John Smith",
  "bio": "Updated bio",
  "notify_email": false,
  "notify_in_site": true
}
```

#### Update Password
```
PUT /api/auth/password
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "current_password": "old_password",
  "new_password": "new_password123"
}
```

### User Endpoints

#### Get User by ID
```
GET /api/users/:id
```

#### Get User by Username
```
GET /api/users/username/:username
```

### Game Endpoints

#### List Games
```
GET /api/games
```

**Query Parameters:**
- `page` (int, default: 1)
- `per_page` (int, default: 20)
- `search` (string)

**Response:**
```json
{
  "success": true,
  "data": {
    "data": [
      {
        "id": 1,
        "name": "Minecraft",
        "slug": "minecraft",
        "description": "Block building game",
        "icon": "minecraft_icon.png",
        "is_active": true,
        "mod_count": 150,
        "created_at": "2023-12-01T10:00:00Z"
      }
    ],
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

#### Get Game by Slug
```
GET /api/games/:slug
```

#### Get Game Mods
```
GET /api/games/:slug/mods
```

**Query Parameters:**
- `page` (int)
- `per_page` (int)
- `sort` (string: "newest", "oldest", "downloads", "likes")
- `tags` (string: comma-separated tag names)

#### Get Game Tags
```
GET /api/games/:slug/tags
```

#### Get Game Moderators
```
GET /api/games/:slug/moderators
```

#### Create Game Request
```
POST /api/games/requests
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "New Game",
  "description": "Description of the game"
}
```

#### Admin: List Game Requests
```
GET /api/games/requests
Authorization: Bearer <token>
Role: admin
```

#### Admin: Create Game
```
POST /api/games
Authorization: Bearer <token>
Role: admin
```

#### Admin: Update Game
```
PUT /api/games/manage/:id
Authorization: Bearer <token>
Role: admin
```

#### Admin: Delete Game
```
DELETE /api/games/manage/:id
Authorization: Bearer <token>
Role: admin
```

### Mod Endpoints

#### Get Mod
```
GET /api/mods/:gameSlug/:modSlug
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Amazing Mod",
    "slug": "amazing-mod",
    "description": "This is an amazing mod",
    "short_description": "Amazing mod for the game",
    "version": "1.0.0",
    "game_version": "1.19.2",
    "downloads": 500,
    "likes": 25,
    "is_liked": false,
    "owner": {
      "id": 1,
      "username": "modder1",
      "display_name": "Modder One"
    },
    "game": {
      "id": 1,
      "name": "Minecraft",
      "slug": "minecraft"
    },
    "tags": [
      {
        "id": 1,
        "name": "Technology",
        "slug": "technology"
      }
    ],
    "files": [
      {
        "id": 1,
        "filename": "mod.jar",
        "file_size": 1024000,
        "is_main": true,
        "created_at": "2023-12-01T10:00:00Z"
      }
    ],
    "created_at": "2023-12-01T10:00:00Z"
  }
}
```

#### Create Mod
```
POST /api/mods
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "My New Mod",
  "description": "Detailed description",
  "short_description": "Brief description",
  "version": "1.0.0",
  "game_version": "1.19.2",
  "game_id": 1,
  "source_website": "https://github.com/user/mod",
  "contact_info": "user@example.com",
  "tags": ["technology", "utility"],
  "dependencies": [2, 3]
}
```

#### Update Mod
```
PUT /api/mods/:id
Authorization: Bearer <token>
```

#### Delete Mod
```
DELETE /api/mods/:id
Authorization: Bearer <token>
```

#### Like Mod
```
POST /api/mods/:id/like
Authorization: Bearer <token>
```

#### Unlike Mod
```
DELETE /api/mods/:id/like
Authorization: Bearer <token>
```

#### Upload Mod File
```
POST /api/mods/:id/files
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Form Data:**
- `file` (file)
- `is_main` (boolean, "true" or "false")

#### Moderator: Approve Mod
```
POST /api/mods/:id/approve
Authorization: Bearer <token>
Role: admin, community_moderator
```

#### Moderator: Reject Mod
```
POST /api/mods/:id/reject
Authorization: Bearer <token>
Role: admin, community_moderator
```

**Request Body:**
```json
{
  "reason": "Rejection reason"
}
```

### Comment Endpoints

#### Get Mod Comments
```
GET /api/comments/mod/:modId
```

**Query Parameters:**
- `page` (int, default: 1)
- `per_page` (int, default: 20)

#### Create Comment
```
POST /api/comments?mod_id=:modId
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "content": "This is a great mod!",
  "parent_id": null
}
```

#### Update Comment
```
PUT /api/comments/:id
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "content": "Updated comment content"
}
```

#### Delete Comment
```
DELETE /api/comments/:id
Authorization: Bearer <token>
```

### Notification Endpoints

#### Get User Notifications
```
GET /api/notifications
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (int, default: 1)
- `per_page` (int, default: 20)

#### Get Unread Count
```
GET /api/notifications/unread-count
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "count": 5
  }
}
```

#### Mark Notification as Read
```
PUT /api/notifications/:id/read
Authorization: Bearer <token>
```

#### Mark All Notifications as Read
```
PUT /api/notifications/read-all
Authorization: Bearer <token>
```

#### Delete Notification
```
DELETE /api/notifications/:id
Authorization: Bearer <token>
```

### Documentation Endpoints

#### Get Game Documentation
```
GET /api/docs/:gameSlug
```

**Query Parameters:**
- `page` (int, default: 1)
- `per_page` (int, default: 20)

#### Get Documentation Article
```
GET /api/docs/:gameSlug/:docSlug
```

#### Create Documentation
```
POST /api/docs/:gameSlug
Authorization: Bearer <token>
Role: admin, community_moderator, wiki_maintainer
```

**Request Body:**
```json
{
  "title": "Getting Started",
  "content": "# Getting Started\n\nThis is the content..."
}
```

#### Update Documentation
```
PUT /api/docs/:gameSlug/:docSlug
Authorization: Bearer <token>
Role: admin, community_moderator, wiki_maintainer
```

#### Delete Documentation
```
DELETE /api/docs/:gameSlug/:docSlug
Authorization: Bearer <token>
Role: admin, community_moderator, wiki_maintainer
```

### Admin Endpoints

#### Get Statistics
```
GET /api/admin/stats
Authorization: Bearer <token>
Role: admin
```

**Response:**
```json
{
  "success": true,
  "data": {
    "users": {
      "total": 1000,
      "active": 800,
      "new_this_month": 50
    },
    "mods": {
      "total": 500,
      "approved": 450,
      "pending": 50
    },
    "system": {
      "uptime": "72h30m",
      "version": "1.0.0"
    }
  }
}
```

#### Get Recent Activity
```
GET /api/admin/activity
Authorization: Bearer <token>
Role: admin
```

**Query Parameters:**
- `limit` (int, default: 50)

#### Create Ban
```
POST /api/admin/bans
Authorization: Bearer <token>
Role: admin
```

**Request Body:**
```json
{
  "user_id": 123,
  "ip_address": "192.168.1.1",
  "game_id": 1,
  "reason": "Violation of terms",
  "duration": 7
}
```

#### List Bans
```
GET /api/admin/bans
Authorization: Bearer <token>
Role: admin
```

**Query Parameters:**
- `page` (int, default: 1)
- `per_page` (int, default: 20)
- `active` (boolean, default: true)

#### Unban User
```
POST /api/admin/bans/:id/unban
Authorization: Bearer <token>
Role: admin
```

### File Serving

#### Static Files
```
GET /files/mods/:filename
GET /files/images/:filename
```

#### Download Mod
```
GET /download/:gameSlug/:modSlug
```

This endpoint serves the main mod file and increments the download counter.

#### Health Check
```
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "service": "azurite-api"
}
```

## Data Models

### User
```typescript
interface User {
  id: number;
  username: string;
  email: string;
  display_name: string;
  avatar: string;
  bio: string;
  role: string;
  is_active: boolean;
  email_verified: boolean;
  notify_email: boolean;
  notify_in_site: boolean;
  created_at: string;
  updated_at: string;
  last_login_at: string;
}
```

### Game
```typescript
interface Game {
  id: number;
  name: string;
  slug: string;
  description: string;
  icon: string;
  is_active: boolean;
  mod_count: number;
  created_at: string;
  updated_at: string;
}
```

### Mod
```typescript
interface Mod {
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
  files?: ModFile[];
  is_liked?: boolean;
}
```

### ModFile
```typescript
interface ModFile {
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
```

### Comment
```typescript
interface Comment {
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
```

### Notification
```typescript
interface Notification {
  id: number;
  user_id: number;
  type: string;
  title: string;
  message: string;
  data: string;
  is_read: boolean;
  created_at: string;
}
```

### Tag
```typescript
interface Tag {
  id: number;
  name: string;
  slug: string;
  game_id: number;
  created_at: string;
}
```

### GameRequest
```typescript
interface GameRequest {
  id: number;
  name: string;
  description: string;
  icon: string;
  requested_by: number;
  status: string;
  admin_notes: string;
  created_at: string;
  updated_at: string;
  user?: User;
}
```

### Documentation
```typescript
interface Documentation {
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
```

### Ban
```typescript
interface Ban {
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
```

### PaginatedResponse
```typescript
interface PaginatedResponse<T> {
  data: T[];
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}
```

## Constants

### User Roles
- `user` - Regular user
- `admin` - Administrator
- `community_moderator` - Community moderator
- `wiki_maintainer` - Wiki maintainer

### Game Request Status
- `pending` - Awaiting review
- `approved` - Approved by admin
- `denied` - Denied by admin

### Notification Types
- `mod_rejected` - Mod was rejected
- `new_comment` - New comment on mod
- `mod_milestone` - Mod reached milestone
- `game_request` - Game request status update
- `mod_approved` - Mod was approved

### Scan Results
- `pending` - Scan in progress
- `clean` - No threats detected
- `threat` - Threat detected

## Rate Limits

The API implements rate limiting on certain endpoints:
- Authentication endpoints: 5 requests per minute per IP
- File upload endpoints: 10 requests per hour per user
- Comment creation: 60 requests per hour per user

## File Upload Limits

- Maximum file size: 100MB (configurable)
- Supported file types: `.jar`, `.zip`, `.rar`, `.7z`, `.tar.gz`
- Image files: `.png`, `.jpg`, `.jpeg`, `.webp` (auto-converted to WebP)

## Security Headers

All responses include security headers:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security: max-age=31536000`

## CORS Configuration

CORS is enabled for frontend development. In production, configure allowed origins appropriately.

## WebSocket Support

The API includes WebSocket support for real-time notifications (implementation details coming soon).

---

*This documentation is for Azurite API v1.0. For the latest updates, check the repository.*
