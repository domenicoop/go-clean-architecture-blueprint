# Add or Swap a Repository

This workflow outlines the steps to add a new data persistence implementation (a "repository") for an existing service. Because the service layer depends on an interface, not a concrete implementation, you can easily swap from one repository (e.g., `inmemory`) to another (e.g., `postgres`) with minimal changes.

### 1. Define the Repository and its Storage Models

This is where you define the concrete data structures and dependencies for your new repository. This layer is responsible for all interactions with your chosen data store.

- **Location:** `internal/repository/your_new_repository/` (e.g., `internal/repository/postgres/`)
- **Action:**
  1. Create a new Go file (e.g., `postgres.go`).
  2. Define a struct for your repository implementation (e.g., `PostgresRepository`). This struct should hold any dependencies it needs, such as a database connection pool (`*sql.DB`).
  3. Define **internal storage models** that map directly to your database schema (e.g., a `productModel` struct with `db` tags). These models are private to the repository and should not be exposed to the service layer.

### 2. Implement the Repository Interface

Your new repository must satisfy the contract required by the service layer. This ensures it can be used as a drop-in replacement for any other repository.

- **Location:** `internal/repository/your_new_repository/`
- **Action:**
  1. Implement all the methods defined in the `service.EntityRepository` interface.
  2. Inside each method, perform the necessary data operations (e.g., writing SQL queries).
  3. Handle the translation between your private storage models and the `domain.Entity` models.
    - When saving data, convert the incoming domain model into your storage model.
    - When retrieving data, convert your storage model into a domain model before returning it.
  4. Translate any database-specific errors into the standard `apperror` types (e.g., convert `sql.ErrNoRows` to `apperror.ErrNotFound`).

### 3. Wire the New Repository

This is the final step where you update the application's "composition root" to use your new repository implementation.

- **Location:** `cmd/server/`
- **Action:**
  1. In `main.go` or `app.go`, find where the dependencies are being instantiated.
  2. Instantiate your new repository, providing it with any required dependencies (like a database connection).
  3. Inject the new repository instance into the service constructor. The service layer will accept it without any changes, as it satisfies the required interface.
