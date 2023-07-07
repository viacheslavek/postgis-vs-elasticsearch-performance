package app

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
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

	points := []storage.Point{
		{Latitude: 51.5074, Longitude: -1.1278},
		{Latitude: 51.5199, Longitude: -2.1238},
		{Latitude: 51.5083, Longitude: -3.1278},
	}

	err = dbPG.AddPointBatch(ctx, points)

	if err != nil {
		fmt.Println("can`t add batch", err)
	}

	fmt.Println("end")

	// здесь будет происходить запуск различных бенчмарков
}
