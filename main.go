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

	imgMats := GetNumMat(img)

	scrbMats := GetNumMat(scrb)

	fmt.Println("Extracting scribble...")

	FG, BG, Alp, scrbMask := ExtractScribble(imgMats, scrbMats)

	gocv.IMWrite("out-fg.jpg", GetCVMat(FG, gocv.MatChannels3))
	gocv.IMWrite("out-bg.jpg", GetCVMat(BG, gocv.MatChannels3))
	gocv.IMWrite("out-alp.jpg", GetCVMat(Alp, gocv.MatChannels1))

	fmt.Println("Exploring neighbor...")

	nFG, nBG := ExploreNeighbor(scrbMask)

	gocv.IMWrite("out-nfg.jpg", SaveNeighborLog(nFG))
	gocv.IMWrite("out-nbg.jpg", SaveNeighborLog(nBG))

}
