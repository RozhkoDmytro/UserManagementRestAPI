# User Management REST API

This project implements a REST API for User Management with the following features:
- User creation with email and password
- Public user profiles
- User profile updates (protected by Basic Auth)
- User deletion
- Paginated list of user profiles

## Database Design

The API uses a single table to store user information:

### Table: users

| Column     | Type             | Constraints                                               |
|------------|------------------|-----------------------------------------------------------|
| id         | INT              | PRIMARY KEY, AUTO_INCREMENT                               |
| email      | VARCHAR(255)     | UNIQUE, NOT NULL                                          |
| password   | VARCHAR(255)     |    NOT NULL                                               |  
| first_name | VARCHAR(255)     |                                                           |
| last_name  | VARCHAR(255)     |                                                           |
| created_at | TIMESTAMP        | DEFAULT CURRENT_TIMESTAMP                                 |
| updated_at | TIMESTAMP        | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP     |
| deleted_at | TIMESTAMP        |                                                           |
|------------|------------------|-----------------------------------------------------------|


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

## Security Notes

- User passwords are hashed before storage in the database
- Basic Auth is required for updating user profiles
- User profiles are considered public information

## Getting Started

[Add instructions for setting up and running the project]

## Dependencies

[List any major dependencies or technologies used]

## Contributing

[Add guidelines for contributing to the project]

## License

[Specify the license for your project]
