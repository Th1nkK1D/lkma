package main

import "gonum.org/v1/gonum/mat"

const errTh = 10

// ExtractScribble -
func ExtractScribble(scrb *mat.Dense) (*mat.Dense, *mat.Dense) {
	nR, nC := scrb.Dims()
	FG, BG := mat.NewDense(nR, nC, make([]float64, nR*nC)), mat.NewDense(nR, nC, make([]float64, nR*nC))

	FG.Apply(func(i, j int, v float64) float64 {
		if v < errTh {
			return 1
		}

		return 0
	}, scrb)

	BG.Apply(func(i, j int, v float64) float64 {
		if v > 255-errTh {
			return 1
		}

		return 0
	}, scrb)

	return FG, BG

}
