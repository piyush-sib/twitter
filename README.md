# Twitter Backend API

This repository contains a simplified backend API for a Twitter-like platform. Users can register, log in, post tweets,
follow other users, and view tweets in their feed. The project is built with Go, using GORM for database interactions
and Uber's `dig` for dependency injection.

---

## Features

- **User Management**: Register and log in with secure authentication.
- **Tweets**: Create, retrieve, and manage tweets.
- **Follow System**: Follow other users and view their tweets in your feed.
- **Authentication**: Middleware to secure endpoints using JWT-based authentication.
- **Extensibility**: Built with scalable architecture and dependency injection using Uber's `dig`.
- **Structured Logging**: JSON-based logger for monitoring API usage and debugging.
- **MySQL Integration**: Persistent storage for user and tweet data using GORM ORM.
- **Testing**: Unit-tested services and handlers for robust functionality.

---

## Getting Started

### Prerequisites

Ensure you have the following installed:

- Go 1.19+
- Docker and Docker Compose
- MySQL Server (it will be started using Docker Compose)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/piyush-sib/twitter.git
   cd twitter
   ```

2. Build and start the application dependencies with Docker Compose:

   ```bash
   docker-compose up --build
   ```
3. Now start the application using following command:
   ```bash
   make start
   ```
4. The API should now be available at `http://localhost:8080`.

---

## API Endpoints

### Complete API documentation along with sample request and response is available in Twitter.postman_collection.json file in the root directory.

### You can easily import this file in Postman and test the API.

## Directory Structure

```plaintext
.
├── cmd/
│   └── twitter-backed/                    # Entry point of the application
├── internal/
│   ├── models/                            # GORM models for database tables
│   ├── twitter-backend->repository/       # Database interaction logic
│   ├── twitter-backend->service/          # Business logic and API handlers
│   ├── twitter-backend->utilities/        # Common utilities functions (e.g., password-hashing, jwt etc)
│   ├── twitter-backend->middlewares/      # Middleware logic (Authentication etc)
│   ├── infrastructure/                    # Infra clients (mysql)
│   └── structuredlogger/                  # JSON-based logger
├── Dockerfile                             # Dockerfile for containerization
├── docker-compose.yml                     # Docker Compose setup
└── README.md                              # Project documentation
└── Twitter.postman_collection.json        # API postman collection
```

---

## Testing

Run tests using the following command:

```bash
make test
```

---

## Future Improvements

- Add WebSocket support for real-time updates (e.g., notifications).
- Implement pagination for large tweet and feed responses.
- Add support for media attachments in tweets.
- Improve test coverage with integration and end-to-end tests.
- Add caching with Redis to improve feed performance.

---

## Contributing

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature-name`).
3. Commit your changes (`git commit -m 'Add feature name'`).
4. Push to the branch (`git push origin feature-name`).
5. Open a Pull Request.

---
