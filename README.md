# Restful API with Gin Framework

This project is a RESTful API built with the Gin framework in Go. It provides CRUD operations for managing articles with support for search, pagination, and logging middleware for enhanced request tracking.

---

## Featrures
- Create, Update, Search, and Get All Articles.
- Middleware to log requests and responses.
- Pagination support for get all articles.
- Modular design for scalability and maintainability.

---

## Prerequisites
Before running the project, ensure you have the following installed:
- Go (version 1.20 or later)
- Git

--- 

## Installation
### 1. Clone the Repository
```sh
git clone https://github.com/brothergiez/restful-api.git
cd restful-api
```

### 2. Install Dependencies

```sh
go mod tidy
```

---

## Configuration
The project uses .env for environment variables. Create a .env file in the root directory with the following content:

```env
APP_PORT=3000
```

or copy from .env.example file by running command 
```sh
cp .env .env.example
```

If no .env file is found, the application will default to port 8080.

---

## Running the Application
Start the server by running:

```sh
go run main.go
```

The server will start on the port specified in the .env file (or 8080 by default).

---

## API Endpoints
### Base URL
```
http://localhost:<APP_PORT>
```

### Endpoints
| Method | Endpoint | Description |
| --- | --- | --- | 
| POST | /articles/create | Create a new article. |
| PUT | /articles/update/:id | Update an article by ID. |
| GET | /articles/search | Search articles by keyword. |
| GET | /articles/get-all | Retrieve articles with pagination. |


---

## Example Usage
### Create an Article

Request:
```sh
curl -X POST http://localhost:8080/articles/create \
-H "Content-Type: application/json" \
-d '{"title": "Learn Go", "content": "Go is an awesome language."}'
```

Response:
```json
{
  "id": 1,
  "title": "Learn Go",
  "content": "Go is an awesome language."
}
```
---
### Update Article

Request :
```sh
curl -X PUT http://localhost:8080/articles/update/1 \
-H "Content-Type: application/json" \
-d '{"title": "Updated Title", "content": "Updated content of the article."}'
```

Response :
```json
{
  "id": 1,
  "title": "Updated Title",
  "content": "Updated content of the article."
}

```
---

### Search Article

Request : 
```sh
curl -X GET "http://localhost:8080/articles/search?keyword=Go"
```

Response : 
```json
[
  {
    "id": 1,
    "title": "Go Programming Basics",
    "content": "Learn the basics of Go programming."
  },
  {
    "id": 2,
    "title": "Advanced Go Patterns",
    "content": "Explore advanced patterns in Go."
  }
]

```

---

### Get All Articles (with pagination)

Request : 
```sh
curl -X GET "http://localhost:8080/articles/get-all?page=1&limit=5"
```

Response :
```json
{
  "page": 1,
  "limit": 5,
  "total": 10,
  "totalPages": 2,
  "articles": [
    {
      "id": 1,
      "title": "Learn Go",
      "content": "Go is an awesome language."
    },
    ...
  ]
}

```

--- 

## Testing

The project includes unit tests for handlers, middleware, routes, and repository layers.

### Run All Tests
```sh
go test ./... -v
```

### Run Tests for a Specific Package
```sh
go test ./handlers
```

---

## License
This project is licensed under the MIT License.