package benchmark

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
	"log"
)

// TODO: сделать мапку с общим временем рабботы бенчмарков и сделать более общую функцию, которая позволит повторять
// запуск бенчмарка N раз и больше - для точных операций
// Завести для этого набор констант

// TODO: замер потраченной памяти

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

	if err := runPolygonBenchDBInitAndAdd(ctx, s, db, countPolygon/100); err != nil {
		return fmt.Errorf("can't run init and add polygon banch: %w\n", err)
	}

	if err := runBenchPolygonSearch(ctx, s, db); err != nil {
		return fmt.Errorf("can't run init and add polygon banch: %w\n", err)
	}

	return nil
}

func HowBadAddSinglePoint(ctx context.Context, s storage.Storage, db string, countPoints int) error {
	_, err := benchDropPoint(ctx, s)
	if err != nil {
		return err
	}

	log.Printf("testing point db: %s\n", db)
	dur, err := benchInitPoint(ctx, s)
	if err != nil {
		return err
	}
	log.Printf("time to Init: %s", dur.String())

	pointGen := genpoint.SimplePointGenerator{}

	points := pointGen.GeneratePoints(countPoints)

	dur, err = benchAddPoint(ctx, s, points)
	if err != nil {
		return err
	}
	log.Printf("time to Add: %s", dur.String())

	return nil
}
