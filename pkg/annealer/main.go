package annealer

import (
	"fmt"
	svg "github.com/ajstarks/svgo"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

type ResultHandler interface {
	StoreReport() (err error)
}
type ShakeResult struct {
	Id     string
	Energy float64
	Tick   int
	Temp   float64
	Draw   string
	Closed bool
}

func (sr *ShakeResult) Write(b []byte) (n int, err error) {
	if sr.Closed {
		sr.Draw = ""
		sr.Closed = false
	}
	sr.Draw = sr.Draw + string(b)
	return len(b), nil
}
func (sr *ShakeResult) Read(b []byte) (n int, err error) {
	b = []byte(sr.Draw)
	return len(sr.Draw), nil
}
func (sr *ShakeResult) Close() error {
	sr.Closed = true
	return nil
}
func (sr *ShakeResult) StoreReport() (err error) {
	//f, err := os.Create("layout" + sr.Id + ".svg")
	f, err := os.Create("layout.svg")
	if err != nil {
		fmt.Errorf("error creating file. %v", err)
		return err
	}
	f.Write([]byte(sr.Draw))
	return nil
}

type Shaker interface {
	Shake() (energy float64, err error)
	Take()
	GetEnergy() float64
	GetResult() ShakeResult
}

type Annealer interface {
	Run(id string) (ShakeResult, error)
}

type AnnealingHist struct {
	t float64
	e float64
}

type AnnealingMachine struct {
	id       string
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

func (am *AnnealingMachine) Run(id string) (result ShakeResult, err error) {
	am.id = id
	am.reset()
	var minE float64 = 1
	var bestSR ShakeResult

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
			bestSR = am.shaker.GetResult()
			bestSR.Id = am.id
			bestSR.Tick = am.tick
			bestSR.Temp = am.temp
		} else if am.shouldITransit(dE) {
			//fmt.Println("random", e, am.temp, am.tick)
			am.transit()
		}
		am.decreaseTemperature()

		if am.tick == am.stopTick || am.temp <= am.stopTemp {
			//fmt.Println(bestSR.Id, fmt.Sprintf("\t"), bestSR.Energy, bestSR.Tick, bestSR.Temp)
			fmt.Printf("%v\t%v %v %v\n", bestSR.Id, bestSR.Energy, bestSR.Tick, bestSR.Temp)
			break
		}
		am.tick += 1
	}

	return bestSR, nil
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

func (am *AnnealingMachine) StoreReport() (err error) {
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
