# Rules for the `http` Handler Directory

This document outlines the best practices and conventions for creating HTTP handlers within this project. The `http` directory is responsible for handling incoming HTTP requests, delegating to the business layer for processing, and formatting the HTTP response.

## Best Practices

### Do's

- **Keep Handlers Thin**: A handler's primary role is to be a translator. It translates an incoming HTTP request into a call to the business layer and then translates the result from the business layer into an HTTP response. It should be lean and focused on the task of handling HTTP communication.
- **Dependency Injection**: Handlers should receive their dependencies (primarily services from the `service` layer) via their constructor (e.g., `NewEntityHandler`). They should not create their own dependencies. This promotes loose coupling and testability.
- **Use Request/Response DTOs**: Handlers must define their own request and response structs (Data Transfer Objects, or DTOs). These DTOs are specific to the API and contain the necessary serialization tags (e.g., `json:"..."`). This decouples the API contract from the internal domain representation.
- **Perform Model Conversion**: The handler is responsible for converting incoming request DTOs into domain models before calling the service layer. It is also responsible for converting the domain models returned by the service into response DTOs before sending them to the client.
- **Use Standard `net/http` Types**: Handler methods should use the standard `http.ResponseWriter` and `*http.Request` types to remain compatible with Go's standard library and various routers. This ensures interoperability and allows for the use of standard middleware.
- **Error Handling**:
    - Handle errors returned from the service layer and map them to appropriate HTTP status codes.
    - Use a centralized error handling mechanism to avoid repetitive error handling logic in each handler.
    - Return error responses in a consistent format (e.g., JSON with an `error` field).
- **Use Context for Request-Scoped Values**: Use the `context.Context` from the `*http.Request` to pass request-scoped data, such as request IDs, authentication information, or deadlines.

### Don'ts

- **No Business Logic**: Absolutely no business logic should be present in a handler. All business rules, data manipulation, and orchestration of operations belong in the `service` layer. For example, a handler should not decide if a user has permission to perform an action; it should call a service that makes that decision.
- **No Direct Data Access**: Handlers must not interact directly with the database or any other data persistence layer. All data access must go through the `service` layer, which in turn uses repository interfaces.
- **Do Not Use Domain Models Directly in API**: Handlers should not accept or return the core domain models directly in their public-facing API. They must use DTOs.
- **Avoid Global State**: Handlers should be stateless. All necessary context and dependencies should be passed in via the request context or the handler's struct fields. Avoid using global variables.
- **Don't Forget `context.Context`**: Do not neglect the request's context. It is crucial for handling cancellations, timeouts, and passing request-scoped data.