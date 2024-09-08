# User Management REST API

This project implements a REST API for User Management with the following features:
- User creation with email and password
- Public user profiles
- User profile updates (protected by Basic Auth)
- User deletion
- Paginated list of user profiles
- Voting
  
## Dependencies

- Go: Programming language used for the API.
- gorilla/mux: Router for handling HTTP requests.
- gorm: ORM for working with the PostgreSQL database.
- Redis: Caching layer for storing session information.
- bcrypt: For hashing passwords securely.
- docker-compose: For container orchestration.

## Database Design

The API uses a single table to store user information:

### Table: users

| Column     | Type             | Constraints                                               |
|------------------|------------------|-----------------------------------------------------------|
| id               | INT              | PRIMARY KEY, AUTO_INCREMENT                               |
| email            | VARCHAR(255)     | UNIQUE, NOT NULL                                          |
| password         | VARCHAR(255)     |    NOT NULL                                               |  
| first_name       | VARCHAR(255)     |                                                           |
| last_name        | VARCHAR(255)     |                                                           |
| created_at       | TIMESTAMP        | DEFAULT CURRENT_TIMESTAMP                                 |
| updated_at       | TIMESTAMP        | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP     |
| deleted_at       | TIMESTAMP        |                                                           |
| vote_updated_at  | TIMESTAMP        |                                                           |
| rating           | INT              |                                                           |


## API Endpoints

### Create User
- **URL:** `/user`
- **Method:** POST
- **Request Body:**
  ```json
  {
    "email"     : "string",
    "password"  : "string",
    "first_name": "string",
    "last_name" : "string",
    "nick_name" : "string"
  }

- Response: 201 Created with the created user ID

### Get User Profile
- **URL:** `/user/{id}`
- **Method:** GET
- **Response:**
  ```json
  {
    "id": "integer",
    "email": "string",
    "first_name": "string",
    "last_name": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```

### Update User Profile
- **URL:** `/users/{id}`
- **Method:** PUT
- **Authentication:** Basic Auth
- **Request Body:**
  ```json
  {
    "first_name": "string",
    "last_name": "string",
    "password": "string"
  }
  ```
- **Response:** 200 OK

### Delete User
- **URL:** `/users/{id}`
- **Method:** DELETE
- **Response:** 204 No Content

### List Users with Pagination
- **URL:** `/users`
- **Method:** GET
- **Query Parameters:** 
  - `page` (default: 1)
  - `limit` (default: 10)
- **Response:**
  ```json
  {
    "users": [
      {
        "id": "integer",
        "email": "string",
        "first_name": "string",
        "last_name": "string"
      },
      ...
    ],
    "total": "integer",
    "page": "integer",
    "limit": "integer"
  }
  ```
### Like User
- **URL:** `/user/like/{id}`
- **Method:** POST
- **Description:** Adds one like to a user's profile.
- **Response:**
  ```json
  {
  "message": "User liked successfully"
  }
  ```
### Dislike User
- **URL:** `/user/dislike/{id}`
- **Method:** POST
- **Description:** Adds one dislike to a user's profile.
- **Response:**
  ```json
  {
  "message": "User disliked successfully"
  }
  ```
### Revoke
- **URL:** `/user/revoke/{id}`
- **Method:** DELETE
- **Description:** Revokes any vote (like or dislike) from a user's profile.
- **Response:**
  ```json
  {
  "message": "Vote revoked successfully"
  }
  ```
  
## Security Notes

- User passwords are hashed before storage in the database
- Basic Auth is required for updating user profiles
- User profiles are considered public information

## Getting Started
- Prerequisites
- Docker (for containerized setup)
- Go 1.16+ (for local setup)
- PostgreSQL
- Redis
- Docker

## Running with Docker
Ensure you have Docker installed.
Run the following command to start up the containers:
```
docker-compose up --build
```
The API will be accessible at http://localhost:50052. (or http://localhost:[APP_PORT])

## License

This project is licensed under the MIT License
