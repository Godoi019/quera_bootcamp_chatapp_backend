# Chat App Backend

A real-time chat application backend built with Go, Fiber v3, WebSocket, and PostgreSQL with clean architecture and service layer pattern.

## Features

- **User Authentication**: Username/password login with PASETO v4 tokens
- **User Management**: Full CRUD operations for user profiles
- **Chat Management**: Create group chats and direct messages
- **Real-time Messaging**: WebSocket support for instant messaging
- **Message Management**: Send, edit, delete messages
- **Member Management**: Add/remove members from group chats
- **Clean Architecture**: Service layer separates business logic from HTTP handlers
- **CLI Interface**: Cobra-based command-line interface with config management

## Tech Stack

- **Go 1.25+**: Programming language
- **Fiber v3**: High-performance web framework
- **Ent**: Type-safe ORM for database operations
- **PostgreSQL**: Relational database
- **PASETO v4**: Modern secure token-based authentication
- **WebSocket**: Real-time bidirectional communication
- **Cobra**: CLI framework for command-line interface
- **Viper**: Configuration management with YAML support
- **Docker Compose**: Container orchestration
- **Air**: Hot reload for development
- **Makefile**: Build automation

## Architecture

This project follows clean architecture principles with clear separation of concerns:

```
┌─────────────┐
│   cmd/      │  → CLI commands & initialization
└──────┬──────┘
       │
┌──────▼──────────────────────────┐
│   internal/                     │
│   ├── server/    → HTTP setup   │
│   ├── handler/   → HTTP layer   │
│   ├── service/   → Business     │
│   ├── auth/      → Auth logic   │
│   └── repository/→ Data access  │
└─────────────────────────────────┘
```

## Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose
- Make

## Project Structure

```
.
├── cmd/
│   ├── server/          # Cobra commands (start, etc.)
│   └── system/          # System maintenance commands
├── config/
│   ├── config.go        # Config type definitions
│   ├── read.go          # Config loading logic
│   └── defaults.go      # Default configuration values
├── internal/
│   ├── auth/            # Authentication service (PASETO v4)
│   ├── handler/         # HTTP and WebSocket handlers (presentation layer)
│   ├── middleware/      # HTTP middleware (auth, logging, etc.)
│   ├── model/           # Request/response models & DTOs
│   ├── server/          # HTTP server setup & routing
│   ├── service/         # Business logic layer
│   │   ├── user.go      # User business logic
│   │   ├── chat.go      # Chat business logic
│   │   └── message.go   # Message business logic
│   └── repository/      # Data access layer
│       ├── ent/         # Generated Ent code
│       └── schema/      # Ent schema definitions
├── pkg/
│   ├── constants/       # Application constants
│   ├── fiber/           # Fiber utilities & validation
│   └── utils/           # Helper functions
├── config.yaml          # Configuration file
├── config.sample.yaml   # Configuration template
├── docker-compose.yml   # PostgreSQL setup
├── Makefile            # Build commands
└── .air.toml           # Air hot reload config
```

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Hossara/quera_bootcamp_chatapp_backend.git
cd quera_bootcamp_chatapp_backend
```

### 2. Set up configuration

```bash
cp config.sample.yaml config.yaml
# Edit config.yaml with your settings
```

### 3. Start PostgreSQL with Docker Compose

```bash
make docker-up
```

### 4. Install dependencies

```bash
go mod download
```

### 5. Run the application

For development with hot reload:
```bash
make dev
```

Or run directly:
```bash
make run
# Or with custom config path:
make run CONFIG_PATH=/path/to/config
```

Using the CLI directly:
```bash
# Start server with default config (current directory)
./bin/chatapp-server server start

# Start server with custom config path
./bin/chatapp-server server start --config /path/to/config
./bin/chatapp-server server start -c ./configs

# View help
./bin/chatapp-server --help
./bin/chatapp-server server start --help
```

## Available Make Commands

```bash
make help          # Show all available commands
make run           # Run the application (CONFIG_PATH=. by default)
make build         # Build the application
make test          # Run tests
make clean         # Clean build artifacts
make docker-up     # Start docker containers
make docker-down   # Stop docker containers
make docker-logs   # View docker logs
make ent-generate  # Generate Ent code
make dev           # Run with Air hot reload
make deps          # Download dependencies
```

## API Endpoints

### Health Check

- `GET /health` - Application health status

### Authentication

- `POST /api/v1/auth/register` - Register a new user
  - Body: `{ "username": "string", "password": "string", "display_name": "string" }`
  - Supports both JSON and form data
- `POST /api/v1/auth/login` - Login user
  - Body: `{ "username": "string", "password": "string" }`
  - Returns PASETO v4 token
- `GET /api/v1/auth/me` - Get current user (authenticated)

### Users

- `GET /api/v1/users?limit=50&offset=0` - List all users
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user (own profile only)
  - Body: `{ "display_name": "string", "password": "string" }`
- `DELETE /api/v1/users/:id` - Delete user (own profile only)
- `POST /api/v1/users/last-seen` - Update last seen timestamp

### Chats

- `POST /api/v1/chats` - Create a new chat
  - Body: `{ "name": "string", "is_group": boolean, "member_ids": [int] }`
- `GET /api/v1/chats?limit=50&offset=0` - List user's chats with members
- `GET /api/v1/chats/:id` - Get chat details with members
- `PUT /api/v1/chats/:id` - Update chat name (admin only)
  - Body: `{ "name": "string" }`
- `DELETE /api/v1/chats/:id` - Delete chat (creator only)
- `POST /api/v1/chats/:id/members` - Add members to chat (admin only)
  - Body: `{ "member_ids": [int] }`
- `DELETE /api/v1/chats/:id/members/:memberId` - Remove member (admin only)

### Messages

- `POST /api/v1/messages` - Send a message
  - Body: `{ "chat_id": int, "content": "string" }`
- `GET /api/v1/messages/:id` - Get message by ID
- `GET /api/v1/messages/chat/:chatId?limit=50&offset=0` - List messages in chat
- `PUT /api/v1/messages/:id` - Update message (own message only)
  - Body: `{ "content": "string" }`
- `DELETE /api/v1/messages/:id` - Delete message (own message only)

### WebSocket

- `GET /ws?token=<auth_token>` - WebSocket connection for real-time chat
- `GET /ws/health` - WebSocket health check and statistics

## Architecture Details

### Service Layer

The application uses a service layer pattern to separate business logic from HTTP handlers:

- **UserService**: User management, authentication, token generation
- **ChatService**: Chat creation, membership management, permissions
- **MessageService**: Message CRUD operations, sender verification

Benefits:
- Clean separation of concerns
- Easier testing and mocking
- Reusable business logic
- No database queries in handlers

### Validation

All request validation uses a consistent pattern with `pkg/fiber` utilities:

```go
req := new(model.LoginRequest)
if err := f.ParseRequestBody(c, req); err != nil {
    return f.RespondError(c, fiber.StatusBadRequest, err.Message, err.Errors)
}
```

Features:
- Automatic JSON and form data parsing
- Struct tag validation (go-playground/validator)
- Standardized error responses
- Field-level error messages

### Authentication

Uses PASETO v4 (Platform-Agnostic Security Tokens):
- Symmetric encryption (V4.local)
- No algorithm confusion vulnerabilities
- Built-in expiration handling
- Secure by default

Token structure:
```json
{
  "user_id": 1,
  "username": "john_doe",
  "issued_at": "2024-01-01T12:00:00Z",
  "expire_at": "2024-01-02T12:00:00Z"
}
```

## WebSocket Protocol

### Connection

Connect to WebSocket endpoint with authentication token:
```
ws://localhost:8080/ws?token=YOUR_AUTH_TOKEN
```

### Message Types

#### Send Message
```json
{
  "type": "message",
  "payload": {
    "chat_id": 1,
    "content": "Hello, World!"
  }
}
```

#### Join Chat Room
```json
{
  "type": "join_chat",
  "payload": {
    "chat_id": 1
  }
}
```

#### Leave Chat Room
```json
{
  "type": "leave_chat",
  "payload": {
    "chat_id": 1
  }
}
```

### Receiving Messages

```json
{
  "type": "message",
  "payload": {
    "message_id": 123,
    "content": "Hello, World!",
    "sender_id": 1,
    "username": "john_doe",
    "chat_id": 1,
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### System Messages

```json
{
  "type": "system",
  "payload": {
    "chat_id": 1,
    "message": "User joined the chat"
  }
}
```

## Authentication

All protected endpoints require an `Authorization` header:

```
Authorization: Bearer <your_paseto_token>
```

Tokens are generated using PASETO v4 and expire after 24 hours (configurable in `config.yaml`).

## Database Schema

### User
- `id`: Primary key
- `username`: Unique username (3-50 characters)
- `password`: Bcrypt hashed password
- `display_name`: Display name
- `created_at`: Creation timestamp
- `updated_at`: Update timestamp
- `last_seen`: Last activity timestamp

### Chat
- `id`: Primary key
- `name`: Chat name (max 100 characters)
- `is_group`: Boolean for group chat vs direct message
- `creator_id`: Foreign key to User
- `created_at`: Creation timestamp
- `updated_at`: Update timestamp

### Message
- `id`: Primary key
- `content`: Message content (text)
- `sender_id`: Foreign key to User
- `chat_id`: Foreign key to Chat
- `is_edited`: Whether message was edited
- `created_at`: Creation timestamp
- `updated_at`: Update timestamp

### ChatMember
- `user_id`: Foreign key to User (composite primary key)
- `chat_id`: Foreign key to Chat (composite primary key)
- `is_admin`: Admin privileges in chat
- `joined_at`: Join timestamp

## Development

### Hot Reload

The project uses Air for hot reload during development:

```bash
make dev
```

Air configuration (`.air.toml`) automatically:
- Watches Go files for changes
- Rebuilds the application
- Restarts the server
- Passes `--config .` flag

### Adding New Ent Schemas

1. Create a new schema:
```bash
cd internal/repository
go run -mod=mod entgo.io/ent/cmd/ent new YourSchema
```

2. Edit the schema file in `internal/repository/schema/yourschema.go`

3. Generate Ent code:
```bash
make ent-generate
```

### Adding New Services

1. Create a new service file in `internal/service/`:
```go
package service

type YourService struct {
    client *ent.Client
}

func NewYourService(client *ent.Client) *YourService {
    return &YourService{client: client}
}
```

2. Add business logic methods

3. Use the service in handlers:
```go
type YourHandler struct {
    yourService *service.YourService
}
```

### Code Structure Best Practices

- **Handlers**: Only handle HTTP concerns (parsing, validation, responses)
- **Services**: Contain all business logic and database operations
- **Models**: Define request/response structures with validation tags
- **Middleware**: Cross-cutting concerns (auth, logging, recovery)

## Testing

Run tests:
```bash
make test
```

Run tests with coverage:
```bash
go test -v -cover ./...
```

## Building for Production

Build the binary:
```bash
make build
```

The binary will be created in `bin/chatapp-server`.

Run in production:
```bash
# With config in current directory
./bin/chatapp-server server start

# With custom config path
./bin/chatapp-server server start --config /etc/chatapp

# With environment variable
export CONFIG_PATH=/etc/chatapp
./bin/chatapp-server server start
```

## Docker Support

### Development with Docker Compose

Start the PostgreSQL database:
```bash
make docker-up
```

Stop the database:
```bash
make docker-down
```

View logs:
```bash
make docker-logs
```

### Docker Compose Services

- **postgres**: PostgreSQL 15 database
  - Port: 5432
  - Database: chatapp_db
  - User: chatapp
  - Password: chatapp123

## Environment Variables

The application supports the following environment variables:

- `CONFIG_PATH`: Path to config directory (default: ".")
- `PORT`: Override server port
- `DATABASE_URL`: PostgreSQL connection string

## Performance Tips

1. **Connection Pooling**: PostgreSQL connections are pooled by Ent
2. **WebSocket**: Uses efficient FastHTTP WebSocket implementation
3. **Validation**: Request validation happens at handler level
4. **Service Layer**: Reduces code duplication and improves maintainability

## Security Features

- **PASETO v4**: Secure token authentication
- **Bcrypt**: Password hashing with salt
- **SQL Injection**: Protected by Ent ORM
- **CORS**: Configurable CORS middleware
- **Rate Limiting**: Ready for implementation via Fiber middleware

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
make docker-logs

# Verify database credentials in config.yaml
cat config.yaml
```

### Port Already in Use

```bash
# Change port in config.yaml
# Or use environment variable
PORT=8081 make run
```

### WebSocket Connection Issues

- Ensure token is valid and not expired
- Check CORS settings for cross-origin requests
- Verify WebSocket endpoint: `ws://localhost:8080/ws?token=TOKEN`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Follow the existing code structure (handlers → services → repository)
4. Add tests for new features
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Code Style

- Follow Go conventions and idioms
- Use meaningful variable names
- Add comments for complex logic
- Keep functions focused and small
- Use service layer for business logic

## License

MIT

## Authors

- **Hossara** - Initial work

## Acknowledgments

- Fiber v3 team for the excellent web framework
- Ent team for the powerful ORM
- PASETO specification for secure tokens
- Go community for best practices and patterns
