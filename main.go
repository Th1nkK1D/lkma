package main

import (
	"gocv.io/x/gocv"
)

const imgPath = "dog.jpg"
const scbPath = "dog_scrb.jpg"

func main() {
	img := gocv.IMRead(imgPath, gocv.IMReadGrayScale)

	gocv.IMWrite("gs-out.jpg", img)

	fdx, fdy := GetFirstDerivative(img)

	gocv.IMWrite("fdx-out.jpg", GetCVMat(fdx))
	gocv.IMWrite("fdy-out.jpg", GetCVMat(fdy))
}
