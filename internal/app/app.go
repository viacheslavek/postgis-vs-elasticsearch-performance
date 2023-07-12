package app

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
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

	// генерирую точки
	generator := genpoint.SimplePointGenerator{}
	points := generator.GeneratePoints(5)

	// добавляю в PostGis
	err = dbPG.AddPointBatch(ctx, points)
	if err != nil {
		fmt.Println("can`t add batch", err)
	}

	// тестовая генерация многоугольника
	// заменить на бейнчмарки
	pg := genpoint.PolygonGenerator{}
	fmt.Println(pg.GeneratePolygon(5))

	pointsGet, err := dbPG.GetInRadius(ctx, internal.Point{Latitude: 51.5074, Longitude: -1.1278}, 10e15)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(pointsGet)

	fmt.Println("end")

}
