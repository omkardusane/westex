package entities

type Demand struct {
	Problem   *Problem
	severity  float32 // how critical this problem is
	demand    float32 // calculated demand recorded
	stability float32 // how stable the demand is over time
}
