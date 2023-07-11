package app

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/postgis"
	"log"
)

func Run(dbPG *postgis.Storage) {

	log.Printf("app start\n")

	ctx := context.Background()

	err := dbPG.Init(ctx)

	if err != nil {
		log.Fatalf("can't init postgis %e\n", err)
	}

	// тестовое добавление
	// потом буду брать из генератора

	generator := genpoint.SimplePointGenerator{}

	points := generator.GeneratePoints(5)

	err = dbPG.AddPointBatch(ctx, points)

	if err != nil {
		fmt.Println("can`t add batch", err)
	}

	fmt.Println("end")

	// здесь будет происходить запуск различных бенчмарков
}
