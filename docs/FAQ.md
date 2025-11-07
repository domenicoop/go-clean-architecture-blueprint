# Frequently Asked Questions (FAQ)

### Shouldn't the service layer be responsible for converting storage models to domain models?

While it might seem intuitive to place conversion logic in the service layer, in this design, the **repository layer is the correct place for this responsibility.**

The primary role of the repository is to abstract away all details of data persistence. This includes not only *how* data is fetched (e.g., SQL queries, API calls) but also in *what format* it is stored. The conversion from a storage-specific model to a pure domain model is the final step of that abstraction. Also the architecture dictates that dependencies must flow inwards. The service layer (an inner layer) should not know anything about the implementation details of the repository (an outer layer). If the service had to perform the conversion, it would need to depend on the repository's internal storage models, which violates this core principle. In the end by keeping this logic in the repository, the service layer remains completely decoupled from the persistence framework. You can swap the database from PostgreSQL to MongoDB, and the only layer that needs to change is the repository. The services and their tests, which operate purely on domain models, remain untouched.

### Why does the `service` layer seem to have two different rules for interfaces? It defines the `EntityRepository` interface but provides the `EntityService` interface.

The `service` layer acts as the "core" of the application and interacts with "outer" layers (like `handler` and `repository`) in two distinct ways, each following a specific pattern:

- **As a CONSUMER (using the `repository` layer):**
    - The `service` layer consumes data from the `repository` layer. To protect itself from infrastructure details, it uses the **Dependency Inversion Principle**.
    - It defines its *own* contract, the `EntityRepository` interface, specifying what it *needs* from a data layer.
    - This makes the `service` layer completely independent of any specific database. The `repository` implementation (e.g., `inmemory` or `postgres`) must then implement this interface.

- **As a PROVIDER (used by the `handler` layer):**
    - The `service` layer provides business logic to the `handler` layer. Here, it acts as a "public API" for the rest of the application.
    - It defines the `EntityService` interface as the stable contract that it *provides*.
    - This allows any number of handlers (e.g., `http`, `gRPC`, `CLI`) to use the same business logic, all while depending on a stable interface, not a concrete implementation.

### I see the handlers and repositories are on opposite sides of the service. Are they both "adapters"?

Yes, they are. In a Clean/Hexagonal Architecture, both the `handler` and `repository` layers are considered **Adapters**. They are responsible for translating data between the core application (the `service` and `domain` layers) and the "outside world."

They just adapt for different things:

- **Primary Adapters (Driving Adapters):** This is the **`handler`** layer. It "drives" the application. It takes requests from the outside world (like an HTTP request) and translates them into calls to the `service` layer.
- **Secondary Adapters (Driven Adapters):** This is the **`repository`** layer. It is "driven by" the application. It takes commands from the `service` layer (like "save this entity") and translates them into operations for a specific technology (like an in-memory map or a SQL database).

### If both handlers and repositories are adapters, why do we connect them to the service layer in two different ways?

This is to protect the `service` layer and enforce the **Dependency Inversion Principle**. The two connection patterns are different because the adapters have different relationships with the core logic:

- **For Handlers (Driving Adapters):** The `service` layer acts as a "public API" for the application. It **provides** a stable `EntityService` interface. The `handler` (the "user" in this case) imports the `service` package and depends on that interface. This allows multiple types of handlers (HTTP, gRPC, etc.) to all use the same core business logic.

- **For Repositories (Driven Adapters):** The `service` layer (the core) must **not** depend on any infrastructure (like a database). To achieve this, the `service` layer (the "user") **defines** the `EntityRepository` interface it *needs*. The `repository` adapter *must* then import the `service` package to implement that interface.

This two-pattern approach ensures that all dependencies flow "inward" toward the `service` layer, and the `service` layer itself remains pure, testable, and completely independent of the outside world.

### There should be a single "generic" repository for all my entities, or multiple specific repositories?

The preferred approach is **multiple repositories**, with one repository per entity (or "aggregate root").
This means you would have an `EntityRepository`, a `UserRepository`, a `ProductRepository`, and so on.

There are multiple reasons for this:

- **Interface Segregation Principle:** A `UserService` should only depend on a `UserRepository` interface that *only* has user-related methods (e.g., `GetUser`, `CreateUser`).
- **The Anti-Pattern:** If you have one giant `DatabaseRepository` with `GetUser`, `GetProduct`, and `GetOrder`, your `UserService` is now forced to "see" methods it doesn't care about, which creates unnecessary coupling.
- **Separation of Concerns (SoC):** Each repository should have a single, clear responsibility: managing the persistence of a single domain entity. A generic repository that handles all entities has low cohesion and mixes too many concerns.
- **Go Idiom:** Go strongly favors small, well-defined interfaces. Defining a `UserRepository` interface in your `service` layer for a `postgres` package to implement is the most idiomatic Go way to apply dependency inversion.

### In this architecture, where should clients for external services go?

Clients for external services (like a payment gateway or an email API) are treated as **Driven Adapters**, exactly like repositories.

The implementation must be split into two parts to protect the `service` layer from infrastructure details and uphold the Dependency Inversion Principle:

1. **The Interface (in the `service` layer):**
    Your `service` layer defines an interface for the *business functionality* it needs. This interface should be placed in `internal/service/interfaces.go`.
   - **Example:** `type NotificationSender interface { Send(ctx context.Context, to string, message string) error }`

2. **The Implementation (in a new `client` layer):**
    You create a new package outside the `service` layer (e.g. `internal/client/sendgrid`). This package contains the concrete struct, holds the actual HTTP client or SDK, and implements the `NotificationSender` interface defined in the `service` layer.

This way, your `service` layer just depends on the `NotificationSender` interface and has no idea what specific external service is being used, keeping it pure and testable.

### Why use "Package-by-Layer" instead of "Package-by-Feature"?

This project is organized by layer (`domain`, `service`, `handler`) rather than by feature (`products`, `users`). This is a deliberate choice with specific trade-offs.

The primary goal of this template is to *teach* Clean Architecture. A layer-based structure makes the boundaries and the **Dependency Rule** visually explicit. You can see at a glance that `service` does not import `handler`. While a "package-by-feature" approach is also used, this layer-based approach is a pure and classic implementation of Clean Architecture, making it a perfect starting point. Both patterns are valid, but this one aligns best with the educational goals of the project.