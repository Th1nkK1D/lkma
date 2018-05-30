package main

import "gonum.org/v1/gonum/mat"

const errTh = 10

// ExtractScribble -
func ExtractScribble(imgMats, scrb ColorMat) (ColorMat, ColorMat, ColorMat, *mat.Dense) {
	nR, nC := scrb[0].Dims()
	chs := len(imgMats)

	FG, BG := NewColorMat(nR, nC, chs, GetBlankFloats(nR, nC, chs)), NewColorMat(nR, nC, chs, GetBlankFloats(nR, nC, chs))
	Alp := NewColorMat(nR, nC, 1, GetBlankFloats(nR, nC, chs))
	ScrbMask := mat.NewDense(nR, nC, make([]float64, nR*nC))

	ScrbMask.Apply(func(i, j int, v float64) float64 {
		if v < errTh {
			// FG
			CloneColorMatPixel(FG, imgMats, i, j)
			Alp[0].Set(i, j, 255)

			return 1
		}

		if v > 255-errTh {
			// BG
			CloneColorMatPixel(BG, imgMats, i, j)

			return -1
		}

		// Unknown
		return 0
	}, scrb[0])

	return FG, BG, Alp, ScrbMask

}
