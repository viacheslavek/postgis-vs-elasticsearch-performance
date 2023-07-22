package benchmark

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
	"log"
	"math/rand"
	"time"
)

func benchGetInRadiusPolygon(ctx context.Context, s storage.PolygonStorage,
	polygon internal.Polygon, radius int) (time.Duration, error) {
	start := time.Now()

	_, err := s.GetInRadiusPolygon(ctx, polygon, radius)

	if err != nil {
		return 0, fmt.Errorf("can't search radius in polygon bench db %w\n", err)
	}

	endBench := time.Since(start)

	return endBench, nil
}

func benchGetInPolygonPolygon(ctx context.Context, s storage.PolygonStorage,
	polygon internal.Polygon) (time.Duration, error) {
	start := time.Now()

	_, err := s.GetInPolygonPolygon(ctx, polygon)

	if err != nil {
		return 0, fmt.Errorf("can't search polygon in polygon bench db %w\n", err)
	}

	endBench := time.Since(start)

	return endBench, nil
}

func benchGetIntersectionPolygon(ctx context.Context, s storage.PolygonStorage,
	polygon internal.Polygon) (time.Duration, error) {

	start := time.Now()

	_, err := s.GetIntersectionPolygon(ctx, polygon)

	if err != nil {
		return 0, fmt.Errorf("can't search intersection polygon in polygon bench db %w\n", err)
	}

	endBench := time.Since(start)

	return endBench, nil
}

func benchGetIntersectionPoint(ctx context.Context, s storage.PolygonStorage,
	point internal.Point) (time.Duration, error) {

	start := time.Now()

	_, err := s.GetIntersectionPoint(ctx, point)

	if err != nil {
		return 0, fmt.Errorf("can't search intersection point in polygon bench db %w\n", err)
	}

	endBench := time.Since(start)

	return endBench, nil
}

func runBenchPolygonSearch(ctx context.Context, s storage.PolygonStorage, db string) error {

	radius := rand.Intn(1e5)

	polyGen := genpoint.PolygonGenerator{}

	polygon := internal.Polygon{Vertical: polyGen.GeneratePolygon(rand.Intn(20) + 3)}

	log.Printf("testing polygon search db: %s\n", db)

	dur, err := benchGetInRadiusPolygon(ctx, s, polygon, radius)
	if err != nil {
		return err
	}
	log.Printf("time to search in polygon radius: %s", dur.String())

	dur, err = benchGetInPolygonPolygon(ctx, s, polygon)
	if err != nil {
		return err
	}
	log.Printf("time to search in polygon polygon: %s", dur.String())

	dur, err = benchGetIntersectionPolygon(ctx, s, polygon)
	if err != nil {
		return err
	}
	log.Printf("time to search intersection polygon: %s", dur.String())

	spg := genpoint.SimplePointGenerator{}

	point := spg.GeneratePoints(1)[0]

	dur, err = benchGetIntersectionPoint(ctx, s, point)
	if err != nil {
		return err
	}
	log.Printf("time to search intersection point in polygon: %s", dur.String())

	return nil
}
