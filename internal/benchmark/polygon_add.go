package benchmark

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
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

func runPolygonBenchDBInitAndAdd(ctx context.Context, s storage.PolygonStorage, bf *BenchFile) error {

	polyGen := genpoint.PolygonGenerator{}
	polygons := polyGen.GeneratePolygons(bf.CountPolygonAdd)

	_, err := benchDropPolygon(ctx, s)
	if err != nil {
		return err
	}

	dur, err := benchInitPolygon(ctx, s)
	if err != nil {
		return err
	}
	bf.Durations[PolygonInit] += dur

	dur, err = benchDropPolygon(ctx, s)
	if err != nil {
		return err
	}
	bf.Durations[PolygonDrop] += dur

	_, err = benchInitPolygon(ctx, s)
	if err != nil {
		return err
	}

	dur, err = benchAddPolygonBatch(ctx, s, polygons)
	if err != nil {
		return err
	}
	bf.Durations[PolygonAddBatch] += dur

	return nil
}
