package benchmark

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
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

func runBenchDBInitAndAdd(ctx context.Context, s storage.Storage, db string, ps []internal.Point) error {

	_, err := benchDropPoint(ctx, s)
	if err != nil {
		return err
	}

	log.Printf("testing db: %s\n", db)
	dur, err := benchInitPoint(ctx, s)
	if err != nil {
		return err
	}
	log.Printf("time to Init: %s", dur.String())

	dur, err = benchAddPoint(ctx, s, ps)
	if err != nil {
		return err
	}
	log.Printf("time to Add: %s", dur.String())

	dur, err = benchDropPoint(ctx, s)
	if err != nil {
		return err
	}
	log.Printf("time to Drop: %s", dur.String())

	dur, err = benchInitPoint(ctx, s)
	if err != nil {
		return err
	}
	log.Printf("time to Init: %s", dur.String())

	dur, err = benchAddPointBatch(ctx, s, ps)
	if err != nil {
		return err
	}
	log.Printf("time to Add : %s", dur.String())

	return nil
}
