package postgis

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"github.com/jackc/pgx/v5"
	"strings"
)

func (s *Storage) InitPolygon(ctx context.Context) error {
	q := `
		CREATE EXTENSION IF NOT EXISTS postgis;
		CREATE TABLE IF NOT EXISTS moscow_region_polygon (
    		id SERIAL PRIMARY KEY,
    		geom GEOMETRY(Polygon, 4326)
		);
	`

	return s.initBase(ctx, q)
}

func (s *Storage) DropPolygon(ctx context.Context) error {
	q := `
		DROP TABLE IF EXISTS moscow_region_polygon;
	`

	return s.drop(ctx, q)
}

func (s *Storage) AddPolygon(ctx context.Context, polygon internal.Polygon) error {

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return fmt.Errorf("can't begin transaction: %w\n", err)
	}

	q := `
		INSERT INTO moscow_region_polygon (geom) 
			VALUES (ST_SetSRID(ST_GeomFromText($1), 4326))
	`

	_, err = tx.Exec(ctx, q, translatePolygonToWKT(polygon.Vertical))

	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("can't add a point: %w\n", err)
	}

	err = tx.Commit(ctx)

	if err != nil {
		return fmt.Errorf("can't commit transactions %w\n", err)
	}

	return nil
}

func (s *Storage) AddPolygonBatch(ctx context.Context, polygons []internal.Polygon) error {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("can`t acquire connection from pool: %w\n", err)
	}

	defer conn.Release()

	batch := &pgx.Batch{}

	q := `
		INSERT INTO moscow_region_polygon (geom) 
			VALUES (ST_SetSRID(ST_GeomFromText($1), 4326))
	`

	for _, p := range polygons {
		batch.Queue(q, translatePolygonToWKT(p.Vertical))
	}

	return s.addBatch(ctx, conn, batch, len(polygons))
}

func translatePolygonToWKT(polygon []internal.Point) string {
	wktPoints := make([]string, len(polygon))
	for i, p := range polygon {
		wktPoints[i] = fmt.Sprintf("%f %f", p.Longitude, p.Latitude)
	}

	return fmt.Sprintf("POLYGON((%s))", strings.Join(wktPoints, ", "))
}
