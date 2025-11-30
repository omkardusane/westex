# Go Architecture Guide: Path to Principal Architect

This document explains the architectural decisions in this project and the Go concepts you need to master to become a Principal Architect.

## 📁 Project Structure Philosophy

### The Standard Go Layout

```
/project-root
├── /cmd                    # Executable commands (main packages)
│   └── /app-name          # One directory per executable
│       └── main.go        # Entry point
├── /pkg                    # Public libraries (importable by external projects)
│   └── /feature           # Organized by feature/domain
├── /internal              # Private application code (not importable externally)
│   └── /feature
├── /api                    # API definitions (OpenAPI, gRPC, etc.)
├── /web                    # Web assets (if applicable)
├── /configs               # Configuration files
├── /scripts               # Build and deployment scripts
├── /test                  # Additional test data
└── go.mod                 # Module definition
```

### Why This Structure?

1. **`/cmd`**: Multiple entry points for different use cases
   - CLI tool: `/cmd/cli`
   - Web server: `/cmd/server`
   - Worker: `/cmd/worker`
   
2. **`/pkg`**: Reusable, well-tested library code
   - Can be imported by external projects
   - Must have stable, well-documented APIs
   - Think of this as your "SDK"

3. **`/internal`**: Application-specific code
   - Go compiler enforces: cannot be imported by external projects
   - Use for implementation details
   - Protects your internal APIs from external dependencies

## 🎯 Key Go Concepts for Architects

### 1. **Modules and Packages**

```go
// go.mod defines your module
module simulation-engine

go 1.21

// Import paths are based on module name
import "simulation-engine/pkg/core"
```

**Best Practices:**
- One module per repository (usually)
- Package names should be short, lowercase, no underscores
- Package should have a single, clear purpose

### 2. **Interfaces: The Heart of Go Architecture**

```go
// Define behavior, not implementation
type Engine interface {
    Run(ticks int)
    GetState() State
}

// Multiple implementations
type SimpleEngine struct { ... }
type DistributedEngine struct { ... }

// Both satisfy the interface
func (e *SimpleEngine) Run(ticks int) { ... }
func (e *DistributedEngine) Run(ticks int) { ... }
```

**Why This Matters:**
- Dependency Inversion: depend on interfaces, not concrete types
- Testing: easy to mock
- Flexibility: swap implementations without changing clients

### 3. **Dependency Injection**

```go
// BAD: Hard-coded dependency
type Service struct {
    db *PostgresDB  // Tightly coupled
}

// GOOD: Inject interface
type Service struct {
    db Database  // Interface
}

// Constructor with DI
func NewService(db Database) *Service {
    return &Service{db: db}
}
```

### 4. **Error Handling**

```go
// Go's explicit error handling
func DoSomething() (Result, error) {
    if err := validate(); err != nil {
        return Result{}, fmt.Errorf("validation failed: %w", err)
    }
    return result, nil
}

// Caller must handle
result, err := DoSomething()
if err != nil {
    log.Printf("error: %v", err)
    return
}
```

**Best Practices:**
- Always check errors
- Wrap errors with context: `fmt.Errorf("context: %w", err)`
- Use custom error types for specific cases
- Don't panic in libraries (only in main)

### 5. **Concurrency Patterns**

```go
// Goroutines: lightweight threads
go processInBackground()

// Channels: communication between goroutines
ch := make(chan Result)
go func() {
    ch <- computeResult()
}()
result := <-ch

// Worker pool pattern
func WorkerPool(jobs <-chan Job, results chan<- Result, workers int) {
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    wg.Wait()
    close(results)
}
```

### 6. **Context for Cancellation and Timeouts**

```go
func ProcessWithTimeout(ctx context.Context, data Data) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    select {
    case result := <-process(data):
        return result
    case <-ctx.Done():
        return ctx.Err() // Timeout or cancellation
    }
}
```

## 🏗️ Architectural Patterns

### 1. **Layered Architecture**

```
Presentation Layer (cmd/)
    ↓
Business Logic Layer (pkg/core, pkg/interaction)
    ↓
Data Layer (pkg/model)
```

### 2. **Repository Pattern**

```go
// Abstract data access
type RegionRepository interface {
    Save(region *Region) error
    FindByID(id string) (*Region, error)
    FindAll() ([]*Region, error)
}

// Implementations
type InMemoryRegionRepo struct { ... }
type PostgresRegionRepo struct { ... }
```

### 3. **Service Layer**

```go
// Business logic encapsulation
type EconomyService struct {
    regionRepo RegionRepository
    logger     Logger
}

func (s *EconomyService) SimulateEconomy(regionID string, ticks int) error {
    region, err := s.regionRepo.FindByID(regionID)
    if err != nil {
        return err
    }
    // Business logic here
    return s.regionRepo.Save(region)
}
```

## 🧪 Testing Strategies

### 1. **Unit Tests**

```go
func TestEngine_Run(t *testing.T) {
    // Arrange
    engine := NewEngine(mockRegion, 10.0, 2.0, 50.0)
    
    // Act
    engine.Run(5)
    
    // Assert
    if engine.CurrentTick != 5 {
        t.Errorf("expected 5 ticks, got %d", engine.CurrentTick)
    }
}
```

### 2. **Table-Driven Tests**

```go
func TestCalculateWage(t *testing.T) {
    tests := []struct {
        name     string
        hours    float64
        rate     float64
        expected float64
    }{
        {"standard", 8.0, 10.0, 80.0},
        {"overtime", 10.0, 15.0, 150.0},
        {"zero", 0.0, 10.0, 0.0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateWage(tt.hours, tt.rate)
            if result != tt.expected {
                t.Errorf("got %.2f, want %.2f", result, tt.expected)
            }
        })
    }
}
```

### 3. **Mocking with Interfaces**

```go
type MockLogger struct {
    LoggedMessages []string
}

func (m *MockLogger) LogEvent(msg string) {
    m.LoggedMessages = append(m.LoggedMessages, msg)
}

func TestWithMock(t *testing.T) {
    mock := &MockLogger{}
    service := NewService(mock)
    service.DoSomething()
    
    if len(mock.LoggedMessages) != 1 {
        t.Error("expected 1 log message")
    }
}
```

## 🚀 API Design

### REST API Example (Future Enhancement)

```go
// /cmd/server/main.go
func main() {
    router := chi.NewRouter()
    
    // Middleware
    router.Use(middleware.Logger)
    router.Use(middleware.Recoverer)
    
    // Routes
    router.Post("/api/v1/simulations", createSimulation)
    router.Get("/api/v1/simulations/{id}", getSimulation)
    router.Post("/api/v1/simulations/{id}/run", runSimulation)
    
    http.ListenAndServe(":8080", router)
}

// Handler
func createSimulation(w http.ResponseWriter, r *http.Request) {
    var req CreateSimulationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Business logic
    sim := service.CreateSimulation(req)
    
    json.NewEncoder(w).Encode(sim)
}
```

### gRPC API Example

```protobuf
// api/simulation.proto
syntax = "proto3";

service SimulationService {
    rpc RunSimulation(RunRequest) returns (RunResponse);
    rpc GetState(GetStateRequest) returns (State);
}

message RunRequest {
    string region_id = 1;
    int32 ticks = 2;
}
```

## 📊 Performance and Optimization

### 1. **Profiling**

```go
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    // Your application code
}

// Access profiles at http://localhost:6060/debug/pprof/
```

### 2. **Benchmarking**

```go
func BenchmarkProcessProduction(b *testing.B) {
    engine := setupEngine()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        engine.processProduction()
    }
}

// Run: go test -bench=. -benchmem
```

### 3. **Memory Management**

```go
// Use sync.Pool for frequently allocated objects
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func process() {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf)
    buf.Reset()
    // Use buffer
}
```

## 🔐 Configuration Management

```go
// /configs/config.go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Logging  LoggingConfig
}

// Load from environment or file
func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.AddConfigPath("./configs")
    viper.AutomaticEnv()
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

## 📚 Learning Path to Principal Architect

### Phase 1: Fundamentals (You are here!)
- ✅ Project structure
- ✅ Packages and modules
- ✅ Structs and methods
- ✅ Basic error handling
- ✅ Testing basics

### Phase 2: Intermediate
- [ ] Interfaces and polymorphism
- [ ] Dependency injection
- [ ] Advanced error handling (custom errors)
- [ ] Table-driven tests
- [ ] HTTP servers and routing

### Phase 3: Advanced
- [ ] Concurrency patterns (goroutines, channels)
- [ ] Context usage
- [ ] Middleware patterns
- [ ] Database integration (SQL, NoSQL)
- [ ] Caching strategies

### Phase 4: Architecture
- [ ] Microservices design
- [ ] Event-driven architecture
- [ ] CQRS and Event Sourcing
- [ ] API design (REST, gRPC, GraphQL)
- [ ] Service mesh concepts

### Phase 5: Production
- [ ] Observability (metrics, tracing, logging)
- [ ] Performance optimization
- [ ] Security best practices
- [ ] CI/CD pipelines
- [ ] Deployment strategies (Docker, Kubernetes)

## 🎓 Recommended Resources

1. **Books:**
   - "The Go Programming Language" by Donovan & Kernighan
   - "Go in Action" by William Kennedy
   - "Cloud Native Go" by Matthew Titmus

2. **Online:**
   - [Effective Go](https://golang.org/doc/effective_go)
   - [Go Blog](https://blog.golang.org/)
   - [Go by Example](https://gobyexample.com/)

3. **Practice:**
   - Build this simulation engine further
   - Contribute to open-source Go projects
   - Implement design patterns in Go

---

**Next Steps for This Project:**
1. Add interfaces for `Engine`, `Repository`, etc.
2. Implement persistence (save/load simulations)
3. Add a REST API server
4. Implement concurrency for large simulations
5. Add comprehensive logging and metrics
