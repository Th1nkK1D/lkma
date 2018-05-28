package main

import (
	"image"
	"math"

	"gocv.io/x/gocv"
	"gonum.org/v1/gonum/mat"
)

const blurSize = 5.0

var fdxMetrix = mat.NewDense(5, 5, []float64{0, 0, 0, 0, 0, 0, -2, -1, 1, 2, 0, -4, -2, 2, 4, 0, -2, -1, 1, 2, 0, 0, 0, 0, 0})
var fdyMetrix = mat.NewDense(5, 5, []float64{0, 0, 0, 0, 0, 0, -2, -4, -2, 0, 0, -1, -2, -1, 0, 0, 1, 2, 1, 0, 0, 2, 4, 2, 0})

// Get valid image index
func getValidIndex(i, bound int) int {
	if i < 0 {
		return 0
	}
	if i >= bound {
		return bound - 1
	}

	return i
}

// Convolution
func convolute(imgMat, filter *mat.Dense) *mat.Dense {
	nRow, nCol := imgMat.Dims()
	fR, fC := filter.Dims()
	filterRad := fR / 2

	table := mat.NewDense(nRow, nCol, make([]float64, nRow*nCol))

	// Loop through all pixel
	for y := 0; y < nRow; y++ {
		for x := 0; x < nCol; x++ {
			sum := 0.0

			// Loop in filter
			for j := 0; j < fR; j++ {
				for i := 0; i < fC; i++ {
					sum += filter.At(j, i) * imgMat.At(getValidIndex(y-j+filterRad, nRow), getValidIndex(x-i+filterRad, nCol))
				}
			}

			table.Set(y, x, math.Abs(sum))
		}
	}

	return table
}

// GetFirstDerivative -
func GetFirstDerivative(img gocv.Mat) (*mat.Dense, *mat.Dense) {
	gocv.GaussianBlur(img, img, image.Point{blurSize, blurSize}, math.Floor(blurSize/2.0), math.Floor(blurSize/2.0), gocv.BorderReflect)

	imgMat := GetNumMat(img)

	return convolute(imgMat, fdxMetrix), convolute(imgMat, fdyMetrix)

}
