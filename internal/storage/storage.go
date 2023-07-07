package storage

type Storage interface {
	AddPoint(p Point) error
	//DeletePoint(p Point) error
	//GetPoint() (Point, error)

	// Основные для поиска
	// GetInRadius(radius float64) ([]Point, error)
	// GetInPolygon(polygon []float64) ([]Point, error)

}

type Point struct {
	Latitude  float64
	Longitude float64
}
