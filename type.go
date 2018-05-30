package main

import (
	"gonum.org/v1/gonum/mat"
)

// ColorMat -
type ColorMat []*mat.Dense

// NeighborLog -
type NeighborLog struct {
	i    int
	j    int
	dist float64
}
