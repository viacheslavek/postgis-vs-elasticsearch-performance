package postgis

import (
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
)

func (s *Storage) GetInRadius(p storage.Point, radius float64) ([]storage.Point, error) {
	return nil, nil
}

func (s *Storage) GetInPolygon(polygon []float64) ([]storage.Point, error) {
	return nil, nil
}
