package benchmark

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
	"log"
	"time"
)

func RunBench(ctx context.Context, s storage.Storage, db string, N int) error {
	fmt.Println("go to the logs")

	start := time.Now()
	pointG := genpoint.SimplePointGenerator{}
	points := pointG.GeneratePoints(N)
	log.Println("generate points: ", time.Since(start))

	if err := runBenchDBInitAndAdd(ctx, s, db, points); err != nil {
		return fmt.Errorf("can't run init banch: %w\n", err)
	}

	// TODO: для поиска в радиусе

	// TODO: для поиска в многоугольниках (от 1 до N)

	return nil

}
