package seamcarve

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

// Do n iterations of seam carving
func Carve(img image.Image, n int) (image.Image, error) {

	width := img.Bounds().Max.X
	if n > width {
		return nil, fmt.Errorf("number of seams %d greater than image width", n)
	}

	for n > 0 {
		img = CarveOne(img)
		n--
	}

	return img, nil

}

// Do one iteration of seam carving
func CarveOne(img image.Image) image.Image {

	width, height := img.Bounds().Max.X, img.Bounds().Max.Y

	// Compute gradients
	grads := make([][]float64, height)
	for i := 0; i < height-1; i++ {
		// First column is +Inf, removed afterwards
		row := make([]float64, width+1)
		row[0] = math.Inf(1)
		for j := 0; j < width-1; j++ {
			pix := img.At(j, i)
			rPix := img.At(j+1, i)
			dPix := img.At(j, i+1)

			row[j+1] = pixDif(pix, rPix) + pixDif(pix, dPix)
		}
		row[width-1] = math.Inf(1)
		grads[i] = row
	}
	lRow := make([]float64, width)
	for j := 0; j < width; j++ {
		lRow[j] = 0.0
	}
	grads[height-1] = lRow

	// Compute accumulatives
	accu := make([][]float64, height)
	accu[0] = grads[0]
	for i := 1; i < height; i++ {
		row := make([]float64, width+1)
		row[0] = math.Inf(1)
		row[width] = math.Inf(1)
		for j := 1; j < width; j++ {
			prevL := accu[i-1][j-1]
			prevU := accu[i-1][j]
			prevR := accu[i-1][j+1]

			min := math.Min(prevL, math.Min(prevU, prevR))
			row[j] = grads[i][j] + min
		}
		accu[i] = row
	}

	// Compute min seam
	seam := make([]int, height)
	lastMin, minIdx := accu[height-1][0], 0
	for j := 1; j < width; j++ {
		if accu[height-1][j] < lastMin {
			lastMin, minIdx = accu[height-1][j], j
		}
	}
	seam[height-1] = minIdx - 1

	for i := height - 2; i >= 0; i-- {
		seamJ := seam[i+1]
		prevL := accu[i][seamJ]
		prevU := accu[i][seamJ+1]
		prevR := accu[i][seamJ+2]

		if prevL <= prevU && prevL <= prevR {
			seam[i] = seamJ - 1
		} else if prevU <= prevL && prevU <= prevR {
			seam[i] = seamJ
		} else {
			seam[i] = seamJ + 1
		}
	}

	newImg := image.NewNRGBA(image.Rect(0, 0, width-1, height))
	for i := 0; i < height; i++ {
		seamJ := seam[i]
		for j := 0; j < seamJ; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			newImg.Pix[newImg.PixOffset(j, i)] = uint8(r / 256)
			newImg.Pix[newImg.PixOffset(j, i)+1] = uint8(g / 256)
			newImg.Pix[newImg.PixOffset(j, i)+2] = uint8(b / 256)
			newImg.Pix[newImg.PixOffset(j, i)+3] = 255
		}
		for j := seamJ + 1; j < width-1; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			newImg.Pix[newImg.PixOffset(j-1, i)] = uint8(r / 256)
			newImg.Pix[newImg.PixOffset(j-1, i)+1] = uint8(g / 256)
			newImg.Pix[newImg.PixOffset(j-1, i)+2] = uint8(b / 256)
			newImg.Pix[newImg.PixOffset(j-1, i)+3] = 255
		}
	}

	return newImg
}

func pixDif(p1 color.Color, p2 color.Color) float64 {
	r1, g1, b1, _ := p1.RGBA()
	r2, g2, b2, _ := p2.RGBA()

	return math.Pow(float64(r1-r2), 2) +
		math.Pow(float64(g1-g2), 2) +
		math.Pow(float64(b1-b2), 2)
}
