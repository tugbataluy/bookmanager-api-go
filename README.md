# bookmanager-api-go
Rest API implementation for book manager database facilitates UX with CLI support
## Features

- RESTful API for managing books
- CLI support for common operations
- Modular project structure (handlers, db, api, cli vb.)
- Easy integration with PostgreSQL

## Getting Started

1. **Clone the repository:**
    ```bash
    git clone https://github.com/yourusername/bookmanager-api-go.git
    cd bookmanager-api-go
    ```

2. **Install dependencies:**
    ```bash
    go mod download
    ```

3. **Run the API server:**
    ```bash
    go run bookmanager/api/main.go
    ```

4. **Run the CLI:**
    ```bash
    go run bookmanager/cmd/bookmanager/main.go
    ```
5. **Optional Build the API server:**
    ```bash
    go build -o bookmanager_api bookmanager/api/main.go
    ```

6. **Optional Build the CLI:**
    ```bash
    go build -o bookmanager bookmanager/cmd/bookmanager/main.go
    ```
## Usage 

- Access API endpoints at `http://localhost:8080/api/v1`
- Use the CLI for quick operations (see [cmd/bookmanager/CLI.md](/bookmanager/cmd/bookmanager/CLI.md))


## Project Structure

- `api/` - API handlers and routes ([api/Rest-API.md](/bookmanager/api/Rest-API.md))
- `cli/` - Command-line interface ([bookmanager/CLI.md](/bookmanager/cmd/bookmanager/CLI.md))
