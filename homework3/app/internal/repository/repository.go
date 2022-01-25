package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const MotivationSelect = `select * from motivations order by random() limit 1`

type Repository struct {
	pool   *pgxpool.Pool
	tracer opentracing.Tracer
}

func (r *Repository) GetRandomMotivation(ctx context.Context, dbpool *pgxpool.Pool) (m motivation, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "Repository.GetRandomMotivation")
	defer span.Finish()
	span.LogFields(
		log.String("query", MotivationSelect),
	)

	row := dbpool.QueryRow(ctx, MotivationSelect)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	err = row.Scan(&m.Id, &m.Content, &m.Author, &m.UserId)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	return
}

func NewRepository(pool *pgxpool.Pool, tracer opentracing.Tracer) *Repository {
	return &Repository{pool: pool, tracer: tracer}
}
