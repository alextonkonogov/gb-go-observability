package main

import (
	"context"
	"fmt"
	"go.uber.org/zap/zapcore"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	aConfig "github.com/alextonkonogov/gb-go-observability/homework3/app/internal/config"
	"github.com/alextonkonogov/gb-go-observability/homework3/app/internal/repository"
	"github.com/alextonkonogov/gb-go-observability/homework3/app/internal/storage"
	jTracer "github.com/alextonkonogov/gb-go-observability/homework3/app/internal/tracer"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = logger.Sync() }()

	cnfg, err := aConfig.NewAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	tracer, closer := jTracer.InitJaeger("motivation", logger)
	defer closer.Close()

	ctx := context.Background()
	dbpool, err := storage.InitDBConn(ctx, cnfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	err = storage.InitTables(ctx, dbpool)
	if err != nil {
		log.Fatal(err)
	}

	a := newApp(ctx, dbpool, logger, tracer)

	r := mux.NewRouter()
	r.Path("/").HandlerFunc(a.startPage)

	http.ListenAndServe("0.0.0.0:8080", nethttp.Middleware(a.tracer, r))
}

type app struct {
	ctx        context.Context
	dbpool     *pgxpool.Pool
	logger     *zap.Logger
	tracer     opentracing.Tracer
	repository *repository.Repository
}

func newApp(ctx context.Context, dbpool *pgxpool.Pool, logger *zap.Logger, tracer opentracing.Tracer) *app {
	return &app{ctx, dbpool, logger, tracer, repository.NewRepository(dbpool, tracer)}
}

func (a app) startPage(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(r.Context(), a.tracer, "startPageHandler")
	defer span.Finish()

	a.logger.Info("motivationHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})

	motivation, err := a.repository.GetRandomMotivation(ctx, a.dbpool)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	content := fmt.Sprintf("\"%s\" %s\n", motivation.Content, motivation.Author)
	tmpl, err := template.New("example").Parse(content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tmpl.Execute(w, content)
}
