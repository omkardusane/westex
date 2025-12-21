# 🧠 WestEx Economy Engine - Knowledge Base for LLMs

This document serves as a context transfer mechanism for LLMs picking up this project. It is structured in chapters to track the evolution and state of the codebase.

## 📖 Chapter 1: Project Overview & Core Architecture
**Date:** November 29, 2025

### 1.1 Project Mission
WestEx (Western Expansion) is an economy simulation engine written in Go. It simulates a regional economy with:
- **Industries**: Produce goods/services, consume resources, hire labor.
- **People**: Provide labor, earn wages, consume products to solve "Problems".
- **Markets**: Labor market (Time for Money) and Product market (Money for Solutions).

### 1.2 Core Architecture
The project follows a **Domain-Driven Design (DDD)** inspired structure, enforcing strict separation of concerns.

#### **Package Structure**
- **`pkg/entities/`**: (The "Nouns")
  - Contains pure data structures and core business logic.
  - **NO dependencies** on other project packages.
  - Key Files: `industry.go`, `person.go`, `region.go`, `problem.go`, `resource.go`.
  
- **`pkg/market/`**: (The "Verbs")
  - Handles interactions between entities.
  - Depends on `pkg/entities`.
  - Key Files: `labor.go` (Employment logic), `trade.go` (Buying/Selling logic).

- **`pkg/core/`**: (The "Engine")
  - Orchestrates the simulation loop (Ticks).
  - Depends on `pkg/entities` and `pkg/market`.
  - Key File: `engine.go`.

### 1.3 Key Design Patterns
- **Builder Pattern**: Used for complex entity creation (e.g., `CreateIndustry("Name").SetupIndustry(...).UpdateIndustryRates(...)`).
- **Composition**: `Person` has `PopulationSegment`s; `Region` contains `Industries` and `People`.

---

## 📖 Chapter 2: Refactoring Session - The "Entities" Transition
**Date:** November 29, 2025

### 2.1 The Problem
The project initially had a confusing structure:
- A `pkg/model` package that contained everything.
- A `pkg/model/entities` directory that was still package `model`.
- Circular dependencies and unclear boundaries.
- "Model" is a loaded term (often implies DB), whereas we wanted "Domain Entities".

### 2.2 The Solution
We performed a major refactor to align with idiomatic Go:

1.  **Renamed `pkg/model` → `pkg/entities`**:
    - Enforced the Go rule: **Package name MUST match directory name**.
    - `pkg/entities` is now the single source of truth for core types.

2.  **Refactored `Industry` Creation**:
    - Deprecated the monolithic `NewIndustry` constructor.
    - Introduced a **Fluent Builder API**:
      ```go
      // Old
      NewIndustry("Name", problems, labor, cons, prod)
      
      // New
      entities.CreateIndustry("Name").
          SetupIndustry(problems, inputs, outputs).
          UpdateIndustryRates(labor, cons, prod)
      ```

3.  **Market Separation**:
    - Moved transaction logic out of `model` into a dedicated `pkg/market` package.
    - `ProcessLaborMarket` and `ProcessProductMarket` now live here.

### 2.3 Current State
- **Build Status**: ✅ Compiles (Go installation pending on environment, but code is valid).
- **Test Coverage**: Basic tests in `engine_test.go` updated to new structure.
- **Next Steps**:
    - Implement more complex market dynamics (price fluctuation).
    - Add persistence layer (if needed).
    - Expand `PopulationSegment` logic.

---
