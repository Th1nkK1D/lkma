package main

import (
	"math"
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

type stateRes struct {
	value  float64
	e      float64
	eTotal float64
	prop   float64
}

const pdWeight = 1
const cdWeight = 1

// const eThreshold = 20
// const captureEach = 20

const tStartRatio = 0.25

// const tDecline = 0.95

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
func GetInitEnergy(I, FG, BG, A ColorMat, S *mat.Dense, nFG, nBG [][]NeighborLog) (*mat.Dense, float64, float64) {
	nRow, nCol := I[0].Dims()

	E := mat.NewDense(nRow, nCol, make([]float64, nRow*nCol))

	e, eMax := 0.0, 0.0

	// Calculate E
	for i := nRow - 1; i >= 0; i-- {
		for j := nCol - 1; j >= 0; j-- {
			ce := getEnergyAt(i, j, I, FG, BG, A, S, nFG, nBG)
			E.Set(i, j, ce)

			if ce > eMax {
				eMax = ce
			}

			e += ce
		}
	}

	tStart := eMax * tStartRatio

	return E, e, tStart
}

// Calculate Probability
func calProp(states []stateRes, eTotalMin, t float64) ([]stateRes, []stateRes, float64) {
	eSmallest := -t * math.Log(math.MaxFloat64)
	pSum := 0.0

	updatedStates := make([]stateRes, 0)

	for s := range states {
		states[s].eTotal -= eTotalMin

		// Overflow handle
		if states[s].e < eSmallest {
			states[s].prop = 0
		} else {
			states[s].prop = math.Pow(math.E, -states[s].eTotal/t)

		}

		if states[s].prop > 0 {
			updatedStates = append(updatedStates, states[s])
			pSum += states[s].prop
		}
	}

	return states, updatedStates, pSum
}

// Sampling from PDF
func samplingPDF(states, updatedStates []stateRes, pSum float64) stateRes {
	randVal := rand.Float64() * pSum
	nState := len(updatedStates)
	i := 0

	if len(updatedStates) == 0 {
		return states[rand.Intn(len(states)-1)]
	}

	for stack := 0.0; i < nState && stack < randVal; i++ {
		stack += updatedStates[i].prop
	}

	if i >= nState {
		return updatedStates[nState-1]
	}

	return updatedStates[i]
}

// // Get F updated value with new E and e
// func getUpdatedValue(Il, Ir, D, lE fmat, le, t float64) (fmat, fmat, float64) {
// 	nRow, nCol := len(Il), len(Il[0])
// 	dispMax := math.Floor(maxDispRatio * float64(nCol))
// 	nD, nE := make(fmat, nRow), make(fmat, nRow)
// 	ne := 0.0

// 	// Each point
// 	for r := 0; r < nRow; r++ {
// 		nD[r], nE[r] = make([]float64, nCol), make([]float64, nCol)

// 		for c := 0; c < nCol; c++ {
// 			eTotalMin := math.MaxFloat64
// 			ls := D[r][c]

// 			states := make([]stateRes, int(dispMax)+1)
// 			i := 0

// 			// State space energy calculating
// 			for s := 0.0; s <= dispMax; s += 1.0 {
// 				states[i].value = s
// 				D[r][c] = s

// 				states[i].e = getEnergyAt(Il, Ir, D, r, c, dispMax)

// 				states[i].eTotal = le + states[i].e - lE[r][c]

// 				if r > 0 {
// 					states[i].eTotal += getEnergyAt(Il, Ir, D, r-1, c, dispMax) - lE[r-1][c]
// 				}

// 				if c > 0 {
// 					states[i].eTotal += getEnergyAt(Il, Ir, D, r, c-1, dispMax) - lE[r][c-1]
// 				}

// 				// Check for min
// 				if states[i].eTotal < eTotalMin {
// 					eTotalMin = states[i].eTotal
// 				}
// 				i++
// 			}

// 			D[r][c] = ls

// 			// Update new Value
// 			sTarget := samplingPDF(calProp(states, eTotalMin, t))

// 			nD[r][c] = sTarget.value
// 			nE[r][c] = sTarget.e
// 			ne += sTarget.e
// 		}
// 	}

// 	return nD, nE, ne
// }

// // Main function
// func main() {
// 	// Init value
// 	Il, Ir, D, E, e, t := initMatrix(gocv.IMRead(os.Args[1], gocv.IMReadGrayScale), gocv.IMRead(os.Args[2], gocv.IMReadGrayScale))

// 	le, l2e := 0.0, 0.0
// 	de := e
// 	i := 0

// 	// init graph data holder
// 	points := make(plotter.XYs, 0)

// 	// Gradient descent looping
// 	for ; de > eThreshold; i++ {
// 		fmt.Printf("%v: E = %v, dE_avg = %v (%v), t = %v\n", i, e, de, eThreshold, t)

// 		// Save for graph plotting
// 		p := make(plotter.XYs, 1)

// 		p[0].X = float64(i)
// 		p[0].Y = e

// 		points = append(points, p[0])

// 		// Print preview
// 		if i%captureEach == 0 {
// 			gocv.IMWrite("out-"+strconv.Itoa(i)+".jpg", getMatfromFmat(D))
// 		}

// 		// Update Value
// 		l2e = le
// 		le = e

// 		D, E, e = getUpdatedValue(Il, Ir, D, E, le, t)

// 		de = (math.Abs(e-le) + math.Abs(le-l2e)) / 2

// 		t *= tDecline
// 	}

// 	fmt.Printf("%v (Final): E = %v, dE_avg = %v (%v), t = %v\n", i, e, de, eThreshold, t)

// 	// Save for graph plotting
// 	p := make(plotter.XYs, 1)

// 	p[0].X = float64(i)
// 	p[0].Y = e

// 	points = append(points, p[0])

// 	// Write output
// 	gocv.IMWrite("out-final.jpg", getMatfromFmat(D))

// 	// Plot graph
// 	plots, err := plot.New()
// 	if err != nil {
// 		panic(err)
// 	}

// 	plots.Title.Text = "Gibb's Sampling Energy"
// 	plots.X.Label.Text = "Iteration"
// 	plots.Y.Label.Text = "Energy"

// 	err = plotutil.AddLines(plots, "", points)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err := plots.Save(8*vg.Inch, 5*vg.Inch, "energy.png"); err != nil {
// 		panic(err)
// 	}
// }
