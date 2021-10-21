package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"postapi/app"
	"postapi/app/database"

	// "postapi/app/otl"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

func main() {
	ctx := context.Background()

	driverOpts := []otlpgrpc.Option{
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(fmt.Sprintf("%s:%d", "127.0.0.1", 4317)),
	}
	driver := otlpgrpc.NewDriver(
		driverOpts...,
	)

	exporter, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		log.Fatalf("could not start otl exporter: %v", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String("posts-api-dev")),
	)

	if err != nil {
		log.Fatalf("could not create otl resource: %v", err)
	}

	// excludedRoutes := []string{"/health [GET]"}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		// trace.WithSampler(otl.CustomSampler{ExcludedRoutes: excludedRoutes, Desc: "customSampler"}),
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
		trace.WithResource(res),
	)

	defer func() { _ = tp.Shutdown(ctx) }()

	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)
	otel.SetTracerProvider(tp)

	app := app.New()
	app.DB = &database.DB{}
	err = app.DB.Open()
	check(err)

	defer app.DB.Close()

	r := app.Router

	s := http.Server{
		Addr: "localhost:9000",
		Handler: otelhttp.NewHandler(r, "request",
			otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
				return fmt.Sprintf("%s [%s]", r.URL.Path, r.Method)
			}),
		),
	}

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("could not start server: %v", err)
	}

	log.Println("App running..")
	err = http.ListenAndServe(":9000", nil)
	check(err)
}

func check(e error) {
	if e != nil {
		log.Println(e)
		os.Exit(1)
	}
}
