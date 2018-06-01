package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

const imgPath = "cat_sm.jpg"
const scrbPath = "cat_sm_scrb.jpg"

func main() {
	fmt.Println("Reading input...")
	img := gocv.IMRead(imgPath, gocv.IMReadColor)
	scrb := gocv.IMRead(scrbPath, gocv.IMReadGrayScale)

	I := GetNumMat(img)

	scrbMats := GetNumMat(scrb)

	fmt.Println("Extracting scribble...")

	FG, BG, A, S := ExtractScribble(I, scrbMats)

	gocv.IMWrite("out-alp.jpg", GetCVMat(A, gocv.MatChannels1))
	gocv.IMWrite("out-fg.jpg", GetCVMat(FG, gocv.MatChannels3))
	gocv.IMWrite("out-bg.jpg", GetCVMat(BG, gocv.MatChannels3))

	fmt.Println("Exploring neighbor...")

	nFG, nBG := ExploreNeighbor(S)

	gocv.IMWrite("out-nfg.jpg", SaveNeighborLog(nFG))
	gocv.IMWrite("out-nbg.jpg", SaveNeighborLog(nBG))

	MimicNeighbor(I, FG, BG, S, nFG, nBG)

	gocv.IMWrite("out-fg-mimic.jpg", GetCVMat(FG, gocv.MatChannels3))
	gocv.IMWrite("out-bg-mimic.jpg", GetCVMat(BG, gocv.MatChannels3))

	// RunGradientDescent(I, FG, BG, A, S, nFG, nBG)

}
