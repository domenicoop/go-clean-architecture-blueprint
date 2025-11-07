# Rules for the `server` Directory

## Directory Organization

The `server` directory is the entry point of the application. It is responsible for the following:

- Loading configuration
- Initializing dependencies
- Setting up the HTTP router
- Starting and gracefully shutting down the HTTP server

Each file in this directory has a specific responsibility:

- `main.go`: The main entry point of the application. It orchestrates the setup and teardown of the application.
- `app.go`: Defines the `application` struct, which holds all the application's dependencies. This is used for dependency injection.
- `config.go`: Handles loading and parsing of application configuration from environment variables.
- `router.go`: Defines the HTTP routes and wires up the handlers.
- `server.go`: Configures and runs the HTTP server, including graceful shutdown logic.

## Best Practices

### Do's

- **Keep `main.go` minimal**: The `main` function should only be responsible for orchestrating the application's startup and shutdown.
- **Use the `application` struct for dependency injection**: Pass the `application` struct to functions that need access to dependencies like the logger, database connection, or configuration.
- **Separate concerns**: Each file should have a single, well-defined responsibility.
- **Handle errors gracefully**: Don't let the application crash on recoverable errors. Use the structured logger to log errors with context.
- **Use graceful shutdown**: Ensure the server shuts down gracefully, allowing active requests to complete.
- **Configure timeouts**: Always configure timeouts for the HTTP server to prevent resource exhaustion.

### Don'ts

- **Do not add business logic to this directory**: Business logic should reside in the `internal/service` package.
- **Do not add database queries to this directory**: Database queries should be in the `internal/repository` package.
- **Do not add HTTP handlers to this directory**: HTTP handlers should be in the `internal/handler/http` package.
- **Avoid global variables**: Use the `application` struct to manage dependencies.
- **Don't panic**: Use `os.Exit(1)` in the `main` function for fatal errors during startup, but avoid `panic` elsewhere.
