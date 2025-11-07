# Getting Started

This guide provides a comprehensive walkthrough to get you started with this Go application template. This will cover the core principles, the project structure, and a step-by-step workflow to implement a feature.

To better understand the project, you can read more about the overall architecture [here](ARCHITECTURE.md).

## A Step-by-Step Workflow

The best way to understand the architecture is to build something.

### 1. Define the Core Model (Domain Layer)

This is the foundation of the feature. It represents the core business concept.

- **Location:** `internal/domain/`
- **Action:** Define the primary domain model struct (e.g., `Product`). This struct should only contain the essential fields and business logic methods (e.g., validation) for that entity.

### 2. Define the Contracts (Service Layer)

This layer defines the application's capabilities and the requirements for outer layers.

- **Location:** `internal/service/interfaces.go`
- **Action:**
    1. **Define the Repository Interface:** Specify the contract for data persistence (e.g., `ProductRepository`) with methods like `Create` or `FindByID`. These methods must accept and return domain models from the previous step.
    2. **Define the Service Interface:** Specify the contract for the business logic (e.g., `ProductService`) that defines the use cases for the feature, such as `CreateNewProduct`.

### 3. Implement the Data Access (Repository Layer)

This layer provides the concrete implementation for data storage, fulfilling the repository interface contract.

- **Location:** `internal/repository/inmemory/` (or a new one like `internal/repository/postgres/`)
- **Action:**
    - Create a struct for the repository implementation (e.g., `ProductRepository`).
    - Implement the methods defined in the repository interface.
    - This is where you will interact with the database. You are responsible for **converting incoming domain models into storage-specific models** before saving, and **converting storage models back into domain models** after retrieval.

### 4. Implement the Business Logic (Service Layer)

This layer contains the core application logic and orchestrates the flow of data.

- **Location:** `internal/service/`
- **Action:**
    - Create a struct for the service implementation. It must hold a reference to the **repository interface** you defined earlier.
    - Implement the methods from the service interface.
    - Orchestrate the business logic by calling methods on the domain models and the repository.

### 5. Create the Handler

This is the outermost layer that exposes the feature to the world.

- **Location:** `internal/handler/http/`
- **Action:**
    1. Define request and response structs (DTOs) specific to this endpoint.
    2. Create a handler struct that holds a reference to the **service interface**.
    3. Implement the handler functions for the API endpoints. Each function is responsible for:
        - Decoding and validating the incoming request DTO.
        - **Converting the request DTO into a domain model.**
        - Calling the appropriate method on the service interface.
        - **Converting the resulting domain model (or error) into a response DTO.**
        - Sending the response to the client.

### 6. Wire All Components Together

This is the composition root where the application starts and the dependency graph is built.

- **Location:** `cmd/server/`
- **Action:**
    1. In `app.go` or `main.go`, instantiate the concrete implementations: repository, then service, then handler.
    2. In `router.go`, register the new API route and map it to the corresponding method in the handler.

---

## Have Questions?

For more details on the design choices and architectural reasoning, please check out the [Frequently Asked Questions (FAQ)](./FAQ.md).
