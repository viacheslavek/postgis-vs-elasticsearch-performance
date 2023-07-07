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

	err = dbPG.AddPoint(ctx, storage.Point{Latitude: 51.5199, Longitude: -0.1238})
	if err != nil {
		log.Printf("can't add to dbPG %e\n", err)
	}

	err = dbPG.AddPoint(ctx, storage.Point{Latitude: 51.5083, Longitude: -0.1278})
	if err != nil {
		log.Printf("can't add to dbPG %e\n", err)
	}

	err = dbPG.DeletePoint(ctx, storage.Point{Latitude: 51.5083, Longitude: -0.1278})
	if err != nil {
		log.Printf("can't delete point in dbPG %e\n", err)
	}

	is, err := dbPG.IsPoint(ctx, storage.Point{Latitude: 51.5199, Longitude: -0.1238})
	fmt.Println(is, "is")
	if err != nil {
		log.Printf("can't is dbPG %e\n", err)
	}

	notIs, err := dbPG.IsPoint(ctx, storage.Point{Latitude: 51.5083, Longitude: -0.1278})
	fmt.Println(notIs, "notIs")
	if err != nil {
		log.Printf("can't is dbPG %e\n", err)
	}

	err = dbPG.Drop(ctx)

	if err != nil {
		log.Printf("can't drop dbPG %e\n", err)
	}

	fmt.Println("end")

	// здесь будет происходить запуск различных бенчмарков
}
