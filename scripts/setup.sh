#!/bin/bash

# StackWhiz Portfolio Backend Setup Script
# This script sets up the development environment

set -e

echo "🚀 Setting up StackWhiz Portfolio Backend..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.21 or later.${NC}"
    echo "Visit: https://golang.org/doc/install"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo -e "${RED}❌ Go version $GO_VERSION is too old. Please install Go 1.21 or later.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go version $GO_VERSION is installed${NC}"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}⚠️  Docker is not installed. You'll need it for the database and Redis.${NC}"
    echo "Visit: https://docs.docker.com/get-docker/"
else
    echo -e "${GREEN}✅ Docker is installed${NC}"
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${YELLOW}⚠️  Docker Compose is not installed. You'll need it for easy setup.${NC}"
    echo "Visit: https://docs.docker.com/compose/install/"
else
    echo -e "${GREEN}✅ Docker Compose is installed${NC}"
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo -e "${BLUE}📝 Creating .env file...${NC}"
    cp env.example .env
    echo -e "${GREEN}✅ .env file created${NC}"
    echo -e "${YELLOW}⚠️  Please review and update the .env file with your configuration${NC}"
else
    echo -e "${GREEN}✅ .env file already exists${NC}"
fi

# Install Go dependencies
echo -e "${BLUE}📦 Installing Go dependencies...${NC}"
go mod download
go mod tidy
echo -e "${GREEN}✅ Go dependencies installed${NC}"

# Install development tools
echo -e "${BLUE}🛠️  Installing development tools...${NC}"
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/swaggo/swag/cmd/swag@latest
echo -e "${GREEN}✅ Development tools installed${NC}"

# Start database and Redis with Docker Compose
if command -v docker-compose &> /dev/null; then
    echo -e "${BLUE}🐳 Starting database and Redis with Docker Compose...${NC}"
    docker-compose up -d postgres redis
    
    # Wait for services to be ready
    echo -e "${BLUE}⏳ Waiting for services to be ready...${NC}"
    sleep 10
    
    # Check if services are running
    if docker-compose ps | grep -q "Up"; then
        echo -e "${GREEN}✅ Database and Redis are running${NC}"
    else
        echo -e "${RED}❌ Failed to start database and Redis${NC}"
        echo "Please check the logs: docker-compose logs"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠️  Docker Compose not available. Please start PostgreSQL and Redis manually.${NC}"
fi

# Generate API documentation
echo -e "${BLUE}📚 Generating API documentation...${NC}"
swag init -g main.go -o ./docs
echo -e "${GREEN}✅ API documentation generated${NC}"

# Run tests
echo -e "${BLUE}🧪 Running tests...${NC}"
go test ./...
echo -e "${GREEN}✅ Tests passed${NC}"

# Build the application
echo -e "${BLUE}🔨 Building the application...${NC}"
go build -o bin/stackwhiz-portfolio-backend main.go
echo -e "${GREEN}✅ Application built successfully${NC}"

echo ""
echo -e "${GREEN}🎉 Setup completed successfully!${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo -e "1. Review and update the .env file with your configuration"
echo -e "2. Start the application: ${YELLOW}make run${NC} or ${YELLOW}go run main.go${NC}"
echo -e "3. Access the API: ${YELLOW}http://localhost:8080${NC}"
echo -e "4. View API docs: ${YELLOW}http://localhost:8080/swagger/index.html${NC}"
echo -e "5. Health check: ${YELLOW}http://localhost:8080/health${NC}"
echo ""
echo -e "${BLUE}Available commands:${NC}"
echo -e "• ${YELLOW}make help${NC} - Show all available commands"
echo -e "• ${YELLOW}make run${NC} - Run the application"
echo -e "• ${YELLOW}make test${NC} - Run tests"
echo -e "• ${YELLOW}make docker-compose-up${NC} - Start all services"
echo -e "• ${YELLOW}make docker-compose-down${NC} - Stop all services"
echo ""
echo -e "${GREEN}Happy coding! 🚀${NC}"
