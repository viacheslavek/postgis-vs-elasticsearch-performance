package benchmark

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
	"log"
	"time"
)

type BenchFile struct {
	Durations              map[int]time.Duration
	DurationPointInPolygon []time.Duration
	dbName                 string
	countPoints            int
	countPolygon           int
	countShapes            int
	countChecks            int
}

func RunBenchNCheck(ctx context.Context, ss storage.Storage, sp storage.PolygonStorage, db string,
	countPoints, countPolygon, countShapes, countChecks int) (BenchFile, error) {

	bf := BenchFile{
		Durations:              make(map[int]time.Duration),
		DurationPointInPolygon: make([]time.Duration, countPolygon),
		dbName:                 db,
		countPoints:            countPoints,
		countPolygon:           countPolygon,
		countShapes:            countShapes,
		countChecks:            countChecks,
	}

	for i := 0; i < countChecks; i++ {
		log.Println("step number:", i)
		if err := RunBenchPoint(ctx, ss, &bf); err != nil {
			return bf, fmt.Errorf("can't do checks in bench point %w", err)
		}
		if err := RunBenchPolygon(ctx, sp, &bf); err != nil {
			return bf, fmt.Errorf("can't do checks in bench polygon %w", err)
		}
	}

	return bf, nil

}

func RunBenchPoint(ctx context.Context, s storage.Storage, bf *BenchFile) error {

	if err := runPointBenchDBInitAndAdd(ctx, s, bf); err != nil {
		return fmt.Errorf("can't run init banch: %w\n", err)
	}

	if err := runBenchPointSearch(ctx, s, bf); err != nil {
		return fmt.Errorf("can't run search banch: %w\n", err)
	}

	return nil
}

func RunBenchPolygon(ctx context.Context, s storage.PolygonStorage, bf *BenchFile) error {

	if err := runPolygonBenchDBInitAndAdd(ctx, s, bf); err != nil {
		return fmt.Errorf("can't run init and add polygon banch: %w\n", err)
	}

	if err := runBenchPolygonSearch(ctx, s, bf); err != nil {
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
