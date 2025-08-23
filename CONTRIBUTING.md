## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Run the test suite (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Guidelines

- Follow Go best practices and idioms
- Write tests for new functionality
- Update documentation as needed
- Use the provided Makefile commands
- Follow the existing code structure and patterns

## Development

### Database

The application uses the SQLite database. More info here needed.

### Configuration

Configuration is handled through environment variables, please refer to the [Environment Variables](backend/.env) for more information.

## Deployment

### Production Deployment

1. **Build the application**
```bash
cd backend/
make build-prod
```

2. **Set production environment variables**
```bash
export ENV=production
```

3. **Run the server**
```bash
./bin/azurite-server
```

### Docker Deployment

```bash
make docker-build
make docker-run
```

### Reverse Proxy Setup (Nginx)

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /files/ {
        proxy_pass http://localhost:8080;
    }

    location / {
        # Frontend static files
        root /path/to/frontend/build;
        try_files $uri $uri/ /index.html;
    }
}
```
