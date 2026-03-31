Shortly_URL URL Shortener

## Project Structure
    Shortly_URL/
      ├── cmd/api/main.go
      ├── internal/
      │   ├── config/
      │   ├── controllers/
      │   ├── middleware/
      │   ├── models/
      │   ├── services/
      │   └── utils/
      ├── docker-compose.yml
      └── README.md

A clean, fast, and secure URL shortener built with **Go (Gin)** and **MongoDB**.  
Each user can create, manage, and track their own shortened links with proper authentication.

---

## ✨ Features

- Lightning-fast URL shortening with custom short codes
- Secure user authentication using JWT
- User-specific link management (only see your own links)
- Click analytics and tracking
- 301 redirects with proper handling
- Fully containerized with Docker
- Persistent data storage
- Clean architecture with proper separation of concerns

---

## 🛠 Tech Stack

- **Backend**: Go + Gin Framework
- **Database**: MongoDB v8
- **Authentication**: JWT + bcrypt
- **Containerization**: Docker & Docker Compose
- **Driver**: Official MongoDB Go Driver v2

---

## 🚀 Quick Start
### 1. Clone the repository
    git.clone git@github.com:Daddy-senpaii/Shortly_URL.git

### 2. Start the application using Docker
    docker compose up -d
### 3. Run the Go server
    go run ./cmd/api
## Api EndPoints
### Authentication
    Method,Endpoint,Description
    POST,/api/register,Register new user
    POST,/api/login,Login and get JWT
### URL designs
    Method,Endpoint,Description,Auth Required
    POST,/api/shorten,Create a new short URL,Yes
    GET,/:code,Redirect to original URL,No
    GET,/api/my-links,Get all my shortened links,Yes
    DELETE,/api/my-links/:id,Delete a shortened link,Yes
