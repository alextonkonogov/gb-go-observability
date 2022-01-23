package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	aConfig "github.com/alextonkonogov/gb-go-observability/homework3/app/internal/config"
	"github.com/alextonkonogov/gb-go-observability/homework3/app/internal/log"
	"github.com/alextonkonogov/gb-go-observability/homework3/app/internal/repository"
	"github.com/alextonkonogov/gb-go-observability/homework3/app/internal/storage"
	jTracer "github.com/alextonkonogov/gb-go-observability/homework3/app/internal/tracer"
)

func main() {
	logger := log.NewLogWithConfiguration()

	cnfg, err := aConfig.NewAppConfig()
	if err != nil {
		logger.WithError(err).Fatal()
	}

	tracer, closer, err := jTracer.InitJaeger("motivation", logger)
	defer closer.Close()

	ctx := context.Background()
	dbpool, err := storage.InitDBConn(ctx, cnfg)
	if err != nil {
		logger.WithError(err).Fatal()
	}
	defer dbpool.Close()

	err = storage.InitTables(ctx, dbpool)
	if err != nil {
		logger.WithError(err).Fatal()
	}

	a := newApp(ctx, dbpool, logger, tracer)

	r := mux.NewRouter()
	r.Path("/").HandlerFunc(a.startPage)

	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: r}
	srv.ListenAndServe()
}

type app struct {
	ctx        context.Context
	dbpool     *pgxpool.Pool
	logger     *logrus.Logger
	tracer     opentracing.Tracer
	repository *repository.Repository
}

func newApp(ctx context.Context, dbpool *pgxpool.Pool, logger *logrus.Logger, tracer opentracing.Tracer) *app {
	return &app{ctx, dbpool, logger, tracer, repository.NewRepository(dbpool, tracer)}
}

func (a app) startPage(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(r.Context(), a.tracer, "startPageHandler")
	defer span.Finish()

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
