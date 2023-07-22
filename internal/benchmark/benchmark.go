package benchmark

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
)

func RunBenchPoint(ctx context.Context, s storage.Storage, db string,
	countPoints, countPolygon, countShapes int) error {
	fmt.Println("go to the logs")

	if err := runPointBenchDBInitAndAdd(ctx, s, db, countPoints); err != nil {
		return fmt.Errorf("can't run init banch: %w\n", err)
	}

	if err := runBenchPointSearch(ctx, s, db, countPolygon, countShapes); err != nil {
		return fmt.Errorf("can't run search banch: %w\n", err)
	}

	return nil
}

func RunBenchPolygon(ctx context.Context, s storage.PolygonStorage, db string, countPolygon int) error {

	if err := runPolygonBenchDBInitAndAdd(ctx, s, db, countPolygon); err != nil {
		return fmt.Errorf("can't run init and add polygon banch: %w\n", err)
	}

	if err := runBenchPolygonSearch(ctx, s, db); err != nil {
		return fmt.Errorf("can't run init and add polygon banch: %w\n", err)
	}

	return nil
}
