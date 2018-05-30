package main

import (
	"math"

	"gocv.io/x/gocv"

	"gonum.org/v1/gonum/mat"
)

func getDist(i, j, x, y int) float64 {
	return math.Sqrt(math.Pow(float64(i-x), 2) + math.Pow(float64(j-y), 2))
}

func explore(neighbors [][]NeighborLog, mask *mat.Dense, i, j, si, sj, w, h int, start bool) [][]NeighborLog {
	dist := getDist(i, j, si, sj)

	if start || (mask.At(i, j) == 0 && (neighbors[i][j].dist == 0 || dist < neighbors[i][j].dist)) {
		// Update value
		if mask.At(i, j) == 0 {
			neighbors[i][j].i = si
			neighbors[i][j].j = sj
			neighbors[i][j].dist = dist
		}

		// Water flow
		if i+1 < h {
			neighbors = explore(neighbors, mask, i+1, j, si, sj, w, h, false)
		}
		if j+1 < w {
			neighbors = explore(neighbors, mask, i, j+1, si, sj, w, h, false)
		}
		if i-1 >= 0 {
			neighbors = explore(neighbors, mask, i-1, j, si, sj, w, h, false)
		}
		if j-1 >= 0 {
			neighbors = explore(neighbors, mask, i, j-1, si, sj, w, h, false)
		}
	}

	return neighbors
}

// ExploreNeighbor -
func ExploreNeighbor(mask *mat.Dense) ([][]NeighborLog, [][]NeighborLog) {
	nRow, nCol := mask.Dims()

	// Init neighbor logs
	neighborFG, neighborBG := make([][]NeighborLog, nRow), make([][]NeighborLog, nRow)

	for n := range neighborFG {
		neighborFG[n] = make([]NeighborLog, nCol)
		neighborBG[n] = make([]NeighborLog, nCol)
	}

	// Start
	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			if mask.At(i, j) == 1 {
				// Marked FG
				neighborFG = explore(neighborFG, mask, i, j, i, j, nCol, nRow, true)
			} else if mask.At(i, j) == -1 {
				// Marked BG
				neighborBG = explore(neighborBG, mask, i, j, i, j, nCol, nRow, true)
			}
		}
	}

	return neighborFG, neighborBG
}

// SaveNeighborLog -
func SaveNeighborLog(neighbors [][]NeighborLog) gocv.Mat {
	nRow, nCol := len(neighbors), len(neighbors[0])
	bytes := make([]byte, nRow*nCol)
	b := 0
	max := 0.0

	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			if v := neighbors[i][j].dist; v > max {
				max = v
			}
		}
	}

	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			bytes[b] = byte(neighbors[i][j].dist * 255 / max)
			b++
		}
	}

	mat, _ := gocv.NewMatFromBytes(nRow, nCol, gocv.MatChannels1, bytes)

	return mat
}
