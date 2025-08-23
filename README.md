# Azurite
[![CI](https://github.com/AbstractMelon/Azurite/actions/workflows/ci.yml/badge.svg)](https://github.com/AbstractMelon/Azurite/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/AbstractMelon/Azurite)](https://goreportcard.com/report/github.com/AbstractMelon/Azurite)
[![codecov](https://codecov.io/gh/AbstractMelon/Azurite/branch/main/graph/badge.svg)](https://codecov.io/gh/AbstractMelon/Azurite)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Frontend: SvelteKit](https://img.shields.io/badge/Frontend-SvelteKit-FF3E00?logo=svelte&logoColor=white)](https://kit.svelte.dev/)
[![pnpm](https://img.shields.io/badge/Package%20Manager-pnpm-%23F69220?logo=pnpm&logoColor=white)](https://pnpm.io/)

A complete modding website platform built for community-driven mod management.

## Overview

Azurite is a full-stack modding platform that allows users to:
- **Discover** mods for their favorite games
- **Share** their own mods with the community
- **Manage** mod versions, dependencies, and files
- **Interact** with other modders through comments and ratings
- **Document** games with community-maintained wikis
- **Moderate** content with role-based permissions

## Features

- **User Authentication**: Email/password and OAuth (GitHub, Discord, Google)
- **Game Management**: Complete game catalog with request system
- **Mod System**: Upload, version, dependency management with malware scanning
- **Community Features**: Comments, likes, ratings, and user profiles
- **Documentation**: Markdown-based game wikis with collaborative editing
- **Search & Discovery**: Advanced filtering, sorting, and search capabilities
- **Role-Based Access Control**: Users, Moderators, Wiki Maintainers, Admins
- **Notification System**: In-app and email notifications
- **File Management**: Custom CDN with WebP image conversion
- **Moderation Tools**: Content approval, user/IP banning, soft deletion
- **Analytics**: Download tracking, user statistics, activity monitoring
- **API-First**: Complete REST API for third-party integrations

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP router)
- **Database**: SQLite3
- **Authentication**: JWT with OAuth integration
- **File Storage**: Local filesystem with custom CDN
- **Email**: SMTP integration for notifications

### Frontend
- **Framework**: SvelteKit
- **Styling**: Tailwind CSS
- **Package Manager**: pnpm

## Backend Setup

### Prerequisites

- Go 1.21 or higher
- Git

1. **Enter the backend directory**
```bash
cd backend
```

2. **Setup environment**
```bash
make setup-dev
```

3. **Configure environment variables**
Edit the `.env` file with your settings.

4. **Start the development server**
```bash
make dev
```

The API will be available at `http://localhost:8080`

## Frontend Setup

### Prerequisites

- Node.js 18 or higher
- pnpm

1. **Enter the frontend directory**
```bash
cd frontend
```

2. **Setup environment**
```bash
pnpm install
```

3. **Start the development server**
```bash
pnpm dev
```

The frontend will be available at `http://localhost:5173`

## Development

Please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file for details on how to contribute to this project.

## License

This project is licensed under the GNU General Public License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: Check the README files in each directory
- **Issues**: Report bugs and request features via GitHub Issues
- **Discussions**: Join our community discussions
- **Email**: Contact the maintainers directly
