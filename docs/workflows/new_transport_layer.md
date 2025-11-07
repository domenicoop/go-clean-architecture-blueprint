# Add a New Transport Layer

### 1. Define the API Contract

This is the "schema" for your new API. For gRPC, this is your `.proto` file. This is the source of truth for your gRPC service's methods and message types.

- **Location:** `docs/proto/v1/`
- **Action:**
    1. Create your `.proto` file (e.g., `entity.proto`).
    2. Define your `service` (e.g., `EntityService`) with its `rpc` methods (e.g., `GetByID`, `Create`).
    3. Define the request and response `message` types (e.g., `GetByIDRequest`, `EntityResponse`).
    4. Generate the Go code using `protoc`. This will create the necessary Go interfaces and structs (e.g., `entity_grpc.pb.go` and `entity.pb.go`).

### 2. Implement the New Handler

This new handler is the gRPC equivalent of your existing `internal/handler/http` package. It acts as the translator between the gRPC transport and the core business logic.

- **Location:** `internal/handler/grpc/`
- **Action:**
    1. Create a new handler struct (e.g., `EntityGRPCHandler`).
    2. This struct **must** hold a dependency on the **service interface** (e.g., `service.EntityService`). This is the *exact same interface* your HTTP handler uses.
    3. Implement the gRPC-generated server interface (e.g., `pb.UnimplementedEntityServiceServer`).
    4. Implement the gRPC methods (e.g., `GetByID`, `Create`). Each method is responsible for:
        * **Converting the incoming gRPC request message into a domain model.**
        * Calling the appropriate method on the **service interface** with the domain model.
        * **Converting the resulting domain model (or error) back into the appropriate gRPC response message.**

### 3. Wire the New Server

This is where you update your application's entry point to instantiate and run the new gRPC server.

  - **Location:** `cmd/server/`
  - **Action:**
    1. **Configuration:** In `config.go`, add new configuration variables for the gRPC server (e.g., `grpc_port`).
    2. **Dependency Injection:** In `main.go` or `router.go`, instantiate your new gRPC handler. You will inject the *exact same service instance* that your HTTP handler uses.**
    3. **Server Registration:** In `main.go` or a new `grpc_server.go`, create the gRPC server and register your handler with it.