package main

import (
	"annealing/pkg/cooler"
	"annealing/pkg/packer"
	"fmt"
)

func main() {

	type RectData struct {
		id     uint64
		width  float64
		height float64
	}
	rectanglesData := []RectData{
		{id: 1, width: 498, height: 1486}, {id: 2, width: 498, height: 1486}, {id: 3, width: 517, height: 1302}, {id: 4, width: 517, height: 1302}, {id: 5, width: 517, height: 1302}, {id: 6, width: 517, height: 1302}, {id: 7, width: 535, height: 1237}, {id: 8, width: 535, height: 1237}, {id: 9, width: 488, height: 1260}, {id: 10, width: 488, height: 1260}, {id: 11, width: 413, height: 1198}, {id: 12, width: 413, height: 1198}, {id: 13, width: 413, height: 1198}, {id: 14, width: 413, height: 1198}, {id: 15, width: 422, height: 1362}, {id: 16, width: 422, height: 1362}, {id: 17, width: 372, height: 1362}, {id: 18, width: 372, height: 1362}, {id: 19, width: 372, height: 1362}, {id: 20, width: 372, height: 1362}, {id: 21, width: 413, height: 1198}, {id: 22, width: 413, height: 1198}, {id: 23, width: 155, height: 1237}, {id: 24, width: 155, height: 1237}, {id: 25, width: 155, height: 637}, {id: 26, width: 155, height: 637}, {id: 27, width: 498, height: 650}, {id: 28, width: 498, height: 650}, {id: 29, width: 498, height: 650}, {id: 30, width: 498, height: 650}, {id: 31, width: 498, height: 650}, {id: 32, width: 498, height: 650}, {id: 33, width: 498, height: 650}, {id: 34, width: 498, height: 650}, {id: 35, width: 498, height: 650}, {id: 36, width: 498, height: 650}, {id: 37, width: 498, height: 650}, {id: 38, width: 498, height: 650}, {id: 39, width: 498, height: 650}, {id: 40, width: 498, height: 650}, {id: 41, width: 498, height: 650}, {id: 42, width: 498, height: 650}, {id: 43, width: 498, height: 650}, {id: 44, width: 498, height: 650}, {id: 45, width: 498, height: 650}, {id: 46, width: 498, height: 650}, {id: 47, width: 498, height: 650}, {id: 48, width: 629, height: 1289}, {id: 49, width: 629, height: 1289}, {id: 50, width: 629, height: 1289}, {id: 51, width: 550, height: 1330}, {id: 52, width: 550, height: 1330}, {id: 53, width: 1410, height: 481}, {id: 54, width: 1410, height: 481}, {id: 55, width: 1410, height: 481}, {id: 56, width: 1410, height: 481}, {id: 57, width: 1410, height: 481}, {id: 58, width: 1410, height: 481}, {id: 59, width: 1410, height: 481}, {id: 60, width: 1410, height: 481}, {id: 61, width: 1410, height: 481}, {id: 62, width: 1410, height: 481}, {id: 63, width: 1410, height: 481}, {id: 64, width: 1410, height: 481}, {id: 65, width: 1410, height: 481}, {id: 66, width: 1410, height: 481}, {id: 67, width: 1410, height: 481}, {id: 68, width: 1410, height: 481}, {id: 69, width: 1410, height: 481}, {id: 70, width: 1410, height: 481}, {id: 71, width: 1410, height: 481}, {id: 72, width: 1410, height: 481}, {id: 73, width: 1410, height: 481}, {id: 74, width: 1410, height: 481}, {id: 75, width: 1410, height: 481}, {id: 76, width: 1410, height: 481}, {id: 77, width: 1410, height: 481}, {id: 78, width: 1410, height: 481}, {id: 79, width: 1410, height: 481}, {id: 80, width: 1322, height: 1457}, {id: 81, width: 1322, height: 1457}, {id: 82, width: 1322, height: 1457}, {id: 83, width: 1322, height: 1457}, {id: 84, width: 1322, height: 1457}, {id: 85, width: 1322, height: 1457}, {id: 86, width: 1322, height: 1457}, {id: 87, width: 1322, height: 1457}, {id: 88, width: 1322, height: 1457}, {id: 89, width: 1322, height: 1457}, {id: 90, width: 1322, height: 1457}, {id: 91, width: 1322, height: 1457}, {id: 92, width: 1322, height: 1457}, {id: 93, width: 1322, height: 1457}, {id: 94, width: 1322, height: 1457}, {id: 95, width: 1322, height: 1457}, {id: 96, width: 1322, height: 1457}, {id: 97, width: 1322, height: 1457}, {id: 98, width: 1322, height: 1457}, {id: 99, width: 1322, height: 1457}, {id: 100, width: 1322, height: 1457}, {id: 101, width: 782, height: 1457}, {id: 102, width: 782, height: 1457}, {id: 103, width: 782, height: 1457}, {id: 104, width: 782, height: 1457}, {id: 105, width: 782, height: 1457}, {id: 106, width: 782, height: 1457}, {id: 107, width: 1340, height: 1360}, {id: 108, width: 1340, height: 1360}, {id: 109, width: 1250, height: 1360}, {id: 110, width: 1250, height: 1360}, {id: 111, width: 1360, height: 431}, {id: 112, width: 1360, height: 431}, {id: 113, width: 1360, height: 431}, {id: 114, width: 902, height: 1457}, {id: 115, width: 902, height: 1457}, {id: 116, width: 902, height: 1457}, {id: 117, width: 902, height: 1457}, {id: 118, width: 902, height: 1457}, {id: 119, width: 902, height: 1457}, {id: 120, width: 902, height: 1457}, {id: 121, width: 902, height: 1457}, {id: 122, width: 902, height: 1457}, {id: 123, width: 902, height: 1457}, {id: 124, width: 902, height: 1457}, {id: 125, width: 902, height: 1457}, {id: 126, width: 902, height: 1457}, {id: 127, width: 902, height: 1457}, {id: 128, width: 902, height: 1457}, {id: 129, width: 902, height: 1457}, {id: 130, width: 902, height: 1457}, {id: 131, width: 902, height: 1457}, {id: 132, width: 902, height: 1457}, {id: 133, width: 902, height: 1457}, {id: 134, width: 902, height: 1457}, {id: 135, width: 902, height: 1457}, {id: 136, width: 902, height: 1457}, {id: 137, width: 902, height: 1457}, {id: 138, width: 902, height: 1457}, {id: 139, width: 902, height: 1457}, {id: 140, width: 902, height: 1457}, {id: 141, width: 538, height: 1405}, {id: 142, width: 538, height: 1405}, {id: 143, width: 538, height: 1405}, {id: 144, width: 538, height: 1405}, {id: 145, width: 538, height: 1405}, {id: 146, width: 538, height: 1405}, {id: 147, width: 538, height: 1405}, {id: 148, width: 538, height: 1405}, {id: 149, width: 538, height: 1405}, {id: 150, width: 538, height: 1405}, {id: 151, width: 538, height: 1405}, {id: 152, width: 538, height: 1405}, {id: 153, width: 538, height: 1405}, {id: 154, width: 538, height: 1405}, {id: 155, width: 538, height: 1405}, {id: 156, width: 538, height: 1405}, {id: 157, width: 538, height: 1405}, {id: 158, width: 538, height: 1405}, {id: 159, width: 538, height: 1405}, {id: 160, width: 538, height: 1405}, {id: 161, width: 538, height: 1405}, {id: 162, width: 538, height: 1405}, {id: 163, width: 538, height: 1405}, {id: 164, width: 538, height: 1405},
	}

	var rectangles []packer.Rectangle
	for _, rd := range rectanglesData {
		r, _ := packer.NewRectangle(rd.id, rd.width, rd.height)
		rectangles = append(rectangles, *r)
	}

	cutout, err := packer.NewCutoutLayout(3150, 6000)
	if err != nil {
		fmt.Errorf("error creating layout: %v\n", err.Error())
	}
	cutout.SetRectangles(rectangles)
	//cutout.Shake()
	//cutout.Compile()

	//cutout.Shake()
	//cutout.Shake()
	//cutout.Shake()
	//cutout.Shake()

	////fmt.Println(cutout.GetPlacedRectangles())
	//cutout.StoreReport()

	am := cooler.NewAnnealingMachine(cutout, 0.0015, 20000)
	err = am.Run()
	if err != nil {
		fmt.Errorf("error runing Annealer: %v\n", err.Error())
	}

	//rs := rand.NewSource(time.Now().UnixNano())
	//fmt.Println(rand.New(rs).Intn(100))
}
