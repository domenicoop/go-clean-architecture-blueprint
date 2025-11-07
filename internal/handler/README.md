# Rules for the `handler` Directory

## Directory Objective and Scope

In the context of the used architecture, the `handler` directory is a key component of the **Interface Adapters** layer. Its primary responsibility is to act as a bridge between the external world (e.g., a web framework, a command-line interface) and the application's core business logic (the `service` layer).

The main objective of this directory is to:

1. **Receive and Parse Incoming Requests**: It takes raw requests from an external source (like an HTTP request) and parses them into a structured format that the application can understand.
2. **Delegate to the Service Layer**: It calls the appropriate methods in the `service` layer to execute the requested business logic.
3. **Format and Send Responses**: It takes the data returned by the `service` layer and formats it into a response suitable for the external source (like a JSON response for an HTTP request).

## Key Responsibilities

- **Data Transfer Objects (DTOs)**: Handlers use DTOs to define the structure of incoming requests and outgoing responses. This decouples the external API from the internal domain models.
- **Model Conversion**: Handlers are responsible for converting request DTOs into domain models to be passed to the `service` layer, and for converting domain models returned by the `service` layer into response DTOs.
- **Input Validation**: Handlers should perform basic validation on the incoming request data to ensure it is in the correct format before passing it to the `service` layer.
- **Error Handling**: Handlers are responsible for catching errors returned by the `service` layer and mapping them to appropriate responses (e.g., HTTP status codes).

## Guiding Principles

- **No Business Logic**: The `handler` layer should be completely devoid of business logic. All business rules and logic must reside in the `service` layer.
- **Dependency Inversion**: The `handler` layer depends on the `service` layer's interfaces, not on its concrete implementations. This allows the business logic to remain independent of the delivery mechanism.
- **Single Responsibility**: Each handler should have a single responsibility, typically corresponding to a specific set of related endpoints or actions.
