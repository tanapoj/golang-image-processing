package main

import (
	"fmt"
	"github.com/thoas/go-funk"

	"image"
	"image/png"
	"image/jpeg"
	"os"
	"image/draw"
)

func test0() {
	fmt.Println("Hello, Go")

	r := funk.Map([]int{1, 2, 3, 4}, func(x int) int {
		return x * 2
	})
	fmt.Println(r)
}

func test1() {

	// Create a blank image 10 pixels wide by 4 pixels tall
	img := image.NewRGBA(image.Rect(0, 0, 10, 4))

	// You can access the pixels through myImage.Pix[i]
	// One pixel takes up four bytes/uint8. One for each: RGBA
	// So the first pixel is controlled by the first 4 elements
	// Values for color are 0 black - 255 full color
	// Alpha value is 0 transparent - 255 opaque
	img.Pix[0] = 255 // 1st pixel red
	img.Pix[1] = 0   // 1st pixel green
	img.Pix[2] = 0   // 1st pixel blue
	img.Pix[3] = 255 // 1st pixel alpha

	// myImage.Pix contains all the pixels
	// in a one-dimensional slice
	fmt.Println(img.Pix)

	// Stride is how many bytes take up 1 row of the image
	// Since 4 bytes are used for each pixel, the stride is
	// equal to 4 times the width of the image
	// Since all the pixels are stored in a 1D slice,
	// we need this to calculate where pixels are on different rows.
	fmt.Println(img.Stride) // 40 for an image 10 pixels wide

	outputFile, err := os.Create("test.png")
	if err != nil {
		// Handle error
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, img)

	// Don't forget to close files
	outputFile.Close()
}

func test2() {
	existingImageFile, err := os.Open("fuji.jpg")
	if err != nil {
		// Handle error
	}
	defer existingImageFile.Close()

	// Calling the generic image.Decode() will tell give us the data
	// and type of image it is as a string. We expect "png"
	imageData, imageType, err := image.Decode(existingImageFile)
	if err != nil {
		// Handle error
	}
	fmt.Println(imageData)
	fmt.Println(imageType)

	// We only need this because we already read from the file
	// We have to reset the file pointer back to beginning
	existingImageFile.Seek(0, 0)

	// Alternatively, since we know it is a png already
	// we can call png.Decode() directly
	loadedImage, err := jpeg.Decode(existingImageFile)
	if err != nil {
		// Handle error
	}
	fmt.Println(loadedImage)
}

func test3() {
	existingImageFile, err := os.Open("fuji.jpg")
	if err != nil {
		// Handle error
	}
	defer existingImageFile.Close()
	im, err := jpeg.Decode(existingImageFile)
	if err != nil {
		// Handle error
	}

	b := im.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), im, b.Min, draw.Src)

	fmt.Println(len(img.Pix))

	var gs [][]int
	for i := 0; i < b.Dy(); i++ {
		row := make([]int, b.Dx())
		gs = append(gs, row)
	}

	for i := 0; i < len(img.Pix); i += 4 {
		pixelAt := int(i / 4)
		row := pixelAt / b.Dx()
		col := pixelAt % b.Dx()

		var r int = int(img.Pix[i+0])
		var g int = int(img.Pix[i+1])
		var b int = int(img.Pix[i+2])
		sum := r + g + b
		avg := sum / 3
		gs[row][col] = avg
	}

	//gs = frameArray(gs, filterBlur)
	gs = frameArray(gs, filterEdge)

	for i := 0; i < len(img.Pix); i += 4 {
		pixelAt := int(i / 4)
		row := pixelAt / b.Dx()
		col := pixelAt % b.Dx()

		img.Pix[i+0] = uint8(gs[row][col])
		img.Pix[i+1] = uint8(gs[row][col])
		img.Pix[i+2] = uint8(gs[row][col])
		img.Pix[i+3] = 255
	}
	//for i := 0; i < len(img.Pix); i += 4 {
	//	var r int = int(img.Pix[i+0])
	//	var g int = int(img.Pix[i+1])
	//	var b int = int(img.Pix[i+2])
	//	sum := r + g + b
	//	avg := sum / 3
	//
	//	if i%1000 == 0 {
	//		fmt.Printf("%d %d %d = %d %d\n", img.Pix[i+0], img.Pix[i+1], img.Pix[i+2], sum, avg)
	//	}
	//	if avg >= 128 {
	//		img.Pix[i+0] = 255
	//		img.Pix[i+1] = 255
	//		img.Pix[i+2] = 255
	//	} else {
	//		img.Pix[i+0] = 0
	//		img.Pix[i+1] = 0
	//		img.Pix[i+2] = 0
	//	}
	//
	//	bb := img.Pix[i+0]
	//
	//	img.Pix[i+0] = bb
	//	img.Pix[i+1] = bb
	//	img.Pix[i+2] = bb
	//}

	img.Pix[0] = 255 // 1st pixel red
	img.Pix[1] = 0   // 1st pixel green
	img.Pix[2] = 0   // 1st pixel blue
	img.Pix[3] = 255 // 1st pixel alpha

	img.Pix[4] = 0
	img.Pix[5] = 255
	img.Pix[6] = 0
	img.Pix[7] = 255

	outputFile, err := os.Create("test2.png")
	if err != nil {
		// Handle error
	}

	png.Encode(outputFile, img)

	outputFile.Close()
}

//type Matrix [...][...]int

func getCell(matrix [][]int, i int, j int) int {
	if i < 0 || i >= len(matrix) {
		return -1
	}
	if j < 0 || j >= len(matrix[i]) {
		return -1
	}
	return matrix[i][j]
}

func dupArray(matrix [][]int) [][]int {
	duplicate := make([][]int, len(matrix))
	for i := range matrix {
		duplicate[i] = make([]int, len(matrix[i]))
		for j := 0; j < len(duplicate[i]); j++ {
			duplicate[i][j] = matrix[i][j]
		}
	}
	return duplicate
}

func frameArray(matrix [][]int, method func(matrix [][]int, i int, j int) int) [][]int {
	output := dupArray(matrix)
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			output[i][j] = method(matrix, i, j)
		}
	}
	return output
}

func filterBlur(matrix [][]int, i int, j int) int {
	at := func(i int, j int) int {
		return getCell(matrix, i, j)
	}
	frame := [9]int{
		at(i-1, j-1), at(i-1, j), at(i-1, j+1),
		at(i, j-1), at(i, j), at(i, j+1),
		at(i+1, j-1), at(i+1, j), at(i+1, j+1),
	}
	r := funk.Filter(frame, func(x int) bool {
		return x >= 0
	}).([]int)
	sum := funk.Reduce(r, func(x int, y int) int {
		return x + y
	}, 0)
	//var rr []int
	//From(frame).WhereT(func(x int) bool {
	//	return x >= 0
	//}).ToSlice(&r)
	return int(sum) / len(r)
}
func filterEdge(matrix [][]int, i int, j int) int {
	at := func(i int, j int, then int) int {
		cell := getCell(matrix, i, j)
		if cell < 0 {
			return -1000
		}
		return cell * then
	}
	M := -1
	Z := 0
	O := 1
	//frame := [9]int{
	//	at(i-1, j-1, M), at(i-1, j, Z), at(i-1, j+1, O),
	//	at(i, j-1, M), at(i, j, Z), at(i, j+1, O),
	//	at(i+1, j-1, M), at(i+1, j, Z), at(i+1, j+1, O),
	//}
	frame := [9]int{
		at(i-1, j-1, M), at(i-1, j, M), at(i-1, j+1, M),
		at(i, j-1, Z), at(i, j, Z), at(i, j+1, Z),
		at(i+1, j-1, O), at(i+1, j, O), at(i+1, j+1, O),
	}
	r := funk.Filter(frame, func(x int) bool {
		return x != -1000
	}).([]int)
	sum := funk.Reduce(r, func(x int, y int) int {
		return x + y
	}, 0)
	//var rr []int
	//From(frame).WhereT(func(x int) bool {
	//	return x >= 0
	//}).ToSlice(&r)
	if int(sum)/len(r) > 128 {
		return 255
	} else {
		return 0
	}
}

func main() {

	test3()

	arr := [][]int{
		{0, 0, 5, 1, 10, 10},
		{0, 0, 5, 1, 10, 10},
		{0, 0, 5, 1, 10, 10},
	}
	fmt.Println(arr)
	fmt.Println(frameArray(arr, filterBlur))
}
