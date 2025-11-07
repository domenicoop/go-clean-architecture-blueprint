# Rules for the `domain` Directory

## Directory Organization

The `domain` directory contains the core domain models of the application. These are plain Go structs that represent the fundamental business objects and concepts.

- `{model_name}.go` files: These files define the structs for the core entities. All domain definitions should reside here.

## Best Practices

### Do's

- **Define Plain Go Structs**: Entities must be simple, plain Go structs. They represent the data and the intrinsic behavior of a business object.
- **Keep them Pure**: These models should be pure and self-contained. They represent *what* the application is about, not *how* it works.
- **Add Intrinsic Methods**: You can add methods to domain structs, but only if they represent behavior that is inherent to the domain itself (e.g., a validation method like `user.IsValidEmail()` or a calculation like `order.CalculateTotal()`).

### Don'ts

- **No Serialization Tags**: Domain models must **not** contain struct tags for serialization or database mapping (e.g., `json:"..."`, `db:"..."`). This keeps the core domain completely independent of infrastructure details.
- **No External Dependencies**: Entities must **not** have dependencies on any other part of the application or external frameworks. This means no dependencies on:
 - The database (`*sql.DB`)
 - Loggers (`*slog.Logger`)
 - Web frameworks (Chi, Gin, etc.)
 - Repositories or Services
- **No Infrastructure Logic**: Do not include any code related to infrastructure concerns. This includes:
 - Database connection logic
 - SQL queries
 - API handling logic (e.g., marshalling/unmarshalling JSON)
- **No Complex Business Logic**: Do not place complex business logic that orchestrates multiple entities or involves external services here. That logic belongs in the `service`(service) layer. For example, a function that creates a new user and sends a welcome email does not belong here.
