## Plan: Java to Go Migration (Nemo Project)

This plan outlines the strategy to migrate the existing Java Spring Boot modular monolith to a Go microservice/monolith, ensuring strict API and data compatibility.

### Core Strategy
*   **Architecture**: Maintain the **Modular Monolith** structure initially to simplify migration, but use Go's package system (`internal/service/...`).
*   **Framework**: Use **Gin** for the HTTP layer (high performance, middleware support similar to Spring Interceptors).
*   **ORM**: Use **GORM** for database interaction (supports MySQL, easy migration from MyBatis).
*   **Config**: Use **Viper** to parse the existing `application.yml`.
*   **Compatibility**: Ensure JSON response structures (`ResultView`) and HTTP status codes match exactly.

### Steps

#### 1. Project Scaffolding & Infrastructure
1.  Initialize Go module (`go mod init netease-kit/nemo`).
2.  Create standard layout: `cmd/server`, `internal/config`, `internal/model`, `internal/service`, `internal/controller`, `pkg/utils`.
3.  Implement **Global Response Wrapper**: Create a Gin middleware/helper to replicate `@RestResponseBody` behavior, wrapping returns in `code`, `msg`, `data`.
4.  Implement **Context Middleware**: Port `Context.get()` logic to Gin middleware, extracting `AppKey`, `AppSecret`, and `Token` from headers and storing in `gin.Context`.
5.  Setup **Global Error Handling**: Create a middleware to catch custom errors (equivalent to `BsException`) and return standard error JSON.

#### 2. Data Layer Migration
1.  **Model Generation**: Use a tool (or manual mapping) to create Go structs from the `init.sql` schema. Ensure JSON tags match the Java DTOs if they differ from DB columns.
2.  **Database Connection**: Configure GORM with MySQL and `go-redis` for Redis.
3.  **Repository Pattern**: Implement data access functions in `internal/repo`, replacing MyBatis Mappers.

#### 3. Service & Logic Porting (Iterative)
1.  **Common Utilities**: Port `nemo-common` utils (Checksums, Encryption, DateUtils) to `pkg/utils`.
2.  **User Service**: Port `nemo-user-service`. Focus on `initAppAndUser` logic first.
3.  **EntLive Service**: Port `nemo-entlive-service` (Live records, Chorus, etc.).
4.  **Social & Game Services**: Port remaining logic.
5.  **External Integrations**: Port HTTP clients for NetEase IM/RTC server-side APIs.

#### 4. API Implementation
1.  **Controller Porting**: Create Gin handlers for each Spring `@RestController`.
    *   Example: `NemoInitController.java` -> `internal/controller/init_controller.go`.
2.  **Validation**: Use `go-playground/validator` to replace Java Bean Validation (`@Valid`).
3.  **Routing**: Define all routes in `internal/router/router.go`, matching the `@RequestMapping` paths exactly.

#### 5. Verification & Cutover
1.  **Integration Testing**: Create a test suite that calls both Java and Go endpoints with the same payload and asserts identical JSON responses.
2.  **Dockerization**: Create a `Dockerfile` for the Go binary (multi-stage build, scratch/alpine image).
3.  **Compose Update**: Update `docker-compose.yml` to swap the Java service with the Go service.

### Further Considerations
1.  **JSON Compatibility**: Java's `Long` often serializes to numbers, but sometimes strings in JS if too large. Go's `int64` handles this, but verify `Gson` vs `encoding/json` behavior for nulls/empty fields.
2.  **Redis Serialization**: Ensure the Go app reads/writes Redis keys using the same serialization format (likely JSON or String) as the Java app if they need to coexist during transition.
3.  **Async Tasks**: Identify any `@Async` or `Scheduled` tasks in Java and implement them using Go goroutines or a cron library (`robfig/cron`).
