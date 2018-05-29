package main

import (
	"gocv.io/x/gocv"
	"gonum.org/v1/gonum/mat"
)

// GetNumMat - get image gonum matrix
func GetNumMat(img gocv.Mat) *mat.Dense {
	bytes := img.ToBytes()

	floats := make([]float64, len(bytes))

	for i, b := range bytes {
		floats[i] = float64(b)
	}

	return mat.NewDense(img.Rows(), img.Cols(), floats)
}

// GetCVMat -
func GetCVMat(imgMat *mat.Dense) gocv.Mat {
	nR, nC := imgMat.Dims()
	bytes := make([]byte, nR*nC)

	max := mat.Max(imgMat)

	i := 0

	for r := 0; r < nR; r++ {
		for c := 0; c < nC; c++ {
			bytes[i] = byte(imgMat.At(r, c) * 255 / max)
			// fmt.Printf("%v -> %v\n", imgMat.At(r, c), bytes[i])
			i++
		}
	}

	return gocv.NewMatFromBytes(nR, nC, gocv.MatTypeCV8U, bytes)
}
