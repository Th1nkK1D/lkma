package main

import (
	"gonum.org/v1/gonum/mat"
)

// ColorMat -
type ColorMat []*mat.Dense
type NeighborLog struct {
	i    int
	j    int
	dist float64
}
