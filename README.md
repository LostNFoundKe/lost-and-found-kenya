# Lost and Found Kenya

A national lost and found database for Kenyans to help reconnect people with their lost items and help those who have found items to return them to their rightful owners.

## Project Overview

This open-source platform allows users to:

- Report lost items with detailed descriptions and images
- Report found items to help owners locate them
- Search for items by categories, locations, and keywords
- Claim ownership of found items
- Communicate securely through the platform

## Tech Stack

- **Backend**: Go with Gin framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT
- **Storage**: Google Cloud Storage for images
- **Caching**: Redis (optional)

## Project Structure

```
/lost-and-found-kenya
├── cmd/                  # Application entry points
│   └── api/              # API server
│       └── main.go       # Main application entry point
├── internal/             # Private application code
│   ├── config/           # Configuration handling
│   ├── models/           # Domain models (database models)
│   ├── repository/       # Database access layers
│   ├── service/          # Business logic
│   ├── handler/          # HTTP API handlers
│   ├── middleware/       # HTTP middleware
│   ├── router/           # Route configuration
│   └── util/             # Utility functions
│       └── storage/      # Storage utilities
├── pkg/                  # Public libraries that can be used by external applications
├── migrations/           # Database migrations
├── docs/                 # Documentation
├── scripts/              # Build and deployment scripts
├── .github/              # GitHub workflows for CI/CD
├── .gitignore
├── go.mod
├── go.sum
|── .env
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.23 or higher
- PostgreSQL 14 or higher
- Google Cloud account with Storage access
- Docker and Docker Compose (optional)

### Setup Instructions

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/lost-and-found-kenya.git
cd lost-and-found-kenya
```

2. **Install dependencies**

```bash
go mod download
```

3. **Set up Google Cloud Storage**

- Create a bucket in Google Cloud Storage
- Create a service account with Storage Admin permissions
- Download the service account key JSON file

4. **Set up environment variables**

Create a `.env` file in the root directory:

```
# Server
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info

# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/lostandfound?sslmode=disable

# JWT
JWT_SECRET=your-256-bit-secret
JWT_EXPIRATION_HOURS=24

# Google Cloud Storage
GCS_BUCKET_NAME=lostandfound-kenya
GCS_PROJECT_ID=your-gcp-project-id
GCS_CREDENTIALS_FILE=path/to/credentials.json

# Redis (Optional)
REDIS_URL=redis://localhost:6379/0
```

5. **Set up the database**

```bash
# Create database in PostgreSQL
createdb lostandfound
```

6. **Run the application**

```bash
go run cmd/api/main.go
```

The server will start at http://localhost:9080.

### Running with Docker

```bash
# Build and start containers
docker-compose up -d

# Stop containers
docker-compose down
```

## API Documentation

### Authentication

| Method | Endpoint        | Description        |
|--------|-----------------|-------------------|
| POST   | /api/v1/register | Register new user  |
| POST   | /api/v1/login    | Login to system    |

### Items

| Method | Endpoint          | Description       |
|--------|-------------------|-------------------|
| POST   | /api/v1/items     | Create new item   |
| GET    | /api/v1/items     | List items        |
| GET    | /api/v1/items/:id | Get item by ID    |
| PUT    | /api/v1/items/:id | Update item       |
| DELETE | /api/v1/items/:id | Delete item       |
| POST   | /api/v1/items/:id/images | Upload image |

### Users

| Method | Endpoint          | Description         |
|--------|-------------------|---------------------|
| GET    | /api/v1/users/me  | Get user profile    |
| PUT    | /api/v1/users/me  | Update user profile |

## Contributing

We welcome contributions from the community! Please see our [Contributing Guidelines](docs/CONTRIBUTING.md) for more details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

Project Lead - [Robinson Thairu](mailto:robinsonthairu@gmail.com)

Project Link: [https://github.com/LostNFoundKe/lost-and-found-kenya](https://github.com/LostNFoundKe/lost-and-found-kenya)

## Acknowledgments

* [Gin Web Framework](https://github.com/gin-gonic/gin)
* [GORM](https://gorm.io/)
* [Google Cloud Storage](https://cloud.google.com/storage)
* All contributors who participate in this project