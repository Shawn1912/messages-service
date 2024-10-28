# Messages Service

## Description

A simple RESTful API for managing messages and determining if they are palindromes.

## Project Structure

messages-service/ <br />
├── main.go  <br />
├── handlers.go  <br />
├── models.go  <br />
├── db_connection.go  <br />
├── utils <br /> &emsp;&emsp;
    └── palindrome.go <br />
└── go.mod

## Build and Run

### Prerequisites

- Go 1.16+
- PostgreSQL

### Steps

1. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/messages-service.git
   cd messages-service
   ```

2. **Set up the database** 

    Ensure PostgreSQL is running and create the `messages` database.

    ``` bash
    psql -U postgres
    CREATE DATABASE messages;
    \q
    ```
    Run the schema script:
    ``` bash
    psql -U postgres -d messages -f schema.sql
    ```

3. **Run the application**
    ``` bash
    go run .
    ```

## API Endpoints

- `POST /message`: Create a new message.
- `GET /messages`: List messages (max 100).
- `GET /message/{id}`: Retrieve a message.
- `PUT /message/{id}`: Update a message.
- `DELETE /message/{id}`: Delete a message.

## Testing
Run unit tests:
``` bash
go test ./...
```
