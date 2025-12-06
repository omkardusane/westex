# 🗺️ Economy Simulation Engine - Roadmap

## 🎯 Vision
Create a realistic economic simulation that models:
- **Real-world dynamics**: Supply chains, labor markets, resource scarcity
- **Human survival**: Basic needs vs pleasures, poverty, mortality
- **Government role**: Resource allocation, population sustainability
- **Economic complexity**: Industries, free resources, production chains

---

## 🏗️ Core Concepts

### **1. People (Agents)**
- Have **basic needs** (food, water, shelter) - survival requirements
- Have **pleasures** (entertainment, luxury) - quality of life
- **Earn money** through labor
- **Spend money** on needs and pleasures
- **Enter poverty** if needs unmet for X time
- **Die** if in poverty too long → affects population (natural death after N weeks and starvation after M weeks & Birth rate replenishes population)

### **2. Industries**
- **Consume resources** (inputs: labor, raw materials)
- **Produce goods** (outputs: food, products, services)
- **Pay wages** to workers
- **Set prices** = cost of production + 10% markup
- **Linked in supply chains** (e.g., agriculture → food processing → retail)

### **3. Free Resources**
- **Limited quantity**: Land, water, minerals, forests
- **Allocated by government**: Leasing, rights, permits
- **Depleted over time**: Mining, harvesting, usage
- **Essential for production**: Industries need access

### **4. Government**
- **Goal**: Keep population stable/growing
- **Powers**:
  - Allocate free resources to industries
  - Collect taxes (future)
  - Provide welfare (future)
  - Regulate industries (future)
- **Metrics**: Population growth, poverty rate, resource depletion

### **5. Markets**
- **Labor Market**: People sell labor, industries buy
- **Product Market**: People buy goods, industries sell
- **Resource Market**: Government allocates, industries lease

---

## 📅 Implementation Phases

### **Phase 1: Foundation** ✅ (CURRENT)
**Status**: In Progress  
**Goal**: Basic production and wealth tracking

- [x] Project structure (entities, core, utils)
- [x] Builder pattern for entities
- [x] Demand-based production targets
- [x] Initial wealth tracking
- [x] Population segments
- [ ] **Production phase with labor payments** ← NEXT
    - [ ] Extract production calculation to separate function

---

### **Phase 2: Labor & Production** 🎯 (NEXT - Weeks 1-2)
**Goal**: Complete production cycle with labor payments

#### **2.1 Production Phase** (Week 1)
- [ ] Industries consume input resources
- [ ] Industries produce output goods
- [ ] **Industries pay wages to workers immediately**
- [ ] Track labor hours used per industry
- [ ] Calculate production costs (labor + resources)

#### **2.2 Labor Market** (Week 1-2)
- [ ] Worker population segment
- [ ] Employment allocation (who works where)
- [ ] Wage calculation (hours × wage rate)
- [ ] Unemployment tracking
- [ ] Labor hours per person per tick

**Deliverable**: Industries produce goods and pay workers

---

### **Phase 3: Trade & Consumption** (Weeks 3-4)
**Goal**: People buy goods to satisfy needs

#### **3.1 Pricing System**
- [ ] Calculate production cost per unit
- [ ] Set price = production cost × 1.10 (10% markup)
- [ ] Track prices per product

#### **3.2 Product Market**
- [ ] People identify their needs (basic vs pleasure)
- [ ] People buy products based on:
  - Need severity
  - Spending power (money available)
  - Product availability
- [ ] Industries sell products and earn revenue

#### **3.3 Needs & Satisfaction**
- [ ] Track need satisfaction per person
- [ ] Differentiate basic needs vs pleasures
- [ ] Unsatisfied needs accumulate over time

**Deliverable**: Complete economic cycle (produce → pay → buy → consume)

---

### **Phase 4: Survival & Mortality** (Weeks 5-6)
**Goal**: Realistic consequences of poverty

#### **4.1 Poverty Mechanics**
- [ ] Define poverty threshold (e.g., can't afford basic needs)
- [ ] Track time in poverty per person
- [ ] Poverty effects:
  - Reduced labor productivity
  - Health deterioration
  - Death after X ticks in poverty

#### **4.2 Population Dynamics**
- [ ] Death removes person from simulation
- [ ] Track population over time
- [ ] Birth rate (future - Phase 6)
- [ ] Migration (future - Phase 7)

**Deliverable**: Population can decline if economy fails

---

### **Phase 5: Free Resources** (Weeks 7-8)
**Goal**: Limited resources that industries need

#### **5.1 Resource Types**
- [ ] Land (for agriculture, buildings)
- [ ] Water (for production, consumption)
- [ ] Minerals (for manufacturing)
- [ ] Forests (for lumber, paper)

#### **5.2 Resource Mechanics**
- [ ] Finite quantities
- [ ] Depletion over time
- [ ] Regeneration (e.g., forests regrow slowly)
- [ ] Quality/accessibility (some resources easier to extract)

#### **5.3 Resource Allocation**
- [ ] Industries request resources
- [ ] Government allocates based on:
  - Industry priority (food > luxury)
  - Payment/lease fees
  - Sustainability goals

**Deliverable**: Industries compete for limited resources

---

### **Phase 6: Government** (Weeks 9-10)
**Goal**: Governing body that manages economy

#### **6.1 Resource Management**
- [ ] Leasing system for land, minerals, etc.
- [ ] Lease prices (revenue for government)
- [ ] Allocation algorithm (priority-based)

#### **6.2 Population Goal**
- [ ] Track population trend
- [ ] Government interventions:
  - Subsidize food if poverty high
  - Restrict resource extraction if depleting
  - Incentivize employment

#### **6.3 Basic Governance**
- [ ] Government budget (from leases, future taxes)
- [ ] Spending on welfare, infrastructure
- [ ] Policy decisions (simple rules-based)

**Deliverable**: Government actively manages economy

---

### **Phase 7: Supply Chains** (Weeks 11-12)
**Goal**: Industries depend on each other

#### **7.1 Industry Dependencies**
- [ ] Agriculture produces raw food
- [ ] Food processing converts raw → packaged food
- [ ] Retail sells to consumers
- [ ] Manufacturing needs minerals from mining

#### **7.2 Supply Chain Mechanics**
- [ ] Industries buy from other industries
- [ ] B2B transactions
- [ ] Supply shortages ripple through chain
- [ ] Price propagation (input costs affect output prices)

**Deliverable**: Complex, interconnected economy

---

### **Phase 8: Advanced Features** (Weeks 13+)
**Goal**: Sophistication and realism

- [ ] Multiple regions with trade
- [ ] Technology/productivity improvements
- [ ] Education system (skilled vs unskilled labor)
- [ ] Financial system (loans, credit, banks)
- [ ] Taxation system
- [ ] Inflation/deflation
- [ ] Economic cycles (boom/bust)

---

## 🎯 Near-Term Actionables (Next 2 Weeks)

### **Week 1: Production with Labor Payments**

#### **Day 1-2: Refactor Production**
- [ ] Extract production calculation to `calculateProduction()` function
- [ ] Add labor cost tracking
- [ ] Implement immediate wage payments in production phase

#### **Day 3-4: Resource Consumption**
- [ ] Industries consume input resources during production
- [ ] Check resource availability before producing
- [ ] Track resource depletion

#### **Day 5-7: Labor Market Foundation**
- [ ] Create worker pool from population
- [ ] Allocate workers to industries
- [ ] Track employment status per person

### **Week 2: Complete Production Cycle**

#### **Day 1-3: Production Costs**
- [ ] Calculate total production cost (labor + resources)
- [ ] Track cost per unit produced
- [ ] Store historical costs

#### **Day 4-5: Pricing System**
- [ ] Implement cost-plus pricing (cost × 1.10)
- [ ] Update prices each tick based on costs
- [ ] Display prices in logs

#### **Day 6-7: Testing & Refinement**
- [ ] Test with different scenarios
- [ ] Balance parameters (wages, production rates)
- [ ] Add logging for debugging

---

## 📊 Success Metrics (Future - Phase 8+)

### **Economic Health**
- **GDP**: Total value of goods produced
- **Unemployment rate**: % of workers without jobs
- **Inflation rate**: Price changes over time
- **Gini coefficient**: Wealth inequality

### **Population Health**
- **Population growth rate**: Births - deaths
- **Poverty rate**: % below poverty line
- **Average lifespan**: How long people survive
- **Needs satisfaction**: % of needs met

### **Resource Health**
- **Resource depletion rate**: How fast resources consumed
- **Sustainability index**: Regeneration vs consumption
- **Resource scarcity**: Availability vs demand

### **Government Performance**
- **Budget balance**: Revenue vs spending
- **Policy effectiveness**: Impact of interventions
- **Population satisfaction**: Aggregate wellbeing

---

## 🏛️ Architecture Principles

### **1. Modularity**
Each phase should be a separate, testable module:
- `pkg/entities` - Data structures
- `pkg/production` - Production logic
- `pkg/labor` - Labor market
- `pkg/trade` - Product market
- `pkg/resources` - Resource management
- `pkg/government` - Government logic

### **2. Configurability**
All parameters should be configurable:
```go
type EconomyConfig struct {
    WageRate           float32
    ProfitMargin       float32
    PovertyThreshold   float32
    TicksToStarvation  int
    ResourceRegenRate  float32
}
```

### **3. Observability**
Track everything for analysis:
- Tick-by-tick logs
- Historical data storage
- Metrics calculation
- Visualization-ready output (JSON/CSV)

### **4. Realism**
Prioritize realistic behavior:
- People act rationally (buy cheapest, work for highest wage)
- Industries maximize profit
- Resources are finite
- Government has constraints

---

## 🔄 Iteration Strategy

For each phase:
1. **Design**: Define data structures and interfaces
2. **Implement**: Write core logic
3. **Test**: Create scenarios to verify behavior
4. **Balance**: Tune parameters for realism
5. **Document**: Update docs with new features
6. **Refactor**: Clean up before moving to next phase

---

## 📚 Documentation Plan

- `ROADMAP.md` (this file) - Long-term vision
- `CURRENT_PHASE.md` - Detailed current work
- `DESIGN_DECISIONS.md` - Why we chose X over Y
- `PARAMETERS.md` - All configurable values
- `SCENARIOS.md` - Test scenarios and expected outcomes

---

## 🎮 Example Scenarios to Test

### **Scenario 1: Balanced Economy**
- 1000 people, 5 industries
- Sufficient resources
- Expected: Stable population, low poverty

### **Scenario 2: Resource Scarcity**
- Limited land for agriculture
- Expected: Food prices rise, some starvation

### **Scenario 3: Unemployment Crisis**
- More workers than jobs
- Expected: Poverty increases, population declines

### **Scenario 4: Supply Chain Failure**
- Mining industry fails → manufacturing stops
- Expected: Cascading failures

---

## 🚀 Getting Started

**Right now, focus on:**
1. Implementing labor payments in production phase
2. Extracting production calculation
3. Setting up labor market foundation

**Next document to create:** `CURRENT_PHASE.md` with detailed tasks for Week 1

---

**This is an ambitious, realistic simulation. Take it one phase at a time, and we'll build something incredible!** 🌟
