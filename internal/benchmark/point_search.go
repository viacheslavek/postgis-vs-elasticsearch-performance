package benchmark

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
	"math/rand"
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

	durations := make([]time.Duration, 0)
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

func runBenchPointSearch(ctx context.Context, s storage.Storage, bf *BenchFile) error {

	spg := genpoint.SimplePointGenerator{}

	point := spg.GeneratePoints(1)[0]

	radius := rand.Intn(1e5)

	dur, err := benchSearchInRadius(ctx, s, point, radius)
	if err != nil {
		return err
	}
	bf.Durations[PointSearchInRadius] += dur

	durs, err := benchSearchInPolygon(ctx, s, bf.countPolygon)
	if err != nil {
		return err
	}
	bf.DurationPointInPolygon = durs

	dur, err = benchSearchInShapes(ctx, s, bf.countShapes)
	if err != nil {
		return err
	}
	bf.Durations[PointSearchInShapes] += dur

	return nil
}
