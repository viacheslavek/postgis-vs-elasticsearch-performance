package storage

type Storage interface {
	AddPoint(p Point) error
	DeletePoint(p Point) error
	IsPoint(p Point) (bool, error)

	GetInRadius(p Point, radius float64) ([]Point, error)
	GetInPolygon(polygon []float64) ([]Point, error)
}

type Point struct {
	Latitude  float64
	Longitude float64
}
