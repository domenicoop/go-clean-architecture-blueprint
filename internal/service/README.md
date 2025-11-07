# Rules for the `service` Directory

## Directory Organization

The `service` directory is responsible for implementing the application's business rules and use cases. It acts as an intermediary between the presentation layer (handlers) and the data access layer (repositories).

- `*{service_name}.go`: These files contain the service implementations. Each service encapsulates a specific area of business logic (e.g., `EntityService` for entity-related operations). Services are the primary components of this layer.
- `interfaces.go`: This file defines the interfaces that the business services depend on. These interfaces act as contracts for external components, such as repositories. This adheres to the Dependency Inversion Principle, allowing the business logic to be independent of the data layer's implementation.

## Best Practices

### Do's

- **Encapsulate Business Logic**: All core business logic, validation rules, and use case orchestration should be placed within services in this directory.
- **Depend on Interfaces**: Services must depend on interfaces (defined in `interfaces.go`), not on concrete implementations of repositories or other external components. This makes the business logic decoupled and easier to test.
- **Be Stateless**: Services should be stateless. They receive all the data they need as arguments to their methods and should not hold state between calls.
- **Orchestrate Data Flow**: A service's role is to orchestrate the flow of data. It calls repository methods to fetch or persist data and executes business rules on that data.
- **Return Plain Entities**: Services should work with and return the core `entity` objects defined in the `/internal/entity` package.
- **Use Context**: Pass `context.Context` as the first argument to all service methods for cancellation, deadlines, and passing request-scoped values.
- **Validation**: Any parameter passed to a service method should be validated before proceeding with the business logic.

### Don'ts

- **No Transport Layer Code**: Do not include any code related to the transport layer (e.g., no `http.Request`, `http.ResponseWriter`, or framework-specific context objects). The business layer should be completely independent of how the application is exposed (e.g., HTTP, gRPC, CLI).
- **No Data Persistence Code**: Do not write any database-specific code (e.g., no `*sql.DB`, SQL queries, or ORM-specific code). All data access should be done through the repository interfaces.
- **No Direct Configuration Access**: Services should not directly access application configuration. Any required values should be passed in during the service's initialization.
- **No Circular Dependencies**: A service in the `service` layer should not depend on the `handler` layer.
- **Missing validation**: Ensure that all input parameters are validated before proceeding with the business logic.
