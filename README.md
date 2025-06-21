# Moon API

A RESTful API built with Go, Gin, and GORM following Clean Architecture principles.

## Features

- 🔐 JWT Authentication
- 🗄️ MySQL Database with GORM
- 📝 Structured Logging with Zap
- 🐳 Docker & Docker Compose support
- 🏗️ Clean Architecture
- ⚙️ Configuration Management
- 🔒 Password Hashing with bcrypt
- 🧪 Testing Support

## Project Structure

```
Moon/
├── cmd/
│   └── main.go              # Application entry point
├── configs/
│   └── config.yaml          # Configuration file
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database connection
│   ├── domain/              # Domain models
│   │   ├── user/
│   │   └── product/
│   ├── handler/             # HTTP handlers
│   ├── middleware/          # Middleware functions
│   ├── repository/          # Data access layer
│   └── usecase/             # Business logic
├── migrations/              # Database migrations
├── pkg/                     # Shared packages
│   ├── hash/                # Password hashing
│   ├── jwt/                 # JWT utilities
│   ├── logger/              # Logging utilities
│   ├── utils/               # Utility functions
│   └── validator/           # Validation utilities
├── static/                  # Static files
├── Dockerfile               # Docker configuration
├── docker-compose.yml       # Docker Compose configuration
├── Makefile                 # Build and development commands
└── README.md               # This file
```

## Prerequisites

- Go 1.24 or higher
- MySQL 8.0 or higher
- Redis (optional)
- Docker & Docker Compose (optional)

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd Moon
```

2. Start the application with Docker Compose:
```bash
docker-compose up -d
```

The application will be available at `http://localhost:8080`

### Manual Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd Moon
```

2. Install dependencies:
```bash
make deps
```

3. Set up your environment variables (create a `.env` file):
```env
APP_PORT=8080
APP_MODE=debug
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_NAME=moon_db
JWT_SECRET=your-super-secret-jwt-key
```

4. Start MySQL and Redis (if using)

5. Run the application:
```bash
make run
```

## API Endpoints

### Health Check
- `GET /ping` - Basic health check
- `GET /api/v1/health` - Detailed health status

### Authentication (TODO)
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login

### Users (TODO)
- `GET /api/v1/users/profile` - Get user profile (protected)
- `PUT /api/v1/users/profile` - Update user profile (protected)

### Products (TODO)
- `GET /api/v1/products` - List products
- `POST /api/v1/products` - Create product (admin only)
- `GET /api/v1/products/:id` - Get product by ID
- `PUT /api/v1/products/:id` - Update product (admin only)
- `DELETE /api/v1/products/:id` - Delete product (admin only)

## Development

### Available Make Commands

```bash
make build          # Build the application
make run            # Run the application
make test           # Run tests
make test-coverage  # Run tests with coverage
make clean          # Clean build artifacts
make deps           # Install dependencies
make lint           # Run linter
make fmt            # Format code
make vet            # Vet code
make docker-build   # Build Docker image
make docker-run     # Run Docker container
```

### Database Migrations

Database migrations are stored in the `migrations/` directory. To run migrations:

```bash
make migrate
```

### Testing

Run tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

## Configuration

The application uses YAML configuration files located in `configs/`. You can override configuration values using environment variables.

### Configuration Structure

```yaml
app:
  name: "Moon API"
  version: "1.0.0"
  port: 8080
  mode: "debug"

database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  name: "moon_db"

jwt:
  secret: "your-secret-key"
  expires_in: 24

redis:
  host: "localhost"
  port: 6379

logger:
  level: "info"
  format: "json"
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | Application port | 8080 |
| `APP_MODE` | Application mode (debug/release) | debug |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 3306 |
| `DB_USERNAME` | Database username | root |
| `DB_PASSWORD` | Database password | password |
| `DB_NAME` | Database name | moon_db |
| `JWT_SECRET` | JWT secret key | - |
| `JWT_EXPIRES_IN` | JWT expiration hours | 24 |
| `REDIS_HOST` | Redis host | localhost |
| `REDIS_PORT` | Redis port | 6379 |
| `LOG_LEVEL` | Log level | info |

## Docker

### Build Image
```bash
make docker-build
```

### Run Container
```bash
make docker-run
```

### Using Docker Compose
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop all services
docker-compose down
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 