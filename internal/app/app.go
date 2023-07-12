package app

import (
	"context"
	"fmt"
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

	fmt.Println("опять я не убрал закоменченный код в прошлый раз")

	fmt.Println("начинаю работу с эластиком")

	fmt.Println("end")

}
