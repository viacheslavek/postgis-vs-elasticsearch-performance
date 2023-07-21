package app

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/benchmark"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/elasticsearch"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/postgis"
	"log"
)

func Run(pg *postgis.Storage, es *elasticsearch.Storage) {

	log.Printf("app start\n")

	ctx := context.Background()

	if err := benchmark.RunBench(ctx, pg, "postgis", 20e1, 5, 10000); err != nil {
		log.Fatalf("can't do postgis bench %e\n", err)
	}

	if err := benchmark.RunBench(ctx, es, "elasticseacrh", 20e1, 5, 10000); err != nil {
		log.Fatalf("can't do es bench %e\n", err)
	}

	polygonGen := genpoint.PolygonGenerator{}

	polygons := polygonGen.GeneratePolygons(2e4)

	err := pg.DropPolygon(ctx)
	if err != nil {
		log.Fatalf("can't init postgis polygon %e\n", err)
	}

	err = pg.InitPolygon(ctx)
	if err != nil {
		log.Fatalf("can't init postgis polygon %e\n", err)
	}

	err = pg.AddPolygonBatch(ctx, polygons)
	if err != nil {
		log.Fatalf("can't add postgis polygon %e\n", err)
	}

	err = pg.AddPolygon(ctx, polygons[0])
	if err != nil {
		log.Fatalf("can't add postgis polygon %e\n", err)
	}

	fmt.Println("all add")

	_, err = pg.GetInPolygonPolygon(ctx, polygons[0].Vertical)
	if err != nil {
		log.Fatalf("can't add postgis polygon %e\n", err)
	}

	_, err = pg.GetInRadiusPolygon(ctx, polygons[0], 10)
	if err != nil {
		log.Fatalf("can't add postgis polygon %e\n", err)
	}

	log.Printf("end\n")

}

//curl -X GET "http://localhost:9200/moscow_region/_search" -H "Content-Type: application/json" -d'
//{
// "query": {
//   "match_all": {}
// },
// "size": 100
//}
//'
