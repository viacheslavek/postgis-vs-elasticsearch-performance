package postgis

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"log"
	"strconv"
	"strings"
)

func (s *Storage) GetInRadius(ctx context.Context, p internal.Point, radius int) ([]internal.Point, error) {

	q := `
		SELECT ST_AsText(geom)
		FROM moscow_region
		WHERE ST_DWithin(
		geom,
		ST_SetSRID(ST_Point($1, $2), 4326),
		$3
	);
	`

	rows, err := s.db.Query(ctx, q, p.Longitude, p.Latitude, radius)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("can't get in radius %w\n", err)
	}

	nearestPoints := make([]internal.Point, 0)

	for rows.Next() {
		var geometryData string
		err = rows.Scan(&geometryData)
		if err != nil {
			return nil, fmt.Errorf("can't scan row in nearestPoint %w\n", err)
		}

		geometryData = strings.TrimLeft(geometryData, "POINT(")
		geometryData = strings.TrimRight(geometryData, ")")
		coordinates := strings.Split(geometryData, " ")

		if len(coordinates) != 2 {
			log.Printf("too much argument\n")
			continue
		}

		latitude, err := strconv.ParseFloat(coordinates[1], 64)
		if err != nil {
			continue
		}
		longitude, err := strconv.ParseFloat(coordinates[0], 64)
		if err != nil {
			continue
		}

		nearestPoints = append(nearestPoints, internal.Point{Latitude: latitude, Longitude: longitude})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows %w\n", err)
	}

	return nearestPoints, nil
}

func (s *Storage) GetInPolygon(polygon []internal.Point) ([]internal.Point, error) {
	return nil, nil
}
