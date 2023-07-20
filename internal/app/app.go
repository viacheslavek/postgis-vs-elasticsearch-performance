package app

import (
	"context"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/elasticsearch"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/postgis"
	"log"
)

func Run(pg *postgis.Storage, es *elasticsearch.Storage) {

	log.Printf("app start\n")
	ctx := context.Background()

	//if err := benchmark.RunBench(ctx, pg, "postgis", 20e1, 5, 10000); err != nil {
	//	log.Fatalf("can't do postgis bench %e\n", err)
	//}
	//
	//if err := benchmark.RunBench(ctx, es, "elasticseacrh", 20e1, 5, 10000); err != nil {
	//	log.Fatalf("can't do es bench %e\n", err)
	//}

	err := es.DropPolygon(ctx)
	if err != nil {
		log.Fatalf("%e", err)
	}

	err = es.InitPolygon(ctx)
	if err != nil {
		log.Fatalf("%e", err)
	}

	pgen := genpoint.PolygonGenerator{}

	polygons := pgen.GeneratePolygons(5)

	err = es.AddPolygonBatch(ctx, polygons)
	if err != nil {
		log.Fatalf("opps %e", err)
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
