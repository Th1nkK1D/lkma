package main

import (
	"fmt"

	"gocv.io/x/gocv"
	"gonum.org/v1/gonum/mat"
)

const imgPath = "cat.jpg"
const scrbPath = "cat_scrb.jpg"

func main() {
	img := gocv.IMRead(imgPath, gocv.IMReadColor)
	scrb := gocv.IMRead(scrbPath, gocv.IMReadGrayScale)

	imgMats := GetNumMat(img)
	fmt.Println(len(imgMats))

	for i := range imgMats {
		fmt.Println(mat.Sum(imgMats[i]))
	}

	scrbMats := GetNumMat(scrb)
	fmt.Println(len(scrbMats))

	for i := range scrbMats {
		fmt.Println(mat.Sum(scrbMats[i]))
	}

	gocv.IMWrite("scrb-out.jpg", GetCVMat(scrbMats, gocv.MatChannels1))
	gocv.IMWrite("test-out.jpg", GetCVMat(imgMats, gocv.MatChannels3))

}
