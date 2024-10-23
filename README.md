# File Watcher in Go

## Overview

This project implements a file watcher algorithm in Go. The application monitors a specified directory for changes to files and displays notifications in real-time on a web page. This project was inspired by live server extensions that automatically reload web pages when changes are made to HTML files.

## Features

- Monitors a specified directory for file changes (creation, modification, deletion).
- Displays notifications in real-time on a web page when changes occur.
- Built with Go's goroutines and channels for efficient concurrent processing.
- Simple HTML frontend that updates dynamically using Server-Sent Events (SSE).

## Getting Started

### Prerequisites

- Go 1.16 or later
- A terminal or command prompt

### Installation

1. Clone the repository:

   ```bash

   git clone https://github.com/Avnish5/File-Watcher-Golang-.git

2. Navigate to the project directory::

   ```bash
   cd file-watcher


3. Run the application::

   ```bash
   go run main.go
## How It Works

The application continuously checks the specified directory for any changes to files. It uses:

- **Goroutines**: To run file watching in a separate thread, allowing the main thread to handle incoming HTTP requests.
- **Channels**: To communicate between the file watcher and connected clients, ensuring that all notifications are sent in real-time.

When a file is created, modified, or deleted, the application sends a notification to all connected clients through a dedicated channel, and the web page updates automatically.

## Usage

- Modify files in the specified directory to see real-time updates on the web page.
- Refresh the page to establish a new connection to the server and continue receiving updates.

## Acknowledgements

- Golang for its powerful concurrency model.
- Inspiration from live server extensions for web development.

## Contributing

Feel free to submit issues or pull requests if you'd like to contribute to the project!


