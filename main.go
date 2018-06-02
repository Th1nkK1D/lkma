package main

import (
	"fmt"
	"os"

	"gocv.io/x/gocv"
)

func main() {
	fmt.Println("Reading input...")
	img := gocv.IMRead(os.Args[1], gocv.IMReadColor)
	scrb := gocv.IMRead(os.Args[2], gocv.IMReadGrayScale)

	I := GetNumMat(img)

	scrbMats := GetNumMat(scrb)

	fmt.Println("Extracting scribble...")

	FG, BG, A, S := ExtractScribble(I, scrbMats)

	gocv.IMWrite("out-alp.jpg", GetCVMat(A, gocv.MatChannels1))
	gocv.IMWrite("out-fg.jpg", GetCVMat(FG, gocv.MatChannels3))
	gocv.IMWrite("out-bg.jpg", GetCVMat(BG, gocv.MatChannels3))

	fmt.Println("Exploring neighbour...")

	nFG, nBG := ExploreNeighbour(S)

	gocv.IMWrite("out-nfg.jpg", SaveNeighbourLog(nFG))
	gocv.IMWrite("out-nbg.jpg", SaveNeighbourLog(nBG))

	fmt.Println("Mimicing neighbour...")

	MimicNeighbour(I, FG, BG, S, nFG, nBG)

	gocv.IMWrite("out-fg-mimic.jpg", GetCVMat(FG, gocv.MatChannels3))
	gocv.IMWrite("out-bg-mimic.jpg", GetCVMat(BG, gocv.MatChannels3))

	RunGradientDescent(I, FG, BG, A, S, nFG, nBG)

}
