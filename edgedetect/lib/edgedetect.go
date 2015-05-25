package edgedetect

import (
  "image"
  "image/color"
)

// Do entire edge detection routine
func Process(img image.Image, t float64, g int) image.Image {

  // width, height := img.Bounds().Max.X, img.Bounds().Max.Y

  bwImg := ConvertBW(img)

  // matrix should be odd width and height
  matrix := [][]int{[]int{-1, -1, 0},
                    []int{-1,  0, 1},
                    []int{ 0,  1, 1}}
  mapImg := MapMatrix(bwImg, matrix, t)

  filImg := Filter(mapImg, g)

  return filImg
}

// Converts input image color model to GrayColorModel
func ConvertBW(img image.Image) *image.Gray {

  width, height := img.Bounds().Max.X, img.Bounds().Max.Y
  bwImg := image.NewGray(image.Rect(0, 0, width, height))
  for x := 0; x < width; x++ {
    for y := 0; y < height; y++ {
      bwImg.Set(x, y, img.At(x, y))
    }
  }

  return bwImg
}

// Maps the matrix over the bw image, setting boundary to black
func MapMatrix(img *image.Gray, matrix [][]int, t float64) *image.Gray {

  width, height := img.Bounds().Max.X, img.Bounds().Max.Y
  mw, mh := len(matrix[0]), len(matrix)

  // Compute threshold based on matrix and t
  max := 0
  for x := 0; x < mw; x++ {
    for y := 0; y < mh; y++ {
      if matrix[y][x] > 0 {
        max += matrix[y][x] * 255
      }
    }
  }
  thresh := int(float64(max) * t)

  mapImg := image.NewGray(image.Rect(0, 0, width, height))
  for x := 0; x < width; x++ {
    for y := 0; y < height; y++ {

      // If matrix cannot be computed (on boundary)
      if ((x < mw / 2) ||
          (y < mh / 2) ||
          (width - x <= mw / 2) ||
          (height - y <= mh / 2)) {
        mapImg.Set(x, y, color.Gray{0})

      // Normal computation otherwise
      } else {
        sum := 0
        for xo := 0; xo < mw; xo++ {
          for yo := 0; yo < mh; yo++ {
            ma := matrix[yo][xo]
            im := img.GrayAt(x - mw / 2 + xo, y - mh / 2 + yo).Y
            sum += ma * int(im)
          }
        }

        if sum > thresh {
          mapImg.Set(x, y, color.Gray{255})
        } else {
          mapImg.Set(x, y, color.Gray{0})
        }
      }
    }
  }

  return mapImg
}

// Filter out grain by looking for connected regions of size g or larger
func Filter(img *image.Gray, g int) *image.Gray {

  // No filtering when grain is minimum
  if g == 1 {
    return img
  }

  width, height := img.Bounds().Max.X, img.Bounds().Max.Y
  mark := make([][]bool, height)
  for y := 0; y < height; y ++ {
    mark[y] = make([]bool, width)
  }

  filImg := image.NewGray(image.Rect(0, 0, width, height))
  for x := 0; x < width; x++ {
    for y := 0; y < height; y++ {
      // If marked, means it has been set already
      if !mark[y][x] {

        if img.GrayAt(x, y).Y == 0 {
          mark[y][x] = true
          // filImg.Set(x, y, color.Gray{0})
        } else {
          n := _FloodCount(img, mark, x, y)
          if n >= g {
            _FloodFill(img, filImg, x, y)
          }
        }
      }
    }
  }

  return filImg
}

// Helper function, count connected white pixels
func _FloodCount(img *image.Gray, mark [][]bool, x, y int) int {

  width, height := img.Bounds().Max.X, img.Bounds().Max.Y
  mark[y][x] = true
  sum := 1
  if x > 0 && !mark[y][x-1] && img.GrayAt(x-1, y).Y == 255 {
    sum += _FloodCount(img, mark, x-1, y)
  }
  if y > 0 && !mark[y-1][x] && img.GrayAt(x, y-1).Y == 255 {
    sum += _FloodCount(img, mark, x, y-1)
  }
  if x < width - 1 && !mark[y][x+1] && img.GrayAt(x+1, y).Y == 255 {
    sum += _FloodCount(img, mark, x+1, y)
  }
  if x < height - 1 && !mark[y+1][x] && img.GrayAt(x, y+1).Y == 255 {
    sum += _FloodCount(img, mark, x, y+1)
  }

  return sum
}

// Helper function, fill the regions that pass the grain size
func _FloodFill(img, filImg *image.Gray, x, y int) {

  width, height := img.Bounds().Max.X, img.Bounds().Max.Y
  filImg.Set(x, y, color.Gray{255})
  if x > 0 && filImg.GrayAt(x-1, y).Y != 255 && img.GrayAt(x-1, y).Y == 255 {
    _FloodFill(img, filImg, x-1, y)
  }
  if x > 0 && filImg.GrayAt(x, y-1).Y != 255 && img.GrayAt(x, y-1).Y == 255 {
    _FloodFill(img, filImg, x, y-1)
  }
  if x < width - 1 && filImg.GrayAt(x+1, y).Y != 255 && img.GrayAt(x+1, y).Y == 255 {
    _FloodFill(img, filImg, x+1, y)
  }
  if x < height - 1 && filImg.GrayAt(x, y+1).Y != 255 && img.GrayAt(x, y+1).Y == 255 {
    _FloodFill(img, filImg, x, y+1)
  }
}
