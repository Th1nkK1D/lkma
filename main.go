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
	// scrb := gocv.IMRead(scrbPath, gocv.IMReadGrayScale)

	mats := GetNumMat(img)
	fmt.Println(len(mats))

	for i := range mats {
		fmt.Println(mat.Sum(mats[i]))
	}

	gocv.IMWrite("test-out.jpg", GetCVMat(mats))

}
