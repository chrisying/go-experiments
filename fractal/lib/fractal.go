package fractal

import ()

type ComplexNumber struct {
	Re float32
	Im float32
}

func Add(c1, c2 ComplexNumber) ComplexNumber {
	return ComplexNumber{c1.Re + c2.Re, c1.Im + c2.Im}
}

func Multiply(c1, c2 ComplexNumber) ComplexNumber {
	return ComplexNumber{
		c1.Re*c2.Re - c1.Im*c2.Im,
		c1.Re*c2.Im + c1.Im*c2.Re}
}

func MagSquare(c ComplexNumber) float32 {
	return c.Re*c.Re + c.Im*c.Im
}

func Mandlebrot(c ComplexNumber) int {
  ESCAPE := 200
  n := 0
  z := ComplexNumber{c.Re, c.Im}
  for MagSquare(z) <= 4.0 && n <= ESCAPE {
    z = Add(Multiply(z, z), c)
    n++
  }
  return n
}

func Julia(z ComplexNumber) int {
  ESCAPE := 200
  n := 0
  c := ComplexNumber{-0.8, 0.156}
  for MagSquare(z) <= 4.0 && n <= ESCAPE {
    z = Add(Multiply(z, z), c)
    n++
  }
  return n
}
