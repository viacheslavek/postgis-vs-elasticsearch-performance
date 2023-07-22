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

func benchInitPoint(ctx context.Context, s storage.Storage) (time.Duration, error) {
	start := time.Now()

	if err := s.Init(ctx); err != nil {
		return 0, fmt.Errorf("can't init bench db %w\n", err)
	}

	return time.Since(start), nil
}

func benchDropPoint(ctx context.Context, s storage.Storage) (time.Duration, error) {
	start := time.Now()

	if err := s.Drop(ctx); err != nil {
		return 0, fmt.Errorf("can't drop bench db %w\n", err)
	}

	endBench := time.Since(start)

	return endBench, nil
}

func benchAddPoint(ctx context.Context, s storage.Storage, ps []internal.Point) (time.Duration, error) {
	start := time.Now()

	for _, p := range ps {
		if err := s.AddPoint(ctx, p); err != nil {
			return 0, fmt.Errorf("can't add single bench db %w\n", err)
		}
	}

	endBench := time.Since(start)

	return endBench, nil
}

func benchAddPointBatch(ctx context.Context, s storage.Storage, ps []internal.Point) (time.Duration, error) {
	start := time.Now()

	if err := s.AddPointBatch(ctx, ps); err != nil {
		return 0, fmt.Errorf("can't add batch bench db %w\n", err)
	}

	endBench := time.Since(start)

	return endBench, nil
}

func runPointBenchDBInitAndAdd(ctx context.Context, s storage.Storage, db string, countPoints int) error {

	log.Printf("testing point db: %s\n", db)

	start := time.Now()
	pointGen := genpoint.SimplePointGenerator{}

	points := pointGen.GeneratePoints(countPoints)
	log.Println("generate points: ", time.Since(start))

	dur, err := benchDropPoint(ctx, s)
	if err != nil {
		return err
	}
	log.Printf("time to Drop: %s", dur.String())

	_, err = benchInitPoint(ctx, s)
	if err != nil {
		return err
	}

	dur, err = benchAddPointBatch(ctx, s, points)
	if err != nil {
		return err
	}
	log.Printf("time to Add batch: %s", dur.String())

	return nil
}
