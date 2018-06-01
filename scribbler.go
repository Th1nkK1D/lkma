package main

import (
	"math"
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

const errTh = 10

func randPixelVal() float64 {
	return math.Floor(rand.Float64() * 255)
}

// ExtractScribble -
func ExtractScribble(imgMats, scrb ColorMat) (ColorMat, ColorMat, ColorMat, *mat.Dense) {
	nR, nC := scrb[0].Dims()
	chs := len(imgMats)

	FG, BG := NewColorMat(nR, nC, chs, GetBlankFloats(nR, nC, chs)), NewColorMat(nR, nC, chs, GetBlankFloats(nR, nC, chs))
	Alp := NewColorMat(nR, nC, 1, GetBlankFloats(nR, nC, 1))
	ScrbMask := mat.NewDense(nR, nC, make([]float64, nR*nC))

	ScrbMask.Apply(func(i, j int, v float64) float64 {
		if v < errTh {
			// FG
			CloneColorMatPixel(FG, i, j, imgMats, i, j)
			Alp[0].Set(i, j, 255)

			return 1
		}

		if v > 255-errTh {
			// BG
			CloneColorMatPixel(BG, i, j, imgMats, i, j)

			return -1
		}

		// Unknown - random
		Alp[0].Set(i, j, randPixelVal())

		for c := 0; c < chs; c++ {
			FG[c].Set(i, j, randPixelVal())
			BG[c].Set(i, j, randPixelVal())
		}

		return 0
	}, scrb[0])

	return FG, BG, Alp, ScrbMask

}
