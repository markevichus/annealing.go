package packer

import (
	"errors"
	"fmt"
	svg "github.com/ajstarks/svgo"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

type HasArea interface {
	Area() float64
}

type Rectangle struct {
	Id     uint64
	Width  float64
	Height float64
	X      float64
	Y      float64
}

func NewRectangle(id uint64, width float64, height float64) (rectangle Rectangle, err error) {
	if width < 0 || height < 0 {
		return Rectangle{}, errors.New("Rectangle W or H < 0")
	}
	return Rectangle{
		Id:     id,
		Width:  width,
		Height: height,
	}, nil
}

func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Packer interface {
	SetDimensions(width float64, height float64)
	SetRectangles(rectangles []Rectangle)
	Compile() float64
	GetEnergy() float64
	GetPlacedRectangles() []Rectangle
	StoreDraw()
}

type CutoutLayout struct {
	width            float64
	height           float64
	rectangles       []Rectangle
	placedRectangles []Rectangle
	energy           float64
}

func NewCutoutLayout(width float64, height float64) (cutout CutoutLayout, err error) {
	if width < 0 || height < 0 {
		return CutoutLayout{}, errors.New("CutoutLayout W or H < 0")
	}
	return CutoutLayout{
		width:  width,
		height: height,
	}, nil
}

func (c *CutoutLayout) SetRectangles(rs []Rectangle) {
	// TODO validation
	// Create border rectangles
	c.placedRectangles = []Rectangle{
		{Id: 0, Y: 0, X: 0, Width: 0, Height: c.height},
		{Id: 0, Y: 0, X: 0, Width: c.width, Height: 0},
	}
	c.rectangles = rs
}

func (c *CutoutLayout) Compile() (energy float64, err error) {
	for _, r := range c.rectangles {
		// Place input rectangle at top right corner
		r.X = c.width - r.Width
		r.Y = c.height

		placedRectangle, isPlaced, err := c.placeRectangle(r)
		if err != nil {
			return 0, err
		}
		if isPlaced {
			c.placedRectangles = append(c.placedRectangles, placedRectangle)
		}
	}

	c.calculateEnergy()
	return c.energy, nil
}

func (c *CutoutLayout) GetEnergy() float64 {
	return c.energy
}

func (c *CutoutLayout) GetPlacedRectangles() []Rectangle {
	return c.placedRectangles
}

func (c *CutoutLayout) StoreDraw() {
	divider := 10
	labelShiftX := 5
	labelShiftY := 5
	labelHeight := 17
	scale := func(num float64) int {
		return int(num / float64(divider))
	}
	rs := rand.NewSource(time.Now().UnixNano())
	random := rand.New(rs)
	rInt := func(max int) int {
		return random.Intn(max)
	}
	rStr := func(max int) string {
		return strconv.Itoa(rInt(max))
	}
	randomColor := func() string {
		r := rStr(200)
		g := rStr(200)
		b := rStr(200)
		return "rgb(" + r + "," + g + "," + b + ", 0.5)"
	}

	f, err := os.Create("/tmp/dat.svg")
	if err != nil {
		fmt.Errorf("error creating file. %v", err)
	}

	canvas := svg.New(f)
	width := 10000
	height := 10000
	canvas.Start(width, height)

	canvas.Rect(0, 0, scale(c.width), scale(c.height), "fill:none;stroke:red;stroke_width:2")

	for _, r := range c.placedRectangles {
		canvas.Rect(scale(r.X), scale(r.Y), scale(r.Width), scale(r.Height),
			"fill:"+randomColor()+";stroke:black;stroke_width:2")
		canvas.Text(scale(r.X)+labelShiftX, scale(r.Y+r.Height)-labelShiftY,
			strconv.Itoa(int(r.Id)), "font-weight:normal;font-size:"+strconv.Itoa(labelHeight)+"px")
	}

	canvas.End()
}

func (c *CutoutLayout) placeRectangle(r Rectangle) (placedRectangle Rectangle, isPlaced bool, err error) {

	ySortedRectangles := make([]Rectangle, len(c.placedRectangles))
	copy(ySortedRectangles, c.placedRectangles)
	sort.Slice(ySortedRectangles, func(i, j int) bool {
		return ySortedRectangles[i].Y+ySortedRectangles[i].Height > ySortedRectangles[j].Y+ySortedRectangles[j].Height
	})

	xSortedRectangles := make([]Rectangle, len(c.placedRectangles))
	copy(xSortedRectangles, c.placedRectangles)
	sort.Slice(xSortedRectangles, func(i, j int) bool {
		return xSortedRectangles[i].X+xSortedRectangles[i].Width > xSortedRectangles[j].X+xSortedRectangles[j].Width
	})

	for _, placedR := range ySortedRectangles {
		if placedR.Y+placedR.Height > r.Y {
			continue
		}
		// Find X-crossing
		if (placedR.X+placedR.Width <= r.X && placedR.X+placedR.Width > r.X) ||
			(placedR.X+placedR.Width > r.X && placedR.X < r.X+r.Width) {
			// If we're standstill on it already
			if r.Y == placedR.Y+placedR.Height {
				if r.Y+r.Height <= c.height {
					return r, true, nil
				} else {
					return Rectangle{}, false, nil
				}
			} else {
				r.Y = placedR.Y + placedR.Height
				break
			}
		}
	}

	for _, placedR := range xSortedRectangles {
		if placedR.X+placedR.Width > r.X {
			continue
		}
		// Find Y-crossing
		if (placedR.Y+placedR.Height <= r.Y && placedR.Y+placedR.Height > r.Y) ||
			(placedR.Y+placedR.Height > r.Y && placedR.Y < r.Y+r.Height) {
			if r.X == placedR.X+placedR.Width {
				if r.X+r.Width <= c.width {
					return r, true, nil
				} else {
					return Rectangle{}, false, nil
				}
			} else {
				r.X = placedR.X + placedR.Width
				break
			}
		}
	}

	fmt.Println(r)
	return c.placeRectangle(r)
}

func (c *CutoutLayout) calculateEnergy() {
	sumArea := float64(0)

	for _, r := range c.placedRectangles {
		sumArea += r.Area()
	}

	e := 1 - sumArea/(c.width*c.height)
	c.energy = math.Round(e*100) / 100
}
