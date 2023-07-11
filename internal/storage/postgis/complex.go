package postgis

import (
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
)

func (s *Storage) GetInRadius(p internal.Point, radius float64) ([]internal.Point, error) {
	return nil, nil
}

func (s *Storage) GetInPolygon(polygon []float64) ([]internal.Point, error) {
	return nil, nil
}
