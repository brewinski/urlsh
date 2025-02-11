
# URL Shortener

A simple and efficient URL shortening service built with Go, HTMX, and SQLite. This project demonstrates how to create a modern web application with minimal JavaScript using HTMX for dynamic interactions.

## Demo

https://github.com/user-attachments/assets/f6c64d0a-a73e-43eb-8fe2-795a2bc93abc

## Features

- Shorten long URLs into compact, shareable links
- Automatic redirection from short URLs to original destinations
- Clean, responsive user interface
- Real-time feedback using HTMX
- Persistent storage using SQLite
- No JavaScript required (besides HTMX)

## Technology Stack

- Backend: Go
- Frontend: HTML + HTMX
- Database: SQLite
- CSS: Custom styling

## Getting Started

### Prerequisites

- Go 1.22 or later
- SQLite


### Installation

1. Clone the repository:
```bash
git clone https://github.com/brewinski/urlsh.git
cd url-shortener
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
make start
```

The application will be available at http://localhost:8801/

## Usage

1. Open your browser and navigate to http://localhost:8801/
2. Enter a long URL in the input field
3. Click "Shorten URL"
4. Copy and share the generated short URL

## Development

The project uses a Makefile for common tasks:

- `make run`: Start the Go server
- `make start`: Open the browser and start the server
- `make open`: Open the application in your default browser

## How It Works

1. When a user submits a URL, it's hashed using SHA-256
2. The hash is base64 encoded to create a unique identifier
3. The original URL and its identifier are stored in SQLite
4. When accessing a short URL, the system looks up the original URL and redirects

