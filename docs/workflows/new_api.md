# Implement a new API Workflow

### 1. Define the Core Model

This is the foundation of your feature. It represents the core business concept and has no dependencies on any other part of your application.

- **Location:** `internal/domain/`
- **Action:**
 - Define the primary **domain model struct** for your feature (e.g., `Product`).
 - Include the essential fields that represent the entity.
 - Add methods to this struct that enforce **invariant business rules** (e.g., a validation method that ensures a price cannot be negative).

### 2. Define the Contracts (Service Layer Interfaces)

This layer defines the application's capabilities and the requirements for outer layers. It depends only on the **domain** layer.

- **Location:** `internal/service/` (or `internal/business/` as per your original structure)
- **Action:**
 1. **Define the Repository Interface:** Specify the contract for data persistence. Create an interface (e.g., `ProductRepository`) with methods that the data access layer must implement, such as `Create`, `FindByID`, or `Update`. These methods should accept and return **domain models**.
 2. **Define the Service Interface:** Specify the contract for the business logic. Create an interface (e.g., `ProductService`) that defines the use cases for your feature, such as `CreateNewProduct` or `GetProductDetails`.

### 3. Implement the Data Access

This layer provides the concrete implementation for data storage and retrieval, fulfilling the contract defined by the repository interface.

- **Location:** `internal/repository/postgres/`
- **Action:**
    - Create a struct that will act as the repository implementation and hold a database connection.
    - **Define internal storage models** that represent how data is structured in your database. These can include `db` tags.
    - Implement all the methods defined in the **repository interface** from the service layer.
    - This is where you write the actual database queries (e.g., SQL). Your methods will interact directly with the database, **converting incoming domain models into storage models before saving, and converting storage models into domain models after retrieval.**

### 4. Implement the Business Logic

This layer contains the core application logic and orchestrates the flow of data. It fulfills the service interface contract.

- **Location:** `internal/service/`
- **Action:**
    - Create a struct for the service implementation. It must hold a reference to the **repository interface** (not the concrete implementation).
    - Implement all the methods defined in the **service interface**.
    - Within these methods, orchestrate the business logic by calling methods on the domain models (for validation) and the repository interface (for data persistence).

### 5. Create the API Endpoint

This is the outermost layer that exposes your feature to the outside world (e.g., via a REST API). It translates external requests into application commands.

- **Location:** `internal/handler/http/`
- **Action:**
 1. Define request and response structs, also known as Data Transfer Objects (DTOs). These are specific to the API and are used for serializing and deserializing data like JSON.
 2. Create a handler struct that holds a reference to the **service interface**.
 3. Implement the handler functions for your API endpoints. Each function is responsible for:
    - Binding and validating the incoming request DTO.
    - **Converting the request DTO into a domain model.**
    - Calling the appropriate method on the **service interface** with the domain model.
    - **Converting the resulting domain model (or error) into an appropriate response DTO.**
    - Sending the response DTO and HTTP status code to the client.

### 6. Wire All Components Together

This is the composition root where the application starts. It's responsible for building the dependency graph and starting the server.

- **Location:** `cmd/server/`
- **Action:**
 1. In your main application setup (`main.go` or `app.go`), instantiate your concrete implementations in the correct order, from outer to inner dependencies:
    - Establish the database connection.
    - Create an instance of your concrete **repository**, passing it the database connection.
    - Create an instance of your concrete **service**, passing it the repository instance.
    - Create an instance of your **handler**, passing it the service instance.
 2. In your router setup (`router.go`), register the new API route and map it to the corresponding method in your handler instance.
