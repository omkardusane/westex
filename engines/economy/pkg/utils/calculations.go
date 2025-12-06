package utils

import "math/rand/v2"

// makeRandomfloat32 generates a random float32 between min and max
func makeRandomfloat32(min, max float32) float32 {
	return min + (max-min)*rand.Float32()
}

func ProbableChance(probablity float32) bool {
	return rand.Float32() < probablity
}
