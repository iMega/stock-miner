package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/containerd/containerd/log"
	"github.com/gorilla/websocket"
	"github.com/imega/daemon"
	"github.com/imega/daemon/configuring/env"
	httpserver "github.com/imega/daemon/http-server"
	"github.com/imega/daemon/logging"
	logenv "github.com/imega/daemon/logging/env"
	"github.com/imega/stock-miner/broker"
	"github.com/imega/stock-miner/graph"
	"github.com/imega/stock-miner/graph/generated"
	health_http "github.com/imega/stock-miner/health/http"
	"github.com/imega/stock-miner/httpwareclient"
	"github.com/imega/stock-miner/market"
	"github.com/imega/stock-miner/session"
	"github.com/imega/stock-miner/storage"
	"github.com/imega/stock-miner/teacher"
	"github.com/imega/stock-miner/yahooprovider"
	"github.com/improbable-eng/go-httpwares"
	http_logrus "github.com/improbable-eng/go-httpwares/logging/logrus"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

const (
	shutdownTimeout = 15 * time.Second
	dbFilename      = "./data.db"
	bufferSize      = 1024
	pingInterval    = 15
	lruSize         = 1000
	lruSmallSize    = 100
)

// nolint
var isDevMode = "false"

func main() {
	logger := logging.New(logenv.ReadConfig())

	httpwareclient.WithLogger(logger.(*logrus.Entry))

	if err := prepareDatabase(); err != nil {
		logger.Fatalf("failed to prepare database, %s", err)
	}

	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		logger.Errorf("failed to open db-file, %s", err)
	}

	s := storage.New(storage.WithSqllite(db))

	mux := http.NewServeMux()

	teacher.New(mux)

	rootEmail, _ := env.Read("ROOT_EMAIL")
	clientID, _ := env.Read("GOOGLE_CLIENTID")
	clientSecret, _ := env.Read("GOOGLE_CLIENT_SECRET")
	callbackURL, _ := env.Read("GOOGLE_CALLBACK_URL")
	session := session.New(
		session.WithClintID(clientID),
		session.WithClientSecret(clientSecret),
		session.WithCallbackURL(callbackURL),
		session.WithDevMode(isDevMode),
		session.WithUserStorage(s),
		session.WithRootEmail(rootEmail),
	)
	session.AppendHandlers(mux)

	mux.Handle(
		"/",
		loggerToContext(
			logger,
			session.DefenceHandler(fileServerWithCustom404(Assets)),
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

	h := httpserver.New(
		"stock-miner",
		httpserver.WithLogger(logger),
		httpserver.WithHandler(mux),
		httpserver.WithLogrusOptions(http_logrus.WithDecider(
			func(w httpwares.WrappedResponseWriter, r *http.Request) bool {
				return r.URL.Path != "/healthcheck"
			},
		)),
	)
	cr := env.Once(h.WatcherConfigFunc)

	d, err := daemon.New(logger, cr)
	if err != nil {
		logger.Fatal(err)
	}

	// Market
	marketURL, _ := env.Read("MARKET_TINKOFF_URL")
	marketToken, _ := env.Read("MARKET_TINKOFF_TOKEN")
	marketInstance := market.New(marketURL, marketToken)

	// SMAStack
	smastack := broker.NewSMAStack()

	// Broker
	yfURL, _ := env.Read("YAHOO_FINANCE_URL")
	b := broker.New(
		broker.WithLogger(logger),
		broker.WithStockStorage(s),
		broker.WithPricer(yahooprovider.New(yfURL)),
		broker.WithMarket(marketInstance),
		broker.WithSettingsStorage(s),
		broker.WithStack(s),
		broker.WithSMAStack(smastack),
		broker.WithSetDevMode(isDevMode),
	)

	// handler.WebsocketUpgrader(websocket.Upgrader{
	//     CheckOrigin: func(r *http.Request) bool {
	//       return true
	//     },
	//   }),

	// srv := handler.NewDefaultServer(
	srv := handler.New(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					UserStorage:      s,
					StockStorage:     s,
					MainerController: b,
					Market:           marketInstance,
					SettingsStorage:  s,
					SMAStack:         smastack,
				},
			},
		),
	)
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  bufferSize,
			WriteBufferSize: bufferSize,
		},
		KeepAlivePingInterval: pingInterval,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New(lruSize))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(lruSmallSize),
	})

	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}
	// _ = corsOptions

	mux.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	mux.Handle(
		"/query",
		loggerToContext(
			logger,
			cors.New(corsOptions).Handler(session.DefenceHandler(srv)),
			// session.DefenceHandler(srv),
		),
	)

	d.RegisterShutdownFunc(
		b.ShutdownFunc(),
		func() { db.Close() },
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

func fileServerWithCustom404(fs http.FileSystem) http.Handler {
	next := http.FileServer(fs)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			r.URL.Path = "/index.htm"
		}

		next.ServeHTTP(w, r)
	})
}

func prepareDatabase() error {
	if err := storage.CreateDatabase(dbFilename); err != nil {
		return fmt.Errorf("failed to create database, %w", err)
	}

	if err := storage.MigrateDatabase(dbFilename); err != nil {
		return fmt.Errorf("failed to migrate database, %w", err)
	}

	return nil
}
