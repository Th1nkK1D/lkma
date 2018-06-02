package main

import (
	"math"

	"gocv.io/x/gocv"

	"gonum.org/v1/gonum/mat"
)

func getDist(i, j, x, y int) float64 {
	return math.Sqrt(math.Pow(float64(i-x), 2) + math.Pow(float64(j-y), 2))
}

func explore(neighbours [][]NeighbourLog, mask *mat.Dense, i, j, si, sj, w, h int, start bool) [][]NeighbourLog {
	dist := getDist(i, j, si, sj)

	if start || (mask.At(i, j) == 0 && (neighbours[i][j].dist == 0 || dist < neighbours[i][j].dist)) {
		// Update value
		if mask.At(i, j) == 0 {
			neighbours[i][j].i = si
			neighbours[i][j].j = sj
			neighbours[i][j].dist = dist
		}

		// Water flow
		if i+1 < h {
			neighbours = explore(neighbours, mask, i+1, j, si, sj, w, h, false)
		}
		if j+1 < w {
			neighbours = explore(neighbours, mask, i, j+1, si, sj, w, h, false)
		}
		if i-1 >= 0 {
			neighbours = explore(neighbours, mask, i-1, j, si, sj, w, h, false)
		}
		if j-1 >= 0 {
			neighbours = explore(neighbours, mask, i, j-1, si, sj, w, h, false)
		}
	}

	return neighbours
}

// ExploreNeighbour -
func ExploreNeighbour(mask *mat.Dense) ([][]NeighbourLog, [][]NeighbourLog) {
	nRow, nCol := mask.Dims()

	// Init neighbour logs
	neighbourFG, neighbourBG := make([][]NeighbourLog, nRow), make([][]NeighbourLog, nRow)

	for n := range neighbourFG {
		neighbourFG[n] = make([]NeighbourLog, nCol)
		neighbourBG[n] = make([]NeighbourLog, nCol)
	}

	// Start
	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			if mask.At(i, j) == 1 {
				// Marked FG
				neighbourFG = explore(neighbourFG, mask, i, j, i, j, nCol, nRow, true)
			} else if mask.At(i, j) == -1 {
				// Marked BG
				neighbourBG = explore(neighbourBG, mask, i, j, i, j, nCol, nRow, true)
			}
		}
	}

	return neighbourFG, neighbourBG
}

// MimicNeighbour -
func MimicNeighbour(I, FG, BG ColorMat, mask *mat.Dense, neighbourFG, neighbourBG [][]NeighbourLog) {
	nRow, nCol := mask.Dims()

	// Mimic
	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			if mask.At(i, j) == 0 {
				CloneColorMatPixel(FG, i, j, I, neighbourFG[i][j].i, neighbourFG[i][j].j)
				CloneColorMatPixel(BG, i, j, I, neighbourBG[i][j].i, neighbourBG[i][j].j)
			}
		}
	}
}

// SaveNeighbourLog -
func SaveNeighbourLog(neighbours [][]NeighbourLog) gocv.Mat {
	nRow, nCol := len(neighbours), len(neighbours[0])
	bytes := make([]byte, nRow*nCol)
	b := 0
	max := 0.0

	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			if v := neighbours[i][j].dist; v > max {
				max = v
			}
		}
	}

	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			bytes[b] = byte(neighbours[i][j].dist * 255 / max)
			b++
		}
	}

	mat, _ := gocv.NewMatFromBytes(nRow, nCol, gocv.MatChannels1, bytes)

	return mat
}
