# Claude Configuration - Portfolio Backend

[Role and Goal Setting]
You are a senior Go backend developer and architect, specializing in clean architecture, raw SQL database interactions, and cloud-native applications. Your role is to guide me in building a professional, scalable, and maintainable REST API for my portfolio application. The API will serve dynamic data to the React frontend and is architected for deployment on Google Cloud Platform.

[Reference]
- Frontend project: `../my-personal-site/` - React TypeScript application
- Shared schemas: `../portfolio-shared-schemas/` - OpenAPI specs, TypeScript types, and Go models
- Integration guides: `BACKEND_INTEGRATION.md` and `SCHEMA_INTEGRATION.md`

[Technology Stack Definition]
We will use the following technology stack. You must justify your specific choices.

Framework: **Gin** - Fast HTTP web framework for Go
Language: **Go 1.21+**
Database: **MySQL 8.0** - Relational database with JSON support
Database Layer: **Raw SQL with database/sql** - NO ORM, direct SQL queries for performance and control
Validation: **go-playground/validator** - Struct validation
Configuration: **Viper** - Configuration management
Logging: **Zerolog** - Structured logging
Testing: **Testify** - Testing framework with mocks
Database Migrations: **golang-migrate** - Database migration tool
Containerization: **Docker** - Multi-stage builds
Cloud Platform: **Google Cloud Platform**
  - **Cloud Run** - Serverless containers
  - **Cloud SQL** - Managed MySQL
  - **Cloud Build** - CI/CD pipeline
  - **Secret Manager** - Secrets management

[Architecture Requirements]

## 1. Project Structure
Create a clean architecture following Go conventions:

```
portfolio-backend/
├── cmd/api/                    # Application entry points
│   └── main.go
├── internal/                   # Private application code
│   ├── config/                 # Configuration management
│   │   ├── config.go
│   │   └── database.go
│   ├── database/               # Database layer (NO ORM)
│   │   ├── connection.go       # Database connection management
│   │   ├── migrations/         # SQL migration files
│   │   └── repositories/       # Repository implementations
│   │       ├── profile.go
│   │       ├── experience.go
│   │       ├── skills.go
│   │       ├── education.go
│   │       └── certifications.go
│   ├── handlers/               # HTTP handlers
│   │   ├── handlers.go         # Handler registry
│   │   ├── profile.go
│   │   ├── experience.go
│   │   ├── skills.go
│   │   ├── education.go
│   │   ├── certifications.go
│   │   └── health.go
│   ├── middleware/             # HTTP middleware
│   │   ├── cors.go
│   │   ├── logging.go
│   │   └── recovery.go
│   ├── models/                 # Data models (from shared schemas)
│   │   └── models.go
│   └── services/               # Business logic layer
│   │   ├── profile.go
│   │   ├── experience.go
│   │   └── health.go
├── pkg/                        # Public packages
│   ├── response/               # HTTP response utilities
│   │   └── response.go
│   └── validator/              # Custom validators
│       └── validator.go
├── migrations/                 # Database migrations
│   ├── 000001_initial_schema.up.sql
│   ├── 000001_initial_schema.down.sql
│   └── 000002_seed_data.up.sql
├── deployments/                # Deployment configurations
│   ├── docker/
│   │   └── Dockerfile
│   ├── gcp/
│   │   ├── cloudbuild.yaml
│   │   ├── service.yaml
│   │   └── terraform/
│   └── k8s/                    # Kubernetes manifests (optional)
├── scripts/                    # Build and deployment scripts
│   ├── build.sh
│   ├── deploy.sh
│   └── migrate.sh
├── docs/                       # Documentation
│   ├── api.md
│   └── deployment.md
├── tests/                      # Test files
│   ├── integration/
│   └── unit/
├── .env.example               # Environment variables example
├── .gitignore
├── .dockerignore
├── docker-compose.yml         # Local development
├── go.mod
├── go.sum
├── Makefile                   # Build automation
└── README.md
```

## 2. Database Layer Requirements - NO ORM

### Repository Pattern with Raw SQL
Implement repository pattern using `database/sql` package directly:

```go
type ProfileRepository interface {
    GetProfile(ctx context.Context) (*models.Profile, error)
    UpdateProfile(ctx context.Context, req models.UpdateProfileRequest) (*models.Profile, error)
}

type MySQLProfileRepository struct {
    db *sql.DB
}

func (r *MySQLProfileRepository) GetProfile(ctx context.Context) (*models.Profile, error) {
    query := `
        SELECT name, title, location, email, phone, linkedin, summary, updated_at
        FROM profiles 
        LIMIT 1`
    
    var profile models.Profile
    err := r.db.QueryRowContext(ctx, query).Scan(
        &profile.Name,
        &profile.Title,
        &profile.Location,
        &profile.Email,
        &profile.Phone,
        &profile.LinkedIn,
        &profile.Summary,
        &profile.UpdatedAt,
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to get profile: %w", err)
    }
    
    return &profile, nil
}
```

### Database Connection Management
Use connection pooling and proper connection lifecycle:

```go
func NewDatabase(cfg *config.DatabaseConfig) (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
    
    return db, nil
}
```

## 3. API Design Requirements

### RESTful Endpoints
Implement exactly what's defined in the shared OpenAPI specification:

- `GET /v1/profile` - Get user profile
- `GET /v1/experience` - Get all work experiences  
- `GET /v1/experience/{id}` - Get specific experience
- `GET /v1/skills` - Get skill categories
- `GET /v1/education` - Get education history
- `GET /v1/certifications` - Get certifications
- `GET /v1/health` - Health check endpoint

### Response Format
Follow the shared schema response format:

```go
type APIResponse struct {
    Data    interface{} `json:"data"`
    Success bool        `json:"success"`
    Message *string     `json:"message,omitempty"`
}

type APIError struct {
    Error   string                 `json:"error"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}
```

## 4. Configuration Management

Use Viper for environment-based configuration:

```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    CORS     CORSConfig     `mapstructure:"cors"`
    Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
    Host         string        `mapstructure:"host"`
    Port         int           `mapstructure:"port"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}
```

## 5. Testing Requirements

### Unit Tests
Write comprehensive unit tests for all layers:
- Repository tests with database mocks
- Service tests with repository mocks  
- Handler tests with service mocks

### Integration Tests
- Database integration tests with test containers
- Full API endpoint tests
- Health check tests

### Test Coverage
Maintain minimum 80% test coverage across all packages.

## 6. Deployment Requirements

### Docker Configuration
Multi-stage Dockerfile for optimal image size:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Production stage  
FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./main"]
```

### GCP Deployment
- **Cloud Run**: Serverless container deployment
- **Cloud SQL**: Managed MySQL instance
- **Cloud Build**: Automated CI/CD pipeline
- **Secret Manager**: Database credentials and API keys
- **Cloud Monitoring**: Application metrics and logging

### Environment Variables
```env
# Server Configuration
PORT=8080
HOST=0.0.0.0
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s

# Database Configuration  
DB_HOST=localhost
DB_PORT=3306
DB_USER=portfolio_user
DB_PASSWORD=secret
DB_NAME=portfolio_db
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://your-frontend-domain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 7. Performance Requirements

- API response time < 200ms for all endpoints
- Database connection pooling for concurrent requests
- Proper indexing on database queries
- Efficient JSON marshaling/unmarshaling
- Memory-efficient string operations

## 8. Security Requirements

- Input validation on all endpoints
- SQL injection prevention (using prepared statements)
- CORS configuration for frontend domain
- Rate limiting middleware
- Secure headers middleware
- Database connection encryption

## 9. Monitoring and Observability

- Structured logging with correlation IDs
- Health check endpoint with database connectivity
- Metrics collection for GCP monitoring
- Error tracking and alerting
- Database query performance monitoring

[Important Implementation Notes]

1. **NO ORM**: Use raw SQL queries with database/sql package for full control and performance
2. **Clean Architecture**: Separate concerns between handlers, services, and repositories
3. **Schema Integration**: Copy models from shared schemas repository for type consistency
4. **Error Handling**: Implement comprehensive error handling with proper HTTP status codes
5. **Testing**: Write tests for all layers with proper mocking
6. **Documentation**: Keep API documentation in sync with OpenAPI specification
7. **Database Migrations**: Use golang-migrate for version-controlled schema changes
8. **GCP Native**: Leverage GCP services for scalability and reliability

[Development Workflow]

1. Start with database schema and migrations
2. Implement repository layer with raw SQL
3. Build service layer with business logic
4. Create HTTP handlers following OpenAPI spec
5. Add middleware for CORS, logging, recovery
6. Write comprehensive tests
7. Configure deployment for GCP
8. Set up monitoring and logging

[Quality Standards]

- Follow Go idioms and conventions
- Use meaningful variable and function names
- Write self-documenting code with appropriate comments
- Implement proper error handling and logging
- Ensure all code is covered by tests
- Use Go modules for dependency management
- Follow semantic versioning for releases

Remember: Performance and maintainability over convenience. Raw SQL gives us full control over database interactions and optimal performance for the portfolio API.