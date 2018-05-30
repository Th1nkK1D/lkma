package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

const imgPath = "cat.jpg"
const scrbPath = "cat_scrb.jpg"

func main() {
	fmt.Println("Reading input...")
	img := gocv.IMRead(imgPath, gocv.IMReadColor)
	scrb := gocv.IMRead(scrbPath, gocv.IMReadGrayScale)

	imgMats := GetNumMat(img)
	// fmt.Println(len(imgMats))

	// for i := range imgMats {
	// 	fmt.Println(mat.Sum(imgMats[i]))
	// }

	scrbMats := GetNumMat(scrb)
	// fmt.Println(len(scrbMats))

	// for i := range scrbMats {
	// 	fmt.Println(mat.Sum(scrbMats[i]))
	// }

	fmt.Println("Extracting scribble...")

	FG, BG, Alp, scrbMask := ExtractScribble(imgMats, scrbMats)

	gocv.IMWrite("out-fg.jpg", GetCVMat(FG, gocv.MatChannels3))
	gocv.IMWrite("out-bg.jpg", GetCVMat(BG, gocv.MatChannels3))
	gocv.IMWrite("out-alp.jpg", GetCVMat(Alp, gocv.MatChannels1))

	fmt.Println("Exploring neighbor...")

	nFG, nBG := ExploreNeighbor(scrbMask)

	gocv.IMWrite("out-nfg.jpg", SaveNeighborLog(nFG))
	gocv.IMWrite("out-nbg.jpg", SaveNeighborLog(nBG))

	// fmt.Println(ScrbMask)

}
