# Rules for the `repository` Directory

## Directory Organization

The `repository` directory contains concrete implementations of the data access interfaces defined in the `service` layer. It acts as a bridge between the application's business logic and the underlying data store.

- `/inmemory` (or other data store specific directories like `/mysql`, `/mongodb`, `/postgres`): Each subdirectory represents a specific data store implementation. This allows the application to easily switch between different database technologies.
- `{model_name}.go`: Inside a specific implementation directory (e.g., `/inmemory`), these files contain the concrete repository structs and methods. For example, `entity.go` provides the InMemory-specific implementation for storing and retrieving `Entity` domain models.

## Best Practices

### Do's

- **Implement Business Layer Interfaces**: Repositories **must** implement the corresponding repository interfaces defined in the `service` layer (e.g., `business.EntityRepository`). This is crucial for maintaining the separation of concerns and the dependency inversion principle.
- **Encapsulate Data Access Logic**: All logic related to a specific data store should be contained within its repository. This includes:
 - Writing SQL queries (or using a query builder/ORM).
 - Managing database transactions.
 - Handling database-specific error codes and translating them into generic application errors if necessary.
- **Perform Data Mapping**: The repository is responsible for mapping data between the database's representation (e.g., table rows) and the application's `domain` models. It should define its own internal storage-specific models and perform the conversion to and from the domain models.
- **Use Dependency Injection**: The database connection handle (e.g., `*sql.DB`) should be passed into the repository's constructor (e.g., `NewEntityRepository`). The repository should not create or manage the database connection itself.
- **Use** `context.Context`: Pass `context.Context` as the first argument to all repository methods to support cancellation, deadlines, and tracing.

### Don'ts

- **No Business Logic**: Repositories should be completely devoid of business logic. Their only concern is data persistence. For example, a repository should save a user to the database, but it should not know *why* the user is being saved or perform any validation on the user's data (that's the service's job).
- **No Interaction with Handlers**: Repositories must never be called directly from the `http` handler layer. The flow is always `handler` -> `service` -> `repository`.
- **Don't Leak Implementation Details**: The methods of a repository should accept and return the pure `domain` models. Do not expose database-specific models or structs to the business layer. The business layer should be completely unaware of the underlying database schema or technology.
- **Avoid Generic Repositories**: While generic repository patterns can seem appealing, they often lead to leaky abstractions. It's better to create specific repository interfaces and implementations for each aggregate root in your domain.
