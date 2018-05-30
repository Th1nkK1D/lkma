package main

import (
	"gocv.io/x/gocv"
	"gonum.org/v1/gonum/mat"
)

// GetBlankDatas -
func GetBlankFloats(r, c, chs int) [][]float64 {
	floats := make([][]float64, chs)

	for f := range floats {
		floats[f] = make([]float64, r*c)
	}

	return floats
}

// NewColorMat -
func NewColorMat(r, c, chs int, datas [][]float64) ColorMat {
	mats := make(ColorMat, chs)

	for m := range mats {
		mats[m] = mat.NewDense(r, c, datas[m])
	}

	return mats
}

// GetNumMat - get image gonum matrix
func GetNumMat(img gocv.Mat) ColorMat {
	bytes := img.ToBytes()
	nPixel := img.Cols() * img.Rows()
	chs := img.Channels()

	floats := make([][]float64, chs)

	for f := range floats {
		floats[f] = make([]float64, nPixel)
	}

	for i, b := range bytes {
		floats[i%chs][i/chs] = float64(b)
	}

	return NewColorMat(img.Cols(), img.Rows(), chs, floats)
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
				i++
			}
		}
	}

	newMat, _ := gocv.NewMatFromBytes(nR, nC, matType, bytes)

	return newMat
}

// CloneColorMatPixel -
func CloneColorMatPixel(dst, src ColorMat, i, j int) {
	chs := len(src)

	for c := 0; c < chs; c++ {
		dst[c].Set(i, j, src[c].At(i, j))
	}
}
