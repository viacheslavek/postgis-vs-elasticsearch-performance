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

func benchInitPolygon(ctx context.Context, s storage.PolygonStorage) (time.Duration, error) {
	start := time.Now()

	if err := s.InitPolygon(ctx); err != nil {
		return 0, fmt.Errorf("can't init polygon bench db %w\n", err)
	}

	return time.Since(start), nil
}

func benchDropPolygon(ctx context.Context, s storage.PolygonStorage) (time.Duration, error) {
	start := time.Now()

	if err := s.DropPolygon(ctx); err != nil {
		return 0, fmt.Errorf("can't drop polygon bench db %w\n", err)
	}

	endBench := time.Since(start)

	return endBench, nil
}

func benchAddPolygon(ctx context.Context, s storage.PolygonStorage, polygons []internal.Polygon) (time.Duration, error) {
	start := time.Now()

	for _, p := range polygons {
		if err := s.AddPolygon(ctx, p); err != nil {
			return 0, fmt.Errorf("can't add single polygon bench db %w\n", err)
		}
	}

	endBench := time.Since(start)

	return endBench, nil
}

func benchAddPolygonBatch(ctx context.Context, s storage.PolygonStorage, polygons []internal.Polygon) (time.Duration, error) {
	start := time.Now()

	if err := s.AddPolygonBatch(ctx, polygons); err != nil {
		return 0, fmt.Errorf("can't add polygon batch bench db %w\n", err)
	}

	endBench := time.Since(start)

	return endBench, nil
}

func runPolygonBenchDBInitAndAdd(ctx context.Context, s storage.PolygonStorage, db string, countPolygons int) error {

	polyGen := genpoint.PolygonGenerator{}
	polygons := polyGen.GeneratePolygons(countPolygons)

	_, err := benchDropPolygon(ctx, s)
	if err != nil {
		return err
	}

	log.Printf("testing polygon db: %s\n", db)

	dur, err := benchInitPolygon(ctx, s)
	if err != nil {
		return err
	}
	log.Printf("time to Init: %s", dur.String())

	dur, err = benchAddPolygon(ctx, s, polygons)
	if err != nil {
		return err
	}
	log.Printf("time to Add: %s", dur.String())

	dur, err = benchDropPolygon(ctx, s)
	if err != nil {
		return err
	}
	log.Printf("time to Drop: %s", dur.String())

	_, err = benchInitPolygon(ctx, s)
	if err != nil {
		return err
	}

	dur, err = benchAddPolygonBatch(ctx, s, polygons)
	if err != nil {
		return err
	}
	log.Printf("time to Add batch: %s", dur.String())

	return nil
}
