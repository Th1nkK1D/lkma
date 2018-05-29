package main

import (
	"gonum.org/v1/gonum/mat"
)

const hillTh = 50

func flow(ground, hillx, hilly *mat.Dense, i, j, w, h int, dist float64, start bool) {
	wet := ground.At(i, j)

	if start || wet < 1 {
		// // Dry zone
		// if hillx.At(i, j) > hillTh && hilly.At(i, j) > hillTh {
		// 	// Hill -> less wet
		dist *= 0.9
		// 	return
		// }

		if start || wet == 0 || dist > wet {
			// Update value
			ground.Set(i, j, dist)
			// fmt.Printf("Set %v,%v -> %v\n", i, j, dist)

			// Water flow
			if i+1 < h {
				flow(ground, hillx, hilly, i+1, j, w, h, dist, false)
			}
			if j+1 < w {
				flow(ground, hillx, hilly, i, j+1, w, h, dist, false)
			}
			if i-1 >= 0 {
				flow(ground, hillx, hilly, i-1, j, w, h, dist, false)
			}
			if j-1 >= 0 {
				flow(ground, hillx, hilly, i, j-1, w, h, dist, false)
			}
		}

	}
}

// StartRaining -
func StartRaining(ground, hillx, hilly *mat.Dense) {
	nRow, nCol := ground.Dims()

	// Start rain drop
	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			if ground.At(i, j) == 1 {
				flow(ground, hillx, hilly, i, j, nCol, nRow, 1, true)
			}
		}
	}

	// Normalize value
	// ground.Apply(func(i, j int, v float64) float64 {
	// 	if v > 0 {
	// 		return 1 / v
	// 	}
	// 	return 0
	// }, ground)
}
