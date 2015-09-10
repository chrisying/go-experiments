package main

import (
	"flag"
	"fmt"
	"go-experiments/edgedetect/lib"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
)

func main() {
	outPtr := flag.String("o", "out.png", "output image")
	thrPtr := flag.Float64("t", 0.1, "edge threshold [0, 1) (0.1 default)")
	graPtr := flag.Int("g", 1, "grain value [1,) (1 default = no grain)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Incorrect number of arguments")
		fmt.Println("Usage: ./main [cl args] [input image]")
		panic(fmt.Errorf("Input error"))
	}

	// Read input image
	inName := args[0]
	inImgF, err := os.Open(inName)
	defer inImgF.Close()
	if err != nil {
		panic(err.Error())
	}

	inImg, _, err := image.Decode(inImgF)
	if err != nil {
		panic(err.Error())
	}

	// Compute edge detect
	outImg := edgedetect.Process(inImg, *thrPtr, *graPtr)

	// Write output image
	outImgF, err := os.Create(*outPtr)
	defer outImgF.Close()
	if err != nil {
		panic(err.Error())
	}

	err = png.Encode(outImgF, outImg)
	if err != nil {
		panic(err.Error())
	}

}
