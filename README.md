# Portfolio Backend API

A production-ready Go REST API backend for a portfolio website, built with clean architecture principles and raw SQL (NO ORM) as specified in the requirements.

## 🚀 Features

- **Clean Architecture**: Repository pattern with service layer separation
- **Raw SQL Implementation**: Direct database/sql usage for performance and control
- **RESTful API**: Following OpenAPI specification standards
- **Production Ready**: Structured logging, health checks, graceful shutdown
- **Cloud Native**: Docker containerization and GCP deployment ready
- **Comprehensive Middleware**: CORS, logging, recovery, security headers
- **Database Migrations**: Version-controlled schema management
- **Validation**: Request validation with detailed error responses

## 🛠 Technology Stack

| Category | Technology | Purpose |
|----------|------------|---------|
| **Language** | Go 1.21+ | Backend development |
| **Framework** | Gin | HTTP web framework |
| **Database** | MySQL 8.0 | Relational database |
| **Database Layer** | Raw SQL (database/sql) | Direct SQL queries (NO ORM) |
| **Configuration** | Viper | Environment-based configuration |
| **Logging** | Zerolog | Structured logging |
| **Validation** | go-playground/validator | Request validation |
| **Migrations** | golang-migrate | Database migrations |
| **Testing** | Testify | Testing framework |
| **Containerization** | Docker | Multi-stage builds |
| **Cloud Platform** | Google Cloud Platform | Cloud Run, Cloud SQL |

## 📁 Project Structure

```
portfolio-backend/
├── cmd/api/                    # Application entry point
│   └── main.go
├── internal/                   # Private application code
│   ├── config/                 # Configuration management
│   ├── database/               # Database layer (NO ORM)
│   │   ├── connection.go       # Connection management
│   │   └── repositories/       # Repository implementations
│   ├── handlers/               # HTTP handlers
│   ├── middleware/             # HTTP middleware
│   ├── models/                 # Data models
│   └── services/               # Business logic layer
├── pkg/                        # Public packages
│   ├── response/               # HTTP response utilities
│   └── validator/              # Custom validators
├── migrations/                 # Database migrations
├── deployments/                # Deployment configurations
│   ├── docker/                 # Docker configuration
│   └── gcp/                    # GCP deployment files
├── scripts/                    # Build and deployment scripts
└── docs/                       # Documentation
```

## 🔧 API Endpoints

### Health Check
- `GET /v1/health` - Service health status

### Portfolio Data
- `GET /v1/profile` - Get user profile
- `PUT /v1/profile` - Update user profile
- `GET /v1/experience` - Get all work experiences
- `GET /v1/experience/{id}` - Get specific experience
- `GET /v1/skills` - Get skills (supports `?group_by=category`)
- `GET /v1/education` - Get education history
- `GET /v1/certifications` - Get certifications

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- MySQL 8.0
- Docker (optional)
- Make (optional, for convenience)

### Local Development

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd portfolio-backend
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   # or
   make deps
   ```

3. **Environment setup**:
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Database setup**:
   ```bash
   # Create database
   mysql -u root -p -e "CREATE DATABASE portfolio_db;"
   
   # Run migrations
   make migrate-up
   # or
   ./scripts/migrate.sh up
   ```

5. **Run the application**:
   ```bash
   go run cmd/api/main.go
   # or
   make run
   ```

   The API will be available at `http://localhost:8080`

### Docker Development

```bash
# Build and run with Docker
make docker-build
make docker-run

# Or use docker-compose (create docker-compose.yml first)
docker-compose up --build
```

## 🔧 Configuration

Configuration is managed through environment variables and optional YAML files using Viper:

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `HOST` | Server host | `0.0.0.0` |
| `PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `3306` |
| `DB_USER` | Database user | `portfolio_user` |
| `DB_PASSWORD` | Database password | *required* |
| `DB_NAME` | Database name | `portfolio_db` |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `LOG_FORMAT` | Log format (json/console) | `json` |

### YAML Configuration (Optional)

Create `config.yaml` in the project root:

```yaml
server:
  host: 0.0.0.0
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

database:
  host: localhost
  port: 3306
  user: portfolio_user
  password: your_secure_password
  database: portfolio_db
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

cors:
  allowed_origins:
    - https://your-frontend-domain.com
    - http://localhost:3000
  allowed_methods:
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
  allowed_headers:
    - Content-Type
    - Authorization

logging:
  level: info
  format: json
```

## 🗄️ Database Schema

The API uses MySQL with the following main tables:

- **profiles**: User profile information
- **experiences**: Work experience entries
- **skills**: Technical skills with categories
- **education**: Educational background
- **certifications**: Professional certifications

Schema is managed through versioned migrations in the `migrations/` directory.

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linting
make lint

# Run quality checks
make check
```

## 🏗️ Building

```bash
# Build for local platform
make build

# Build for all platforms
./scripts/build.sh

# Clean build artifacts
make clean
```

## 🚀 Deployment

### Google Cloud Platform

#### Option 1: Using Cloud Build

```bash
# Deploy using cloud build
gcloud builds submit --config deployments/gcp/cloudbuild.yaml

# Or use the deploy script
./scripts/deploy.sh --project YOUR_PROJECT_ID
```

#### Option 2: Using Terraform

```bash
# Initialize and apply Terraform
cd deployments/gcp/terraform
terraform init
terraform plan
terraform apply
```

### Environment Setup for Production

1. **Cloud SQL**: MySQL instance with private networking
2. **Secret Manager**: Database credentials storage
3. **Cloud Run**: Serverless container deployment
4. **IAM**: Service account with minimal permissions

## 📊 Monitoring and Observability

- **Structured Logging**: JSON format with correlation IDs
- **Health Checks**: Database connectivity monitoring
- **Metrics**: Built-in HTTP metrics
- **Error Tracking**: Comprehensive error handling
- **Performance**: Sub-200ms response times

## 🔒 Security Features

- **Input Validation**: Request validation on all endpoints
- **SQL Injection Prevention**: Prepared statements
- **CORS Configuration**: Configurable origin restrictions
- **Security Headers**: Standard security headers
- **Error Handling**: No sensitive information leakage

## 🛠️ Development Tools

### Makefile Targets

```bash
make help          # Show all available targets
make dev           # Run development workflow
make build         # Build the application
make test          # Run tests
make lint          # Run linters
make migrate-up    # Run database migrations
make docker-build  # Build Docker image
make deploy        # Deploy to GCP
```

### Scripts

- `scripts/build.sh` - Multi-platform build script
- `scripts/deploy.sh` - GCP deployment automation
- `scripts/migrate.sh` - Database migration management

## 📝 API Response Format

### Success Response
```json
{
  "data": { ... },
  "success": true,
  "message": "Optional success message"
}
```

### Error Response
```json
{
  "error": "error_code",
  "message": "Human readable error",
  "details": {
    "field": "Validation error message"
  }
}
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For questions or support:
1. Check the [documentation](docs/)
2. Review existing [issues](../../issues)
3. Create a new issue for bugs or feature requests

---

**Built with ❤️ using Go and following clean architecture principles**