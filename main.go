package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/containerd/containerd/log"
	"github.com/imega/daemon"
	"github.com/imega/daemon/configuring/env"
	httpserver "github.com/imega/daemon/http-server"
	"github.com/imega/daemon/logging"
	"github.com/imega/stock-miner/broker"
	"github.com/imega/stock-miner/graph"
	"github.com/imega/stock-miner/graph/generated"
	health_http "github.com/imega/stock-miner/health/http"
	"github.com/imega/stock-miner/session"
	"github.com/imega/stock-miner/storage"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

const shutdownTimeout = 15 * time.Second
const dbFilename = "./data.db"

var devMode = "false"

func main() {
	logger := logging.New(logging.Config{
		Channel: "stock-miner",
		Level:   "debug",
	})

	if err := storage.CreateDatabase(dbFilename); err != nil {
		logger.Fatalf("failed to create database, ", err)
	}

	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		logger.Errorf("failed to open db-file, %s", err)
	}

	s := storage.New(storage.WithSqllite(db))

	mux := http.NewServeMux()

	clientID, _ := env.Read("GOOGLE_CLIENTID")
	clientSecret, _ := env.Read("GOOGLE_CLIENT_SECRET")
	callbackURL, _ := env.Read("GOOGLE_CALLBACK_URL")
	session := session.New(
		session.WithClintID(clientID),
		session.WithClientSecret(clientSecret),
		session.WithCallbackURL(callbackURL),
		session.WithDevMode(devMode),
		session.WithUserStorage(s),
	)
	session.AppendHandlers(mux)

	mux.Handle(
		"/",
		loggerToContext(
			logger,
			session.DefenceHandler(http.FileServer(Assets)),
		),
	)
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

	h := httpserver.New("stock-miner", logger, mux)
	cr := env.Once(h.WatcherConfigFunc)

	d, err := daemon.New(logger, cr)
	if err != nil {
		logger.Fatal(err)
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					UserStorage: s,
				}},
		),
	)

	mux.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", loggerToContext(logger, session.DefenceHandler(srv)))

	b := broker.New(broker.WithStorage(s), broker.WithLogger(logger))

	d.RegisterShutdownFunc(
		b.ShutdownFunc(),
		func() {
			db.Close()
		},
	)

	logger.Info("stock-miner is started")

	if err := d.Run(shutdownTimeout); err != nil {
		logger.Errorf("failed to loop until shutdown: %s", err)
	}

	logger.Info("stock-miner is stopped")
}

func loggerToContext(l logrus.FieldLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := log.WithLogger(req.Context(), l.(*logrus.Entry))

		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
