package main

import (
	"gocv.io/x/gocv"
	"gonum.org/v1/gonum/mat"
)

// GetNumMat - get image gonum matrix
func GetNumMat(img gocv.Mat) ColorMat {
	bytes := img.ToBytes()
	nPixel := img.Cols() * img.Rows()
	chs := img.Channels()

	floats := make([][]float64, chs)
	mats := make(ColorMat, chs)

	for f := range floats {
		floats[f] = make([]float64, nPixel)
	}

	for i, b := range bytes {
		floats[i%chs][i/chs] = float64(b)
	}

	for m := range mats {
		mats[m] = mat.NewDense(img.Rows(), img.Cols(), floats[m])
	}

	return mats
}

// GetCVMat -
func GetCVMat(mats ColorMat, matType gocv.MatType) gocv.Mat {
	nR, nC := mats[0].Dims()
	bytes := make([]byte, nR*nC*len(mats))

	i := 0

	for r := 0; r < nR; r++ {
		for c := 0; c < nC; c++ {
			for m := range mats {
				bytes[i] = byte(mats[m].At(r, c))
				// fmt.Printf("%v -> %v\n", imgMat.At(r, c), bytes[i])
				i++
			}
		}
	}

	newMat, _ := gocv.NewMatFromBytes(nR, nC, matType, bytes)

	return newMat
}
