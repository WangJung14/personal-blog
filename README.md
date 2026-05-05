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

## Deployment

The easiest way to deploy this application to a VPS is using Docker Compose.

1.  **Clone the repository to your server**
2.  **Ensure Docker is installed**
3.  **Deploy everything with one command**:
    ```bash
    docker-compose up -d --build
    ```
    This command will:
    - Spin up a PostgreSQL database.
    - Build the Go application into a lightweight Docker image.
    - Start the application and link it to the database.
    - Persist your database data and uploaded images even if the containers are restarted.

### Deploy to Render

Render is a great platform for hosting Go applications.

1.  **Push your code to GitHub**.
2.  **Go to [Render Dashboard](https://dashboard.render.com/)**.
3.  **Click "New" -> "Blueprint"**.
4.  **Connect your GitHub repository**.
5.  Render will automatically detect the `render.yaml` file and set up:
    - A managed PostgreSQL database.
    - Your Go web application.
    - Automatic environment variable configuration.

> [!WARNING]
> **Note on Image Uploads**: Render's free tier uses ephemeral storage. This means images uploaded to the blog will be deleted whenever the service restarts. For production use, it is recommended to use an external storage service like Cloudinary or AWS S3, or attach a [Render Disk](https://render.com/docs/disks) (paid feature).
