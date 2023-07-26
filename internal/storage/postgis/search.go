package postgis

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/jackc/pgx/v5"
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
	)
		LIMIT 10000;
	`

	rows, err := s.db.Query(ctx, q, p.Longitude, p.Latitude, radius)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("can't get in radius %w\n", err)
	}

	return translateRowsPoint(rows)

}

func (s *Storage) GetInPolygon(ctx context.Context, polygon []internal.Point) ([]internal.Point, error) {

	polygonWKT := translatePolygonToWKT(polygon)

	q := `
		SELECT ST_AsText(geom)
		FROM moscow_region
		WHERE ST_Within(
		    geom, 
		    ST_SetSRID(ST_GeomFromText($1), 4326)
		)
		LIMIT 10000;
	`

	rows, err := s.db.Query(ctx, q, polygonWKT)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("can't querry points in polygon %w\n", err)
	}

	return translateRowsPoint(rows)
}

func (s *Storage) GetInRadiusPolygon(ctx context.Context, p internal.Polygon, radius int) ([]internal.Polygon, error) {

	centralPoint, newRadius := genpoint.GetCentralPolygonPointWithRadius(p.Vertical, radius)

	q := `
		SELECT ST_AsText(geom)
		FROM moscow_region_polygon
		WHERE ST_Within(
		geom,
		ST_Buffer(ST_SetSRID(ST_Point($1, $2), 4326),
		$3)
		)
		LIMIT 10000;
	`

	rows, err := s.db.Query(ctx, q, centralPoint.Longitude, centralPoint.Latitude, newRadius)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("can't get in radius %w\n", err)
	}

	return translateRowsPolygon(rows)
}

func (s *Storage) GetInPolygonPolygon(ctx context.Context, polygon internal.Polygon) ([]internal.Polygon, error) {
	polygonWKT := translatePolygonToWKT(polygon.Vertical)

	q := `
		SELECT ST_AsText(geom)
		FROM moscow_region_polygon
		WHERE ST_Within(
		    ST_SetSRID(ST_GeomFromText($1), 4326),
		    geom
		)
		LIMIT 10000;
	`

	rows, err := s.db.Query(ctx, q, polygonWKT)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("can't querry points in polygon %w\n", err)
	}

	return translateRowsPolygon(rows)
}

func (s *Storage) GetIntersectionPolygon(ctx context.Context, polygon internal.Polygon) ([]internal.Polygon, error) {
	polygonWKT := translatePolygonToWKT(polygon.Vertical)

	q := `
		SELECT ST_AsText(geom)
		FROM moscow_region_polygon
		WHERE ST_Intersects(
		    ST_SetSRID(ST_GeomFromText($1), 4326),
		    geom
		)
		LIMIT 10000;
	`

	rows, err := s.db.Query(ctx, q, polygonWKT)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("can't querry points in polygon %w\n", err)
	}

	return translateRowsPolygon(rows)
}

func (s *Storage) GetIntersectionPoint(ctx context.Context, point internal.Point) ([]internal.Polygon, error) {

	q := `
		SELECT ST_AsText(geom)
		FROM moscow_region_polygon
		WHERE ST_Contains(
		geom,
		(ST_SetSRID(ST_Point($1, $2), 4326))
	) 
		LIMIT 10000;
	`

	rows, err := s.db.Query(ctx, q, point.Longitude, point.Latitude)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("can't querry points in polygon %w\n", err)
	}

	return translateRowsPolygon(rows)
}

func (s *Storage) GetInShapes(ctx context.Context, shapes internal.Shapes) ([]internal.Point, error) {

	// Так делать плохо при тех данных, которые вводит пользователь, так как возможна sql инъекция
	// Здесь же я сам задаю то, что хочу передать и можно опустить неправильное использование переменных
	q := fmt.Sprintf("SELECT ST_AsText(geom) FROM moscow_region WHERE %s LIMIT 10000;", getStrShapesSQL(shapes))

	rows, err := s.db.Query(ctx, q)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("can't querry points in polygon %w\n", err)
	}

	return translateRowsPoint(rows)
}

func getStrShapesSQL(shapes internal.Shapes) string {

	lenPolygon := len(shapes.Polygons)
	lenCircle := len(shapes.Circles)

	queries := make([]string, lenPolygon+lenCircle)

	for i := 0; i < lenPolygon; i++ {
		queries[i] = convertPolygonToSQL(shapes.Polygons[i])
	}

	for i := 0; i < lenCircle; i++ {
		queries[i+lenPolygon] = convertCircleToSQL(shapes.Circles[i])
	}

	return strings.Join(queries, " OR ")
}

func convertPolygonToSQL(polygon internal.Polygon) string {
	return fmt.Sprintf("ST_Intersects(geom, (ST_SetSRID(ST_GeomFromText('%s'), 4326)))",
		translatePolygonToWKT(polygon.Vertical))
}

func convertCircleToSQL(circle internal.Circle) string {
	return fmt.Sprintf("ST_Intersects(geom, ST_Buffer(ST_SetSRID(ST_MakePoint(%f, %f), 4326), %d))",
		circle.Centre.Longitude, circle.Centre.Latitude, circle.Radius)
}

func translateRowsPolygon(rows pgx.Rows) ([]internal.Polygon, error) {
	nearestPolygons := make([]internal.Polygon, 0)

	for rows.Next() {
		var geometryData string
		err := rows.Scan(&geometryData)
		if err != nil {
			return nil, fmt.Errorf("can't scan row in nearestPoint %w\n", err)
		}

		geometryData = strings.TrimLeft(geometryData, "POLYGON(")
		geometryData = strings.TrimRight(geometryData, ")")
		polygonPoints := strings.Split(geometryData, ",")

		polygon := make([]internal.Point, len(polygonPoints))

		for i, p := range polygonPoints {
			polygon[i] = convertPointStrToFloat(strings.Split(p, " "))
		}

		nearestPolygons = append(nearestPolygons, internal.Polygon{Vertical: polygon})

	}

	return nearestPolygons, nil
}

func translateRowsPoint(rows pgx.Rows) ([]internal.Point, error) {

	nearestPoints := make([]internal.Point, 0)

	for rows.Next() {
		var geometryData string
		err := rows.Scan(&geometryData)
		if err != nil {
			return nil, fmt.Errorf("can't scan row in nearestPoint %w\n", err)
		}

		geometryData = strings.TrimLeft(geometryData, "POINT(")
		geometryData = strings.TrimRight(geometryData, ")")
		coordinates := strings.Split(geometryData, " ")

		nearestPoints = append(nearestPoints, convertPointStrToFloat(coordinates))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows %w\n", err)
	}

	return nearestPoints, nil
}

func convertPointStrToFloat(coordinates []string) internal.Point {
	if len(coordinates) != 2 {
		log.Printf("too much argument\n")
		return internal.Point{}
	}

	latitude, err1 := strconv.ParseFloat(coordinates[1], 64)
	if err1 != nil {
		log.Printf("too much argument\n")
		return internal.Point{}
	}
	longitude, err1 := strconv.ParseFloat(coordinates[0], 64)
	if err1 != nil {
		log.Printf("too much argument\n")
		return internal.Point{}
	}

	return internal.Point{Latitude: latitude, Longitude: longitude}
}
