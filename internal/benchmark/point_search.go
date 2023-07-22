package benchmark

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
	"log"
	"time"
)

func benchSearchInRadius(ctx context.Context, s storage.Storage, p internal.Point, radius int) (time.Duration, error) {
	start := time.Now()

	_, err := s.GetInRadius(ctx, p, radius)

	if err != nil {
		return 0, fmt.Errorf("can't search radius in bench db %w\n", err)
	}

	endBench := time.Since(start)

	return endBench, nil
}

func benchSearchInPolygon(ctx context.Context, s storage.Storage, countPolygon int) ([]time.Duration, error) {

	durations := make([]time.Duration, countPolygon)
	genPolygon := genpoint.PolygonGenerator{}

	for i := 3; i < countPolygon; i++ {
		polygon := genPolygon.GeneratePolygon(i)
		dur, err := getInPolygon(ctx, s, polygon)
		if err != nil {
			return nil, fmt.Errorf("can't search radius in bench db %w\n", err)
		}
		durations = append(durations, dur)
	}

	return durations[3:], nil
}

func getInPolygon(ctx context.Context, s storage.Storage, polygon []internal.Point) (time.Duration, error) {
	start := time.Now()
	_, err := s.GetInPolygon(ctx, polygon)

	if err != nil {
		return 0, fmt.Errorf("can't search radius in bench db %w\n", err)
	}

	endBench := time.Since(start)

	return endBench, nil
}

func benchSearchInShapes(ctx context.Context, s storage.Storage, countShapes int) (time.Duration, error) {

	shapes := genpoint.GenerateShapes(countShapes)

	start := time.Now()

	_, err := s.GetInShapes(ctx, shapes)

	if err != nil {
		return 0, fmt.Errorf("can't search radius in bench db %w\n", err)
	}

	endBench := time.Since(start)

	return endBench, nil
}

func runBenchSearch(ctx context.Context, s storage.Storage, db string,
	p internal.Point, radius, countPolygon, countShapes int) error {

	log.Printf("testing db: %s\n", db)

	dur, err := benchSearchInRadius(ctx, s, p, radius)
	if err != nil {
		return err
	}
	log.Printf("time to search in radius: %s", dur.String())

	durs, err := benchSearchInPolygon(ctx, s, countPolygon)
	if err != nil {
		return err
	}
	log.Printf("time to search in polygon: %s", dur.String())

	for i, d := range durs {
		log.Println("count:", i, "-", d)
	}

	dur, err = benchSearchInShapes(ctx, s, countShapes)
	if err != nil {
		return err
	}
	log.Printf("time to search in radius: %s", dur.String())

	return nil
}
