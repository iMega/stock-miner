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
	health_http "github.com/imega/stock-miner/health/http"
	"github.com/imega/stock-miner/session_store"
	"github.com/imega/stock-miner/storage"
)

const shutdownTimeout = 15 * time.Second
const dbFilename = "./data.db"

func main() {
	log := logging.New(logging.Config{
		Channel: "stock-miner",
		Level:   "debug",
	})

	if err := storage.CreateDatabase(dbFilename); err != nil {
		log.Fatalf("failed to create database, ", err)
	}

	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		log.Errorf("failed to open db-file, %s", err)
	}

	mux := http.NewServeMux()

	clientID, _ := env.Read("GOOGLE_CLIENTID")
	clientSecret, _ := env.Read("GOOGLE_CLIENT_SECRET")
	callbackURL, _ := env.Read("GOOGLE_CALLBACK_URL")
	session := session_store.New(
		session_store.WithClintID(clientID),
		session_store.WithClientSecret(clientSecret),
		session_store.WithCallbackURL(callbackURL),
	)
	session.AppendHandlers(mux)

	mux.Handle("/", session.DefenceHandler(http.FileServer(Assets)))
	mux.HandleFunc(
		"/healthcheck",
		health_http.HandlerFunc(
			health_http.WithHealthCheckFuncs(
				func() bool {
					if err := db.Ping(); err != nil {
						return false
					}

					return true
				},
			),
		),
	)

	h := httpserver.New("stock-miner", log, mux)
	cr := env.Once(h.WatcherConfigFunc)

	d, err := daemon.New(log, cr)
	if err != nil {
		log.Fatal(err)
	}

	s := storage.New(storage.WithSqllite(db))
	b := broker.New(broker.WithStorage(s), broker.WithLogger(log))

	d.RegisterShutdownFunc(
		b.ShutdownFunc(),
		func() {
			db.Close()
		},
	)

	log.Info("stock-miner is started")

	if err := d.Run(shutdownTimeout); err != nil {
		log.Errorf("failed to loop until shutdown: %s", err)
	}

	log.Info("stock-miner is stopped")
}
