package main

import (
  "fmt"
  "flag"
  "image"
  _ "image/jpeg"
  "image/png"
  "os"
  "seamcarve/lib"
)

func main () {
  outPtr := flag.String("o", "out.png", "output image")
  numPtr := flag.Int("n", 100, "number of seams")
  flag.Parse()

  args := flag.Args()
  if len(args) != 1 {
    fmt.Println("Incorrect number of arguments")
    fmt.Println("Usage: ./main [cl args] [input image]")
    os.Exit(2)
  }

  // Read input image
  inName := args[0]
  inImgF, err := os.Open(inName)
  defer inImgF.Close()
  if err != nil {
    fmt.Println(err);
    os.Exit(2)
  }

  inImg, _, err := image.Decode(inImgF)
  if err != nil {
    fmt.Println(err);
    os.Exit(2)
  }

  // Compute seam carve
  outImg, err := seamcarve.Carve(inImg, *numPtr)

  // Write output image
  outImgF, err := os.Create(*outPtr)
  defer outImgF.Close()
  if err != nil {
    fmt.Println(err);
    os.Exit(2)
  }

  err = png.Encode(outImgF, outImg)
  if err != nil {
    fmt.Println(err)
    os.Exit(2)
  }

}
