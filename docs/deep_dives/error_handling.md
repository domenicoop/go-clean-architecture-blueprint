# Error Handling

In this architecture, we treat errors as first-class values that flow through the system. They are not just exceptional cases but a predictable part of the application's logic.

### Known Errors

These are predictable errors that are part of the business domain or application contract. They are errors you *expect* to happen.

**Examples:**
    - `apperror.ErrNotFound`: A requested resource was not found.
    - `apperror.ErrInvalidInput`: User-provided data failed validation.
    - `apperror.ErrConflict`: A resource creation failed due to a conflict (e.g., duplicate email).

These errors are translated into specific, client-friendly responses (e.g., HTTP `404`, `400`, `409`).

### Unknown Errors 

These are unexpected, internal failures.

**Examples:**
    - A database connection is lost.
    - A nil pointer is dereferenced.
    - A disk is full.

These errors must **never** be sent to the client. They are logged in detail internally, and a generic, safe response (e.g., HTTP `500 Internal Server Error`) is sent to the client.

---

## The Error Flow Layer-by-Layer

Each layer has a distinct responsibility in the error handling flow. The error flows "up" from the innermost layer, getting wrapped with context along the way.

### 1. The `repository` Layer

- It **translates** infrastructure-specific errors into standard `apperror` types.
- It **must not** leak database-specific errors (like `sql.ErrNoRows`) to the service layer. This would violate the rule that the service layer is independent of the data layer.
- It catches errors like `sql.ErrNoRows` and returns `apperror.ErrNotFound`.
- It catches a "unique constraint violation" and returns `apperror.ErrConflict`.
- If it receives an error it *cannot* translate (e.g., connection timed out), it returns the original error, which will be treated as an "Unknown Error" by the central handler.

### 2. The `service` Layer

- It performs business-level validation (e.g., "is `entity.Name` empty?") and is the primary source for `apperror.ErrInvalidInput`.
- When it receives an error from the repository (e.g., `apperror.ErrNotFound`), it **wraps** it with business-level context using `fmt.Errorf`.

### 3. The `handler` Layer

- This layer is the **central point for translating application errors into transport-specific responses** (e.g., HTTP status codes).
- It receives errors from the `service` layer (e.g., `apperror.ErrNotFound`).
- It also generates its own errors for *transport-level* issues (e.g., malformed JSON), which it should wrap in `apperror.ErrInvalidInput`.
- It uses a **central error helper** (e.g., `handleError`) defined *within* the `handler` package (e.g., in `internal/handler/http/errors.go`). This helper is injected with a logger during handler creation.
- This `handleError` function is the only place in the handler that contains the `switch` statement to inspect errors.
    - It checks for **Known Errors** (using `errors.Is`) and writes the appropriate `4xx` response.
    - It logs **Unknown Errors** (the `default` case) and writes a generic `500` response.

### 4. The `cmd` (Main) Layer

- This is the application's entry point and "composition root".
- Its primary role in error handling is to **inject the application logger** (`*slog.Logger`) into the `handler` layer when it's created.
- It does **not** contain any error-to-HTTP-response mapping logic. It simply wires up the components and starts the server.

---

### Example Flow: "Entity Not Found"

1. **Handler:** The user requests `GET /entities/123`. The handler calls `h.service.GetByID(ctx, "123")`.
2. **Service:** The service validates the ID and calls `s.repo.FindByID(ctx, "123")`.
3. **Repository:** The repository (e.g., a SQL implementation) executes `SELECT ... WHERE id=123`. The database returns `sql.ErrNoRows`. The repository **translates** this and returns `apperror.ErrNotFound`.
4. **Service:** The service receives `apperror.ErrNotFound`. It wraps it: `return nil, fmt.Errorf("service: failed to find entity with id 123: %w", err)`.
5. **Handler:** The handler receives the wrapped error. It calls its *own* internal helper: `h.handleError(w, r, wrappedErr)`.
6. **Handler (handleError):** The helper receives the error.
    - It checks: `errors.Is(wrappedErr, apperror.ErrNotFound)`.
    - This check returns `true`.
    - The helper writes a `404 Not Found` response to the client.