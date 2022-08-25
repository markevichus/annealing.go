package packer

import (
	"annealing/pkg/annealer"
	"errors"
	"fmt"
	svg "github.com/ajstarks/svgo"
	"io"
	"math/rand"
	"strconv"
	"time"
)

type HasArea interface {
	CalArea() float64
}

type Rectangle struct {
	Id     uint64
	Width  float64
	Height float64
	X      float64
	Y      float64
	Area   float64
}

func NewRectangle(id uint64, width float64, height float64) (rectangle *Rectangle, err error) {
	if width < 0 || height < 0 {
		return &Rectangle{}, errors.New("Rectangle W or H < 0")
	}
	return &Rectangle{
		Id:     id,
		Width:  width,
		Height: height,
	}, nil
}

func (r *Rectangle) CalArea() float64 {
	area := r.Width * r.Height
	r.Area = area
	return area
}

type Packer interface {
	SetDimensions(width float64, height float64)
	SetRectangles(rectangles []Rectangle)
	Compile() (float64, error)
	GetPlacedRectangles() []Rectangle
	StoreDraw()
}

type layoutChange struct {
	aIndex int
	bIndex int
	rotate bool
	rIndex int
}

type CutoutLayout struct {
	width             float64
	height            float64
	rectangles        []Rectangle
	placedRectangles  []Rectangle
	xSortedRectangles []Rectangle
	ySortedRectangles []Rectangle
	changes           []layoutChange
	energy            float64
	rs                rand.Source
}

func NewCutoutLayout(width float64, height float64) (cutout *CutoutLayout, err error) {
	if width < 0 || height < 0 {
		return &CutoutLayout{}, errors.New("CutoutLayout W or H < 0")
	}
	return &CutoutLayout{
		width:  width,
		height: height,
		rs:     rand.NewSource(time.Now().UnixNano()),
	}, nil
}

func (c *CutoutLayout) SetRectangles(rs []Rectangle) {
	// TODO validation
	c.reset()
	c.rectangles = make([]Rectangle, len(rs))
	for i, r := range rs {
		c.rectangles[i] = r
		c.rectangles[i].CalArea()
	}
}

func (c *CutoutLayout) reset() {
	// Border rectangles
	c.placedRectangles = []Rectangle{
		{Id: 0, Y: 0, X: 0, Width: 0, Height: c.height},
		{Id: 0, Y: 0, X: 0, Width: c.width, Height: 0},
	}
	c.xSortedRectangles = []Rectangle{}
	c.ySortedRectangles = []Rectangle{}
	for _, r := range c.placedRectangles {
		c.insertIntoSortedRectangles(r)
	}
	c.revertChanges()
}

func (c *CutoutLayout) Compile() error {
	for _, r := range c.rectangles {
		// Place input rectangle at top right corner
		r.X = c.width - r.Width
		r.Y = c.height

		_, _, err := c.placeRectangle(r)
		if err != nil {
			return err
		}
	}
	c.calculateEnergy()

	return nil
}

func (c *CutoutLayout) GetPlacedRectangles() []Rectangle {
	return c.placedRectangles
}

func (c *CutoutLayout) StoreDraw(w io.Writer) (err error) {
	divider := 10
	labelShiftX := 5
	labelShiftY := 5
	labelHeight := 17
	scale := func(num float64) int {
		return int(num / float64(divider))
	}
	rs := rand.NewSource(time.Now().UnixNano())
	rndInt := func(max int) int {
		return rand.New(rs).Intn(max)
	}
	rndStr := func(max int) string {
		return strconv.Itoa(rndInt(max))
	}
	randomColor := func() string {
		r := rndStr(200)
		g := rndStr(200)
		b := rndStr(200)
		return "rgb(" + r + "," + g + "," + b + ", 0.5)"
	}

	canvas := svg.New(w)
	canvas.Start(scale(c.width), scale(c.height))

	canvas.Rect(0, 0, scale(c.width), scale(c.height), "fill:none;stroke:red;stroke_width:2")
	canvas.Text(4*labelShiftX, 4*labelShiftY,
		fmt.Sprintf("%f", c.energy), "font-weight:bold;font-size:"+strconv.Itoa(int(float64(labelHeight)*1.5))+"px")

	for _, r := range c.placedRectangles {
		canvas.Rect(scale(r.X), scale(r.Y), scale(r.Width), scale(r.Height),
			"fill:"+randomColor()+";stroke:black;stroke_width:2")
		canvas.Text(scale(r.X)+labelShiftX, scale(r.Y+r.Height)-labelShiftY,
			strconv.Itoa(int(r.Id)), "font-weight:normal;font-size:"+strconv.Itoa(labelHeight)+"px")

		canvas.Text(scale(r.X)+labelShiftX, scale(r.Y)+labelShiftY*2,
			strconv.Itoa(int(r.X)), "font-weight:normal;font-size:10px")
		canvas.Text(scale(r.X)+labelShiftX, scale(r.Y)+labelShiftY*4,
			strconv.Itoa(int(r.Y+r.Height)), "font-weight:normal;font-size:10px")
	}

	canvas.End()
	return nil
}

func (c *CutoutLayout) placeRectangle(r Rectangle) (placedRectangle Rectangle, isPlaced bool, err error) {
	placed := false
	outOfHeight := false

	for _, placedR := range c.ySortedRectangles {
		if placedR.Y+placedR.Height > r.Y {
			continue
		}
		// Find X-crossing
		if placedR.X+placedR.Width > r.X && placedR.X < r.X+r.Width {
			// If we're standstill on it already
			if r.Y == placedR.Y+placedR.Height {
				if r.Y+r.Height <= c.height {
					placed = true
				} else {
					outOfHeight = true
				}
			} else {
				r.Y = placedR.Y + placedR.Height
			}
			break
		}
	}

	if !placed {
		for _, placedR := range c.xSortedRectangles {
			if placedR.X+placedR.Width > r.X {
				continue
			}
			// Find Y-crossing
			if placedR.Y+placedR.Height > r.Y && placedR.Y < r.Y+r.Height {
				if r.X == placedR.X+placedR.Width {
					if r.Y+r.Height <= c.height {
						placed = true
					} else {
						outOfHeight = true
					}
				} else {
					r.X = placedR.X + placedR.Width
				}
				break
			}
		}
	}

	if placed {
		c.placedRectangles = append(c.placedRectangles, r)
		c.insertIntoSortedRectangles(r)
		return r, true, nil
	} else if outOfHeight {
		return Rectangle{}, false, nil
	} else {
		return c.placeRectangle(r)
	}

}

func (c *CutoutLayout) insertIntoSortedRectangles(placedRectangle Rectangle) {
	var iIndex int
	found := false
	for i, sR := range c.ySortedRectangles {
		if placedRectangle.Y+placedRectangle.Height > sR.Y+sR.Height {
			iIndex = i
			found = true
			break
		}
	}
	if found {
		c.ySortedRectangles = append(c.ySortedRectangles[:iIndex+1], c.ySortedRectangles[iIndex:]...)
		c.ySortedRectangles[iIndex] = placedRectangle
	} else {
		c.ySortedRectangles = append(c.ySortedRectangles, placedRectangle)
	}

	found = false
	for i, sR := range c.xSortedRectangles {
		if placedRectangle.X+placedRectangle.Width > sR.X+sR.Width {
			iIndex = i
			found = true
			break
		}
	}
	if found {
		c.xSortedRectangles = append(c.xSortedRectangles[:iIndex+1], c.xSortedRectangles[iIndex:]...)
		c.xSortedRectangles[iIndex] = placedRectangle
	} else {
		c.xSortedRectangles = append(c.xSortedRectangles, placedRectangle)
	}
}

func (c *CutoutLayout) calculateEnergy() {
	sumArea := float64(0)
	for _, r := range c.placedRectangles {
		sumArea += r.Area
	}
	c.energy = 1 - sumArea/(c.width*c.height)
}

// Shaker

func (c *CutoutLayout) Shake() (energy float64, err error) {
	c.reset()

	// Find positions to change
	if len(c.rectangles) >= 2 {
		// Already has two border rectangles
		shakeUntilIndex := len(c.placedRectangles) - 2
		if shakeUntilIndex < 2 {
			shakeUntilIndex = len(c.rectangles) - 1
		}

		var aIndex, bIndex, rIndex int
		var rotate bool

		unchanged := true
		for unchanged {
			if rand.New(c.rs).Float64() > 0.2 {
				unchanged = false
				validPair := false
				for !validPair {
					aIndex = rand.New(c.rs).Intn(shakeUntilIndex + 1)
					bIndex = rand.New(c.rs).Intn(shakeUntilIndex + 1)
					if !c.areRectanglesSame(aIndex, bIndex) {
						validPair = true
					}
				}
			}
			if rand.New(c.rs).Float64() > 0.3 {
				unchanged = false
				rIndex = rand.New(c.rs).Intn(shakeUntilIndex + 1)
				rotate = true
			}
		}

		c.change(layoutChange{
			aIndex: aIndex,
			bIndex: bIndex,
			rotate: rotate,
			rIndex: rIndex,
		})
	}

	err = c.Compile()
	if err != nil {
		return 0, err
	}

	return c.energy, nil
}

func (c *CutoutLayout) areRectanglesSame(aIndex int, bIndex int) bool {
	// TODO compare shapes
	return aIndex == bIndex ||
		(c.rectangles[aIndex].Width == c.rectangles[bIndex].Width &&
			c.rectangles[aIndex].Height == c.rectangles[bIndex].Height)
}

func (c *CutoutLayout) Take() {
	c.changes = []layoutChange{}
}

func (c *CutoutLayout) change(ch layoutChange) {
	tmpR := c.rectangles[ch.aIndex]
	c.rectangles[ch.aIndex] = c.rectangles[ch.bIndex]
	c.rectangles[ch.bIndex] = tmpR

	if ch.rotate {
		width := c.rectangles[ch.rIndex].Width
		c.rectangles[ch.rIndex].Width = c.rectangles[ch.rIndex].Height
		c.rectangles[ch.rIndex].Height = width
	}
	c.changes = append(c.changes, ch)
}

func (c *CutoutLayout) revertChanges() {
	if len(c.changes) > 0 {
		for _, ch := range c.changes {
			if ch.rotate {
				c.rectangles[ch.rIndex].Width, c.rectangles[ch.rIndex].Height = c.rectangles[ch.rIndex].Height, c.rectangles[ch.rIndex].Width
			}
			if ch.aIndex != ch.bIndex {
				tmpR := c.rectangles[ch.aIndex]
				c.rectangles[ch.aIndex] = c.rectangles[ch.bIndex]
				c.rectangles[ch.bIndex] = tmpR
			}
		}
		c.changes = []layoutChange{}
	}
}

func (c *CutoutLayout) GetEnergy() float64 {
	return c.energy
}

func (c *CutoutLayout) GetResult() annealer.ShakeResult {
	sr := annealer.ShakeResult{Energy: c.energy}
	c.StoreDraw(&sr)
	return sr
}
