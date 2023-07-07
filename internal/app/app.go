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

	err = dbPG.AddPoint(ctx, storage.Point{Latitude: 51.5074, Longitude: -0.1278})
	if err != nil {
		log.Printf("can't add to dbPG %e\n", err)
	}

	fmt.Println("end")

	// здесь будет происходить запуск различных бенчмарков
}
