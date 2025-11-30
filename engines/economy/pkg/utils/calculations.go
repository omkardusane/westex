package utils

import "math/rand/v2"

// makeRandomFloat64 generates a random float64 between min and max
func makeRandomFloat64(min, max float64) float64 {
	return min + (max-min)*rand.Float64()
}

func ProbableChance(probablity float32) bool {
	return rand.Float32() < probablity
}
