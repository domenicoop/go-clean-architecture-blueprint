# Architecture Overview

This document provides a detailed overview of the architectural foundation of this project, which is based on the principles of **Clean Architecture**, as popularized by Robert C. Martin.

The primary goal of this architecture is to create a system that is:
- **Independent of Frameworks:** The core business logic has no knowledge of external frameworks like web servers or databases.
- **Testable:** The business rules can be tested without the UI, database, or any other external element.
- **Independent of UI:** The UI can change easily without changing the rest of the system.
- **Independent of Database:** You can swap out your persistence layer without affecting the business rules.

---

## The Layers

The architecture is organized into a series of layers, each with a distinct responsibility.

![Architecture diagram](architecture.svg)

### 1. Domain Layer (`internal/domain`)

- **Responsibility:** This is the innermost layer. It contains the core business entities and the rules that are fundamental to the application (enterprise-wide business rules).
- **Content:** Plain Go structs representing the core concepts (e.g., `Product`, `User`). These structs may have methods that enforce invariants (e.g., `product.IsValid()`).
- **Dependencies:** It has **zero** dependencies on any other layer in the application.

### 2. Service Layer (`internal/service`)

- **Responsibility:** This layer contains the application-specific business logic. It orchestrates the domain entities to perform use cases. For example, a `CreateNewProduct` use case would live here.
- **Content:**
    - **Service Implementations:** Structs that contain the logic for the use cases.
    - **Interfaces:** This layer defines the interfaces (contracts) that outer layers **must** implement. This is a key part of the Dependency Inversion Principle. For example, it defines `ProductRepository`, specifying what data access operations the application needs, without knowing *how* they will be implemented.
- **Dependencies:** It depends only on the **Domain** layer.

### 3. Adapter Layers

These layers are responsible for converting data between the format most convenient for the Service layer and the format most convenient for external agencies like the database or the web.

#### Handlers (`internal/handler/http`, `internal/handler/grpc`)

- **Type:** Primary or "Driving" Adapters. They drive the application.
- **Responsibility:** To adapt incoming requests from the outside world (e.g., an HTTP request) into calls to the Service layer.
- **Content:** Web handlers, gRPC servers, or CLI commands. They are responsible for:
    - Parsing incoming requests.
    - Using Data Transfer Objects (DTOs) for request and response bodies.
    - Converting DTOs to and from Domain models.
    - Calling the appropriate Service layer method.
    - Handling transport-specific error responses.
- **Dependencies:** They depend on the **Service** layer.

#### Repositories (`internal/repository/postgres`, `internal/repository/inmemory`)

- **Type:** Secondary or "Driven" Adapters. They are driven by the application.
- **Responsibility:** To implement the repository interfaces defined in the Service layer. They handle all the details of data persistence.
- **Content:** Concrete implementations for interacting with a database. They are responsible for:
    - Writing SQL queries or using an ORM.
    - Defining internal storage models (e.g., structs with `db` tags).
    - Converting these storage models to and from Domain models.
    - Translating database-specific errors into application-defined errors.
- **Dependencies:** They depend on the **Service** layer (to implement its interfaces).

---

## The Dependency Rule in Practice

The strict, inward-facing direction of dependencies is the cornerstone of this architecture.

- The **Service** layer defines an interface it needs, like `ProductRepository`.
- The **Repository** layer (e.g., `postgres.ProductRepository`) implements that interface.

This means the `postgres` package **imports** the `service` package. The dependency points **inward**. The `service` layer has no knowledge of PostgreSQL, allowing us to swap it for a different database without changing a single line of code in the `service` or `domain` layers.

## Data and Control Flow

Here is the typical flow for an incoming HTTP request:

1. **Request In:** An HTTP request hits a route, which is handled by a method in the `http.Handler`.
2. **Handler:**
    - The handler decodes the JSON request body into a request DTO.
    - It validates the request DTO.
    - It converts the DTO into a `domain.Entity` model.
    - It calls the `service.EntityService` with the domain model.
3. **Service:**
    - The service executes its business logic.
    - It calls the `repository.EntityRepository` interface to persist or retrieve data.
4. **Repository:**
    - The concrete repository implementation receives the domain model.
    - It converts the `domain.Entity` into an internal storage model.
    - It executes a SQL query.
    - It converts the result from the storage model back to a `domain.Entity` and returns it to the service.
5. **Response Out:** The flow reverses. The service returns the domain model to the handler. The handler converts it into a response DTO, marshals it to JSON, and writes the HTTP response to the client.
