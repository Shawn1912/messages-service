# Messages Service

## Description

A simple Go-based microservice that provides CRUD operations for messages and includes palindrome detection functionality.

## Architecture
![Test Image 4](/architecture.png)

## Project Structure

messages-service/ <br />
├── database <br /> &emsp;&emsp;
    ├── db_connection_test.go <br />&emsp;&emsp;
    ├── db_connection.go <br />&emsp;&emsp;
    ├── models.go <br />&emsp;&emsp;
    └── schema.sql  <br />
├── handlers <br /> &emsp;&emsp;
    ├── handlers.go <br />&emsp;&emsp;
    └── handlers_test.go  <br />
├── utils <br /> &emsp;&emsp;
    ├── palindrome.go <br />&emsp;&emsp;
    └── palindrome_test.go  <br />
├── go.mod  <br />
├── go.sum  <br />
├── main.go  <br />
└── main_test.go

## Build and Run

### Prerequisites

- Go (version 1.17 or higher)
- Docker
- PostgreSQL

### Steps

1. **Clone the repository**

   ```bash
   git clone https://github.com/shawn1912/messages-service.git
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
    psql -U postgres -d messages -f ./database/schema.sql
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

### Example: Creating a message
``` bash
curl -X POST http://localhost:8080/messages \
  -H 'Content-Type: application/json' \
  -d '{"content": "A man a plan a canal Panama"}'
```

### Response
``` json
{
  "id": 29,
  "content": "A man a plan a canal Panama",
  "isPalindrome": true,
  "createdAt": "2024-10-28T12:00:00Z",
  "updatedAt": "2024-10-28T12:00:00Z"
}
```

## Testing
Run unit tests:
``` bash
go test ./...
```
