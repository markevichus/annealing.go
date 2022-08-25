package cooler

import (
	"fmt"
	svg "github.com/ajstarks/svgo"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

type Shaker interface {
	Shake() (energy float64, err error)
	Take()
	GetEnergy() float64
	StoreReport()
	GetResult() ShakeResult
}

type ShakeResult struct {
	Energy float64
	Tick   int
}

type Annealer interface {
	Run() (ShakeResult, error)
}

type AnnealingHist struct {
	t float64
	e float64
}

type AnnealingMachine struct {
	temp     float64
	stopTemp float64
	tick     int
	stopTick int
	shaker   Shaker
	hist     map[int]AnnealingHist
	rs       rand.Source
}

func NewAnnealingMachine(s Shaker, stopTemp float64, stopTick int) *AnnealingMachine {
	if stopTemp == 0 {
		stopTemp = 0.001
	}
	if stopTick == 0 {
		stopTick = 10000
	}
	am := AnnealingMachine{
		stopTemp: stopTemp,
		stopTick: stopTick,
		shaker:   s,
		rs:       rand.NewSource(time.Now().UnixNano()),
	}
	return &am
}

func (am *AnnealingMachine) Run() (result ShakeResult, err error) {
	am.reset()
	var minE float64 = 1
	var bestResult ShakeResult

	for {
		e, err := am.shaker.Shake()
		if err != nil {
			return ShakeResult{}, err
		}

		dE := e - minE
		if dE < 0 {
			minE = e
			//fmt.Println("minE", e, am.temp)
			am.transit()
			// Get Shaker report data
			//am.shaker.StoreReport()
			bestResult = am.shaker.GetResult()
			bestResult.Tick = am.tick
		} else if am.shouldITransit(dE) {
			//fmt.Println("random", e, am.temp, am.tick)
			am.transit()
		}
		am.decreaseTemperature()

		if am.tick == am.stopTick || am.temp <= am.stopTemp {
			fmt.Println("STOP", bestResult, am.tick, am.temp)
			break
		}
		am.tick += 1
	}
	//fmt.Println("from am", bestResult)
	//ch <- bestResult
	//am.storeReport()

	return bestResult, nil
}

func (am *AnnealingMachine) reset() {
	am.temp = 0.98
	am.tick = 1
	am.hist = make(map[int]AnnealingHist)
}

func (am *AnnealingMachine) transit() {
	am.shaker.Take()

	am.hist[am.tick] = AnnealingHist{
		t: am.temp,
		e: am.shaker.GetEnergy(),
	}
}

func (am *AnnealingMachine) decreaseTemperature() {
	am.temp *= 1.0 - 1.0/(float64(am.tick)*1.2+30.0)
}

func (am *AnnealingMachine) shouldITransit(dE float64) bool {
	probability := math.Pow(math.E, -dE/am.temp)
	if rand.New(am.rs).Float64() <= probability {
		return true
	} else {
		return false
	}
}

func (am *AnnealingMachine) storeReport() (err error) {
	var divider int = 10
	//labelShiftX, labelShiftY, labelHeight := 5, 5, 17
	var graphXMul, graphYMul float64 = 0.1, 400
	var graphWidth, graphHeight float64 = 20000, 3000
	scale := func(num float64) int {
		return int(num / float64(divider))
	}

	f, err := os.Create("annealing.svg")
	if err != nil {
		fmt.Errorf("error creating file. %v", err)
		return err
	}

	canvas := svg.New(f)
	canvas.Start(scale(graphWidth), scale(graphHeight*2))
	canvas.Rect(0, 0, scale(graphWidth), scale(graphHeight), "fill:none;stroke:green;stroke_width:2")

	ticks := make([]int, 0, len(am.hist))
	for tick := range am.hist {
		ticks = append(ticks, tick)
	}
	sort.Ints(ticks)

	pt0 := [2]int{0, 0}
	pe0 := [2]int{0, 0}
	for _, tick := range ticks {
		pt1 := [2]int{int(float64(tick) * graphXMul), int(am.hist[tick].t * graphYMul)}
		pe1 := [2]int{int(float64(tick) * graphXMul), int(am.hist[tick].e * graphYMul)}
		canvas.Line(pt0[0], pt0[1], pt1[0], pt1[1], "stroke:black;stroke_width:1")
		canvas.Line(pe0[0], pe0[1], pe1[0], pe1[1], "stroke:red;stroke_width:1")
		pt0 = pt1
		pe0 = pe1
	}

	canvas.End()
	return nil
}
