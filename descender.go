package main

import (
	"fmt"
	"math"
	"strconv"

	"gocv.io/x/gocv"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/plotter"
)

type stateRes struct {
	value  float64
	e      float64
	eTotal float64
	prop   float64
}

const pdWeight = 1
const cdWeight = 1

const eThreshold = 0.00001
const captureEach = 5

// const tStartRatio = 0.25
// const tDecline = 0.95

const step = 0.001

// Calculate E at a specific point
func getEnergyAt(i, j int, I, FG, BG, A ColorMat, S *mat.Dense, nFG, nBG [][]NeighborLog) float64 {
	if S.At(i, j) != 0 {
		return 0
	}

	// NN Pixel distance
	e := pdWeight * math.Pow(A[0].At(i, j)/256-nFG[i][j].dist/(nFG[i][j].dist+nBG[i][j].dist), 2)

	// NN Color space distance
	fgd, bgd := GetColorDistance(I, i, j, nFG[i][j].i, nFG[i][j].j), GetColorDistance(I, i, j, nBG[i][j].i, nBG[i][j].j)

	e += cdWeight * math.Pow(A[0].At(i, j)/256-fgd/(fgd+bgd), 2)

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

			// if ce > eMax {
			// 	eMax = ce
			// }

			e += ce
		}
	}

	// tStart := eMax * tStartRatio

	return E, e
}

func updateValue(I, FG, BG, A ColorMat, S, E *mat.Dense, nFG, nBG [][]NeighborLog) (ColorMat, ColorMat, ColorMat, *mat.Dense, float64) {
	nRow, nCol := I[0].Dims()
	chs := len(I)
	eps := 0.00001 //math.SmallestNonzeroFloat64

	newFG, newBG := NewColorMat(nRow, nCol, chs, GetBlankFloats(nRow, nCol, chs)), NewColorMat(nRow, nCol, chs, GetBlankFloats(nRow, nCol, chs))
	newA := NewColorMat(nRow, nCol, 1, GetBlankFloats(nRow, nCol, 1))
	newE := mat.NewDense(nRow, nCol, make([]float64, nRow*nCol))

	// Numerical derivative
	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {

			if S.At(i, j) == 0 {
				// FG update
				for ch := 0; ch < chs; ch++ {
					x := FG[ch].At(i, j)

					FG[ch].Set(i, j, x+eps)

					fpx := (getEnergyAt(i, j, I, FG, BG, A, S, nFG, nBG) - E.At(i, j)) / eps

					newFG[ch].Set(i, j, x-step*fpx)
					FG[ch].Set(i, j, x)
				}

				// BG update
				for ch := 0; ch < chs; ch++ {
					x := BG[ch].At(i, j)

					BG[ch].Set(i, j, x+eps)

					fpx := (getEnergyAt(i, j, I, FG, BG, A, S, nFG, nBG) - E.At(i, j)) / eps

					newBG[ch].Set(i, j, x-step*fpx)
					BG[ch].Set(i, j, x)
				}

				// Alpha update
				x := A[0].At(i, j)

				A[0].Set(i, j, x+eps)

				fpx := (getEnergyAt(i, j, I, FG, BG, A, S, nFG, nBG) - E.At(i, j)) / eps

				newA[0].Set(i, j, x-step*fpx)
				A[0].Set(i, j, x)

				// fmt.Printf("A(%v,%v), %v -> %v\n", i, j, A[0].At(i, j), newA[0].At(i, j))

				// Update energy
				newE.Set(i, j, getEnergyAt(i, j, I, FG, BG, A, S, nFG, nBG))
			}
		}
	}

	return newFG, newBG, newA, newE, mat.Sum(newE)
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
			gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-fg.jpg", GetCVMat(FG, gocv.MatChannels3))
			gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-bg.jpg", GetCVMat(BG, gocv.MatChannels3))
			gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-a.jpg", GetCVMat(A, gocv.MatChannels1))
		}

		// Update Value
		l2e = le
		le = e

		FG, BG, A, E, e = updateValue(I, FG, BG, A, S, E, nFG, nBG)

		de = (math.Abs(e-le) + math.Abs(le-l2e)) / 2

	}

	fmt.Printf("%v (Final): E = %v, dE_avg = %v (%v)\n", i, e, de, eThreshold)

	// Save for graph plotting
	p := make(plotter.XYs, 1)

	p[0].X = float64(i)
	p[0].Y = e

	points = append(points, p[0])

	// Write output
	gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-final-fg.jpg", GetCVMat(FG, gocv.MatChannels3))
	gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-final-bg.jpg", GetCVMat(BG, gocv.MatChannels3))
	gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-final-a.jpg", GetCVMat(A, gocv.MatChannels3))

	// // Plot graph
	// plots, err := plot.New()
	// if err != nil {
	// 	panic(err)
	// }

	// plots.Title.Text = "Gibb's Sampling Energy"
	// plots.X.Label.Text = "Iteration"
	// plots.Y.Label.Text = "Energy"

	// err = plotutil.AddLines(plots, "", points)
	// if err != nil {
	// 	panic(err)
	// }

	// if err := plots.Save(8*vg.Inch, 5*vg.Inch, "energy.png"); err != nil {
	// 	panic(err)
	// }
}
