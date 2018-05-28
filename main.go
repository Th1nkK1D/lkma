package main

import (
	"gocv.io/x/gocv"
)

const imgPath = "dog.jpg"
const scrbPath = "dog_scrb.jpg"

func main() {
	img := gocv.IMRead(imgPath, gocv.IMReadGrayScale)
	scrb := gocv.IMRead(scrbPath, gocv.IMReadGrayScale)

	gocv.IMWrite("gs-out.jpg", img)

	// I := GetNumMat(img)

	fdx, fdy := GetFirstDerivative(img)

	gocv.IMWrite("fdx-out.jpg", GetCVMat(fdx))
	gocv.IMWrite("fdy-out.jpg", GetCVMat(fdy))

	FG, BG := ExtractScribble(GetNumMat(scrb))

	gocv.IMWrite("scrb-fg-out.jpg", GetCVMat(FG))
	gocv.IMWrite("scrb-bg-out.jpg", GetCVMat(BG))
}
