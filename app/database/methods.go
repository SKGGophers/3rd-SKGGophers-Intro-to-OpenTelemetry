package database

import (
	"context"
	"postapi/app/models"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (d *DB) CreatePost(ctx context.Context, p *models.Post) error {
	span := trace.SpanFromContext(ctx)
	_, sp := span.Tracer().Start(ctx, "Insert [PG]")
	defer sp.End()

	configureSpan(sp, "INSERT", insertPostSchema)

	res, err := d.db.Exec(insertPostSchema, p.Title, p.Content, p.Author)
	if err != nil {
		return err
	}

	res.LastInsertId()
	return err
}

func (d *DB) GetPosts(ctx context.Context) ([]*models.Post, error) {
	span := trace.SpanFromContext(ctx)
	_, sp := span.Tracer().Start(ctx, "Select [PG]")
	defer sp.End()

	configureSpan(sp, "SELECT", selectPostsSchema)

	var posts []*models.Post
	err := d.db.Select(&posts, selectPostsSchema)
	if err != nil {
		return posts, err
	}

	return posts, nil
}

func configureSpan(s trace.Span, dbOp string, statement string) {
	// https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/trace/semantic_conventions
	s.SetAttributes(
		attribute.KeyValue{Key: "db.system", Value: attribute.StringValue("postgresql")},
		attribute.KeyValue{Key: "db.operation", Value: attribute.StringValue(dbOp)},
		attribute.KeyValue{Key: "db.statement", Value: attribute.StringValue(statement)},
		attribute.KeyValue{Key: "span.kind", Value: attribute.StringValue("CLIENT")},
	)

}
