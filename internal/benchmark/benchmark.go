package benchmark

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
)

// TODO: перенести генерацию точек и внутренних вещей в низ лежащие функции, чтобы не было зависимостей на всех уровнях

func RunBenchPoint(ctx context.Context, s storage.Storage, db string,
	countPoints, countPolygon, countShapes int) error {
	fmt.Println("go to the logs")

	pointG := genpoint.SimplePointGenerator{}

	if err := runPointBenchDBInitAndAdd(ctx, s, db, countPoints); err != nil {
		return fmt.Errorf("can't run init banch: %w\n", err)
	}

	point := pointG.GeneratePoints(1)

	radius := 1000

	if err := runBenchSearch(ctx, s, db, point[0], radius, countPolygon, countShapes); err != nil {
		return fmt.Errorf("can't run search banch: %w\n", err)
	}

	return nil
}

func RunBenchPolygon(ctx context.Context, s storage.PolygonStorage, db string, countPolygon int) error {

	if err := runPolygonBenchDBInitAndAdd(ctx, s, db, countPolygon); err != nil {
		return fmt.Errorf("can't run init and add polygon banch: %w\n", err)
	}

	return nil
}
