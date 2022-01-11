package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"html/template"
	"log"
	"net/http"

	"github.com/alextonkonogov/gb-go-observability/homework1/internal/config"
	"github.com/alextonkonogov/gb-go-observability/homework1/internal/motivations"
	"github.com/alextonkonogov/gb-go-observability/homework1/internal/storage"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cnfg, err := config.NewAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	dbpool, err := storage.InitDBConn(ctx, cnfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	h := new(handler)
	h.dbpool = dbpool
	h.ctx = ctx

	r := mux.NewRouter()
	r.Path("/metrics").Handler(promhttp.Handler())
	r.Path("/").HandlerFunc(h.startPage)

	srv := &http.Server{Addr: "localhost:1234", Handler: r}
	srv.ListenAndServe()
}

type handler struct {
	ctx    context.Context
	dbpool *pgxpool.Pool
}

func (h handler) startPage(w http.ResponseWriter, r *http.Request) {
	motivation, err := motivations.GetRandomMotivation(h.ctx, h.dbpool)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	content := fmt.Sprintf("\"%s\" %s\n", motivation.Content, motivation.Author)
	tmpl, err := template.New("example").Parse(content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	tmpl.Execute(w, content)
}
