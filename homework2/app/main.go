package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/alextonkonogov/gb-go-observability/homework2/app/internal/config"
	"github.com/alextonkonogov/gb-go-observability/homework2/app/internal/motivations"
	"github.com/alextonkonogov/gb-go-observability/homework2/app/internal/storage"
)

func main() {
	log := logrus.New()
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(customFormatter)

	log.Info("Preparing configs...")
	cnfg, err := config.NewAppConfig()
	if err != nil {
		log.WithError(err).Fatal("failed to get configs")
	}
	log.Info("Configs are ready")

	ctx := context.Background()

	log.Info("Connecting to database...")
	dbpool, err := storage.InitDBConn(ctx, cnfg)
	if err != nil {
		log.WithError(err).Fatal("failed to init DB connection")
	}
	defer dbpool.Close()

	log.Info("Creating tables...")
	err = storage.InitTables(ctx, dbpool)
	if err != nil {
		log.WithError(err).Fatal("failed to init DB tables")
	}
	log.Info("Database is ready")

	h := new(handler)
	h.dbpool = dbpool
	h.ctx = ctx

	r := mux.NewRouter()
	r.Path("/metrics").Handler(promhttp.Handler())
	r.Path("/").HandlerFunc(h.startPage)

	srv := &http.Server{Addr: "0.0.0.0:9000", Handler: r}
	srv.ListenAndServe()
	log.Info("App is running!")
}

type handler struct {
	ctx    context.Context
	dbpool *pgxpool.Pool
}

func (h handler) startPage(w http.ResponseWriter, r *http.Request) {
	motivation, err := motivations.GetRandomMotivation(h.ctx, h.dbpool)
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
