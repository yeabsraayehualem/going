# going

A personal web framework for Go, inspired by Django's project structure and workflow.

## Features

- **Project Structure**: Follows Django's project and app structure
- **Configuration**: YAML-based configuration
- **Database**: SQLite support out of the box
- **Authentication**: Session-based authentication
- **Password Hashing**: Secure password hashing using Argon2id
- **Routing**: Flexible routing with gorilla/mux
- **ORM**: Database operations with GORM

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yeabsraayehualem/going.git
   cd going
   ```

2. Initialize a new project:
   ```bash
   go run cmd/going/main.go -init
   ```

3. Start the development server:
   ```bash
   go run cmd/going/main.go
   ```

   The server will start on `http://localhost:8080` by default.

## Project Structure

```
going/
├── cmd/
│   └── going/          # Main application entry point
├── config/              # Configuration files
├── internal/
│   ├── app/             # Application core
│   ├── auth/            # Authentication and authorization
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and models
│   ├── middleware/      # HTTP middleware
│   └── session/         # Session management
├── pkg/
│   └── utils/           # Utility functions
├── apps/                # Your applications go here
├── migrations/          # Database migrations
├── static/              # Static files (CSS, JS, images)
└── templates/           # HTML templates
```

## Creating a New App

1. Create a new directory under `apps/` for your app:
   ```bash
   mkdir -p apps/myapp
   ```

2. Create a `models.go` file in your app directory:
   ```go
   package myapp

   import "going/internal/database"

   type MyModel struct {
       ID   uint   `gorm:"primaryKey"`
       Name string `gorm:"size:255"`
   }

   func init() {
       database.RegisterModels(&MyModel{})
   }
   ```

## Configuration

Edit `config/config.yaml` to configure your application:

```yaml
# Database configuration
database:
  driver: sqlite3
  name: going.db
  path: ./data
  log_level: info

# Server configuration
server:
  host: 0.0.0.0
  port: 8080

# Session configuration
session:
  name: going_session
  secret: change-this-to-a-secure-secret-key
  lifetime: 120  # Session lifetime in minutes (2 hours)
```

## Authentication

### Password Hashing

```go
import "going/internal/auth"

// Hash a password
hashedPassword, err := auth.HashPassword("mysecurepassword")
if err != nil {
    // handle error
}

// Verify a password
match, err := auth.VerifyPassword("mysecurepassword", hashedPassword)
if err != nil {
    // handle error
}
if match {
    // password is correct
}
```

### Session Management

```go
// Create a new session
session := app.Session.CreateSession()
session.Values["user_id"] = userID

// Set session cookie
app.Session.SetSessionCookie(w, session.ID)

// Get session from request
session, err := app.Session.GetSessionFromRequest(r)
if err != nil {
    // handle error (no session or session expired)
}

// Get user ID from session
userID, ok := session.Values["user_id"].(int)
if !ok {
    // handle error (invalid user ID)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
