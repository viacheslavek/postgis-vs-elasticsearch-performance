package internal

type Point struct {
	Latitude  float64
	Longitude float64
}

type Polygon struct {
	Vertical []Point
}

type Circle struct {
	Centre Point
	Radius int
}

type Shapes struct {
	Polygons []Polygon
	Circles  []Circle
}
