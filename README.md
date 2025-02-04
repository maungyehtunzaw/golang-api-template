# Golang API Template

A **RESTful API** template written in Go (Golang), featuring:

- **Gin** as the HTTP framework  
- **GORM** with MySQL  
- **Redis** for caching/refresh tokens  
- **JWT-based** authentication (access & refresh tokens)  
- **Multi-language (i18n)** support (via embedded JSON files)  
- **Paging** utility for listing resources  
- **AutoMigrate** for database migrations  
- **Vendor** folder (optional) for dependency management  
- A standard **JSON response** format (code, status, message, data)

---

## Table of Contents

1. [Features](#features)  
2. [Project Structure](#project-structure)  
3. [Prerequisites](#prerequisites)  
4. [Getting Started](#getting-started)  
5. [Environment Variables](#environment-variables)  
6. [Running the Server](#running-the-server)  
7. [Running Migrations (Optional)](#running-migrations-optional)  
8. [Available Endpoints](#available-endpoints)  
9. [Using Vendoring (Optional)](#using-vendoring-optional)  
10. [License](#license)

---

## Features

- **User Registration & Authentication**  
  - Stores user credentials (hashed) in MySQL.  
  - Issues **short-lived** access tokens and **long-lived** refresh tokens (stored in Redis).

- **Standard JSON Response**  
  - Returns consistent response objects with `code`, `status`, `message`, and optional `data`.

- **Multi-language (i18n)**  
  - Includes an example of embedding `.json` translation files (e.g., `en.json`, `es.json`).

- **Pagination**  
  - Utility helpers to parse page/limit from query parameters and compute offsets.

- **Auto Migration**  
  - Automatically migrates database tables using GORM’s `AutoMigrate`.

---

## Project Structure

```plaintext
golang-api-template/
├── cmd/
│   ├── server/
│   │   └── main.go              // Entry point for starting the server
│   └── migrate/
│       └── main.go              // (Optional) Command for running DB migrations
├── internal/
│   ├── config/
│   │   ├── config.go            // General app config (env, DB, secrets, etc.)
│   │   └── redis.go             // Redis config & connection
│   ├── handlers/
│   │   ├── auth_handler.go      // Auth endpoints (login, refresh, logout)
│   │   └── user_handler.go      // User endpoints (register, get user, etc.)
│   ├── i18n/
│   │   ├── en.json
│   │   ├── es.json
│   │   └── i18n.go
│   ├── middlewares/
│   │   ├── auth_middleware.go   // JWT validation for protected routes
│   │   └── locale_middleware.go // Detects user locale (Accept-Language)
│   ├── models/
│   │   ├── user.go
│   │   └── product.go
│   ├── repository/
│   │   ├── user_repo.go         // MySQL queries for user
│   │   └── product_repo.go      // MySQL queries for product
│   ├── router/
│   │   └── router.go            // Sets up Gin routes & groups
│   ├── service/
│   │   ├── auth_service.go      // Auth logic (issue tokens, logout, etc.)
│   │   ├── user_service.go      // Business logic for user
│   │   └── product_service.go   // Business logic for product
│   └── utils/
│       └── pagination.go        // Helper for pagination
├── pkg/
│   └── response/
│       └── response.go          // Standardized JSON response format
├── vendor/                      // (Optional) Populated by `go mod vendor`
├── go.mod
└── go.sum
