package main

import (
	"annealing/pkg/packer"
	"fmt"
)

func main() {
	var rectangles []packer.Rectangle

	r1, _ := packer.NewRectangle(1, 1024, 518)
	r2, _ := packer.NewRectangle(2, 1200, 450)
	r3, _ := packer.NewRectangle(3, 2046, 1567)
	rectangles = append(rectangles, r1, r2, r3)

	cutout, err := packer.NewCutoutLayout(3150, 6000)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	cutout.SetRectangles(rectangles)
	energy, err := cutout.Compile()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Println("Energy", energy)

	//for _, r := range cutout.placedRectangles {
	//	fmt.Println(r)
	//}
}
