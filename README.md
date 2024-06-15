# Movie Recommendation System - Go Part

This part of the project provides a REST API to fetch movie recommendations for users.

## Table of Contents

- [Overview](#overview)
- [Structure](#structure)
- [Setup](#setup)
- [Usage](#usage)
- [License](#license)

## Overview

The Go part of the project sets up a REST API that reads from a SQLite database containing movie recommendations and provides endpoints to fetch recommendations for users.

## Structure

    go/
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── handlers.go
    └── handlers_test.go


## Setup

### Prerequisites

- Go 1.18 or higher
- Docker

### Building and Running the Docker Container

1. Build the Docker image:
```bash
docker build -t recommendations-api .
```

2. Run the Docker container:
```bash
docker run -p 8080:8080 recommendations-api
```

## Usage

### Fetch Recommendations

You can fetch movie recommendations for a user by sending a GET request to the API.

Example:
```bash
curl http://localhost:8080/recommendations/1/5
```

### Running Tests
```bash
go test ./...
```

## License
This project is licensed under the MIT License.