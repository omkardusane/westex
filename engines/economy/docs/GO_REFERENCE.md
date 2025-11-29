# 📖 Go Quick Reference Guide

A cheat sheet for Go concepts with examples from your economy simulation engine.

---

## 🎯 Table of Contents

1. [Package & Imports](#package--imports)
2. [Variables & Types](#variables--types)
3. [Structs & Methods](#structs--methods)
4. [Interfaces](#interfaces)
5. [Error Handling](#error-handling)
6. [Collections](#collections)
7. [Concurrency](#concurrency)
8. [Testing](#testing)
9. [Common Patterns](#common-patterns)

---

## Package & Imports

```go
// Package declaration (must be first)
package model

// Import single package
import "fmt"

// Import multiple packages
import (
    "fmt"
    "errors"
    "simulation-engine/pkg/model"
)

// Import with alias
import (
    m "simulation-engine/pkg/model"
    log "simulation-engine/pkg/logging"
)

// Blank import (for side effects only)
import _ "github.com/lib/pq"
```

### Visibility Rules
```go
// Exported (public) - starts with capital letter
type Region struct { }
func NewRegion() *Region { }

// Unexported (private) - starts with lowercase
type internalCache struct { }
func calculateTotal() float64 { }
```

---

## Variables & Types

### Declaration

```go
// Declare and initialize
var name string = "Alice"

// Type inference
var age = 25

// Short declaration (inside functions only)
count := 100

// Multiple variables
var (
    x int = 10
    y int = 20
)

// Constants
const MaxPlayers = 100
const (
    StatusActive   = "active"
    StatusInactive = "inactive"
)
```

### Basic Types

```go
// Numbers
var i int = 42
var f float64 = 3.14
var u uint = 100

// Strings
var s string = "hello"
var multiline = `This is a
multi-line string`

// Booleans
var b bool = true

// Pointers
var p *int = &i
```

### Type Conversion

```go
var i int = 42
var f float64 = float64(i)
var u uint = uint(i)
```

---

## Structs & Methods

### Struct Definition

```go
type Person struct {
    Name       string      // Exported field
    age        int         // Unexported field
    Money      float64
    Segments   []*Segment  // Slice of pointers
}
```

### Constructor Pattern

```go
// Constructor function (convention: NewTypeName)
func NewPerson(name string, money, laborHours float64) *Person {
    return &Person{
        Name:       name,
        Money:      money,
        LaborHours: laborHours,
        Segments:   make([]*Segment, 0),  // Initialize slice
    }
}
```

### Methods

```go
// Value receiver - receives a copy
func (p Person) GetName() string {
    return p.Name
}

// Pointer receiver - receives the original
func (p *Person) SetName(name string) {
    p.Name = name
}

// Method with return value
func (p Person) CanAfford(amount float64) bool {
    return p.Money >= amount
}

// Method with error return
func (p *Person) SpendMoney(amount float64) error {
    if !p.CanAfford(amount) {
        return fmt.Errorf("insufficient funds")
    }
    p.Money -= amount
    return nil
}
```

### Embedding (Composition)

```go
type Employee struct {
    Person          // Embedded struct
    EmployeeID int
}

// Can access Person fields directly
emp := Employee{
    Person:     Person{Name: "Alice"},
    EmployeeID: 123,
}
fmt.Println(emp.Name)  // Access embedded field
```

---

## Interfaces

### Definition

```go
// Interface defines behavior
type Logger interface {
    Log(message string)
    LogError(err error)
}

// Empty interface (any type)
var anything interface{}
anything = 42
anything = "hello"
anything = &Person{}
```

### Implementation (Implicit)

```go
// No "implements" keyword needed
type ConsoleLogger struct{}

func (c *ConsoleLogger) Log(message string) {
    fmt.Println(message)
}

func (c *ConsoleLogger) LogError(err error) {
    fmt.Println("ERROR:", err)
}

// ConsoleLogger now implements Logger automatically
```

### Interface Composition

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Compose interfaces
type ReadWriter interface {
    Reader
    Writer
}
```

### Type Assertion

```go
var i interface{} = "hello"

// Type assertion
s := i.(string)
fmt.Println(s)

// Safe type assertion
s, ok := i.(string)
if ok {
    fmt.Println(s)
}

// Type switch
switch v := i.(type) {
case string:
    fmt.Println("string:", v)
case int:
    fmt.Println("int:", v)
default:
    fmt.Println("unknown type")
}
```

---

## Error Handling

### Basic Pattern

```go
result, err := DoSomething()
if err != nil {
    return err  // Propagate error
}
// Use result
```

### Creating Errors

```go
// Simple error
err := fmt.Errorf("invalid input: %s", input)

// From errors package
err := errors.New("something went wrong")

// Wrapping errors (Go 1.13+)
err := fmt.Errorf("failed to process: %w", originalErr)
```

### Custom Error Types

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Usage
return &ValidationError{
    Field:   "email",
    Message: "invalid format",
}
```

### Error Checking

```go
// Check if error is specific type
var validationErr *ValidationError
if errors.As(err, &validationErr) {
    fmt.Println("Validation error:", validationErr.Field)
}

// Check if error is specific value
if errors.Is(err, ErrNotFound) {
    fmt.Println("Not found")
}
```

---

## Collections

### Arrays (Fixed Size)

```go
// Declare array
var arr [5]int
arr[0] = 1

// Initialize
arr := [5]int{1, 2, 3, 4, 5}

// Let compiler count
arr := [...]int{1, 2, 3}
```

### Slices (Dynamic)

```go
// Create slice
var s []int

// Initialize
s := []int{1, 2, 3}

// Make with length and capacity
s := make([]int, 5)      // len=5, cap=5
s := make([]int, 0, 10)  // len=0, cap=10

// Append
s = append(s, 4)
s = append(s, 5, 6, 7)

// Slice operations
s[1:3]   // Elements 1 and 2
s[:3]    // First 3 elements
s[2:]    // From element 2 to end
s[:]     // All elements (copy)

// Length and capacity
len(s)
cap(s)
```

### Maps

```go
// Create map
var m map[string]int

// Initialize
m := make(map[string]int)

// Literal
m := map[string]int{
    "alice": 100,
    "bob":   200,
}

// Set
m["charlie"] = 300

// Get
value := m["alice"]

// Get with existence check
value, exists := m["alice"]
if exists {
    fmt.Println(value)
}

// Delete
delete(m, "alice")

// Iterate
for key, value := range m {
    fmt.Println(key, value)
}
```

### Iteration

```go
// Slice/Array
for i, value := range slice {
    fmt.Println(i, value)
}

// Just index
for i := range slice {
    fmt.Println(i)
}

// Just value
for _, value := range slice {
    fmt.Println(value)
}

// Map
for key, value := range myMap {
    fmt.Println(key, value)
}

// Traditional for loop
for i := 0; i < 10; i++ {
    fmt.Println(i)
}
```

---

## Concurrency

### Goroutines

```go
// Start goroutine
go doSomething()

// Anonymous function
go func() {
    fmt.Println("Running in goroutine")
}()

// With parameters
go func(name string) {
    fmt.Println("Hello", name)
}("Alice")
```

### Channels

```go
// Create channel
ch := make(chan int)

// Buffered channel
ch := make(chan int, 10)

// Send
ch <- 42

// Receive
value := <-ch

// Receive with ok
value, ok := <-ch
if !ok {
    fmt.Println("Channel closed")
}

// Close channel
close(ch)

// Range over channel
for value := range ch {
    fmt.Println(value)
}
```

### Select

```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case msg := <-ch2:
    fmt.Println("Received from ch2:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout")
default:
    fmt.Println("No message")
}
```

### WaitGroup

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Println("Worker", id)
    }(i)
}

wg.Wait()  // Wait for all goroutines
```

### Mutex

```go
var (
    mu    sync.Mutex
    count int
)

func increment() {
    mu.Lock()
    defer mu.Unlock()
    count++
}
```

---

## Testing

### Basic Test

```go
package model

import "testing"

func TestPerson_CanAfford(t *testing.T) {
    person := NewPerson("Alice", 100.0, 8.0)
    
    if !person.CanAfford(50.0) {
        t.Error("should be able to afford 50")
    }
    
    if person.CanAfford(150.0) {
        t.Error("should not be able to afford 150")
    }
}
```

### Table-Driven Tests

```go
func TestCalculateWage(t *testing.T) {
    tests := []struct {
        name     string
        hours    float64
        rate     float64
        expected float64
    }{
        {"standard", 8.0, 10.0, 80.0},
        {"overtime", 12.0, 10.0, 120.0},
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

### Test Helpers

```go
func TestMain(m *testing.M) {
    // Setup
    setup()
    
    // Run tests
    code := m.Run()
    
    // Teardown
    teardown()
    
    os.Exit(code)
}

func setup() {
    // Initialize test data
}

func teardown() {
    // Clean up
}
```

### Benchmarks

```go
func BenchmarkProcessProduction(b *testing.B) {
    engine := setupEngine()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        engine.processProduction()
    }
}

// Run: go test -bench=.
```

---

## Common Patterns

### Constructor Pattern

```go
func NewPerson(name string, money float64) *Person {
    return &Person{
        Name:     name,
        Money:    money,
        Segments: make([]*Segment, 0),
    }
}
```

### Functional Options

```go
type Option func(*Engine)

func WithLogger(logger Logger) Option {
    return func(e *Engine) {
        e.Logger = logger
    }
}

func WithWageRate(rate float64) Option {
    return func(e *Engine) {
        e.WagePerHour = rate
    }
}

func NewEngine(region *Region, opts ...Option) *Engine {
    e := &Engine{
        Region:      region,
        WagePerHour: 10.0,  // Default
    }
    
    for _, opt := range opts {
        opt(e)
    }
    
    return e
}

// Usage
engine := NewEngine(
    region,
    WithLogger(logger),
    WithWageRate(15.0),
)
```

### Builder Pattern

```go
type RegionBuilder struct {
    region *Region
}

func NewRegionBuilder(name string) *RegionBuilder {
    return &RegionBuilder{
        region: NewRegion(name),
    }
}

func (b *RegionBuilder) AddIndustry(industry *Industry) *RegionBuilder {
    b.region.AddIndustry(industry)
    return b
}

func (b *RegionBuilder) AddPerson(person *Person) *RegionBuilder {
    b.region.AddPerson(person)
    return b
}

func (b *RegionBuilder) Build() *Region {
    return b.region
}

// Usage
region := NewRegionBuilder("Silicon Valley").
    AddIndustry(farmCo).
    AddIndustry(funZone).
    AddPerson(alice).
    Build()
```

### Singleton Pattern

```go
var (
    instance *Config
    once     sync.Once
)

func GetConfig() *Config {
    once.Do(func() {
        instance = &Config{
            // Initialize
        }
    })
    return instance
}
```

---

## 🎯 Quick Tips

### When to Use Pointers

✅ **Use pointers when:**
- Modifying the receiver
- Struct is large
- Consistency (if one method uses pointer, all should)

❌ **Don't use pointers when:**
- Struct is small (int, bool, small structs)
- You want immutability
- Method doesn't modify receiver

### Naming Conventions

```go
// Variables: camelCase
var userName string
var totalCount int

// Exported: PascalCase
type Person struct { }
func NewPerson() { }

// Unexported: camelCase
type internalCache struct { }
func calculateTotal() { }

// Acronyms: all caps or all lowercase
var userID int    // Good
var userHTTPClient // Good
var userId int    // Bad
```

### Common Mistakes

```go
// ❌ Forgetting to check errors
result, _ := DoSomething()  // Don't ignore errors!

// ✅ Always check
result, err := DoSomething()
if err != nil {
    return err
}

// ❌ Modifying slice in range
for _, item := range items {
    item.Value = 10  // Modifies copy, not original!
}

// ✅ Use index
for i := range items {
    items[i].Value = 10
}

// ❌ Closing channel from receiver
close(ch)  // Should be done by sender

// ✅ Sender closes
go func() {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch)  // Sender closes
}()
```

---

## 📚 Resources

- [Effective Go](https://golang.org/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Go Playground](https://play.golang.org/)
- [Go Standard Library](https://pkg.go.dev/std)

---

**Keep this reference handy while working through the exercises!**
