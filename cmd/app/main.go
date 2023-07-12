package main

import (
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/app"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/postgis"
	"log"
)

// инициализация логгера и конфигов. Дальше передаю управление другому пакету

func main() {

	log.Printf("app starts launching\n")

	postGis, err := postgis.New()

	if err != nil {
		log.Fatalf("can't connect to db %e\n", err)
	}

	app.Run(postGis)
}
