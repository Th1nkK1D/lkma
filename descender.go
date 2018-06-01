package main

import (
	"fmt"
	"math"
	"strconv"

	"gocv.io/x/gocv"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const saWeight = 1
const pdWeight = 1
const cdWeight = 1
const ieWeight = 2

const eThreshold = 100
const captureEach = 100

// const tStartRatio = 0.25
// const tDecline = 0.95

const step = 0.001

// Calculate E at a specific point
func getEnergyAt(i, j int, I, FG, BG, A ColorMat, S *mat.Dense, nFG, nBG [][]NeighborLog) float64 {
	nRow, nCol := S.Dims()
	e := 0.0

	// Smoothness constrain
	sCount, saSum := 0.0, 0.0

	if i < nRow-1 {
		// sgSum += GetColorDistance(FG, i, j, i+1, j)/2 + GetColorDistance(BG, i, j, i+1, j)/2
		saSum += math.Pow(GetColorDistance(I, i, j, i+1, j)-GetColorDistance(A, i, j, i+1, j), 2)
		sCount++
	}

	if j < nCol-1 {
		// sgSum += GetColorDistance(FG, i, j, i, j+1)/2 + GetColorDistance(BG, i, j, i, j+1)/2
		saSum += math.Pow(GetColorDistance(I, i, j, i, j+1)-GetColorDistance(A, i, j, i, j+1), 2)
	}

	if sCount > 0 {
		e += (saWeight * saSum) / sCount
	}

	if S.At(i, j) != 0 {
		return e
	}

	// NN Pixel distance
	e += pdWeight * math.Pow(A[0].At(i, j)/255-nBG[i][j].dist/(nFG[i][j].dist+nBG[i][j].dist), 2)

	// NN Color space distance
	fgd, bgd := GetColorDistance(I, i, j, nFG[i][j].i, nFG[i][j].j), GetColorDistance(I, i, j, nBG[i][j].i, nBG[i][j].j)

	e += cdWeight * math.Pow(A[0].At(i, j)/255-bgd/(fgd+bgd), 2)

	// Image error
	chs := len(I)
	ie := 0.0
	a := A[0].At(i, j) / 255

	for ch := 0; ch < chs; ch++ {
		ie += math.Pow(I[ch].At(i, j)-(a*FG[ch].At(i, j)-(1-a)*BG[ch].At(i, j)), 2)
	}

	e += ieWeight * ie / float64(255*chs)

	return e
}

// GetInitEnergy - Initialize energy matrix
func getInitEnergy(I, FG, BG, A ColorMat, S *mat.Dense, nFG, nBG [][]NeighborLog) (*mat.Dense, float64) {
	nRow, nCol := I[0].Dims()

	E := mat.NewDense(nRow, nCol, make([]float64, nRow*nCol))

	e := 0.0

	// Calculate E
	for i := nRow - 1; i >= 0; i-- {
		for j := nCol - 1; j >= 0; j-- {
			ce := getEnergyAt(i, j, I, FG, BG, A, S, nFG, nBG)
			E.Set(i, j, ce)

			e += ce
		}
	}

	return E, e
}

func updateValue(I, FG, BG, A ColorMat, S, E *mat.Dense, nFG, nBG [][]NeighborLog) (ColorMat, *mat.Dense, float64) {
	nRow, nCol := I[0].Dims()
	// chs := len(I)
	eps := math.Sqrt(2.2 * math.Pow(10, -16))

	// newFG, newBG := NewColorMat(nRow, nCol, chs, GetBlankFloats(nRow, nCol, chs)), NewColorMat(nRow, nCol, chs, GetBlankFloats(nRow, nCol, chs))
	newA := NewColorMat(nRow, nCol, 1, GetBlankFloats(nRow, nCol, 1))
	newE := mat.NewDense(nRow, nCol, make([]float64, nRow*nCol))

	// Numerical derivative
	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {

			if S.At(i, j) == 0 {
				// // FG update
				// for ch := 0; ch < chs; ch++ {
				// 	x := FG[ch].At(i, j)

				// 	FG[ch].Set(i, j, x+eps)

				// 	fpx := (getEnergyAt(i, j, I, FG, BG, A, S, nFG, nBG) - E.At(i, j))

				// 	if i > 0 {
				// 		fpx += getEnergyAt(i-1, j, I, FG, BG, A, S, nFG, nBG) - E.At(i-1, j)
				// 	}

				// 	if j > 0 {
				// 		fpx += getEnergyAt(i, j-1, I, FG, BG, A, S, nFG, nBG) - E.At(i, j-1)
				// 	}

				// 	fpx /= eps

				// 	newFG[ch].Set(i, j, x-step*fpx)
				// 	FG[ch].Set(i, j, x)
				// }

				// // BG update
				// for ch := 0; ch < chs; ch++ {
				// 	x := BG[ch].At(i, j)

				// 	BG[ch].Set(i, j, x+eps)

				// 	fpx := (getEnergyAt(i, j, I, FG, BG, A, S, nFG, nBG) - E.At(i, j))

				// 	if i > 0 {
				// 		fpx += getEnergyAt(i-1, j, I, FG, BG, A, S, nFG, nBG) - E.At(i-1, j)
				// 	}

				// 	if j > 0 {
				// 		fpx += getEnergyAt(i, j-1, I, FG, BG, A, S, nFG, nBG) - E.At(i, j-1)
				// 	}

				// 	fpx /= eps

				// 	newBG[ch].Set(i, j, x-step*fpx)
				// 	BG[ch].Set(i, j, x)
				// }

				// Alpha update
				x := A[0].At(i, j)

				A[0].Set(i, j, x+eps)

				fpx := (getEnergyAt(i, j, I, FG, BG, A, S, nFG, nBG) - E.At(i, j))

				if i > 0 {
					fpx += getEnergyAt(i-1, j, I, FG, BG, A, S, nFG, nBG) - E.At(i-1, j)
				}

				if j > 0 {
					fpx += getEnergyAt(i, j-1, I, FG, BG, A, S, nFG, nBG) - E.At(i, j-1)
				}

				fpx /= eps

				newA[0].Set(i, j, x-step*fpx)
				A[0].Set(i, j, x)

				// fmt.Printf("A(%v,%v), %v -> %v\n", i, j, A[0].At(i, j), newA[0].At(i, j))
			} else {
				// Copy value
				// CloneColorMatPixel(newFG, FG, i, j)
				// CloneColorMatPixel(newBG, BG, i, j)
				CloneColorMatPixel(newA, i, j, A, i, j)
			}

		}
	}

	// Update energy
	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			newE.Set(i, j, getEnergyAt(i, j, I, FG, BG, newA, S, nFG, nBG))
		}
	}

	return newA, newE, mat.Sum(newE)
}

// RunGradientDescent -
func RunGradientDescent(I, FG, BG, A ColorMat, S *mat.Dense, nFG, nBG [][]NeighborLog) {
	E, e := getInitEnergy(I, FG, BG, A, S, nFG, nBG)

	le, l2e := 0.0, 0.0
	de := e
	i := 0

	// init graph data holder
	points := make(plotter.XYs, 0)

	// Gradient descent looping
	for ; de > eThreshold; i++ {
		fmt.Printf("%v: E = %v, dE_avg = %v (%v)\n", i, e, de, eThreshold)

		// Save for graph plotting
		p := make(plotter.XYs, 1)

		p[0].X = float64(i)
		p[0].Y = e

		points = append(points, p[0])

		// Print preview
		if i%captureEach == 0 {
			// gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-fg.jpg", GetCVMat(FG, gocv.MatChannels3))
			// gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-bg.jpg", GetCVMat(BG, gocv.MatChannels3))
			gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-a.jpg", GetCVMat(A, gocv.MatChannels1))
		}

		// Update Value
		l2e = le
		le = e

		A, E, e = updateValue(I, FG, BG, A, S, E, nFG, nBG)

		de = (math.Abs(e-le) + math.Abs(le-l2e)) / 2

	}

	fmt.Printf("%v (Final): E = %v, dE_avg = %v (%v)\n", i, e, de, eThreshold)

	// Save for graph plotting
	p := make(plotter.XYs, 1)

	p[0].X = float64(i)
	p[0].Y = e

	points = append(points, p[0])

	// Write output
	// gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-final-fg.jpg", GetCVMat(FG, gocv.MatChannels3))
	// gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-final-bg.jpg", GetCVMat(BG, gocv.MatChannels3))
	gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-final-a.jpg", GetCVMat(A, gocv.MatChannels1))

	// Plot graph
	plots, err := plot.New()
	if err != nil {
		panic(err)
	}

	plots.Title.Text = "Gibb's Sampling Energy"
	plots.X.Label.Text = "Iteration"
	plots.Y.Label.Text = "Energy"

	err = plotutil.AddLines(plots, "", points)
	if err != nil {
		panic(err)
	}

	if err := plots.Save(8*vg.Inch, 5*vg.Inch, "out-energy.png"); err != nil {
		panic(err)
	}
}
