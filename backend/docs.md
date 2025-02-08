# API Documentation

## Authentication

### Register

- **POST /api/auth/register**
  - Body: `username`, `email`, `password`, `role` (optional)
  - Response: `token`, `user`

### Login

- **POST /api/auth/login**
  - Body: `email`, `password`
  - Response: `token`, `user`

### Logout

- **POST /api/auth/logout**
  - Response: `message`

## Games

### Create Game

- **POST /api/games**
  - Body: `name`, `description`, `shortDescription`, `tags`, `coverImageUrl` (optional), `websiteUrl` (optional), `supportedModTypes`, `latestVersion`, `supportedVersions`, `isActive`
  - Response: `game`

### Get Games

- **GET /api/games**
  - Query: `page`, `limit`, `search`, `tags`, `sort`, `order`
  - Response: `games`, `pagination`

### Get Game

- **GET /api/games/:gameId**
  - Path: `gameId`
  - Response: `game`

### Update Game

- **PATCH /api/games/:gameId**
  - Path: `gameId`
  - Body: `name`, `description`, `shortDescription`, `tags`, `coverImageUrl`, `websiteUrl`, `supportedModTypes`, `latestVersion`, `supportedVersions`, `isActive`
  - Response: `game`

### Delete Game

- **DELETE /api/games/:gameId**
  - Path: `gameId`
  - Response: `message`

### Toggle Game Status

- **POST /api/games/:gameId/toggle-status**
  - Path: `gameId`
  - Response: `game`

### Add Category

- **POST /api/games/:gameId/categories**
  - Path: `gameId`
  - Body: `category`
  - Response: `game`

### Remove Category

- **DELETE /api/games/:gameId/categories**
  - Path: `gameId`
  - Body: `category`
  - Response: `message`

### Add Tag

- **POST /api/games/:gameId/tags**
  - Path: `gameId`
  - Body: `tag`
  - Response: `game`

### Remove Tag

- **DELETE /api/games/:gameId/tags**
  - Path: `gameId`
  - Body: `tag`
  - Response: `message`

## Mods

### Create Mod

- **POST /api/mods**
  - Body: `name`, `description`, `shortDescription`, `tags`, `gameId`, `file`
  - Response: `mod`

### Get Mods

- **GET /api/mods**
  - Query: `page`, `limit`, `search`, `tags`, `sort`, `order`, `gameId`
  - Response: `mods`, `pagination`

### Get Mod

- **GET /api/mods/:modId**
  - Path: `modId`
  - Response: `mod`

### Update Mod

- **PATCH /api/mods/:modId**
  - Path: `modId`
  - Body: `name`, `description`, `shortDescription`, `tags`, `file`
  - Response: `mod`

### Upload Mod File

- **POST /api/mods/:modId/file**
  - Path: `modId`
  - Body: `file`
  - Response: `mod`

### Upload Screenshots

- **POST /api/mods/:modId/screenshots**
  - Path: `modId`
  - Body: `screenshots`
  - Response: `mod`

### Publish Mod

- **POST /api/mods/:modId/publish**
  - Path: `modId`
  - Response: `mod`

### Add Rating

- **POST /api/mods/:modId/ratings**
  - Path: `modId`
  - Body: `rating`
  - Response: `mod`

### Add Comment

- **POST /api/mods/:modId/comments**
  - Path: `modId`
  - Body: `comment`
  - Response: `mod`

### Increment Downloads

- **POST /api/mods/:modId/downloads**
  - Path: `modId`
  - Response: `mod`

## Types

### User

- `id`, `username`, `email`, `password`, `role`, `displayName` (optional), `avatarUrl` (optional), `bio` (optional), `mods` (optional), `favorites` (optional)

### Game

- `id`, `name`, `description`, `shortDescription`, `tags`, `coverImageUrl` (optional), `websiteUrl` (optional), `supportedModTypes`, `latestVersion`, `supportedVersions`, `isActive`, `categories`

### Mod

- `id`, `name`, `description`, `shortDescription`, `version`, `gameId`, `creatorId`, `filePath`, `fileSize`, `tags`, `downloads`, `ratings`, `averageRating`, `comments`, `changelog`, `screenshots`, `requirements` (optional), `isPublished`

### Pagination

- `page`, `limit`, `total`, `totalPages`
