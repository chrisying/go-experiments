package main

import (
	"fmt"
	"go-experiments/fractal/lib"
	"image"
	"image/png"
	"os"
	"runtime"
)

// Image size in pixels
var SIZE int

// Output file name
var OUTPUT string

// Constants used for randomizing coloring scheme
var RMIX int
var GMIX int
var BMIX int

func main() {

	runtime.GOMAXPROCS(4)
	// Feel free to change any of these, arbitrarity chosen
	SIZE, OUTPUT, RMIX, GMIX, BMIX = 800, "output.png", 273, 89, 171

	//img := runSeq()
	img := runPara()

	// Save image to file
	stdin, _ := os.Create(OUTPUT)
	defer stdin.Close()
	err := png.Encode(stdin, img)
	if err != nil {
		fmt.Println(err)
	}
}

// Sequential implementation
func runSeq() *image.NRGBA {

	// Create image array using image builtin package
	img := image.NewNRGBA(image.Rect(0, 0, SIZE, SIZE))
	for row := img.Rect.Min.Y; row < img.Rect.Max.Y; row++ {
		for col := img.Rect.Min.X; col < img.Rect.Max.X; col++ {

			X := float32(col)/(float32(SIZE)/4.0) - 2.0
			Y := float32(row)/(float32(SIZE)/4.0) - 2.0

			n := fractal.Julia(fractal.ComplexNumber{X, Y})
			// n := fractal.Mandlebrot(fractal.ComplexNumber{X, Y})
			// fmt.Printf("X:%f, Y:%f, n:%d\n", X, Y, n)

			// Notice that img.Pix is a 1-D array!
			img.Pix[img.PixOffset(col, row)] = uint8((n * RMIX) % 256)
			img.Pix[img.PixOffset(col, row)+1] = uint8((n * GMIX) % 256)
			img.Pix[img.PixOffset(col, row)+2] = uint8((n * BMIX) % 256)
			img.Pix[img.PixOffset(col, row)+3] = 255
		}
	}

	return img

}

// Parallel implementation
// For demonstration purposes, uses 4 threads
// In practice, the parallel implementation is actually slower

type Packet struct {
	n   int
	row int
	col int
}

func runQuadrant(lowX, highX, lowY, highY int, c chan Packet) {

	for row := lowY; row < highY; row++ {
		for col := lowX; col < highX; col++ {

			X := float32(col)/(float32(SIZE)/4.0) - 2.0
			Y := float32(row)/(float32(SIZE)/4.0) - 2.0

			n := fractal.Julia(fractal.ComplexNumber{X, Y})
			// n := fractal.Mandlebrot(fractal.ComplexNumber{X, Y})

			c <- Packet{n, row, col}
		}
	}

}

func runPara() *image.NRGBA {

	// Create image array using image builtin package
	img := image.NewNRGBA(image.Rect(0, 0, SIZE, SIZE))

	maxX := img.Rect.Max.X
	maxY := img.Rect.Max.Y
	midX := maxX / 2
	midY := maxY / 2
	c := make(chan Packet)

	go runQuadrant(0, midX, 0, midY, c)
	go runQuadrant(0, midX, midY, maxY, c)
	go runQuadrant(midX, maxX, 0, midY, c)
	go runQuadrant(midX, maxX, midY, maxY, c)

	for i := 0; i < SIZE*SIZE; i++ {
		pack := <-c
		row := pack.row
		col := pack.col
		n := pack.n

		// Notice that img.Pix is a 1-D array!
		img.Pix[img.PixOffset(col, row)] = uint8((n * RMIX) % 256)
		img.Pix[img.PixOffset(col, row)+1] = uint8((n * GMIX) % 256)
		img.Pix[img.PixOffset(col, row)+2] = uint8((n * BMIX) % 256)
		img.Pix[img.PixOffset(col, row)+3] = 255
	}

	return img
}
