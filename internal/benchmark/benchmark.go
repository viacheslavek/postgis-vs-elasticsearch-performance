package benchmark

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
	"log"
	"time"
)

// TODO: перенести генерацию точек и внутренних вещей в самые нижние функции, чтобы не было зависимостей на всех уровнях
// А так же радиус перенести из главной функции в низ

func RunBench(ctx context.Context, s storage.Storage, db string,
	countPoints, countPolygon, radius, countShapes int) error {
	fmt.Println("go to the logs")

	start := time.Now()
	pointG := genpoint.SimplePointGenerator{}
	points := pointG.GeneratePoints(countPoints)
	log.Println("generate points: ", time.Since(start))

	if err := runBenchDBInitAndAdd(ctx, s, db, points); err != nil {
		return fmt.Errorf("can't run init banch: %w\n", err)
	}

	point := pointG.GeneratePoints(1)

	if err := runBenchSearch(ctx, s, db, point[0], radius, countPolygon, countShapes); err != nil {
		return fmt.Errorf("can't run search banch: %w\n", err)
	}

	return nil

}
