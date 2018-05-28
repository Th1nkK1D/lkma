package main

import (
	"gocv.io/x/gocv"
)

const imgPath = "dog_sm.jpg"
const scbPath = "dog_sm_scrb.jpg"

func main() {
	img := gocv.IMRead(imgPath, gocv.IMReadGrayScale)

	fd := GetFirstDerivative(img)

	gocv.IMWrite("fd-out.jpg", GetCVMat(fd))
}
