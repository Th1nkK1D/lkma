package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"gocv.io/x/gocv"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const sgWeight = 1
const saWeight = 1
const pdWeight = 1
const cdWeight = 1
const ieWeight = 1

const eThreshold = 1000
const captureEach = 5

const tStartRatio = 0.25
const tDecline = 0.9

// Calculate E at a specific point
func getEnergyAt(i, j int, I, FG, BG, A ColorMat, S *mat.Dense, nFG, nBG [][]NeighborLog) float64 {
	nRow, nCol := S.Dims()
	e := 0.0

	// Smoothness constrain
	sCount, sgSum, saSum := 0.0, 0.0, 0.0

	if i < nRow-1 {
		sgSum += GetColorDistance(FG, i, j, i+1, j)/2 + GetColorDistance(BG, i, j, i+1, j)/2
		saSum += math.Abs(GetColorDistance(I, i, j, i+1, j) - GetColorDistance(A, i, j, i+1, j))
		sCount++
	}

	if j < nCol-1 {
		sgSum += GetColorDistance(FG, i, j, i, j+1)/2 + GetColorDistance(BG, i, j, i, j+1)/2
		saSum += math.Abs(GetColorDistance(I, i, j, i, j+1) - GetColorDistance(A, i, j, i, j+1))
	}

	if sCount > 0 {
		e += (sgWeight*sgSum + saWeight*saSum) / sCount
	}

	if S.At(i, j) != 0 {
		return e
	}

	// NN Pixel distance
	e += pdWeight * math.Pow(A[0].At(i, j)/256-nBG[i][j].dist/(nFG[i][j].dist+nBG[i][j].dist), 2)

	// NN Color space distance
	fgd, bgd := GetColorDistance(I, i, j, nFG[i][j].i, nFG[i][j].j), GetColorDistance(I, i, j, nBG[i][j].i, nBG[i][j].j)

	e += cdWeight * math.Pow(A[0].At(i, j)/256-bgd/(fgd+bgd), 2)

	// Image error
	chs := len(I)
	ie := 0.0
	a := A[0].At(i, j)

	for ch := 0; ch < chs; ch++ {
		ie += math.Pow(I[ch].At(i, j)-(a*FG[ch].At(i, j)-(1-a)*BG[ch].At(i, j)), 2)
	}

	e += ieWeight * ie / float64(256*chs)

	return e
}

// GetInitEnergy - Initialize energy matrix
func getInitEnergy(I, FG, BG, A ColorMat, S *mat.Dense, nFG, nBG [][]NeighborLog) (*mat.Dense, float64, float64) {
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
func calProp(states []StateRes, eTotalMin, t float64) ([]StateRes, []StateRes, float64) {
	eSmallest := -t * math.Log(math.MaxFloat64)
	pSum := 0.0

	updatedStates := make([]StateRes, 0)

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
func samplingPDF(states, updatedStates []StateRes, pSum float64) StateRes {
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

func updateFG(i, j, ch int, I, FG, BG, A ColorMat, S, E *mat.Dense, nFG, nBG [][]NeighborLog, e, t float64) float64 {
	eTotalMin := math.MaxFloat64
	cv := FG[ch].At(i, j)

	states := make([]StateRes, 256)
	si := 0

	// State space energy calculating
	for s := 0.0; s < 256; s += 1.0 {
		states[si].value = s
		FG[ch].Set(i, j, s)

		states[si].e = getEnergyAt(i, j, I, FG, BG, A, S, nFG, nBG)

		states[si].eTotal = e + states[si].e - E.At(i, j)

		if i > 0 {
			states[si].eTotal += getEnergyAt(i-1, j, I, FG, BG, A, S, nFG, nBG) - E.At(i-1, j)
		}

		if j > 0 {
			states[si].eTotal += getEnergyAt(i, j-1, I, FG, BG, A, S, nFG, nBG) - E.At(i, j-1)
		}

		// Check for min
		if states[si].eTotal < eTotalMin {
			eTotalMin = states[si].eTotal
		}
		si++
	}

	FG[ch].Set(i, j, cv)
	// fmt.Printf("FG (%v,%v) = %v\n", i, j, FG[ch].At(i, j))

	// Update new Value
	sTarget := samplingPDF(calProp(states, eTotalMin, t))
	// fmt.Println(sTarget)

	return sTarget.value
}

func updateBG(i, j, ch int, I, FG, BG, A ColorMat, S, E *mat.Dense, nFG, nBG [][]NeighborLog, e, t float64) float64 {
	eTotalMin := math.MaxFloat64
	cv := BG[ch].At(i, j)

	states := make([]StateRes, 256)
	si := 0

	// State space energy calculating
	for s := 0.0; s < 256; s += 1.0 {
		states[si].value = s
		BG[ch].Set(i, j, s)

		states[si].e = getEnergyAt(i, j, I, FG, BG, A, S, nFG, nBG)

		states[si].eTotal = e + states[si].e - E.At(i, j)

		if i > 0 {
			states[si].eTotal += getEnergyAt(i-1, j, I, FG, BG, A, S, nFG, nBG) - E.At(i-1, j)
		}

		if j > 0 {
			states[si].eTotal += getEnergyAt(i, j-1, I, FG, BG, A, S, nFG, nBG) - E.At(i, j-1)
		}

		// Check for min
		if states[si].eTotal < eTotalMin {
			eTotalMin = states[si].eTotal
		}
		si++
	}

	BG[ch].Set(i, j, cv)

	// Update new Value
	sTarget := samplingPDF(calProp(states, eTotalMin, t))

	return sTarget.value
}

func updateA(i, j int, I, FG, BG, A ColorMat, S, E *mat.Dense, nFG, nBG [][]NeighborLog, e, t float64) float64 {
	eTotalMin := math.MaxFloat64
	cv := A[0].At(i, j)

	states := make([]StateRes, 256)
	si := 0

	// State space energy calculating
	for s := 0.0; s < 256; s += 1.0 {
		states[si].value = s
		A[0].Set(i, j, s)

		states[si].e = getEnergyAt(i, j, I, FG, BG, A, S, nFG, nBG)

		states[si].eTotal = e + states[si].e - E.At(i, j)

		if i > 0 {
			states[si].eTotal += getEnergyAt(i-1, j, I, FG, BG, A, S, nFG, nBG) - E.At(i-1, j)
		}

		if j > 0 {
			states[si].eTotal += getEnergyAt(i, j-1, I, FG, BG, A, S, nFG, nBG) - E.At(i, j-1)
		}

		// Check for min
		if states[si].eTotal < eTotalMin {
			eTotalMin = states[si].eTotal
		}
		si++
	}

	A[0].Set(i, j, cv)

	// Update new Value
	sTarget := samplingPDF(calProp(states, eTotalMin, t))

	return sTarget.value
}

func updateValue(I, FG, BG, A ColorMat, S, E *mat.Dense, nFG, nBG [][]NeighborLog, e, t float64) (ColorMat, ColorMat, ColorMat, *mat.Dense, float64) {
	nRow, nCol := I[0].Dims()
	chs := len(I)
	newFG, newBG := NewColorMat(nRow, nCol, chs, GetBlankFloats(nRow, nCol, chs)), NewColorMat(nRow, nCol, chs, GetBlankFloats(nRow, nCol, chs))
	newA := NewColorMat(nRow, nCol, 1, GetBlankFloats(nRow, nCol, 1))
	newE := mat.NewDense(nRow, nCol, make([]float64, nRow*nCol))
	newe := 0.0

	// Get update values
	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			if S.At(i, j) == 0 {
				for ch := 0; ch < chs; ch++ {
					newFG[ch].Set(i, j, updateFG(i, j, ch, I, FG, BG, A, S, E, nFG, nBG, e, t))
					newBG[ch].Set(i, j, updateBG(i, j, ch, I, FG, BG, A, S, E, nFG, nBG, e, t))
				}

				newA[0].Set(i, j, updateA(i, j, I, FG, BG, A, S, E, nFG, nBG, e, t))

			} else {
				CloneColorMatPixel(newFG, i, j, FG, i, j)
				CloneColorMatPixel(newBG, i, j, BG, i, j)
				CloneColorMatPixel(newA, i, j, A, i, j)
			}
		}
	}

	// Update E
	for i := 0; i < nRow; i++ {
		for j := 0; j < nCol; j++ {
			ce := getEnergyAt(i, j, I, newFG, BG, A, S, nFG, nBG)
			newE.Set(i, j, ce)
			newe += ce
		}
	}

	return newFG, newBG, newA, newE, newe
}

// RunGradientDescent -
func RunGradientDescent(I, FG, BG, A ColorMat, S *mat.Dense, nFG, nBG [][]NeighborLog) {
	E, e, t := getInitEnergy(I, FG, BG, A, S, nFG, nBG)

	le, l2e := 0.0, 0.0
	de := e
	i := 0

	// init graph data holder
	points := make(plotter.XYs, 0)

	// Gradient descent looping
	for ; de > eThreshold; i++ {
		fmt.Printf("%v: E = %v, dE_avg = %v (%v), t = %v\n", i, e, de, eThreshold, t)

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

		FG, BG, A, E, e = updateValue(I, FG, BG, A, S, E, nFG, nBG, e, t)

		de = (math.Abs(e-le) + math.Abs(le-l2e)) / 2

		t *= tDecline
	}

	fmt.Printf("%v (Final): E = %v, dE_avg = %v (%v), t = %v\n", i, e, de, eThreshold, t)

	// Save for graph plotting
	p := make(plotter.XYs, 1)

	p[0].X = float64(i)
	p[0].Y = e

	points = append(points, p[0])

	// Write output
	gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-final-fg.jpg", GetCVMat(FG, gocv.MatChannels3))
	gocv.IMWrite("out-gd-"+strconv.Itoa(i)+"-final-bg.jpg", GetCVMat(BG, gocv.MatChannels3))
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
