# Economy Simulation Engine

A high-level economy simulation engine built in Go. This engine simulates interactions between regions, industries, people, problems, and resources.

## 🏗️ Architecture

This project follows a clean, modular architecture designed to be the core engine for future web or 3D game applications.

### Project Structure

```
/economy
├── /cmd
│   └── /sim-cli              # CLI entrypoint
│       └── main.go
├── /pkg                      # Core reusable library
│   ├── /core                 # Simulation engine
│   │   └── engine.go
│   ├── /model                # Data models
│   │   ├── problem.go
│   │   ├── resource.go
│   │   ├── entity.go
│   │   └── region.go
│   ├── /interaction          # Business logic
│   │   ├── trade.go
│   │   └── labor.go
│   └── /logging              # Utilities
│       └── logger.go
└── go.mod
```

## 🎮 Entities

- **Region**: Container for all entities in a geographic/economic area
- **Industry**: Produces goods/services, solves 1-2 problems
- **Person**: Individual with needs, belongs to population segments
- **Problem**: High-level needs (food, water, entertainment, etc.)
- **Resource**: Materials/products with quantities that change over time

## 🔄 Simulation Flow

Each tick of the simulation runs through these phases:

1. **Production Phase**: Industries produce goods
2. **Labor Market Phase**: People work for industries and earn wages
3. **Product Market Phase**: People buy products to solve their problems
4. **Reset**: Labor hours reset for the next tick

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher

### Installation

```bash
# Navigate to the project directory
cd ./westex/engines/economy

# Run the simulation
go run ./cmd/sim-cli
```

### Running Tests

```bash
go test ./...
```

## 📊 Example Output

The simulation logs all interactions:
- Industries producing goods
- People working and earning wages
- People buying products
- Final economic summary with wealth distribution

## 🎯 Project demonstrates

This project demonstrates:

1. **Project Structure**: Clean separation of concerns (cmd, pkg, internal)
2. **Modularity**: Reusable packages that can be imported by different frontends
3. **Data Types**: Proper use of structs, pointers, and interfaces
4. **Entrypoints**: Multiple possible entrypoints (CLI, web server, etc.)
5. **Testing**: Unit tests for business logic
6. **APIs**: Clear public APIs in pkg/ for external consumption
7. **Concurrency**: (Future) Goroutines for parallel simulation
8. **Performance**: (Future) Profiling and optimization

## 🔮 Future Enhancements

- [ ] REST API server (web game backend)
- [ ] WebSocket support for real-time updates
- [ ] Persistence layer (database)
- [ ] Concurrency for large-scale simulations
- [ ] Event sourcing for replay capability
- [ ] Metrics and monitoring
- [ ] Configuration files (YAML/JSON)
- [ ] Plugin system for custom industries/behaviors

## 📚 Go Concepts Covered

- Package organization
- Struct methods
- Pointers vs values
- Slices and maps
- Encapsulation (public/private)
- Constructor patterns
- Error handling
- Logging and observability

## 🏛️ Architectural Principles

1. **Separation of Concerns**: Engine logic is independent of UI
2. **Dependency Inversion**: Core doesn't depend on external systems
3. **Single Responsibility**: Each package has one clear purpose
4. **Open/Closed**: Easy to extend with new entity types
5. **Interface Segregation**: (Future) Small, focused interfaces

---

Built with ❤️ to learn Go and game engine architecture
