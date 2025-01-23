# Gorbit API Framework

## Overview
Gorbit is a comprehensive, batteries-included API framework for Go, designed to simplify and accelerate API development with multi-database support and robust configuration management.

## Features
- Multi-database support (MySQL, PostgreSQL, MongoDB)
- Redis caching
- Docker containerization
- Swagger documentation
- Flexible configuration management
- Modular project structure

## Prerequisites
- Go 1.21+
- Docker
- Docker Compose

## Getting Started

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/Gorbit.git
cd Gorbit
```

### 2. Replace Module Name
Replace `your-module-name` in all files with your actual module name:
```bash
find . -type f -name "*.go" -exec sed -i 's/your-module-name/your-actual-module-name/g' {} +
```

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Build and Run
```bash
# Build Docker containers
make build

# Start the application
make run
```

## Development Commands
- `make build`: Build Docker containers
- `make run`: Start the application
- `make stop`: Stop the containers
- `make clean`: Remove all containers and volumes
- `make test`: Run application tests
- `make swagger`: Regenerate Swagger documentation

## Configuration
Configuration files are located in the `configs/` directory:
- `config.yaml`: Server settings
- `database.yaml`: Database configurations
- `redis.yaml`: Redis connection details

## Documentation
Access Swagger UI at: `http://localhost:8080/swagger/index.html`

## License
[Choose an appropriate license]