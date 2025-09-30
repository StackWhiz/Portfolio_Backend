# Portfolio Backend API

A high-performance, scalable backend API for portfolio. Built with Go, featuring Redis caching, PostgreSQL database, and comprehensive REST API endpoints.

## 🚀 Features

- **High Performance**: Built with Go for optimal performance and low latency
- **Scalable Architecture**: Microservices-ready with clean separation of concerns
- **Caching Layer**: Redis integration for fast data retrieval
- **Database**: PostgreSQL with GORM for robust data management
- **Security**: JWT authentication, rate limiting, CORS protection
- **API Documentation**: Swagger/OpenAPI documentation
- **Containerized**: Docker and Docker Compose support
- **Monitoring**: Health checks and structured logging

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Nginx         │    │   API Gateway   │
│   (React/Vue)   │◄──►│   (Reverse      │◄──►│   (Rate Limit,  │
│                 │    │    Proxy)       │    │    CORS, Auth)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                       ┌─────────────────┐              │
                       │   Go API        │◄─────────────┘
                       │   (Gin Router)  │
                       └─────────────────┘
                                │
                ┌───────────────┼───────────────┐
                │               │               │
        ┌───────▼──────┐ ┌──────▼──────┐ ┌─────▼─────┐
        │  PostgreSQL  │ │   Redis     │ │  Logging  │
        │  (Database)  │ │  (Cache)    │ │  System   │
        └──────────────┘ └─────────────┘ └───────────┘
```

## 📋 API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/profile` | Get profile information |
| GET | `/api/v1/experiences` | Get work experiences |
| GET | `/api/v1/skills` | Get technical skills |
| GET | `/api/v1/projects` | Get portfolio projects |
| POST | `/api/v1/contact` | Submit contact form |
| GET | `/health` | Health check |

### Admin Endpoints (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| PUT | `/api/v1/admin/profile` | Update profile |
| POST | `/api/v1/admin/experiences` | Create experience |
| PUT | `/api/v1/admin/experiences/:id` | Update experience |
| DELETE | `/api/v1/admin/experiences/:id` | Delete experience |
| POST | `/api/v1/admin/skills` | Create skill |
| PUT | `/api/v1/admin/skills/:id` | Update skill |
| DELETE | `/api/v1/admin/skills/:id` | Delete skill |
| POST | `/api/v1/admin/projects` | Create project |
| PUT | `/api/v1/admin/projects/:id` | Update project |
| DELETE | `/api/v1/admin/projects/:id` | Delete project |
| GET | `/api/v1/admin/contacts` | Get contact submissions |
| PUT | `/api/v1/admin/contacts/:id/status` | Update contact status |

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | User login |

## 🛠️ Technology Stack

- **Backend**: Go 1.21, Gin Web Framework
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **ORM**: GORM
- **Authentication**: JWT
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker, Docker Compose
- **Monitoring**: Health checks, structured logging

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd stackwhiz_backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

4. **Start PostgreSQL and Redis**
   ```bash
   # Using Docker Compose
   docker-compose up -d postgres redis
   
   # Or start them manually
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

6. **Access the API**
   - API: http://localhost:8080
   - Health Check: http://localhost:8080/health
   - Swagger Docs: http://localhost:8080/swagger/index.html

### Docker Deployment

1. **Build and run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

2. **Check service status**
   ```bash
   docker-compose ps
   ```

3. **View logs**
   ```bash
   docker-compose logs -f api
   ```

## 📊 Database Schema

### Profile
- Personal information, contact details, professional summary

### Experience
- Work history with achievements, technologies, and time periods

### Skills
- Technical skills categorized by type (Languages, Frameworks, Tools, etc.)

### Projects
- Portfolio projects with descriptions, technologies, and links

### Contact
- Contact form submissions with status tracking

### User
- Admin users for content management

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Environment (development/production) | development |
| `DATABASE_URL` | PostgreSQL connection string | postgres://user:password@localhost:5432/portfolio_db |
| `REDIS_URL` | Redis connection string | redis://localhost:6379 |
| `JWT_SECRET` | JWT signing secret | your-secret-key-change-in-production |
| `PORT` | Server port | 8080 |
| `RATE_LIMIT` | Requests per second limit | 100 |

### Database Configuration

The application uses GORM for database operations with automatic migrations. The database schema is created automatically on startup.

### Redis Configuration

Redis is used for caching API responses to improve performance. Cache keys include:
- `profile` - Profile information
- `experiences` - Work experiences
- `skills` - Technical skills
- `projects` - Portfolio projects
- `projects:featured` - Featured projects only

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication for admin endpoints
- **Rate Limiting**: Configurable rate limiting to prevent abuse
- **CORS Protection**: Configurable CORS policies
- **Security Headers**: XSS protection, content type sniffing prevention
- **Input Validation**: Request validation using Gin's binding
- **SQL Injection Protection**: GORM provides protection against SQL injection

## 📈 Performance Features

- **Redis Caching**: Reduces database load and improves response times
- **Connection Pooling**: Optimized database connection management
- **Structured Logging**: Efficient logging with structured data
- **Health Checks**: Built-in health monitoring
- **Graceful Shutdown**: Proper cleanup on application termination

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestProfileService ./internal/service
```

## 📝 API Documentation

The API documentation is automatically generated using Swagger/OpenAPI. Access it at:
- Development: http://localhost:8080/swagger/index.html
- Production: https://yourdomain.com/swagger/index.html

## 🚀 Deployment

### Production Deployment

1. **Set production environment variables**
   ```bash
   export ENVIRONMENT=production
   export DATABASE_URL=postgres://user:password@db-host:5432/portfolio_db
   export REDIS_URL=redis://redis-host:6379
   export JWT_SECRET=your-super-secure-jwt-secret
   ```

2. **Build for production**
   ```bash
   CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
   ```

3. **Run with process manager**
   ```bash
   # Using systemd, PM2, or similar
   ./main
   ```

### Docker Production Deployment

```bash
# Build production image
docker build -t portfolio-api:latest .

# Run with production configuration
docker run -d \
  --name portfolio-api \
  -p 8080:8080 \
  -e ENVIRONMENT=production \
  -e DATABASE_URL=postgres://user:password@db-host:5432/portfolio_db \
  -e REDIS_URL=redis://redis-host:6379 \
  -e JWT_SECRET=your-super-secure-jwt-secret \
  portfolio-api:latest
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

** StackWhize **
- GitHub: [@StackWhiz](https://github.com/StackWhiz)
- Telegram: [@galacticdot](https://t.me/galacticdot)

## 🙏 Acknowledgments

- Built with Go and the amazing Go ecosystem
- Inspired by modern backend architecture patterns
- Thanks to the open-source community for the excellent tools and libraries
