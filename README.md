# Relay Hub

Relay Hub is a lightweight, self-hosted HTTP API testing workspace and client written in Go. It enables developers to execute, organize, and inspect HTTP requests with built-in user authentication, request collections, and request history logs.

## Features

- **HTTP Request Execution**: Proxy HTTP requests (GET, POST, PUT, DELETE, PATCH) with custom headers, query params, and JSON request bodies.
- **User Authentication**: Multi-tenant workspace support with user registration, secure bcrypt password hashing, and token-based session management.
- **Request Collections & Saved Requests**: Organize API requests into named collections and reusable templates.
- **Request History & Metrics**: Automatic request logging with status codes, response headers, body formatting, and high-precision execution timing (`elapsed_ms`).
- **Single Binary & Embedded UI**: Simple web interface served directly from the Go server.

## Architecture

- **Language & Runtime**: Go
- **Database**: MySQL (Auto-migrating schema for users, sessions, collections, saved requests, and history)


## API Endpoints

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `/api/auth/register` | `POST` | Register a new user account |
| `/api/auth/login` | `POST` | Authenticate and receive a session token |
| `/api/auth/logout` | `POST` | Invalidate current session token |
| `/api/auth/me` | `GET` | Fetch authenticated user profile |
| `/api/request` | `POST` | Execute proxy HTTP request and record history |
| `/api/collections` | `GET` / `POST` | List or create request collections |
| `/api/saved-requests` | `GET` / `POST` | List or save request templates |
| `/api/saved-requests/:id` | `DELETE` | Remove a saved request |
| `/api/history` | `GET` / `DELETE` | List recent request logs or clear history |
| `/api/history/:id` | `DELETE` | Delete a single history record |
