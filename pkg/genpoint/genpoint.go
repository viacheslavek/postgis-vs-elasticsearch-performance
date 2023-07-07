package genpoint

type getPoint interface {
	generate(N int) []Point
}

type Point struct {
	x int
	y int
}
