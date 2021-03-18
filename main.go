package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/imega/daemon"
	"github.com/imega/daemon/configuring/env"
	httpserver "github.com/imega/daemon/http-server"
	"github.com/imega/daemon/logging"
	"github.com/imega/stock-miner/broker"
	"github.com/imega/stock-miner/storage"
)

const shutdownTimeout = 15 * time.Second
const dbFilename = "./data.db"

func main() {
	log := logging.New(logging.Config{
		Channel: "stock-miner",
		Level:   "debug",
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	h := httpserver.New("stock-miner", log, mux)
	cr := env.Once(h.WatcherConfigFunc)

	d, err := daemon.New(log, cr)
	if err != nil {
		log.Fatal(err)
	}

	if err := storage.CreateDatabase(dbFilename); err != nil {
		log.Fatalf("failed to create database, ", err)
	}

	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		log.Errorf("failed to open db-file, %s", err)
	}

	s := storage.New(storage.WithSqllite(db))
	b := broker.New(broker.WithStorage(s))

	d.RegisterShutdownFunc(b.ShutdownFunc())

	log.Info("stock-miner is started")

	if err := d.Run(shutdownTimeout); err != nil {
		log.Errorf("failed to loop until shutdown: %s", err)
	}

	log.Info("stock-miner is stopped")
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok."))
}
