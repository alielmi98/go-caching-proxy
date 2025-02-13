# go-caching-proxy

This project implements a caching proxy server in Go that forwards requests to a specified origin URL and caches the responses to improve performance.

This project is part of the backend section from [roadmap.sh](https://roadmap.sh/projects/caching-server).

## Features

- Forwards HTTP requests to an origin server.
- Caches responses to reduce load on the origin server and improve response times.
- Supports cache clearing via command-line interface.

## Project Structure

```
go-caching-proxy
├── cmd
│   └── main.go        # Entry point of the application
├── pkg
│   ├── cache
│   │   └── cache.go   # Caching functionality
│   ├── proxy
│   │   └── proxy.go   # Proxy server logic
│   └── utils
│       └── utils.go   # Utility functions
├── go.mod              # Module definition
└── README.md           # Project documentation
```

## Installation

To install the project, clone the repository and navigate to the project directory:

```bash
git clone https://github.com/alielmi98/go-caching-proxy.git
cd go-caching-proxy
```

Then, run the following command to download the dependencies:

```bash
go mod tidy
```

## Usage

To run the caching proxy server, use the following command:

```bash
go run cmd/main.go --port <port> --origin <origin-url>
```

Replace `<port>` with the desired port number and `<origin-url>` with the URL of the origin server.

To clear the cache before starting the server, use the `--clear-cache` flag:

```bash
go run cmd/main.go --port <port> --origin <origin-url> --clear-cache
```

## Clearing the Cache

To clear the cache while the server is running, you can send a request to the cache clearing endpoint:

```bash
curl http://localhost:<port>/clear-cache
```

Replace `<port>` with the port number on which the server is running.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.