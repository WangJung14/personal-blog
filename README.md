# Personal Blog

A lightweight, robust, and modern personal blog built with Go and PostgreSQL.

## Features

- **Blazing Fast**: Powered by Go and `go-chi` router.
- **Modern Interface**: Clean public-facing design and a secure Admin Dashboard.
- **Rich Text Editor**: Integrated with Quill.js for a seamless content authoring experience.
- **Secure Authentication**: Stateless session management using signed JWTs.
- **Automated Slug Generation**: Automatically generates URL-friendly slugs for posts, handling duplicates gracefully.
- **Containerized Database**: Easy-to-start local development using Docker Compose.

## Tech Stack

- **Backend**: Go (latest)
- **Router**: `github.com/go-chi/chi/v5`
- **Database**: PostgreSQL (`github.com/lib/pq`)
- **Authentication**: JWT (`github.com/golang-jwt/jwt/v5`)
- **Frontend**: Go `html/template`, raw CSS, JS.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/)
- [Docker & Docker Compose](https://www.docker.com/)

### Installation & Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/WangJung14/personal-blog.git
   cd personal-blog
   ```

2. **Configure Environment Variables**
   The project requires a `.env` file in the root directory. You can create one with the following default values:
   ```env
   PORT=8080
   ADMIN_USER=admin
   ADMIN_PASS=admin123
   SECRET_KEY=supersecretkey123
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/blog?sslmode=disable
   ```

3. **Start the Database**
   Run the provided Docker Compose configuration to start a local PostgreSQL instance:
   ```bash
   docker-compose up -d
   ```

4. **Run the Application**
   ```bash
   go run main.go
   ```
   *The application will automatically connect to the database and run the necessary schema migrations on startup.*

### Usage

- **Public Site**: Open [http://localhost:8080](http://localhost:8080) to view the blog.
- **Admin Panel**: Go to [http://localhost:8080/admin](http://localhost:8080/admin) and log in using the credentials defined in your `.env` file (default: `admin` / `admin123`). From there, you can create, edit, and delete posts.
