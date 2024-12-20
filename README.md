### Cache Proxy Server

This is a simple and efficient cache proxy server built using Golang. The server acts as an intermediary for requests from clients seeking resources from other servers. It caches responses to improve performance and reduce the load on the backend servers.

#### Features:
- **Caching**: Stores responses to reduce load on backend servers and improve response times.
- **Configurable**: Easily configurable via command-line arguments for URL and port.
- **Graceful Shutdown**: Handles server shutdown gracefully to ensure all requests are processed.
- **Modular Design**: Separation of concerns for better maintainability and testability.

#### Usage:
1. **Install the package**:
    ```sh
    go install github.com/balajiss36/cache-proxy
    ```
2. **Start the server**:
    ```sh
    cache-proxy start --url http://example.com --port 8080
    ```

#### Dependencies:
```go.mod
module github.com/balajiss36/cache-proxy

go 1.22.9

require (
    github.com/inconshreveable/mousetrap v1.1.0 // indirect
    github.com/spf13/cobra v1.8.1 // indirect
    github.com/spf13/pflag v1.0.5 // indirect
)