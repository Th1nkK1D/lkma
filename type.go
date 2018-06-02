package main

import (
	"gonum.org/v1/gonum/mat"
)

// ColorMat -
type ColorMat []*mat.Dense

// NeighbourLog -
type NeighbourLog struct {
	i    int
	j    int
	dist float64
}

// FloatMat -
type FloatMat [][]float64

// StateRes -
type StateRes struct {
	value  float64
	e      float64
	eTotal float64
	prop   float64
}
